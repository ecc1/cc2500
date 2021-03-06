package cc2500

import (
	"errors"
	"unsafe"
)

var (
	// ErrRXFIFOOverflow indicates a RXFIFO overflow condition.
	ErrRXFIFOOverflow = errors.New("RXFIFO overflow")

	// ErrTXFIFOUnderflow indicates a TXFIFO underflow condition.
	ErrTXFIFOUnderflow = errors.New("TXFIFO underflow")
)

// Bytes returns the RFConfiguration as a byte slice.
func (config *RFConfiguration) Bytes() []byte {
	return (*[TEST0 - IOCFG2 + 1]byte)(unsafe.Pointer(config))[:]
}

// ReadConfiguration reads the current RFConfiguration from the radio.
func (r *Radio) ReadConfiguration() *RFConfiguration {
	if r.Error() != nil {
		return nil
	}
	regs := r.hw.ReadBurst(IOCFG2, TEST0-IOCFG2+1)
	return (*RFConfiguration)(unsafe.Pointer(&regs[0]))
}

// WriteConfiguration writes the given RFConfiguration to the radio.
func (r *Radio) WriteConfiguration(config *RFConfiguration) {
	r.hw.WriteBurst(IOCFG2, config.Bytes())
}

// InitRF initializes the radio to communicate with
// a Dexcom G4 continuous glucose monitor at the given frequency.
func (r *Radio) InitRF(frequency uint32) {
	rf := ResetRFConfiguration
	fb := frequencyToRegisters(frequency)

	// Asserts when sync word has been sent/received,
	// and de-asserts at the end of the packet.
	rf.IOCFG0 = 0x06

	rf.SYNC1 = 0xD3
	rf.SYNC0 = 0x91

	rf.PKTCTRL1 = PKTCTRL1_APPEND_STATUS | PKTCTRL1_ADR_CHK_NONE
	rf.PKTCTRL0 = PKTCTRL0_CRC_EN | PKTCTRL0_LENGTH_CONFIG_VARIABLE

	// Intermediate frequency
	// 0x09 * 26 MHz / 2^10 == 228515 Hz
	rf.FSCTRL1 = 0x09

	rf.FREQ2 = fb[0]
	rf.FREQ1 = fb[1]
	rf.FREQ0 = fb[2]

	// See table 20 in data sheet.
	// CHANBW_E = 1, CHANBW_M = 1, DRATE_E = 10
	// Channel BW = 26 MHz / (8 * (4 + CHANBW_M) * 2^CHANBW_E) == 325 kHz
	// See data sheet section 13 and Dexcom's FCC filing at
	// https://apps.fcc.gov/eas/GetApplicationAttachment.html?id=1373548
	rf.MDMCFG4 = 1<<MDMCFG4_CHANBW_E_SHIFT |
		1<<MDMCFG4_CHANBW_M_SHIFT |
		10<<MDMCFG4_DRATE_E_SHIFT

	// DRATE_M = 248 (0xF8)
	// Data rate = (256 + DRATE_M) * 2^DRATE_E * 26 MHz / 2^28 == 49987 Baud
	rf.MDMCFG3 = 0xF8

	rf.MDMCFG2 = MDMCFG2_DEM_DCFILT_ON |
		MDMCFG2_MOD_FORMAT_MSK |
		MDMCFG2_SYNC_MODE_30_32

	// CHANSPC_E = 3
	rf.MDMCFG1 = MDMCFG1_FEC_DIS |
		MDMCFG1_NUM_PREAMBLE_2 |
		3<<MDMCFG1_CHANSPC_E_SHIFT

	// CHANSPC_M = 59 (0x3B)
	// Channel spacing = (256 + CHANSPC_M) * 2^CHANSPC_E * 26 MHz / 2^18 == 249938 Hz
	rf.MDMCFG0 = 0x3B

	rf.DEVIATN = 0x40

	rf.MCSM2 = MCSM2_RX_TIME_END_OF_PACKET

	rf.MCSM1 = MCSM1_CCA_MODE_ALWAYS |
		MCSM1_RXOFF_MODE_IDLE |
		MCSM1_TXOFF_MODE_IDLE

	rf.MCSM0 = MCSM0_FS_AUTOCAL_FROM_IDLE

	rf.FOCCFG = FOCCFG_FOC_PRE_K_2K |
		FOCCFG_FOC_POST_K_PRE_K |
		FOCCFG_FOC_LIMIT_BW_OVER_4

	rf.BSCFG = BSCFG_BS_PRE_KI_2KI |
		BSCFG_BS_PRE_KP_3KP |
		BSCFG_BS_POST_KI_PRE_KI_OVER_2 |
		BSCFG_BS_POST_KP_PRE_KP |
		BSCFG_BS_LIMIT_0

	rf.AGCCTRL2 = AGCCTRL2_MAX_DVGA_GAIN_BUT_1 |
		AGCCTRL2_MAX_LNA_GAIN_0 |
		AGCCTRL2_MAGN_TARGET_36dB

	rf.AGCCTRL1 = AGCCTRL1_AGC_LNA_PRIORITY_0 |
		AGCCTRL1_CARRIER_SENSE_REL_THR_DISABLE |
		AGCCTRL1_CARRIER_SENSE_ABS_THR_0DB

	rf.AGCCTRL0 = AGCCTRL0_HYST_LEVEL_MEDIUM |
		AGCCTRL0_WAIT_TIME_32 |
		AGCCTRL0_AGC_FREEZE_NORMAL |
		AGCCTRL0_FILTER_LENGTH_32

	rf.FREND1 = 2<<FREND1_LNA_CURRENT_SHIFT |
		3<<FREND1_LNA2MIX_CURRENT_SHIFT |
		1<<FREND1_LODIV_BUF_CURRENT_RX_SHIFT |
		2<<FREND1_MIX_CURRENT_SHIFT

	// Use PA_TABLE 1 for transmitting '1' in ASK
	// (PA_TABLE 0 is always used for '0')
	rf.FREND0 = 1<<FREND0_LODIV_BUF_CURRENT_TX_SHIFT |
		0<<FREND0_PA_POWER_SHIFT

	rf.FSCAL3 = 2<<6 | 2<<4 | 0x09
	rf.FSCAL2 = 0x0A
	rf.FSCAL1 = 0x00
	rf.FSCAL0 = 0x20

	rf.TEST2 = TEST2_RX_LOW_DATA_RATE_MAGIC
	rf.TEST1 = TEST1_RX_LOW_DATA_RATE_MAGIC
	rf.TEST0 = 2<<2 | 1<<1 | 1<<0

	r.WriteConfiguration(&rf)

	// Power amplifier output settings (see section 24 of the data sheet)
	r.hw.WriteRegister(PATABLE, 0xBB)
}

