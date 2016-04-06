///////////////////////////////////////////////////////////////////////////
// Copyright 2016 Siva Chandra
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
///////////////////////////////////////////////////////////////////////////

package golf

import (
	"encoding/binary"
	"fmt"
	"os"
)

// Values of type ELFClass represent the class (32-bit or 64-bit) of an ELF
// file.
type ELFClass byte

// Values of type ELFEndianess represent the endianess of an ELF file.
type ELFEndianess byte

// Value of type OSABI represent the operating system ABI used in an ELF file.
type OSABI byte

type ELFIdent struct {
	// The magic number
	MagicNumber [4]byte

	// The class of the ELF file, 32-bit or 64-bit.
	// It takes one of the constant values Class32 or Class64.
	Class ELFClass

	// The endianess of the file.
	// It takes one of the constant values LittleEndian or BigEndian.
	Endianess ELFEndianess

	// The ELF version of the file.
	ELFVersion byte

	// The operating system ABI used in the file.
	// It takes one of the constant value ABI*.
	ABI OSABI

	// The version of the ABI used.
	ABIVersion byte

	// Unused padding.
	Padding [7]byte
}

const (
	Mag0 byte = 0x7F
	Mag1 byte = 'E'
	Mag2 byte = 'L'
	Mag3 byte = 'F'
)

const (
	// Invalid class.
	ClassNone ELFClass = ELFClass(0)

	// Class of a 32-bit ELF file.
	Class32 ELFClass = ELFClass(1)

	// Class of a 64-bit ELF file.
	Class64 ELFClass = ELFClass(2)
)

const (
	// LittleEndian denotes that the endianess used in an ELF file is
	// little endian.
	LittleEndian ELFEndianess = ELFEndianess(1)

	// BigEndian denotes that the endianess used in an ELF file is
	// bit endian.
	BigEndian ELFEndianess = ELFEndianess(2)
)

var endianMap = map[ELFEndianess]binary.ByteOrder{
	LittleEndian: binary.LittleEndian,
	BigEndian:    binary.BigEndian,
}

const (
	ABINone       OSABI = 0
	ABISystemV    OSABI = OSABI(0)
	ABIHPUX       OSABI = OSABI(1)
	ABINetBSD     OSABI = OSABI(2)
	ABIGnu        OSABI = OSABI(3)
	ABILinux      OSABI = OSABI(3)
	ABISolaris    OSABI = OSABI(6)
	ABIAIX        OSABI = OSABI(7)
	ABIIRIX       OSABI = OSABI(8)
	ABIFreeBSD    OSABI = OSABI(9)
	ABITru64      OSABI = OSABI(10)
	ABIModesto    OSABI = OSABI(11)
	ABIOpenBSD    OSABI = OSABI(12)
	ABIArmAEABI   OSABI = OSABI(64)
	ABIArm        OSABI = OSABI(97)
	ABIStandalone OSABI = OSABI(255)
)

// ELFType values denote the type of the ELF file.
// For example a value of 'TypeExecutable' specifies that the ELF file
// is an executable file.
type ELFType uint16

// MachineArch values denote the machine architecture (or the ISA) used
// in the ELF file.
type MachineArch uint16

type ELFHeader interface {
	ELFIdent() *ELFIdent
	Type() ELFType
	Machine() MachineArch
	Version() uint32
	EntryPoint() uint64
	ProgHdrTblOffset() uint64
	SectHdrTblOffset() uint64
	Flags() uint32
	HeaderSize() uint16
	ProgHdrTblEntrySize() uint16
	ProgHdrCount() uint16
	SectHdrTblEntrySize() uint16
	SectHdrCount() uint16
	StrTblIndex() uint16
}

const (
	TypeNone              ELFType = ELFType(0)
	TypeRelocatable       ELFType = ELFType(1)
	TypeExecutable        ELFType = ELFType(2)
	TypeShared            ELFType = ELFType(3)
	TypeCore              ELFType = ELFType(4)
	TypeCountDefined      ELFType = ELFType(5)
	TypeStartOSSpecific   ELFType = ELFType(0xfe00)
	TypeEndOSSpecific     ELFType = ELFType(0xfeff)
	TypeStartProcSpecific ELFType = ELFType(0xff00)
	TypeEndProcSpecific   ELFType = ELFType(0xffff)
)

const (
	MachineSPARC   MachineArch = MachineArch(0x02)
	MachineX86     MachineArch = MachineArch(0x03)
	MachineMIPS    MachineArch = MachineArch(0x08)
	MachinePowerPC MachineArch = MachineArch(0x14)
	MachineARM     MachineArch = MachineArch(0x28)
	MachineSuperH  MachineArch = MachineArch(0x2A)
	MachineIA64    MachineArch = MachineArch(0x32)
	MachineX86_64  MachineArch = MachineArch(0x3E)
	MachineAArch64 MachineArch = MachineArch(0xB7)
)

type header32 struct {
	// The struct value capturing the ELF file indentifier.
	ident ELFIdent

	platformSpecific struct {
		// The file type.
		// It takes one of the constant values TypeNone, TypeRelocatable,
		// TypeExecutable, TypeShared, or TypeCore
		Type ELFType

		// The machine type of the instruction set architecture used in the
		// file.
		Machine MachineArch

		Version          uint32 // The elf version
		EntryPoint       uint32
		ProgHdrTblOffset uint32
		SectHdrTblOffset uint32

		Flags uint32

		// The size of the ELF header table as on disk.
		// Note, it is not the byte size of this struct.
		HeaderSize uint16

		ProgHdrTblEntrySize  uint16
		ProgHdrTblEntryCount uint16
		SectHdrTblEntrySize  uint16
		SectHdrTblEntryCount uint16
		StrTblIndex          uint16
	}
}

