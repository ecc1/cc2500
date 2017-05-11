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
	// Channel contains information about a channel register and frequency offset.
	Channel struct {
		number uint8 // CHANNR value
		offset uint8 // FSCTRL0 value
	}

	// Reading is a pointer to a Packet.
	Reading *Packet
)

var (
	// Channels contains the channel numbers for the listed frequencies,
	// assuming 250 kHz channel spacing.
	// The initial FSCTRL0 offsets were determined empirically.
	Channels = []Channel{
		{000, 0xBE}, // 2425 MHz
		{100, 0xBE}, // 2450 MHz
		{199, 0xBE}, // 2474.75 MHz
		{209, 0xBE}, // 2477.25 MHz
	}
)

func (r *Radio) changeChannel(i int) {
	c := Channels[i]
	if verbose {
		log.Printf("changing to channel %d", i)
		printFrequency("offset ", c.offset)
	}
	r.hw.WriteRegister(CHANNR, c.number)
	r.hw.WriteRegister(FSCTRL0, c.offset)
}

func (r *Radio) adjustFrequency(i int) {
	c := &Channels[i]
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
		for n := range Channels {
			if verbose {
				log.Printf("listening on channel %d; sync = %v", n, inSync)
			}
			r.changeChannel(n)
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
				now := time.Now()
				inSync = true
				lastReading = now.Add(-time.Duration(n) * channelInterval)
				r.adjustFrequency(n)
				p = makePacket(now, n, data, rssi)
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

// ReceiveReadings starts a goroutine to listen for incoming packets
// and returns a channel that can be used to receive them.
func (r *Radio) ReceiveReadings() <-chan Reading {
	readings := make(chan Reading, 10)
	go r.scanChannels(readings)
	return readings
}
