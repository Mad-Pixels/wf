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
	nmDestWirelessAccessPoints = nmDest + ".Device.Wireless.GetAllAccessPoints"
)

type dbusNm struct {
	conn *dbus.Conn

	devicesWireless []dbus.ObjectPath
}

// NewDbusNm return dbusNm object.
func NewDbusNm() (*dbusNm, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}

	var (
		devicesAll      = []dbus.ObjectPath{}
		devicesWireless = []dbus.ObjectPath{}
	)
	if err := conn.
		Object(nmDest, "/org/freedesktop/NetworkManager").
		Call(nmDestDevices, 0).
		Store(&devicesAll); err != nil {
		return nil, err
	}
	for _, device := range devicesAll {
		var deviceType uint32
		if err := conn.
			Object(nmDest, device).
			Call(nmPropGet, 0, nmDestDevice, "DeviceType").
			Store(&deviceType); err != nil {
			continue
		}
		switch deviceType {
		case 2:
			devicesWireless = append(devicesWireless, device)
		default:
			continue
		}
	}

	return &dbusNm{
		conn: conn,

		devicesWireless: devicesWireless,
	}, nil
}
