package speedtest

import (
	"net"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type TraceHop struct {
	TTL      int
	Address  string
	Latency  time.Duration
	Reached  bool
}

func Traceroute(host string, maxHops int, timeout time.Duration) ([]TraceHop, error) {
	addr, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		return nil, err
	}

	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	packetConn := conn.IPv4PacketConn()

	hops := make([]TraceHop, 0, maxHops)

	for ttl := 1; ttl <= maxHops; ttl++ {
		err = packetConn.SetTTL(ttl)
		if err != nil {
			return hops, err
		}

		msg := icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   1234,
				Seq:  ttl,
				Data: []byte("wifi-traceroute"),
			},
		}

		data, err := msg.Marshal(nil)
		if err != nil {
			return hops, err
		}

		start := time.Now()

		_, err = conn.WriteTo(data, addr)
		if err != nil {
			continue
		}

		err = conn.SetReadDeadline(start.Add(timeout))
		if err != nil {
			return hops, err
		}

		buffer := make([]byte, 1500)

		n, source, err := conn.ReadFrom(buffer)
		if err != nil {
			hops = append(hops, TraceHop{
				TTL: ttl,
			})
			continue
		}

		reply, err := icmp.ParseMessage(1, buffer[:n])
		if err != nil {
			continue
		}

		hop := TraceHop{
			TTL:     ttl,
			Address: source.String(),
			Latency: time.Since(start),
		}

		if reply.Type == ipv4.ICMPTypeEchoReply {
			hop.Reached = true
		}

		hops = append(hops, hop)

		if hop.Reached {
			break
		}
	}

	return hops, nil
}