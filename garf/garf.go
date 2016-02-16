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
	"bytes"
	"encoding/binary"
	"fmt"
)

import (
	"eureka/golf"
	"eureka/utils/leb128"
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
	Tag        DwTag
	Attributes map[DwAt]Attribute
	Parent     *DIE
	Children   []*DIE

	// Offset of the first byte of the contribution of this DIE.
	debugInfoOffsetStart uint64

	// Offset of the first byte after the end of the contribution of this DIE.
	debugInfoOffsetEnd uint64
}

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

	// Size of the unit in the .debug_info section.
	// It is not the same as the initial length feild in the unit's
	// header.
	Size uint64

	// The DWARF version of this unit.
	Version uint16

	// The offset into the .debug_abbrev section where the info for this
	// unit begins.
	DebugAbbrevOffset uint64

	// The size of the DW_AT_addr attributes in this unit.
	AddressSize byte

	// Offset into the .debug_info section at which the data for
	// DIE tree of this unit begins. It is NOT the offset at which the
	// header for this unit begins.
	DebugInfoOffset uint64

	// The complete DIE tree of this unit. Will be nil until a call to the
	// DIETree method.
	dieTree *DIE

	// The line number program for this unit. Will be nil until a call to the
	// LnInfo method.
	lnInfo *LnInfo
}

func (u DwUnit) DIETree() (*DIE, error) {
	if u.dieTree != nil {
		return u.dieTree, nil
	}

	var err error
	u.dieTree, err = u.Parent.readDIETree(&u, u.DebugInfoOffset)
	return u.dieTree, err
}

func (u DwUnit) LineNumberInfo() (*LnInfo, error) {
	if u.lnInfo != nil {
		return u.lnInfo, nil
	}

	var err error
	u.lnInfo, err = u.Parent.readLineNumberInfo(&u)
	return u.lnInfo, err
}

type DebugStrMap map[uint64]string

type DwData struct {
	fileName       string
	elf            *golf.ELF
	abbrevTableMap map[uint64]AbbrevTable
	debugStrMap    DebugStrMap
	compUnits      []DwUnit
	typeUnits      []DwUnit

	// Mapping from offset into .debug_info section to the DIE at that
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
	dwData.abbrevTableMap = make(map[uint64]AbbrevTable)

	return dwData, nil
}

func (d *DwData) ELFData() *golf.ELF {
	return d.elf
}

func (d *DwData) FileName() string {
	return d.fileName
}

func (d *DwData) AbbrevTable(offset uint64) (AbbrevTable, error) {
	abbrevTable, exists := d.abbrevTableMap[offset]
	if exists {
		return abbrevTable, nil
	}

	sectMap := d.elf.SectMap()
	debugAbbrevSections, exists := sectMap[".debug_abbrev"]
	if !exists {
		return nil, fmt.Errorf(".debug_abbrev section is not present.", nil)
	}

	if len(debugAbbrevSections) > 1 {
		return nil, fmt.Errorf("More than one .debug_abbrev sections.", nil)
	}

	reader, err := debugAbbrevSections[0].NewReader()
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

	d.abbrevTableMap[offset] = table
	return table, nil
}

func (d *DwData) CompUnits() ([]DwUnit, error) {
	if d.compUnits != nil {
		return d.compUnits, nil
	}

	sectMap := d.elf.SectMap()
	debugInfoSections, exists := sectMap[".debug_info"]
	if !exists {
		return nil, fmt.Errorf(".debug_info section is not present.", nil)
	}

	if len(debugInfoSections) > 1 {
		return nil, fmt.Errorf("More than one .debug_info sections.", nil)
	}

	reader, err := debugInfoSections[0].NewReader()
	if err != nil {
		return nil, fmt.Errorf("Error fetching .debug_info section reader.", err)
	}

	d.compUnits = make([]DwUnit, 0)
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
			var cu DwUnit

			cu.Parent = d
			cu.Type = unitType

			if format == DwFormat64 {
				cu.Size = length + 12
			} else {
				cu.Size = length + 4
			}

			cu.Format = format
			cu.Version = version
			cu.DebugAbbrevOffset = debugAbbrevOffset
			cu.AddressSize = addrSize
			cu.DebugInfoOffset = uint64(reader.Size() - int64(reader.Len()))
			d.compUnits = append(d.compUnits, cu)
			reader.Seek(int64(cu.Size+headerOffset), 0)
		}
	}

	return d.compUnits, nil
}

func (d *DwData) DebugStr() (DebugStrMap, error) {
	if d.debugStrMap != nil {
		return d.debugStrMap, nil
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

	var str []uint8
	offset := uint64(0)
	debugStrMap := make(DebugStrMap)
	for offset < uint64(len(debugStrData)) {
		c := debugStrData[offset]
		str = append(str, c)
		if c == 0 {
			l := len(str)
			debugStrMap[offset+1-uint64(l)] = string(str[0 : l-1])
			str = nil
		}
		offset++
	}

	d.debugStrMap = debugStrMap
	return d.debugStrMap, nil
}

func (d *DwData) readDIETree(u *DwUnit, offset uint64) (*DIE, error) {
	sectMap := d.elf.SectMap()
	debugInfoSections, exists := sectMap[".debug_info"]
	if !exists {
		return nil, fmt.Errorf(".debug_info section is not present.", nil)
	}

	if len(debugInfoSections) > 1 {
		return nil, fmt.Errorf("More than one .debug_info sections.", nil)
	}

	reader, err := debugInfoSections[0].NewReader()
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
	debugInfoOffset := uint64(r.Size() - int64(r.Len()))
	die, exists := d.dieMap[debugInfoOffset]
	if exists {
		die.Parent = parent
		r.Seek(int64(die.debugInfoOffsetEnd), 0)

		return die, nil
	}

	abrrevTable, err := d.AbbrevTable(u.DebugAbbrevOffset)
	if err != nil {
		err = fmt.Errorf("Error getting abbrev table while reading a DIE tree.", err)
		return nil, err
	}

	abbrevCode, err := leb128.ReadUnsigned(r)
	if err != nil {
		return nil, fmt.Errorf("Error reading abbrev code of a DIE.", err)
	}

	// Return if its a NULL entry
	if abbrevCode == 0 {
		return nil, nil
	}

	abbrevEntry, exists := abrrevTable[abbrevCode]
	if !exists {
		return nil, fmt.Errorf("Invalid abbrev code for a DIE.", nil)
	}

	attributes := make(map[DwAt]Attribute)
	for _, attrForm := range abbrevEntry.AttrForms {
		attr, err := d.readAttr(u, r, attrForm.Name, attrForm.Form, en)
		if err != nil {
			msg := fmt.Sprintf(
				"Error reading value of attribute %d.\n%s",
				attrForm.Name, err.Error())
			err = fmt.Errorf(msg)
			return nil, err
		}
		attributes[attr.Name] = attr
	}

	die = new(DIE)
	die.Tag = abbrevEntry.Tag
	die.Parent = parent
	die.Attributes = attributes
	die.debugInfoOffsetStart = debugInfoOffset

	for abbrevEntry.HasChildren {
		childDie, err := d.readDIETreeHelper(u, r, en, die)
		if err != nil {
			err = fmt.Errorf(
				"Error reading child DIE tree of tag %d.\n%s",
				abbrevEntry.Tag, err.Error())
			return nil, err
		}

		if childDie == nil {
			break
		}

		die.Children = append(die.Children, childDie)
	}
	die.debugInfoOffsetEnd = uint64(r.Size() - int64(r.Len()))

	d.dieMap[debugInfoOffset] = die
	return die, nil
}
