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

func TestDebugInfoMultipleCU(t *testing.T) {
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

	die, err := compUnits[0].DIETree()
	if err != nil {
		t.Errorf("Error reading DIE tree of comp unit 0.\n%s", err.Error())
		return
	}

	if die.Tag != DW_TAG_compile_unit {
		t.Errorf("Wrong DIE tag for comp unit 0.")
	}

	die, err = compUnits[1].DIETree()
	if err != nil {
		t.Errorf("Error reading DIE tree of comp unit 1.\n%s", err.Error())
		return
	}

	if die.Tag != DW_TAG_compile_unit {
		t.Errorf("Wrong DIE tag for comp unit 1.")
	}

	if len(die.Children) != 3 {
		t.Errorf("Wrong number of children for the root of the DIE tree of comp unit 1.")
	}

	childDie := die.Children[2].Children[0]
	if childDie.Tag != DW_TAG_formal_parameter {
		t.Errorf("Wrong tag for a DIE in comp unit 1.")
	}

	typeDie := childDie.Attributes[DW_AT_type].Value.(*DIE)
	if typeDie.Attributes[DW_AT_name].Value.(string) != "int" {
		t.Errorf("Wrong type name for type DIE in comp unit 1.")
	}

	die, err = compUnits[2].DIETree()
	if err != nil {
		t.Errorf("Error reading DIE tree of comp unit 2.\n%s", err.Error())
		return
	}

	if die.Tag != DW_TAG_compile_unit {
		t.Errorf("Wrong DIE tag for comp unit 2.")
	}
}
