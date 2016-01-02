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
	"fmt"
	"io"
)

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
