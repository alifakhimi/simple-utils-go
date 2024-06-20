package simutils

import (
	"net/url"
	"reflect"
	"testing"
)

var (
	urls = map[string]bool{
		"https":                         false,
		"https://":                      false,
		"":                              false,
		"http://www":                    true,
		"http://www.dumpsters.com":      true,
		"https://www.dumpsters.com:443": true,
		"/testing-path":                 false,
		"testing-path":                  false,
		"alskjff#?asf//dfas":            false,
	}
)

func TestURL_IsValid(t *testing.T) {
	type fields struct {
		URL url.URL
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{}

	for k, want := range urls {
		parsed, _ := url.Parse(k)
		if parsed == nil {
			parsed = &url.URL{}
		}
		tests = append(tests, struct {
			name   string
			fields fields
			want   bool
		}{
			name: k,
			fields: fields{
				URL: *parsed,
			},
			want: want,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &URL{
				URL: tt.fields.URL,
			}
			if got := u.IsValid(); got != tt.want {
				t.Errorf("URL.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestURL_Scan(t *testing.T) {
	type fields struct {
		URL url.URL
	}
	type args struct {
		b interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "empty string",
			fields: fields{},
			args: args{
				b: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &URL{
				URL: tt.fields.URL,
			}
			if err := u.Scan(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("URL.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestURL_MarshalJSON(t *testing.T) {
	type fields struct {
		URL url.URL
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name:    "empty string",
			fields:  fields{},
			want:    []byte{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := URL{
				URL: tt.fields.URL,
			}
			got, err := u.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("URL.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("URL.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