// Frequency returns the radio's current frequency, in Hertz.
func (r *Radio) Frequency() uint32 {
	return registersToFrequency(r.hw.ReadBurst(FREQ2, 3))
}

func registersToFrequency(freq []byte) uint32 {
	f := uint32(freq[0])<<16 + uint32(freq[1])<<8 + uint32(freq[2])
	return uint32(uint64(f) * FXOSC >> 16)
}

// SetFrequency sets the radio to the given frequency, in Hertz.
func (r *Radio) SetFrequency(freq uint32) {
	r.hw.WriteBurst(FREQ2, frequencyToRegisters(freq))
}

func frequencyToRegisters(freq uint32) []byte {
	f := (uint64(freq)<<16 + FXOSC/2) / FXOSC
	return []byte{byte(f >> 16), byte(f >> 8), byte(f)}
}

func registerToFrequencyOffset(offset byte) int32 {
	return int32(int8(offset)) * FXOSC >> 14
}

func frequencyOffsetToRegister(offset int32) byte {
	return byte((int64(offset)<<14 + FXOSC/2) / FXOSC)
}

// ReadIF returns the radio's intermediate frequency, in Hertz.
func (r *Radio) ReadIF() uint32 {
	f := r.hw.ReadRegister(FSCTRL1)
	return uint32(uint64(f) * FXOSC >> 10)
}

// ReadChannelParams returns the radio's channel bandwidth and data rate.
func (r *Radio) ReadChannelParams() (uint32, uint32) {
	m4 := r.hw.ReadRegister(MDMCFG4)
	chanbwExp := (m4 >> MDMCFG4_CHANBW_E_SHIFT) & 0x3
	chanbwMant := (m4 >> MDMCFG4_CHANBW_M_SHIFT) & 0x3
	drateExp := (m4 >> MDMCFG4_DRATE_E_SHIFT) & 0xF
	drateMant := r.hw.ReadRegister(MDMCFG3)
	chanbw := uint32(FXOSC / ((4 + uint64(chanbwMant)) << (chanbwExp + 3)))
	drate := uint32(((256 + uint64(drateMant)) << drateExp * FXOSC) >> 28)
	return chanbw, drate
}

