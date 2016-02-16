// #############################################################################
// This file is part of the "garf" package of the "Eureka" project.
// It is distributed under the MIT License. Refer to the LICENSE file for more
// information.
//
// Website: http://www.github.com/sivachandra/eureka
// #############################################################################

// Package garf provides API to read DWARF debug info from ELF files.
package garf

import (
	"encoding/binary"
	"fmt"
)

import (
	"eureka/utils"
	"eureka/utils/leb128"
)

func (d *DwData) readLineNumberInfo(u *DwUnit) (*LnInfo, error) {
	uDie, err := u.DIETree()
	if err != nil {
		err = fmt.Errorf(
			"Cannot read line number info due to error fetching DIE for the unit.\n%s",
			err.Error())
		return nil, err
	}

	stmtListAt, exists := uDie.Attributes[DW_AT_stmt_list]
	if !exists {
		err = fmt.Errorf(
			"Cannot read line number info as unit does not have DW_AT_stmt_list attr.")
		return nil, err
	}

	var offset uint64
	switch stmtListAt.Value.(type) {
	case uint64:
		offset = stmtListAt.Value.(uint64)
	default:
		err = fmt.Errorf("Unknown value type of units DW_AT_stmt_list attr.")
		return nil, err
	}

	elf := d.ELFData()
	debugLineSect, exists := elf.SectMap()[".debug_line"]
	if !exists {
		err = fmt.Errorf("Cannot read line number info as .debug_line section is missing.")
		return nil, err
	}
	if len(debugLineSect) > 1 {
	}

	sectReader, err := debugLineSect[0].NewReader()
	if err != nil {
		err = fmt.Errorf(
			"Unable to get a SectReader for .debug_line section.\n%s", err.Error())
		return nil, err
	}

	_, err = sectReader.Seek(int64(offset), 0)
	if err != nil {
		err = fmt.Errorf(
			"Unable to seek to the offset for the unit in the .debug_line section.\n%s",
			err.Error())
		return nil, err
	}
	initLen := sectReader.Len()

	lnInfo := new(LnInfo)
	endianess := d.elf.Endianess()

	var len64 uint64
	var len32 uint32
	var lenSize = uint64(4)
	err = binary.Read(sectReader, endianess, &len32)
	if err != nil {
		return nil, fmt.Errorf("Error reading initial length of line info.")
	}
	if len32 == 0xFFFFFFFF {
		err = binary.Read(sectReader, endianess, &len64)
		if err != nil {
			return nil, fmt.Errorf("Error reading initial length of line info.")
		}
		lenSize += 8
	} else {
		len64 = uint64(len32)
	}
	lnInfo.Size = lenSize + len64

	err = binary.Read(sectReader, endianess, &lnInfo.Version)
	if err != nil {
		return nil, fmt.Errorf("Error reading version from line info header.")
	}

	if lnInfo.Version >= 5 {
		err = binary.Read(sectReader, endianess, &lnInfo.AddressSize)
		if err != nil {
			err = fmt.Errorf(
				"Error reading address size from line info header.\n%s",
				err.Error())
			return nil, err
		}
		err = binary.Read(sectReader, endianess, &lnInfo.SegmentSelectorSize)
		if err != nil {
			err = fmt.Errorf(
				"Error reading segment selector size from line info header.\n%s",
				err.Error())
			return nil, err
		}
	}

	// Skip over the header length field.
	if u.Format == DwFormat32 {
		err = binary.Read(sectReader, endianess, &len32)
	} else {
		err = binary.Read(sectReader, endianess, &len64)
	}
	if err != nil {
		err = fmt.Errorf(
			"Error skipping over header length field of line info header.\n%s",
			err.Error())
		return nil, err
	}

	err = binary.Read(sectReader, endianess, &lnInfo.minInstrLength)
	if err != nil {
		err = fmt.Errorf(
			"Error reading minimum instruction length from line info header.\n%s",
			err.Error())
		return nil, err
	}

	if lnInfo.Version >= 4 {
		err = binary.Read(sectReader, endianess, &lnInfo.maxOprPerInstr)
		if err != nil {
			err = fmt.Errorf(
				"Error reading max opers per instr from line info header.\n%s",
				err.Error())
			return nil, err
		}
	}

	err = binary.Read(sectReader, endianess, &lnInfo.defaultIsStmt)
	if err != nil {
		err = fmt.Errorf(
			"Error reading 'default_is_stmt' from line info header.\n%s",
			err.Error())
		return nil, err
	}
	err = binary.Read(sectReader, endianess, &lnInfo.lineBase)
	if err != nil {
		err = fmt.Errorf(
			"Error reading line base from line info header.\n%s",
			err.Error())
		return nil, err
	}
	err = binary.Read(sectReader, endianess, &lnInfo.lineRange)
	if err != nil {
		err = fmt.Errorf(
			"Error reading line range from line info header.\n%s",
			err.Error())
		return nil, err
	}
	err = binary.Read(sectReader, endianess, &lnInfo.opcodeBase)
	if err != nil {
		err = fmt.Errorf(
			"Error reading opcode base from line info header.\n%s",
			err.Error())
		return nil, err
	}

	for i := uint8(1); i < lnInfo.opcodeBase; i++ {
		c, err := sectReader.ReadByte()
		if err != nil {
			err = fmt.Errorf(
				"Error reading opcode count table from line info header.\n%s",
				err.Error())
			return nil, err
		}
		lnInfo.operandCountTbl = append(lnInfo.operandCountTbl, c)
	}

	if lnInfo.Version >= 5 {
		// TODO: Add support for DWARF 5 line number info.
		return lnInfo, nil
	}

	// Read directory entries
	for true {
		dir, err := utils.ReadCString(sectReader)
		if err != nil {
			err = fmt.Errorf(
				"Error reading directory entry from line info header.\n%s",
				err.Error())
			return nil, err
		}
		if len(dir) == 0 {
			break
		}

		lnInfo.Directories = append(lnInfo.Directories, dir)
	}

	// Read file entries
	for true {
		var fileEntry LnFileEntry

		fileEntry.Path, err = utils.ReadCString(sectReader)
		if err != nil {
			err = fmt.Errorf(
				"Error reading file name from line info header.\n%s",
				err.Error())
			return nil, err
		}
		if len(fileEntry.Path) == 0 {
			break
		}

		fileEntry.DirIndex, err = leb128.ReadUnsigned(sectReader)
		if err != nil {
			err = fmt.Errorf(
				"Error reading directory index of a file in line info header.\n%s",
				err.Error())
			return nil, err
		}

		fileEntry.Timestamp, err = leb128.ReadUnsigned(sectReader)
		if err != nil {
			err = fmt.Errorf(
				"Error reading time stamp of a file in line info header.\n%s",
				err.Error())
			return nil, err
		}

		fileEntry.Size, err = leb128.ReadUnsigned(sectReader)
		if err != nil {
			err = fmt.Errorf(
				"Error reading size of a file in line info header.\n%s",
				err.Error())
			return nil, err
		}

		lnInfo.Files = append(lnInfo.Files, fileEntry)
	}

	// Read the program until the end of the line info for the unit.
	for uint64(initLen-sectReader.Len()) < lnInfo.Size {
		b, err := sectReader.ReadByte()
		if err != nil {
			err = fmt.Errorf(
				"Error reading opcode of a line program instruction.\n%s",
				err.Error())
			return nil, err
		}

		var instr LnInstr
		if b == 0 {
			// Extension opcode
			// Read out the size of the instruction first
			_, err := leb128.Read(sectReader)
			if err != nil {
				err = fmt.Errorf(
					"Error reading extension opcode instruction size.\n%s",
					err.Error())
				return nil, err
			}

			b, err = sectReader.ReadByte()
			if err != nil {
				err = fmt.Errorf(
					"Error reading extension opcode from line program.\n%s",
					err.Error())
				return nil, err
			}

			instr.Opcode = DwLnOpcode(b)
			instr.OpcodeType = DwLnOpcodeExt

			switch DwLnOpcode(b) {
			case DW_LNE_end_sequence:
				break
			case DW_LNE_set_address:
				addrSize := d.elf.AddressSize()
				err = nil
				var addr uint64

				switch addrSize {
				case 1:
					var addr8 uint8
					addr8, err = sectReader.ReadByte()
					addr = uint64(addr8)
				case 2:
					var addr16 uint16
					err = binary.Read(sectReader, endianess, &addr16)
					addr = uint64(addr16)
				case 4:
					var addr32 uint32
					err = binary.Read(sectReader, endianess, &addr32)
					addr = uint64(addr32)
				case 8:
					err = binary.Read(sectReader, endianess, &addr)
				default:
					err = fmt.Errorf(
						"Unsupported address size in DW_LNE_set_address.")
				}

				if err != nil {
					err = fmt.Errorf(
						"Error reading operand of DW_LNE_set_address.\n%s",
						err.Error())
					return nil, err
				}

				operand, err := leb128.Encode(addr)
				if err != nil {
					err = fmt.Errorf(
						"Error encoding operand of DW_LNE_set_address.\n%s",
						err.Error())
					return nil, err
				}

				instr.Operands = append(instr.Operands, operand)
			case DW_LNE_define_file:
				err = fmt.Errorf(
					"Unsupported extended opcode in line number program.")
				return nil, err
			case DW_LNE_set_discriminator:
				operand, err := leb128.Read(sectReader)
				if err != nil {
					msg := "Error reading operand of DW_LNE_set_discriminator."
					err = fmt.Errorf("%s\n%s", msg, err.Error())
					return nil, err
				}

				instr.Operands = append(instr.Operands, operand)
			}
		} else if b < lnInfo.opcodeBase {
			// Standard opcode
			instr.Opcode = DwLnOpcode(b)
			instr.OpcodeType = DwLnOpcodeStd
			switch DwLnOpcode(b) {
			case DW_LNS_copy:
				fallthrough
			case DW_LNS_negate_stmt:
				fallthrough
			case DW_LNS_set_basic_block:
				fallthrough
			case DW_LNS_const_add_pc:
				fallthrough
			case DW_LNS_set_prologue_end:
				fallthrough
			case DW_LNS_set_epilogue_begin:
				break
			case DW_LNS_advance_pc:
				fallthrough
			case DW_LNS_advance_line:
				fallthrough
			case DW_LNS_set_file:
				fallthrough
			case DW_LNS_set_column:
				fallthrough
			case DW_LNS_set_isa:
				operand, err := leb128.Read(sectReader)
				if err != nil {
					msg := "Error reading operand of std line program opcode."
					err = fmt.Errorf("%s\n%s", msg, err.Error())
					return nil, err
				}

				instr.Operands = append(instr.Operands, operand)
			case DW_LNS_fixed_advance_pc:
				var operand16 uint16
				err = binary.Read(sectReader, endianess, &operand16)
				if err != nil {
					msg := "Error reading operand of DW_LNS_fixed_advance_pc."
					err = fmt.Errorf("%s\n%s", msg, err.Error())
					return nil, err
				}

				operand, err := leb128.Encode(operand16)
				if err != nil {
					msg := "Error encoding operand of DW_LNS_fixed_advance_pc."
					err = fmt.Errorf("%s\n%s", msg, err.Error())
					return nil, err
				}

				instr.Operands = append(instr.Operands, operand)
			}
		} else {
			// Special opcode
			instr.Opcode = DwLnOpcode(b)
			instr.OpcodeType = DwLnOpcodeSpecial
		}

		lnInfo.Program = append(lnInfo.Program, instr)
	}

	return lnInfo, nil
}
