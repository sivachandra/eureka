// #############################################################################
// This file is part of the "golf" package of the "Eureka" project.
// It is distributed under the MIT License. Refer to the LICENSE file for more
// information.
//
// Website: http://www.github.com/sivachandra/eureka
// #############################################################################

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

	debugStrMap, err := dwData.DebugStr()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if len(debugStrMap) != 4 {
		t.Errorf("Incorrect number of entries in .debug_str.")
		return
	}

	str, exists := debugStrMap[0]
	if !exists {
		t.Errorf("Expected string entry at offset 0.")
		return
	} else {
		if str != "main" {
			t.Errorf("Expected string 'main' at offset 0.")
			return
		}
	}

	str, exists = debugStrMap[5]
	if !exists {
		t.Errorf("Expected string entry at offset 5.")
		return
	} else {
		if str != "main.c" {
			t.Errorf("Expected string 'main.c' at offset 5.")
			return
		}
	}

	str, exists = debugStrMap[12]
	if !exists {
		t.Errorf("Expected string entry at offset 12.")
		return
	} else {
		if str != "GNU C 4.8.2 -mtune=generic -march=x86-64 -g -fstack-protector" {
			t.Errorf("Unexpected string at offset 12.")
			return
		}
	}

	str, exists = debugStrMap[74]
	if !exists {
		t.Errorf("Expected string entry at offset 74.")
		return
	} else {
		if str != "/home/sivachandra/LAB/c++" {
			t.Errorf("Unexpected string at offset 74.")
			return
		}
	}
}
