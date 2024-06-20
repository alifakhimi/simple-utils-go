package rest

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type (
	Clients map[string]*Client

	Client struct {
		ClientConfig
		Client *resty.Client `json:"-" gorm:"-:all"`
	}

	ClientConfig struct {
		BaseURL               string            `json:"base_url,omitempty" mapstructure:"base_url"`
		QueryParam            url.Values        `json:"query_param,omitempty" mapstructure:"query_param"`
		FormData              url.Values        `json:"form_data,omitempty" mapstructure:"form_data"`
		PathParams            map[string]string `json:"path_params,omitempty" mapstructure:"path_params"`
		Header                http.Header       `json:"header,omitempty" mapstructure:"header"`
		UserInfo              *User             `json:"user_info,omitempty" mapstructure:"user_info"`
		Token                 string            `json:"token,omitempty" mapstructure:"token"`
		AuthScheme            string            `json:"auth_scheme,omitempty" mapstructure:"auth_scheme"`
		Cookies               []*http.Cookie    `json:"cookies,omitempty" mapstructure:"cookies"`
		Debug                 bool              `json:"debug,omitempty" mapstructure:"debug"`
		DisableWarn           bool              `json:"disable_warn,omitempty" mapstructure:"disable_warn"`
		AllowGetMethodPayload bool              `json:"allow_get_method_payload,omitempty" mapstructure:"allow_get_method_payload"`
		RetryCount            int               `json:"retry_count,omitempty" mapstructure:"retry_count"`
		RetryWaitTime         time.Duration     `json:"retry_wait_time,omitempty" mapstructure:"retry_wait_time"`
		RetryMaxWaitTime      time.Duration     `json:"retry_max_wait_time,omitempty" mapstructure:"retry_max_wait_time"`
		// Proxy sets the Proxy URL and Port for Resty client.
		//
		// like: http://proxyserver:8888
		//
		// You could also set Proxy via environment variable.
		//
		// Refer to godoc `http.ProxyFromEnvironment`.
		Proxy    string `json:"proxy,omitempty"`
		UseProxy bool   `json:"use_proxy,omitempty"`
	}
)

// User type is to hold an username and password information
type User struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (r *Client) UnmarshalJSON(b []byte) error {
	var tmp = &Client{}
	if err := json.Unmarshal(b, &tmp.ClientConfig); err != nil {
		return err
	}

	restyClient := resty.New().
		SetBaseURL(tmp.ClientConfig.BaseURL).
		SetPathParams(tmp.ClientConfig.PathParams).
		SetAuthToken(tmp.ClientConfig.Token).
		SetAuthScheme(tmp.ClientConfig.AuthScheme).
		SetCookies(tmp.ClientConfig.Cookies).
		SetDebug(tmp.ClientConfig.Debug).
		SetDisableWarn(tmp.ClientConfig.DisableWarn).
		SetAllowGetMethodPayload(tmp.ClientConfig.AllowGetMethodPayload).
		SetRetryCount(tmp.ClientConfig.RetryCount)

	if tmp.ClientConfig.RetryMaxWaitTime > 0 {
		restyClient = restyClient.
			SetRetryWaitTime(tmp.ClientConfig.RetryWaitTime)
	}
	if tmp.ClientConfig.RetryMaxWaitTime > 0 {
		restyClient = restyClient.
			SetRetryMaxWaitTime(tmp.ClientConfig.RetryMaxWaitTime)
	}
	if tmp.ClientConfig.UseProxy {
		restyClient = restyClient.
			SetProxy(tmp.ClientConfig.Proxy)
	}
	for k, v := range tmp.ClientConfig.QueryParam {
		for _, s := range v {
			restyClient.QueryParam.Add(k, s)
		}
	}
	for k, v := range tmp.ClientConfig.FormData {
		for _, s := range v {
			restyClient.FormData.Add(k, s)
		}
	}
	for k, v := range tmp.ClientConfig.Header {
		for _, s := range v {
			restyClient.Header.Add(k, s)
		}
	}

	tmp.Client = restyClient

	*r = *tmp

	return nil
}

func (c Client) Value() (value driver.Value, err error) {
	var b []byte
	if b, err = json.Marshal(c); err != nil {
		return
	}

	return b, nil
}

func (c *Client) Scan(value interface{}) (err error) {
	if value == nil {
		*c = Client{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := Client{}
	if err = json.Unmarshal(bytes, &result); err != nil {
		return
	}

	*c = result

	return nil
}

// GormDBDataType gorm db data type
func (Client) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}

func (cs Clients) Get(name string) *resty.Client {
	return nil
}
