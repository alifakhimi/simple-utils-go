package simscheme

import (
	"fmt"
	"reflect"

	simutils "github.com/alifakhimi/simple-utils-go"
)

type Document struct {
	Label     Label             `json:"label,omitempty"`
	Relations map[Key]*Relation `json:"relations,omitempty"`
	Nodes     map[Key]*Node     `json:"nodes,omitempty"`
}

func (schema *Schema) BuildDocumentLabel(keys ...Key) *Label {
	return schema.Label.Append(documentName, keys...)
}

func (schema *Schema) NewDocumentWithKeys(keys ...Key) *Document {
	return schema.NewDocumentWithLabel(
		schema.BuildDocumentLabel(keys...),
	)
}

func (schema *Schema) NewDocumentWithType(t any) *Document {
	return schema.NewDocumentWithLabel(
		schema.BuildDocumentLabel(Key(simutils.GetTableName(t))),
	)
}

func (schema *Schema) NewDocument() *Document {
	return schema.NewDocumentWithKeys(defDocumentLabel)
}

func (schema *Schema) NewDocumentWithLabel(label *Label) *Document {
	if label == nil {
		return schema.NewDocument()
	}

	doc := &Document{
		Label:     *label,
		Relations: make(map[Key]*Relation),
		Nodes:     make(map[Key]*Node),
	}

	return doc
}

func (doc *Document) GetData(dst any) error {
	v := reflect.ValueOf(dst)

	if v.Kind() != reflect.Ptr && v.Kind() != reflect.Pointer {
		// } else if v.Kind() == reflect.Array || v.Kind() == reflect.Slice {
		// 	v = reflect.ValueOf(&dst)
		// } else {
		return fmt.Errorf("error in dst type (needs pointer of slice) got %T expected *%T", dst, dst)
	}

	v = v.Elem()
	// fmt.Printf("dst type %T %s %v", dst, v.Kind(), v.Type())

	// if !v.CanSet() {
	// 	return fmt.Errorf("dst is not settable %T", dst)
	// }

	for _, node := range doc.Nodes {
		v.Set(reflect.Append(v, reflect.ValueOf(node.Data)))
	}

	return nil
}

func (doc *Document) AddNode(value any, relations ...any) *Node {
	// Add node to document
	node := doc.NewNode(value)
	key := node.Label.GetKey()

	if node, exists := doc.Nodes[key]; exists {
		return node
	}

	doc.Nodes[key] = node

	// Add relations to document
	for _, m := range relations {
		relNode := doc.AddNode(m)
		_ = doc.AddRelation(node, relNode)
	}

	return node
}

func (doc *Document) AddNodes(values any) *Document {
	v := reflect.ValueOf(values)

	if v.Kind() != reflect.Array && v.Kind() != reflect.Slice {
		// panic("Invalid data-type")
		return doc
	}

	for i := 0; i < v.Len(); i++ {
		value := v.Index(i).Interface()
		doc.AddNode(value)
	}

	return doc
}

func (doc *Document) GetOrInitNode(value any) *Node {
	if value == nil {
		return nil
	}

	return doc.AddNode(value)
}

func (doc *Document) GetNode(value any) *Node {
	if value == nil {
		return nil
	}

	label := doc.BuildNodeLabel(value)
	return doc.Nodes[label.GetKey()]
}

func (doc *Document) AddRelation(from, to *Node, values ...any) *Relation {
	newRel := doc.NewRelation(from, to, values)

	if rel, exists := doc.Relations[newRel.Label.GetKey()]; exists {
		return rel
	}

	doc.Relations[newRel.Label.GetKey()] = newRel

	return newRel
}
