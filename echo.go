package simutils

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/alifakhimi/simple-utils-go/multierror"
	"github.com/labstack/echo/v4"
)

// Define custom errors for specific scenarios
var (
	// Error when casting echo.Context to custom Context fails
	ErrBindContextFailed = errors.New("binding context failed")
	// Error when type assertion of RequestModel fails
	ErrRequestModelTypeAssertion = errors.New("request model type assertion failed")
)

// Context is a custom structure that embeds echo.Context and adds a RequestModel field
type Context struct {
	echo.Context
	// RequestModel can hold any type of data
	RequestModel any
}

// Binder attempts to bind a given object to the custom Context
func Binder(echoContext echo.Context, i any) (*Context, error) {
	// Check if echoContext can be cast to *Context
	if ctx, ok := echoContext.(*Context); !ok {
		return nil, ErrBindContextFailed // Return error if casting fails
	} else if err := ctx.Bind(i); err != nil { // Attempt to bind the object to the context
		return nil, err // Return error if binding fails
	} else {
		ctx.RequestModel = i // Assign the bound object to RequestModel
		return ctx, nil      // Return the updated context
	}
}

// GetRequestModel extracts the RequestModel from Context and performs a type assertion
func GetRequestModel[T any](ctx *Context) (T, error) {
	// Check if RequestModel is not nil
	if ctx.RequestModel != nil {
		// Attempt to cast RequestModel to the desired type T
		if value, ok := ctx.RequestModel.(T); ok {
			return value, nil // Return the casted value if successful
		} else {
			return value, ErrRequestModelTypeAssertion // Return error if type assertion fails
		}
	} else {
		// If RequestModel is nil, return the zero value of type T and an error
		var zv T
		return zv, errors.New("type assertion failed")
	}
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
		template = ResponseOk(content, multierror.Join(err), meta)
	case http.StatusCreated:
		template = ResponseCreated(content, meta)
	case http.StatusBadRequest:
		template = ResponseBadRequest(content, multierror.Join(err))
	case http.StatusInternalServerError:
		template = ResponseInternalServerError(content, multierror.Join(err))
	case http.StatusNotFound:
		template = ResponseNotFound(content, multierror.Join(err))
	case http.StatusUnprocessableEntity:
		template = ResponseUnprocessableEntity(content, multierror.Join(err))
	case http.StatusMethodNotAllowed:
		template = ResponseMethodNotAllowed(content, multierror.Join(err))
	case http.StatusUnauthorized:
		template = ResponseUnauthorized(content, multierror.Join(err))
	case http.StatusForbidden:
		template = ResponseForbidden(content, multierror.Join(err))
	case http.StatusGatewayTimeout:
		template = ResponseGatewayTimeOut(content, multierror.Join(err))
	case http.StatusLocked:
		template = ResponseLocked(content, multierror.Join(err))
	case http.StatusNotAcceptable:
		template = ResponseNotAcceptable(content, multierror.Join(err))
	default:
		template = ResponseInternalServerError(content, multierror.Join(err))
	}

	return ctx.JSON(httpStatus, template)
}

// GetWithCode return template with considering error code
func GetWithCode(data interface{}, code int, err error) (template *ResponseTemplate) {
	msg := err.Error()

	switch code {
	case http.StatusInternalServerError:
		template = ResponseInternalServerError(data, msg)
	case http.StatusBadRequest:
		template = ResponseBadRequest(data, msg)
	case http.StatusForbidden:
		template = ResponseForbidden(data, msg)
	case http.StatusNotFound:
		template = ResponseNotFound(data, msg)
	case http.StatusUnprocessableEntity:
		template = ResponseUnprocessableEntity(data, msg)
	case http.StatusUnauthorized:
		template = ResponseUnauthorized(data, msg)
	case http.StatusMethodNotAllowed:
		template = ResponseMethodNotAllowed(data, msg)
	default:
		template = ResponseStatusNotImplemented(data, msg)
	}

	return
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

func PopQueryParam[T string | []string](qp url.Values, key string) T {
	var result T
	switch any(result).(type) {
	case string:
		result = any(qp.Get(key)).(T)
	case []string:
		result = any(qp[key]).(T)
	}
	qp.Del(key)
	return result
}
