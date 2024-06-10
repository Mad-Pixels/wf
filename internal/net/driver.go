package net

import "context"

type Driver interface {
	Scan(context.Context) ([]network, error)
}

type network struct {
	ssid    string
	freq    string
	level   string
	quality string
}

func (n network) GetSsid() string {
	return n.ssid
}

func (n network) GetFreq() string {
	return n.freq
}

func (n network) GetLevel() string {
	return n.level
}

func (n network) GetQuality() string {
	return n.quality
}
