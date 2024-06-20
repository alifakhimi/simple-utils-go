package simutils

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

const (
	// Secret ...
	Secret string = "SecretKey" // replace SecretKey with strong key
	// JWT Claims -----------------------------------------------------------------------
	// ClaimsClientKey ...
	ClaimsClientKey string = "key"
	// ClaimsUsername ...
	ClaimsUsername string = "uname"
	// ClaimsRole ...
	ClaimsRole string = "rol"
	// ClaimsExpireTime ...
	ClaimsExpireTime string = "exp"
	// Context Headers ------------------------------------------------------------------
	// HeadersUser ...
	HeadersUser string = "x-user"
	// HeadersUserID ...
	HeadersUserID string = "x-user-id"
	// HeadersClient ...
	HeadersClient string = "x-client"
	// HeadersCustomer ...
	HeadersCustomer string = "x-customer"
	// HeadersAuthenticated ...
	HeadersAuthenticated string = "x-authenticated"
	// HeadersTokenKey ...
	HeadersTokenKey string = "x-token-key"
	// HeadersDatabase ...
	HeadersDatabase string = "x-database"
	// HeadersRestAPI ...
	HeadersRestAPI string = "x-rest-api"
	// HeadersTracingContext ...
	HeadersTracingContext = "x-tracing-context"
	HeaderPlatform        = "x-platform"
	HeaderSource          = "x-source"
	// limit
	CTXLimit string = "x-limit"
	// offset
	CTXOffset string = "x-offset"
	// filters
	CTXFilters string = "x-filters"
	// sorts
	CTXSorts string = "x-sorts"
)

var (
	alphanumericChars = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	digitChars        = []byte("0123456789")
	verificationChars = []byte("123456789")
)

// GetAbsPath ...
func GetAbsPath(filename string) string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	return exPath + "/" + filename
}

func Replace(text string, old []string, new []string) (string, error) {
	if len(old) != len(new) {
		return "", errors.New("the length of old and new must be same")
	}

	var s string

	for i, o := range old {
		s = strings.Replace(text, o, new[i], -1)
	}

	return s, nil
}

func ReplaceMap(text string, m map[string]interface{}) string {
	if len(m) == 0 {
		return text
	}

	for k, v := range m {
		text = strings.Replace(text, k, fmt.Sprintf("%v", v), -1)
	}

	return text
}

// GetUintFromString ...
func GetUintFromString(s string) (uint, error) {
	idP, err := strconv.Atoi(s)
	uID := uint(idP)

	return uID, err
}

// GetIntFromString ...
func GetIntFromString(s string) (int, error) {
	idP, err := strconv.Atoi(s)

	return idP, err
}

// GetStringsDefault ...
func GetStringsDefault(s []string, def []string) []string {
	if len(s) > 0 {
		return s
	}

	return def
}

// GetStringDefault ...
func GetStringDefault(s string, def string) string {
	if strings.TrimSpace(s) == "" {
		return def
	}

	return s
}

// GetIntDefault ...
func GetIntDefault(s string, def int) int {
	v, err := strconv.ParseInt(s, 0, 0)

	if err != nil {
		return def
	}

	return int(v)
}

// GetBoolsDefault ...
func GetBoolsDefault(params []string, def []bool) (values []bool) {
	if len(params) > 0 {
		for _, param := range params {
			if value, err := strconv.ParseBool(param); err == nil {
				values = append(values, value)
			}
		}

		return values
	}

	return def
}

// GetBoolDefault ...
func GetBoolDefault(s string, def bool) bool {
	v, err := strconv.ParseBool(s)

	if err != nil {
		return def
	}

	return v
}

// GetBool ...
func GetBool(s string) (bool, error) {
	if s == "" {
		return false, nil
	}
	return strconv.ParseBool(s)
	// if s == "true" {
	// 	return true, nil
	// } else if s == "false" || s == "" {
	// 	return false, nil
	// } else {
	// 	return false, errors.New("invalid value")
	// }
}

// GenerateVerificationCode ...
func GenerateVerificationCode(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = verificationChars[int(b[i])%len(verificationChars)]
	}
	return string(b)
}

// GenerateCode ...
func GenerateCode(max int) (string, error) {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		return "", err
	}
	for i := 0; i < len(b); i++ {
		b[i] = digitChars[int(b[i])%len(digitChars)]
	}

	return string(b), nil
}

// RandStringCode ...
func RandStringCode(max int) (string, error) {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		return "", err
	}
	for i := range b {
		b[i] = alphanumericChars[int(b[i])%len(alphanumericChars)]
	}
	return string(b), nil
}

// AnyItemExists ...
func AnyItemExists(searchArrayType interface{}, checkArrayType interface{}) bool {
	if searchArrayType == nil || checkArrayType == nil {
		return false
	}

	var (
		searchArr = reflect.ValueOf(searchArrayType)
	)

	if searchArr.Kind() != reflect.Array && searchArr.Kind() != reflect.Slice {
		// panic("Invalid data-type")
		return false
	}

	for i := 0; i < searchArr.Len(); i++ {
		if ItemExists(checkArrayType, searchArr.Index(i).Interface()) {
			return true
		}
	}

	return false
}

