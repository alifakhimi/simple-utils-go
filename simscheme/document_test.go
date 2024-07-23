package simscheme

import (
	"testing"

	simutils "github.com/alifakhimi/simple-utils-go"
)

type unittype struct {
	simutils.Model
}

var unittypes = []*unittype{
	{Model: simutils.Model{ID: 1}},
	{Model: simutils.Model{ID: 2}},
	{Model: simutils.Model{ID: 3}},
	{Model: simutils.Model{ID: 4}},
}

func TestDocument_GetData(t *testing.T) {
	type fields struct {
		Label     Label
		Relations map[Key]*Relation
		Nodes     map[Key]*Node
	}
	type args struct {
		dst any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "slice ([]*model.UnitType)",
			fields: fields{
				Nodes: map[Key]*Node{
					Key("1"): {
						Data: unittypes[0],
					},
					Key("2"): {
						Data: unittypes[1],
					},
					Key("3"): {
						Data: unittypes[2],
					},
				},
			},
			args: args{
				dst: []*unittype{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &Document{
				Label:     tt.fields.Label,
				Relations: tt.fields.Relations,
				Nodes:     tt.fields.Nodes,
			}
			if err := doc.GetData(tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("Document.GetData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
