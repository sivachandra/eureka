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
	"bytes"
	"encoding/binary"
	"fmt"
)

import (
	"eureka/guts/leb128"
)

func operandReadError(op DwOp, i uint8, e error) error {
	e = fmt.Errorf(
		"Error reading operand %d to %s in DWARF expression.\n%s",
		i, DwOpStr[op], e.Error())
	return e
}

func (d *DwData) readSizeAndDwExpr(
	u *DwUnit, r *bytes.Reader, en binary.ByteOrder) (DwExpr, error) {
	l, err := leb128.ReadUnsigned(r)
	if err != nil {
		return nil, fmt.Errorf("Error reading length of exprloc data.\n%s", err.Error())
	}

	return d.readDwExpr(u, r, en, l)
}

func (d *DwData) readDwExpr(
	u *DwUnit, r *bytes.Reader, en binary.ByteOrder, l uint64) (DwExpr, error) {
	var expr DwExpr
	rem := r.Len()
	for uint64(rem-r.Len()) < l {
		b, err := r.ReadByte()
		if err != nil {
			err = fmt.Errorf(
				"Error reading opcode of operation %d in exprloc data.\n%s",
				len(expr),
				err.Error())
			return nil, err
		}

		op := DwOp(b)

		var operation DwOperation
		operation.Op = op
		operation.Operands = make([]interface{}, 0)

		switch op {
		case DW_OP_addr:
			switch d.elf.AddressSize() {
			case 4:
				var addr uint32
				err = binary.Read(r, en, &addr)
				if err != nil {
					return nil, operandReadError(op, 0, err)
				}
				operation.Operands = append(operation.Operands, uint64(addr))
			case 8:
				var addr uint64
				err = binary.Read(r, en, &addr)
				if err != nil {
					return nil, operandReadError(op, 0, err)
				}
				operation.Operands = append(operation.Operands, addr)
			default:
				err = fmt.Errorf(
					"Unknown target address size to read DW_OP_addr operand.")
				return nil, err
			}
		case DW_OP_deref:
		case DW_OP_const1u:
			var c uint8
			err = binary.Read(r, en, &c)
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			operation.Operands = append(operation.Operands, c)
		case DW_OP_const1s:
			var c int8
			err = binary.Read(r, en, &c)
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			operation.Operands = append(operation.Operands, c)
		case DW_OP_const2u:
			var c uint16
			err = binary.Read(r, en, &c)
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			operation.Operands = append(operation.Operands, c)
		case DW_OP_const2s:
			var c int16
			err = binary.Read(r, en, &c)
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			operation.Operands = append(operation.Operands, c)
		case DW_OP_const4u:
			var c uint32
			err = binary.Read(r, en, &c)
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			operation.Operands = append(operation.Operands, c)
		case DW_OP_const4s:
			var c int32
			err = binary.Read(r, en, &c)
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			operation.Operands = append(operation.Operands, c)
		case DW_OP_const8u:
			var c uint64
			err = binary.Read(r, en, &c)
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			operation.Operands = append(operation.Operands, c)
		case DW_OP_const8s:
			var c int64
			err = binary.Read(r, en, &c)
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			operation.Operands = append(operation.Operands, c)
		case DW_OP_constu:
			var c uint64
			c, err = leb128.ReadUnsigned(r)
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			operation.Operands = append(operation.Operands, c)
		case DW_OP_consts:
			var c int64
			c, err = leb128.ReadSigned(r)
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			operation.Operands = append(operation.Operands, c)
		case DW_OP_dup:
		case DW_OP_drop:
		case DW_OP_over:
		case DW_OP_pick:
			p, err := r.ReadByte()
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			operation.Operands = append(operation.Operands, p)
		case DW_OP_swap:
		case DW_OP_rot:
		case DW_OP_xderef:
		case DW_OP_abs:
		case DW_OP_and:
		case DW_OP_div:
		case DW_OP_minus:
		case DW_OP_mod:
		case DW_OP_mul:
		case DW_OP_neg:
		case DW_OP_not:
		case DW_OP_or:
		case DW_OP_plus:
		case DW_OP_plus_uconst:
			var c uint64
			c, err = leb128.ReadUnsigned(r)
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			operation.Operands = append(operation.Operands, c)
		case DW_OP_shl:
		case DW_OP_shr:
		case DW_OP_shra:
		case DW_OP_xor:
		case DW_OP_bra:
			var c int16
			err = binary.Read(r, en, &c)
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			operation.Operands = append(operation.Operands, c)
		case DW_OP_eq:
		case DW_OP_ge:
		case DW_OP_gt:
		case DW_OP_le:
		case DW_OP_lt:
		case DW_OP_ne:
		case DW_OP_skip:
			var c int16
			err = binary.Read(r, en, &c)
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			operation.Operands = append(operation.Operands, c)
		case DW_OP_lit0:
		case DW_OP_lit1:
		case DW_OP_lit2:
		case DW_OP_lit3:
		case DW_OP_lit4:
		case DW_OP_lit5:
		case DW_OP_lit6:
		case DW_OP_lit7:
		case DW_OP_lit8:
		case DW_OP_lit9:
		case DW_OP_lit10:
		case DW_OP_lit11:
		case DW_OP_lit12:
		case DW_OP_lit13:
		case DW_OP_lit14:
		case DW_OP_lit15:
		case DW_OP_lit16:
		case DW_OP_lit17:
		case DW_OP_lit18:
		case DW_OP_lit19:
		case DW_OP_lit20:
		case DW_OP_lit21:
		case DW_OP_lit22:
		case DW_OP_lit23:
		case DW_OP_lit24:
		case DW_OP_lit25:
		case DW_OP_lit26:
		case DW_OP_lit27:
		case DW_OP_lit28:
		case DW_OP_lit29:
		case DW_OP_lit30:
		case DW_OP_lit31:
		case DW_OP_reg0:
		case DW_OP_reg1:
		case DW_OP_reg2:
		case DW_OP_reg3:
		case DW_OP_reg4:
		case DW_OP_reg5:
		case DW_OP_reg6:
		case DW_OP_reg7:
		case DW_OP_reg8:
		case DW_OP_reg9:
		case DW_OP_reg10:
		case DW_OP_reg11:
		case DW_OP_reg12:
		case DW_OP_reg13:
		case DW_OP_reg14:
		case DW_OP_reg15:
		case DW_OP_reg16:
		case DW_OP_reg17:
		case DW_OP_reg18:
		case DW_OP_reg19:
		case DW_OP_reg20:
		case DW_OP_reg21:
		case DW_OP_reg22:
		case DW_OP_reg23:
		case DW_OP_reg24:
		case DW_OP_reg25:
		case DW_OP_reg26:
		case DW_OP_reg27:
		case DW_OP_reg28:
		case DW_OP_reg29:
		case DW_OP_reg30:
		case DW_OP_reg31:
		case DW_OP_breg0:
			fallthrough
		case DW_OP_breg1:
			fallthrough
		case DW_OP_breg2:
			fallthrough
		case DW_OP_breg3:
			fallthrough
		case DW_OP_breg4:
			fallthrough
		case DW_OP_breg5:
			fallthrough
		case DW_OP_breg6:
			fallthrough
		case DW_OP_breg7:
			fallthrough
		case DW_OP_breg8:
			fallthrough
		case DW_OP_breg9:
			fallthrough
		case DW_OP_breg10:
			fallthrough
		case DW_OP_breg11:
			fallthrough
		case DW_OP_breg12:
			fallthrough
		case DW_OP_breg13:
			fallthrough
		case DW_OP_breg14:
			fallthrough
		case DW_OP_breg15:
			fallthrough
		case DW_OP_breg16:
			fallthrough
		case DW_OP_breg17:
			fallthrough
		case DW_OP_breg18:
			fallthrough
		case DW_OP_breg19:
			fallthrough
		case DW_OP_breg20:
			fallthrough
		case DW_OP_breg21:
			fallthrough
		case DW_OP_breg22:
			fallthrough
		case DW_OP_breg23:
			fallthrough
		case DW_OP_breg24:
			fallthrough
		case DW_OP_breg25:
			fallthrough
		case DW_OP_breg26:
			fallthrough
		case DW_OP_breg27:
			fallthrough
		case DW_OP_breg28:
			fallthrough
		case DW_OP_breg29:
			fallthrough
		case DW_OP_breg30:
			fallthrough
		case DW_OP_breg31:
			var c int64
			c, err = leb128.ReadSigned(r)
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			operation.Operands = append(operation.Operands, c)
		case DW_OP_regx:
			var c uint64
			c, err = leb128.ReadUnsigned(r)
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			operation.Operands = append(operation.Operands, c)
		case DW_OP_fbreg:
			var c int64
			c, err = leb128.ReadSigned(r)
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			operation.Operands = append(operation.Operands, c)
		case DW_OP_bregx:
			reg, err := leb128.ReadUnsigned(r)
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			offset, err := leb128.ReadSigned(r)
			if err != nil {
				return nil, operandReadError(op, 1, err)
			}
			operation.Operands = append(operation.Operands, reg)
			operation.Operands = append(operation.Operands, offset)
		case DW_OP_piece:
			var c uint64
			c, err = leb128.ReadUnsigned(r)
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			operation.Operands = append(operation.Operands, c)
		case DW_OP_deref_size:
			s, err := r.ReadByte()
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			operation.Operands = append(operation.Operands, s)
		case DW_OP_xderef_size:
			s, err := r.ReadByte()
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			operation.Operands = append(operation.Operands, s)
		case DW_OP_nop:
		case DW_OP_push_object_address:
		case DW_OP_call2:
			var c uint16
			err = binary.Read(r, en, &c)
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			operation.Operands = append(operation.Operands, c)
		case DW_OP_call4:
			var c uint32
			err = binary.Read(r, en, &c)
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			operation.Operands = append(operation.Operands, c)
		case DW_OP_call_ref:
			if u.Format == DwFormat32 {
				var c uint32
				err = binary.Read(r, en, &c)
				if err != nil {
					return nil, operandReadError(op, 0, err)
				}
				operation.Operands = append(operation.Operands, c)
			} else {
				var c uint64
				err = binary.Read(r, en, &c)
				if err != nil {
					return nil, operandReadError(op, 0, err)
				}
				operation.Operands = append(operation.Operands, c)
			}
		case DW_OP_form_tls_address:
		case DW_OP_call_frame_cfa:
		case DW_OP_bit_piece:
			reg, err := leb128.ReadUnsigned(r)
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			offset, err := leb128.ReadUnsigned(r)
			if err != nil {
				return nil, operandReadError(op, 1, err)
			}
			operation.Operands = append(operation.Operands, reg)
			operation.Operands = append(operation.Operands, offset)
		case DW_OP_implicit_value:
			size, err := leb128.ReadUnsigned(r)
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			value := make([]byte, size)
			_, err = r.Read(value)
			if err != nil {
				return nil, operandReadError(op, 1, err)
			}
			operation.Operands = append(operation.Operands, size)
			operation.Operands = append(operation.Operands, value)
		case DW_OP_stack_value:
		case DW_OP_implicit_pointer, DW_OP_GNU_implicit_pointer:
			var offset uint64
			if u.Format == DwFormat32 {
				var c uint32
				err = binary.Read(r, en, &c)
				if err != nil {
					return nil, operandReadError(op, 0, err)
				}
				offset = uint64(c)
			} else {
				err = binary.Read(r, en, &offset)
				if err != nil {
					return nil, operandReadError(op, 0, err)
				}
			}
			constant, err := leb128.ReadUnsigned(r)
			if err != nil {
				return nil, operandReadError(op, 1, err)
			}
			operation.Operands = append(operation.Operands, offset)
			operation.Operands = append(operation.Operands, constant)
		case DW_OP_addrx:
			var c uint64
			c, err = leb128.ReadUnsigned(r)
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			operation.Operands = append(operation.Operands, c)
		case DW_OP_constx:
			var c uint64
			c, err = leb128.ReadUnsigned(r)
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			operation.Operands = append(operation.Operands, c)
		case DW_OP_entry_value, DW_OP_GNU_entry_value:
			size, err := leb128.ReadUnsigned(r)
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			value := make([]byte, size)
			_, err = r.Read(value)
			if err != nil {
				return nil, operandReadError(op, 1, err)
			}
			operation.Operands = append(operation.Operands, size)
			operation.Operands = append(operation.Operands, value)
		case DW_OP_const_type, DW_OP_GNU_const_type:
			return nil, fmt.Errorf("Unsupport opcode %s in DWARF expr.", DwOpStr[op])
		case DW_OP_regval_type, DW_OP_GNU_regval_type:
			reg, err := leb128.ReadUnsigned(r)
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			offset, err := leb128.ReadUnsigned(r)
			if err != nil {
				return nil, operandReadError(op, 1, err)
			}
			operation.Operands = append(operation.Operands, reg)
			operation.Operands = append(operation.Operands, offset)
		case DW_OP_deref_type, DW_OP_GNU_deref_type:
			size, err := r.ReadByte()
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			offset, err := leb128.ReadUnsigned(r)
			if err != nil {
				return nil, operandReadError(op, 1, err)
			}
			operation.Operands = append(operation.Operands, size)
			operation.Operands = append(operation.Operands, offset)
		case DW_OP_xderef_type:
			size, err := r.ReadByte()
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			offset, err := leb128.ReadUnsigned(r)
			if err != nil {
				return nil, operandReadError(op, 1, err)
			}
			operation.Operands = append(operation.Operands, size)
			operation.Operands = append(operation.Operands, offset)
		case DW_OP_convert, DW_OP_GNU_convert:
			var c uint64
			c, err = leb128.ReadUnsigned(r)
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			operation.Operands = append(operation.Operands, c)
		case DW_OP_reinterpret, DW_OP_GNU_reinterpret:
			var c uint64
			c, err = leb128.ReadUnsigned(r)
			if err != nil {
				return nil, operandReadError(op, 0, err)
			}
			operation.Operands = append(operation.Operands, c)
		case DW_OP_GNU_push_tls_address:
			return nil, fmt.Errorf("Unsupport opcode %s in DWARF expr.", DwOpStr[op])
		case DW_OP_GNU_uninit:
			return nil, fmt.Errorf("Unsupport opcode %s in DWARF expr.", DwOpStr[op])
		case DW_OP_GNU_encoded_addr:
			return nil, fmt.Errorf("Unsupport opcode %s in DWARF expr.", DwOpStr[op])
		case DW_OP_GNU_parameter_ref:
			return nil, fmt.Errorf("Unsupport opcode %s in DWARF expr.", DwOpStr[op])
		default:
			err = fmt.Errorf("Unknown opcode %d while reading DWARF expression.", op)
			return nil, err
		}

		expr = append(expr, operation)
	}

	return expr, nil
}
