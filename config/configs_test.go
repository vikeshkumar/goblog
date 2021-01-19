package config

import (
	"reflect"
	"testing"
	"time"
)

func Test_cfg_GetString(t *testing.T) {
	cfg, e1 := Init("conf_test.toml", ".", "toml")
	tests := []struct {
		name    string
		cfg     Configuration
		args    string
		want    string
		error   error
		wantErr bool
	}{
		{"can read z", cfg, "z", "z", e1, false},
		{"can read a.b", cfg, "a.b", "b", e1, false},
		{"can read a.c", cfg, "a.c", "c", e1, false},
		{"can read d.e.f", cfg, "d.e.f", "f", e1, false},
		{"can read d.e.g", cfg, "d.e.g", "g", e1, false},
		{"a.b.c.d should return empty string", cfg, "a.b.c.d", "", e1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.wantErr {
			case false:
				if got := tt.cfg.GetString(tt.args); got != tt.want || tt.error != nil {
					t.Errorf("GetString() = %v, want %v", got, tt.want)
				}
			case true:
				if tt.error == nil || tt.cfg != nil {
					t.Errorf("Got error = %v, Got Config = %c", tt.error, tt.cfg)
				}
			}

		})
	}
}

func TestInit(t *testing.T) {
	cfg, e1 := Init("conf_test.toml", ".", "toml")
	eCfg, e2 := Init("conf_test2.toml", ".", "toml")
	tests := []struct {
		name    string
		cfg     Configuration
		error   error
		wantErr bool
	}{
		{"a.b.c.d should return empty string", cfg, e1, false},
		{"a.b.c.d should return empty string", eCfg, e2, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.wantErr {
			case false:
				if tt.error != nil || tt.cfg == nil {
					t.Errorf("Got error = %v, wanted = %v", tt.error, tt.wantErr)
				}
			case true:
				if tt.error == nil || tt.cfg != nil {
					t.Errorf("wanted error, got = %v, wanted nil config, got = %c", tt.error, tt.cfg)
				}
			}

		})
	}
}

func Test_cfg_GetDuration(t *testing.T) {
	cfg, _ := Init("conf_test.toml", ".", "toml")
	tests := []struct {
		name   string
		fields string
		config Configuration
		want   time.Duration
	}{
		{"able to read duration", "server.duration", cfg, 10},
		{"able to read duration", "duration", cfg, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.GetDuration(tt.fields); got != tt.want {
				t.Errorf("GetDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name   string
		pre    func()
		create bool
	}{
		{"should return nil", func() {}, false},
		{"should return c1", func() {
			Init("conf_test.toml", ".", "toml")
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pre()
			if got := Get(); !reflect.DeepEqual(got, c) {
				t.Errorf("Get() = %v, want %v", got, c)
			}
		})
	}
}
