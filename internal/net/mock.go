//go:build darwin
// +build darwin

package net

import (
	"math/rand"
	"strconv"
)

// access point data.
type accessPoint struct {
	strength   uint8
	frequency  uint32
	maxBitrate uint32
	flags      uint32
	wpaFlags   uint32
	rsnFlags   uint32
	mode       uint32
	hwAddress  string
	ssid       string
}

// GetSsid return access point name.
func (ap accessPoint) GetSsid() string {
	return ap.ssid
}

// GetQuality return access point quality level.
func (ap accessPoint) GetQuality() string {
	return strconv.Itoa(int(ap.strength))
}

// GetChannel return access point channel.
func (ap accessPoint) GetChannel() string {
	return strconv.Itoa(int(ap.frequency))
}

// GetFreq return access point frequency in MHz.
func (ap accessPoint) GetFreq() string {
	return strconv.Itoa(int(ap.frequency)) + " MHz"
}

// GetMaxBitrate return access point bitrate in Mbps.
func (ap accessPoint) GetMaxBitrate() string {
	return strconv.Itoa(int(ap.maxBitrate)) + " Mbps"
}

// GetMacAddr return access point mac address.
func (ap accessPoint) GetMacAddr() string {
	return ap.hwAddress
}

// GetMode return access point mode type.
func (ap accessPoint) GetMode() string {
	switch ap.mode {
	case 1:
		return "Ad-hoc"
	case 2:
		return "Infra"
	case 3:
		return "AP"
	default:
		return "NaN"
	}
}

// GetAccessType return access point acces type.
func (ap accessPoint) GetAccessType() string {
	if ap.rsnFlags != 0 {
		return "WPA2"
	}
	if ap.wpaFlags != 0 {
		return "WPA"
	}
	if ap.flags&0x1 != 0 {
		return "WEP"
	}
	return "OPEN"
}

// mock object
type mockDriver struct {
	currentAp *accessPoint
}

func NewDriver() (driverInterface, error) {
	return &mockDriver{}, nil
}

// driverInterface method which send request for Re-Scan wireless networks.
func (md *mockDriver) wirelessScan() error {
	return nil
}

// driverInterface method which return current wireless access point.
func (md *mockDriver) currentConnetcion() (AccessPoint, error) {
	return md.currentAp, nil
}

// direverInterface method which init connection to wireless network.
func (md *mockDriver) wirelessConnect(ssid, password string) error {
	md.currentAp = randomNetwork().(*accessPoint)
	md.currentAp.ssid = ssid
	return nil
}

// driverInterface method which return available access points data as []AccessPoint object.
func (md *mockDriver) wirelessAccessPoints() ([]AccessPoint, error) {
	numNetworks := rand.Intn(30) + 1

	networks := make([]AccessPoint, numNetworks)
	for i := 0; i < numNetworks; i++ {
		networks[i] = randomNetwork()
	}
	return networks, nil
}

func randomNetwork() AccessPoint {
	frequencies := []uint32{
		2412,
		2417,
		2422,
		2427,
		2432,
		2437,
		5500,
		5510,
		5520,
		5940,
		5950,
		5960,
		5970,
		5980,
		5990,
	}
	return &accessPoint{
		strength:   uint8(rand.Intn(100)),
		frequency:  frequencies[rand.Intn(len(frequencies))],
		maxBitrate: uint32(rand.Intn(1000) + 100),
		flags:      uint32(rand.Intn(4)),
		wpaFlags:   uint32(rand.Intn(2)),
		rsnFlags:   uint32(rand.Intn(2)),
		mode:       uint32(rand.Intn(4)),
		hwAddress:  randomMACAddress(),
		ssid:       randomString(8),
	}
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func randomMACAddress() string {
	const hexDigits = "0123456789ABCDEF"
	mac := make([]byte, 17)
	for i := 0; i < 17; i++ {
		switch i {
		case 2, 5, 8, 11, 14:
			mac[i] = ':'
		default:
			mac[i] = hexDigits[rand.Intn(16)]
		}
	}
	return string(mac)
}
