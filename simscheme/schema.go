package simscheme

import simutils "github.com/alifakhimi/simple-utils-go"

type Schema struct {
	Label     Label             `json:"label,omitempty"`
	Relations map[Key]*Relation `json:"relations,omitempty"`
	Documents map[Key]*Document `json:"documents,omitempty"`
}

func NewSchema(label *Label) *Schema {
	if label == nil {
		label = BuildSchemaLabel()
	}

	schema := &Schema{
		Label:     *label,
		Relations: make(map[Key]*Relation),
		Documents: make(map[Key]*Document),
	}

	_ = schema.AddNewDocument()

	return schema
}

func BuildSchemaLabel(keys ...Key) *Label {
	return BuildLabel(schemaName, keys...)
}

func (schema *Schema) AddNewDocument() *Document {
	return schema.AddDocument(schema.NewDocument())
}

func (schema *Schema) AddNewDocumentWithKey(keys ...Key) *Document {
	return schema.AddDocument(
		schema.NewDocumentWithKeys(keys...),
	)
}

func (schema *Schema) AddNewDocumentWithType(t any) *Document {
	return schema.AddDocument(
		schema.NewDocumentWithType(t),
	)
}

func (schema *Schema) AddDocument(doc *Document) *Document {
	if doc == nil {
		return schema.AddNewDocument()
	}

	if doc, exists := schema.Documents[doc.Label.GetKey()]; exists {
		return doc
	}

	schema.Documents[doc.Label.GetKey()] = doc

	return doc
}

func (schema *Schema) GetDocument() *Document {
	return schema.GetDocumentByLabel(schema.BuildDocumentLabel(defDocumentLabel))
}

func (schema *Schema) GetDocumentByType(t any) *Document {
	return schema.GetDocumentByLabel(
		schema.BuildDocumentLabel(Key(simutils.GetTableName(t))),
	)
}

func (doc *Document) ReadAll(dst any) error {
	return nil
}

func (schema *Schema) GetDocumentByLabel(label *Label) *Document {
	if label == nil {
		return schema.GetDocument()
	}

	return schema.Documents[label.GetKey()]
}
