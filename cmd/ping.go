package cmd

import (
	"fmt"
	"os"
	"os/signal"
	// "net"
	"time"

	
	"github.com/go-ping/ping"
	"github.com/urfave/cli"

)

var Ping = cli.Command{
	Name:        "ping",
	Usage:       "This command ping ip or domain",
	Description: `ping ip`,
	Action:      PingRun,
	Flags: []cli.Flag{
		stringFlag("ip, i", "127.0.0.1", "ip or domain address"),
		intFlag("timeout, t", 3, "timeout"),
	},
}

func PingRun(c *cli.Context) error {

	ip := c.String("ip")
	timeout := c.Int64("timeout")

	pinger, err := ping.NewPinger(ip)
	if err != nil {
		panic(err)
	}

	signal_os := make(chan os.Signal, 1)
	signal.Notify(signal_os, os.Interrupt)
	go func() {
		for _ = range signal_os {
			pinger.Stop()
		}
	}()

	pinger.OnRecv = func(pkt *ping.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d ttl=%v time=%v \n", pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Ttl, pkt.Rtt)
	}

	pinger.OnDuplicateRecv = func(pkt *ping.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v (DUP!)\n", pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.Ttl)
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	pinger.Size = 56
	pinger.Timeout = time.Duration(timeout) * time.Second
	err = pinger.Run()
	if err != nil {
		panic(err)
	}
	return nil
}



