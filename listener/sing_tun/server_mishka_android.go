//go:build mishka

package sing_tun

import (
	tun "github.com/metacubex/sing-tun"
	"golang.org/x/sys/unix"
)

// fd 是 runtime 的稳定来源，listener 只持有副本；重建栈不能关掉下一次 reload 要用的 fd。
func tunNew(options tun.Options) (tun.Tun, error) {
	if options.FileDescriptor == 0 {
		return tun.New(options)
	}
	fd, err := unix.FcntlInt(uintptr(options.FileDescriptor), unix.F_DUPFD_CLOEXEC, 0)
	if err != nil {
		return nil, err
	}
	options.FileDescriptor = fd
	device, err := tun.New(options)
	if err != nil {
		_ = unix.Close(fd)
	}
	return device, err
}
