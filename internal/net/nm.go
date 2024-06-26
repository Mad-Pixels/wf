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

func (ap accessPoint) GetSsid() string {
	return ap.ssid
}

func (ap accessPoint) GetQuality() uint8 {
	return ap.strength
}

func (ap accessPoint) GetFreq() uint32 {
	switch {
	case ap.frequency >= 2412 && ap.frequency <= 2484:
		return (ap.frequency - 2407) / 5
	case ap.frequency >= 5170 && ap.frequency <= 5825:
		return (ap.frequency - 5000) / 5
	default:
		return 0
	}
}

func (ap accessPoint) GetMaxBitrate() uint32 {
	return ap.maxBitrate
}

func (ap accessPoint) GetMacAddr() string {
	return ap.hwAddress
}

func (ap accessPoint) GetSecType() string {
	if ap.rsnFlags != 0 {
		return "WPA2"
	}
	if ap.wpaFlags != 0 {
		return "WPA"
	}
	if ap.flags&0x1 != 0 {
		return "WEP"
	}
	return "OPEN"
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

func (nm dbusNm) WirelessAccessPoints() ([]AccessPoint, error) {
	apList := []AccessPoint{}

	for _, device := range nm.devicesWireless {
		var apPaths []dbus.ObjectPath

		if err := device.
			Call(nmDestWirelessAccessPoints, 0).
			Store(&apPaths); err != nil {
			return nil, err
		}

		for _, apPath := range apPaths {
			apObject, err := dbusObjAp{
				object: nm.conn.Object(nmDest, apPath),
			}.AccessPointAsync()
			if err != nil {
				return nil, err
			}
			apList = append(apList, apObject)
		}
	}

	return apList, nil
}

type dbusObjAp struct {
	object dbus.BusObject
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

	err := ap.object.Call(nmPropGet, 0, nmDestWirelessAccessPoint, "Ssid").Store(&val)
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func (ap dbusObjAp) strength() (uint8, error) {
	var val uint8

	err := ap.object.Call(nmPropGet, 0, nmDestWirelessAccessPoint, "Strength").Store(&val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (ap dbusObjAp) frequency() (uint32, error) {
	var val uint32

	err := ap.object.Call(nmPropGet, 0, nmDestWirelessAccessPoint, "Frequency").Store(&val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (ap dbusObjAp) hwAddress() (string, error) {
	var val string

	err := ap.object.Call(nmPropGet, 0, nmDestWirelessAccessPoint, "HwAddress").Store(&val)
	if err != nil {
		return "", err
	}
	return val, nil
}

func (ap dbusObjAp) mode() (string, error) {
	var val string

	err := ap.object.Call(nmPropGet, 0, nmDestWirelessAccessPoint, "Mode").Store(&val)
	if err != nil {
		return "", err
	}
	return val, nil
}

func (ap dbusObjAp) maxBitrate() (uint32, error) {
	var val uint32

	err := ap.object.Call(nmPropGet, 0, nmDestWirelessAccessPoint, "MaxBitrate").Store(&val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (ap dbusObjAp) flags() (uint32, error) {
	var val uint32

	err := ap.object.Call(nmPropGet, 0, nmDestWirelessAccessPoint, "Flags").Store(&val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (ap dbusObjAp) wpaFlags() (uint32, error) {
	var val uint32

	err := ap.object.Call(nmPropGet, 0, nmDestWirelessAccessPoint, "WpaFlags").Store(&val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (ap dbusObjAp) rsnFlags() (uint32, error) {
	var val uint32

	err := ap.object.Call(nmPropGet, 0, nmDestWirelessAccessPoint, "RsnFlags").Store(&val)
	if err != nil {
		return 0, err
	}
	return val, nil
}
