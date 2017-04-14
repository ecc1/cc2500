package cc2500

// Configuration for Intel Edison.

const (
	spiDevice    = "/dev/spidev5.1"
	spiSpeed     = 6000000 // Hz
	interruptPin = 48      // GPIO for receive interrupts
)
