package ch16

import (
	"reflect"
	"testing"
)

func GetMemInfoReflect(v any) (info StructMemInfo, err error) {
	// Check if input type is structure.
	vt := reflect.TypeOf(v)
	if vt.Kind() != reflect.Struct {
		return info, ErrNonStruct
	}

	// Check each field in a structure.
	info.SizeOf = vt.Size()
	info.FieldOffSets = make(map[string]uintptr, vt.NumField())

	for i := 0; i < vt.NumField(); i += 1 {
		field := vt.Field(i)
		info.FieldOffSets[field.Name] = vt.Field(i).Offset
	}

	return info, nil
}

func TestGetMemInfo(t *testing.T) {
	// Test big.
	oib := OrderInfoBig{}
	oib_mi_unsafe := GetMemInfoBigUnsafe(oib)
	oib_mi_reflect, err := GetMemInfoReflect(oib)
	if err != nil {
		t.Fatal("getting big memory info using reflection failed")
	}

	if !reflect.DeepEqual(oib_mi_unsafe, oib_mi_reflect) {
		t.Error("unsafe and reflection big memory info mismatch")
	}

	// Test small.
	ois := OrderInfoSmall{}
	ois_mi_unsafe := GetMemInfoSmallUnsafe(ois)
	ois_mi_reflect, err := GetMemInfoReflect(ois)
	if err != nil {
		t.Fatal("getting small memory info using reflection failed")
	}

	if !reflect.DeepEqual(ois_mi_unsafe, ois_mi_reflect) {
		t.Error("unsafe and reflection small memory info mismatch")
	}

	// Check sizes.
	if oib_mi_unsafe.SizeOf < ois_mi_unsafe.SizeOf {
		t.Error("big order info smaller than small order info")
	}
}
