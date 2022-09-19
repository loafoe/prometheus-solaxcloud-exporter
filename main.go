package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/loafoe/prometheus-solaxcloud-exporter/solaxcloud"
)

var listenAddr string
var debug bool

func main() {
	flag.BoolVar(&debug, "debug", false, "Enable debugging")
	flag.StringVar(&listenAddr, "listen", "0.0.0.0:8887", "Listen address for HTTP metrics")
	flag.Parse()

	sn := os.Getenv("SOLAXCLOUD_SN")
	tokenId := os.Getenv("SOLAXCLOUD_TOKEN_ID")

	resp, err := solaxcloud.GetRealtimeInfo(
		solaxcloud.WithSNAndTokenID(sn, tokenId),
		solaxcloud.WithDebug(debug))
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("%+v\n", resp)
}
