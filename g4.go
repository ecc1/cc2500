package cc2500

import (
	"fmt"
	"log"
	"os"
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

	verboseG4 = false
)

type (
	// Channel contains information about a channel register and frequency offset.
	Channel struct {
		number uint8 // CHANNR value
		offset uint8 // FSCTRL0 value
	}
)

var (
	// Channels contains the channel numbers for the listed frequencies,
	// assuming 250 kHz channel spacing.
	// The initial FSCTRL0 offsets were determined empirically.
	Channels = []Channel{
		{000, 0xFD}, // 2425 MHz
		{100, 0xFD}, // 2450 MHz
		{199, 0xFD}, // 2474.75 MHz
		{209, 0xFD}, // 2477.25 MHz
	}
)

func (r *Radio) changeChannel(i int) {
	c := Channels[i]
	if verboseG4 {
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
	if verboseG4 {
		printFrequency("FREQEST", freqEst)
		printFrequency("FSCTRL0", offset)
		printFrequency("offset ", c.offset)
	}
}

func printFrequency(label string, f byte) {
	log.Printf("%s = %d Hz (%X)", label, registerToFrequencyOffset(f), f)
}

func (r *Radio) scanChannels(readings chan<- *Packet, sync bool) {
	inSync := false
	lastReading := time.Time{}
	r.Init(baseFrequency)
	for {
		waitTime := slowWait
		var p *Packet
		for n := range Channels {
			if verboseG4 {
				log.Printf("listening on channel %d; sync = %v", n, inSync)
			}
			r.changeChannel(n)
			if n == 0 && inSync {
				syncSleep(lastReading)
				waitTime = syncWait
			}
			data, rssi := r.Receive(waitTime)
			p = r.checkPacket(n, data, rssi)
			err := r.Error()
			if err != nil && err != ErrReceiveTimeout {
				log.Print(err)
			}
			r.SetError(nil)
			if !sync {
				break
			}
			if p != nil {
				inSync = true
				lastReading = p.Timestamp.Add(-time.Duration(n) * channelInterval)
				r.adjustFrequency(n)
				break
			}
			waitTime = fastWait
		}
		readings <- p
		if p == nil {
			inSync = false
		}
	}
}

func syncSleep(lastReading time.Time) {
	t := time.Now().Add(wakeupMargin)
	sleepTime := lastReading.Add(readingInterval).Sub(t)
	if sleepTime > 0 {
		if verboseG4 {
			log.Printf("sleeping for %v", sleepTime)
		}
		time.Sleep(sleepTime)
	}
}

func (r *Radio) checkPacket(channel int, data []byte, rssi int) *Packet {
	if r.Error() != nil {
		return nil
	}
	if len(data) != packetLength {
		r.SetError(fmt.Errorf("unexpected %d-byte packet: % X", len(data), data))
		return nil
	}
	pktCRC := data[packetLength-1]
	calcCRC := CRC8(data[11 : packetLength-1])
	if calcCRC != pktCRC {
		r.SetError(fmt.Errorf("computed CRC %02X but received %02X", calcCRC, pktCRC))
		return nil
	}
	p := unmarshalPacket(time.Now(), channel, data, rssi)
	if p.TransmitterID != transmitterID && transmitterID != "" {
		r.SetError(fmt.Errorf("ignoring packet from transmitter %s", p.TransmitterID))
		return nil
	}
	return p
}

// ReceiveReadings starts a goroutine to listen for incoming packets
// and returns a channel that can be used to receive them.
func (r *Radio) ReceiveReadings() <-chan *Packet {
	var sync bool
	if transmitterID == "" {
		log.Printf("receiving readings from any G4 transmitter (%s environment variable not set)", transmitterIDEnvVar)
		sync = false
	} else {
		log.Printf("receiving readings from G4 transmitter %s", transmitterID)
		sync = true
	}
	readings := make(chan *Packet, 10)
	go r.scanChannels(readings, sync)
	return readings
}

const (
	transmitterIDEnvVar = "DEXCOM_G4_XMTR_ID"
)

var (
	transmitterID = ""
)

func init() {
	transmitterID = os.Getenv(transmitterIDEnvVar)
}
