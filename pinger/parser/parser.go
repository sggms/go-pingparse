package parser

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	ErrNotEnoughLines       = errors.New("not enough lines")
	ErrHeaderMismatch       = errors.New("header mismatch")
	ErrUnrecognizedLine     = errors.New("unrecognized ping reply line")
	ErrMalformedStatsHeader = errors.New("malformed stats header")
	ErrMalformedStatsLine1  = errors.New("malformed stats line 1")
	ErrMalformedStatsLine2  = errors.New("malformed stats line 2")
)

var (
	headerRx         = regexp.MustCompile(`^PING (\d+\.\d+\.\d+\.\d+) \((\d+\.\d+\.\d+\.\d+)\) (\d+)\((\d+)\) bytes of data\.$`)
	lineRx           = regexp.MustCompile(`^(\d+) bytes from (\d+\.\d+\.\d+\.\d+): icmp_seq=(\d+) ttl=(\d+) time=(.*)$`)
	statsSeparatorRx = regexp.MustCompile(`^--- (\d+\.\d+\.\d+\.\d+) ping statistics ---$`)
	statsLine1       = regexp.MustCompile(`^(\d+) packets transmitted, (\d+) received, (\d+)% packet loss, time (.*)$`)
	statsLine2       = regexp.MustCompile(`^rtt min/avg/max/mdev = ([^/]+)/([^/]+)/([^/]+)/([^ ]+) (.*)$`)
	pipeNo           = regexp.MustCompile(`([^,]+), pipe (\d+)$`)
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

// Parse will parse the specified ping output and return all the information in a a PingOutput object.
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
			last = i + 2
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

		pr.Time, err = time.ParseDuration(strings.Replace(m[5], " ", "", -1))
		if err != nil {
			return nil, err
		}

		po.Replies = append(po.Replies, pr)
	}

	// parse header
	m = statsSeparatorRx.FindStringSubmatch(lines[last])
	if len(m) != 2 {
		return nil, ErrMalformedStatsHeader
	}
	po.Stats.IPAddress = m[1]

	// parse stats line 1
	last++
	m = statsLine1.FindStringSubmatch(lines[last])
	if len(m) != 5 {
		return nil, ErrMalformedStatsLine1
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

	po.Stats.Time, err = time.ParseDuration(m[4])
	if err != nil {
		return nil, err
	}

	if len(po.Replies) == 0 {
		return &po, nil
	}

	// parse stats line 2
	last++
	m = statsLine2.FindStringSubmatch(lines[last])
	if len(m) != 6 {
		return nil, ErrMalformedStatsLine2
	}

	unit := m[5]
	pm := pipeNo.FindStringSubmatch(unit)
	if len(pm) > 1 {
		unit = pm[1]
		// pipe number in pm[2] is ignored
	}

	po.Stats.RoundTrip, err = time.ParseDuration(m[1] + unit)
	if err != nil {
		return nil, err
	}
	po.Stats.Average, err = time.ParseDuration(m[2] + unit)
	if err != nil {
		return nil, err
	}
	po.Stats.Max, err = time.ParseDuration(m[3] + unit)
	if err != nil {
		return nil, err
	}
	po.Stats.MeanDeviation, err = time.ParseDuration(m[4] + unit)
	if err != nil {
		return nil, err
	}

	return &po, nil
}
