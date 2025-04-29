package models

import "gorm.io/gorm"

type OperationMode string
type DeviceType string

const (
	OperationModeUltraWideBand OperationMode = "ULTRA_WIDE_BAND"
	OperationModeWifi          OperationMode = "WIFI"
	DeviceTypeAnchor           DeviceType    = "ANCHOR"
	DeviceTypeTag              DeviceType    = "TAG"
)

type StationConfiguration struct {
	gorm.Model
	OperationMode
}
