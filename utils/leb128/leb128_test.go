// #############################################################################
// This file is part of the "leb128" package of the "Eureka" project.
// It is distributed under the MIT License. Refer to the LICENSE file for more
// information.
//
// Website: http://www.github.com/sivachandra/eureka
// #############################################################################

package leb128

import (
	"bytes"
	"testing"
)

func TestReadSigned(t *testing.T) {
	b := []byte{0x9b, 0xf1, 0x59}
	r := bytes.NewReader(b)

	res, err := ReadSigned(r)
	if err != nil {
		t.Errorf("Error testing ReadSigned:\n%s", err.Error())
		return
	}
	if res != -624485 {
		t.Errorf("ReadSigned result wrong. Expected -624485, got %d", res)
		return
	}
}

func TestReadUnsigned(t *testing.T) {
	b := []byte{0xE5, 0x8E, 0x26}
	r := bytes.NewReader(b)

	res, err := ReadUnsigned(r)
	if err != nil {
		t.Errorf("Error testing ReadUnsigned:\n%s", err.Error())
		return
	}
	if res != 624485 {
		t.Errorf("ReadUnsigned result wrong. Expected 624485, got %d", res)
		return
	}
}

func TestReadLEB128Signed(t *testing.T) {
	b := []byte{0x9b, 0xf1, 0x59}
	r := bytes.NewReader(b)

	n, err := Read(r)
	if err != nil {
		t.Errorf("Error testing Read:\n%s", err.Error())
		return
	}
	res, err := n.AsSigned()
	if err != nil {
		t.Errorf("Error testing LEB128.AsSigned:\n%s", err.Error())
		return
	}
	if res != -624485 {
		t.Errorf("LEB128.AsSigned result wrong. Expected -624485, got %d", res)
		return
	}
}

func TestReadLEB128Unsigned(t *testing.T) {
	b := []byte{0xE5, 0x8E, 0x26}
	r := bytes.NewReader(b)

	n, err := Read(r)
	if err != nil {
		t.Errorf("Error testing Read:\n%s", err.Error())
		return
	}
	res, err := n.AsUnsigned()
	if err != nil {
		t.Errorf("Error testing LEB128.AsUnsigned:\n%s", err.Error())
		return
	}
	if res != 624485 {
		t.Errorf("LEB128.AsUnsigned result wrong. Expected 624485, got %d", res)
		return
	}
}

func TestEncode(t *testing.T) {
	val8 := int8(-5)
	valU8 := uint8(5)
	val16 := int16(-5)
	valU16 := uint16(5)
	val32 := int32(-5)
	valU32 := uint32(5)
	val64 := int64(-5)
	valU64 := uint64(5)

	l, err := Encode(val8)
	if err != nil {
		t.Errorf("Error encoding int8 number.\n%s", err.Error())
	}
	r, err := l.AsSigned()
	if err != nil {
		t.Errorf("Error reading signed number from encoded int8.\n%s", err.Error())
	}
	if r != -5 {
		t.Errorf("Wrong number read from encoded int8: %d", r)
		t.Error("Encoding: ", l)
	}

	l, err = Encode(valU8)
	if err != nil {
		t.Errorf("Error encoding uint8 number.\n%s", err.Error())
	}
	ru, err := l.AsUnsigned()
	if err != nil {
		t.Errorf("Error reading signed number from encoded uint8.\n%s", err.Error())
	}
	if ru != 5 {
		t.Errorf("Wrong number read from encoded uint8: %d", r)
		t.Error("Encoding: ", l)
	}

	l, err = Encode(val16)
	if err != nil {
		t.Errorf("Error encoding int16 number.\n%s", err.Error())
	}
	r, err = l.AsSigned()
	if err != nil {
		t.Errorf("Error reading signed number from encoded int16.\n%s", err.Error())
	}
	if r != -5 {
		t.Errorf("Wrong number read from encoded int16: %d", r)
		t.Error("Encoding: ", l)
	}

	l, err = Encode(valU16)
	if err != nil {
		t.Errorf("Error encoding uint16 number.\n%s", err.Error())
	}
	ru, err = l.AsUnsigned()
	if err != nil {
		t.Errorf("Error reading signed number from encoded uint16.\n%s", err.Error())
	}
	if ru != 5 {
		t.Errorf("Wrong number read from encoded uint16: %d", r)
		t.Error("Encoding: ", l)
	}

	l, err = Encode(val32)
	if err != nil {
		t.Errorf("Error encoding int32 number.\n%s", err.Error())
	}
	r, err = l.AsSigned()
	if err != nil {
		t.Errorf("Error reading signed number from encoded int32.\n%s", err.Error())
	}
	if r != -5 {
		t.Errorf("Wrong number read from encoded int32: %d", r)
		t.Error("Encoding: ", l)
	}

	l, err = Encode(valU32)
	if err != nil {
		t.Errorf("Error encoding uint32 number.\n%s", err.Error())
	}
	ru, err = l.AsUnsigned()
	if err != nil {
		t.Errorf("Error reading signed number from encoded uint32.\n%s", err.Error())
	}
	if ru != 5 {
		t.Errorf("Wrong number read from encoded uint32: %d", r)
		t.Error("Encoding: ", l)
	}

	l, err = Encode(val64)
	if err != nil {
		t.Errorf("Error encoding int64 number.\n%s", err.Error())
	}
	r, err = l.AsSigned()
	if err != nil {
		t.Errorf("Error reading signed number from encoded int64.\n%s", err.Error())
	}
	if r != -5 {
		t.Errorf("Wrong number read from encoded int64: %d", r)
		t.Error("Encoding: ", l)
	}

	l, err = Encode(valU64)
	if err != nil {
		t.Errorf("Error encoding uint64 number.\n%s", err.Error())
	}
	ru, err = l.AsUnsigned()
	if err != nil {
		t.Errorf("Error reading signed number from encoded uint64.\n%s", err.Error())
	}
	if ru != 5 {
		t.Errorf("Wrong number read from encoded uint64: %d", r)
		t.Error("Encoding: ", l)
	}

	valU64 = 0x8000000000000000
	l, err = Encode(valU64)
	if err != nil {
		t.Errorf("Error encoding uint64 number.\n%s", err.Error())
	}
	ru, err = l.AsUnsigned()
	if err != nil {
		t.Errorf("Error reading signed number from encoded uint64.\n%s", err.Error())
	}
	if ru != 0x8000000000000000 {
		t.Errorf("Wrong number read from encoded uint64: %d", r)
		t.Error("Encoding: ", l)
	}
}
