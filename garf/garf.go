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

// Package garf provides API to read DWARF debug info from ELF files.
package garf

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

import (
	"eureka/golf"
	"eureka/guts/leb128"
	"eureka/guts/ruts"
)

type DwFormat uint8

const (
	DwFormat32 = DwFormat(0)
	DwFormat64 = DwFormat(1)
)

type AttrForm struct {
	Name DwAt
	Form DwForm
}

type AbbrevEntry struct {
	AbbrevCode  uint64
	Tag         DwTag
	HasChildren bool
	AttrForms   []AttrForm
}

type AbbrevTable map[uint64]AbbrevEntry

type Attribute struct {
	Name  DwAt
	Value interface{}
}

type DIE struct {
	// The debug info tag of this DIE.
	Tag DwTag

	// A map of attributes of this DIE.
	Attributes map[DwAt]Attribute

	// The parent DIE of this DIE.
	Parent *DIE

	// The children of this DIE.
	Children []*DIE

	// The unit to which this DIE belongs to.
	Unit *DwUnit

	// Offset of the first byte of the contribution of this DIE
	// in the .debug_info section.
	startOffset uint64

	// Offset of the first byte after the end of the contribution of this DIE.
	endOffset uint64
}

// DwOperation is an operation in a DWARF expression with an opcode and its operands.
type DwOperation struct {
	// The opcode of the operation
	Op DwOp

	// The operands of the operation. They are one of the integral types:
	//   int8, uint8, int16, uint16, int32, uint32, int64, uint64
	// The DWARF standard prescribes the type of operands an opcode takes.
	// For convenience, LEB128 and ULEB128 numbers are stored as int64 and
	// uint64 numbers respectively. Also, operands which denote addresses on
	// the target architecture are always stored as uint64 values.
	Operands []interface{}
}

type DwExpr []DwOperation

type LocListEntryType uint8

const (
	LocListEntryTypeNormal            = LocListEntryType(1)
	LocListEntryTypeDefault           = LocListEntryType(2)
	LocListEntryTypeBaseAddrSelection = LocListEntryType(3)
	LocListEntryTypeEndOfList         = LocListEntryType(4)
)

type LocListEntry interface {
	LocListEntryType() LocListEntryType
}

type NormalLocListEntry struct {
	Begin uint64
	End   uint64
	Loc   DwExpr
}

func (e NormalLocListEntry) LocListEntryType() LocListEntryType {
	return LocListEntryTypeNormal
}

type DefaultLocListEntry DwExpr

func (e DefaultLocListEntry) LocListEntryType() LocListEntryType {
	return LocListEntryTypeDefault
}

type BaseAddrSelectionLocListEntry uint64

func (e BaseAddrSelectionLocListEntry) LocListEntryType() LocListEntryType {
	return LocListEntryTypeBaseAddrSelection
}

type EndOfListLocListEntry struct {
}

func (e EndOfListLocListEntry) LocListEntryType() LocListEntryType {
	return LocListEntryTypeEndOfList
}

type LocList []LocListEntry

type LnInfoTimestamp interface {
}

type LnFileEntry struct {
	Path      string
	DirIndex  uint64
	Size      uint64
	Timestamp uint64
	MD5       [16]byte
}

type DwLnOpcodeType uint8

const (
	DwLnOpcodeSpecial = DwLnOpcodeType(0x01)
	DwLnOpcodeStd     = DwLnOpcodeType(0x02)
	DwLnOpcodeExt     = DwLnOpcodeType(0x03)
)

type LnInstr struct {
	// The instruction opcode. Its type (special, standard or extension) is
	// determined by the OpcodeType field.
	Opcode DwLnOpcode

	// OpcodeType takes one of DwLnOpcodeSpecial, DwLnOpcodeStd or
	// DwLnOpcodeExt value.
	OpcodeType DwLnOpcodeType

	// Note 1: DW_LNS_fixed_advance_pc instruction takes a single usigned
	// 16-bit operand. That operand will also be stored as an unsigned LEB128
	// number.
	//
	// Note 2: The opcode DW_LNE_define_file which is deprecated in DWARF 5
	// takes a string operand. There are no known producers which emit this
	// opcode. Hence, we do not support it here.
	Operands []leb128.LEB128
}

type LnInfo struct {
	Size                uint64
	Version             uint16
	AddressSize         uint8
	SegmentSelectorSize uint8

	Directories []string
	Files       []LnFileEntry

	minInstrLength  uint8
	maxOprPerInstr  uint8
	defaultIsStmt   uint8
	lineBase        int8
	lineRange       uint8
	opcodeBase      uint8
	operandCountTbl []uint8

	Program []LnInstr
}

