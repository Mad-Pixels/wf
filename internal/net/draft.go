//go:build linux

package net

// func stat() (*network, error) {

// 	// conn, err := dbus.SystemBus()
// 	// if err != nil {
// 	// 	return nil, fmt.Errorf("failed to connect to system bus: %v", err)
// 	// }

// 	// nm := conn.Object("org.freedesktop.NetworkManager", "/org/freedesktop/NetworkManager")
// 	// var activeConnections []dbus.ObjectPath
// 	// err = nm.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.NetworkManager", "ActiveConnections").Store(&activeConnections)
// 	// if err != nil {
// 	// 	return nil, fmt.Errorf("failed to get active connections: %v", err)
// 	// }

// 	// for _, connPath := range activeConnections {
// 	// 	connObj := conn.Object("org.freedesktop.NetworkManager", connPath)
// 	// 	var devices []dbus.ObjectPath
// 	// 	err = connObj.Call("org.freedesktop.NetworkManager.Connection.Active.GetDevices", 0).Store(&devices)
// 	// 	if err != nil {
// 	// 		continue
// 	// 	}

// 	// 	for _, devicePath := range devices {
// 	// 		deviceObj := conn.Object("org.freedesktop.NetworkManager", devicePath)
// 	// 		var deviceType uint32
// 	// 		err = deviceObj.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.NetworkManager.Device", "DeviceType").Store(&deviceType)
// 	// 		if err != nil {
// 	// 			continue
// 	// 		}

// 	// 		if deviceType == 2 { // 2 is the type for Wi-Fi devices
// 	// 			var apPath dbus.ObjectPath
// 	// 			err = deviceObj.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.NetworkManager.Device.Wireless", "ActiveAccessPoint").Store(&apPath)
// 	// 			if err != nil || apPath == "/" {
// 	// 				continue
// 	// 			}

// 	// 			apObj := conn.Object("org.freedesktop.NetworkManager", apPath)

// 	// 			var ssid []byte
// 	// 			err = apObj.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.NetworkManager.AccessPoint", "Ssid").Store(&ssid)
// 	// 			if err != nil {
// 	// 				continue
// 	// 			}

// 	// 			var freq uint32
// 	// 			err = apObj.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.NetworkManager.AccessPoint", "Frequency").Store(&freq)
// 	// 			if err != nil {
// 	// 				continue
// 	// 			}

// 	// 			var strength uint8
// 	// 			err = apObj.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.NetworkManager.AccessPoint", "Strength").Store(&strength)
// 	// 			if err != nil {
// 	// 				continue
// 	// 			}

// 	// 			return &network{
// 	// 				ssid:  string(ssid),
// 	// 				freq:  fmt.Sprintf("%d MHz", freq),
// 	// 				level: string(strength),
// 	// 			}, nil
// 	// 		}
// 	// 	}
// 	// }

// 	return nil, nil
// }

// func conn(ssid, password string) error {
// 	// conn, err := dbus.SystemBus()
// 	// if err != nil {
// 	// 	return fmt.Errorf("failed to connect to system bus: %v", err)
// 	// }

// 	// nm := conn.Object("org.freedesktop.NetworkManager", "/org/freedesktop/NetworkManager")
// 	// settings := conn.Object("org.freedesktop.NetworkManager", "/org/freedesktop/NetworkManager/Settings")

// 	// connection := map[string]map[string]dbus.Variant{
// 	// 	"802-11-wireless": {
// 	// 		"ssid": dbus.MakeVariant([]byte(ssid)),
// 	// 		"mode": dbus.MakeVariant("infrastructure"),
// 	// 	},
// 	// 	"802-11-wireless-security": {
// 	// 		"key-mgmt": dbus.MakeVariant("wpa-psk"),
// 	// 		"psk":      dbus.MakeVariant(password),
// 	// 	},
// 	// 	"connection": {
// 	// 		"type":        dbus.MakeVariant("802-11-wireless"),
// 	// 		"id":          dbus.MakeVariant(ssid),
// 	// 		"uuid":        dbus.MakeVariant(uuid.NewString()),
// 	// 		"autoconnect": dbus.MakeVariant(true),
// 	// 	},
// 	// 	"ipv4": {
// 	// 		"method": dbus.MakeVariant("auto"),
// 	// 	},
// 	// 	"ipv6": {
// 	// 		"method": dbus.MakeVariant("ignore"),
// 	// 	},
// 	// }

// 	// var path dbus.ObjectPath
// 	// err = settings.Call("org.freedesktop.NetworkManager.Settings.AddConnection", 0, connection).Store(&path)
// 	// if err != nil {
// 	// 	return fmt.Errorf("failed to add connection: %v", err)
// 	// }

// 	// devicePath, err := getWiFiDevicePath(conn)
// 	// if err != nil {
// 	// 	return fmt.Errorf("failed to get Wi-Fi device path: %v", err)
// 	// }

// 	// var activeConnPath dbus.ObjectPath
// 	// err = nm.Call("org.freedesktop.NetworkManager.ActivateConnection", 0, path, devicePath, "/").Store(&activeConnPath)
// 	// if err != nil {
// 	// 	return fmt.Errorf("failed to activate connection: %v", err)
// 	// }

// 	return nil
// }

// func getWiFiDevicePath(conn *dbus.Conn) (dbus.ObjectPath, error) {
// 	// nm := conn.Object("org.freedesktop.NetworkManager", "/org/freedesktop/NetworkManager")
// 	// var devices []dbus.ObjectPath
// 	// err := nm.Call("org.freedesktop.NetworkManager.GetDevices", 0).Store(&devices)
// 	// if err != nil {
// 	// 	return "", fmt.Errorf("failed to get devices: %v", err)
// 	// }

// 	// for _, devicePath := range devices {
// 	// 	device := conn.Object("org.freedesktop.NetworkManager", devicePath)
// 	// 	var deviceType uint32
// 	// 	err = device.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.NetworkManager.Device", "DeviceType").Store(&deviceType)
// 	// 	if err != nil {
// 	// 		continue
// 	// 	}

// 	// 	if deviceType == 2 { // 2 is the type for Wi-Fi devices
// 	// 		return devicePath, nil
// 	// 	}
// 	// }

// 	// return "", fmt.Errorf("no Wi-Fi device found")
// 	return "", nil
// }
