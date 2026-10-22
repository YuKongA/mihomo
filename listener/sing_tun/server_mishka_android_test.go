//go:build mishka

package sing_tun

import (
	"io"
	"os"
	"testing"

	tun "github.com/metacubex/sing-tun"
)

func TestVPNTUNRecreationKeepsSourceDescriptor(t *testing.T) {
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	defer writer.Close()
	for _, value := range []byte{1, 2} {
		device, err := tunNew(tun.Options{FileDescriptor: int(reader.Fd())})
		if err != nil {
			t.Fatal(err)
		}
		if err := device.Close(); err != nil {
			t.Fatal(err)
		}
		if _, err := writer.Write([]byte{value}); err != nil {
			t.Fatalf("listener close invalidated runtime fd: %v", err)
		}
		var got [1]byte
		if _, err := io.ReadFull(reader, got[:]); err != nil || got[0] != value {
			t.Fatalf("runtime fd after recreation: got=%v err=%v", got, err)
		}
	}
}
