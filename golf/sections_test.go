// This file is part of the "golf" package of the "Eureka" project.
// It is distributed under the MIT License. Refer to the LICENSE file for more
// information.
//
// Website: http://www.github.com/sivachandra/eureka
// #############################################################################

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
	strTblData, err := sectNameSect.Data(elf.Header().ELFIdent().Endianess)
	strTbl := strTblData.(StrTbl)
	if len(strTbl) != 34 {
		t.Errorf(
			"Incorrect entry count in section name table. Expecting 34, found %d.\n",
			len(strTbl))
		return
	}

	symTabSect := sectMap[NameSymTab][0]
	symTab, err := symTabSect.Data(elf.Header().ELFIdent().Endianess)
	if err != nil {
		t.Errorf("Unable to read .symtab.\n%s", err.Error())
		return
	}
	symCount := 0
	for _, symList := range symTab.(SymTab) {
		symCount += len(symList)
	}
	if symCount != 69 {
		t.Errorf("Incorrect size of symbol table.\nExpected 69, got %d.\n", symCount)
		return
	}

	strTblSect := sectMap[NameSymNameTbl][0]
	strTblData, err = strTblSect.Data(elf.Header().ELFIdent().Endianess)
	if err != nil {
		t.Errorf("Unable to read .strtab.\n%s", err.Error())
		return
	}
	strTbl = strTblData.(StrTbl)
	if len(strTbl) != 36 {
		t.Errorf("Incorrect size of .strtab.\nExpected 36, got %d\n", len(strTbl))
		return
	}

	symTabSect = sectMap[NameDynSymTab][0]
	symTab, err = symTabSect.Data(elf.Header().ELFIdent().Endianess)
	if err != nil {
		t.Errorf("Unable to read .symtab.\n%s", err.Error())
		return
	}

	strTblSect = sectMap[NameDynSymNameTbl][0]
	strTblData, err = strTblSect.Data(elf.Header().ELFIdent().Endianess)
	if err != nil {
		t.Errorf("Unable to read .dynstr.\n%s", err.Error())
		return
	}
	strTbl = strTblData.(StrTbl)
	if len(strTbl) != 5 {
		t.Errorf("Incorrect size of .dynstr.\nExpected 5, got %d\n", len(strTbl))
		return
	}
}
