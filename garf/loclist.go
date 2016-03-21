// #############################################################################
// This file is part of the "garf" package of the "Eureka" project.
// It is distributed under the MIT License. Refer to the LICENSE file for more
// information.
//
// Website: http://www.github.com/sivachandra/eureka
// #############################################################################

// Package garf provides API to read DWARF debug info from ELF files.
package garf

import (
	"encoding/binary"
	"fmt"
	"math"
)

func (d *DwData) readLocList(u *DwUnit, offset uint64, en binary.ByteOrder) (LocList, error) {
	sectMap := d.elf.SectMap()
	s, exists := sectMap[".debug_loc"]
	if !exists {
		return nil, fmt.Errorf(".debug_loc section missing in ELF data.")
	}

	r, err := s[0].NewReader()
	if err != nil {
		return nil, fmt.Errorf("Error creating .debug_loc section reader.\n%", err.Error())
	}

	_, err = r.Seek(int64(offset), 0)
	if err != nil {
		err = fmt.Errorf(
			"Unable to seek the loc list offset in .debug_loc.\n%s", err.Error())
		return nil, err
	}

	addressSize := d.elf.AddressSize()
	var locList LocList
	for {
		var begin, end uint64
		if addressSize == 4 {
			var begin32 uint32
			err = binary.Read(r, en, &begin32)
			if err != nil {
				break
			}

			if begin32 == math.MaxUint32 {
				begin = math.MaxUint64
			} else {
				begin = uint64(begin32)
			}

			var end32 uint32
			err = binary.Read(r, en, &end32)
			if err != nil {
				break
			}

			if end32 == math.MaxUint32 {
				end = math.MaxUint64
			} else {
				end = uint64(begin32)
			}
		} else {
			err = binary.Read(r, en, &begin)
			if err != nil {
				break
			}

			err = binary.Read(r, en, &end)
			if err != nil {
				break
			}
		}

		if begin == 0 && end == 0 {
			// End of list entry
			var entry EndOfListLocListEntry
			locList = append(locList, entry)
			break
		} else if begin == math.MaxUint64 {
			// Base address selection entry
			locList = append(locList, BaseAddrSelectionLocListEntry(end))
		} else if begin == 0 && end == math.MaxUint64 {
			// Default loc list entry
			var expr DwExpr
			expr, err = d.readSizeAndDwExpr(u, r, en)
			if err != nil {
				err = fmt.Errorf(
					"Error reading DWARF expr from default loc list entry.\n%s",
					err.Error())
				break
			}

			locList = append(locList, DefaultLocListEntry(expr))
		} else {
			// Normal loc list entry
			var size uint16
			err = binary.Read(r, en, &size)
			if err != nil {
				err = fmt.Errorf(
					"Error reading size of normal loc list entry.\n%s",
					err.Error())
				break
			}

			var expr DwExpr
			expr, err = d.readDwExpr(u, r, en, uint64(size))
			if err != nil {
				err = fmt.Errorf(
					"Error reading DWARF expr from normal loc list entry.\n%s",
					err.Error())
				break
			}

			locList = append(locList, NormalLocListEntry{begin, end, expr})
		}
	}

	if err != nil {
		return nil, err
	}

	return locList, nil
}
