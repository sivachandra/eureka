// This file is part of the "garf" package of the "Eureka" project.
// It is distributed under the MIT License. Refer to the LICENSE file for more
// information.
//
// Website: http://www.github.com/sivachandra/eureka
// #############################################################################

package garf

import (
	"testing"
)

func TestLineInfoSingleCU(t *testing.T) {
	dwData, err := LoadDwData("test_data/single_cu_linux_x86_64.exe")
	if err != nil {
		t.Errorf("Error loading DWARF from file.\n%s", err.Error())
		return
	}

	compUnits, err := dwData.CompUnits()
	if err != nil {
		t.Errorf("Error reading comp units.\n%s", err.Error())
		return
	}
	if len(compUnits) != 1 {
		t.Errorf("Wrong number of comp units: %d", len(compUnits))
		return
	}

	lnInfo, err := compUnits[0].LineNumberInfo()
	if err != nil {
		t.Errorf("Error getting comp unit line number info.\n%s", err.Error())
		return
	}

	if lnInfo.Version != 2 {
		t.Errorf("Wrong version of line info.")
		return
	}

	if lnInfo.minInstrLength != 1 {
		t.Errorf(
			"Wrong minimum instruction length. Expected 1, got %d.",
			lnInfo.minInstrLength)
		return
	}

	if lnInfo.defaultIsStmt == 0 {
		t.Errorf(
			"Wrong default_is_stmt value. Expected non-zero, got %d.",
			lnInfo.defaultIsStmt)
		return
	}

	if lnInfo.lineBase != -5 {
		t.Errorf("Wrong line base value. Expected -5, got %d.", lnInfo.lineBase)
		return
	}

	if lnInfo.lineRange != 14 {
		t.Errorf("Wrong line range value. Expected 14, got %d.", lnInfo.lineRange)
		return
	}

	if lnInfo.opcodeBase != 13 {
		t.Errorf("Wrong opcode base value. Expected 13, got %d.", lnInfo.opcodeBase)
		return
	}

        if len(lnInfo.operandCountTbl) != 12 {
		t.Errorf(
			"Wrong length of operand count table. Expected 12, got %d.",
			len(lnInfo.operandCountTbl))
		return
	}

	if len(lnInfo.Directories) != 0 {
		t.Errorf(
			"Wrong number of directory entries. Expected 0, got %d.",
			len(lnInfo.Directories))
		return
	}

	if len(lnInfo.Files) != 1 {
		t.Errorf(
			"Wrong number of file entries. Expected 1, got %d.",
			len(lnInfo.Files))
		return
	}

	if lnInfo.Files[0].Path != "main.c" {
		t.Errorf(
			"Wrong file name in file entry. Expected 'main.c', got '%s'.",
			lnInfo.Files[0].Path)
		return
	}

	if len(lnInfo.Program) != 6 {
		t.Errorf(
			"Wrong number of instrs in line number program. Expected 6, got %d.",
			len(lnInfo.Program))
		return
	}
}
