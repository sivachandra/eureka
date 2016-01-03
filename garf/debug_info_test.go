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

func TestDebugInfoSingleCU(t *testing.T) {
	dwData, err := LoadDwData("test_data/linux_x86_64.exe")
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

	die, err := compUnits[0].DIETree()
	if err != nil {
		t.Errorf("Error fetching DIE tree for comp unit.\n%s", err.Error())
		return
	}
	if die == nil {
		t.Errorf("Empty DIE tree for comp unit.")
		return
	}

	if len(die.Children) != 2 {
		t.Errorf("Wrong number of children for the root DIE.")
		return
	}

	if len(die.Children[0].Children) != 0 || len(die.Children[1].Children) != 0 {
		t.Errorf("Wrong number of children for the leaf DIEs.")
		return
	}

	if len(die.Attributes) != 7 {
		t.Errorf("Wrong number of attributes for the root DIE.")
		return
	}

	if die.Tag != DW_TAG_compile_unit {
		t.Errorf("Wrong tag for the root DIE.")
	}

	attrs := die.Attributes
	for _, a := range attrs {
		switch a.Name {
		case DW_AT_producer:
			val := a.Value.(string)
			if val != "GNU C 4.8.2 -mtune=generic -march=x86-64 -g -fstack-protector" {
				t.Errorf("Wrong value for DW_AT_producer attribute of root DIE.")
			}
		case DW_AT_name:
			val := a.Value.(string)
			if val != "main.c" {
				t.Errorf("Wrong value for DW_AT_name attribute of root DIE.")
			}
		case DW_AT_language:
			val := a.Value.(DwLang)
			if val != DW_LANG_C89 {
				t.Errorf("Wrong value for DW_AT_language attribute of root DIE.")
			}
		case DW_AT_stmt_list:
			val := a.Value.(uint64)
			if val != 0 {
				t.Errorf("Wrong value of DW_AT_stmt_list attribute of root DIE.")
			}
		case DW_AT_comp_dir:
			val := a.Value.(string)
			if val != "/home/sivachandra/LAB/c++" {
				t.Errorf("Wrong value for DW_AT_name attribute of root DIE.")
			}
		case DW_AT_low_pc:
			val := a.Value.(uint64)
			if val != 0x4004ed {
				t.Errorf("Wrong value of DW_AT_stmt_list attribute of root DIE.")
			}
		case DW_AT_high_pc:
			val := a.Value.(uint64)
			if val != 0xb {
				t.Errorf("Wrong value of DW_AT_stmt_list attribute of root DIE.")
			}
		default:
			t.Errorf("Unexpected attribute for root DIE.")
		}
	}

	child := die.Children[0]
	if child.Tag != DW_TAG_subprogram {
		t.Errorf("Wrong tag for child DIE.")
	}
	if len(child.Attributes) != 9 {
		t.Errorf("Wrong number of attributes for child DIE.")
	}
	attrs = child.Attributes
	for _, a := range attrs {
		switch a.Name {
		case DW_AT_external:
			val := a.Value.(bool)
			if !val {
				t.Errorf("Unexpected value for attr DW_AT_external of child DIE.")
			}
		case DW_AT_name:
			val := a.Value.(string)
			if val != "main" {
				t.Errorf("Wrong value for DW_AT_name attribute of child DIE.")
			}
		case DW_AT_decl_file:
			val := a.Value.(uint32)
			if val != 1 {
				t.Errorf("Unexpected value for attr DW_AT_decl_file of child DIE.")
			}
		case DW_AT_decl_line:
			val := a.Value.(uint32)
			if val != 2 {
				t.Errorf("Unexpected value for attr DW_AT_decl_line of child DIE.")
			}
		case DW_AT_type:
			val := a.Value.(*DIE)
			if val != die.Children[1] {
				t.Errorf("Incorrect reference DIE linking.")
			}
		case DW_AT_low_pc:
			val := a.Value.(uint64)
			if val != 0x4004ed {
				t.Errorf("Unexpected value for attr DW_AT_low_pc of child DIE.")
			}
		case DW_AT_high_pc:
			val := a.Value.(uint64)
			if val != 0xb {
				t.Errorf("Unexpected value for attr DW_AT_high_pc of child DIE.")
			}
		case DW_AT_frame_base:
			val := a.Value.([]byte)
			if len(val) != 1 {
				t.Errorf("Unexpected length of byte slice for DW_AT_frame_base.")
			}
		case DW_AT_GNU_all_call_sites:
			val := a.Value.(bool)
			if !val {
				t.Errorf("Unexpected DW_AT_GNU_all_call_sites val of child DIE.")
			}
		default:
			t.Errorf("Unexpected attribute for child DIE.")
		}
	}

	child = die.Children[1]
	if child.Tag != DW_TAG_base_type {
		t.Errorf("Wrong tag for child DIE.")
	}
	if len(child.Attributes) != 3 {
		t.Errorf("Wrong number of attributes for child DIE.")
	}
	attrs = child.Attributes
	for _, a := range attrs {
		switch a.Name {
		case DW_AT_byte_size:
			val := a.Value.(uint32)
			if val != 4 {
				t.Errorf("Unexpected value for attr DW_AT_encoding of child DIE.")
			}
		case DW_AT_encoding:
			val := a.Value.(DwAte)
			if val != DW_ATE_signed {
				t.Errorf("Unexpected value of attr DW_AT_encoding of child DIE.")
			}
		case DW_AT_name:
			val := a.Value.(string)
			if val != "int" {
				t.Errorf("Wrong value for DW_AT_name attribute of child DIE.")
			}
		default:
			t.Errorf("Unexpected attribute for child DIE.")
		}
	}
}
