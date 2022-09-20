package solaxcloud

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-resty/resty/v2"
)

func GetRealtimeInfo(ctx context.Context, opts ...OptionFunc) (*Response, error) {
	client := resty.New()

	request := client.R()
	request.Method = http.MethodGet
	request, _ = WithDefaultURL()(client, request)

	for _, o := range opts {
		r, err := o(client, request)
		if err != nil {
			return nil, err
		}
		request = r
	}
	resp, err := request.SetContext(ctx).Send()
	if err != nil {
		return nil, err
	}
	var jsonResponse Response
	err = json.Unmarshal(resp.Body(), &jsonResponse)
	if err != nil {
		return nil, err
	}
	return &jsonResponse, nil
}
