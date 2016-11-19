package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ecc1/cc2500"
)

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.LUTC)
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s frequency", os.Args[0])
	}
	frequency := getFrequency(os.Args[1])
	r := cc2500.Open().(*cc2500.Radio)
	if r.Error() != nil {
		log.Fatal(r.Error())
	}
	log.Printf("setting frequency to %d", frequency)
	r.Init(frequency)
	log.Printf("actual frequency: %d", r.Frequency())
	for r.Error() == nil {
		data, rssi := r.Receive(time.Hour)
		log.Printf("% X (RSSI = %d)", data, rssi)
	}
	log.Fatal(r.Error())
}

func getFrequency(s string) uint32 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatal(err)
	}
	if 2.4 <= f && f <= 2.5 {
		return uint32(f * 1e9)
	}
	if 2400 <= f && f <= 2500 {
		return uint32(f * 1e6)
	}
	if 2400000000 <= f && f <= 2500000000 {
		return uint32(f)
	}
	log.Fatalf("%s: invalid pump frequency", s)
	panic("unreachable")
}
