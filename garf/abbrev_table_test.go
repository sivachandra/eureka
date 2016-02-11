// #############################################################################
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

func TestSingleCUAbbrevTable(t *testing.T) {
	dwData, err := LoadDwData("test_data/single_cu_linux_x86_64.exe")
	if err != nil {
		t.Errorf("Error loading DWARF from file.\n%s", err.Error())
		return
	}

	// We will just fetch the abbrev table at offset 0 as the input ELF file
	// has only 1 CU.
	abbrevTable, err := dwData.AbbrevTable(0)
	if err != nil {
		t.Errorf("Error loading abbrev table.\n", err.Error())
		return
	}

	if len(abbrevTable) != 3 {
		t.Errorf("Incorrect length of abbrev table. Expected 3, found %d", len(abbrevTable))
		return
	}

	entry1, exists := abbrevTable[1]
	if !exists {
		t.Errorf("Entry with abbrev code 1 is missing.")
	}
	if entry1.AbbrevCode != 1 {
		t.Errorf("Wrong abbrev code in entry wih abbrev code 1.")
	}
	if entry1.Tag != DW_TAG_compile_unit {
		t.Errorf("Wrong tag for entry with abbrev code 1.")
	}
	if !entry1.HasChildren {
		t.Errorf("Wrong children description entry for entry wih abrev code 1.")
	}
	if len(entry1.AttrForms) != 7 {
		t.Errorf("Wrong number of attributes for entry with abbrev code 1.")
	} else {
		if entry1.AttrForms[0].Name != DW_AT_producer {
			t.Errorf("Wrong 0th attr for abbrev code 1.")
		}
		if entry1.AttrForms[6].Name != DW_AT_stmt_list {
			t.Errorf("Wrong 6th attr for abbrev code 1.")
		}
		if entry1.AttrForms[1].Form != DW_FORM_data1 {
			t.Errorf("Wrong form for 1st attr of entry with abbrev code 1.")
		}
	}

	entry2, exists := abbrevTable[2]
	if !exists {
		t.Errorf("Entry with abbrev code 2 is missing.")
	}
	if entry2.AbbrevCode != 2 {
		t.Errorf("Wrong abbrev code in entry wih abbrev code 2.")
	}
	if entry2.Tag != DW_TAG_subprogram {
		t.Errorf("Wrong tag for entry with abbrev code 2.")
	}
	if entry2.HasChildren {
		t.Errorf("Wrong children description entry for entry wih abrev code 2.")
	}
	if len(entry2.AttrForms) != 9 {
		t.Errorf("Wrong number of attributes for entry with abbrev code 2.")
	} else {
		if entry2.AttrForms[0].Name != DW_AT_external {
			t.Errorf("Wrong 0th attr for abbrev code 2.")
		}
		if entry2.AttrForms[6].Name != DW_AT_high_pc {
			t.Errorf("Wrong 6th attr for abbrev code 2.")
		}
		if entry2.AttrForms[1].Form != DW_FORM_strp {
			t.Errorf("Wrong form for 1st attr of entry with abbrev code 2.")
		}
	}

	entry3, exists := abbrevTable[3]
	if !exists {
		t.Errorf("Entry with abbrev code 3 is missing.")
	}
	if entry3.AbbrevCode != 3 {
		t.Errorf("Wrong abbrev code in entry wih abbrev code 3.")
	}
	if entry3.Tag != DW_TAG_base_type {
		t.Errorf("Wrong tag for entry with abbrev code 3.")
	}
	if entry3.HasChildren {
		t.Errorf("Wrong children description entry for entry wih abrev code 3.")
	}
	if len(entry3.AttrForms) != 3 {
		t.Errorf("Wrong number of attributes for entry with abbrev code 3.")
	} else {
		if entry3.AttrForms[0].Name != DW_AT_byte_size {
			t.Errorf("Wrong 0th attr for abbrev code 3.")
		}
		if entry3.AttrForms[2].Name != DW_AT_name {
			t.Errorf("Wrong 2nd attr for abbrev code 3.")
		}
		if entry3.AttrForms[1].Form != DW_FORM_data1 {
			t.Errorf("Wrong form for 1st attr of entry with abbrev code 3.")
		}
	}
}

func TestMultipleCUAbbrevTable(t *testing.T) {
	dwData, err := LoadDwData("test_data/multiple_cu_linux_x86_64.exe")
	if err != nil {
		t.Errorf("Error loading DWARF from file.\n%s", err.Error())
		return
	}

	abbrevTable1, err := dwData.AbbrevTable(0)
	if err != nil {
		t.Errorf("Error loading abbrev table.\n", err.Error())
		return
	}

	if len(abbrevTable1) != 11 {
		t.Errorf(
			"Incorrect length of abbrev table. Expected 11, found %d",
			len(abbrevTable1))
		return
	}

	entry, exists := abbrevTable1[2]
	if !exists {
		t.Errorf("Entry with abbrev code 2 is missing.")
	}

	if len(entry.AttrForms) != 10 {
		t.Errorf(
			"Wrong number of attributes for entry with abbrev code 2: %d.",
			len(entry.AttrForms))
		return
	}

	abbrevTable2, err := dwData.AbbrevTable(0x9d)
	if err != nil {
		t.Errorf("Error loading abbrev table.\n", err.Error())
		return
	}

	if len(abbrevTable2) != 8 {
		t.Errorf(
			"Incorrect length of abbrev table. Expected 8, found %d",
			len(abbrevTable2))
		return
	}

	abbrevTable3, err := dwData.AbbrevTable(0x118)
	if err != nil {
		t.Errorf("Error loading abbrev table.\n", err.Error())
		return
	}

	if len(abbrevTable3) != 4 {
		t.Errorf(
			"Incorrect length of abbrev table. Expected 8, found %d",
			len(abbrevTable3))
		return
	}
}
