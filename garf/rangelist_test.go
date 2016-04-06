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

func TestRangeLists(t *testing.T) {
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

	attr, exists := die.Attributes[DW_AT_ranges]
	if !exists {
		t.Errorf("Missing attribute DW_AT_ranges for comp unit 0.")
		return
	}

	rangeList := attr.Value.(RangeList)
	if len(rangeList) != 2 {
		t.Errorf("Wrong size of range list for comp unit 0.")
		return
	}

	entry := rangeList[0]
	normalEntry := entry.(RangeListEntryNormal)
	if normalEntry.Begin != 0x400400 {
		t.Errorf("Incorrect begin offset in normal range list entry for comp unit 0.")
		return
	}
	if normalEntry.End != 0x400419 {
		t.Errorf("Incorrect end offset in normal range list entry for comp unit 0.")
		return
	}

	_ = rangeList[1].(RangeListEntryEndOfList)

	varDie := die.Children[0].Children[0]
	attr, exists = varDie.Attributes[DW_AT_ranges]
	if !exists {
		t.Errorf("Missing attribute DW_AT_ranges on a lexical block DIE in comp unit 0.")
		return
	}

	rangeList = attr.Value.(RangeList)
	if len(rangeList) != 3 {
		t.Errorf("Wrong size of range list for a lexical block DIE in comp unit 0.")
		return
	}

	entry = rangeList[0]
	normalEntry = entry.(RangeListEntryNormal)
	if normalEntry.Begin != 0x400404 {
		t.Errorf("Incorrect begin offset in normal range list entry for lexical block DIE.")
		return
	}
	if normalEntry.End != 0x40040e {
		t.Errorf("Incorrect end offset in normal range list entry for lexical block DIE.")
		return
	}

	entry = rangeList[1]
	normalEntry = entry.(RangeListEntryNormal)
	if normalEntry.Begin != 0x400412 {
		t.Errorf("Incorrect begin offset in normal range list entry for lexical block DIE.")
		return
	}
	if normalEntry.End != 0x400419 {
		t.Errorf("Incorrect end offset in normal range list entry for lexical block DIE.")
		return
	}

	_ = rangeList[2].(RangeListEntryEndOfList)
}
