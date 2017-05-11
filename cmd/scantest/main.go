package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ecc1/cc2500"
)

func main() {
	r := cc2500.Open()
	hw := r.Hardware()
	log.Printf("connected to %s radio on %s", hw.Name(), hw.Device())
	hours := time.Tick(1 * time.Hour)
	readings := r.ReceiveReadings()
	numReadings := 0
	for {
		if r.Error() != nil {
			log.Fatal(r.Error())
		}
		select {
		case <-hours:
			fmt.Printf("%d readings in previous hour\n", numReadings)
			numReadings = 0
		case reading := <-readings:
			if reading != nil {
				print(reading)
				numReadings++
			}
		}
	}
}

func print(r cc2500.Reading) {
	b, err := json.MarshalIndent(*r, "", "  ")
	if err != nil {
		fmt.Println(err)
		fmt.Printf("%+v\n", *r)
	} else {
		fmt.Println(string(b))
	}
}
