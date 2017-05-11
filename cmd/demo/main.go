package main

import (
	"log"

	"github.com/ecc1/cc2500"
)

func main() {
	r := cc2500.Open()
	if r.Error() != nil {
		log.Fatal(r.Error())
	}

	log.Printf("Resetting radio")
	r.Reset()
	r.DumpRF()

	freq := uint32(2400000000)
	log.Println("")
	log.Printf("Initializing radio to %d Hz", freq)
	r.InitRF(freq)
	r.DumpRF()

	log.Println("")
	freq += 500000
	log.Printf("Changing frequency to %d", freq)
	r.SetFrequency(freq)
	r.DumpRF()
}
