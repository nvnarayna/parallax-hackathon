package diagnostics

import (
	"net"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type PingResponse struct {
	Sequence int
	Latency  time.Duration
}

func Ping(host string, count int, timeout time.Duration) ([]PingResponse, error) {
	addr, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		return nil, err
	}

	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	responses := make([]PingResponse, 0, count)

	for i := 0; i < count; i++ {
		msg := icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   1234,
				Seq:  i,
				Data: []byte("wifi-diagnostic"),
			},
		}

		data, err := msg.Marshal(nil)
		if err != nil {
			return responses, err
		}

		start := time.Now()

		_, err = conn.WriteTo(data, addr)
		if err != nil {
			continue
		}

		err = conn.SetReadDeadline(start.Add(timeout))
		if err != nil {
			return responses, err
		}

		buffer := make([]byte, 1500)

		n, _, err := conn.ReadFrom(buffer)
		if err != nil {
			continue
		}

		reply, err := icmp.ParseMessage(1, buffer[:n])
		if err != nil {
			continue
		}

		if reply.Type == ipv4.ICMPTypeEchoReply {
			responses = append(responses, PingResponse{
				Sequence: i,
				Latency:  time.Since(start),
			})
		}
	}

	return responses, nil
}