type DwUnit struct {
	// The parent DwData from which this unit was read from.
	Parent *DwData

	// Type of the unit.
	Type DwUnitType

	// Format of the unit.
	Format DwFormat

	// The size of the DW_AT_addr attributes in this unit.
	AddressSize byte

	// The DWARF version of this unit.
	Version uint16

	// Size of the unit in the .debug_info section.
	// It is not the same as the initial length feild in the unit's
	// header.
	size uint64

	// The offset into the .debug_abbrev section where the info for this
	// unit begins.
	debugAbbrevOffset uint64

	// Offset of this units header in the .debug_info section.
	headerOffset uint64

	// Offset into the .debug_info section at which the data for
	// DIE tree of this unit begins, after this unit's header.
	dataOffset uint64

	// The abbreviation table for this unit. Will be nil until a call to the
	// DIETree method.
	abbrevTable AbbrevTable

	// The complete DIE tree of this unit. Will be nil until a call to the
	// DIETree method.
	dieTree *DIE

	// The line number program for this unit. Will be nil until a call to the
	// LnInfo method.
	lnInfo *LnInfo
}

func (u *DwUnit) DIETree() (*DIE, error) {
	if u.dieTree != nil {
		return u.dieTree, nil
	}

	var err error
	u.dieTree, err = u.Parent.readDIETree(u, u.dataOffset)
	return u.dieTree, err
}

func (u *DwUnit) LineNumberInfo() (*LnInfo, error) {
	if u.lnInfo != nil {
		return u.lnInfo, nil
	}

	var err error
	u.lnInfo, err = u.Parent.readLineNumberInfo(u)
	return u.lnInfo, err
}

// DebugStrTbl encapsulates the data in the .debug_str section.
type DebugStrTbl struct {
	data []byte
}

// Reads a string from the .debug_str data at the specified offset.
//
// Die attributes can refer to full or partial strings. Hence, we do not
// prefetch the full strings. Each string is read out on demand from the
// specified offset.
func (t *DebugStrTbl) ReadStr(offset uint64) (string, error) {
	if offset >= uint64(len(t.data)) {
		return "", fmt.Errorf("Invalid .debug_str offset.")
	}

	r := bytes.NewReader(t.data)
	_, err := r.Seek(int64(offset), 0)
	if err != nil {
		return "", fmt.Errorf("Unable to seek to .debug_str offset.\n%s", err.Error())
	}

	return ruts.ReadCString(r)
}

type DwData struct {
	fileName    string
	elf         *golf.ELF
	debugStrTbl *DebugStrTbl
	compUnits   []*DwUnit
	typeUnits   []*DwUnit

	// Mapping from offset into the .debug_info section to the DIE at that
	// offset.
	dieMap map[uint64]*DIE
}

func LoadDwData(fileName string) (*DwData, error) {
	dwData := new(DwData)
	var err error

	dwData.fileName = fileName
	dwData.elf, err = golf.Read(fileName)
	if err != nil {
		err = fmt.Errorf("Error loading ELF info from '%s'.\n%s", fileName, err.Error())
		return nil, err
	}

	dwData.dieMap = make(map[uint64]*DIE)

	return dwData, nil
}

func (d *DwData) ELFData() *golf.ELF {
	return d.elf
}

func (d *DwData) FileName() string {
	return d.fileName
}

func (d *DwData) AbbrevTable(offset uint64) (AbbrevTable, error) {
	sectMap := d.elf.SectMap()
	sections, exists := sectMap[".debug_abbrev"]
	if !exists {
		return nil, fmt.Errorf(".debug_abbrev section is not present.", nil)
	}

	if len(sections) > 1 {
		return nil, fmt.Errorf("More than one .debug_abbrev sections.", nil)
	}

	reader, err := sections[0].NewReader()
	if err != nil {
		return nil, fmt.Errorf("Error fetching .debug_abbrev reader.", err)
	}

	_, err = reader.Seek(int64(offset), 0)
	if err != nil {
		return nil, fmt.Errorf("Error seeking to .debug_abbrev offset.")
	}

	table := make(AbbrevTable)
	for true {
		abbrevCode, err := leb128.ReadUnsigned(reader)
		if err != nil {
			return nil, fmt.Errorf("Error reading abbreviation code.", nil)
		}
		if abbrevCode == NullAbbrevEntry {
			break
		}

		tag, err := leb128.ReadUnsigned(reader)
		if err != nil {
			msg := fmt.Sprintf("Error reading tag for abbrev code %d.", abbrevCode)
			return nil, fmt.Errorf(msg, err)
		}

		hasChildren, err := reader.ReadByte()
		if err != nil {
			msg := fmt.Sprintf(
				"Error reading child determination entry for abbrev code %d.",
				abbrevCode)
			return nil, fmt.Errorf(msg, err)
		}

		var entry AbbrevEntry
		entry.AbbrevCode = abbrevCode
		entry.Tag = DwTag(tag)
		entry.HasChildren = (hasChildren == DW_CHILDREN_yes)
		entry.AttrForms = make([]AttrForm, 0)

		for true {
			attr, err := leb128.ReadUnsigned(reader)
			if err != nil {
				msg := fmt.Sprintf(
					"Error reading an attr name of entry with abbrev code %d.",
					abbrevCode)
				return nil, fmt.Errorf(msg, err)
			}

			form, err := leb128.ReadUnsigned(reader)
			if err != nil {
				msg := fmt.Sprintf(
					"Error reading an attr form of entry with abbrev code %d.",
					abbrevCode)
				return nil, fmt.Errorf(msg, err)
			}

			if form == 0 && attr == 0 {
				break
			}

			var pair AttrForm
			pair.Name = DwAt(attr)
			pair.Form = DwForm(form)
			entry.AttrForms = append(entry.AttrForms, pair)
		}

		table[entry.AbbrevCode] = entry
	}

	return table, nil
}

