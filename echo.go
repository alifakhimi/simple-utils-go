package simutils

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/labstack/echo/v4"
)

func IsDigit(s string) bool {
	if strings.TrimSpace(s) == "" {
		return false
	}

	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func ReplyTemplate(ctx echo.Context, httpStatus int, err error, template interface{}, meta interface{}) error {
	data, er := json.Marshal(template) // Convert to a json string
	if er != nil {
		return er
	}

	content := make(map[string]interface{})

	er = json.Unmarshal(data, &content) // Convert to a map
	if err != nil {
		return er
	}

	return Reply(ctx, httpStatus, err, content, meta)
}

// Reply ...
func Reply(ctx echo.Context, httpStatus int, err error, content map[string]interface{}, meta interface{}) error {
	var template *ResponseTemplate

	switch httpStatus {
	case http.StatusOK:
		template = Ok(content, err, meta)
	case http.StatusCreated:
		template = Created(content, meta)
	case http.StatusBadRequest:
		template = BadRequest(content, err.Error())
	case http.StatusInternalServerError:
		template = InternalServerError(content, err.Error())
	case http.StatusNotFound:
		template = NotFound(content, err.Error())
	case http.StatusUnprocessableEntity:
		template = UnprocessableEntity(content, err.Error())
	case http.StatusMethodNotAllowed:
		template = MethodNotAllowed(content, err.Error())
	case http.StatusUnauthorized:
		template = Unauthorized(content, err.Error())
	case http.StatusForbidden:
		template = Forbidden(content, err.Error())
	case http.StatusGatewayTimeout:
		template = GatewayTimeOut(content, err.Error())
	case http.StatusLocked:
		template = Locked(content, err.Error())
	case http.StatusNotAcceptable:
		template = NotAcceptable(content, err.Error())
	default:
		template = InternalServerError(content, errors.New("invalid reply request"))
	}

	return ctx.JSON(httpStatus, template)
}

// ExportRoutes ...
func ExportRoutes(e *echo.Echo, prefix string) ([]*echo.Route, error) {
	var apiRoutes []*echo.Route
	routes := e.Routes()
	for _, route := range routes {
		if strings.Index(route.Path, prefix) == 0 {
			apiRoutes = append(apiRoutes, route)
		}
	}

	// data, err := json.Marshal(apiRoutes)
	// if err != nil {
	// 	return nil, err
	// }
	// //ioutil.WriteFile("routes.json", data, 0644)
	// return data, err
	return apiRoutes, nil
}

func DurationToHumanity(d time.Duration) (days int, hours int, minutes int) {
	minutes = int(d.Minutes()) % 60
	hours = int(d.Hours()) % 24
	days = int(d.Hours()) / 24
	return days, hours, minutes
}

// ParseFilterQuery ...
func ParseFilterQuery(ctx echo.Context, key string) []uint {

	_url := ctx.QueryString()
	qs, _ := url.ParseQuery(_url)
	filterArr := qs[key]

	filters := make([]uint, len(filterArr))
	for i := 0; i < len(filterArr); i++ {
		v, _ := strconv.Atoi(filterArr[i])
		filters[i] = uint(v)
	}

	return filters
}

// ParseIDsQuery ...
func ParseIDsQuery(ctx echo.Context, key string) PIDs {
	_url := ctx.QueryString()
	qs, _ := url.ParseQuery(_url)
	filterArr := qs[key]
	ids := strings.Split(filterArr[0], ",")

	filters := make(PIDs, len(ids))
	for i := 0; i < len(ids); i++ {
		if pid, err := ParsePID(ids[i]); err == nil {
			filters[i] = pid
		}
	}

	return filters
}

func ArrayElementExists(a []string, v string) bool {
	for _, s := range a {
		if s == v {
			return true
		}
	}

	return false
}

func ErrorToHttpStatusCode(err error) (status int) {
	switch err {
	case ErrNotFound, ErrRecordNotFound:
		status = http.StatusNotFound
	case ErrInvalidRequest:
		status = http.StatusBadRequest
	case ErrAlreadyExist:
		status = http.StatusNotAcceptable

	default:
		status = http.StatusNotImplemented
	}

	return
}
