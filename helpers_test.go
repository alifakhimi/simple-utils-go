package simutils

import "testing"

func TestAnyItemExists(t *testing.T) {
	type args struct {
		searchArrayType interface{}
		checkArrayType  interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "array chack false",
			args: args{
				searchArrayType: []string{"a", "b"},
				checkArrayType:  []string{"c"},
			},
			want: false,
		},
		{
			name: "array check true",
			args: args{
				searchArrayType: []string{"a", "b", "c"},
				checkArrayType:  []string{"a"},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AnyItemExists(tt.args.searchArrayType, tt.args.checkArrayType); got != tt.want {
				t.Errorf("AnyItemExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetToNilIfZeroValue(t *testing.T) {
	var (
		x  int
		y  float64
		z  string
		xp *int
	)

	type args struct {
		value any
	}
	tests := []struct {
		name string
		args args
		want *any
	}{
		{
			name: "int",
			args: args{
				value: x,
			},
			want: nil,
		},
		{
			name: "float64",
			args: args{
				value: y,
			},
			want: nil,
		},
		{
			name: "string",
			args: args{
				value: z,
			},
			want: nil,
		},
		{
			name: "*int",
			args: args{
				value: xp,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetToNilIfZeroValue(tt.args.value); got != tt.want {
				t.Errorf("SetToNilIfZeroValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
