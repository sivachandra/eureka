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

func TestStdVectorStringGcc(t *testing.T) {
	dwData, err := LoadDwData("test_data/std_vector_string_gcc-4.8.4.exe")
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

	_, err = compUnits[0].DIETree()
	if err != nil {
		t.Errorf("Error reading DIE tree of comp unit 0.\n%s", err.Error())
		return
	}

	_, err = compUnits[0].LineNumberInfo()
	if err != nil {
		t.Errorf("Error reading line number info comp unit 0.\n%s", err.Error())
		return
	}
}
