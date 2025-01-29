package ch16

import (
	"unsafe"
)

type OrderInfoBig struct {
	OrderCode   rune
	Amount      int
	OrderNumber uint16
	Items       []string
	IsReady     bool
}

type OrderInfoSmall struct {
	IsReady     bool
	OrderNumber uint16
	OrderCode   rune
	Amount      int
	Items       []string
}

type StructMemInfo struct {
	SizeOf       uintptr
	FieldOffSets map[string]uintptr
}

func GetMemInfoBigUnsafe(oi OrderInfoBig) (info StructMemInfo) {
	info.SizeOf = unsafe.Sizeof(oi)
	info.FieldOffSets = make(map[string]uintptr, 5)
	info.FieldOffSets["OrderCode"] = unsafe.Offsetof(oi.OrderCode)
	info.FieldOffSets["Amount"] = unsafe.Offsetof(oi.Amount)
	info.FieldOffSets["OrderNumber"] = unsafe.Offsetof(oi.OrderNumber)
	info.FieldOffSets["Items"] = unsafe.Offsetof(oi.Items)
	info.FieldOffSets["IsReady"] = unsafe.Offsetof(oi.IsReady)
	return info
}

func GetMemInfoSmallUnsafe(oi OrderInfoSmall) (info StructMemInfo) {
	info.SizeOf = unsafe.Sizeof(oi)
	info.FieldOffSets = make(map[string]uintptr, 5)
	info.FieldOffSets["IsReady"] = unsafe.Offsetof(oi.IsReady)
	info.FieldOffSets["OrderNumber"] = unsafe.Offsetof(oi.OrderNumber)
	info.FieldOffSets["OrderCode"] = unsafe.Offsetof(oi.OrderCode)
	info.FieldOffSets["Amount"] = unsafe.Offsetof(oi.Amount)
	info.FieldOffSets["Items"] = unsafe.Offsetof(oi.Items)
	return info
}
