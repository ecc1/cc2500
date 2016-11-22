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

type Channel struct {
	number uint8 // CHANNR value
	offset uint8 // FSCTRL0 value
}

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

func (r *Radio) ChangeChannel(c Channel) {
	r.hw.WriteRegister(CHANNR, c.number)
	r.hw.WriteRegister(FSCTRL0, c.offset)
}

func (r *Radio) ScanChannels() {
	r.Init(baseFrequency)
	for {
		waitTime := slowWait
		for n, c := range channels {
			log.Printf("listening on channel %d", n)
			r.ChangeChannel(c)
			data, rssi := r.Receive(waitTime)
			if r.Error() != nil {
				log.Print(r.Error())
				r.SetError(nil)
				continue
			}
			log.Printf("% X (RSSI = %d)", data, rssi)
			r.adjustFrequency(c)
			waitTime = fastWait
		}
	}
}

func (r *Radio) adjustFrequency(c Channel) {
	freqEst := r.hw.ReadRegister(FREQEST)
	offset := r.hw.ReadRegister(FSCTRL0)
	c.offset = offset + freqEst
	r.hw.WriteRegister(FSCTRL0, c.offset)
	if verbose {
		log.Printf("FREQEST = %d Hz (%X)", registerToFrequencyOffset(freqEst), freqEst)
		log.Printf("FSCTRL0 = %d Hz (%X)", registerToFrequencyOffset(offset), offset)
		log.Printf("offset  = %d Hz (%X)", registerToFrequencyOffset(c.offset), c.offset)
	}
}
