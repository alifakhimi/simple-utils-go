package simutils

import "testing"

func TestPID_IsValid(t *testing.T) {
	tests := []struct {
		name string
		id   PID
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.id.IsValid(); got != tt.want {
				t.Errorf("PID.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullPID_IsValid(t *testing.T) {
	type fields struct {
		PID   PID
		Valid bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := NullPID{
				PID:   tt.fields.PID,
				Valid: tt.fields.Valid,
			}
			if got := id.IsValid(); got != tt.want {
				t.Errorf("NullPID.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	pid := PID(12)
	var nilPointerPID *PID

	type args struct {
		id interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "nil *PID",
			args: args{
				id: nilPointerPID,
			},
			want: false,
		},
		{
			name: "Valid PID",
			args: args{
				id: PID(12345),
			},
			want: true,
		},
		{
			name: "Valid PID (1..9 digits)",
			args: args{
				id: 123456789,
			},
			want: true,
		},
		{
			name: "Valid PID (string)",
			args: args{
				id: "12345",
			},
			want: true,
		},
		{
			name: "Invalid PID (non-digit characters)",
			args: args{
				id: "124578B",
			},
			want: false,
		},
		{
			name: "Invalid PID (big string)",
			args: args{
				id: "12345678901234567890",
			},
			want: false,
		},
		{
			name: "Pointer PID",
			args: args{
				id: &pid,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValid(tt.args.id); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
