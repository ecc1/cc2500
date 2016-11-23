package cc2500

import (
	"log"
	"time"
)

const (
	baseFrequency = 2425000000
	slowWait      = 6 * time.Minute
	fastWait      = 550 * time.Millisecond
)

type (
	Channel struct {
		number uint8 // CHANNR value
		offset uint8 // FSCTRL0 value
	}

	Packet struct {
		Body []byte
		RSSI int
	}

	Reading []Packet
)

var (
	// With 250 kHz channel spacing, these channel numbers
	// correspond to the frequencies below.  The initial
	// FSCTRL0 offsets were determined empirically.
	channels = []Channel{
		{000, 0xBE}, // 2425 MHz
		{100, 0xBE}, // 2450 MHz
		{199, 0xBE}, // 2474.75 MHz
		{209, 0xBE}, // 2477.25 MHz
	}
)

func (r *Radio) changeChannel(c Channel) {
	r.hw.WriteRegister(CHANNR, c.number)
	r.hw.WriteRegister(FSCTRL0, c.offset)
}

func (r *Radio) adjustFrequency(c Channel) {
	freqEst := r.hw.ReadRegister(FREQEST)
	offset := r.hw.ReadRegister(FSCTRL0)
	c.offset = offset + freqEst
	r.hw.WriteRegister(FSCTRL0, c.offset)
	if verbose {
		printFrequency("FREQEST", freqEst)
		printFrequency("FSCTRL0", offset)
		printFrequency("offset ", c.offset)
	}
}

func printFrequency(label string, f byte) {
	log.Printf("%s = %d Hz (%X)", label, registerToFrequencyOffset(f), f)
}

func (r *Radio) scanChannels(readings chan<- Reading) {
	r.Init(baseFrequency)
	for {
		waitTime := slowWait
		v := []Packet{}
		for n, c := range channels {
			if verbose {
				log.Printf("listening on channel %d", n)
			}
			r.changeChannel(c)
			data, rssi := r.Receive(waitTime)
			if r.Error() != nil {
				log.Print(r.Error())
				r.SetError(nil)
				continue
			}
			r.adjustFrequency(c)
			v = append(v, Packet{Body: data, RSSI: rssi})
			waitTime = fastWait
		}
		readings <- v
	}
}

func (r *Radio) ReceiveReadings() <-chan Reading {
	readings := make(chan Reading, 10)
	go r.scanChannels(readings)
	return readings
}
