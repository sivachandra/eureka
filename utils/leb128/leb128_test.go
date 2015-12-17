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
	}
	if res != -624485 {
		t.Errorf("ReadSigned result wrong. Expected 624485, got %d", res)
	}
}

func TestReadUnsigned(t *testing.T) {
	b := []byte{0xE5, 0x8E, 0x26}
	r := bytes.NewReader(b)

	res, err := ReadUnsigned(r)
	if err != nil {
		t.Errorf("Error testing ReadUnsigned:\n%s", err.Error())
	}
	if res != 624485 {
		t.Errorf("ReadUnsigned result wrong. Expected 624485, got %d", res)
	}
}
