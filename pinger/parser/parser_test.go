package parser

import (
	"testing"
	"time"
)

const (
	payload1 = `PING 127.0.0.1 (127.0.0.1) 56(84) bytes of data.
64 bytes from 127.0.0.1: icmp_seq=1 ttl=64 time=0.026 ms
64 bytes from 127.0.0.1: icmp_seq=2 ttl=64 time=0.021 ms
64 bytes from 127.0.0.1: icmp_seq=3 ttl=64 time=0.031 ms

--- 127.0.0.1 ping statistics ---
3 packets transmitted, 3 received, 0% packet loss, time 10021ms
rtt min/avg/max/mdev = 0.021/0.026/0.031/0.004 ms
`
	payload2 = `PING 172.17.0.1 (172.17.0.1) 56(84) bytes of data.
64 bytes from 172.17.0.1: icmp_seq=1 ttl=64 time=0.098 ms
64 bytes from 172.17.0.1: icmp_seq=2 ttl=64 time=0.090 ms

--- 172.17.0.1 ping statistics ---
3 packets transmitted, 2 received, 33% packet loss, time 10292ms
rtt min/avg/max/mdev = 0.090/0.094/0.098/0.004 ms
`
	payload3 = `PING 172.17.0.1 (172.17.0.1) 56(84) bytes of data.

--- 172.17.0.1 ping statistics ---
3 packets transmitted, 0 received, 100% packet loss, time 10205ms
`
	payload4 = `PING 172.17.0.2 (172.17.0.2) 56(84) bytes of data.
64 bytes from 172.17.0.2: icmp_seq=4 ttl=63 time=286 ms
64 bytes from 172.17.0.2: icmp_seq=5 ttl=63 time=111 ms

--- 172.17.0.2 ping statistics ---
5 packets transmitted, 2 received, 60% packet loss, time 4055ms
rtt min/avg/max/mdev = 111.409/198.736/286.063/87.327 ms
`
)

func TestSimplePing(t *testing.T) {
	po, err := Parse(payload1)
	if err != nil {
		t.Fatal(err)
	}

	if len(po.Replies) != 3 {
		t.Fatal("invalid number of replies found")
	}

	if po.Replies[1].SequenceNumber != 2 {
		t.Error("invalid sequence number found")
	}

	if po.Replies[2].Time != time.Duration(31)*time.Microsecond {
		t.Error("invalid time of reply 2 found")
	}

	// check some fields
	if po.Stats.RoundTrip != time.Duration(21)*time.Microsecond {
		t.Error("invalid RTT found")
	}
	if po.Stats.Average != time.Duration(26)*time.Microsecond {
		t.Error("invalid avg found")
	}
	if po.Stats.Max != time.Duration(31)*time.Microsecond {
		t.Error("invalid max found")
	}
	if po.Stats.MeanDeviation != time.Duration(4)*time.Microsecond {
		t.Error("invalid max found")
	}
}

func TestShakyPing(t *testing.T) {
	po, err := Parse(payload2)
	if err != nil {
		t.Fatal(err)
	}

	if len(po.Replies) != 2 {
		t.Fatal("invalid number of replies found")
	}

	if po.Replies[1].SequenceNumber != 2 {
		t.Error("invalid sequence number found")
	}

	if po.Replies[0].Time != time.Duration(98)*time.Microsecond {
		t.Error("invalid time of reply 2 found")
	}

	// check some fields
	if po.Stats.RoundTrip != time.Duration(90)*time.Microsecond {
		t.Error("invalid RTT found")
	}
	if po.Stats.Average != time.Duration(94)*time.Microsecond {
		t.Error("invalid avg found")
	}
	if po.Stats.Max != time.Duration(98)*time.Microsecond {
		t.Error("invalid max found")
	}
	if po.Stats.MeanDeviation != time.Duration(4)*time.Microsecond {
		t.Error("invalid max found")
	}
}

func TestFailedPing(t *testing.T) {
	po, err := Parse(payload3)
	if err != nil {
		t.Fatal(err)
	}

	if len(po.Replies) != 0 {
		t.Fatal("invalid number of replies found")
	}

	// check some fields
	if po.Stats.PacketLossPercent != 100 {
		t.Error("invalid packet loss percentage found")
	}
	if po.Stats.PacketsTransmitted != 3 {
		t.Error("invalid packets transmitted found")
	}
	if po.Stats.Time != time.Duration(10205)*time.Millisecond {
		t.Error("invalid stats time found")
	}
}

func TestPingDifferent(t *testing.T) {
	po, err := Parse(payload4)
	if err != nil {
		t.Fatal(err)
	}

	if len(po.Replies) != 2 {
		t.Fatal("invalid number of replies found")
	}

	if po.Replies[1].SequenceNumber != 5 {
		t.Error("invalid sequence number found")
	}

	if po.Replies[0].Time != time.Duration(286)*time.Millisecond {
		t.Error("invalid time of reply found")
	}
}
