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
	"eureka/utils"
	"eureka/utils/leb128"
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
	case DW_AT_producer:
		fallthrough
	case DW_AT_comp_dir:
		fallthrough
	case DW_AT_linkage_name:
		fallthrough
	case DW_AT_name:
		attr.Value, err = d.readAttrStr(u, r, form, en)
	case DW_AT_language:
		attr.Value, err = d.readAttrUint16(u, r, form, en)
		attr.Value = DwLang(attr.Value.(uint16))
	case DW_AT_low_pc:
		fallthrough
	case DW_AT_high_pc:
		attr.Value, err = d.readAttrUint64(u, r, form, en)
	case DW_AT_stmt_list:
		attr.Value, err = d.readAttrUint64(u, r, form, en)
	case DW_AT_external:
		attr.Value, err = d.readAttrFlag(u, r, form, en)
	case DW_AT_decl_file:
		fallthrough
	case DW_AT_decl_line:
		attr.Value, err = d.readAttrUint32(u, r, form, en)
	case DW_AT_sibling:
		fallthrough
	case DW_AT_abstract_origin:
		fallthrough
	case DW_AT_specification:
		fallthrough
	case DW_AT_type:
		attr.Value, err = d.readAttrRef(u, r, form, en)
	case DW_AT_frame_base:
		attr.Value, err = d.readAttrByteSlice(u, r, form, en)
	case DW_AT_byte_size:
		attr.Value, err = d.readAttrUint32(u, r, form, en)
	case DW_AT_encoding:
		attr.Value, err = r.ReadByte()
		attr.Value = DwAte(attr.Value.(byte))
	case DW_AT_ranges:
		attr.Value, err = d.readAttrUint64(u, r, form, en)
	case DW_AT_location:
		if form == DW_FORM_sec_offset {
			attr.Value, err = d.readAttrUint64(u, r, form, en)
		} else if form == DW_FORM_exprloc {
			attr.Value, err = d.readAttrByteSlice(u, r, form, en)
		} else {
			err = fmt.Errorf("Unsupported form %d for DW_AT_Location.", form)
		}
	case DW_AT_declaration:
		attr.Value, err = d.readAttrFlag(u, r, form, en)
	case DW_AT_GNU_tail_call:
		fallthrough
	case DW_AT_GNU_all_call_sites:
		attr.Value, err = d.readAttrFlag(u, r, form, en)
	case DW_AT_GNU_call_site_value:
		attr.Value, err = d.readAttrByteSlice(u, r, form, en)
	default:
		attr.Value, err = d.readAttrByteSlice(u, r, form, en)
	}

	return attr, err
}

func (d *DwData) readAttrStr(
	u *DwUnit, r *bytes.Reader, form DwForm, en binary.ByteOrder) (string, error) {
	switch form {
	case DW_FORM_string:
		str, err := utils.ReadUntil(r, byte(0))
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

		debugStrMap, err := d.DebugStr()
		if err != nil {
			return "", fmt.Errorf("Error reading .debug_str.", err)
		}

		str, exists := debugStrMap[offset]
		if !exists {
			return "", fmt.Errorf("Invalid .debug_str offset", err)
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
		return 0, fmt.Errorf("Cannot read data of form %d as int16.", f)
	}

	err = fmt.Errorf("Error reading data of form %d.\n%s", f, err.Error())
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
		return 0, fmt.Errorf("Cannot read data of form %d as int32.", f)
	}

	err = fmt.Errorf("Error reading data of form %d.\n%s", f, err.Error())
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
		return 0, fmt.Errorf("Cannot read data of form %d as int64.", f)
	}

	err = fmt.Errorf("Error reading data of form %d.\n%s", f, err.Error())
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
		return 0, fmt.Errorf("Cannot read data of form %d as int16.", f)
	}

	err = fmt.Errorf("Error reading data of form %d.\n%s", f, err.Error())
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
		return 0, fmt.Errorf("Cannot read data of form %d as uint32.", f)
	}

	err = fmt.Errorf("Error reading data of form %d.\n%s", f, err.Error())
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
		return 0, fmt.Errorf("Cannot read data of form %d as uint64.", f)
	}

	err = fmt.Errorf("Error reading data of form %d.\n%s", f, err.Error())
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
		return false, fmt.Errorf("Cannot read data of form %d as a flag.", f)
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
	case DW_FORM_ref_udata:
		var i uint64

		i, err = leb128.ReadUnsigned(r)
		if err != nil {
			break
		}

		offset = uint64(i)
	default:
		return nil, fmt.Errorf("Cannot read form %d data as a reference.")
	}

	if err != nil {
		return nil, fmt.Errorf("Error reading form %d data.\n%s", err.Error())
	}

	dieTree, err := d.readDIETree(u, offset)
	if err != nil {
		err = fmt.Errorf(
			"Error reading DIE tree at offset %d specified by form %d.\n%s",
			offset, f, err.Error())
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
		return nil, fmt.Errorf("Cannot read form %d data a block of bytes.", f)
	}

	if err != nil {
		return nil, fmt.Errorf("Error reading block size of form %d data.\n%s", err.Error())
	}

	b := make([]byte, size)
	_, err = r.Read(b)
	if err != nil {
		err = fmt.Errorf(
			"Error reading %d block of data for form %d.\n%s", size, f, err.Error())
		return nil, err
	}

	return b, nil
}
