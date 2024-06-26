package net

import (
	"context"
	"fmt"
)

type Driver interface {
	Conn(context.Context, string, string) error
	Scan(context.Context) ([]network, error)
	Stat(context.Context) (*network, error)
}

// ---------------------------------------------

// AccessPoint ...
type AccessPoint interface {
	GetAccessType() string
	GetMaxBitrate() string
	GetMacAddr() string
	GetChannel() string
	GetQuality() string
	GetFreq() string
	GetSsid() string
}

// ---------------------------------------------

func Items() ([]AccessPoint, error) {
	n, err := NewDbusNm()
	if err != nil {
		return nil, err
	}
	return n.WirelessAccessPoints()

}

type network struct {
	bssid    string
	ssid     string
	mode     string
	channel  uint32
	rate     string
	signal   uint8
	bars     string
	security string
}

type Network struct{}

func NewNetwork() *Network {
	return &Network{}
}

func (n Network) Scan(ctx context.Context) ([]network, error) {
	return scan()
}

func (n Network) Stat(ctx context.Context) (*network, error) {
	return stat()
}

func (n Network) Conn(ctx context.Context, ssid, password string) error {
	return conn(ssid, password)
}

func (n network) GetBssid() string {
	return n.bssid
}

func (n network) GetSsid() string {
	return n.ssid
}

func (n network) GetMode() string {
	return n.mode
}

func (n network) GetChannel() string {
	return fmt.Sprintf("%d", n.channel)
}

func (n network) GetRate() string {
	return n.rate
}

func (n network) GetSignal() string {
	return fmt.Sprintf("%d", n.signal)
}

func (n network) GetBars() string {
	return fmt.Sprintf("%v", n.bars)
}

func (n network) GetSecurity() string {
	return n.security
}

// ---
