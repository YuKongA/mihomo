//go:build mishka

package config

import (
	"reflect"
	"testing"
)

func TestMishkaVPNPatchDoesNotChangeRoot(t *testing.T) {
	for _, fd := range []int{0, 42} {
		cfg := DefaultRawConfig()
		cfg.DNS.Enable = true
		cfg.Interface = "wlan0"
		cfg.RoutingMark = 123
		cfg.ExternalControllerTLS = "127.0.0.1:9443"
		cfg.RedirPort = 7892
		cfg.TProxyPort = 7893
		cfg.Tun.FileDescriptor = fd
		cfg.Tun.CongestionController = "bbr"
		cfg.Tun.IncludePackage = []string{"example.include"}
		cfg.Tun.ExcludePackage = []string{"example.exclude"}
		cfg.Listeners = []map[string]any{
			{"type": "http", "name": "http"},
			{"type": "tun", "name": "tun"},
			{"type": "tproxy", "name": "tproxy"},
			{"type": "redir", "name": "redir"},
			{"type": "socks", "name": "socks"},
		}
		before := *cfg
		patchMishka(cfg)
		if fd == 0 {
			if !reflect.DeepEqual(*cfg, before) {
				t.Fatal("ROOT configuration changed")
			}
			continue
		}
		if cfg.Interface != "" || cfg.RoutingMark != 0 || cfg.ExternalControllerTLS != "" {
			t.Fatal("VPN inherited interface, mark, or TLS controller")
		}
		if cfg.Tun.CongestionController != "" || len(cfg.Tun.IncludePackage) != 0 || len(cfg.Tun.ExcludePackage) != 0 {
			t.Fatal("VPN retained unsupported TUN fields")
		}
		if len(cfg.Listeners) != 2 || cfg.Listeners[0]["type"] != "http" || cfg.Listeners[1]["type"] != "socks" {
			t.Fatalf("VPN listeners = %#v", cfg.Listeners)
		}
	}
}
