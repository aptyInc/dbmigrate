package target

import "testing"

func Test_getSQLInfo(t *testing.T) {
	tests := []struct {
		name    string
		dbURL   string
		want    string
		wantErr bool
	}{
		{
			name:    "validURL",
			dbURL:   "postgres://postgres:pass@localhost:5432/test",
			want:    "host=localhost port=5432 user=postgres password=pass dbname=test sslmode=disable",
			wantErr: false,
		},
		{
			name:    "InvalidDB",
			dbURL:   "my://postgres:pass@localhost:5432/test",
			want:    "",
			wantErr: true,
		},
		{
			name:    "InvalidURL",
			dbURL:   "my://postgres:localhost:5432/test",
			want:    "",
			wantErr: true,
		},
		{
			name:    "InvalidURL",
			dbURL:   "my://postgres:localhost:5432/test",
			want:    "",
			wantErr: true,
		},
		{
			name:    "InvalidPort",
			dbURL:   "postgres://postgres:pass@localhost5432/test",
			want:    "",
			wantErr: true,
		},
		{
			name:    "EmptyPassword",
			dbURL:   "postgres://postgres@localhost:5432/test",
			want:    "",
			wantErr: true,
		},
		{
			name:    "NoDBName",
			dbURL:   "postgres://postgres:pass@localhost:5432",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getSQLInfo(tt.dbURL, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("getSQLInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getSQLInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
