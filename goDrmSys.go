package goDrmSys

/*
#cgo pkg-config: libdrm
#include <stdlib.h>
#include <xf86drm.h>
int getNumDevices(void)
{
    return drmGetDevices(0, 0);
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

//constants mapping constants from libdrm
const (
	DRM_BUS_PCI      = 0
	DRM_BUS_USB      = 1
	DRM_BUS_PLATFORM = 2
	DRM_BUS_HOST1X   = 3
)

//DrmPciBusInfo maps to drmPciBusInfo of libdrm
type DrmPciBusInfo struct {
	Domain uint16
	Bus    uint8
	Dev    uint8
	Func   uint8
}

//DrmPciDeviceInfo maps to drmPciDeviceInfo of libdrm
type DrmPciDeviceInfo struct {
	VendorID    uint16
	DeviceID    uint16
	SubVendorID uint16
	SubDeviceID uint16
	RevisionID  uint8
}

//DrmUsbBusInfo maps to drmUsbBusInfo of libdrm
type DrmUsbBusInfo struct {
	Bus uint8
	Dev uint8
}

//DrmUsbDeviceInfo maps to drmUsbDeviceInfo of libdrm
type DrmUsbDeviceInfo struct {
	Vendor  uint16
	Product uint16
}

//DrmPlatformBusInfo maps to drmPlatformBusInfo of libdrm
type DrmPlatformBusInfo struct {
	FullName string
}

//DrmPlatformDeviceInfo maps to drmPlatformDeviceInfo of libdrm
type DrmPlatformDeviceInfo struct {
	Compatible **C.char //TODO: conversion
}

//DrmHost1xBusInfo maps to drmHost1xBusInfo of libdrm
type DrmHost1xBusInfo struct {
	FullName string
}

//DrmHost1xDeviceInfo maps to drmHost1xDeviceInfo of libdrm
type DrmHost1xDeviceInfo struct {
	Compatible **C.char //TODO: conversion
}

//DeviceNodes contains strings for Primary, Render, Contol part of GPU
type DeviceNodes struct {
	NodeMap map[string]string
}

//BusIn defines the interface for all types of BusInfo
type BusIn interface {
	GetBusInfo() string
}

//GetBusInfo for PciInfo
func (b DrmPciBusInfo) GetBusInfo() string {
	return (fmt.Sprintf("%d:%d:%d.%d", b.Domain, b.Bus, b.Dev, b.Func))
}

//GetBusInfo for UsbInfo
func (b DrmUsbBusInfo) GetBusInfo() string {
	return (fmt.Sprintf("%d:%d", b.Bus, b.Dev))
}

//GetBusInfo for PlatformInfo
func (b DrmPlatformBusInfo) GetBusInfo() string {
	return (fmt.Sprintf("%s", b.FullName))
}

//GetBusInfo for Host1xInfo
func (b DrmHost1xBusInfo) GetBusInfo() string {
	return (fmt.Sprintf("%s", b.FullName))
}

//DevIn defines the interface for all types of DeviceInfo
type DevIn interface {
	GetDevInfo() string
}

//GetDevInfo for PciInfo
func (b DrmPciDeviceInfo) GetDevInfo() string {
	return (fmt.Sprintf("0x%X:0x%X:0x%X:0x%X:%d", b.VendorID, b.DeviceID,
		b.SubVendorID, b.SubDeviceID, b.RevisionID))
}

//GetDevInfo for UsbInfo
func (b DrmUsbDeviceInfo) GetDevInfo() string {
	return (fmt.Sprintf("0x%X:0x%X", b.Vendor, b.Product))
}

//GetDevInfo for PlatformInfo
func (b DrmPlatformDeviceInfo) GetDevInfo() string {
	//add
	/*list := []string{}
	clist := (*[1 << 30]*C.char)(unsafe.Pointer(b.Compatible))[:C.DRM_NODE_MAX]
	for i := 0; i < len(clist); i++ {
		list = append(list, C.GoString(clist[i]))
	}
	ret := strings.Join(list, ":")*/
	return ""
}

