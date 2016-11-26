package cc2500

import (
	"log"
	"time"
)

const (
	baseFrequency = 2425000000

	readingInterval = 5 * time.Minute
	channelInterval = 500 * time.Millisecond

	wakeupMargin = 100 * time.Millisecond

	slowWait = readingInterval + 1*time.Minute
	fastWait = channelInterval + 50*time.Millisecond
	syncWait = wakeupMargin + 100*time.Millisecond
)

type (
	Channel struct {
		number uint8 // CHANNR value
		offset uint8 // FSCTRL0 value
	}

	Reading *Packet
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
	inSync := false
	lastReading := time.Time{}
	r.Init(baseFrequency)
	for {
		waitTime := slowWait
		p := (*Packet)(nil)
		for n, c := range channels {
			if verbose {
				log.Printf("listening on channel %d; sync = %v", n, inSync)
			}
			r.changeChannel(c)
			if n == 0 && inSync {
				t := time.Now().Add(wakeupMargin)
				sleepTime := lastReading.Add(readingInterval).Sub(t)
				if sleepTime > 0 {
					if verbose {
						log.Printf("sleeping for %v in %s state", sleepTime, r.State())
					}
					time.Sleep(sleepTime)
				}
				waitTime = syncWait
			}
			data, rssi := r.Receive(waitTime)
			if r.Error() == nil {
				inSync = true
				lastReading = time.Now().Add(-time.Duration(n) * channelInterval)
				r.adjustFrequency(c)
				p = makePacket(data, rssi)
				break
			}
			log.Print(r.Error())
			r.SetError(nil)
			waitTime = fastWait
		}
		readings <- p
		if p == nil {
			inSync = false
		}
	}
}

func (r *Radio) ReceiveReadings() <-chan Reading {
	readings := make(chan Reading, 10)
	go r.scanChannels(readings)
	return readings
}
