package main

import (
	"fmt"
	"os"

	"github.com/sparrc/go-ping"
)

func myPinger(addr string) {
	fmt.Println("Pinger started!")
	pinger, err := ping.NewPinger(addr)
	pinger.SetPrivileged(true)

	if err != nil {

		fmt.Printf("ERROR: %s\n", err.Error())

		return

	}

	pinger.OnRecv = func(pkt *ping.Packet) {

		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",

			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)

	}

	pinger.OnFinish = func(stats *ping.Statistics) {

		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)

		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",

			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)

		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",

			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)

	}

	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())

	pinger.Run()
}

func main() {
	fmt.Println("\nICMP-client succesfully started!\n")
	var addr string
	fmt.Print("Enter address: ")
	fmt.Fscan(os.Stdin, &addr)
	var n int
	fmt.Print("Enter number of threads: ")
	fmt.Fscan(os.Stdin, &n)
	i := 0
	//n := 10
	for i < n {
		go myPinger("www.google.com")
		i++
	}
	myPinger("www.google.com")
}
