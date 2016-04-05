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
	"eureka/guts/ruts"
)

func (d *DwData) readAttr(
	u *DwUnit,
	r *bytes.Reader,
	at DwAt,
	form DwForm,
	en binary.ByteOrder) (Attribute, error) {
	var attr Attribute
	var err error

	attr.Name = at

	switch at {
	case DW_AT_sibling:
		attr.Value, err = d.readAttrRef(u, r, form, en)
	case DW_AT_location:
		if form.IsLocListPtr() {
			attr.Value, err = d.readAttrLocList(u, r, form, en)
		} else if form.IsExprLoc() {
			attr.Value, err = d.readSizeAndDwExpr(u, r, en)
		} else {
			err = fmt.Errorf("Unsupported form %s for DW_AT_Location.", DwFormStr[form])
		}
	case DW_AT_name:
		attr.Value, err = d.readAttrStr(u, r, form, en)
	case DW_AT_ordering:
		var v uint8
		v, err = r.ReadByte()
		attr.Value = DwOrder(v)
	case DW_AT_byte_size:
		fallthrough
	case DW_AT_bit_offset:
		fallthrough
	case DW_AT_bit_size:
		if form.IsConstant() {
			attr.Value, err = d.readAttrUint32(u, r, form, en)
		} else if form.IsExprLoc() {
			attr.Value, err = d.readSizeAndDwExpr(u, r, en)
		} else if form.IsRef() {
			attr.Value, err = d.readAttrRef(u, r, form, en)
		} else {
			err = fmt.Errorf(
				"Unsupported form %s for attribute %s.",
				DwFormStr[form],
				DwAtStr[at])
		}
	case DW_AT_stmt_list:
		if !form.IsLinePtr() {
			err = fmt.Errorf(
				"Unsupported form %s for DW_AT_stmt_list.", DwFormStr[form])
			break
		}
		attr.Value, err = d.readAttrUint64(u, r, form, en)
	case DW_AT_low_pc:
		if !form.IsAddress() {
			err = fmt.Errorf(
				"Unsupported form %s for DW_AT_low_pc.", DwFormStr[form])
			break
		}
		attr.Value, err = d.readAttrUint64(u, r, form, en)
	case DW_AT_high_pc:
		if !form.IsAddress() && !form.IsConstant() {
			err = fmt.Errorf(
				"Unsupported form %s for DW_AT_high_pc.", DwFormStr[form])
			break
		}
		attr.Value, err = d.readAttrUint64(u, r, form, en)
	case DW_AT_language:
		if !form.IsConstant() {
			err = fmt.Errorf(
				"Unsupported form %s for DW_AT_language.", DwFormStr[form])
			break
		}
		attr.Value, err = d.readAttrUint16(u, r, form, en)
		attr.Value = DwLang(attr.Value.(uint16))
	case DW_AT_visibility:
		var v uint8
		v, err = r.ReadByte()
		attr.Value = DwVis(v)
	case DW_AT_import:
		attr.Value, err = d.readAttrRef(u, r, form, en)
	case DW_AT_string_length:
		if form.IsExprLoc() {
			attr.Value, err = d.readSizeAndDwExpr(u, r, en)
		} else if form.IsLocListPtr() {
			attr.Value, err = d.readAttrLocList(u, r, form, en)
		} else {
			err = fmt.Errorf(
				"Unsupported form %s for DW_AT_string_length.", DwFormStr[form])
		}
	case DW_AT_common_reference:
		attr.Value, err = d.readAttrRef(u, r, form, en)
	case DW_AT_comp_dir:
		attr.Value, err = d.readAttrStr(u, r, form, en)
	case DW_AT_const_value:
		if form.IsBlock() {
			attr.Value, err = d.readAttrByteSlice(u, r, form, en)
		} else if form.IsConstant() {
			attr.Value, err = d.readAttrInt64(u, r, form, en)
		} else if form.IsString() {
			attr.Value, err = d.readAttrStr(u, r, form, en)
		} else {
			err = fmt.Errorf(
				"Unsupported form %s for DW_AT_const_value.", DwFormStr[form])
		}
	case DW_AT_containing_type:
		attr.Value, err = d.readAttrRef(u, r, form, en)
	case DW_AT_default_value:
		attr.Value, err = d.readAttrRef(u, r, form, en)
	case DW_AT_inline:
		var v uint8
		v, err = r.ReadByte()
		attr.Value = DwInl(v)
	case DW_AT_is_optional:
		attr.Value, err = d.readAttrFlag(u, r, form, en)
	case DW_AT_lower_bound:
		if form.IsExprLoc() {
			attr.Value, err = d.readSizeAndDwExpr(u, r, en)
		} else if form.IsConstant() {
			attr.Value, err = d.readAttrInt64(u, r, form, en)
		} else if form.IsRef() {
			attr.Value, err = d.readAttrRef(u, r, form, en)
		} else {
			err = fmt.Errorf(
				"Unsupported form %s for DW_AT_lower_bound.", DwFormStr[form])
		}
	case DW_AT_producer:
		attr.Value, err = d.readAttrStr(u, r, form, en)
	case DW_AT_prototyped:
		attr.Value, err = d.readAttrFlag(u, r, form, en)
	case DW_AT_return_addr:
		if form.IsExprLoc() {
			attr.Value, err = d.readSizeAndDwExpr(u, r, en)
		} else if form.IsLocListPtr() {
			attr.Value, err = d.readAttrLocList(u, r, form, en)
		} else {
			err = fmt.Errorf(
				"Unsupported form %s for DW_AT_return_addr.", DwFormStr[form])
		}
	case DW_AT_start_scope:
		if form.IsRangeListPtr() {
			attr.Value, err = d.readAttrRangeList(u, r, form, en)
		} else if form.IsConstant() {
			attr.Value, err = d.readAttrInt64(u, r, form, en)
		} else {
			err = fmt.Errorf(
				"Unsupported form %s for DW_AT_start_scope.", DwFormStr[form])
		}
	case DW_AT_bit_stride:
		if form.IsExprLoc() {
			attr.Value, err = d.readSizeAndDwExpr(u, r, en)
		} else if form.IsConstant() {
			attr.Value, err = d.readAttrInt64(u, r, form, en)
		} else if form.IsRef() {
			attr.Value, err = d.readAttrRef(u, r, form, en)
		} else {
			err = fmt.Errorf(
				"Unsupported form %s for DW_AT_bit_stride.", DwFormStr[form])
		}
	case DW_AT_upper_bound:
		if form.IsExprLoc() {
			attr.Value, err = d.readSizeAndDwExpr(u, r, en)
		} else if form.IsConstant() {
			attr.Value, err = d.readAttrInt64(u, r, form, en)
		} else if form.IsRef() {
			attr.Value, err = d.readAttrRef(u, r, form, en)
		} else {
			err = fmt.Errorf(
				"Unsupported form %s for DW_AT_upper_bound.", DwFormStr[form])
		}
	case DW_AT_abstract_origin:
		attr.Value, err = d.readAttrRef(u, r, form, en)
	case DW_AT_accessibility:
		var v uint8
		v, err = r.ReadByte()
		attr.Value = DwAccess(v)
	case DW_AT_address_class:
		// TODO: Is reading a byte enough??
		attr.Value, err = r.ReadByte()
	case DW_AT_artificial:
		attr.Value, err = d.readAttrFlag(u, r, form, en)
	case DW_AT_base_types:
		attr.Value, err = d.readAttrRef(u, r, form, en)
	case DW_AT_data_member_location:
		if form == DW_FORM_sec_offset {
			attr.Value, err = d.readAttrUint64(u, r, form, en)
		} else if form == DW_FORM_exprloc {
			attr.Value, err = d.readAttrByteSlice(u, r, form, en)
		} else {
			attr.Value, err = d.readAttrInt64(u, r, form, en)
		}
	case DW_AT_decl_file:
		attr.Value, err = d.readAttrUint32(u, r, form, en)
	case DW_AT_decl_line:
		attr.Value, err = d.readAttrUint32(u, r, form, en)
	case DW_AT_declaration:
		attr.Value, err = d.readAttrFlag(u, r, form, en)
	case DW_AT_encoding:
		attr.Value, err = r.ReadByte()
		attr.Value = DwAte(attr.Value.(byte))
	case DW_AT_external:
		attr.Value, err = d.readAttrFlag(u, r, form, en)
	case DW_AT_frame_base:
		if form.IsExprLoc() {
			attr.Value, err = d.readSizeAndDwExpr(u, r, en)
		} else if form.IsLocListPtr() {
			attr.Value, err = d.readAttrLocList(u, r, form, en)
		} else {
			err = fmt.Errorf(
				"Unsupported form %s for DW_AT_frame_base.", DwFormStr[form])
		}
	case DW_AT_friend:
		attr.Value, err = d.readAttrRef(u, r, form, en)
	case DW_AT_namelist_item:
		attr.Value, err = d.readAttrRef(u, r, form, en)
	case DW_AT_priority:
		attr.Value, err = d.readAttrRef(u, r, form, en)
	case DW_AT_segment:
		if form.IsExprLoc() {
			attr.Value, err = d.readSizeAndDwExpr(u, r, en)
		} else if form.IsLocListPtr() {
			attr.Value, err = d.readAttrLocList(u, r, form, en)
		} else {
			err = fmt.Errorf(
				"Unsupported form %s for DW_AT_segment.", DwFormStr[form])
		}
	case DW_AT_specification:
		attr.Value, err = d.readAttrRef(u, r, form, en)
	case DW_AT_static_link:
		if form.IsExprLoc() {
			attr.Value, err = d.readSizeAndDwExpr(u, r, en)
		} else if form.IsLocListPtr() {
			attr.Value, err = d.readAttrLocList(u, r, form, en)
		} else {
			err = fmt.Errorf(
				"Unsupported form %s for DW_AT_static_link.", DwFormStr[form])
		}
	case DW_AT_type:
		attr.Value, err = d.readAttrRef(u, r, form, en)
	case DW_AT_variable_parameter:
		attr.Value, err = d.readAttrFlag(u, r, form, en)
	case DW_AT_virtuality:
		var v uint8
		v, err = r.ReadByte()
		attr.Value = DwVirtuality(v)
	case DW_AT_vtable_elem_location:
		if form.IsExprLoc() {
			attr.Value, err = d.readSizeAndDwExpr(u, r, en)
		} else if form.IsLocListPtr() {
			attr.Value, err = d.readAttrLocList(u, r, form, en)
		} else {
			err = fmt.Errorf(
				"Unsupported form %s for DW_AT_vtable_elem_location.",
				DwFormStr[form])
		}
	case DW_AT_ranges:
		attr.Value, err = d.readAttrRangeList(u, r, form, en)
	case DW_AT_picture_string:
		attr.Value, err = d.readAttrStr(u, r, form, en)
	case DW_AT_mutable:
		attr.Value, err = d.readAttrFlag(u, r, form, en)
	case DW_AT_threads_scaled:
		attr.Value, err = d.readAttrFlag(u, r, form, en)
	case DW_AT_explicit:
		attr.Value, err = d.readAttrFlag(u, r, form, en)
	case DW_AT_object_pointer:
		attr.Value, err = d.readAttrRef(u, r, form, en)
	case DW_AT_endianity:
		var v uint8
		v, err = r.ReadByte()
		attr.Value = DwEnd(v)
	case DW_AT_elemental:
		attr.Value, err = d.readAttrFlag(u, r, form, en)
	case DW_AT_pure:
		attr.Value, err = d.readAttrFlag(u, r, form, en)
	case DW_AT_recursive:
		attr.Value, err = d.readAttrFlag(u, r, form, en)
	case DW_AT_signature:
		attr.Value, err = d.readAttrRef(u, r, form, en)
	case DW_AT_main_subprogram:
		attr.Value, err = d.readAttrFlag(u, r, form, en)
	case DW_AT_const_expr:
		attr.Value, err = d.readAttrFlag(u, r, form, en)
	case DW_AT_enum_class:
		attr.Value, err = d.readAttrFlag(u, r, form, en)
	case DW_AT_linkage_name:
		attr.Value, err = d.readAttrStr(u, r, form, en)

	// GNU extension attributes
	case DW_AT_GNU_tail_call:
		attr.Value, err = d.readAttrFlag(u, r, form, en)
	case DW_AT_GNU_all_tail_call_sites:
		attr.Value, err = d.readAttrFlag(u, r, form, en)
	case DW_AT_GNU_all_call_sites:
		attr.Value, err = d.readAttrFlag(u, r, form, en)
	case DW_AT_GNU_call_site_value:
		if form.IsExprLoc() {
			attr.Value, err = d.readSizeAndDwExpr(u, r, en)
		} else {
			err = fmt.Errorf(
				"Unsupported form %s for %s.",
				DwFormStr[form], DwAtStr[DW_AT_GNU_call_site_value])
		}
	default:
		attr.Value, err = d.readAttrByteSlice(u, r, form, en)
	}

	return attr, err
}

