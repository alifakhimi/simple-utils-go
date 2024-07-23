package simscheme

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	simutils "github.com/alifakhimi/simple-utils-go"
)

type Node struct {
	Label Label `json:"label,omitempty"`
	Data  any   `json:"data,omitempty"`
	Meta  any   `json:"meta,omitempty"`
}

func (node *Node) Scope() string {
	return node.Label.Scope
}

func (node *Node) SetData(t any) *Node {
	if t == nil {
		return node
	}

	// is already unmarshaled
	if reflect.TypeOf(node.Data) == reflect.TypeOf(t) {
		// v := reflect.ValueOf(t)
		// v.Set(reflect.ValueOf(node.Data))
		return node
	}

	data, err := json.Marshal(node.Data)
	if err != nil {
		return node
	}

	err = json.Unmarshal(data, &t)
	if err != nil {
		return node
	}

	node.Data = t

	return node
}

func GetScope(value any) string {
	return simutils.GetTableName(value)
}

func (doc *Document) BuildNodeLabel(value any) *Label {
	scope := GetScope(value)
	if len(doc.Label.Keys) > 0 {
		scope = string(doc.Label.Keys[0])
	}
	return doc.Label.Append(scope, GetKeys(value)...)
}

func (doc *Document) NewNode(value any) *Node {
	return doc.NewNodeWithLabel(
		doc.BuildNodeLabel(value),
		value,
	)
}

func (doc *Document) NewNodeWithLabel(label *Label, value any) *Node {
	if label == nil {
		return doc.NewNode(value)
	}

	node := &Node{
		Label: *label,
		Data:  value,
	}

	return node
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
