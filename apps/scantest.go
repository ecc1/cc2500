package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ecc1/cc2500"
)

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.LUTC)
	r := cc2500.Open().(*cc2500.Radio)
	for v := range r.ReceiveReadings() {
		if r.Error() != nil {
			log.Fatal(r.Error())
		}
		fmt.Printf("%v:\n", time.Now())
		for _, p := range v {
			fmt.Printf("  % X (RSSI = %d)\n", p.Body, p.RSSI)
		}
		fmt.Printf("\n")
	}
}
