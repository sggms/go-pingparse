package pinger

import (
	"testing"
	"time"
)

func TestPingLocalhost(t *testing.T) {
	_, err := Ping("127.0.0.1", time.Second, time.Second*3, 56)
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestPingInvalidHost(t *testing.T) {
	po, err := Ping("128.0.0.0", time.Second, time.Second*3, 56)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if po.Stats.PacketLossPercent != 0x64 {
		t.Error("invalid packet loss percent found")
	}
}
