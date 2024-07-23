package simscheme

import (
	"reflect"
	"testing"
)

func TestParseLabel(t *testing.T) {
	type args struct {
		t      string
		parent *Label
	}
	tests := []struct {
		name    string
		args    args
		want    Key
		wantErr bool
	}{
		{
			name: "parent and multi keys omit",
			args: args{
				t:      "(nod:[1])",
				parent: nil,
			},
			want:    Key("(nod:[1])"),
			wantErr: false,
		},
		{
			name: "with parent and multi keys",
			args: args{
				t:      "(doc:[default])" + scopeSeparator + "(nod:[1;2;123])",
				parent: BuildLabel("doc", "def"),
			},
			want:    Key("(doc:[default])" + scopeSeparator + "(nod:[1;2;123])"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLabel(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLabel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.GetKey(), tt.want) {
				t.Errorf("ParseLabel() = %v, want %v", got, tt.want)
			}
		})
	}
}
