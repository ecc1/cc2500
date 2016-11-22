package cc2500

import (
	"fmt"
	"log"
	"math"
	"time"
)

const (
	verbose  = true
	fifoSize = 64
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
	minRSSI      = math.MinInt8
	deassertPoll = 2 * time.Millisecond
	maxWaitCount = 5
)

func (r *Radio) Receive(timeout time.Duration) ([]byte, int) {
	r.changeState(SRX, STATE_RX)
	defer r.changeState(SIDLE, STATE_IDLE)
	if verbose {
		log.Printf("waiting for interrupt in %s state", r.State())
	}
	r.hw.AwaitInterrupt(timeout)
	for count := 0; r.Error() == nil && r.hw.ReadInterrupt(); count++ {
		n := r.ReadNumRxBytes()
		if verbose {
			log.Printf("  interrupt still asserted; %d bytes in FIFO", n)
		}
		if n >= fifoSize {
			break
		}
		time.Sleep(deassertPoll)
	}
	if r.Error() != nil {
		return nil, minRSSI
	}
	numBytes := int(r.ReadNumRxBytes())
	data := r.hw.ReadBurst(RXFIFO, numBytes)
	if r.hw.ReadInterrupt() {
		r.SetError(fmt.Errorf("interrupt still asserted with %d bytes in FIFO", numBytes))
		return nil, minRSSI
	}
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