func (d *DwData) readAttrStr(
	u *DwUnit, r *bytes.Reader, form DwForm, en binary.ByteOrder) (string, error) {
	switch form {
	case DW_FORM_string:
		str, err := ruts.ReadCString(r)
		if err != nil {
			err = fmt.Errorf("Error reading inline string attribute value.", err)
			return "", err
		}
		return string(str), nil
	case DW_FORM_strp:
		var offset uint64
		if u.Format == DwFormat32 {
			var offset32 uint32

			err := binary.Read(r, en, &offset32)
			if err != nil {
				err = fmt.Errorf("Error reading .debug_str 32-bit offset.", err)
				return "", err
			}

			offset = uint64(offset32)
		} else {
			err := binary.Read(r, en, &offset)
			if err != nil {
				err = fmt.Errorf("Error reading .debug_str 64-bit offset.", err)
				return "", err
			}
		}

		debugStrTbl, err := d.DebugStr()
		if err != nil {
			return "", fmt.Errorf("Error reading .debug_str.", err)
		}

		str, err := debugStrTbl.ReadStr(offset)
		if err != nil {
			return "", fmt.Errorf("Error reading string: %s.", err.Error)
		}

		return str, nil
	default:
		err := fmt.Errorf(
			fmt.Sprintf("Cannot read data of form %d as string data.", form), nil)
		return "", err
	}
}