// ReadModemConfig returns the radio's modem configuration:
// whether FEC is enabled, the minimum preamble length, and the channel spacing.
func (r *Radio) ReadModemConfig() (bool, uint8, uint32) {
	m1 := r.hw.ReadRegister(MDMCFG1)
	fec := m1&MDMCFG1_FEC_EN != 0
	minPreamble := numPreamble[(m1&MDMCFG1_NUM_PREAMBLE_MASK)>>4]
	chanspcExp := m1 & MDMCFG1_CHANSPC_E_MASK
	chanspcMant := r.hw.ReadRegister(MDMCFG0)
	chanspc := uint32(((256 + uint64(chanspcMant)) << chanspcExp * FXOSC) >> 18)
	return fec, minPreamble, chanspc
}

func registerToRSSI(rssi byte) int {
	const rssiOffset = 72 // see data sheet section 17.3
	d := int(rssi)
	if d >= 128 {
		d -= 256
	}
	return d/2 - rssiOffset
}

// ReadRSSI returns the radio's RSSI, in dBm.
func (r *Radio) ReadRSSI() int {
	return registerToRSSI(r.hw.ReadRegister(RSSI))
}

// ReadPATable returns the contents of PATABLE.
func (r *Radio) ReadPATable() []byte {
	return r.hw.ReadBurst(PATABLE, 8)
}

// ReadNumRXBytes reads the RXBYTES register
// and detects RXFIFO overflow.
func (r *Radio) ReadNumRXBytes() byte {
	n := r.hw.ReadRegister(RXBYTES)
	if n&RXFIFO_OVERFLOW != 0 {
		r.Strobe(SFRX)
		r.err = ErrRXFIFOOverflow
	}
	return n & NUM_RXBYTES_MASK
}

// ReadNumTXBytes reads the TXBYTES register
// and detects TXFIFO underflow.
func (r *Radio) ReadNumTXBytes() byte {
	n := r.hw.ReadRegister(TXBYTES)
	if n&TXFIFO_UNDERFLOW != 0 {
		r.Strobe(SFTX)
		r.err = ErrTXFIFOUnderflow
	}
	return n & NUM_TXBYTES_MASK
}

// State returns the radio's current state as a string.
func (r *Radio) State() string {
	return StateName(r.ReadState())
}

// ReadState returns the radio's current state.
func (r *Radio) ReadState() byte {
	status := r.Strobe(SNOP)
	return (status >> STATE_SHIFT) & STATE_MASK
}

// StateName converts a state value to a string.
func StateName(state byte) string {
	return stateName[state]
}

// ReadMARCState returns the radio's MARC state.
func (r *Radio) ReadMARCState() byte {
	return r.hw.ReadRegister(MARCSTATE) & MARCSTATE_MASK
}

// MARCStateName converts a MARC state value to a string.
func MARCStateName(state byte) string {
	return marcState[state]
}

var (
	stateName = []string{
		"IDLE",
		"RX",
		"TX",
		"FSTXON",
		"CALIBRATE",
		"SETTLING",
		"RXFIFO_OVERFLOW",
		"TXFIFO_UNDERFLOW",
	}
	marcState = []string{
		"SLEEP",
		"IDLE",
		"XOFF",
		"VCOON_MC",
		"REGON_MC",
		"MANCAL",
		"VCOON",
		"REGON",
		"STARTCAL",
		"BWBOOST",
		"FS_LOCK",
		"IFADCON",
		"ENDCAL",
		"RX",
		"RX_END",
		"RX_RST",
		"TXRX_SWITCH",
		"RXFIFO_OVERFLOW",
		"FSTXON",
		"TX",
		"TX_END",
		"RXTX_SWITCH",
		"TXFIFO_UNDERFLOW",
	}
	strobeString = []string{
		"SRES",
		"SFSTXON",
		"SXOFF",
		"SCAL",
		"SRX",
		"STX",
		"SIDLE",
		"SAFC",
		"SWOR",
		"SPWD",
		"SFRX",
		"SFTX",
		"SWORRST",
		"SNOP",
	}
	numPreamble = []uint8{2, 3, 4, 6, 8, 12, 16, 24}
)