func (d *DwData) CompUnits() ([]*DwUnit, error) {
	if d.compUnits != nil {
		return d.compUnits, nil
	}

	sectMap := d.elf.SectMap()
	sections, exists := sectMap[".debug_info"]
	if !exists {
		return nil, fmt.Errorf(".debug_info section is not present.", nil)
	}

	if len(sections) > 1 {
		return nil, fmt.Errorf("More than one .debug_info sections.", nil)
	}

	reader, err := sections[0].NewReader()
	if err != nil {
		return nil, fmt.Errorf("Error fetching .debug_info section reader.", err)
	}

	d.compUnits = make([]*DwUnit, 0)
	en := d.elf.Endianess()
	for true {
		if reader.Len() == 0 {
			break
		}

		var length uint64
		var format DwFormat
		var size32 uint32

		headerOffset := uint64(reader.Size() - int64(reader.Len()))

		err := binary.Read(reader, en, &size32)
		if err != nil {
			err = fmt.Errorf(
				"Error reading first 32 bits of length of a unit in .debug_info.",
				err)
			return nil, err
		}

		if size32 == 0xffffffff {
			format = DwFormat64
			var size64 uint64

			err := binary.Read(reader, en, &size64)
			if err != nil {
				err = fmt.Errorf(
					"Error reading 64-bit length of a unit in .debug_info.",
					err)
				return nil, err
			}

			length = size64
		} else {
			format = DwFormat32
			length = uint64(size32)
		}

		var version uint16
		err = binary.Read(reader, en, &version)
		if err != nil {
			err = fmt.Errorf("Error reading version of a unit in .debug_info.", err)
			return nil, err
		}

		unitType := DW_UT_compile
		if version >= 5 {
			err = binary.Read(reader, en, &unitType)
			if err != nil {
				err = fmt.Errorf(
					"Error reading unit type of a unit in .debug_info.", err)
				return nil, err
			}
		}

		var debugAbbrevOffset uint64
		if format == DwFormat32 {
			var offset uint32
			err = binary.Read(reader, en, &offset)
			if err != nil {
				err = fmt.Errorf(
					"Error reading 32-bit debug abbrev offset of a unit.", err)
				return nil, err
			}

			debugAbbrevOffset = uint64(offset)
		} else {
			err = binary.Read(reader, en, &debugAbbrevOffset)
			if err != nil {
				err = fmt.Errorf(
					"Error reading 64-bit debug abbrev offset of a unit.", err)
				return nil, err
			}
		}

		var addrSize byte
		err = binary.Read(reader, en, &addrSize)
		if err != nil {
			err = fmt.Errorf(
				"Error reading address size from a unit header in .debug_info.",
				err)
			return nil, err
		}

		if unitType == DW_UT_type {
		} else {
			cu := new(DwUnit)

			cu.Parent = d
			cu.Type = unitType

			if format == DwFormat64 {
				cu.size = length + 12
			} else {
				cu.size = length + 4
			}

			cu.Format = format
			cu.Version = version
			cu.headerOffset = headerOffset
			cu.debugAbbrevOffset = debugAbbrevOffset
			cu.AddressSize = addrSize
			cu.dataOffset = uint64(reader.Size() - int64(reader.Len()))
			cu.abbrevTable = nil
			d.compUnits = append(d.compUnits, cu)
			reader.Seek(int64(cu.size+headerOffset), 0)
		}
	}

	return d.compUnits, nil
}

