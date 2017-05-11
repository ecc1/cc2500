package cc2500

import (
	"fmt"
	"log"
	"math"
	"time"
)

const (
	verbose  = false
	fifoSize = 64
)

func init() {
	if verbose {
		log.SetFlags(log.Ltime | log.Lmicroseconds | log.LUTC)
	}
}

const (
	minRSSI      = math.MinInt8
	deassertPoll = 2 * time.Millisecond
)

// Receive listens with the given timeout for an incoming packet.
// It returns the packet and the associated RSSI.
func (r *Radio) Receive(timeout time.Duration) ([]byte, int) {
	r.changeState(SRX, STATE_RX)
	defer r.changeState(SIDLE, STATE_IDLE)
	if verbose {
		log.Printf("waiting for interrupt in %s state", r.State())
	}
	r.hw.AwaitInterrupt(timeout)
	for count := 0; r.Error() == nil && r.hw.ReadInterrupt(); count++ {
		n := r.ReadNumRXBytes()
		if verbose {
			log.Printf("  interrupt still asserted; %d bytes in FIFO", n)
		}
		if n >= fifoSize {
			break
		}
		time.Sleep(deassertPoll)
	}
	numBytes := int(r.ReadNumRXBytes())
	data := r.hw.ReadBurst(RXFIFO, numBytes)
	if r.hw.ReadInterrupt() {
		r.SetError(fmt.Errorf("interrupt still asserted with %d bytes in FIFO", numBytes))
	}
	if r.Error() != nil {
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
	if numBytes < 3 {
		r.SetError(fmt.Errorf("invalid %d-byte packet", numBytes))
		return nil, minRSSI
	}
	lenByte := int(data[0])
	rssi := registerToRSSI(data[numBytes-2])
	status := data[numBytes-1]
	crcOK := status&(1<<7) != 0
	lqi := status &^ (1 << 7)
	packet := data[1 : numBytes-2]
	if !crcOK {
		r.SetError(fmt.Errorf("invalid CRC for %d-byte packet", len(packet)))
		return nil, rssi
	}
	if lenByte != len(packet) {
		r.SetError(fmt.Errorf("incorrect length byte (%d) for %d-byte packet", lenByte, len(packet)))
		return nil, rssi
	}
	if verbose {
		log.Printf("received %d-byte packet with LQI = %d", len(packet), lqi)
	}
	r.stats.Packets.Received++
	r.stats.Bytes.Received += numBytes
	return packet, rssi
}

// Send transmits the given packet.
func (r *Radio) Send(data []byte) {
	panic("unimplemented")
}
