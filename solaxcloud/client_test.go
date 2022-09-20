package solaxcloud_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/loafoe/prometheus-solaxcloud-exporter/solaxcloud"
	"github.com/stretchr/testify/assert"
)

var (
	muxSolaxCloud    *http.ServeMux
	serverSolaxCloud *httptest.Server
)

func setup(_ *testing.T) func() {
	muxSolaxCloud = http.NewServeMux()
	serverSolaxCloud = httptest.NewServer(muxSolaxCloud)

	return func() {
		serverSolaxCloud.Close()
	}
}

func getURL(host string) string {
	return host + "/proxyApp/proxy/api/getRealtimeInfo.do"
}

func TestGetRealtimeInfo(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	muxSolaxCloud.HandleFunc("/proxyApp/proxy/api/getRealtimeInfo.do", func(w http.ResponseWriter, r *http.Request) {
		if !assert.Equal(t, http.MethodGet, r.Method) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
  "success": true,
  "exception": "Query success!",
  "result": {
    "inverterSN": "XM0000000000",
    "sn": "S123456789",
    "acpower": 163,
    "yieldtoday": 3.5,
    "yieldtotal": 5653.8,
    "feedinpower": 0,
    "feedinenergy": 0,
    "consumeenergy": 0,
    "feedinpowerM2": 0,
    "soc": 0,
    "peps1": 0,
    "peps2": 0,
    "peps3": 0,
    "inverterType": "4",
    "inverterStatus": "102",
    "uploadTime": "2022-09-19 17:42:20",
    "batPower": 0,
    "powerdc1": 173,
    "powerdc2": 0,
    "powerdc3": null,
    "powerdc4": null,
    "batStatus": "0"
  }
}`)
	})

	ctx := context.Background()

	resp, err := solaxcloud.GetRealtimeInfo(ctx, solaxcloud.WithURL(getURL(serverSolaxCloud.URL)))
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.True(t, resp.Success)
	assert.Equal(t, "XM0000000000", resp.Result.InverterSN)
}

func TestGetRealtimeInfoError(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	muxSolaxCloud.HandleFunc("/proxyApp/proxy/api/getRealtimeInfo.do", func(w http.ResponseWriter, r *http.Request) {
		if !assert.Equal(t, http.MethodGet, r.Method) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
  "success": false,
  "exception": "Query success!",
  "result": "this sn did not access!"
}`)
	})

	ctx := context.Background()

	resp, err := solaxcloud.GetRealtimeInfo(ctx, solaxcloud.WithURL(getURL(serverSolaxCloud.URL)))
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Equal(t, "this sn did not access!", resp.Result.Error)
	assert.False(t, resp.Success)
}
