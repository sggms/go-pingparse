package parser

import (
	"fmt"
	"testing"
	"time"
)

var (
	expectedTestCases = []PingOutput{
		// 0
		PingOutput{
			Host:              `127.0.0.1`,
			ResolvedIPAddress: `127.0.0.1`,
			PayloadSize:       56,
			PayloadActualSize: 84,
			Replies: []PingReply{
				PingReply{64, `127.0.0.1`, 1, 64, 26 * time.Microsecond, "", false},
				PingReply{64, `127.0.0.1`, 2, 64, 21 * time.Microsecond, "", false},
				PingReply{64, `127.0.0.1`, 3, 64, 31 * time.Microsecond, "", false},
			},
			Stats: PingStatistics{
				IPAddress:          `127.0.0.1`,
				PacketsTransmitted: 3,
				PacketsReceived:    3,
				Time:               10021 * time.Millisecond,
				RoundTripMin:       21 * time.Microsecond,
				RoundTripMax:       31 * time.Microsecond,
				RoundTripAverage:   26 * time.Microsecond,
				RoundTripDeviation: 4 * time.Microsecond,
			},
		},
		// 1
		PingOutput{
			Host:              `172.17.0.1`,
			ResolvedIPAddress: `172.17.0.1`,
			PayloadSize:       56,
			PayloadActualSize: 84,
			Replies: []PingReply{
				PingReply{64, `172.17.0.1`, 1, 64, 98 * time.Microsecond, "", false},
				PingReply{64, `172.17.0.1`, 2, 64, 90 * time.Microsecond, "", false},
			},
			Stats: PingStatistics{
				IPAddress:          `172.17.0.1`,
				PacketsTransmitted: 3,
				PacketsReceived:    2,
				PacketLossPercent:  33,
				Time:               10292 * time.Millisecond,
				RoundTripMin:       90 * time.Microsecond,
				RoundTripMax:       98 * time.Microsecond,
				RoundTripAverage:   94 * time.Microsecond,
				RoundTripDeviation: 4 * time.Microsecond,
			},
		},
		// 2
		PingOutput{
			Host:              `172.17.0.1`,
			ResolvedIPAddress: `172.17.0.1`,
			PayloadSize:       56,
			PayloadActualSize: 84,
			Stats: PingStatistics{
				PacketsTransmitted: 3,
				PacketsReceived:    0,
				PacketLossPercent:  100,
				Time:               10205 * time.Millisecond,
			},
		},
		// 3
		PingOutput{
			Host:              `172.17.0.2`,
			ResolvedIPAddress: `172.17.0.2`,
			PayloadSize:       56,
			PayloadActualSize: 84,
			Replies: []PingReply{
				PingReply{64, `172.17.0.2`, 4, 63, 286 * time.Millisecond, "", false},
				PingReply{64, `172.17.0.2`, 5, 63, 111 * time.Millisecond, "", false},
			},
			Stats: PingStatistics{
				IPAddress:          `172.17.0.2`,
				PacketsTransmitted: 5,
				PacketsReceived:    2,
				PacketLossPercent:  60,
				Time:               4055 * time.Millisecond,
				RoundTripMin:       111409 * time.Microsecond,
				RoundTripMax:       286063 * time.Microsecond,
				RoundTripAverage:   198736 * time.Microsecond,
				RoundTripDeviation: 87327 * time.Microsecond,
			},
		},
		// 4
		PingOutput{
			Host:              `172.17.0.2`,
			ResolvedIPAddress: `172.17.0.2`,
			PayloadSize:       56,
			PayloadActualSize: 84,
			Replies: []PingReply{
				PingReply{64, `172.17.0.2`, 4, 63, 286 * time.Millisecond, "", false},
				PingReply{64, `172.17.0.2`, 5, 63, 111 * time.Millisecond, "", false},
			},
			Stats: PingStatistics{
				IPAddress:          `172.17.0.2`,
				PacketsTransmitted: 5,
				PacketsReceived:    2,
				PacketLossPercent:  60,
				Time:               4055 * time.Millisecond,
				RoundTripMin:       111409 * time.Microsecond,
				RoundTripMax:       286063 * time.Microsecond,
				RoundTripAverage:   198736 * time.Microsecond,
				RoundTripDeviation: 87327 * time.Microsecond,
			},
		},
		// 5
		PingOutput{
			Host:              `172.17.0.3`,
			ResolvedIPAddress: `172.17.0.3`,
			PayloadSize:       56,
			PayloadActualSize: 84,
			Replies: []PingReply{
				PingReply{0, ``, 2, 0, 0, "Destination Host Unreachable", false},
			},
			Stats: PingStatistics{
				IPAddress:          `172.17.0.3`,
				Errors:             1,
				PacketsTransmitted: 4,
				PacketsReceived:    0,
				PacketLossPercent:  100,
				Time:               3055 * time.Millisecond,
			},
		},
		// 6
		PingOutput{
			Host:              `127.0.0.1`,
			ResolvedIPAddress: `127.0.0.1`,
			PayloadSize:       56,
			Replies: []PingReply{
				PingReply{64, `127.0.0.1`, 0, 64, 61 * time.Microsecond, "", false},
				PingReply{64, `127.0.0.1`, 1, 64, 57 * time.Microsecond, "", false},
				PingReply{64, `127.0.0.1`, 2, 64, 108 * time.Microsecond, "", false},
			},
			Stats: PingStatistics{
				IPAddress:          `127.0.0.1`,
				PacketsTransmitted: 3,
				PacketsReceived:    3,
				RoundTripMin:       57 * time.Microsecond,
				RoundTripMax:       108 * time.Microsecond,
				RoundTripAverage:   75 * time.Microsecond,
				RoundTripDeviation: 23 * time.Microsecond,
			},
		},
		// 7
		PingOutput{
			Host:              `172.17.0.4`,
			ResolvedIPAddress: `172.17.0.4`,
			PayloadSize:       56,
			Stats: PingStatistics{
				IPAddress:          `172.17.0.4`,
				PacketsTransmitted: 6,
				PacketsReceived:    0,
				PacketLossPercent:  100,
			},
		},
	}
	payloads = []string{
		// 0
		`PING 127.0.0.1 (127.0.0.1) 56(84) bytes of data.
64 bytes from 127.0.0.1: icmp_seq=1 ttl=64 time=0.026 ms
64 bytes from 127.0.0.1: icmp_seq=2 ttl=64 time=0.021 ms
64 bytes from 127.0.0.1: icmp_seq=3 ttl=64 time=0.031 ms

--- 127.0.0.1 ping statistics ---
3 packets transmitted, 3 received, 0% packet loss, time 10021ms
rtt min/avg/max/mdev = 0.021/0.026/0.031/0.004 ms
`,
		// 1
		`PING 172.17.0.1 (172.17.0.1) 56(84) bytes of data.
64 bytes from 172.17.0.1: icmp_seq=1 ttl=64 time=0.098 ms
64 bytes from 172.17.0.1: icmp_seq=2 ttl=64 time=0.090 ms

--- 172.17.0.1 ping statistics ---
3 packets transmitted, 2 received, 33% packet loss, time 10292ms
rtt min/avg/max/mdev = 0.090/0.094/0.098/0.004 ms
`,
		// 2
		`PING 172.17.0.1 (172.17.0.1) 56(84) bytes of data.

--- 172.17.0.1 ping statistics ---
3 packets transmitted, 0 received, 100% packet loss, time 10205ms
`,
		// 3
		`PING 172.17.0.2 (172.17.0.2) 56(84) bytes of data.
64 bytes from 172.17.0.2: icmp_seq=4 ttl=63 time=286 ms
64 bytes from 172.17.0.2: icmp_seq=5 ttl=63 time=111 ms

--- 172.17.0.2 ping statistics ---
5 packets transmitted, 2 received, 60% packet loss, time 4055ms
rtt min/avg/max/mdev = 111.409/198.736/286.063/87.327 ms
`,
		// 4
		`PING 172.17.0.2 (172.17.0.2) 56(84) bytes of data.
64 bytes from 172.17.0.2: icmp_seq=4 ttl=63 time=286 ms
64 bytes from 172.17.0.2: icmp_seq=5 ttl=63 time=111 ms

--- 172.17.0.2 ping statistics ---
5 packets transmitted, 2 received, 60% packet loss, time 4055ms
rtt min/avg/max/mdev = 111.409/198.736/286.063/87.327 ms, pipe 2
`,
		// 5
		`PING 172.17.0.3 (172.17.0.3) 56(84) bytes of data.
From 93.184.216.34 icmp_seq=2 Destination Host Unreachable

--- 172.17.0.3 ping statistics ---
4 packets transmitted, 0 received, +1 errors, 100% packet loss, time 3055ms
pipe 3
`,
		// 6
		`PING 127.0.0.1 (127.0.0.1): 56 data bytes, id 0x0001 = 1
64 bytes from 127.0.0.1: icmp_seq=0 ttl=64 time=0.061 ms
64 bytes from 127.0.0.1: icmp_seq=1 ttl=64 time=0.057 ms
64 bytes from 127.0.0.1: icmp_seq=2 ttl=64 time=0.108 ms
--- 127.0.0.1 ping statistics ---
3 packets transmitted, 3 packets received, 0% packet loss
round-trip min/avg/max/stddev = 0.057/0.075/0.108/0.023 ms
`,
		// 7
		`PING 172.17.0.4 (172.17.0.4): 56 data bytes, id 0x05e1 = 1505
--- 172.17.0.4 ping statistics ---
6 packets transmitted, 0 packets received, 100% packet loss
`,
	}
)

