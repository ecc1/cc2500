package main

import (
	"log"

	"github.com/ecc1/cc2500"
)

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.LUTC)
	r := cc2500.Open().(*cc2500.Radio)
	if r.Error() != nil {
		log.Fatal(r.Error())
	}
	r.ScanChannels()
}
