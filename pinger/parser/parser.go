package parser

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	ErrNotEnoughLines        = errors.New("not enough lines")
	ErrHeaderMismatch        = errors.New("header mismatch")
	ErrUnexpectedFloatFormat = errors.New("unexpected float format")
	ErrUnrecognizedLine      = errors.New("unrecognized ping reply line")
	ErrMalformedStats        = errors.New("malformed stats")
)

var (
	headerRx         = regexp.MustCompile(`^PING (\d+\.\d+\.\d+\.\d+) \((\d+\.\d+\.\d+\.\d+)\) (\d+)\((\d+)\) bytes of data\.$`)
	lineRx           = regexp.MustCompile(`^(\d+) bytes from (\d+\.\d+\.\d+\.\d+): icmp_seq=(\d+) ttl=(\d+) time=(\d+\.\d+) ms$`)
	statsSeparatorRx = regexp.MustCompile(`^--- (\d+\.\d+\.\d+\.\d+) ping statistics ---$`)
	statsLine1       = regexp.MustCompile(`^(\d+) packets transmitted, (\d+) received, (\d+)% packet loss, time (\d+)ms$`)
	statsLine2       = regexp.MustCompile(`^rtt min/avg/max/mdev = (\d+\.\d+)/(\d+\.\d+)/(\d+\.\d+)/(\d+\.\d+) ms$`)
)

// PingOutput contains the whole ping operation output.
type PingOutput struct {
	IPAddress         string
	ResolvedIPAddress string
	PayloadSize       uint
	PayloadActualSize uint
	Replies           []PingReply
	Stats             PingStatistics
}

// PingReply contains an individual ping reply line.
type PingReply struct {
	Size           uint
	FromAddress    string
	SequenceNumber uint
	TTL            uint
	Time           time.Duration
}

// PingStatistics contains the statistics of the whole ping operation.
type PingStatistics struct {
	IPAddress          string
	PacketsTransmitted uint
	PacketsReceived    uint
	PacketLossPercent  uint8
	Time               time.Duration
	RoundTrip          time.Duration
	Average            time.Duration
	Max                time.Duration
	MeanDeviation      time.Duration
}

func Parse(s string) (*PingOutput, error) {
	var po PingOutput

	// separate full output text into lines
	lines := strings.Split(s, "\n")
	if len(lines) < 5 {
		return nil, ErrNotEnoughLines
	}

	m := headerRx.FindStringSubmatch(lines[0])
	if len(m) != 5 {
		return nil, ErrHeaderMismatch
	}
	po.IPAddress = m[1]
	po.ResolvedIPAddress = m[2]
	payloadSize, err := strconv.ParseUint(m[3], 10, 64)
	if err != nil {
		return nil, err
	}
	po.PayloadSize = uint(payloadSize)
	payloadActualSize, err := strconv.ParseUint(m[4], 10, 64)
	if err != nil {
		return nil, err
	}
	po.PayloadActualSize = uint(payloadActualSize)

	var last int
	for i, line := range lines[1:] {
		if line == "" {
			last = i
			break
		}
		m := lineRx.FindStringSubmatch(line)
		if len(m) != 6 {
			return nil, ErrUnrecognizedLine
		}

		var pr PingReply

		replySize, err := strconv.ParseUint(m[1], 10, 64)
		if err != nil {
			return nil, err
		}
		pr.Size = uint(replySize)
		pr.FromAddress = m[2]
		replySeqNo, err := strconv.ParseUint(m[3], 10, 64)
		if err != nil {
			return nil, err
		}
		pr.SequenceNumber = uint(replySeqNo)
		replyTTL, err := strconv.ParseUint(m[4], 10, 64)
		if err != nil {
			return nil, err
		}
		pr.TTL = uint(replyTTL)

		pr.Time, err = parseMs(m[5])
		if err != nil {
			return nil, err
		}

		po.Replies = append(po.Replies, pr)
	}

	header := lines[last+2]
	m = statsSeparatorRx.FindStringSubmatch(header)
	if len(m) != 2 {
		return nil, ErrMalformedStats
	}
	po.Stats.IPAddress = m[1]

	m = statsLine1.FindStringSubmatch(lines[last+3])
	if len(m) != 5 {
		return nil, ErrMalformedStats
	}
	packetsTransmitted, err := strconv.ParseUint(m[1], 10, 64)
	if err != nil {
		return nil, err
	}
	po.Stats.PacketsTransmitted = uint(packetsTransmitted)

	packetsReceived, err := strconv.ParseUint(m[2], 10, 64)
	if err != nil {
		return nil, err
	}
	po.Stats.PacketsReceived = uint(packetsReceived)

	packetLossPcent, err := strconv.ParseUint(m[3], 10, 64)
	if err != nil {
		return nil, err
	}
	po.Stats.PacketLossPercent = uint8(packetLossPcent)

	t, err := strconv.ParseUint(m[4], 10, 64)
	if err != nil {
		return nil, err
	}
	po.Stats.Time = time.Millisecond * time.Duration(t)

	m = statsLine2.FindStringSubmatch(lines[last+4])
	if len(m) != 5 {
		return nil, ErrMalformedStats
	}

	po.Stats.RoundTrip, err = parseMs(m[1])
	if err != nil {
		return nil, err
	}
	po.Stats.Average, err = parseMs(m[2])
	if err != nil {
		return nil, err
	}
	po.Stats.Max, err = parseMs(m[3])
	if err != nil {
		return nil, err
	}
	po.Stats.MeanDeviation, err = parseMs(m[4])
	if err != nil {
		return nil, err
	}

	return &po, nil
}

func parseMs(s string) (time.Duration, error) {
	if len(s) != 5 {
		return 0, ErrUnexpectedFloatFormat
	}

	// remove dot and leading zeroes
	s = strings.TrimLeft(strings.Replace(s, ".", "", -1), "0")
	us, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return time.Microsecond * time.Duration(us), nil
}