func (d *DwData) readAttrInt16(
	u *DwUnit, r *bytes.Reader, f DwForm, en binary.ByteOrder) (int16, error) {
	var err error

	switch f {
	case DW_FORM_data1:
		var i int8

		err = binary.Read(r, en, &i)
		if err != nil {
			break
		}

		return int16(i), nil
	case DW_FORM_data2:
		var i int16

		err = binary.Read(r, en, &i)
		if err != nil {
			break
		}

		return i, nil
	case DW_FORM_sdata:
		var i int64
		i, err = leb128.ReadSigned(r)
		if err != nil {
			break
		}

		return int16(i), nil
	default:
		err = fmt.Errorf(
			"Cannot read data of form %s as int16 value", DwFormStr[f])
		return 0, err
	}

	err = fmt.Errorf("Error reading data of form %s.\n%s", DwFormStr[f], err.Error())
	return 0, err
}

func (d *DwData) readAttrInt32(
	u *DwUnit, r *bytes.Reader, f DwForm, en binary.ByteOrder) (int32, error) {
	var err error

	switch f {
	case DW_FORM_data1:
		var i int8

		err = binary.Read(r, en, &i)
		if err != nil {
			break
		}

		return int32(i), nil
	case DW_FORM_data2:
		var i int16

		err = binary.Read(r, en, &i)
		if err != nil {
			break
		}

		return int32(i), nil
	case DW_FORM_data4:
		var i int32

		err = binary.Read(r, en, &i)
		if err != nil {
			break
		}

		return i, nil
	case DW_FORM_sdata:
		var i int64
		i, err = leb128.ReadSigned(r)
		if err != nil {
			break
		}

		return int32(i), nil
	default:
		return 0, fmt.Errorf("Cannot read data of form %s as int32.", DwFormStr[f])
	}

	err = fmt.Errorf("Error reading data of form %s.\n%s", DwFormStr[f], err.Error())
	return 0, err
}

