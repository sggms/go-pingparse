package parser

import (
	"testing"
	"time"
)

const (
	payloadPing1 = `PING 127.0.0.1 (127.0.0.1) 56(84) bytes of data.
64 bytes from 127.0.0.1: icmp_seq=1 ttl=64 time=0.026 ms
64 bytes from 127.0.0.1: icmp_seq=2 ttl=64 time=0.021 ms
64 bytes from 127.0.0.1: icmp_seq=3 ttl=64 time=0.031 ms

--- 127.0.0.1 ping statistics ---
3 packets transmitted, 3 received, 0% packet loss, time 10021ms
rtt min/avg/max/mdev = 0.021/0.026/0.031/0.004 ms
`
)

func TestSimplePing(t *testing.T) {
	po, err := Parse(payloadPing1)
	if err != nil {
		t.Fatal(err)
	}

	if len(po.Replies) != 3 {
		t.Error("invalid number of replies found")
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
