// #############################################################################
// This file is part of the "garf" package of the "Eureka" project.
// It is distributed under the MIT License. Refer to the LICENSE file for more
// information.
//
// Website: http://www.github.com/sivachandra/eureka
// #############################################################################

// Package garf provides API to read DWARF debug info from ELF files.
package garf

func (f DwForm) IsAddress() bool {
	return f == DW_FORM_addr
}

func (f DwForm) IsBlock() bool {
	switch f {
	case DW_FORM_block1, DW_FORM_block2, DW_FORM_block4:
		fallthrough
	case DW_FORM_block:
		return true
	default:
		return false
	}
}

func (f DwForm) IsFixedWidthConst() bool {
	switch f {
	case DW_FORM_data1, DW_FORM_data2, DW_FORM_data4, DW_FORM_data8:
		return true
	default:
		return false
	}
}

func (f DwForm) IsSignedVarWidthConst() bool {
	return f == DW_FORM_sdata
}

func (f DwForm) IsUnsignedVarWidthConst() bool {
	return f == DW_FORM_udata
}

func (f DwForm) IsConstant() bool {
	return f.IsFixedWidthConst() || f.IsSignedVarWidthConst() || f.IsUnsignedVarWidthConst()
}

func (f DwForm) IsExprLoc() bool {
	return f == DW_FORM_exprloc
}

func (f DwForm) IsFlag() bool {
	return f == DW_FORM_flag || f == DW_FORM_flag_present
}

func (f DwForm) IsLinePtr() bool {
	return f == DW_FORM_sec_offset
}

func (f DwForm) IsLocListPtr() bool {
	return f == DW_FORM_sec_offset
}

func (f DwForm) IsMacPtr() bool {
	return f == DW_FORM_sec_offset
}

func (f DwForm) IsRangeListPtr() bool {
	return f == DW_FORM_sec_offset
}

func (f DwForm) IsCompUnitRef() bool {
	switch f {
	case DW_FORM_ref1, DW_FORM_ref2, DW_FORM_ref4, DW_FORM_ref8:
		fallthrough
	case DW_FORM_ref_udata:
		return true
	default:
		return false
	}
}

func (f DwForm) IsGlobalRef() bool {
	return f == DW_FORM_ref_addr
}

func (f DwForm) IsSupRef() bool {
	return f == DW_FORM_ref_sup
}

func (f DwForm) IsTypeUnitRef() bool {
	return f == DW_FORM_ref_sig8
}

func (f DwForm) IsRef() bool {
	return f.IsCompUnitRef() || f.IsGlobalRef() || f.IsSupRef() || f.IsTypeUnitRef()
}

func (f DwForm) IsString() bool {
	switch f {
	case DW_FORM_string, DW_FORM_strp, DW_FORM_strx, DW_FORM_str_sup:
		return true
	default:
		return false
	}
}
