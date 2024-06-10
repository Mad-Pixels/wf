package net

import "context"

type NetworkManager struct{}

func NewNetworkManager() Driver {
	return &NetworkManager{}
}

func (n NetworkManager) Scan(ctx context.Context) ([]network, error) {
	f := []network{
		{
			ssid:    "ssid1",
			freq:    "freq1",
			level:   "level1",
			quality: "quality1",
		},
		{
			ssid:    "ssid2",
			freq:    "freq2",
			level:   "level2",
			quality: "quality2",
		},
	}
	return f, nil
}
