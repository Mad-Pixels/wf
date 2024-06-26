package net

import "github.com/godbus/dbus/v5"

const (
	nmDest = "org.freedesktop.NetworkManager"
	nmProp = "org.freedesktop.DBus.Properties"
)

var (
	nmPropGet = nmProp + ".Get"

	nmDestDevice               = nmDest + ".Device"
	nmDestDevices              = nmDest + ".GetDevices"
	nmDestWirelessAccessPoint  = nmDest + ".AccessPoint"
	nmDestWirelessAccessPoints = nmDest + ".Device.Wireless.GetAllAccessPoints"
)

type dbusNm struct {
	conn *dbus.Conn

	devicesWireless []dbus.BusObject
}

// NewDbusNm return dbusNm object.
func NewDbusNm() (*dbusNm, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}

	var (
		devicesPaths    = []dbus.ObjectPath{}
		devicesWireless = []dbus.BusObject{}
	)
	if err := conn.
		Object(nmDest, "/org/freedesktop/NetworkManager").
		Call(nmDestDevices, 0).
		Store(&devicesPaths); err != nil {
		return nil, err
	}

	for _, devicePath := range devicesPaths {
		var (
			deviceObj  = conn.Object(nmDest, devicePath)
			deviceType uint32
		)
		if err := deviceObj.
			Call(nmPropGet, 0, nmDestDevice, "DeviceType").
			Store(&deviceType); err != nil {
			continue
		}
		switch deviceType {
		case 2:
			devicesWireless = append(devicesWireless, deviceObj)
		default:
			continue
		}
	}

	return &dbusNm{
		conn: conn,

		devicesWireless: devicesWireless,
	}, nil
}

func (nm dbusNm) WirelessAccessPoints() (interface{}, error) {
	var apPaths []dbus.ObjectPath

	for _, device := range nm.devicesWireless {
		if err := device.
			Call(nmDestWirelessAccessPoints, 0).
			Store(&apPaths); err != nil {
			return nil, err
		}

		for _, apPath := range apPaths {
			apObject := dbusObjAp{object: nm.conn.Object(nmDest, apPath)} //newDbusObjAp(nm.conn, apPath)

		}

	}
	return nil, nil
}

type dbusObjAp struct {
	object dbus.BusObject
}

func (ap dbusObjAp) ssid() (string, error) {
	var val []byte

	if err := ap.store("Ssid", val); err != nil {
		return "", err
	}
	return string(val), nil
}

func (ap dbusObjAp) strength() (uint8, error) {
	var val uint8

	if err := ap.store("Strength", val); err != nil {
		return 0, err
	}
	return val, nil
}

func (ap dbusObjAp) frequency() (uint32, error) {
	var val uint32

	if err := ap.store("Frequency", val); err != nil {
		return 0, err
	}
	return val, nil
}

func (ap dbusObjAp) hwAddress() (string, error) {
	var val string

	if err := ap.store("HwAddress", val); err != nil {
		return "", err
	}
	return val, nil
}

func (ap dbusObjAp) mode() (string, error) {
	var val string

	if err := ap.store("Mode", val); err != nil {
		return "", err
	}
	return val, nil
}

func (ap dbusObjAp) maxBitrate() (uint32, error) {
	var val uint32

	if err := ap.store("MaxBitrate", val); err != nil {
		return 0, err
	}
	return val, nil
}

func (ap dbusObjAp) flags() (uint32, error) {
	var val uint32

	if err := ap.store("Flags", val); err != nil {
		return 0, err
	}
	return val, nil
}

func (ap dbusObjAp) wpaFlags() (uint32, error) {
	var val uint32

	if err := ap.store("WpaFlags", val); err != nil {
		return 0, err
	}
	return val, nil
}

func (ap dbusObjAp) rsnFlags() (uint32, error) {
	var val uint32

	if err := ap.store("RsnFlags", val); err != nil {
		return 0, err
	}
	return val, nil
}

func (ap dbusObjAp) store(field string, value any) error {
	return ap.object.
		Call(nmPropGet, 0, nmDestWirelessAccessPoint, field).
		Store(value)
}
