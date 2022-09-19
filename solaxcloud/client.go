package solaxcloud

import (
	"encoding/json"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type Result result

type result struct {
	Error          string  `json:"-"`
	InverterSN     string  `json:"inverterSN"`
	SN             string  `json:"sn"`
	ACPower        int     `json:"acpower"`
	YieldToday     float64 `json:"yieldtoday"`
	YieldTotal     float64 `json:"yieldtotal"`
	FeedInPower    int     `json:"feedinpower"`
	FeedInEnergy   int     `json:"feedinenergy"`
	ConsumeEnergy  int     `json:"consumeenergy"`
	FeedInPowerM2  int     `json:"feedinpowerM2"`
	Soc            int     `json:"soc"`
	Peps1          int     `json:"peps1"`
	Peps2          int     `json:"peps2"`
	Peps3          int     `json:"peps3"`
	InverterType   string  `json:"inverterType"`
	InverterStatus string  `json:"inverterStatus"`
	UploadTime     string  `json:"uploadTime"`
	BatPower       int     `json:"batPower"`
	PowerDC1       int     `json:"powerdc1"`
	PowerDC2       int     `json:"powerdc2"`
	PowerDC3       *int    `json:"powerdc3,omitempty"`
	PowerDC4       *int    `json:"powerdc4,omitempty"`
	BatStatus      string  `json:"batStatus"`
}

type OptionFunc func(r *resty.Request) (*resty.Request, error)

func (r *Result) UnmarshalJSON(data []byte) (err error) {
	var str string
	var res result
	err = json.Unmarshal(data, &str)
	if err == nil {
		r.Error = str
		return nil
	}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return err
	}
	r.BatPower = res.BatPower
	r.InverterType = res.InverterType
	r.SN = res.SN
	r.BatStatus = res.BatStatus
	r.ACPower = res.ACPower
	r.ConsumeEnergy = res.ConsumeEnergy
	r.InverterSN = res.InverterSN
	r.FeedInEnergy = res.FeedInEnergy
	r.FeedInPower = res.FeedInPower
	r.FeedInPowerM2 = res.FeedInPowerM2
	r.Peps1 = res.Peps1
	r.Peps2 = res.Peps2
	r.Peps3 = res.Peps3
	r.PowerDC1 = res.PowerDC1
	r.PowerDC2 = res.PowerDC2
	r.PowerDC3 = res.PowerDC3
	r.PowerDC4 = res.PowerDC4
	r.Soc = res.Soc
	r.YieldTotal = res.YieldTotal
	r.YieldToday = res.YieldToday
	r.InverterStatus = res.InverterStatus
	r.UploadTime = res.UploadTime
	return nil
}

type Response struct {
	Success   bool   `json:"success"`
	Exception string `json:"exception"`
	Result    Result `json:"result"`
}

func WithDefaultURL() OptionFunc {
	return func(r *resty.Request) (*resty.Request, error) {
		r.URL = "https://www.solaxcloud.com/proxyApp/proxy/api/getRealtimeInfo.do"
		return r, nil
	}
}

func WithURL(url string) OptionFunc {
	return func(r *resty.Request) (*resty.Request, error) {
		r.URL = url
		return r, nil
	}
}

func WithSNAndTokenID(sn, tokenID string) OptionFunc {
	return func(r *resty.Request) (*resty.Request, error) {
		return r.SetQueryParams(map[string]string{
			"sn":      sn,
			"tokenId": tokenID,
		}), nil
	}
}

func GetRealtimeInfo(opts ...OptionFunc) (*Response, error) {
	client := resty.New()

	request := client.R()
	request.Method = http.MethodGet
	request, _ = WithDefaultURL()(request)

	for _, o := range opts {
		r, err := o(request)
		if err != nil {
			return nil, err
		}
		request = r
	}
	resp, err := request.Send()
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
