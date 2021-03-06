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
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"time"
)

// Value of SectType represent the different types of sections in an ELF file.
type SectType uint32

// A SectHdr interface capturing a class independent section header
// entry in the section header table. The actual layout and sizes of the
// different fields of a section header are as in the table below. One
// should use the Class method to figure the class of section header at
// hand. The sizes of the values returned by the methods are of the largest
// size that can accomodate both elf32 and elf64 sizes.
//
//  ======================================================================
//  ||                        ||       ELF32       ||       ELF64       ||
//  ||                        ============================================
//  || Field Name             ||  Size  |  Offset  ||  Size  |  Offset  ||
//  ======================================================================
//  || name                   || uint32 |    0     || uint32 |    0     ||
//  ======================================================================
//  || type                   || uint32 |    4     || uint32 |    4     ||
//  ======================================================================
//  || flags                  || uint32 |    8     || uint64 |    8     ||
//  ======================================================================
//  || addr                   || uint32 |    12    || uint64 |    16    ||
//  ======================================================================
//  || offset                 || uint32 |    16    || uint64 |    24    ||
//  ======================================================================
//  || size                   || uint32 |    20    || uint64 |    32    ||
//  ======================================================================
//  || link                   || uint32 |    24    || uint32 |    40    ||
//  ======================================================================
//  || info                   || uint32 |    28    || uint32 |    44    ||
//  ======================================================================
//  || addralign              || uint32 |    32    || uint64 |    48    ||
//  ======================================================================
//  || entsize                || uint32 |    36    || uint64 |    56    ||
//  ======================================================================
type SectHdr interface {
	// Return the class of the ELF file to which this section header belongs.
	Class() ELFClass

	// Return the index into the section name string table where the name
	// of the section can be found. This is the byte index in the string
	// table section.
	NameIndex() uint32

	// Return the type of the section.
	Type() SectType

	// Return the sections flags.
	Flags() uint64

	// Return the section address.
	Address() uint64

	// Return the offset of the section in the ELF file.
	Offset() uint64

	// Return the byte size of the section.
	Size() uint64

	// Return the link data of the section.
	Link() uint32

	// Return the info data of the section.
	Info() uint32

	// Return the alignment of the section.
	Alignment() uint64

	// Return the size of elements in case of tabular sections.
	EntrySize() uint64
}

const (
	SectTypeUnused            SectType = SectType(0)
	SectTypeProgBits          SectType = SectType(1)
	SectTypeSymTab            SectType = SectType(2)
	SectTypeStrTab            SectType = SectType(3)
	SectTypeRelA              SectType = SectType(4)
	SectTypeHashTab           SectType = SectType(5)
	SectTypeDynamic           SectType = SectType(6)
	SectTypeNotes             SectType = SectType(7)
	SectTypeNoBits            SectType = SectType(8)
	SectTypeRel               SectType = SectType(9)
	SectTypeDynSym            SectType = SectType(11)
	SectTypeInitArray         SectType = SectType(14)
	SectTypeFinalizeArray     SectType = SectType(15)
	SectTypePreInitArray      SectType = SectType(16)
	SectTypeGroup             SectType = SectType(17)
	SectTypeExtSectIndeces    SectType = SectType(18)
	SectTypeNumDefinedTypes   SectType = SectType(19)
	SectTypeStartOSSpecific   SectType = SectType(0x60000000)
	SectTypeEndOSSpecific     SectType = SectType(0x6fffffff)
	SectTypeStartProcSpecific SectType = SectType(0x70000000)
	SectTypeEndProcSpecific   SectType = SectType(0x7fffffff)
	SectTypeStartAppSpecific  SectType = SectType(0x80000000)
	SectTypeEndAppSpecific    SectType = SectType(0x8fffffff)
)

const (
	SectIndexSectNameTblExt    uint16 = 0xFFFF
	SectIndexStartReserved     uint16 = 0xFF00
	SectIndexStartProcSpecific uint16 = 0xFF00
	SectIndexStartOSSpecific   uint16 = 0xFF20
	SectIndexEndOSSpecific     uint16 = 0xFF3F
	SectIndexAbsSym            uint16 = 0xFFF1
	SectIndexCommonSym         uint16 = 0xFFF2
	SectIndexEndProcSpecific   uint16 = 0xFF1F
	SectIndexEndReserved       uint16 = 0xFFFF
)

