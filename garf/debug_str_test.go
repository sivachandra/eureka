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

func TestDebugStr(t *testing.T) {
	dwData, err := LoadDwData("test_data/single_cu_linux_x86_64.exe")
	if err != nil {
		t.Errorf("Error loading DWARF from file.\n%s", err.Error())
		return
	}

	debugStrTbl, err := dwData.DebugStr()
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	str, err := debugStrTbl.ReadStr(0)
	if err != nil {
		t.Errorf("Expected string entry at offset 0.\n%s", err.Error())
		return
	} else {
		if str != "main" {
			t.Errorf("Expected string 'main' at offset 0.")
			return
		}
	}

	str, err = debugStrTbl.ReadStr(5)
	if err != nil {
		t.Errorf("Expected string entry at offset 5.\n%s", err.Error())
		return
	} else {
		if str != "main.c" {
			t.Errorf("Expected string 'main.c' at offset 5.")
			return
		}
	}

	str, err = debugStrTbl.ReadStr(12)
	if err != nil {
		t.Errorf("Expected string entry at offset 12.\n%s", err.Error())
		return
	} else {
		if str != "GNU C 4.8.2 -mtune=generic -march=x86-64 -g -fstack-protector" {
			t.Errorf("Unexpected string at offset 12.")
			return
		}
	}

	str, err = debugStrTbl.ReadStr(74)
	if err != nil {
		t.Errorf("Expected string entry at offset 74.\n%s", err.Error())
		return
	} else {
		if str != "/home/sivachandra/LAB/c++" {
			t.Errorf("Unexpected string at offset 74.")
			return
		}
	}
}