func (d *DwData) readAttrInt64(
	u *DwUnit, r *bytes.Reader, f DwForm, en binary.ByteOrder) (int64, error) {
	var err error

	switch f {
	case DW_FORM_data1:
		var i int8

		err = binary.Read(r, en, &i)
		if err != nil {
			break
		}

		return int64(i), nil
	case DW_FORM_data2:
		var i int16

		err = binary.Read(r, en, &i)
		if err != nil {
			break
		}

		return int64(i), nil
	case DW_FORM_data4:
		var i int32

		err = binary.Read(r, en, &i)
		if err != nil {
			break
		}

		return int64(i), nil
	case DW_FORM_data8:
		var i int64

		err = binary.Read(r, en, &i)
		if err != nil {
			break
		}

		return i, nil
	case DW_FORM_sdata:
		var i int64
		i, err = leb128.ReadSigned(r)
		if err != nil {
			break
		}

		return i, nil
	default:
		return 0, fmt.Errorf("Cannot read data of form %s as int64.", DwFormStr[f])
	}

	err = fmt.Errorf("Error reading data of form %s.\n%s", DwFormStr[f], err.Error())
	return 0, err
}

func (d *DwData) readAttrUint16(
	u *DwUnit, r *bytes.Reader, f DwForm, en binary.ByteOrder) (uint16, error) {
	var err error

	switch f {
	case DW_FORM_data1:
		var i uint8

		err = binary.Read(r, en, &i)
		if err != nil {
			break
		}

		return uint16(i), nil
	case DW_FORM_data2:
		var i uint16

		err = binary.Read(r, en, &i)
		if err != nil {
			break
		}

		return i, nil
	case DW_FORM_udata:
		var i uint64
		i, err = leb128.ReadUnsigned(r)
		if err != nil {
			break
		}

		return uint16(i), nil
	default:
		return 0, fmt.Errorf("Cannot read data of form %s as int16.", DwFormStr[f])
	}

	err = fmt.Errorf("Error reading data of form %s.\n%s", DwFormStr[f], err.Error())
	return 0, err
}

