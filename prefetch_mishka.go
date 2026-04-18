//go:build mishka

package main

import (
	"runtime"
	"sync"

	"github.com/metacubex/mihomo/config"
	C "github.com/metacubex/mihomo/constant"
	P "github.com/metacubex/mihomo/constant/provider"
	"github.com/metacubex/mihomo/hub/executor"
	"github.com/metacubex/mihomo/log"
)

// runPrefetch 解析配置、并发下载所有 HTTP Provider 到本地后退出。
func runPrefetch(configBytes []byte) int {
	var (
		cfg *config.Config
		err error
	)
	if len(configBytes) != 0 {
		cfg, err = executor.ParseWithBytes(configBytes)
	} else {
		cfg, err = executor.Parse()
	}
	if err != nil {
		log.Errorln("prefetch parse %s error: %v", C.Path.Config(), err.Error())
		return 1
	}
	defer runtime.KeepAlive(cfg)

	var wg sync.WaitGroup
	for name, pv := range cfg.Providers {
		if pv.VehicleType() != P.HTTP {
			continue
		}
		wg.Add(1)
		go func(name string, pv P.ProxyProvider) {
			defer wg.Done()
			if err := pv.Update(); err != nil {
				log.Warnln("prefetch proxy provider %s error: %v", name, err)
				return
			}
			log.Infoln("prefetch proxy provider %s done", name)
		}(name, pv)
	}
	for name, pv := range cfg.RuleProviders {
		if pv.VehicleType() != P.HTTP {
			continue
		}
		wg.Add(1)
		go func(name string, pv P.RuleProvider) {
			defer wg.Done()
			if err := pv.Update(); err != nil {
				log.Warnln("prefetch rule provider %s error: %v", name, err)
				return
			}
			log.Infoln("prefetch rule provider %s done", name)
		}(name, pv)
	}
	wg.Wait()
	return 0
}

func init() {
	prefetchRunner = runPrefetch
}