func TestPings(t *testing.T) {
	for i := 0; i < len(payloads); i++ {
		// capture range variables
		payload := payloads[i]
		expected := expectedTestCases[i]
		t.Run(fmt.Sprintf("payload #%d", i), func(t *testing.T) {
			t.Parallel()

			po, err := Parse(payload)
			if err != nil {
				t.Fatal(err)
			}

			if po.Host != expected.Host {
				t.Errorf("expected host %q, but got %q", expected.Host, po.Host)
			}

			if po.ResolvedIPAddress != expected.ResolvedIPAddress {
				t.Errorf("expected resolved IP address %q, but got %q", expected.ResolvedIPAddress, po.ResolvedIPAddress)
			}

			if po.PayloadSize != expected.PayloadSize {
				t.Errorf("expected payload size %v, but got %v", expected.PayloadSize, po.PayloadSize)
			}
			if po.PayloadActualSize != expected.PayloadActualSize {
				t.Errorf("expected payload actual size %v, but got %v", expected.PayloadActualSize, po.PayloadActualSize)
			}

			///
			/// check stats fields
			///

			if po.Stats.IPAddress != expected.Stats.IPAddress {
				t.Errorf("expected stats IP address %q, but got %q", expected.Stats.IPAddress, po.Stats.IPAddress)
			}

			if po.Stats.Errors != expected.Stats.Errors {
				t.Errorf("expected errors %v, but got %v", expected.Stats.Errors, po.Stats.Errors)
			}
			if po.Stats.PacketLossPercent != expected.Stats.PacketLossPercent {
				t.Errorf("expected packet loss percent %v, but got %v", expected.Stats.PacketLossPercent, po.Stats.PacketLossPercent)
			}
			if po.Stats.Time != expected.Stats.Time {
				t.Errorf("expected time %v, but got %v", expected.Stats.Time, po.Stats.Time)
			}

			if po.Stats.RoundTripMin != expected.Stats.RoundTripMin {
				t.Errorf("expected rtt min %v, but got %v", expected.Stats.RoundTripMin, po.Stats.RoundTripMin)
			}
			if po.Stats.RoundTripMax != expected.Stats.RoundTripMax {
				t.Errorf("expected rtt max %v, but got %v", expected.Stats.RoundTripMax, po.Stats.RoundTripMax)
			}
			if po.Stats.RoundTripAverage != expected.Stats.RoundTripAverage {
				t.Errorf("expected rtt avg %v, but got %v", expected.Stats.RoundTripAverage, po.Stats.RoundTripAverage)
			}
			if po.Stats.RoundTripDeviation != expected.Stats.RoundTripDeviation {
				t.Errorf("expected rtt mdev %v, but got %v", expected.Stats.RoundTripDeviation, po.Stats.RoundTripDeviation)
			}

			if po.Stats.PacketsReceived != expected.Stats.PacketsReceived {
				t.Errorf("expected packets received %v, but got %v", expected.Stats.PacketsReceived, po.Stats.PacketsReceived)
			}
			if po.Stats.PacketsTransmitted != expected.Stats.PacketsTransmitted {
				t.Errorf("expected packets transmitted %v, but got %v", expected.Stats.PacketsTransmitted, po.Stats.PacketsTransmitted)
			}

			if len(expected.Replies) != len(po.Replies) {
				t.Errorf("expected %d replies, but got %d", len(expected.Replies), len(po.Replies))
			}

			for i := 0; i < len(po.Replies); i++ {
				pr := po.Replies[i]
				epr := expected.Replies[i]
				if epr.SequenceNumber != pr.SequenceNumber {
					t.Errorf("reply %d: expected sequence number %d, but got %d", i, epr.SequenceNumber, pr.SequenceNumber)
				}

				if epr.Time != pr.Time {
					t.Errorf("reply %d: expected time %v, but got %v", i, epr.Time, pr.Time)
				}

				if epr.Size != pr.Size {
					t.Errorf("reply %d: expected size %v, but got %v", i, epr.Size, pr.Size)
				}

				if epr.FromAddress != pr.FromAddress {
					t.Errorf("reply %d: expected from address %v, but got %v", i, epr.FromAddress, pr.FromAddress)
				}
				if epr.SequenceNumber != pr.SequenceNumber {
					t.Errorf("reply %d: expected sequence number %v, but got %v", i, epr.SequenceNumber, pr.SequenceNumber)
				}
				if epr.TTL != pr.TTL {
					t.Errorf("reply %d: expected TTL %v, but got %v", i, epr.TTL, pr.TTL)
				}
				if epr.Error != pr.Error {
					t.Errorf("reply %d: expected error %q, but got %q", i, epr.Error, pr.Error)
				}
				if epr.Duplicate != pr.Duplicate {
					t.Errorf("reply %d: expected duplicate %v, but got %v", i, epr.Duplicate, pr.Duplicate)
				}
			}

		})
	}

}
