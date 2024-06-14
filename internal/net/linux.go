//go:build linux

package net

import (
	"fmt"
	"math/rand"

	"github.com/godbus/dbus/v5"
)

func scan() ([]network, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to system bus: %v", err)
	}

	nm := conn.Object("org.freedesktop.NetworkManager", "/org/freedesktop/NetworkManager")
	var devices []dbus.ObjectPath
	err = nm.Call("org.freedesktop.NetworkManager.GetDevices", 0).Store(&devices)
	if err != nil {
		return nil, fmt.Errorf("failed to get devices: %v", err)
	}

	var networks []network
	for _, devicePath := range devices {
		device := conn.Object("org.freedesktop.NetworkManager", devicePath)
		var deviceType uint32
		err = device.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.NetworkManager.Device", "DeviceType").Store(&deviceType)
		if err != nil {
			continue
		}

		// 2 is the type for Wi-Fi devices
		if deviceType == 2 {
			var accessPoints []dbus.ObjectPath
			err = device.Call("org.freedesktop.NetworkManager.Device.Wireless.GetAllAccessPoints", 0).Store(&accessPoints)
			if err != nil {
				continue
			}

			for _, apPath := range accessPoints {
				ap := conn.Object("org.freedesktop.NetworkManager", apPath)

				var ssid []byte
				err = ap.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.NetworkManager.AccessPoint", "Ssid").Store(&ssid)
				if err != nil {
					continue
				}

				var strength uint8
				err = ap.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.NetworkManager.AccessPoint", "Strength").Store(&strength)
				if err != nil {
					continue
				}

				networks = append(networks, network{
					ssid:  string(ssid),
					level: string(strength),
				})
			}
		}
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
