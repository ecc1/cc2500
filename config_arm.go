package cc2500

// Configuration for Raspberry Pi Zero W.

const (
	spiDevice    = "/dev/spidev1.0"
	spiSpeed     = 6000000 // Hz
	customCS     = 0       // default chip select
	interruptPin = 22      // GPIO for receive interrupts
)
