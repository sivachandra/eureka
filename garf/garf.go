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
	Attributes []Attribute
	Parent     *DIE
	Children   []*DIE
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
	// DIE tree of this unit begins.
	DebugInfoOffset uint64

	// The complete DIE tree of this unit.
	dieTree *DIE
}

func (u DwUnit) DIETree() (*DIE, error) {
	if u.dieTree != nil {
		return u.dieTree, nil
	}

	var err error
	u.dieTree, err = u.Parent.readDIETree(&u, u.DebugInfoOffset)
	return u.dieTree, err
}

type DebugStrMap map[uint64]string

type DwData struct {
	fileName    string
	elf         *golf.ELF
	abbrevTable AbbrevTable
	debugStrMap DebugStrMap
	compUnits   []DwUnit
	typeUnits   []DwUnit

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

	return dwData, nil
}

func (d *DwData) ELFData() *golf.ELF {
	return d.elf
}

func (d *DwData) FileName() string {
	return d.fileName
}

func (d *DwData) formatError(msg string, err error) error {
	m := fmt.Sprintf("[%s] %s", d.fileName, msg)

	if err != nil {
		m = fmt.Sprintf("%s\n[%s] %s", m, d.fileName, err.Error())
	}

	return fmt.Errorf("%s", m)
}

func (d *DwData) AbbrevTable() (AbbrevTable, error) {
	if d.abbrevTable != nil {
		return d.abbrevTable, nil
	}

	table := make(AbbrevTable)

	sectMap := d.elf.SectMap()
	debugAbbrevSections, exists := sectMap[".debug_abbrev"]
	if !exists {
		return nil, d.formatError(".debug_abbrev section is not present.", nil)
	}

	if len(debugAbbrevSections) > 1 {
		return nil, d.formatError("More than one .debug_abbrev sections.", nil)
	}

	debugAbbrevData, err := debugAbbrevSections[0].RawData()
	if err != nil {
		return nil, d.formatError("Error fetching .debug_abbrev data.", err)
	}

	reader := bytes.NewReader(debugAbbrevData)
	for true {
		abbrevCode, err := leb128.ReadUnsigned(reader)
		if err != nil {
			return nil, d.formatError("Error reading abbreviation code.", nil)
		}
		if abbrevCode == NullAbbrevEntry {
			break
		}

		tag, err := leb128.ReadUnsigned(reader)
		if err != nil {
			msg := fmt.Sprintf("Error reading tag for abbrev code %d.", abbrevCode)
			return nil, d.formatError(msg, err)
		}

		hasChildren, err := reader.ReadByte()
		if err != nil {
			msg := fmt.Sprintf(
				"Error reading child determination entry for abbrev code %d.",
				abbrevCode)
			return nil, d.formatError(msg, err)
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
				return nil, d.formatError(msg, err)
			}

			form, err := leb128.ReadUnsigned(reader)
			if err != nil {
				msg := fmt.Sprintf(
					"Error reading an attr form of entry with abbrev code %d.",
					abbrevCode)
				return nil, d.formatError(msg, err)
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

	d.abbrevTable = table
	return table, nil
}

func (d *DwData) CompUnits() ([]DwUnit, error) {
	if d.compUnits != nil {
		return d.compUnits, nil
	}

	sectMap := d.elf.SectMap()
	debugInfoSections, exists := sectMap[".debug_info"]
	if !exists {
		return nil, d.formatError(".debug_info section is not present.", nil)
	}

	if len(debugInfoSections) > 1 {
		return nil, d.formatError("More than one .debug_info sections.", nil)
	}

	reader, err := debugInfoSections[0].NewSectReader()
	if err != nil {
		return nil, d.formatError("Error fetching .debug_info section reader.", err)
	}
	defer reader.Finish()

	d.compUnits = make([]DwUnit, 0)
	en := d.elf.Endianess()
	for true {
		if reader.Len() == 0 {
			break
		}

		var length uint64
		var format DwFormat
		var size32 uint32

		err := binary.Read(reader, en, &size32)
		if err != nil {
			err = d.formatError(
				"Error reading first 32 bits of length of a unit in .debug_info.",
				err)
			return nil, err
		}

		if size32 == 0xffffffff {
			format = DwFormat64
			var size64 uint64

			err := binary.Read(reader, en, &size64)
			if err != nil {
				err = d.formatError(
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
			err = d.formatError("Error reading version of a unit in .debug_info.", err)
			return nil, err
		}

		unitType := DW_UT_compile
		if version >= 5 {
			err = binary.Read(reader, en, &unitType)
			if err != nil {
				err = d.formatError(
					"Error reading unit type of a unit in .debug_info.", err)
				return nil, err
			}
		}

		var debugAbbrevOffset uint64
		if format == DwFormat32 {
			var offset uint32
			err = binary.Read(reader, en, &offset)
			if err != nil {
				err = d.formatError(
					"Error reading 32-bit debug abbrev offset of a unit.", err)
				return nil, err
			}

			debugAbbrevOffset = uint64(offset)
		} else {
			err = binary.Read(reader, en, &debugAbbrevOffset)
			if err != nil {
				err = d.formatError(
					"Error reading 64-bit debug abbrev offset of a unit.", err)
				return nil, err
			}
		}

		var addrSize byte
		err = binary.Read(reader, en, &addrSize)
		if err != nil {
			err = d.formatError(
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
			cu.DebugInfoOffset = reader.Size() - reader.Len()
			d.compUnits = append(d.compUnits, cu)
			reader.Seek(int64(cu.Size+cu.DebugInfoOffset), 0)
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
		return nil, d.formatError(".debug_str section is not present.", nil)
	}

	if len(debugStrSections) > 1 {
		return nil, d.formatError("More than one .debug_str sections.", nil)
	}

	debugStrData, err := debugStrSections[0].RawData()
	if err != nil {
		return nil, d.formatError("Error fetching .debug_str data.", err)
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
		return nil, d.formatError(".debug_info section is not present.", nil)
	}

	if len(debugInfoSections) > 1 {
		return nil, d.formatError("More than one .debug_info sections.", nil)
	}

	reader, err := debugInfoSections[0].NewSectReader()
	if err != nil {
		return nil, d.formatError("Error fetching .debug_info section reader.", err)
	}
	defer reader.Finish()

	_, err = reader.Seek(int64(offset), 0)
	if err != nil {
		err = fmt.Errorf(
			"Error seeking to the DIE offset to read the DIE tree.\n%s", err.Error())
		return nil, err
	}

	return d.readDIETreeHelper(u, reader, d.elf.Endianess(), nil)
}

func (d *DwData) readDIETreeHelper(
	u *DwUnit, r *golf.SectReader, en binary.ByteOrder, parent *DIE) (*DIE, error) {
	abrrevTable, err := d.AbbrevTable()
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

	attributes := make([]Attribute, 0)
	for _, attrForm := range abbrevEntry.AttrForms {
		attr, err := d.readAttr(u, r, attrForm.Name, attrForm.Form, en)
		if err != nil {
			err = fmt.Errorf("Error reading an attribute value for a DIE.\n%s", err)
			return nil, err
		}
		attributes = append(attributes, attr)
	}

	die := new(DIE)
	die.Tag = abbrevEntry.Tag
	die.Parent = parent
	die.Attributes = attributes

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

	return die, nil
}
