package main

import (
	"fmt"
	"time"

	"github.com/prometheus-community/pro-bing"
)

func NewSimplePinger(ip string) (p *probing.Pinger, err error) {

	p, err = probing.NewPinger(ip)

	p.SetPrivileged(true)

	if err != nil {

		return
	}

	p.Count = 1
	//p.Size = *size
	p.Interval = time.Second
	p.Timeout = p.Interval * 2
	//p.TTL = *ttl
	//p.InterfaceName = *iface
	//p.SetPrivileged(*privileged)
	//p.SetTrafficClass(uint8(*tclass))

	return
}

func Ping(addr string) (done <-chan *probing.Statistics, err error) {

	doneChan := make(chan *probing.Statistics)
	done = doneChan

	p, err := probing.NewPinger(addr)

	if err != nil {

		return
	}

	p.OnRecv = func(pkt *probing.Packet) {

		fmt.Printf("IP Addr: %s receive, RTT: %v\n", pkt.IPAddr, pkt.Rtt)
	}

	p.OnFinish = func(stat *probing.Statistics) {

		doneChan <- stat
	}

	p.Run()

	return
}

func main() {

	res, err := Ping("10.0.0.1")

	if err != nil {

		return
	}

	stats := <-res

	fmt.Printf("%d\n", stats.AvgRtt)

	select {}
}
