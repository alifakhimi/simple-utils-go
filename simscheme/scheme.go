package simscheme

import (
	"errors"
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
