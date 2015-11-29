// A quick & friendly hack around the go-fastping lib
package pingoo

import (
	"net"
	"time"

	"github.com/tatsushid/go-fastping"
)

type Pingoo struct {
	IPv4Address *net.IPAddr
	FastPing    *fastping.Pinger
}

type response struct {
	addr *net.IPAddr
	rtt  time.Duration
}

func New(addr string) (*Pingoo, error) {
	p := &Pingoo{}

	ra, err := net.ResolveIPAddr("ip4:icmp", addr)
	if err != nil {
		return nil, err
	}

	p.IPv4Address = ra

	p.FastPing = fastping.NewPinger()
	p.FastPing.AddIPAddr(p.IPv4Address)

	return p, nil
}

func (p *Pingoo) Ping(count int, rtt time.Duration) (int, error) {
	onRecv, onIdle := make(chan *response), make(chan bool)

	p.FastPing.OnRecv = func(addr *net.IPAddr, t time.Duration) {
		onRecv <- &response{addr: addr, rtt: t}
	}

	p.FastPing.OnIdle = func() {
		onIdle <- true
	}

	p.FastPing.MaxRTT = rtt
	p.FastPing.RunLoop()

	return p.performPing(count, onRecv, onIdle)
}

// This is pretty gross but gets the job done
func (p *Pingoo) performPing(count int, onRecv chan *response, onIdle chan bool) (int, error) {
	lastReceived := false
	iterations, received := 0, 0

loop:
	for {
		if iterations == count {
			// Max iterations reached, time to bail
			break loop
		}

		select {
		case <-onRecv:
			lastReceived = true
			iterations++
			received++
		case <-onIdle:
			if lastReceived {
				// Received response during last iteration
				lastReceived = false
				continue loop
			} else {
				// Did not receive a response
				iterations++
			}
		case <-p.FastPing.Done():
			if err := p.FastPing.Err(); err != nil {
				return 0, err
			}
			break loop
		}
	}

	return received, nil
}
