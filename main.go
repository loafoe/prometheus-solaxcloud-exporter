package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/loafoe/prometheus-solaxcloud-exporter/solaxcloud"
)

var listenAddr string
var debugLog string

func main() {
	flag.StringVar(&debugLog, "debuglog", "", "The debug log to dump traffic in")
	flag.StringVar(&listenAddr, "listen", "0.0.0.0:8889", "Listen address for HTTP metrics")
	flag.Parse()

	sn := os.Getenv("SOLAXCLOUD_SN")
	tokenId := os.Getenv("SOLAXCLOUD_TOKEN_ID")

	resp, err := solaxcloud.GetRealtimeInfo(solaxcloud.WithSNAndTokenID(sn, tokenId))
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("%+v\n", resp)
}
