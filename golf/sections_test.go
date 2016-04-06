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
	"testing"
)

func TestHeaderCount(t *testing.T) {
	elf, err := Read("test_data/linux_x86_64.exe")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	sectCount := elf.Header().SectHdrCount()
	sectHdrTable := elf.SectHdrTbl()
	if uint16(len(sectHdrTable)) != sectCount {
		t.Error("Mismatch in the section header count.")
		return
	}
}

func TestSections(t *testing.T) {
	elf, err := Read("test_data/linux_x86_64.exe")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	sectMap := elf.SectMap()
	sectNameSect := sectMap[NameSectNameTbl][0]
	strTblData, err := sectNameSect.Data()
	strTbl, err := BuildStrTbl(strTblData)
	if len(strTbl) != 34 {
		t.Errorf(
			"Incorrect entry count in section name table. Expecting 34, found %d.\n",
			len(strTbl))
		return
	}

	symTabSect := sectMap[NameSymTab][0]
	symTabData, err := symTabSect.Data()
	symTab, err := BuildSymTab(
		symTabData, symTabSect.SectHdr(), elf.Header().ELFIdent().Endianess)
	if err != nil {
		t.Errorf("Unable to read .symtab.\n%s", err.Error())
		return
	}
	symCount := 0
	for _, symList := range symTab {
		symCount += len(symList)
	}
	if symCount != 69 {
		t.Errorf("Incorrect size of symbol table.\nExpected 69, got %d.\n", symCount)
		return
	}

	strTblSect := sectMap[NameSymNameTbl][0]
	strTblData, err = strTblSect.Data()
	if err != nil {
		t.Errorf("Unable to read .strtab.\n%s", err.Error())
		return
	}
	strTbl, err = BuildStrTbl(strTblData)
	if len(strTbl) != 36 {
		t.Errorf("Incorrect size of .strtab.\nExpected 36, got %d\n", len(strTbl))
		return
	}

	symTabSect = sectMap[NameDynSymTab][0]
	symTabData, err = symTabSect.Data()
	symTab, err = BuildSymTab(
		symTabData, symTabSect.SectHdr(), elf.Header().ELFIdent().Endianess)
	if err != nil {
		t.Errorf("Unable to read .dynsym.\n%s", err.Error())
		return
	}

	strTblSect = sectMap[NameDynSymNameTbl][0]
	strTblData, err = strTblSect.Data()
	if err != nil {
		t.Errorf("Unable to read .dynstr.\n%s", err.Error())
		return
	}
	strTbl, err = BuildStrTbl(strTblData)
	if len(strTbl) != 5 {
		t.Errorf("Incorrect size of .dynstr.\nExpected 5, got %d\n", len(strTbl))
		return
	}
}
