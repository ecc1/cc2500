// +build customcs

package cc2500

// SPI configuration for Intel Edison with custom chip select.

const (
	spiDevice = "/dev/spidev5.1"
	customCS  = 44 // GPIO for custom chip select
)
