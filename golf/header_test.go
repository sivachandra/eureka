// #############################################################################
// This file is part of the "golf" package of the "Eureka" project.
// It is distributed under the MIT License. Refer to the LICENSE file for more
// information.
//
// Website: http://www.github.com/sivachandra/eureka
// #############################################################################

package golf

import (
	"testing"
)

func testMagicNumber(t *testing.T, elf *ELF) {
	mag := elf.Header().ELFIdent().MagicNumber
	if len(mag) != 4 {
		t.Errorf("Incorrect magic number length.")
		return
	}

	if mag[0] != Mag0 || mag[1] != Mag1 || mag[2] != Mag2 || mag[3] != Mag3 {
		t.Errorf("Incorrect magic number.")
		return
	}
}

func testClass(t *testing.T, elf *ELF) {
	if elf.Header().ELFIdent().Class != Class64 {
		t.Errorf("Incorrect class.")
		return
	}
}

func testEndianess(t *testing.T, elf *ELF) {
	if elf.Header().ELFIdent().Endianess != LittleEndian {
		t.Errorf("Incorrect endianess.")
		return
	}
}

func testType(t *testing.T, elf *ELF) {
	if elf.Header().Type() != TypeExecutable {
		t.Errorf("Incorrect file type %d.", elf.Header().Type())
		return
	}
}

func testMachine(t *testing.T, elf *ELF) {
	if elf.Header().Machine() != MachineX86_64 {
		t.Errorf("Incorrect machine type %d.", elf.Header().Machine())
		return
	}
}

func TestHeader(t *testing.T) {
	elf, err := Read("test_data/linux_x86_64.exe")
	if err != nil {
		t.Errorf("Reading a linux x86_64 exe file failed.")
		return
	}

	testMagicNumber(t, elf)
	testClass(t, elf)
	testEndianess(t, elf)
	testType(t, elf)
	testMachine(t, elf)
}
