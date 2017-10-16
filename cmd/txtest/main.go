package main

import (
	"flag"
	"log"
	"time"

	"github.com/ecc1/cc2500"
)

var (
	count            = flag.Int("n", 0, "send only `count` packets")
	minPacketSize    = flag.Int("min", 1, "minimum packet `size` in bytes")
	maxPacketSize    = flag.Int("max", 30, "maximum packet `size` in bytes")
	frequency        = flag.Uint("f", 2444000000, "frequency in Hz")
	interPacketDelay = flag.Duration("delay", time.Second, "inter-packet delay")
)

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.LUTC)
	flag.Parse()
	r := cc2500.Open()
	if r.Error() != nil {
		log.Fatal(r.Error())
	}
	defer r.Close()
	log.Printf("setting frequency to %d", *frequency)
	r.Init(uint32(*frequency))
	log.Printf("actual frequency: %d", r.Frequency())

	n := *minPacketSize
	pkts := 0
	data := make([]byte, *maxPacketSize)
	for r.Error() == nil {
		if *count != 0 && pkts == *count {
			return
		}
		for i := 0; i < n; i++ {
			data[i] = byte(i + 1)
		}
		packet := data[:n]
		log.Printf("data: % X", packet)
		r.Send(packet)
		pkts++
		n++
		if n > *maxPacketSize {
			n = *minPacketSize
		}
		time.Sleep(*interPacketDelay)
	}
	log.Fatal(r.Error())
}
