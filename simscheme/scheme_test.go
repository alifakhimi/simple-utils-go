package simscheme

import (
	"reflect"
	"testing"

	simutils "github.com/alifakhimi/simple-utils-go"
)

type product struct {
	simutils.Model
}

func TestGetKeys(t *testing.T) {
	type args struct {
		val any
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Get Keys from struct",
			args: args{
				val: &product{Model: simutils.Model{ID: 123}},
			},
			want: []string{
				"123",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetKeys(tt.args.val); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}
