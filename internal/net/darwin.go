//go:build darwin

package net

import (
	"math/rand"
	"time"
)

func scan() ([]network, error) {
	rand.Seed(time.Now().UnixNano())
	numNetworks := rand.Intn(100) + 1

	networks := make([]network, numNetworks)
	for i := 0; i < numNetworks; i++ {
		networks[i] = randomNetwork()
	}
	return networks, nil
}

func stat() (*network, error) {
	nn := randomNetwork()
	return &nn, nil
}

func randomNetwork() network {
	return network{
		ssid:    randomString(5),
		freq:    randomString(4),
		level:   randomString(6),
		quality: randomString(8),
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

func conn(ssid, password string) error {
	return nil
}
