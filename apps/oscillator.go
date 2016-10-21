package main

// Route the CC2500 clock oscillator (divided by 16) to the GDO0 pin,
// so that its frequency can be measured easily with an oscilloscope
// or frequency counter.

// IMPORTANT: disconnect the Edison GPIO used for interrupts before
// running this program.  Otherwise the Edison will be interrupted
// at 1.625 MHz and become non-responsive.

import (
	"log"

	"github.com/ecc1/cc2500"
)

func main() {
	r := cc2500.Open()
	r.Reset()
	// Route CLK_XOSC/24 to GDO0 pin.
	// See data sheet, Table 33.
	r.Hardware().WriteRegister(cc2500.IOCFG0, 0x38)
	if r.Error() != nil {
		log.Fatal(r.Error())
	}
}
