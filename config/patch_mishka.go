//go:build mishka

package config

import (
	"slices"

	C "github.com/metacubex/mihomo/constant"
)

var (
	mishkaDefaultNameServers = []string{
		"223.5.5.5",
		"119.29.29.29",
		"8.8.4.4",
		"1.0.0.1",
	}
	mishkaDefaultFakeIPFilter = []string{
		// STUN
		"+.stun.*.*",
		"+.stun.*.*.*",
		"+.stun.*.*.*.*",
		"+.stun.*.*.*.*.*",
		// Google Voices
		"lens.l.google.com",
		// Nintendo Switch STUN
		"*.n.n.srv.nintendo.net",
		// PlayStation STUN
		"+.stun.playstation.net",
		// Xbox
		"xbox.*.*.microsoft.com",
		"*.*.xboxlive.com",
		// Microsoft Captive Portal
		"*.msftncsi.com",
		"*.msftconnecttest.com",
		// Bilibili CDN
		"*.mcdn.bilivideo.cn",
		// Windows default LAN workgroup
		"WORKGROUP",
	}
	mishkaDefaultFakeIPRange = "28.0.0.0/8"
)

func init() {
	mishkaPatch = patchMishka
}

// 未启用 DNS 时补 fake-ip 配置，否则 TUN 接管查询后无人应答。
// AppendSystemDNS 把 system:// 加入主 nameserver 并发查询，并非失败后才回退。
// fd 模式只保留 Android 可用的 TUN 字段。
func patchMishka(cfg *RawConfig) {
	if !cfg.DNS.Enable {
		cfg.DNS = DefaultRawConfig().DNS
		cfg.DNS.Enable = true
		cfg.DNS.NameServer = mishkaDefaultNameServers
		cfg.DNS.EnhancedMode = C.DNSFakeIP
		cfg.DNS.FakeIPRange = mishkaDefaultFakeIPRange
		cfg.DNS.FakeIPFilter = mishkaDefaultFakeIPFilter
		cfg.ClashForAndroid.AppendSystemDNS = true
	}
	if cfg.ClashForAndroid.AppendSystemDNS && !slices.Contains(cfg.DNS.NameServer, "system://") && !slices.Contains(cfg.DNS.NameServer, "system") {
		cfg.DNS.NameServer = append(cfg.DNS.NameServer, "system://")
	}

	// ROOT 模式（fd == 0）走 mihomo TUN（netlink + tun.New()）路径。
	if cfg.Tun.FileDescriptor > 0 {
		// 白名单外的所有 RawTun 字段归零。
		cfg.Tun = RawTun{
			Enable:               cfg.Tun.Enable,
			Device:               cfg.Tun.Device,
			Stack:                cfg.Tun.Stack,
			DNSHijack:            cfg.Tun.DNSHijack,
			AutoRoute:            cfg.Tun.AutoRoute,
			AutoDetectInterface:  cfg.Tun.AutoDetectInterface,
			Inet6Address:         cfg.Tun.Inet6Address,
			MTU:                  cfg.Tun.MTU,
			CongestionController: cfg.Tun.CongestionController,
			FileDescriptor:       cfg.Tun.FileDescriptor,
			IncludePackage:       cfg.Tun.IncludePackage,
			ExcludePackage:       cfg.Tun.ExcludePackage,
		}
	}
}
