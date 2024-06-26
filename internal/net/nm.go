package net

import (
	"fmt"
	"strconv"

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

// access point data.
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

// GetSsid return access point name.
func (ap accessPoint) GetSsid() string {
	return ap.ssid
}

// GetQuality return access point quality level.
func (ap accessPoint) GetQuality() string {
	return strconv.Itoa(int(ap.strength))
}

// GetChannel return access point channel.
func (ap accessPoint) GetChannel() string {
	switch {
	case ap.frequency >= 2412 && ap.frequency <= 2484:
		return strconv.Itoa(int((ap.frequency - 2407) / 5))
	case ap.frequency >= 5170 && ap.frequency <= 5825:
		return strconv.Itoa(int((ap.frequency - 5000) / 5))
	default:
		return "NaN"
	}
}

// GetFreq return access point frequency in MHz.
func (ap accessPoint) GetFreq() string {
	return strconv.Itoa(int(ap.frequency)) + " MHz"
}

// GetMaxBitrate return access point bitrate in Mbps.
func (ap accessPoint) GetMaxBitrate() string {
	return fmt.Sprintf("%.2f Mbps", float64(ap.maxBitrate)/1000)
}

// GetMacAddr return access point mac address.
func (ap accessPoint) GetMacAddr() string {
	return ap.hwAddress
}

// GetAccessType return access point acces type.
func (ap accessPoint) GetAccessType() string {
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

// dbus NetworkManager object.
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

// WirelessAccessPoints return available access points data as []AccessPoint object.
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
			}.AccessPoint()
			if err != nil {
				return nil, err
			}
			apList = append(apList, apObject)
		}
	}

	return apList, nil
}

// dbus AccessPoint object.
type dbusObjAp struct {
	object dbus.BusObject
}

// AccessPoint return access point data as AccessPoint object.
func (ap dbusObjAp) AccessPoint() (AccessPoint, error) {
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
	go ex("maxBitrate", func() (any, error) { return ap.maxBitrate() })
	go ex("frequency", func() (any, error) { return ap.frequency() })
	go ex("hwAddress", func() (any, error) { return ap.hwAddress() })
	go ex("wpaFlags", func() (any, error) { return ap.wpaFlags() })
	go ex("rsnFlags", func() (any, error) { return ap.rsnFlags() })
	go ex("strength", func() (any, error) { return ap.strength() })
	go ex("flags", func() (any, error) { return ap.flags() })
	go ex("ssid", func() (any, error) { return ap.ssid() })
	go ex("mode", func() (any, error) { return ap.mode() })

	accPoint := &accessPoint{}
	for i := 0; i < fieldCount; i++ {
		res := <-ch
		if res.err != nil {
			return nil, res.err
		}

		switch res.field {
		case "maxBitrate":
			accPoint.maxBitrate = res.value.(uint32)
		case "frequency":
			accPoint.frequency = res.value.(uint32)
		case "hwAddress":
			accPoint.hwAddress = res.value.(string)
		case "wpaFlags":
			accPoint.wpaFlags = res.value.(uint32)
		case "rsnFlags":
			accPoint.rsnFlags = res.value.(uint32)
		case "strength":
			accPoint.strength = res.value.(uint8)
		case "flags":
			accPoint.flags = res.value.(uint32)
		case "ssid":
			accPoint.ssid = res.value.(string)
		case "mode":
			accPoint.mode = res.value.(string)
		}
	}

	return accPoint, nil
}

// run dbus call and return "Ssid" [name] for current AccessPoint path.
func (ap dbusObjAp) ssid() (string, error) {
	var val []byte

	if err := ap.object.
		Call(nmPropGet, 0, nmDestWirelessAccessPoint, "Ssid").
		Store(&val); err != nil {
		return "", err
	}
	return string(val), nil
}

// run dbus call and return "Strength" [level] for current AccessPoint path.
func (ap dbusObjAp) strength() (uint8, error) {
	var val uint8

	if err := ap.object.
		Call(nmPropGet, 0, nmDestWirelessAccessPoint, "Strength").
		Store(&val); err != nil {
		return 0, err
	}
	return val, nil
}

// run dbus call and return "Frequency" for current AccessPoint path.
func (ap dbusObjAp) frequency() (uint32, error) {
	var val uint32

	if err := ap.object.
		Call(nmPropGet, 0, nmDestWirelessAccessPoint, "Frequency").
		Store(&val); err != nil {
		return 0, err
	}
	return val, nil
}

// run dbus call and return "HwAddress" [MAC] for current AccessPoint path.
func (ap dbusObjAp) hwAddress() (string, error) {
	var val string

	if err := ap.object.
		Call(nmPropGet, 0, nmDestWirelessAccessPoint, "HwAddress").
		Store(&val); err != nil {
		return "", err
	}
	return val, nil
}

// run dbus call and return "Mode" [Ad-hoc, Access Point, Station, ...] for current AccessPoint path.
func (ap dbusObjAp) mode() (string, error) {
	var val string

	if err := ap.object.
		Call(nmPropGet, 0, nmDestWirelessAccessPoint, "Mode").
		Store(&val); err != nil {
		return "", err
	}
	return val, nil
}

// run dbus call and return "MaxBitrate" [kB/s] for current AccessPoint path.
func (ap dbusObjAp) maxBitrate() (uint32, error) {
	var val uint32

	if err := ap.object.
		Call(nmPropGet, 0, nmDestWirelessAccessPoint, "MaxBitrate").
		Store(&val); err != nil {
		return 0, err
	}
	return val, nil
}

// run dbus call and return "Flags" [WEP] for current AccessPoint path.
func (ap dbusObjAp) flags() (uint32, error) {
	var val uint32

	if err := ap.object.
		Call(nmPropGet, 0, nmDestWirelessAccessPoint, "Flags").
		Store(&val); err != nil {
		return 0, err
	}
	return val, nil
}

// run dbus call and return "WpaFlags" [WPA] for current AccessPoint path.
func (ap dbusObjAp) wpaFlags() (uint32, error) {
	var val uint32

	if err := ap.object.
		Call(nmPropGet, 0, nmDestWirelessAccessPoint, "WpaFlags").
		Store(&val); err != nil {
		return 0, err
	}
	return val, nil
}

// run dbus call and reurn "RsnFlags" [WPA2] for current AccessPoint path.
func (ap dbusObjAp) rsnFlags() (uint32, error) {
	var val uint32

	if err := ap.object.
		Call(nmPropGet, 0, nmDestWirelessAccessPoint, "RsnFlags").
		Store(&val); err != nil {
		return 0, err
	}
	return val, nil
}
