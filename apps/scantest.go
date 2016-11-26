package main

import (
	"log"
	"time"

	"github.com/ecc1/cc2500"
)

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.LUTC)
	r := cc2500.Open().(*cc2500.Radio)
	hours := time.Tick(1 * time.Hour)
	readings := r.ReceiveReadings()
	numReadings := 0
	for {
		if r.Error() != nil {
			log.Fatal(r.Error())
		}
		select {
		case <-hours:
			log.Printf("%d readings in previous hour", numReadings)
			numReadings = 0
		case p := <-readings:
			if p != nil {
				log.Printf("%+v", *p)
				numReadings++
			}
		}
	}
}