func (header *header32) ELFIdent() *ELFIdent {
	return &header.ident
}

func (header *header32) Type() ELFType {
	return header.platformSpecific.Type
}

func (header *header32) Machine() MachineArch {
	return header.platformSpecific.Machine
}

func (header *header32) Version() uint32 {
	return header.platformSpecific.Version
}

func (header *header32) EntryPoint() uint64 {
	return uint64(header.platformSpecific.EntryPoint)
}

func (header *header32) ProgHdrTblOffset() uint64 {
	return uint64(header.platformSpecific.ProgHdrTblOffset)
}

func (header *header32) SectHdrTblOffset() uint64 {
	return uint64(header.platformSpecific.SectHdrTblOffset)
}

func (header *header32) Flags() uint32 {
	return header.platformSpecific.Flags
}

func (header *header32) HeaderSize() uint16 {
	return header.platformSpecific.HeaderSize
}

func (header *header32) ProgHdrTblEntrySize() uint16 {
	return header.platformSpecific.ProgHdrTblEntrySize
}

func (header *header32) ProgHdrCount() uint16 {
	return header.platformSpecific.ProgHdrTblEntryCount
}

func (header *header32) SectHdrTblEntrySize() uint16 {
	return header.platformSpecific.SectHdrTblEntrySize
}

func (header *header32) SectHdrCount() uint16 {
	return header.platformSpecific.SectHdrTblEntryCount
}

func (header *header32) StrTblIndex() uint16 {
	return header.platformSpecific.StrTblIndex
}

type header64 struct {
	// The struct value capturing the ELF file indentifier.
	ident ELFIdent

	platformSpecific struct {
		// The file type.
		// It takes one of the constant values TypeNone, TypeRelocatable,
		// TypeExecutable, TypeShared, or TypeCore
		Type ELFType

		// The machine type of the instruction set architecture used in the
		// file.
		Machine MachineArch

		Version          uint32 // The elf version
		EntryPoint       uint64
		ProgHdrTblOffset uint64
		SectHdrTblOffset uint64

		Flags uint32

		// The size of the ELF header as on disk.
		// Note, it is not the byte size of this struct.
		HeaderSize uint16

		ProgHdrTblEntrySize  uint16
		ProgHdrTblEntryCount uint16
		SectHdrTblEntrySize  uint16
		SectHdrTblEntryCount uint16
		StrTblIndex          uint16
	}
}

func (header *header64) ELFIdent() *ELFIdent {
	return &header.ident
}

func (header *header64) Type() ELFType {
	return header.platformSpecific.Type
}

func (header *header64) Machine() MachineArch {
	return header.platformSpecific.Machine
}

func (header *header64) Version() uint32 {
	return header.platformSpecific.Version
}

func (header *header64) EntryPoint() uint64 {
	return header.platformSpecific.EntryPoint
}

func (header *header64) ProgHdrTblOffset() uint64 {
	return header.platformSpecific.ProgHdrTblOffset
}

func (header *header64) SectHdrTblOffset() uint64 {
	return header.platformSpecific.SectHdrTblOffset
}

func (header *header64) Flags() uint32 {
	return header.platformSpecific.Flags
}

func (header *header64) HeaderSize() uint16 {
	return header.platformSpecific.HeaderSize
}

func (header *header64) ProgHdrTblEntrySize() uint16 {
	return header.platformSpecific.ProgHdrTblEntrySize
}

func (header *header64) ProgHdrCount() uint16 {
	return header.platformSpecific.ProgHdrTblEntryCount
}

func (header *header64) SectHdrTblEntrySize() uint16 {
	return header.platformSpecific.SectHdrTblEntrySize
}

func (header *header64) SectHdrCount() uint16 {
	return header.platformSpecific.SectHdrTblEntryCount
}

func (header *header64) StrTblIndex() uint16 {
	return header.platformSpecific.StrTblIndex
}

func readHeader(file *os.File) (ELFHeader, error) {
	fileName := file.Name()
	_, err := file.Seek(0, 0)
	if err != nil {
		err = fmt.Errorf("Unable to seek while reading '%s'.\n%s", fileName, err.Error())
		return nil, err
	}

	var ident ELFIdent
	err = binary.Read(file, binary.LittleEndian, &ident)
	if err != nil {
		err = fmt.Errorf("Error reading ELFIdent from '%s'.\n%s", fileName, err.Error())
		return nil, err
	}

	if ident.Class == Class32 {
		header := new(header32)

		header.ident = ident
		err = binary.Read(file, endianMap[ident.Endianess], &header.platformSpecific)
		if err != nil {
			err = fmt.Errorf(
				"Error reading platform specific part of header from '%s'.\n%s",
				fileName,
				err.Error())
			return nil, err
		}

		return header, nil
	} else {
		header := new(header64)

		header.ident = ident
		err = binary.Read(file, endianMap[ident.Endianess], &header.platformSpecific)
		if err != nil {
			err = fmt.Errorf(
				"Error reading platform specific part of header from '%s'.\n%s",
				fileName,
				err.Error())
			return nil, err
		}

		return header, nil
	}
}