func (d *DwData) readAttrUint32(
	u *DwUnit, r *bytes.Reader, f DwForm, en binary.ByteOrder) (uint32, error) {
	var err error

	switch f {
	case DW_FORM_data1:
		var i uint8

		err = binary.Read(r, en, &i)
		if err != nil {
			break
		}

		return uint32(i), nil
	case DW_FORM_data2:
		var i uint16

		err = binary.Read(r, en, &i)
		if err != nil {
			break
		}

		return uint32(i), nil
	case DW_FORM_data4:
		var i uint32

		err = binary.Read(r, en, &i)
		if err != nil {
			break
		}

		return i, nil
	case DW_FORM_udata:
		var i uint64
		i, err = leb128.ReadUnsigned(r)
		if err != nil {
			break
		}

		return uint32(i), nil
	default:
		return 0, fmt.Errorf("Cannot read data of form %s as uint32.", DwFormStr[f])
	}

	err = fmt.Errorf("Error reading data of form %s.\n%s", DwFormStr[f], err.Error())
	return 0, err
}

func (d *DwData) readAttrUint64(
	u *DwUnit, r *bytes.Reader, f DwForm, en binary.ByteOrder) (uint64, error) {
	var err error

	switch f {
	case DW_FORM_data1:
		var i uint8

		err = binary.Read(r, en, &i)
		if err != nil {
			break
		}

		return uint64(i), nil
	case DW_FORM_data2:
		var i uint16

		err = binary.Read(r, en, &i)
		if err != nil {
			break
		}

		return uint64(i), nil
	case DW_FORM_data4:
		var i uint32

		err = binary.Read(r, en, &i)
		if err != nil {
			break
		}

		return uint64(i), nil
	case DW_FORM_data8:
		var i uint64

		err = binary.Read(r, en, &i)
		if err != nil {
			break
		}

		return i, nil
	case DW_FORM_udata:
		var i uint64
		i, err = leb128.ReadUnsigned(r)
		if err != nil {
			break
		}

		return i, nil
	case DW_FORM_addr:
		if u.AddressSize == 4 {
			var i uint32
			err = binary.Read(r, en, &i)
			if err != nil {
				break
			}
			return uint64(i), nil
		} else {
			var i uint64
			err = binary.Read(r, en, &i)
			if err != nil {
				break
			}
			return i, nil
		}
	case DW_FORM_sec_offset:
		if u.Format == DwFormat32 {
			var i uint32

			err := binary.Read(r, en, &i)
			if err != nil {
				break
			}

			return uint64(i), nil
		} else {
			var i uint64

			err := binary.Read(r, en, &i)
			if err != nil {
				break
			}

			return i, nil
		}
	default:
		return 0, fmt.Errorf("Cannot read data of form %s as uint64.", DwFormStr[f])
	}

	err = fmt.Errorf("Error reading data of form %s.\n%s", DwFormStr[f], err.Error())
	return 0, err
}

