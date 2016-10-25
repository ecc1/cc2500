package cc2500

import (
	"log"
	"math"
	"time"
)

const (
	verbose            = true
	fifoSize           = 64

	// Approximate time for one byte to be transmitted, based on the data rate.
	byteDuration = time.Millisecond
)

func init() {
	if verbose {
		log.SetFlags(log.Ltime | log.Lmicroseconds | log.LUTC)
	}
}

var (
	// Accumulated frequency offset from FREQEST values after good packets.
	freqOffset uint8
)

const (
	minRSSI         = math.MinInt8
	deassertPoll    = 1500 * time.Microsecond
	maxDeassertWait = 10 * time.Millisecond
)

func (r *Radio) Receive(timeout time.Duration) ([]byte, int) {
	r.changeState(SRX, STATE_RX)
	defer r.changeState(SIDLE, STATE_IDLE)
	if verbose {
		log.Printf("waiting for interrupt in %s state", r.State())
	}
	r.hw.AwaitInterrupt(timeout)
	waited := time.Duration(0)
	for r.Error() == nil && r.hw.ReadInterrupt() {
		n := r.ReadNumRxBytes()
		if verbose {
			log.Printf("  interrupt still asserted; %d bytes in FIFO", n)
		}
		if n >= fifoSize || waited >= maxDeassertWait {
			break
		}
		time.Sleep(deassertPoll)
		waited += deassertPoll
	}
	if r.Error() != nil {
		return nil, minRSSI
	}
	numBytes := int(r.ReadNumRxBytes())
	if r.hw.ReadInterrupt() {
		if verbose {
			log.Printf("interrupt still asserted after %v with %d bytes in FIFO", maxDeassertWait, numBytes)
		}
		return nil, minRSSI
	}
	data := r.hw.ReadBurst(RXFIFO, numBytes)
	// Enter IDLE state before reading FREQEST.
	// See Design Note DN015 (SWRA159).
	r.changeState(SIDLE, STATE_IDLE)
	return r.verifyPacket(data, numBytes)
}

// Check whether packet has correct length byte and valid CRC.
// Return the body of the packet (or nil if invalid) and the RSSI.
func (r *Radio) verifyPacket(data []byte, numBytes int) ([]byte, int) {
	if r.Error() != nil {
		return nil, minRSSI
	}
	if numBytes < 3 {
		if verbose {
			log.Printf("invalid %d-byte packet", numBytes)
		}
		return nil, minRSSI
	}
	lenByte := int(data[0])
	rssi := registerToRSSI(data[numBytes-2])
	status := data[numBytes-1]
	crcOK := status&(1<<7) != 0
	lqi := status &^ (1 << 7)
	packet := data[1 : numBytes-2]
	if !crcOK {
		if verbose {
			log.Printf("invalid CRC for %d-byte packet", len(packet))
		}
		return nil, rssi
	}
	if lenByte != len(packet) {
		if verbose {
			log.Printf("incorrect length byte (%d) for %d-byte packet", lenByte, len(packet))
		}
		return nil, rssi
	}
	if verbose {
		log.Printf("received %d-byte packet with LQI = %d", len(packet), lqi)
	}
	freqEst := r.hw.ReadRegister(FREQEST)
	freqOffset += freqEst
	if verbose {
		log.Printf("frequency  offset = %d Hz (%X)", registerToFrequencyOffset(freqEst), freqEst)
		log.Printf("cumulative offset = %d Hz (%X)", registerToFrequencyOffset(freqOffset), freqOffset)
	}
	r.hw.WriteRegister(FSCTRL0, freqOffset)
	r.stats.Packets.Received++
	r.stats.Bytes.Received += numBytes
	return packet, rssi
}

func (r *Radio) Send(data []byte) {
	panic("unimplemented")
}
