package net

import (
	"context"
	"math/rand"
	"time"
)

type NetworkManager struct{}

func NewNetworkManager() Driver {
	return &NetworkManager{}
}

func (n NetworkManager) Scan(ctx context.Context) ([]network, error) {
	rand.Seed(time.Now().UnixNano())
	numNetworks := rand.Intn(10) + 1

	networks := make([]network, numNetworks)
	for i := 0; i < numNetworks; i++ {
		networks[i] = randomNetwork()
	}
	return networks, nil

}

func (n NetworkManager) Stat(ctx context.Context) (*network, error) {
	nn := randomNetwork()
	return &nn, nil
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func randomNetwork() network {
	return network{
		ssid:    randomString(5),
		freq:    randomString(4),
		level:   randomString(6),
		quality: randomString(8),
	}
}
