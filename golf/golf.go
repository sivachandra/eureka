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

// Package golf provides API to read ELF files from first principles.
package golf

import (
	"encoding/binary"
	"fmt"
	"os"
)

// ELF encapsulates the data of an ELF file. It is the enrty point to reading
// symbols, strings etc. from an ELF file.
type ELF struct {
	header           ELFHeader
	progHdrTbl       []SegHdr
	sectHdrTbl       []SectHdr
	sectMap          SectMap
	sectNameTblIndex uint32
}

// Returns the ELF header.
func (elf *ELF) Header() ELFHeader {
	return elf.header
}

// Returns the endianess of the data in the ELF file.
func (elf *ELF) Endianess() binary.ByteOrder {
	return endianMap[elf.Header().ELFIdent().Endianess]
}

// Returns the address size of the architecture in bytes.
func (elf *ELF) AddressSize() uint8 {
	switch elf.Header().Machine() {
	case MachineX86:
		return 4
	case MachineARM:
		return 4
	case MachineX86_64:
		return 8
	case MachineAArch64:
		return 8
	default:
		return 0
	}
}

// Returns the program header table.
func (elf *ELF) ProgHdrTbl() []SegHdr {
	return elf.progHdrTbl
}

// Returns the section header table.
func (elf *ELF) SectHdrTbl() []SectHdr {
	return elf.sectHdrTbl
}

// Returns the index of the string table holding section names.
// Note that this is the true index of the table holding section names, and not
// the one found in the ELF header. [The string table index in the header could
// be set to SectNameTblExtIndex of 0xFFFF in case of extended numbering.]
func (elf *ELF) SectNameTblIndex() uint32 {
	return elf.sectNameTblIndex
}

// Returns a mapping from section names to a slice of sections. Since multiple
// sections can have the same name, each name maps to a slice of sections having
// that same name.
func (elf *ELF) SectMap() SectMap {
	return elf.sectMap
}

// Reads in an ELF file whose path is given by the string value fileName.
// If successful, it returns a pointer to the ELF object and nil error.
// If reading the file fails, then nil is returned along with the
// appropriate error message.
func Read(fileName string) (elf *ELF, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("Unable to open file '%s'.\n%s", fileName, err.Error())
	}
	defer file.Close()

	elf = new(ELF)
	elf.header, err = readHeader(file)
	if err != nil {
		return nil, fmt.Errorf("Error reading header from '%s'.\n%s", fileName, err.Error())
	}

	sectHdrTbl, sectNameTblIndex, err := readSectHdrTbl(file, elf.header)
	if err != nil {
		err := fmt.Errorf(
			"Error reading section header table from '%s'.\n%s", fileName, err.Error())
		return nil, err
	}
	elf.sectHdrTbl = sectHdrTbl
	elf.sectNameTblIndex = sectNameTblIndex

	elf.sectMap, err = readSectMap(file, sectHdrTbl, sectNameTblIndex)
	if err != nil {
		return nil, err
	}

	elf.progHdrTbl, err = readSegHdrTbl(file, elf.header)
	if err != nil {
		return nil, err
	}

	return elf, nil
}
