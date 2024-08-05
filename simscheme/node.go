package simscheme

import (
	"encoding/json"
	"reflect"
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
