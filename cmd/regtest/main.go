package main

import (
	"fmt"
	"log"

	"github.com/ecc1/cc2500"
)

func main() {
	r := cc2500.Open()
	if r.Error() != nil {
		log.Fatal(r.Error())
	}
	r.Reset()

	dumpRegs(r)

	fmt.Printf("\nTesting individual writes\n")
	hw := r.Hardware()
	hw.WriteRegister(cc2500.SYNC1, 0x44)
	hw.WriteRegister(cc2500.SYNC0, 0x55)
	readRegs(r)

	r.Reset()
	fmt.Printf("\nTesting burst writes\n")
	hw.WriteBurst(cc2500.SYNC1, []byte{0x66, 0x77})
	readRegs(r)
}

func dumpRegs(r *cc2500.Radio) {
	fmt.Printf("\nConfiguration registers:\n")
	regs := r.ReadConfiguration().Bytes()
	resetValue := cc2500.ResetRFConfiguration.Bytes()
	if r.Error() != nil {
		log.Fatal(r.Error())
	}
	for i, v := range regs {
		fmt.Printf("%02X  %02X  %08b", cc2500.IOCFG2+i, v, v)
		r := resetValue[i]
		if v == r {
			fmt.Printf("\n")
		} else {
			fmt.Printf("  **** SHOULD BE %02X  %08b\n", r, r)
		}
	}
}

func readRegs(r *cc2500.Radio) {
	hw := r.Hardware()
	x := hw.ReadRegister(cc2500.SYNC1)
	y := hw.ReadRegister(cc2500.SYNC0)
	if r.Error() != nil {
		log.Fatal(r.Error())
	}
	fmt.Printf("individual: %X %X\n", x, y)
	v := hw.ReadBurst(cc2500.SYNC1, 2)
	if r.Error() != nil {
		log.Fatal(r.Error())
	}
	fmt.Printf("  burst:    %X %X\n", v[0], v[1])
}