func (d *DwData) DebugStr() (*DebugStrTbl, error) {
	if d.debugStrTbl != nil {
		return d.debugStrTbl, nil
	}

	sectMap := d.elf.SectMap()
	debugStrSections, exists := sectMap[".debug_str"]
	if !exists {
		return nil, fmt.Errorf(".debug_str section is not present.", nil)
	}

	if len(debugStrSections) > 1 {
		return nil, fmt.Errorf("More than one .debug_str sections.", nil)
	}

	debugStrData, err := debugStrSections[0].Data()
	if err != nil {
		return nil, fmt.Errorf("Error fetching .debug_str data.", err)
	}

	d.debugStrTbl = new(DebugStrTbl)
	d.debugStrTbl.data = debugStrData
	return d.debugStrTbl, nil
}

func (d *DwData) readDIETree(u *DwUnit, offset uint64) (*DIE, error) {
	sectMap := d.elf.SectMap()
	sections, exists := sectMap[".debug_info"]
	if !exists {
		return nil, fmt.Errorf(".debug_info section is not present.", nil)
	}

	if len(sections) > 1 {
		return nil, fmt.Errorf("More than one .debug_info sections.", nil)
	}

	reader, err := sections[0].NewReader()
	if err != nil {
		return nil, fmt.Errorf("Error fetching .debug_info section reader.", err)
	}

	_, err = reader.Seek(int64(offset), 0)
	if err != nil {
		err = fmt.Errorf(
			"Error seeking to the DIE offset to read the DIE tree.\n%s", err.Error())
		return nil, err
	}

	return d.readDIETreeHelper(u, reader, d.elf.Endianess(), nil)
}

func (d *DwData) readDIETreeHelper(
	u *DwUnit, r *bytes.Reader, en binary.ByteOrder, parent *DIE) (*DIE, error) {
	// This is the DIE's offset in .debug_info section.
	offset := uint64(r.Size() - int64(r.Len()))

	die, exists := d.dieMap[offset]
	if exists {
		die.Parent = parent
		r.Seek(int64(die.endOffset), 0)

		return die, nil
	}

	if u.abbrevTable == nil {
		var err error
		u.abbrevTable, err = d.AbbrevTable(u.debugAbbrevOffset)
		if err != nil {
			err = fmt.Errorf(
				"Error getting abbrev table while reading a DIE tree.\n%s",
				err.Error())
			return nil, err
		}
	}

	abbrevCode, err := leb128.ReadUnsigned(r)
	if err != nil {
		return nil, fmt.Errorf("Error reading abbrev code of a DIE.", err)
	}

	// Return if its a NULL entry
	if abbrevCode == 0 {
		return nil, nil
	}

	abbrevEntry, exists := u.abbrevTable[abbrevCode]
	if !exists {
		return nil, fmt.Errorf("Invalid abbrev code for a DIE.", nil)
	}

	die = new(DIE)
	die.Tag = abbrevEntry.Tag
	die.Parent = parent
	die.Unit = u
	die.startOffset = offset

	// We register the DIE in the die map even before the the complete die is
	// read out as we want to avoid infinite recursion due to DIEs referring to
	// each other cyclically. For example, a DIE for a type can have a typedef
	// DIE as its sibling, which in turn refers to the DIE itself.
	//
	// The registered DIE should be deleted from the DIE map if an error occurs
	// reading it.
	d.dieMap[offset] = die

	attributes := make(map[DwAt]Attribute)
	for _, attrForm := range abbrevEntry.AttrForms {
		attr, err := d.readAttr(u, r, attrForm.Name, attrForm.Form, en)
		if err != nil {
			delete(d.dieMap, offset)
			msg := fmt.Sprintf(
				"Error reading value of attribute %s of tag %s at offset %x.\n%s",
				DwAtStr[attrForm.Name], DwTagStr[abbrevEntry.Tag],
				offset, err.Error())
			err = fmt.Errorf(msg)
			return nil, err
		}
		attributes[attr.Name] = attr
	}
	die.Attributes = attributes

	for abbrevEntry.HasChildren {
		childDie, err := d.readDIETreeHelper(u, r, en, die)
		if err != nil {
			delete(d.dieMap, offset)
			err = fmt.Errorf(
				"Error reading child DIE tree of tag %x at offset %x.\n%s",
				DwTagStr[abbrevEntry.Tag], offset, err.Error())
			return nil, err
		}

		if childDie == nil {
			break
		}

		die.Children = append(die.Children, childDie)
	}
	die.endOffset = uint64(r.Size() - int64(r.Len()))

	return die, nil
}
