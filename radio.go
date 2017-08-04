package cc2500

import (
	"errors"
	"fmt"
	"log"
	"math"
	"time"
)

const (
	verbose      = false
	fifoSize     = 64
	minRSSI      = math.MinInt8
	deassertPoll = 2 * time.Millisecond
)

func init() {
	if verbose {
		log.SetFlags(log.Ltime | log.Lmicroseconds | log.LUTC)
	}
}

// ErrReceiveTimeout indicates that a Receive operation timed out.
var ErrReceiveTimeout = errors.New("receive timeout")

// Receive listens with the given timeout for an incoming packet.
// It returns the packet and the associated RSSI.
// Packet layout in RX FIFO:
//	0: length byte (n)
//	1..n: packet body
//	n+1: RSSI
//	n+2: CRC OK and LQI
// 2-byte CRC following packet body is checked and stripped in hardware.
func (r *Radio) Receive(timeout time.Duration) ([]byte, int) {
	r.changeState(SRX, STATE_RX)
	defer r.changeState(SIDLE, STATE_IDLE)
	if verbose {
		log.Printf("waiting for interrupt in %s state", r.State())
	}
	r.hw.AwaitInterrupt(timeout)
	for r.Error() == nil && r.hw.ReadInterrupt() {
		n := r.ReadNumRXBytes()
		if verbose {
			log.Printf("  interrupt still asserted with %d bytes in FIFO", n)
		}
		if n >= fifoSize {
			break
		}
		time.Sleep(deassertPoll)
	}
	numBytes := int(r.ReadNumRXBytes())
	if numBytes == 0 {
		r.SetError(ErrReceiveTimeout)
		return nil, minRSSI
	}
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
	if numBytes < 4 {
		r.SetError(fmt.Errorf("invalid %d-byte packet", numBytes))
		return nil, minRSSI
	}
	lenByte := int(data[0])
	rssi := registerToRSSI(data[numBytes-2])
	status := data[numBytes-1]
	crcOK := status&(1<<7) != 0
	if !crcOK {
		r.SetError(fmt.Errorf("invalid CRC: % X (RSSI %d)", data, rssi))
		return nil, rssi
	}
	lqi := status &^ (1 << 7)
	if lenByte != numBytes-3 {
		r.SetError(fmt.Errorf("incorrect length: % X (RSSI %d)", data, rssi))
		return nil, rssi
	}
	packet := data[1 : numBytes-2]
	if verbose {
		log.Printf("received packet with RSSI %d, LQI %02X: % X", rssi, lqi, packet)
	}
	r.stats.Packets.Received++
	r.stats.Bytes.Received += numBytes
	return packet, rssi
}

// Send transmits the given packet.
func (r *Radio) Send(data []byte) {
	panic("unimplemented")
}