const (
	// Name of the section which is a string table containing section names.
	NameSectNameTbl = ".shstrtab"

	// Name of the section which contains the ELF symbol table.
	NameSymTab = ".symtab"

	// Name of the section which is a string table containing names of symbols
	// found in the '.symtab' section.
	NameSymNameTbl = ".strtab"

	// Name of the section which contains the ELF dynamic symbol table.
	NameDynSymTab = ".dynsym"

	// Name of the section which is a string table containing names of symbols
	// found in the '.dynsym' section.
	NameDynSymNameTbl = ".dynstr"
)

type sectHdr32 struct {
	diskData struct {
		NameIndex uint32
		Type      SectType
		Flags     uint32
		Addr      uint32
		Offset    uint32
		Size      uint32
		Link      uint32
		Info      uint32
		AddrAlign uint32
		EntSize   uint32
	}
}

func (sh *sectHdr32) Class() ELFClass {
	return Class32
}

func (sh *sectHdr32) NameIndex() uint32 {
	return sh.diskData.NameIndex
}

func (sh *sectHdr32) Type() SectType {
	return sh.diskData.Type
}

func (sh *sectHdr32) Flags() uint64 {
	return uint64(sh.diskData.Flags)
}

func (sh *sectHdr32) Address() uint64 {
	return uint64(sh.diskData.Addr)
}

func (sh *sectHdr32) Offset() uint64 {
	return uint64(sh.diskData.Offset)
}

func (sh *sectHdr32) Size() uint64 {
	return uint64(sh.diskData.Size)
}

func (sh *sectHdr32) Link() uint32 {
	return sh.diskData.Link
}

func (sh *sectHdr32) Info() uint32 {
	return sh.diskData.Info
}

func (sh *sectHdr32) Alignment() uint64 {
	return uint64(sh.diskData.AddrAlign)
}

func (sh *sectHdr32) EntrySize() uint64 {
	return uint64(sh.diskData.EntSize)
}

type sectHdr64 struct {
	diskData struct {
		NameIndex uint32
		Type      SectType
		Flags     uint64
		Addr      uint64
		Offset    uint64
		Size      uint64
		Link      uint32
		Info      uint32
		AddrAlign uint64
		EntSize   uint64
	}
}

func (sh *sectHdr64) Class() ELFClass {
	return Class64
}

func (sh *sectHdr64) NameIndex() uint32 {
	return sh.diskData.NameIndex
}

func (sh *sectHdr64) Type() SectType {
	return sh.diskData.Type
}

func (sh *sectHdr64) Flags() uint64 {
	return sh.diskData.Flags
}

func (sh *sectHdr64) Address() uint64 {
	return sh.diskData.Addr
}

func (sh *sectHdr64) Offset() uint64 {
	return sh.diskData.Offset
}

func (sh *sectHdr64) Size() uint64 {
	return sh.diskData.Size
}

func (sh *sectHdr64) Link() uint32 {
	return sh.diskData.Link
}

func (sh *sectHdr64) Info() uint32 {
	return sh.diskData.Info
}

func (sh *sectHdr64) Alignment() uint64 {
	return sh.diskData.AddrAlign
}

func (sh *sectHdr64) EntrySize() uint64 {
	return sh.diskData.EntSize
}

func readSectHdrTbl(f *os.File, header ELFHeader) ([]SectHdr, uint32, error) {
	elfIdent := header.ELFIdent()
	class := elfIdent.Class
	e := elfIdent.Endianess
	offset := header.SectHdrTblOffset()
	_, err := f.Seek(int64(offset), 0)
	if err != nil {
		return nil, 0, err
	}

	var sectCount uint64
	var strTblIndex uint32
	n := header.SectHdrCount()
	if n == 0 {
		if class == Class32 {
			var sectHdr32 sectHdr32
			err = binary.Read(f, endianMap[e], &sectHdr32.diskData)
			sectCount = uint64(sectHdr32.diskData.Size)
			strTblIndex = sectHdr32.diskData.Link
		} else {
			var sectHdr64 sectHdr64
			err = binary.Read(f, endianMap[e], &sectHdr64.diskData)
			sectCount = sectHdr64.diskData.Size
			strTblIndex = sectHdr64.diskData.Link
		}
		if err != nil {
			return nil, 0, errors.New("Error reading section header 0.\n" + err.Error())
		}

		// Reset the file position.
		_, err = f.Seek(int64(offset), 0)
		if err != nil {
			return nil, 0, err
		}
	} else {
		sectCount = uint64(n)
		strTblIndex = uint32(header.StrTblIndex())
	}

	sectHdrTbl := make([]SectHdr, sectCount)
	for i := uint64(0); i < sectCount; i++ {
		if class == Class32 {
			sectHdr32 := new(sectHdr32)
			sectHdrTbl[i] = sectHdr32
			err = binary.Read(f, endianMap[e], &sectHdr32.diskData)
		} else {
			sectHdr64 := new(sectHdr64)
			sectHdrTbl[i] = sectHdr64
			err = binary.Read(f, endianMap[e], &sectHdr64.diskData)
		}
		if err != nil {
			return nil, 0, errors.New("Error reading section header.\n" + err.Error())
		}
	}

	return sectHdrTbl, strTblIndex, nil
}

