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
