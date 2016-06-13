package cc1101

import (
	"bytes"
	"log"
	"time"
)

const (
	readFifoUsingBurst = true
	fifoSize           = 64
	maxPacketSize      = 110

	// Approximate time for one byte to be transmitted, based on
	// the data rate.  It was determined empirically so that few
	// if any iterations are needed in drainTxFifo().
	byteDuration = time.Millisecond
)

func (r *Radio) Send(data []byte) {
	if len(data) > maxPacketSize {
		log.Panicf("attempting to send %d-byte packet", len(data))
	}
	if r.Error() != nil {
		return
	}
	if verbose {
		log.Printf("sending %d-byte packet in %s state", len(data), r.State())
	}
	// Terminate packet with zero byte,
	// and pad with another to ensure final bytes
	// are transmitted before leaving TX state.
	packet := make([]byte, len(data), len(data)+2)
	copy(packet, data)
	packet = packet[:cap(packet)]
	defer r.changeState(SIDLE, STATE_IDLE)
	r.transmit(packet)
	if r.Error() == nil {
		r.stats.Packets.Sent++
		r.stats.Bytes.Sent += len(data)
	}
}

func (r *Radio) transmit(data []byte) {
	avail := fifoSize
	for r.Error() == nil {
		if avail > len(data) {
			avail = len(data)
		}
		r.WriteBurst(TXFIFO, data[:avail])
		r.changeState(STX, STATE_TX)
		data = data[avail:]
		if len(data) == 0 {
			break
		}
		// Transmitting a packet that is larger than the TXFIFO size.
		// See TI Design Note DN500 (swra109c).
		// Err on the short side here to avoid TXFIFO underflow.
		time.Sleep(fifoSize / 4 * byteDuration)
		for r.Error() == nil {
			n := r.ReadNumTxBytes()
			if n < fifoSize {
				avail = fifoSize - int(n)
				if avail > len(data) {
					avail = len(data)
				}
				break
			}
		}
	}
	r.finishTx(avail)
}

func (r *Radio) finishTx(numBytes int) {
	time.Sleep(time.Duration(numBytes) * byteDuration)
	for r.Error() == nil {
		n := r.ReadNumTxBytes()
		if n == 0 || r.Error() == TxFifoUnderflow {
			break
		}
		s := r.ReadState()
		if s != STATE_TX && s != STATE_TXFIFO_UNDERFLOW {
			log.Panicf("unexpected %s state during TXFIFO drain", StateName(s))
		}
		if verbose {
			log.Printf("waiting to transmit %d bytes in %s state", n, r.State())
		}
		time.Sleep(byteDuration)
	}
	if verbose {
		log.Printf("TX FIFO drained in %s state", r.State())
	}
}

func (r *Radio) Receive(timeout time.Duration) ([]byte, int) {
	if r.Error() != nil {
		return nil, 0
	}
	r.changeState(SRX, STATE_RX)
	defer r.changeState(SIDLE, STATE_IDLE)
	if verbose {
		log.Printf("waiting for interrupt in %s state", r.State())
	}
	r.err = r.interruptPin.Wait(timeout)
	startedWaiting := time.Time{}
	for r.Error() == nil {
		numBytes := r.ReadNumRxBytes()
		if r.Error() == RxFifoOverflow {
			// Flush RX FIFO and change back to RX.
			r.changeState(SRX, STATE_RX)
			continue
		}
		// Don't read last byte of FIFO if packet is still
		// being received. See Section 20 of data sheet.
		if numBytes < 2 {
			if startedWaiting.IsZero() {
				startedWaiting = time.Now()
			} else if time.Since(startedWaiting) >= timeout {
				break
			}
			time.Sleep(byteDuration)
			continue
		}
		if readFifoUsingBurst {
			data := r.ReadBurst(RXFIFO, int(numBytes))
			if r.Error() != nil {
				break
			}
			i := bytes.IndexByte(data, 0)
			if i == -1 {
				// No zero byte found; packet is still incoming.
				// Append all the data and continue to receive.
				_, r.err = r.receiveBuffer.Write(data)
				continue
			}
			// End of packet.
			_, r.err = r.receiveBuffer.Write(data[:i])
		} else {
			c := r.ReadRegister(RXFIFO)
			if r.Error() != nil {
				break
			}
			if c != 0 {
				r.err = r.receiveBuffer.WriteByte(c)
				continue
			}
		}
		// End of packet.
		rssi := r.ReadRSSI()
		size := r.receiveBuffer.Len()
		if size == 0 {
			if verbose {
				log.Printf("ignoring empty packet in %s state", r.State())
			}
			continue
		}
		r.stats.Packets.Received++
		r.stats.Bytes.Received += size
		p := make([]byte, size)
		_, r.err = r.receiveBuffer.Read(p)
		if r.Error() != nil {
			break
		}
		r.receiveBuffer.Reset()
		if verbose {
			log.Printf("received %d-byte packet in %s state; %d bytes remaining", size, r.State(), r.ReadNumRxBytes())
		}
		return p, rssi
	}
	return nil, 0
}

func (r *Radio) changeState(strobe byte, desired byte) {
	r.SetError(nil)
	s := r.ReadState()
	if s == desired {
		return
	}
	if verbose {
		log.Printf("change from %s to %s", StateName(s), StateName(desired))
	}
	for r.Error() == nil {
		switch s {
		case desired:
			return
		case STATE_RXFIFO_OVERFLOW:
			s = r.Strobe(SFRX)
		case STATE_TXFIFO_UNDERFLOW:
			s = r.Strobe(SFTX)
		default:
			s = r.Strobe(strobe)
		}
		s = (s >> STATE_SHIFT) & STATE_MASK
		if verbose {
			log.Printf("  %s", StateName(s))
		}
	}
}
