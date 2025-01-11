package simutils

import (
	"encoding/json"
	"math"
	"net/http"
)

// ResponseTemplate standard template for http responses
type ResponseTemplate struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
	Meta    interface{} `json:"meta,omitempty"`
	Links   interface{} `json:"links,omitempty"`
}

func (m ResponseTemplate) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *ResponseTemplate) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

// PaginateTemplate ...
type PaginateTemplate struct {
	Pages    int     `json:"pages"`
	Total    int     `json:"total"`
	Limit    int     `json:"limit"`
	Offset   int     `json:"offset"`
	Page     int     `json:"page"`
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"prev"`
}

// CreatePaginateTemplateByPage create pagination using page
func CreatePaginateTemplateByPage(total, page, limit int) *PaginateTemplate {
	return PaginateTemplate{}.Create(total, page*limit, limit)
}

// CreatePaginateTemplate create pagination using offset
func CreatePaginateTemplate(total, offset, limit int) *PaginateTemplate {
	return PaginateTemplate{}.Create(total, offset, limit)
}

// Create create pagination
func (PaginateTemplate) Create(total, offset, limit int) *PaginateTemplate {
	var (
		pages    int
		next     *string
		prev     *string
		tempNext string
		tempPrev string

		pt = &PaginateTemplate{}
	)

	if offset <= 0 {
		offset = 0
	}

	if limit <= 0 {
		limit = 5
	}

	pages = int(math.Ceil(float64(total) / float64(limit)))

	if offset < total-limit {
		tempNext = "has next"
		next = &tempNext
	}

	if offset > 0 && total > limit {
		tempPrev = "has prev"
		prev = &tempPrev
	}

	pt.Next = next
	pt.Limit = limit
	pt.Offset = offset
	pt.Pages = pages
	pt.Page = offset/limit + 1
	pt.Previous = prev
	pt.Total = total

	return pt
}

// BadRequest ...
func ResponseBadRequest(data, msg interface{}) *ResponseTemplate {
	return &ResponseTemplate{
		Code:    http.StatusBadRequest,
		Status:  http.StatusText(http.StatusBadRequest),
		Message: msg,
		Data:    data,
	}
}

// InternalServerError ...
func ResponseInternalServerError(data, msg interface{}) *ResponseTemplate {
	return &ResponseTemplate{
		Code:    http.StatusInternalServerError,
		Status:  http.StatusText(http.StatusInternalServerError),
		Message: msg,
		Data:    data,
	}
}

// StatusNotImplemented ...
func ResponseStatusNotImplemented(data, msg interface{}) *ResponseTemplate {
	return &ResponseTemplate{
		Code:    http.StatusNotImplemented,
		Status:  http.StatusText(http.StatusNotImplemented),
		Message: msg,
		Data:    data,
	}
}

// InternalServerErrorWithData ...
func ResponseInternalServerErrorWithData(data, msg interface{}) *ResponseTemplate {
	return &ResponseTemplate{
		Code:    http.StatusInternalServerError,
		Status:  http.StatusText(http.StatusInternalServerError),
		Message: msg,
		Data:    data,
	}
}

// NotFound ...
func ResponseNotFound(data, msg interface{}) *ResponseTemplate {
	return &ResponseTemplate{
		Code:    http.StatusNotFound,
		Status:  http.StatusText(http.StatusNotFound),
		Message: msg,
		Data:    data,
	}
}

// UnprocessableEntity ...
func ResponseUnprocessableEntity(data, msg interface{}) *ResponseTemplate {
	return &ResponseTemplate{
		Code:    http.StatusUnprocessableEntity,
		Status:  http.StatusText(http.StatusUnprocessableEntity),
		Message: msg,
		Data:    data,
	}
}

// Unauthorized ...
func ResponseUnauthorized(data, msg interface{}) *ResponseTemplate {
	return &ResponseTemplate{
		Code:    http.StatusUnauthorized,
		Status:  http.StatusText(http.StatusUnauthorized),
		Message: msg,
		Data:    data,
	}
}

// GatewayTimeOut ...
func ResponseGatewayTimeOut(data, msg interface{}) *ResponseTemplate {
	return &ResponseTemplate{
		Code:    http.StatusGatewayTimeout,
		Status:  http.StatusText(http.StatusGatewayTimeout),
		Message: msg,
		Data:    data,
	}
}

// Locked ...
func ResponseLocked(data, msg interface{}) *ResponseTemplate {
	return &ResponseTemplate{
		Code:    http.StatusLocked,
		Status:  http.StatusText(http.StatusLocked),
		Message: msg,
		Data:    data,
	}
}

// NotAcceptable ...
func ResponseNotAcceptable(data, msg interface{}) *ResponseTemplate {
	return &ResponseTemplate{
		Code:    http.StatusNotAcceptable,
		Status:  http.StatusText(http.StatusNotAcceptable),
		Message: msg,
		Data:    data,
	}
}

// Forbidden ...
func ResponseForbidden(data, msg interface{}) *ResponseTemplate {
	return &ResponseTemplate{
		Code:    http.StatusForbidden,
		Status:  http.StatusText(http.StatusForbidden),
		Message: msg,
		Data:    data,
	}
}

// MethodNotAllowed ...
func ResponseMethodNotAllowed(data, msg interface{}) *ResponseTemplate {
	return &ResponseTemplate{
		Code:    http.StatusMethodNotAllowed,
		Status:  http.StatusText(http.StatusMethodNotAllowed),
		Message: msg,
		Data:    data,
	}
}

// Ok ...
func ResponseOk(data, msg interface{}, meta interface{}) *ResponseTemplate {
	return &ResponseTemplate{
		Code:    http.StatusOK,
		Status:  http.StatusText(http.StatusOK),
		Message: msg,
		Data:    data,
		Meta:    meta,
	}
}

// Created ...
func ResponseCreated(data, meta interface{}) *ResponseTemplate {
	return &ResponseTemplate{
		Code:    http.StatusCreated,
		Status:  http.StatusText(http.StatusCreated),
		Message: "Created",
		Data:    data,
		Meta:    meta,
	}
}
