// #############################################################################
// This file is part of the "golf" package of the "Eureka" project.
// It is distributed under the MIT License. Refer to the LICENSE file for more
// information.
//
// Website: http://www.github.com/sivachandra/eureka
// #############################################################################

package golf

import (
	"encoding/binary"
	"fmt"
	"os"
)

// Set of constants which specify the type of segment in a program/segment
// header.
const (
	SegTypeNull              = uint32(0)
	SegTypeLoad              = uint32(1)
	SegTypeDynamic           = uint32(2)
	SegTypeInterp            = uint32(3)
	SegTypeNote              = uint32(4)
	SegTypeReserved          = uint32(5)
	SegTypeProgHdr           = uint32(6)
	SegTypeTLS               = uint32(7)
	SegTypeNumDefinedTypes   = uint32(8)
	SegTypeStartOSSpecific   = uint32(0x60000000)
	SegTypeGnuEHFrame        = uint32(0x6474e550)
	SegTypeGnuStack          = uint32(0x6474e551)
	SegTypeGnuRelRO          = uint32(0x6474e552)
	SegTypeSunOSStart        = uint32(0x6ffffffa)
	SegTypeSunWBSS           = uint32(0x6ffffffa)
	SegTypeSunWStack         = uint32(0x6ffffffb)
	SegTypeSunOSEnd          = uint32(0x6fffffff)
	SegTypeEndOSSpecific     = uint32(0x6fffffff)
	SegTypeStartProcSpecific = uint32(0x70000000)
	SegTypeEndProcSpecific   = uint32(0x7fffffff)
)

// Set of constants which represent the segment flags.
// A segment header can include a flag generated from more than one
// of these values using the '|' operator. For example, a flag could
// be generated as the result of SegFlagsExecutable | SegFlagsReadable.
const (
	SegFlagsExecutable       = uint32(1 << 0)
	SegFlagsWritable         = uint32(1 << 1)
	SegFlagsReadable         = uint32(1 << 2)
	SegFlagsOSSpecificMask   = uint32(0x0ff00000)
	SegFlagsProcSpecificMask = uint32(0xf0000000)
)

// A SegHdr represents the entry in the program header table of an ELF file.
type SegHdr interface {
	// Returns the class of the segment, which is the class of the ELF file to
	// which it belongs.
	Class() ELFClass

	// Returns the type of the sgement,
	Type() uint32

	// Returns the offset of the segment in the ELF file.
	Offset() uint64

	// Returns the virtual address of the segment.
	VirtualAddress() uint64

	// Returns the physical address of the segment.
	PhysicalAddress() uint64

	// Returns the size of the segment in the ELF file.
	FileSize() uint64

	// Returns the size of the segment in the memory.
	MemSize() uint64

	// Returns the flags for the segment.
	Flags() uint32

	// Returns the alignment of the segment.
	Alignment() uint64
}

type segHdr32 struct {
	diskData struct {
		Type            uint32
		Offset          uint32
		VirtualAddress  uint32
		PhysicalAddress uint32
		FileSize        uint32
		MemSize         uint32
		Flags           uint32
		Alignment       uint32
	}
}

func (hdr *segHdr32) Class() ELFClass {
	return Class32
}

func (hdr *segHdr32) Type() uint32 {
	return hdr.diskData.Type
}

func (hdr *segHdr32) Offset() uint64 {
	return uint64(hdr.diskData.Offset)
}

func (hdr *segHdr32) VirtualAddress() uint64 {
	return uint64(hdr.diskData.VirtualAddress)
}

func (hdr *segHdr32) PhysicalAddress() uint64 {
	return uint64(hdr.diskData.PhysicalAddress)
}

func (hdr *segHdr32) FileSize() uint64 {
	return uint64(hdr.diskData.FileSize)
}

func (hdr *segHdr32) MemSize() uint64 {
	return uint64(hdr.diskData.MemSize)
}

func (hdr *segHdr32) Flags() uint32 {
	return hdr.diskData.Flags
}

func (hdr *segHdr32) Alignment() uint64 {
	return uint64(hdr.diskData.Alignment)
}

type segHdr64 struct {
	diskData struct {
		Type            uint32
		Flags           uint32
		Offset          uint64
		VirtualAddress  uint64
		PhysicalAddress uint64
		FileSize        uint64
		MemSize         uint64
		Alignment       uint64
	}
}

func (hdr *segHdr64) Class() ELFClass {
	return Class64
}

func (hdr *segHdr64) Type() uint32 {
	return hdr.diskData.Type
}

func (hdr *segHdr64) Offset() uint64 {
	return hdr.diskData.Offset
}

func (hdr *segHdr64) VirtualAddress() uint64 {
	return hdr.diskData.VirtualAddress
}

func (hdr *segHdr64) PhysicalAddress() uint64 {
	return hdr.diskData.PhysicalAddress
}

func (hdr *segHdr64) FileSize() uint64 {
	return hdr.diskData.FileSize
}

func (hdr *segHdr64) MemSize() uint64 {
	return hdr.diskData.MemSize
}

func (hdr *segHdr64) Flags() uint32 {
	return hdr.diskData.Flags
}

func (hdr *segHdr64) Alignment() uint64 {
	return hdr.diskData.Alignment
}

func readSegHdrTbl(file *os.File, header ELFHeader) ([]SegHdr, error) {
	_, err := file.Seek(int64(header.ProgHdrTblOffset()), 0)
	if err != nil {
		err = fmt.Errorf(
			"Unable to seek to the program header table in '%s'.\n%s", file.Name(), err.Error())
		return nil, err
	}

	var segHdrTbl []SegHdr
	for i := uint16(0); i < header.ProgHdrCount(); i++ {
		endianess := header.ELFIdent().Endianess
		var hdr SegHdr
		if header.ELFIdent().Class == Class32 {
			hdr32 := new(segHdr32)
			err = binary.Read(file, endianMap[endianess], &hdr32.diskData)
			hdr = hdr32
		} else {
			hdr64 := new(segHdr64)
			err = binary.Read(file, endianMap[endianess], &hdr64.diskData)
			hdr = hdr64
		}

		if err != nil {
			err = fmt.Errorf(
				"Error reading segment header from '%s'.\n%s", file.Name(), err.Error())
			return nil, err
		}

		segHdrTbl = append(segHdrTbl, hdr)
	}

	return segHdrTbl, nil
}
