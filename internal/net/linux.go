//go:build linux

package net

import (
	"fmt"

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

				var frequency uint32
				err = ap.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.NetworkManager.AccessPoint", "Frequency").Store(&frequency)
				if err != nil {
					continue
				}

				var hwAddress string
				err = ap.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.NetworkManager.AccessPoint", "HwAddress").Store(&hwAddress)
				if err != nil {
					continue
				}

				var mode string
				err = ap.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.NetworkManager.AccessPoint", "Mode").Store(&mode)
				if err != nil {
					continue
				}

				var maxBitrate uint32
				err = ap.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.NetworkManager.AccessPoint", "MaxBitrate").Store(&maxBitrate)
				if err != nil {
					continue
				}

				var flags uint32
				err = ap.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.NetworkManager.AccessPoint", "Flags").Store(&flags)
				if err != nil {
					continue
				}

				var wpaFlags uint32
				err = ap.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.NetworkManager.AccessPoint", "WpaFlags").Store(&wpaFlags)
				if err != nil {
					continue
				}

				var rsnFlags uint32
				err = ap.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.NetworkManager.AccessPoint", "RsnFlags").Store(&rsnFlags)
				if err != nil {
					continue
				}

				security := getSecurity(flags, wpaFlags, rsnFlags)
				bars := getSignalBars(strength)

				networks = append(networks, network{
					bssid:    hwAddress,
					ssid:     string(ssid),
					mode:     mode,
					channel:  frequencyToChannel(frequency),
					rate:     fmt.Sprintf("%d Mbps", maxBitrate/1000),
					signal:   strength,
					bars:     bars,
					security: security,
				})
			}
		}
	}
	return networks, nil
}

// getSecurity determines the security type of the network.
func getSecurity(flags, wpaFlags, rsnFlags uint32) string {
	if rsnFlags != 0 {
		return "WPA2"
	}
	if wpaFlags != 0 {
		return "WPA"
	}
	if flags&0x1 != 0 {
		return "WEP"
	}
	return "Open"
}

// getSignalBars converts signal strength to bars.
func getSignalBars(signal uint8) string {
	switch {
	case signal > 75:
		return "▂▄▆█"
	case signal > 50:
		return "▂▄▆"
	case signal > 25:
		return "▂▄"
	default:
		return "▂"
	}
}

// frequencyToChannel converts frequency to channel.
func frequencyToChannel(freq uint32) uint32 {
	switch {
	case freq >= 2412 && freq <= 2484:
		return (freq - 2407) / 5
	case freq >= 5170 && freq <= 5825:
		return (freq - 5000) / 5
	default:
		return 0
	}
}

// func scan1() ([]network, error) {
// 	conn, err := dbus.SystemBus()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to connect to system bus: %v", err)
// 	}

// 	nm := conn.Object("org.freedesktop.NetworkManager", "/org/freedesktop/NetworkManager")
// 	var devices []dbus.ObjectPath
// 	err = nm.Call("org.freedesktop.NetworkManager.GetDevices", 0).Store(&devices)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get devices: %v", err)
// 	}

// 	var networks []network
// 	for _, devicePath := range devices {
// 		device := conn.Object("org.freedesktop.NetworkManager", devicePath)
// 		var deviceType uint32
// 		err = device.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.NetworkManager.Device", "DeviceType").Store(&deviceType)
// 		if err != nil {
// 			continue
// 		}

// 		// 2 is the type for Wi-Fi devices
// 		if deviceType == 2 {
// 			var accessPoints []dbus.ObjectPath
// 			err = device.Call("org.freedesktop.NetworkManager.Device.Wireless.GetAllAccessPoints", 0).Store(&accessPoints)
// 			if err != nil {
// 				continue
// 			}

// 			for _, apPath := range accessPoints {
// 				ap := conn.Object("org.freedesktop.NetworkManager", apPath)

// 				var ssid []byte
// 				err = ap.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.NetworkManager.AccessPoint", "Ssid").Store(&ssid)
// 				if err != nil {
// 					continue
// 				}

// 				var strength uint8
// 				err = ap.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.NetworkManager.AccessPoint", "Strength").Store(&strength)
// 				if err != nil {
// 					continue
// 				}

// 				networks = append(networks, network{
// 					ssid:  string(ssid),
// 					level: string(strength),
// 				})
// 			}
// 		}
// 	}
// 	return networks, nil
// }

