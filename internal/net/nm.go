package net

import (
	"github.com/godbus/dbus/v5"
)

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

type accessPoint struct {
	strength   uint8
	frequency  uint32
	maxBitrate uint32
	flags      uint32
	wpaFlags   uint32
	rsnFlags   uint32
	hwAddress  string
	ssid       string
	mode       string
}

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

func (ap dbusObjAp) AccessPoint() (AccessPoint, error) {
	ssid, err := ap.ssid()
	if err != nil {
		return nil, err
	}
	strength, err := ap.strength()
	if err != nil {
		return nil, err
	}
	frequency, err := ap.frequency()
	if err != nil {
		return nil, err
	}
	hwAddress, err := ap.hwAddress()
	if err != nil {
		return nil, err
	}
	mode, err := ap.mode()
	if err != nil {
		return nil, err
	}
	maxBitrate, err := ap.maxBitrate()
	if err != nil {
		return nil, err
	}
	flags, err := ap.flags()
	if err != nil {
		return nil, err
	}
	wpaFlags, err := ap.wpaFlags()
	if err != nil {
		return nil, err
	}
	rsnFlags, err := ap.rsnFlags()
	if err != nil {
		return nil, err
	}

	return &accessPoint{
		ssid:       ssid,
		strength:   strength,
		frequency:  frequency,
		hwAddress:  hwAddress,
		mode:       mode,
		maxBitrate: maxBitrate,
		flags:      flags,
		wpaFlags:   wpaFlags,
		rsnFlags:   rsnFlags,
	}, nil
}

func (ap dbusObjAp) AccessPointAsync() (AccessPoint, error) {
	type res struct {
		field string
		err   error
		value any
	}
	var (
		fieldCount = 9

		ch = make(chan res, fieldCount)
		ex = func(field string, f func() (any, error)) {
			value, err := f()
			ch <- res{
				field: field,
				value: value,
				err:   err,
			}

		}
	)
	defer close(ch)

	go ex("ssid", func() (any, error) { return ap.ssid() })
	go ex("strength", func() (any, error) { return ap.strength() })
	go ex("frequency", func() (any, error) { return ap.frequency() })
	go ex("hwAddress", func() (any, error) { return ap.hwAddress() })
	go ex("mode", func() (any, error) { return ap.mode() })
	go ex("maxBitrate", func() (any, error) { return ap.maxBitrate() })
	go ex("flags", func() (any, error) { return ap.flags() })
	go ex("wpaFlags", func() (any, error) { return ap.wpaFlags() })
	go ex("rsnFlags", func() (any, error) { return ap.rsnFlags() })

	accPoint := &accessPoint{}
	for i := 0; i < fieldCount; i++ {
		res := <-ch
		if res.err != nil {
			return nil, res.err
		}

		switch res.field {
		case "ssid":
			accPoint.ssid = res.value.(string)
		case "strength":
			accPoint.strength = res.value.(uint8)
		case "frequency":
			accPoint.frequency = res.value.(uint32)
		case "hwAddress":
			accPoint.hwAddress = res.value.(string)
		case "mode":
			accPoint.mode = res.value.(string)
		case "maxBitrate":
			accPoint.maxBitrate = res.value.(uint32)
		case "flags":
			accPoint.flags = res.value.(uint32)
		case "wpaFlags":
			accPoint.wpaFlags = res.value.(uint32)
		case "rsnFlags":
			accPoint.rsnFlags = res.value.(uint32)
		}
	}

	return accPoint, nil
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
