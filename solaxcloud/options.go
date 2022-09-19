package solaxcloud

import (
	"github.com/go-resty/resty/v2"
)

type OptionFunc func(c *resty.Client, r *resty.Request) (*resty.Request, error)

func WithDefaultURL() OptionFunc {
	return func(c *resty.Client, r *resty.Request) (*resty.Request, error) {
		r.URL = "https://www.solaxcloud.com/proxyApp/proxy/api/getRealtimeInfo.do"
		return r, nil
	}
}

func WithURL(url string) OptionFunc {
	return func(c *resty.Client, r *resty.Request) (*resty.Request, error) {
		r.URL = url
		return r, nil
	}
}

func WithSNAndTokenID(sn, tokenID string) OptionFunc {
	return func(c *resty.Client, r *resty.Request) (*resty.Request, error) {
		return r.SetQueryParams(map[string]string{
			"sn":      sn,
			"tokenId": tokenID,
		}), nil
	}
}

func WithDebug(debug bool) OptionFunc {
	return func(c *resty.Client, r *resty.Request) (*resty.Request, error) {
		c.SetDebug(debug)
		return r, nil
	}
}