func (d *DwData) readAttrFlag(
	u *DwUnit, r *bytes.Reader, f DwForm, en binary.ByteOrder) (bool, error) {
	switch f {
	case DW_FORM_flag:
		b, err := r.ReadByte()

		if err != nil {
			err = fmt.Errorf("Erroring reading DW_FORM_flag value.\n%s", err.Error())
			return false, err
		}

		return b != 0, nil
	case DW_FORM_flag_present:
		return true, nil
	default:
		return false, fmt.Errorf("Cannot read data of form %s as a flag.", DwFormStr[f])
	}
}

func (d *DwData) readAttrRef(
	u *DwUnit, r *bytes.Reader, f DwForm, en binary.ByteOrder) (*DIE, error) {
	var offset uint64
	var err error

	switch f {
	case DW_FORM_ref1:
		var b byte

		b, err = r.ReadByte()
		if err != nil {
			break
		}

		offset = uint64(b)
	case DW_FORM_ref2:
		var i uint16

		err = binary.Read(r, en, &i)
		if err != nil {
			break
		}

		offset = uint64(i)
	case DW_FORM_ref4:
		var i uint32

		err = binary.Read(r, en, &i)
		if err != nil {
			break
		}

		offset = uint64(i)
	case DW_FORM_ref8:
		err = binary.Read(r, en, &offset)
	case DW_FORM_ref_udata:
		var i uint64

		i, err = leb128.ReadUnsigned(r)
		if err != nil {
			break
		}

		offset = uint64(i)
	default:
		return nil, fmt.Errorf("Cannot read form %s data as a reference.", DwFormStr[f])
	}

	if err != nil {
		return nil, fmt.Errorf("Error reading form %s data.\n%s", DwFormStr[f], err.Error())
	}

	dieTree, err := d.readDIETree(u, u.headerOffset+offset)
	if err != nil {
		err = fmt.Errorf(
			"Error reading DIE tree at offset %d specified by form %s.\n%s",
			offset, DwFormStr[f], err.Error())
	}

	return dieTree, nil
}

func (d *DwData) readAttrByteSlice(
	u *DwUnit, r *bytes.Reader, f DwForm, en binary.ByteOrder) ([]byte, error) {
	var size uint64
	var err error

	switch f {
	case DW_FORM_block1:
		var i byte

		i, err = r.ReadByte()
		if err != nil {
			break
		}

		size = uint64(i)
	case DW_FORM_block2:
		var i uint16

		err = binary.Read(r, en, &i)
		if err != nil {
			break
		}

		size = uint64(i)
	case DW_FORM_block4:
		var i uint32

		err = binary.Read(r, en, &i)
		if err != nil {
			break
		}

		size = uint64(i)
	case DW_FORM_exprloc:
		fallthrough
	case DW_FORM_block:
		size, err = leb128.ReadUnsigned(r)
		if err != nil {
			break
		}
	default:
		return nil, fmt.Errorf("Cannot read form %s data a block of bytes.", DwFormStr[f])
	}

	if err != nil {
		err = fmt.Errorf(
			"Error reading block size of form %s data.\n%s", DwFormStr[f], err.Error())
		return nil, err
	}

	b := make([]byte, size)
	_, err = r.Read(b)
	if err != nil {
		err = fmt.Errorf(
			"Error reading %d-byte block of data for form %s.\n%s",
			size, DwFormStr[f], err.Error())
		return nil, err
	}

	return b, nil
}

func (d *DwData) readAttrLocList(
	u *DwUnit, r *bytes.Reader, form DwForm, en binary.ByteOrder) (LocList, error) {
	offset, err := d.readAttrUint64(u, r, form, en)
	if err != nil {
		err = fmt.Errorf("Error reading .debug_loc offset.\n%s", err.Error())
		return nil, err
	}
	return d.readLocList(u, offset, en)
}

func (d *DwData) readAttrRangeList(
	u *DwUnit, r *bytes.Reader, form DwForm, en binary.ByteOrder) (RangeList, error) {
	offset, err := d.readAttrUint64(u, r, form, en)
	if err != nil {
		err = fmt.Errorf("Error reading .debug_ranges offset.\n%s", err.Error())
		return nil, err
	}
	return d.readRangeList(u, offset, en)
}