// ReverseSlice reverece order slice
func ReverseSlice(data interface{}) {
	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Slice {
		panic(errors.New("data must be a slice type"))
	}
	valueLen := value.Len()
	for i := 0; i <= int((valueLen-1)/2); i++ {
		reverseIndex := valueLen - 1 - i
		tmp := value.Index(reverseIndex).Interface()
		value.Index(reverseIndex).Set(value.Index(i))
		value.Index(i).Set(reflect.ValueOf(tmp))
	}
}

// PurgeArray purge array
func PurgeArray(arrayType interface{}) interface{} {
	if arrayType == nil {
		return nil //, errors.New("arrayType is nil")
	}

	arr := reflect.ValueOf(arrayType)

	if arr.Kind() != reflect.Array && arr.Kind() != reflect.Slice {
		return nil //, errors.New("arrayType must be array or slice")
	}

	check := make(map[interface{}]int)

	for i := 0; i < arr.Len(); i++ {
		check[arr.Index(i).Interface()] = 1
	}

	typ := reflect.TypeOf(arrayType).Elem()
	temp := make([]interface{}, len(check))
	for key := range check {
		temp = append(temp, key)
	}

	t := reflect.MakeSlice(reflect.SliceOf(typ), 0, len(check))
	reflect.Copy(t, reflect.ValueOf(temp))

	return t.Interface() //, nil
}

// ItemExists ...
func ItemExists(arrayType interface{}, item interface{}) bool {
	if arrayType == nil {
		return false
	}

	arr := reflect.ValueOf(arrayType)

	// fmt.Println("array:", arr.Kind().String())

	if arr.Kind() != reflect.Array && arr.Kind() != reflect.Slice {
		return false
		// panic("Invalid data-type")
	}

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

// RemoveIndex remove index from slice or array
func RemoveIndex(arrayType interface{}, i int) (interface{}, error) {
	if arrayType == nil {
		return nil, errors.New("array is nil")
	}

	arr := reflect.ValueOf(arrayType)
	if arr.Kind() != reflect.Array && arr.Kind() != reflect.Slice {
		return nil, errors.New("arrayType must be array or slice")
	}

	if i >= arr.Len() {
		return nil, errors.New("index is upper than arrayType length")
	}

	newArr := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(arrayType).Elem()), 0, arr.Len()-1)

	for idx := 0; idx < arr.Len(); idx++ {
		if idx == i {
			continue
		}

		newArr = reflect.Append(newArr, arr.Index(idx).Elem())
	}

	return newArr, nil
}

// TrimString ...
func TrimString(s string) (value string) {
	value = strings.TrimSpace(s)
	value = strings.ToValidUTF8(value, "")
	return value
}

// MergeMaps overwriting duplicate keys, you should handle that if there is a need
func MergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

func JSONTo(src interface{}, dst interface{}) (err error) {
	var js []byte

	if js, err = json.Marshal(src); err != nil {
		return err
	}

	if err = json.Unmarshal(js, dst); err != nil {
		return err
	}

	return nil
}

// Cast ...
func Cast(src interface{}, dst interface{}) {
	dstValue := reflect.ValueOf(dst)
	dstElem := dstValue.Elem()
	countDstElems := dstElem.NumField()
	dstType := dstElem.Type()

	for i := 0; i < countDstElems; i++ {
		dstField := dstElem.Field(i)
		dstFieldName := dstType.Field(i).Name
		srcValue := reflect.ValueOf(src)
		srcElem := srcValue.Elem()
		countSrcElems := srcElem.NumField()
		srcType := srcElem.Type()

		if i == 0 {
			valID := reflect.Indirect(srcValue).Field(i)
			dstField.Set(valID)
			continue
		}

		for j := 1; j < countSrcElems; j++ {
			srcFieldName := srcType.Field(j).Name

			if dstFieldName == srcFieldName {
				// FOUND
				val := reflect.Indirect(srcValue).Field(j)
				dstField.Set(val)
				break
			}
		}
	}
}

func SetToNilIfZeroValue[T any](value T) *T {
	// با استفاده از reflect مقدار و نوع ورودی را دریافت می‌کنیم
	val := reflect.ValueOf(value)
	kind := val.Kind()

	// بررسی مقدار صفر بودن و نوع ورودی
	if val.IsZero() && kind != reflect.Func && kind != reflect.Chan && kind != reflect.Interface {
		// اگر مقدار صفر بود و نوع غیر قابل تغییر بود، پوینتر نال برگردان
		return nil
	}

	// اگر مقدار صفر نبود، یا ورودی قابل تغییر بود، ورودی را به صورت پوینتر برگردان
	return &value
}
