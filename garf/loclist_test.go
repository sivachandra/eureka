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

func TestLocLists(t *testing.T) {
	dwData, err := LoadDwData("test_data/multiple_cu_linux_x86_64.exe")
	if err != nil {
		t.Errorf("Error loading DWARF from file.\n%s", err.Error())
		return
	}

	compUnits, err := dwData.CompUnits()
	if err != nil {
		t.Errorf("Error reading comp units.\n%s", err.Error())
		return
	}
	if len(compUnits) != 3 {
		t.Errorf("Wrong number of comp units: %d", len(compUnits))
		return
	}

	// Items in comp unit 0
	die, err := compUnits[0].DIETree()
	if err != nil {
		t.Errorf("Error reading DIE tree of comp unit 0.\n%s", err.Error())
		return
	}

	if die.Tag != DW_TAG_compile_unit {
		t.Errorf("Wrong DIE tag for comp unit 0.")
	}

	varDie := die.Children[0].Children[0].Children[0]
	attr, exists := varDie.Attributes[DW_AT_location]
	if !exists {
		t.Errorf("Missing attribute DW_AT_location on a var DIE in comp unit 0.")
		return
	}

	locList := attr.Value.(LocList)
	if len(locList) != 2 {
		t.Errorf("Wrong size of loc list in comp unit 0.")
		return
	}

	entry := locList[0]
	normalEntry := entry.(NormalLocListEntry)
	if normalEntry.Begin != 0x40040e {
		t.Errorf("Incorrect begin offset in normal loc list entry in comp unit 0.")
		return
	}
	if normalEntry.End != 0x400418 {
		t.Errorf("Incorrect end offset in normal loc list entry in comp unit 0.")
		return
	}

	if len(normalEntry.Loc) != 1 {
		t.Errorf("Wrong length of loc expr in loc list entry in comp unit 0.")
		return
	}

	if normalEntry.Loc[0].Op != DW_OP_reg0 {
		t.Errorf("Wrong operation in loc expr in loc list entry in comp unit 0.")
		return
	}

	_ = locList[1].(EndOfListLocListEntry)

	// Items in comp unit 1
	die, err = compUnits[1].DIETree()
	if err != nil {
		t.Errorf("Error reading DIE tree of comp unit 1.\n%s", err.Error())
		return
	}

	if die.Tag != DW_TAG_compile_unit {
		t.Errorf("Wrong DIE tag for comp unit 1.")
	}

	paramDie := die.Children[0].Children[0]
	attr, exists = paramDie.Attributes[DW_AT_location]
	if !exists {
		t.Errorf("Missing attribute DW_AT_location on a param DIE in comp unit 1.")
		return
	}

	locList = attr.Value.(LocList)
	if len(locList) != 4 {
		t.Errorf("Wrong size of loc list in comp unit 1.")
		return
	}

	normalEntry = locList[0].(NormalLocListEntry)
	if normalEntry.Loc[0].Op != DW_OP_reg5 {
		t.Errorf("Wrong operation in loc expr in first loc list entry in comp unit 1.")
		return
	}

	normalEntry = locList[1].(NormalLocListEntry)
	if normalEntry.Loc[0].Op != DW_OP_reg3 {
		t.Errorf("Wrong operation in loc expr in second loc list entry in comp unit 1.")
		return
	}

	normalEntry = locList[2].(NormalLocListEntry)
	if normalEntry.Loc[0].Op != DW_OP_GNU_entry_value {
		t.Errorf("Wrong operation in loc expr in third loc list entry in comp unit 1.")
		return
	}
	if normalEntry.Loc[1].Op != DW_OP_stack_value {
		t.Errorf("Wrong operation in loc expr in third loc list entry in comp unit 1.")
		return
	}

	_ = locList[3].(EndOfListLocListEntry)
}
