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

func TestDwExprInMultipleCU(t *testing.T) {
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

	subProgDie := die.Children[0]
	if subProgDie.Tag != DW_TAG_subprogram {
		t.Errorf("Wrong DIE of comp unit 0's first child.")
		return
	}
	attr, exists := subProgDie.Attributes[DW_AT_frame_base]
	if !exists {
		t.Errorf("Missing attribute DW_AT_frame_base in the 'main' function.")
		return
	}
	expr := attr.Value.(DwExpr)
	if len(expr) != 1 {
		t.Errorf("Wrong length of the DWARF expression of DW_AT_frame_base of 'main'.")
		return
	}
	if expr[0].Op != DW_OP_call_frame_cfa {
		t.Errorf("Wrong opcode in the DWARF expression of DW_AT_frame_base of 'main'.")
		return
	}
	if len(expr[0].Operands) != 0 {
		t.Errorf("Wrong operand count in the DWARF expr of DW_AT_frame_base of 'main'.")
		return
	}

	callSiteParamDie := subProgDie.Children[0].Children[2].Children[0]
	if callSiteParamDie.Tag != DW_TAG_GNU_call_site_parameter {
		t.Errorf("Wrong TAG of call site parameter DIE in 'main'.")
		return
	}
	locAttr, exists := callSiteParamDie.Attributes[DW_AT_location]
	if !exists {
		t.Errorf("DW_AT_location missing in call site param DIE in 'main'.")
		return
	}
	expr = locAttr.Value.(DwExpr)
	if len(expr) != 1 {
		t.Errorf("Wrong length of the DWARF expression in call site param DIE in 'main'.")
		return
	}
	if expr[0].Op != DW_OP_reg5 {
		t.Errorf("Wrong opcode in the DWARF expression in call site param DIE in 'main'.")
		return
	}
	if len(expr[0].Operands) != 0 {
		t.Errorf("Wrong operand count in the DWARF expr in call site param DIE in 'main'.")
		return
	}

	locAttr, exists = callSiteParamDie.Attributes[DW_AT_GNU_call_site_value]
	if !exists {
		t.Errorf("DW_AT_call_site_value missing in call site param DIE in 'main'.")
		return
	}
	expr = locAttr.Value.(DwExpr)
	if len(expr) != 1 {
		t.Errorf("Wrong length of the DWARF expression in call site param DIE in 'main'.")
		return
	}
	if expr[0].Op != DW_OP_const1s {
		t.Errorf("Wrong opcode in the DWARF expression in call site param DIE in 'main'.")
		return
	}
	if len(expr[0].Operands) != 1 {
		t.Errorf("Wrong operand count in the DWARF expr in call site param DIE in 'main'.")
		return
	}
	val := expr[0].Operands[0].(int8)
	if val != -20 {
		t.Errorf("Wrong operand value in the DWARF expr in call site param DIE in 'main'.")
		return
	}

	// Items in comp unit 1
	die, err = compUnits[1].DIETree()
	if err != nil {
		t.Errorf("Error reading DIE tree of comp unit 1.\n%s", err.Error())
		return
	}

	if die.Tag != DW_TAG_compile_unit {
		t.Errorf("Wrong DIE tag for comp unit 1.")
	}

	subProgDie = die.Children[0]
	if subProgDie.Tag != DW_TAG_subprogram {
		t.Errorf("Wrong DIE of comp unit 1's first child.")
		return
	}
	attr, exists = subProgDie.Attributes[DW_AT_frame_base]
	if !exists {
		t.Errorf("Missing attribute DW_AT_frame_base in the 'f1' function.")
		return
	}
	expr = attr.Value.(DwExpr)
	if len(expr) != 1 {
		t.Errorf("Wrong length of the DWARF expression of DW_AT_frame_base of 'f1'.")
		return
	}
	if expr[0].Op != DW_OP_call_frame_cfa {
		t.Errorf("Wrong opcode in the DWARF expression of DW_AT_frame_base of 'f1'.")
		return
	}
	if len(expr[0].Operands) != 0 {
		t.Errorf("Wrong operand count in the DWARF expr of DW_AT_frame_base of 'f1'.")
		return
	}

	callSiteParamDie = subProgDie.Children[1].Children[0]
	if callSiteParamDie.Tag != DW_TAG_GNU_call_site_parameter {
		t.Errorf("Wrong TAG of call site parameter DIE in 'main'.")
		return
	}
	locAttr, exists = callSiteParamDie.Attributes[DW_AT_location]
	if !exists {
		t.Errorf("DW_AT_location missing in call site param DIE in 'f1'.")
		return
	}
	expr = locAttr.Value.(DwExpr)
	if len(expr) != 1 {
		t.Errorf("Wrong length of the DWARF expression in call site param DIE in 'f1'.")
		return
	}
	if expr[0].Op != DW_OP_reg5 {
		t.Errorf("Wrong opcode in the DWARF expression in call site param DIE in 'f1'.")
		return
	}
	if len(expr[0].Operands) != 0 {
		t.Errorf("Wrong operand count in the DWARF expr in call site param DIE in 'f1'.")
		return
	}

	locAttr, exists = callSiteParamDie.Attributes[DW_AT_GNU_call_site_value]
	if !exists {
		t.Errorf("DW_AT_call_site_value missing in call site param DIE in 'f1'.")
		return
	}
	expr = locAttr.Value.(DwExpr)
	if len(expr) != 1 {
		t.Errorf("Wrong length of the DWARF expression in call site param DIE in 'f1'.")
		return
	}
	if expr[0].Op != DW_OP_breg3 {
		t.Errorf("Wrong opcode in the DWARF expression in call site param DIE in 'f1'.")
		return
	}
	if len(expr[0].Operands) != 1 {
		t.Errorf("Wrong operand count in the DWARF expr in call site param DIE in 'f1'.")
		return
	}
	offset := expr[0].Operands[0].(int64)
	if offset != 0 {
		t.Errorf("Wrong operand value in the DWARF expr in call site param DIE in 'f1'.")
		return
	}

	// Items in comp unit 2
	die, err = compUnits[2].DIETree()
	if err != nil {
		t.Errorf("Error reading DIE tree of comp unit 2.\n%s", err.Error())
		return
	}

	if die.Tag != DW_TAG_compile_unit {
		t.Errorf("Wrong DIE tag for comp unit 2.")
	}

	subProgDie = die.Children[0]
	if subProgDie.Tag != DW_TAG_subprogram {
		t.Errorf("Wrong DIE of comp unit 2's first child.")
		return
	}
	attr, exists = subProgDie.Attributes[DW_AT_frame_base]
	if !exists {
		t.Errorf("Missing attribute DW_AT_frame_base in the 'f2' function.")
		return
	}
	expr = attr.Value.(DwExpr)
	if len(expr) != 1 {
		t.Errorf("Wrong length of the DWARF expression of DW_AT_frame_base of 'f2'.")
		return
	}
	if expr[0].Op != DW_OP_call_frame_cfa {
		t.Errorf("Wrong opcode in the DWARF expression of DW_AT_frame_base of 'f2'.")
		return
	}
	if len(expr[0].Operands) != 0 {
		t.Errorf("Wrong operand count in the DWARF expr of DW_AT_frame_base of 'f2'.")
		return
	}

	paramDie := subProgDie.Children[0]
	if paramDie.Tag != DW_TAG_formal_parameter {
		t.Errorf("Wrong TAG of param DIE in 'f2'.")
		return
	}
	attr, exists = paramDie.Attributes[DW_AT_location]
	if !exists {
		t.Errorf("DW_AT_location missing in param DIE in 'f2'.")
		return
	}
	expr = attr.Value.(DwExpr)
	if len(expr) != 1 {
		t.Errorf("Wrong length of the DWARF expression in param DIE in 'f2'.")
		return
	}
	if expr[0].Op != DW_OP_reg5 {
		t.Errorf("Wrong opcode in the DWARF expression in param DIE in 'f2'.")
		return
	}
	if len(expr[0].Operands) != 0 {
		t.Errorf("Wrong operand count in the DWARF expr in param DIE in 'f2'.")
		return
	}
}
