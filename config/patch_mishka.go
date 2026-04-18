//go:build mishka

package config

import (
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

// patchMishka adjusts the parsed RawConfig for Android usage:
//   - When the subscription has DNS disabled, inject a sensible default
//     (fake-ip mode, a mix of domestic and international nameservers, and a
//     fake-ip filter list covering STUN, game consoles, captive portals, etc.).
//     Applies to all tunnel modes (VPN / ROOT TUN / ROOT TPROXY) because dns
//     hijack/redirect leaves no responder when DNS is off.
//   - In VPN mode only (FileDescriptor > 0): when AppendSystemDNS is set
//     (either by the subscription or implied by the default injection above),
//     append "system://" so Android system DNS acts as a last-resort fallback.
//     Skipped under ROOT modes — ROOT TUN would re-enter sing-tun's catch-all
//     route and loop, and ROOT TPROXY has no Kotlin-side UpdateSystemDNS feed
//     either, so the entry would be a dead nameserver slot.
func patchMishka(cfg *RawConfig) {
	if !cfg.DNS.Enable {
		cfg.DNS = DefaultRawConfig().DNS
		cfg.DNS.Enable = true
		cfg.DNS.NameServer = mishkaDefaultNameServers
		cfg.DNS.EnhancedMode = C.DNSFakeIP
		cfg.DNS.FakeIPRange = mishkaDefaultFakeIPRange
		cfg.DNS.FakeIPFilter = mishkaDefaultFakeIPFilter
		if cfg.Tun.FileDescriptor > 0 {
			cfg.ClashForAndroid.AppendSystemDNS = true
		}
	}
	if cfg.Tun.FileDescriptor > 0 && cfg.ClashForAndroid.AppendSystemDNS {
		cfg.DNS.NameServer = append(cfg.DNS.NameServer, "system://")
	}
}
