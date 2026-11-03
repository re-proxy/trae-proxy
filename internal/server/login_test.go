package server

import (
	"strings"
	"testing"
)

// resolveCallbackURL 优先级：请求体 callback_url > 服务端 callback_url 配置 > 127.0.0.1:port/authorize。
func TestResolveCallbackURL(t *testing.T) {
	cases := []struct {
		name    string
		cfgURL  string
		reqURL  string
		reqPort string
		want    string
		wantErr bool
	}{
		{
			name: "default_local",
			want: "http://127.0.0.1:18080/authorize",
		},
		{
			name:    "default_custom_port",
			reqPort: "19000",
			want:    "http://127.0.0.1:19000/authorize",
		},
		{
			name:   "server_config_public",
			cfgURL: "https://trae.example.com/authorize",
			want:   "https://trae.example.com/authorize",
		},
		{
			name:    "server_config_wins_over_port",
			cfgURL:  "https://trae.example.com/authorize",
			reqPort: "19000",
			want:    "https://trae.example.com/authorize",
		},
		{
			name:   "request_overrides_server_config",
			cfgURL: "https://trae.example.com/authorize",
			reqURL: "http://10.0.0.5:7864/authorize",
			want:   "http://10.0.0.5:7864/authorize",
		},
		{
			name:    "request_rejects_bare_port",
			reqURL:  "127.0.0.1:18080/authorize",
			wantErr: true,
		},
		{
			name:    "request_rejects_bad_scheme",
			reqURL:  "ftp://trae.example.com/authorize",
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := &Handler{cfg: Config{CallbackURL: tc.cfgURL}}
			got, err := h.resolveCallbackURL(tc.reqURL, tc.reqPort)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("want error, got %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("resolve: %v", err)
			}
			if got != tc.want {
				t.Errorf("callback url=%q want %q", got, tc.want)
			}
		})
	}
}

// 公网回调地址必须原样出现在登录 URL 的 auth_callback_url 参数里（URL 编码后）。
func TestBuildLoginURLUsesGivenCallback(t *testing.T) {
	cb := "https://trae.example.com/authorize"
	u := BuildLoginURL("machine1", "device1", cb)
	if !strings.Contains(u, "auth_callback_url=https%3A%2F%2Ftrae.example.com%2Fauthorize") {
		t.Errorf("login url missing encoded callback: %s", u)
	}
}
