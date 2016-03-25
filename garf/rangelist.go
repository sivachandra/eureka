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

type RangeListEntryType uint8

const (
	RangeListEntryTypeNormal            = RangeListEntryType(1)
	RangeListEntryTypeBaseAddrSelection = RangeListEntryType(2)
	RangeListEntryTypeEndOfList         = RangeListEntryType(3)
)

type RangeListEntry interface {
	RangeListEntryType() RangeListEntryType
}

type RangeListEntryNormal struct {
	Begin uint64
	End   uint64
}

func (e RangeListEntryNormal) RangeListEntryType() RangeListEntryType {
	return RangeListEntryTypeNormal
}

type RangeListEntryBaseAddrSelection uint64

func (e RangeListEntryBaseAddrSelection) RangeListEntryType() RangeListEntryType {
	return RangeListEntryTypeBaseAddrSelection
}

type RangeListEntryEndOfList struct {
}

func (e RangeListEntryEndOfList) RangeListEntryType() RangeListEntryType {
	return RangeListEntryTypeEndOfList
}

type RangeList []RangeListEntry

func (d *DwData) readRangeList(u *DwUnit, offset uint64, en binary.ByteOrder) (RangeList, error) {
	sectMap := d.elf.SectMap()
	s, exists := sectMap[".debug_ranges"]
	if !exists {
		return nil, fmt.Errorf(".debug_ranges section missing in ELF data.")
	}

	r, err := s[0].NewReader()
	if err != nil {
		err = fmt.Errorf("Error creating .debug_ranges section reader.\n%s", err.Error())
		return nil, err
	}

	_, err = r.Seek(int64(offset), 0)
	if err != nil {
		err = fmt.Errorf(
			"Unable to seek the loc list offset in .debug_loc.\n%s", err.Error())
		return nil, err
	}

	addressSize := d.elf.AddressSize()
	var rangeList RangeList
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
			var entry RangeListEntryEndOfList
			rangeList = append(rangeList, entry)
			break
		} else if begin == math.MaxUint64 {
			// Base address selection entry
			rangeList = append(rangeList, RangeListEntryBaseAddrSelection(end))
		} else {
			rangeList = append(rangeList, RangeListEntryNormal{begin, end})
		}
	}

	if err != nil {
		return nil, err
	}

	return rangeList, nil
}
