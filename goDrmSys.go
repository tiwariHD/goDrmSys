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
    "unsafe"
    "fmt"
)

const (
    DRM_BUS_PCI         = 0
    DRM_BUS_USB         = 1
    DRM_BUS_PLATFORM    = 2
    DRM_BUS_HOST1X      = 3
)

type DrmPciBusInfo struct {
    Domain      uint16
    Bus         uint8
    Dev         uint8
    Func        uint8
}

type DrmPciDeviceInfo struct {
    VendorId    uint16
    DeviceId    uint16
    SubVendorid uint16
    SubDeviceId uint16
    RevisionId  uint8
}

type DrmUsbBusInfo struct {
    Bus         uint8
    Dev         uint8
}

type DrmUsbDeviceInfo struct {
    Vendor      uint16
    Product     uint16
}

type DrmPlatformBusInfo struct {
    FullName    string
}

type DrmPlatformDeviceInfo struct {
    Compatible  **C.char    //TODO: conversion
}

type DrmHost1xBusInfo struct {
    FullName    string
}

type DrmHost1xDeviceInfo struct {
    Compatible  **C.char    //TODO: conversion
}

type DeviceNodes struct {
    Primary     string
    Control     string
    Render      string
}

type BusIn interface {
    GetBusInfo() string
}

func (b DrmPciBusInfo) GetBusInfo() string {
    return (fmt.Sprintf("%d:%d:%d.%d", b.Domain, b.Bus, b.Dev, b.Func))
}

func (b DrmUsbBusInfo) GetBusInfo() string {
    return (fmt.Sprintf("%d:%d", b.Bus, b.Dev))
}

func (b DrmPlatformBusInfo) GetBusInfo() string {
    return (fmt.Sprintf("%s", b.FullName))
}

func (b DrmHost1xBusInfo) GetBusInfo() string {
    return (fmt.Sprintf("%s", b.FullName))
}

type DevIn interface {
    GetDevInfo()
}

func (b DrmPciDeviceInfo) GetDevInfo() {
    //add
}

func (b DrmUsbDeviceInfo) GetDevInfo() {
    //add
}

func (b DrmPlatformDeviceInfo) GetDevInfo() {
    //add
}

func (b DrmHost1xDeviceInfo) GetDevInfo() {
    //add
}

type DeviceInfo struct {
    BusInfo     BusIn
    DevInfo     DevIn
}

type Device struct {
    Nodes   DeviceNodes
    Info    DeviceInfo
}

func (nodes *DeviceNodes) from_sys (d C.drmDevicePtr) {
    nds := (*[1<<30]*C.char)(unsafe.Pointer(d.nodes))[:C.DRM_NODE_MAX]

    if (d.available_nodes & (1 << C.DRM_NODE_PRIMARY)) != 0 {
        nodes.Primary = C.GoString(nds[C.DRM_NODE_PRIMARY])
    }

    if (d.available_nodes & (1 << C.DRM_NODE_CONTROL)) != 0 {
        nodes.Control = C.GoString(nds[C.DRM_NODE_CONTROL])
    }

    if (d.available_nodes & (1 << C.DRM_NODE_RENDER)) != 0 {
        nodes.Render = C.GoString(nds[C.DRM_NODE_RENDER])
    }

}

func (info *DeviceInfo) from_sys (d C.drmDevicePtr) {
    switch int(d.bustype) {

    case DRM_BUS_PCI :
        info.BusInfo = *(*(*DrmPciBusInfo))(unsafe.Pointer(&d.businfo))
        info.DevInfo = *(*(*DrmPciDeviceInfo))(unsafe.Pointer(&d.deviceinfo))

    case DRM_BUS_USB :
        info.BusInfo = *(*(*DrmUsbBusInfo))(unsafe.Pointer(&d.businfo))
        info.DevInfo = *(*(*DrmUsbDeviceInfo))(unsafe.Pointer(&d.deviceinfo))

    case DRM_BUS_PLATFORM :
        info.BusInfo = *(*(*DrmPlatformBusInfo))(unsafe.Pointer(&d.businfo))
        info.DevInfo = *(*(*DrmPlatformDeviceInfo))(unsafe.Pointer(&d.deviceinfo))

    case DRM_BUS_HOST1X :
        info.BusInfo = *(*(*DrmHost1xBusInfo))(unsafe.Pointer(&d.businfo))
        info.DevInfo = *(*(*DrmHost1xDeviceInfo))(unsafe.Pointer(&d.deviceinfo))
    }
}

func DrmAvailable() bool {

    return (C.drmAvailable() != 0)
}

func GetDevices() []Device {

    n := C.getNumDevices()
    d := make([]C.drmDevicePtr, n)
    C.drmGetDevices(&d[0], n)

    devices := make([]Device, n)
    for i := 0; i < int(n); i++ {
        devices[i].Nodes.from_sys(d[i])
        devices[i].Info.from_sys(d[i])
    }

    C.drmFreeDevices(&d[0], n)

    return devices

}

