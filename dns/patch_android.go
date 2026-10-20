//go:build android && cmfa

package dns

import (
	"sync"

	"github.com/metacubex/mihomo/component/resolver"
)

var (
	systemResolverMu sync.RWMutex
	systemResolver   []dnsClient
)

// DefaultResolver 已聚合主、节点和直连解析器，不能再逐个重置共享客户端。
func FlushCacheWithDefaultResolver() {
	for _, r := range []resolver.Resolver{
		resolver.DefaultResolver,
		resolver.SystemResolver,
	} {
		if r != nil {
			r.ClearCache()
			r.ResetConnection()
		}
	}
}

func UpdateSystemDNS(addr []string) {
	ns := make([]NameServer, 0, len(addr))
	for _, d := range addr {
		ns = append(ns, NameServer{Addr: d})
	}

	clients := transform(ns, nil)
	systemResolverMu.Lock()
	systemResolver = clients
	systemResolverMu.Unlock()
}

func (c *systemClient) getDnsClients() ([]dnsClient, error) {
	systemResolverMu.RLock()
	clients := systemResolver
	systemResolverMu.RUnlock()
	// 发布后不再修改切片，读者取得快照后无需继续持锁。
	return clients, nil
}

func (c *systemClient) ResetConnection() {
	clients, _ := c.getDnsClients()
	for _, r := range clients {
		r.ResetConnection()
	}
}
