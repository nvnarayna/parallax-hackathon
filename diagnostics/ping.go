package diagnostics

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func CheckPacketLoss(host string, count int, timeout time.Duration) (float64, error) {
	addr, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		return 0, err
	}

	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	received := 0

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
			return 0, err
		}

		_, err = conn.WriteTo(data, addr)
		if err != nil {
			continue
		}

		err = conn.SetReadDeadline(time.Now().Add(timeout))
		if err != nil {
			return 0, err
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
			received++
		}
	}

	lost := count - received
	packetLoss := float64(lost) / float64(count) * 100

	return packetLoss, nil
}

func TestPacketLoss(url string) {
	loss, err := CheckPacketLoss(
		url,
		10,
		2*time.Second,
	)

	if err != nil {
		fmt.Println("packet loss failed:", err)
		return
	}

	fmt.Printf("packet loss: %.2f%%\n", loss)
}