// Section represents a section of an ELF file.
type Section struct {
	name     string
	header   SectHdr
	data     []byte
	fileName string
	modTime  time.Time
}

// Returns the header of the section.
func (section *Section) SectHdr() SectHdr {
	return section.header
}

// Returns the name of the section.
func (section *Section) Name() string {
	return section.name
}

// Returns a reader whose view is the section data.
func (section *Section) NewReader() (*bytes.Reader, error) {
	data, err := section.Data()
	if err != nil {
		err = fmt.Errorf(
			"Unable to read section data before creating section reader.\n%s",
			err.Error())
		return nil, err
	}

	return bytes.NewReader(data), nil
}

// Returns the section data.
// The section data is cached in memory. Only the first call to Data reads the
// section data from memory. All subsequent calls return the cached data.
func (section *Section) Data() ([]byte, error) {
	fileInfo, err := os.Stat(section.fileName)
	if err != nil {
		return nil, fmt.Errorf("Unable to stat '%s'.\n%s", section.fileName, err.Error())
	}

	if section.modTime.Unix() < fileInfo.ModTime().Unix() {
		err = fmt.Errorf(
			"File '%s' modified after loading. Cannot read data for section '%s'",
			section.fileName, section.name)
		return nil, err
	}

	if section.data != nil {
		return section.data, nil
	}

	file, err := os.Open(section.fileName)
	if err != nil {
		err = fmt.Errorf(
			"Unable to open '%s' to read data for section '%s'.\n%s",
			section.fileName, section.name, err.Error())
		return nil, err
	}
	defer file.Close()

	_, err = file.Seek(int64(section.header.Offset()), 0)
	if err != nil {
		err = fmt.Errorf(
			"Unable to seek to section '%s' in '%s' to read raw data.\n%s",
			section.name, file.Name(), err.Error())
		return nil, err
	}

	var data []byte
	for i := uint64(0); i < section.header.Size(); i++ {
		var oneByte byte
		err = binary.Read(file, binary.LittleEndian, &oneByte)
		if err != nil {
			err = fmt.Errorf(
				"Error reading raw data from '%s'.\n%s", file.Name(), err.Error())
			return nil, err
		}
		data = append(data, oneByte)
	}

	section.data = data
	return section.data, nil
}

func newSection(name string, sectHdr SectHdr, fileName string) (*Section, error) {
	section := new(Section)

	section.name = name
	section.header = sectHdr
	section.data = nil
	section.fileName = fileName

	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return nil, fmt.Errorf("Unable to stat '%s'.\n%s", section.fileName, err.Error())
	}

	section.modTime = fileInfo.ModTime()

	return section, nil
}

// StrTbl represents a string table in an ELF file. It is a mapping from byte
// indeces to strings.
type StrTbl map[uint32]string

func BuildStrTbl(data []byte) (StrTbl, error) {
	// Read the first NULL string
	if data[0] != 0 {
		err := fmt.Errorf("First byte in the string table is not NULL.")
		return nil, err
	}

	stringMap := make(map[uint32]string)
	stringMap[uint32(0)] = string(byte(0))

	var str []uint8
	for index, char := range data {
		if index == 0 {
			continue
		}

		str = append(str, char)
		if char == 0 {
			stringMap[uint32(index-(len(str)-1))] = string(str[0 : len(str)-1])
			if len(str) == 1 {
				break
			}
			str = nil
		}
	}

	return StrTbl(stringMap), nil
}

// Section map is a map from names to the list of sections with the same name.
// Since more than one section have the same name, each name maps to a list
// slice of sections.
type SectMap map[string][]*Section

