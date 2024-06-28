package net

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

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
	nmDestWirelessReScan       = nmDest + ".Device.Wireless.RequestScan"
	nmDestActiveConnection     = nmDest + ".ActivateConnection"
	nmDestWirelessConnection   = nmDest + ".Settings.AddConnection"
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

	path   dbus.ObjectPath
	device dbus.ObjectPath
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
	return fmt.Sprintf("%d Mbps", ap.maxBitrate/1000)
}

// GetMacAddr return access point mac address.
func (ap accessPoint) GetMacAddr() string {
	return ap.hwAddress
}

func (ap accessPoint) GetMode() string {
	return ap.mode
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

// return hw address as hex bytes.
func (ap accessPoint) hwAddressHex() []byte {
	hw, _ := hex.DecodeString(
		strings.ReplaceAll(ap.hwAddress, ":", ""),
	)
	return hw
}

// dbus NetworkManager object.
type dbusNm struct {
	conn *dbus.Conn

	devicesWireless []dbus.BusObject
}

// New Driver return dbusNm as driver object.
func NewDriver() (driverInterface, error) {
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

// driverInterface method which send request for Re-Scan wireless networks. (work only from superuser).
func (nm dbusNm) wirelessScan() error {
	for _, device := range nm.devicesWireless {
		if err := device.
			Call(nmDestWirelessReScan, 0, map[string]dbus.Variant{}).
			Err; err != nil {
			return err
		}
	}

	return nil
}

// driverInterface method which return current wireless access point.
func (nm dbusNm) currentConnetcion() (AccessPoint, error) {
	return accessPoint{}, nil
}

// direverInterface method which ....
func (nm dbusNm) wirelessConnect(ssid, password string) error {
	ap, err := nm.getWirelessAccessPoint(ssid)
	if err != nil {
		return err
	}

	connCfg := map[string]map[string]dbus.Variant{
		"802-11-wireless": {
			"ssid":  dbus.MakeVariant([]byte(ap.ssid)),
			"bssid": dbus.MakeVariant(ap.hwAddressHex()),
		},
		"802-11-wireless-security": {
			"psk": dbus.MakeVariant(password),
		},
		"connection": {
			"type":        dbus.MakeVariant("802-11-wireless"),
			"id":          dbus.MakeVariant(ap.ssid),
			"autoconnect": dbus.MakeVariant(true),
		},
		"ipv4": {
			"method": dbus.MakeVariant("auto"),
		},
		"ipv6": {
			"method": dbus.MakeVariant("ignore"),
		},
	}
	switch ap.GetAccessType() {
	case "WPA2":
		connCfg["802-11-wireless-security"]["key-mgmt"] = dbus.MakeVariant("wpa-psk")
		connCfg["802-11-wireless-security"]["proto"] = dbus.MakeVariant([]string{"rsn"})
	case "WPA":
		connCfg["802-11-wireless-security"]["key-mgmt"] = dbus.MakeVariant("wpa-psk")
		connCfg["802-11-wireless-security"]["proto"] = dbus.MakeVariant([]string{"wpa"})
	case "WEP":
		connCfg["802-11-wireless-security"]["key-mgmt"] = dbus.MakeVariant("none")
		connCfg["802-11-wireless-security"]["wep-key0"] = dbus.MakeVariant(password)
		connCfg["802-11-wireless-security"]["wep-tx-keyidx"] = dbus.MakeVariant(0)
	case "OPEN":
		delete(connCfg, "802-11-wireless-security")
	}

	var (
		settings = nm.conn.Object(nmDest, "/org/freedesktop/NetworkManager/Settings")
		call     = settings.Call(nmDestWirelessConnection, 0, connCfg)
		path     dbus.ObjectPath
	)
	if call.Err != nil {
		return call.Err
	}
	if err := call.Store(&path); err != nil {
		return err
	}
	return nm.conn.
		Object(
			nmDest,
			"/org/freedesktop/NetworkManager",
		).
		Call(
			nmDestActiveConnection,
			0,
			path,
			ap.device,
			dbus.ObjectPath("/"),
		).Err
}

// driverInterface method which return available access points data as []AccessPoint object.
func (nm dbusNm) wirelessAccessPoints() ([]AccessPoint, error) {
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
			}.accessPoint()
			if err != nil {
				return nil, err
			}

			apObject.device = device.Path()
			apList = append(apList, apObject)
		}
	}

	return apList, nil
}

// return available access point from list by ssid or error.
func (nm dbusNm) getWirelessAccessPoint(ssid string) (*accessPoint, error) {
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
			}.accessPoint()
			apObject.device = device.Path()

			if err != nil {
				continue
			}
			if apObject.ssid == ssid {
				return apObject, nil
			}
		}
	}

	return nil, errors.New("access point not found")
}

// dbus AccessPoint object.
type dbusObjAp struct {
	object dbus.BusObject
}

// AccessPoint return access point data as AccessPoint object.
func (ap dbusObjAp) accessPoint() (*accessPoint, error) {
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

	accPoint := &accessPoint{
		path: ap.object.Path(),
	}
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
