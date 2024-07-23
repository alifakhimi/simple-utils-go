package simscheme

import "strings"

type Relation struct {
	Label Label `json:"label,omitempty"`
	Data  any   `json:"data,omitempty"`
	From  *Node `json:"-"`
	To    *Node `json:"-"`
}

func (doc *Document) NewRelation(from, to *Node, values ...any) *Relation {
	label := doc.BuildRelationLabel(from, to)

	var data any
	if len(values) > 0 {
		data = values
	}

	return &Relation{
		Label: *label,
		From:  from,
		To:    to,
		Data:  data,
	}
}

func (doc *Document) BuildRelationLabel(from, to *Node) *Label {
	return doc.Label.Append(relationName,
		Key(strings.Join(
			[]string{
				string(from.Label.GetKey()),
				string(to.Label.GetKey()),
			},
			relSeparator,
		)),
	)
}
