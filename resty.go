package simutils

import (
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"

	"github.com/alifakhimi/simple-service-go/utils/rest"
)

type RestyClientInfo struct {
	BaseURL string
	Token   string
}

func NewRestyClient(info *RestyClientInfo) (client *resty.Client, err error) {
	if info == nil {
		return nil, rest.ErrInvalidRestClientInfo
	}

	if _, err = url.Parse(info.BaseURL); err != nil {
		return nil, err
	}

	client = resty.
		New().
		SetBaseURL(info.BaseURL).
		SetAuthToken(info.Token)

	return client, nil
}

func LogResty(reqiest *resty.Request, response *resty.Response, err error) {
	// Explore request object
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Println("Request Info:")
	fmt.Println("  Method	:", reqiest.Method)
	fmt.Println("  Url		:", reqiest.URL)
	fmt.Println("  Body		:\n", reqiest.Body)

	// Explore response object
	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", response.StatusCode())
	fmt.Println("  Status     :", response.Status())
	fmt.Println("  Proto      :", response.Proto())
	fmt.Println("  Time       :", response.Time())
	fmt.Println("  Received At:", response.ReceivedAt())
	fmt.Println("  Body       :\n", response)
	fmt.Println()

	// //Explore trace info
	// logrus.Infoln("Request Trace Info:")
	// ti := response.Request.TraceInfo()
	// fmt.Println("  DNSLookup     :", ti.DNSLookup)
	// fmt.Println("  ConnTime      :", ti.ConnTime)
	// fmt.Println("  TCPConnTime   :", ti.TCPConnTime)
	// fmt.Println("  TLSHandshake  :", ti.TLSHandshake)
	// fmt.Println("  ServerTime    :", ti.ServerTime)
	// fmt.Println("  ResponseTime  :", ti.ResponseTime)
	// fmt.Println("  TotalTime     :", ti.TotalTime)
	// fmt.Println("  IsConnReused  :", ti.IsConnReused)
	// fmt.Println("  IsConnWasIdle :", ti.IsConnWasIdle)
	// fmt.Println("  ConnIdleTime  :", ti.ConnIdleTime)
	// fmt.Println("  RequestAttempt:", ti.RequestAttempt)
	// fmt.Println("  RemoteAddr    :", ti.RemoteAddr.String())
}
