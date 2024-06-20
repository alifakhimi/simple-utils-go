package simutils

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
)

func CloneURL(u *url.URL) *url.URL {
	if u == nil {
		return nil
	}
	u2 := new(url.URL)
	*u2 = *u
	if u.User != nil {
		u2.User = new(url.Userinfo)
		*u2.User = *u.User
	}
	return u2
}

// Valid operators
// eq			, neq			, gt			, gte					, lt			, lte					, like	, nlike		, in	, nin		, cf		, pl			, pr
// equal		, not equal		, greater than	, greater than or equal	, lower than	, lower than or equal	, like	, not like	, in	, not in	, child of	, parent left	, parent right
var ValidOperators = []string{"eq", "neq", "gt", "gte", "lt", "lte", "like", "nlike", "in", "nin", "cf", "pl", "pr"}

// asc: ASCENDING, desc: DESCENDING
var ValidOrders = []string{"asc", "desc"}

type SortValue struct {
	Order string
	Key   string
}

type FilterValue struct {
	Or       bool
	Operator string
	Value    interface{}
}

type URLRequest map[string][]FilterValue

type FilterValues []FilterValue

func (u URLRequest) GetOne(key string) (value *FilterValue) {
	if vs, exists := u[key]; exists {
		return &vs[0]
	}

	return nil
}

func parseSortValue(s string) SortValue {
	var (
		order = "asc"
		key   string

		extract = strings.Split(s, ":")
	)

	if len(extract) > 0 {
		key = extract[0]
	}

	if len(extract) > 1 && ArrayElementExists(ValidOrders, extract[1]) {
		order = extract[1]
	}

	return SortValue{Order: order, Key: key}
}

func parseFilterValue(s string) FilterValue {
	var (
		or       bool
		operator string
		value    interface{}
		extract  = strings.SplitN(s, ":", 2)
	)

	if len(extract) > 0 {
		operator = extract[0]
		if strings.HasPrefix(operator, "+") {
			or = true
			operator = strings.TrimPrefix(operator, "+")
		}
		if ArrayElementExists(ValidOperators, operator) && len(extract) > 1 {
			value = extract[1]
		} else {
			value = operator
			operator = "eq"
		}
	} else {
		value = s
		operator = "eq"
	}

	return FilterValue{Or: or, Operator: operator, Value: value}
}

func ParseURL(ctx echo.Context) (err error) {
	var (
		urlValues     = ctx.QueryParams()
		offset, limit int
		filters       = make(map[string][]FilterValue)
		sorts         []SortValue
	)

	// Remove unnecessary filter values
	for k, vs := range urlValues {
		if k == "limit" {
			if limit, err = strconv.Atoi(vs[0]); err != nil {
				return
			}
		} else if k == "offset" {
			if offset, err = strconv.Atoi(vs[0]); err != nil {
				return err
			}
		} else if k == "sort" {
			if vs[0] == "" {
				continue
			}

			sortsValue := strings.Split(vs[0], ",")
			for _, s := range sortsValue {
				sorts = append(sorts, parseSortValue(s))
			}
		} else {
			// There are filters
			// Create filter map
			for _, v := range vs {
				filters[k] = append(filters[k], parseFilterValue(v))
			}
		}
	}

	ctx.Set(CTXLimit, limit)
	ctx.Set(CTXOffset, offset)
	ctx.Set(CTXFilters, filters)
	ctx.Set(CTXSorts, sorts)

	return
}

// ParsePaginationParams parse pagination query params
func ParsePaginationParams(ctx echo.Context) (limit int, offset int, err error) {
	limit = cast.ToInt(ctx.QueryParam("limit"))
	offset = cast.ToInt(ctx.QueryParam("offset"))

	// if limit < 0 {
	// 	// err = errors.New("limit must be positive")
	// 	limit = 100
	// }

	// Maximum products per request
	if limit == 0 || limit > 100 {
		limit = 100
	}

	return int(limit), int(offset), err
}

func ParseContext(ctx echo.Context) (limit, offset int, filters map[string][]FilterValue, sorts []SortValue) {
	limit = ctx.Get(CTXLimit).(int)
	offset = ctx.Get(CTXOffset).(int)
	filters = ctx.Get(CTXFilters).(map[string][]FilterValue)
	sorts = ctx.Get(CTXSorts).([]SortValue)

	if limit <= 0 {
		limit = 5
	}

	if offset <= 0 {
		offset = 0
	}

	return
}
