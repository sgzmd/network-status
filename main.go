/*
This is a simple ping utility that uses the go-fastping library to ping a list of hosts and returns the average latency
in milliseconds. It can be run as:

sudo go run main.go -count 5 -max-latency 150 google.com yahoo.com 8.8.8.8

It will check that the hosts are reachable and that the average latency is less than 150 ms.
If the average latency is more than 150 ms, or some hosts are unreachable it will exit with a non-zero exit code.
*/

package main

import (
	"flag"
	"fmt"
	"github.com/tatsushid/go-fastping"
	"log"
	"net"
	"os"
	"time"
)

// PingHost sends an ICMP request to the host and returns avg latency in milliseconds or error
// the host could not be reached.
func PingHost(hosts []string, count int, sourceIP string) (time.Duration, error) {
	p := fastping.NewPinger()

	var totalRtt int64
	totalRtt = 0

	responses := make(map[string]int)
	for _, host := range hosts {
		ra, err := net.ResolveIPAddr("ip4:icmp", host)
		log.Printf("Resolving host: %s -> %+v", host, ra)
		if err != nil {
			return 0, err
		}
		p.AddIPAddr(ra)
		responses[fmt.Sprintf("%s", ra)] = 0
	}

	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		log.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
		ipstr := fmt.Sprintf("%s", addr)
		totalRtt += rtt.Milliseconds()
		responses[ipstr] = responses[ipstr] + 1
	}

	p.OnIdle = func() {
		log.Printf("Finish")
	}

	if sourceIP != "" {
		p.Source(sourceIP)
	}

	for i := 0; i < count; i++ {
		err := p.Run()
		if err != nil {
			return 0, err
		}
	}

	for ip, v := range responses {
		if v == 0 {
			return 0, fmt.Errorf("could not reach %+v", ip)
		}
		if v != count {
			return 0, fmt.Errorf("reached %+v only %d times out of %d requested", ip, v, count)
		}
	}

	avgRtt := totalRtt / int64(len(hosts)*count)
	return time.Duration(avgRtt) * time.Millisecond, nil
}

func main() {
	count := flag.Int("count", 5, "Number of pings to send")
	maxLatency := flag.Int("max-latency", 150, "Max latency in milliseconds")
	sourceIP := flag.String("source-ip", "", "Source IP address (default: system default)")

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(-1)
	}

	hosts := flag.Args()

	d, e := PingHost(hosts, *count, *sourceIP)
	if e != nil {
		log.Fatal(e)
	} else {
		if d.Milliseconds() > int64(*maxLatency) {
			log.Printf("Average ping time: %v ms, which is more than %d ms", d.Milliseconds(), *maxLatency)
			os.Exit(1)
		}

		log.Printf("Average ping time: %v ms\n", d.Milliseconds())
		os.Exit(0)
	}
}
