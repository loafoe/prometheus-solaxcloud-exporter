package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/loafoe/prometheus-solaxcloud-exporter/solaxcloud"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var listenAddr string
var debug bool

var (
	metricNamePrefix = "solaxcloud_"
	registry         = prometheus.NewRegistry()
)

var (
	yieldTodayMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricNamePrefix + "yield_today",
		Help: "The yield for today (KWh)",
	}, []string{
		"inverter_sn",
	})

	yieldTotalMetrics = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricNamePrefix + "yield_total",
		Help: "The total yield of the system (KWh)",
	}, []string{
		"inverter_sn",
	})

	acPowerMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricNamePrefix + "ac_power",
		Help: "Current power generation (Wh)",
	}, []string{
		"inverter_sn",
	})
	upMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricNamePrefix + "up",
		Help: "The inverter power on status",
	}, []string{
		"sn",
		"inverter_sn",
	})
)

func init() {
	registry.MustRegister(yieldTotalMetrics)
	registry.MustRegister(yieldTodayMetric)
	registry.MustRegister(acPowerMetric)
	registry.MustRegister(upMetric)
}

func main() {
	flag.BoolVar(&debug, "debug", false, "Enable debugging")
	flag.StringVar(&listenAddr, "listen", "0.0.0.0:8887", "Listen address for HTTP metrics")
	flag.Parse()

	sn := os.Getenv("SOLAXCLOUD_SN")
	tokenId := os.Getenv("SOLAXCLOUD_TOKEN_ID")

	go func() {
		sleep := false
		for {
			if sleep {
				time.Sleep(time.Second * 60) // 5 minute resolution, so we poll every minute for now
			}
			sleep = true
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

			fmt.Printf("calling SolaxCloud...\n")
			resp, err := solaxcloud.GetRealtimeInfo(ctx,
				solaxcloud.WithSNAndTokenID(sn, tokenId),
				solaxcloud.WithDebug(debug))
			cancel()
			if err != nil {
				fmt.Printf("error: %v\n", err)
				upMetric.WithLabelValues(sn, "").Set(0)
				if errors.Is(err, context.DeadlineExceeded) {
					fmt.Printf("not sleeping\n")
					sleep = false
				}
				continue
			}
			yieldTodayMetric.WithLabelValues(resp.Result.InverterSN).Set(resp.Result.YieldToday)
			yieldTotalMetrics.WithLabelValues(resp.Result.InverterSN).Set(resp.Result.YieldTotal)
			acPowerMetric.WithLabelValues(resp.Result.InverterSN).Set(resp.Result.ACPower)
			upMetric.WithLabelValues(sn, resp.Result.InverterSN).Set(1.0)
		}
	}()

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	_ = http.ListenAndServe(listenAddr, nil)
}
