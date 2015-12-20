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
	"fmt"
)

import (
	"eureka/golf"
	"eureka/utils/leb128"
)

type AttrFormPair struct {
	Attr DwAt
	Form DwForm
}

type AbbrevEntry struct {
	AbbrevCode uint64
	Tag DwTag
	HasChildren bool
	Attributes []AttrFormPair
}

type AbbrevTable map[uint64]AbbrevEntry

type DwFile struct {
	fileName string
	elf *golf.ELF
	abbrevTable AbbrevTable
}

func LoadDwFile(fileName string) (*DwFile, error) {
	dwFile := new(DwFile)
	var err error

	dwFile.fileName = fileName
	dwFile.elf, err = golf.Read(fileName)
	if err != nil {
		err = fmt.Errorf("Error loading ELF info from '%s'.\n%s", fileName, err.Error())
		return nil, err
	}

	return dwFile, nil
}

func (d *DwFile) ELFData() *golf.ELF {
	return d.elf
}

func (d *DwFile) FileName() string {
	return d.fileName
}

func (d *DwFile) formatError(msg string, err error) error {
	m := fmt.Sprintf("[%s] %s", d.fileName, msg)

	if err != nil {
		m = fmt.Sprintf("%s\n%s", m, err.Error())
	}

	return fmt.Errorf("%s", m)
}

func (d *DwFile) GetAbbrevTable() (AbbrevTable, error) {
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
		entry.Attributes = make([]AttrFormPair, 0)

		for true {
			attr, err := leb128.ReadUnsigned(reader)
			if err != nil {
				msg := fmt.Sprintf(
					"Error reading an attribute name of entry with abbrev code %d.",
					abbrevCode)
				return nil, d.formatError(msg, err)
			}

			form, err := leb128.ReadUnsigned(reader)
			if err != nil {
				msg := fmt.Sprintf(
					"Error reading an attribute form of entry with abbrev code %d.",
					abbrevCode)
				return nil, d.formatError(msg, err)
			}

			if form == 0 && attr == 0 {
				break
			}

			var pair AttrFormPair
			pair.Attr = DwAt(attr)
			pair.Form = DwForm(form)
			entry.Attributes = append(entry.Attributes, pair)
		}

		table[entry.AbbrevCode] = entry
	}

	d.abbrevTable = table
	return table, nil
}
