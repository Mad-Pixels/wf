package net

import "log"

// Main Network Driver.
var Driver = &driver{}

// Initialize or error main Driver object.
func init() {
	driver, err := NewDriver()
	if err != nil {
		log.Fatal(err)
	}

	Driver.driverInterface = driver
}

// AccessPoint represent network access point.
type AccessPoint interface {
	GetAccessType() string
	GetMaxBitrate() string
	GetMacAddr() string
	GetChannel() string
	GetQuality() string
	GetFreq() string
	GetSsid() string
}

// represent Driver object.
type driverInterface interface {
	wirelessAccessPoints() ([]AccessPoint, error)
	currentConnetcion() (AccessPoint, error)
	wirelessScan() error
}

type driver struct {
	driverInterface
}

// WirelessAccessPoints return list of active access points.
func (d driver) WirelessAccessPoints() ([]AccessPoint, error) {
	return d.wirelessAccessPoints()
}

// WirelessScan manual trigger for re-scan active access points.
func (d driver) WirelessScan() error {
	return d.wirelessScan()
}

// CurrentConnection return current access point.
func (d driver) CurrentConnetcion() (AccessPoint, error) {
	return d.currentConnetcion()
}