//GetDevInfo for Host1xInfo
func (b DrmHost1xDeviceInfo) GetDevInfo() string {
	//add
	/*list := []string{}
	clist := (*[1 << 30]*C.char)(unsafe.Pointer(b.Compatible))[:C.DRM_NODE_MAX]
	for i := 0; i < len(clist); i++ {
		list = append(list, C.GoString(clist[i]))
	}
	ret := strings.Join(list, ":")*/
	return ""
}

//DeviceInfo contains BusInfo and Device info fields of drmDevice
type DeviceInfo struct {
	BusInfo BusIn
	DevInfo DevIn
}

//Device contains all info from libdrm struct drmDevice
type Device struct {
	Nodes DeviceNodes
	Info  DeviceInfo
}

func (nodes *DeviceNodes) fromSys(d C.drmDevicePtr) {
	nodes.NodeMap = make(map[string]string)
	nds := (*[1 << 30]*C.char)(unsafe.Pointer(d.nodes))[:C.DRM_NODE_MAX]

	if (d.available_nodes & (1 << C.DRM_NODE_PRIMARY)) != 0 {
		nodes.NodeMap["primary"] = C.GoString(nds[C.DRM_NODE_PRIMARY])
		//nodes.Primary = C.GoString(nds[C.DRM_NODE_PRIMARY])
	}

	if (d.available_nodes & (1 << C.DRM_NODE_CONTROL)) != 0 {
		nodes.NodeMap["control"] = C.GoString(nds[C.DRM_NODE_CONTROL])
		//nodes.Control = C.GoString(nds[C.DRM_NODE_CONTROL])
	}

	if (d.available_nodes & (1 << C.DRM_NODE_RENDER)) != 0 {
		nodes.NodeMap["render"] = C.GoString(nds[C.DRM_NODE_RENDER])
		//nodes.Render = C.GoString(nds[C.DRM_NODE_RENDER])
	}

}

func (info *DeviceInfo) fromSys(d C.drmDevicePtr) {
	switch int(d.bustype) {

	case DRM_BUS_PCI:
		info.BusInfo = *(*(*DrmPciBusInfo))(unsafe.Pointer(&d.businfo))
		info.DevInfo = *(*(*DrmPciDeviceInfo))(unsafe.Pointer(&d.deviceinfo))

	case DRM_BUS_USB:
		info.BusInfo = *(*(*DrmUsbBusInfo))(unsafe.Pointer(&d.businfo))
		info.DevInfo = *(*(*DrmUsbDeviceInfo))(unsafe.Pointer(&d.deviceinfo))

	case DRM_BUS_PLATFORM:
		info.BusInfo = *(*(*DrmPlatformBusInfo))(unsafe.Pointer(&d.businfo))
		info.DevInfo = *(*(*DrmPlatformDeviceInfo))(unsafe.Pointer(&d.deviceinfo))

	case DRM_BUS_HOST1X:
		info.BusInfo = *(*(*DrmHost1xBusInfo))(unsafe.Pointer(&d.businfo))
		info.DevInfo = *(*(*DrmHost1xDeviceInfo))(unsafe.Pointer(&d.deviceinfo))
	}
}

//DrmAvailable checks whether DRM is available
func DrmAvailable() bool {

	return (C.drmAvailable() != 0)
}

//GetDevices returns array containing info for all devices
func GetDevices() []Device {

	n := C.getNumDevices()
	d := make([]C.drmDevicePtr, n)
	C.drmGetDevices(&d[0], n)

	devices := make([]Device, n)
	for i := 0; i < int(n); i++ {
		devices[i].Nodes.fromSys(d[i])
		devices[i].Info.fromSys(d[i])
	}

	C.drmFreeDevices(&d[0], n)

	return devices

}