func readSectMap(f *os.File, sectHdrTbl []SectHdr, sectNameTblIndex uint32) (SectMap, error) {
	sectMap := make(SectMap, len(sectHdrTbl))

	strTblSect, err := newSection("dummy-name", sectHdrTbl[sectNameTblIndex], f.Name())
	strTblData, err := strTblSect.Data()
	if err != nil {
		err = fmt.Errorf(
			"Error reading string table data from '%s'.\n%s", f.Name(), err.Error())
		return nil, err
	}
	strTbl, err := BuildStrTbl(strTblData)
	if err != nil {
		err = fmt.Errorf(
			"Unable to build string table from string table data.\n%s",
			err.Error())
		return nil, err
	}

	for _, sectHdr := range sectHdrTbl {
		sectName := strTbl[sectHdr.NameIndex()]
		_, exists := sectMap[sectName]
		if !exists {
			sectMap[sectName] = make([]*Section, 0)
		}
		section, err := newSection(sectName, sectHdr, f.Name())
		if err != nil {
			return nil, err
		}

		sectMap[sectName] = append(sectMap[sectName], section)
	}

	return sectMap, nil
}

// Symbol represents an entry for a symbol in a symbol table of an ELF file.
type Symbol interface {
	// Returns the byte index into string table where the name of this
	// symbol can be found.
	NameIndex() uint32

	// Returns the address (or value) of this symbol.
	Addr() uint64

	// Returns the size of the symbol.
	Size() uint64

	// Returns the symbol info.
	Info() uint8

	// Returns the symbol visibility.
	Visibility() uint8

	// Returns the index of the section in which this symbol can be found.
	SectIndex() uint16
}

type symbol32 struct {
	diskData struct {
		NameIndex  uint32
		Addr       uint32
		Size       uint32
		Info       uint8
		Visibility uint8
		SectIndex  uint16
	}
}

func (symbol *symbol32) NameIndex() uint32 {
	return symbol.diskData.NameIndex
}

func (symbol *symbol32) Addr() uint64 {
	return uint64(symbol.diskData.Addr)
}

func (symbol *symbol32) Size() uint64 {
	return uint64(symbol.diskData.Size)
}

func (symbol *symbol32) Info() uint8 {
	return symbol.diskData.Info
}

func (symbol *symbol32) Visibility() uint8 {
	return symbol.diskData.Visibility
}

func (symbol *symbol32) SectIndex() uint16 {
	return symbol.diskData.SectIndex
}

type symbol64 struct {
	diskData struct {
		NameIndex  uint32
		Info       uint8
		Visibility uint8
		SectIndex  uint16
		Addr       uint64
		Size       uint64
	}
}

func (symbol *symbol64) NameIndex() uint32 {
	return symbol.diskData.NameIndex
}

func (symbol *symbol64) Addr() uint64 {
	return symbol.diskData.Addr
}

func (symbol *symbol64) Size() uint64 {
	return symbol.diskData.Size
}

func (symbol *symbol64) Info() uint8 {
	return symbol.diskData.Info
}

func (symbol *symbol64) Visibility() uint8 {
	return symbol.diskData.Visibility
}

func (symbol *symbol64) SectIndex() uint16 {
	return symbol.diskData.SectIndex
}

// SymTab represents a symbol table in an ELF file, It is a mapping
// from the name index (the byte index into the string table containing
// symbol names) to a slice of symbols with the same name.
type SymTab map[uint32][]Symbol

func BuildSymTab(data []byte, sectHdr SectHdr, endianess ELFEndianess) (SymTab, error) {
	reader := bytes.NewReader(data)
	symTab := make(SymTab)
	var symbol Symbol
	var i uint64 = 0
	for ; i < sectHdr.Size(); i += sectHdr.EntrySize() {
		var err error
		if sectHdr.Class() == Class32 {
			sym32 := new(symbol32)
			err = binary.Read(reader, endianMap[endianess], &sym32.diskData)
			symbol = sym32
		} else {
			sym64 := new(symbol64)
			err = binary.Read(reader, endianMap[endianess], &sym64.diskData)
			symbol = sym64
		}

		if err != nil {
			return nil, fmt.Errorf("Error reading symtab from '%s'.\n%s", err.Error())
		} else {
			nameIndex := symbol.NameIndex()
			_, exists := symTab[nameIndex]
			if !exists {
				symTab[nameIndex] = make([]Symbol, 0)
			}
			symTab[symbol.NameIndex()] = append(symTab[nameIndex], symbol)
		}
	}

	return symTab, nil
}
