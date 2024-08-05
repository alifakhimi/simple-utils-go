package simscheme

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	simutils "github.com/alifakhimi/simple-utils-go"
)

// Default Labels
const (
	keySeparator     = ";"
	relSeparator     = "-/"
	scopeSeparator   = "/"
	schemaName       = "sch"
	documentName     = "doc"
	nodeName         = "nod"
	relationName     = "rel"
	defSchemaLabel   = "default"
	defDocumentLabel = "default"
)

// Errors
var (
	ErrNodeAlreadyExists     = errors.New("node already exists")
	ErrRelationAlreadyExists = errors.New("relation already exists")
	ErrDocumentAlreadyExists = errors.New("document already exists")
	ErrLabelParse            = errors.New("label parse process failed")
)

var (
	defSchema = NewSchema(BuildSchemaLabel(defSchemaLabel))
)

type Key string

func Validate() error { panic("not implemented") }

func SetVersion(ver string) { panic("not implemented") }

func SetMeta(value any) { panic("not implemented") }

func Version() string { panic("not implemeneted") }

func GetSchema() *Schema { return defSchema }

func GetDocument() *Document { return defSchema.GetDocument() }

func GetDocumentByLabel(label *Label) *Document { return defSchema.GetDocumentByLabel(label) }

func GetScope(value any) string {
	return simutils.GetTableName(value)
}

func GetKeys(val any) []Key {
	keys := []Key{}
	v := reflect.ValueOf(val)

	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}

	// fmt.Println("Kind", v.Kind())

	// to find the kind
	if v.Kind() == reflect.Struct {
		// fmt.Println("Number of fields", v.NumField())
		for i := 0; i < v.NumField(); i++ {
			t := v.Type()

			// fmt.Printf("Field: %d \t type: %T \t value: %v\n",
			// 	i, v.Field(i), v.Field(i))

			f := t.Field(i)
			// fmt.Println("Field Name", f.Name, "Field Type", f.Type, "Tag", f.Tag)
			gormPrimaryTag := simutils.ItemExists(strings.Split(strings.ToLower(f.Tag.Get("gorm")), ";"), strings.ToLower("primaryKey"))
			if gormPrimaryTag {
				keys = append(keys, Key(fmt.Sprintf("%v", v.Field(i))))
			}
		}
	}

	return keys
}
