// +build !customcs

package cc2500

// SPI configuration for Intel Edison with kernel support for CS0.

const (
	spiDevice = "/dev/spidev5.0"
	customCS  = 0 // default chip select
)
