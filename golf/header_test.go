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
