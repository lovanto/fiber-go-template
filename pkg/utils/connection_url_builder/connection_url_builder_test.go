package connection_url_builder

import (
	"os"
	"testing"
)

func TestConnectionURLBuilder(t *testing.T) {
	orig := os.Getenv
	restore := func(key, val string) {
		if val == "" {
			os.Unsetenv(key)
		} else {
			if err := os.Setenv(key, val); err != nil {
				t.Fatalf("failed to set env %s: %v", key, err)
			}
		}
	}

	tests := []struct {
		name    string
		n       string
		env     map[string]string
		wantURL string
		wantErr bool
	}{
		{
			name: "postgres", n: "postgres", wantErr: false,
			env: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "5432",
				"DB_USER":     "user",
				"DB_PASSWORD": "pass",
				"DB_NAME":     "dbname",
				"DB_SSL_MODE": "disable",
			},
			wantURL: "host=localhost port=5432 user=user password=pass dbname=dbname sslmode=disable",
		},
		{
			name: "mysql", n: "mysql", wantErr: false,
			env: map[string]string{
				"DB_HOST":     "mysql.host",
				"DB_PORT":     "3306",
				"DB_USER":     "muser",
				"DB_PASSWORD": "mpass",
				"DB_NAME":     "mydb",
			},
			wantURL: "muser:mpass@tcp(mysql.host:3306)/mydb",
		},
		{
			name: "redis", n: "redis", wantErr: false,
			env: map[string]string{
				"REDIS_HOST": "redis.host",
				"REDIS_PORT": "6379",
			},
			wantURL: "redis.host:6379",
		},
		{
			name: "fiber", n: "fiber", wantErr: false,
			env: map[string]string{
				"SERVER_HOST": "0.0.0.0",
				"SERVER_PORT": "8080",
			},
			wantURL: "0.0.0.0:8080",
		},
		{
			name: "unsupported", n: "mongo", wantErr: true,
			env:     map[string]string{},
			wantURL: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backups := map[string]string{}
			for k, v := range tt.env {
				backups[k] = orig(k)
				if err := os.Setenv(k, v); err != nil {
					t.Fatalf("failed to set env %s: %v", k, err)
				}
			}

			got, err := ConnectionURLBuilder(tt.n)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error=%v got %v", tt.wantErr, err)
			}
			if got != tt.wantURL {
				t.Fatalf("expected url %q got %q", tt.wantURL, got)
			}

			for k := range tt.env {
				restore(k, backups[k])
			}
		})
	}
}
