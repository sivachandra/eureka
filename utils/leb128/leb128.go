// #############################################################################
// This file is part of the "leb128" package of the "Eureka" project.
// It is distributed under the MIT License. Refer to the LICENSE file for more
// information.
//
// Website: http://www.github.com/sivachandra/eureka
// #############################################################################

// Package leb128 provides API to read LEB128 numbers from io.Reader objects.
package leb128

import (
	"bytes"
	"fmt"
	"io"
)

type LEB128 []byte

func ReadSigned(r io.ByteReader) (int64, error) {
	var res uint64 = 0
	var shift uint = 0
	var lastByte byte

	for true {
		b, err := r.ReadByte()
		if err != nil {
			return 0, fmt.Errorf("Error reading signed LEB128.\n%s", err.Error())
		}

		res |= uint64(b&0x7f) << shift

		lastByte = b
		shift += 7

		if 0x80&b == 0 {
			break
		}
	}

	if shift < 64 && (lastByte&0x40 != 0) {
		res |= 0xFFFFFFFFFFFFFFFF << shift
	}

	return int64(res), nil
}

func ReadUnsigned(r io.ByteReader) (uint64, error) {
	var res uint64 = 0
	var shift uint = 0

	for true {
		b, err := r.ReadByte()
		if err != nil {
			return 0, fmt.Errorf("Error reading unsigned LEB128.\n%s", err.Error())
		}

		res |= uint64(b&0x7f) << shift

		if 0x80&b == 0 {
			break
		}

		shift += 7
	}

	return res, nil
}

func Read(r io.ByteReader) (LEB128, error) {
	n := make([]byte, 0)

	for true {
		b, err := r.ReadByte()
		if err != nil {
			return LEB128(nil), err
		}

		n = append(n, b)

		if b & 0x80 == 0 {
			break
		}
	}

	return n, nil
}

func (n LEB128) AsSigned() (int64, error) {
	r := bytes.NewReader([]byte(n))
	return ReadSigned(r)
}

func (n LEB128) AsUnsigned() (uint64, error) {
	r := bytes.NewReader([]byte(n))
	return ReadUnsigned(r)
}

func Encode(v interface{}) (LEB128, error) {
	var valU64 uint64
	n := byte(0)
	s := false

	switch v.(type) {
	case int8:
		valU64 = uint64(int64(v.(int8)))
		n = 2
		s = true
	case uint8:
		valU64 = uint64(v.(uint8))
		n = 2
	case int16:
		valU64 = uint64(int64(v.(int16)))
		n = 3
		s = true
	case uint16:
		valU64 = uint64(v.(uint16))
		n = 3
	case int32:
		valU64 = uint64(int64(v.(int32)))
		n = 5
		s = true
	case uint32:
		valU64 = uint64(v.(uint32))
		n = 5
	case int64:
		valU64 = uint64(v.(int64))
		n = 10
		s = true
	case uint64:
		valU64 = v.(uint64)
		n = 10
	default:
		return nil, fmt.Errorf("Cannot encode value into an LEB128 number.")
	}

	out := make([]byte, 0)
	for k := byte(1); k <= n; k++ {
		b := byte(valU64 & 0x7F)
		if k != n {
			b |= 0x80
		}

		if k == 10 && s && b == 1 {
			// Sign extend the last byte of a b4-bit signed number.
			b |= 0x7F
		}

		out = append(out, b)
		valU64 = valU64 >> 7
	}

	return LEB128(out), nil
}
