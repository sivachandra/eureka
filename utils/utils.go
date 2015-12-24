// #############################################################################
// This file is part of the "utils" package of the "Eureka" project.
// It is distributed under the MIT License. Refer to the LICENSE file for more
// information.
//
// Website: http://www.github.com/sivachandra/eureka
// #############################################################################

// Package utils provides a utility API.
package utils

import (
	"io"
)

func ReadUntil(r io.ByteReader, delim byte) ([]byte, error) {
	var str []byte

	for true {
		c, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		if c == delim {
			break
		}
		str = append(str, c)
	}

	return str, nil
}

