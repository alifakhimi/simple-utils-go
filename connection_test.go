package simutils

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Test_dialector(t *testing.T) {
	type args struct {
		dbConn *DBConnection
	}
	tests := []struct {
		name string
		args args
		want gorm.Dialector
	}{
		{
			name: "sqlite connection",
			args: args{
				dbConn: &DBConnection{
					DBConfig: DBConfig{
						Driver: SQLite,
						DSN:    "test.db",
						Debug:  true,
					},
				},
			},
			want: &sqlite.Dialector{DSN: "test.db"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dialector(tt.args.dbConn); got != nil {
				t.Errorf("dialector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConnect(t *testing.T) {
	type args struct {
		dbConn *DBConnection
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "connect sqlite",
			args: args{
				dbConn: &DBConnection{
					DBConfig: DBConfig{
						Driver: SQLite,
						DSN:    "test.db",
						Debug:  true,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Connect(tt.args.dbConn); (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// func TestRegexQuery(t *testing.T) {
// 	type args struct {
// 		dbConn *DBConnection
// 	}
// 	tests := []struct {
// 		name      string
// 		args      args
// 		wantError bool
// 	}{
// 		{
// 			name: "regexp (digits only)",
// 			args: args{
// 				dbConn: &DBConnection{
// 					DBConfig: DBConfig{
// 						Driver: SQLite,
// 						DSN:    "test.db",
// 						Debug:  true,
// 					},
// 				},
// 			},
// 			wantError: false,
// 		},
// 	}
// }