func stat() (*network, error) {

	// conn, err := dbus.SystemBus()
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to connect to system bus: %v", err)
	// }

	// nm := conn.Object("org.freedesktop.NetworkManager", "/org/freedesktop/NetworkManager")
	// var activeConnections []dbus.ObjectPath
	// err = nm.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.NetworkManager", "ActiveConnections").Store(&activeConnections)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get active connections: %v", err)
	// }

	// for _, connPath := range activeConnections {
	// 	connObj := conn.Object("org.freedesktop.NetworkManager", connPath)
	// 	var devices []dbus.ObjectPath
	// 	err = connObj.Call("org.freedesktop.NetworkManager.Connection.Active.GetDevices", 0).Store(&devices)
	// 	if err != nil {
	// 		continue
	// 	}

	// 	for _, devicePath := range devices {
	// 		deviceObj := conn.Object("org.freedesktop.NetworkManager", devicePath)
	// 		var deviceType uint32
	// 		err = deviceObj.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.NetworkManager.Device", "DeviceType").Store(&deviceType)
	// 		if err != nil {
	// 			continue
	// 		}

	// 		if deviceType == 2 { // 2 is the type for Wi-Fi devices
	// 			var apPath dbus.ObjectPath
	// 			err = deviceObj.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.NetworkManager.Device.Wireless", "ActiveAccessPoint").Store(&apPath)
	// 			if err != nil || apPath == "/" {
	// 				continue
	// 			}

	// 			apObj := conn.Object("org.freedesktop.NetworkManager", apPath)

	// 			var ssid []byte
	// 			err = apObj.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.NetworkManager.AccessPoint", "Ssid").Store(&ssid)
	// 			if err != nil {
	// 				continue
	// 			}

	// 			var freq uint32
	// 			err = apObj.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.NetworkManager.AccessPoint", "Frequency").Store(&freq)
	// 			if err != nil {
	// 				continue
	// 			}

	// 			var strength uint8
	// 			err = apObj.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.NetworkManager.AccessPoint", "Strength").Store(&strength)
	// 			if err != nil {
	// 				continue
	// 			}

	// 			return &network{
	// 				ssid:  string(ssid),
	// 				freq:  fmt.Sprintf("%d MHz", freq),
	// 				level: string(strength),
	// 			}, nil
	// 		}
	// 	}
	// }

	return nil, nil
}

func conn(ssid, password string) error {
	// conn, err := dbus.SystemBus()
	// if err != nil {
	// 	return fmt.Errorf("failed to connect to system bus: %v", err)
	// }

	// nm := conn.Object("org.freedesktop.NetworkManager", "/org/freedesktop/NetworkManager")
	// settings := conn.Object("org.freedesktop.NetworkManager", "/org/freedesktop/NetworkManager/Settings")

	// connection := map[string]map[string]dbus.Variant{
	// 	"802-11-wireless": {
	// 		"ssid": dbus.MakeVariant([]byte(ssid)),
	// 		"mode": dbus.MakeVariant("infrastructure"),
	// 	},
	// 	"802-11-wireless-security": {
	// 		"key-mgmt": dbus.MakeVariant("wpa-psk"),
	// 		"psk":      dbus.MakeVariant(password),
	// 	},
	// 	"connection": {
	// 		"type":        dbus.MakeVariant("802-11-wireless"),
	// 		"id":          dbus.MakeVariant(ssid),
	// 		"uuid":        dbus.MakeVariant(uuid.NewString()),
	// 		"autoconnect": dbus.MakeVariant(true),
	// 	},
	// 	"ipv4": {
	// 		"method": dbus.MakeVariant("auto"),
	// 	},
	// 	"ipv6": {
	// 		"method": dbus.MakeVariant("ignore"),
	// 	},
	// }

	// var path dbus.ObjectPath
	// err = settings.Call("org.freedesktop.NetworkManager.Settings.AddConnection", 0, connection).Store(&path)
	// if err != nil {
	// 	return fmt.Errorf("failed to add connection: %v", err)
	// }

	// devicePath, err := getWiFiDevicePath(conn)
	// if err != nil {
	// 	return fmt.Errorf("failed to get Wi-Fi device path: %v", err)
	// }

	// var activeConnPath dbus.ObjectPath
	// err = nm.Call("org.freedesktop.NetworkManager.ActivateConnection", 0, path, devicePath, "/").Store(&activeConnPath)
	// if err != nil {
	// 	return fmt.Errorf("failed to activate connection: %v", err)
	// }

	return nil
}

func getWiFiDevicePath(conn *dbus.Conn) (dbus.ObjectPath, error) {
	// nm := conn.Object("org.freedesktop.NetworkManager", "/org/freedesktop/NetworkManager")
	// var devices []dbus.ObjectPath
	// err := nm.Call("org.freedesktop.NetworkManager.GetDevices", 0).Store(&devices)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to get devices: %v", err)
	// }

	// for _, devicePath := range devices {
	// 	device := conn.Object("org.freedesktop.NetworkManager", devicePath)
	// 	var deviceType uint32
	// 	err = device.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.NetworkManager.Device", "DeviceType").Store(&deviceType)
	// 	if err != nil {
	// 		continue
	// 	}

	// 	if deviceType == 2 { // 2 is the type for Wi-Fi devices
	// 		return devicePath, nil
	// 	}
	// }

	// return "", fmt.Errorf("no Wi-Fi device found")
	return "", nil
}
