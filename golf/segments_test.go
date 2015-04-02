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

func testSegHdr(
	t *testing.T,
	segHdr SegHdr, segIndex uint16,
	segType uint32, flags uint32,
	offset uint64, align uint64,
	fileSize uint64, memSize uint64,
	virtAddr uint64, physAddr uint64) {
	if segHdr.Type() != segType {
		t.Errorf("Type of segment header at index %d is incorrect.", segIndex)
	}
	if segHdr.Flags() != flags {
		t.Errorf("Wrong flags in segment header at index %d.", segIndex)
	}
	if segHdr.Offset() != offset {
		t.Errorf("Bad offset in segment header at index %d.", segIndex)
	}
	if segHdr.PhysicalAddress() != physAddr {
		t.Errorf("Bad physical address in segment header at index %d.", segIndex)
	}
	if segHdr.VirtualAddress() != virtAddr {
		t.Errorf("Bad virtual address in segment header at index %d.", segIndex)
	}
	if segHdr.FileSize() != fileSize {
		t.Errorf("Bad file size in segment header at index %d.", segIndex)
	}
	if segHdr.MemSize() != memSize {
		t.Errorf("Bad memory size in segment header at index %d.", segIndex)
	}
	if segHdr.Alignment() != align {
		t.Errorf("Bad alignment in segment header at index %d.", segIndex)
	}
}

func TestProgHdrTbl(t *testing.T) {
	elf, err := Read("test_data/linux_x86_64.exe")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	segCount := elf.Header().ProgHdrCount()
	progHdrTbl := elf.ProgHdrTbl()
	if uint16(len(progHdrTbl)) != segCount {
		t.Error("Mismatch in the section header count.")
		return
	}
	if segCount != uint16(9) {
		t.Error("Wrong number of segments. Expecting 9.")
		return
	}

	testSegHdr(
		t,
		progHdrTbl[0], 0,
		SegTypeProgHdr, SegFlagsReadable|SegFlagsExecutable,
		uint64(0x000040), uint64(0x8),
		uint64(0x0001f8), uint64(0x0001f8),
		uint64(0x0000000000400040), uint64(0x0000000000400040))

	testSegHdr(
		t,
		progHdrTbl[1], 1,
		SegTypeInterp, SegFlagsReadable,
		uint64(0x000238), uint64(0x1),
		uint64(0x00001c), uint64(0x00001c),
		uint64(0x0000000000400238), uint64(0x0000000000400238))

	testSegHdr(
		t,
		progHdrTbl[2], 2,
		SegTypeLoad, SegFlagsReadable|SegFlagsExecutable,
		uint64(0x000000), uint64(0x200000),
		uint64(0x0006ac), uint64(0x0006ac),
		uint64(0x0000000000400000), uint64(0x0000000000400000))

	testSegHdr(
		t,
		progHdrTbl[3], 3,
		SegTypeLoad, SegFlagsReadable|SegFlagsWritable,
		uint64(0x000e10), uint64(0x200000),
		uint64(0x000228), uint64(0x000230),
		uint64(0x0000000000600e10), uint64(0x0000000000600e10))

	testSegHdr(
		t,
		progHdrTbl[4], 4,
		SegTypeDynamic, SegFlagsReadable|SegFlagsWritable,
		uint64(0x000e28), uint64(0x8),
		uint64(0x0001d0), uint64(0x0001d0),
		uint64(0x0000000000600e28), uint64(0x0000000000600e28))

	testSegHdr(
		t,
		progHdrTbl[5], 5,
		SegTypeNote, SegFlagsReadable,
		uint64(0x000254), uint64(0x4),
		uint64(0x000044), uint64(0x000044),
		uint64(0x0000000000400254), uint64(0x0000000000400254))

	testSegHdr(
		t,
		progHdrTbl[6], 6,
		SegTypeGnuEHFrame, SegFlagsReadable,
		uint64(0x000584), uint64(0x4),
		uint64(0x000034), uint64(0x000034),
		uint64(0x0000000000400584), uint64(0x0000000000400584))

	testSegHdr(
		t,
		progHdrTbl[7], 7,
		SegTypeGnuStack, SegFlagsReadable|SegFlagsWritable,
		uint64(0x000000), uint64(0x10),
		uint64(0x000000), uint64(0x000000),
		uint64(0x0000000000000000), uint64(0x0000000000000000))

	testSegHdr(
		t,
		progHdrTbl[8], 8,
		SegTypeGnuRelRO, SegFlagsReadable,
		uint64(0x000e10), uint64(0x1),
		uint64(0x0001f0), uint64(0x0001f0),
		uint64(0x0000000000600e10), uint64(0x0000000000600e10))
}
