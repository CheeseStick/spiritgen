package lantern

import (
	"spiritgen/internal/validation"
	"spiritgen/internal/xlsx"
)

// ValidateAddress checks that the household address is non-blank.
func ValidateAddress(s string) []validation.Error {
	if xlsx.IsBlank(s) {
		return []validation.Error{{
			Code:    ErrMissingAddress,
			Field:   "address",
			Message: "주소는 필수로 입력해야 합니다.",
		}}
	}
	return nil
}

// ValidatePerson checks that a person's name is non-blank. Relation and dharma
// name are optional.
func ValidatePerson(name string) []validation.Error {
	if xlsx.IsBlank(name) {
		return []validation.Error{{
			Code:    ErrMissingName,
			Field:   "name",
			Message: "이름은 필수로 입력해야 합니다.",
		}}
	}
	return nil
}

// ValidatePresenter checks that the 영가등 presenter (복위) is non-blank.
func ValidatePresenter(s string) []validation.Error {
	if xlsx.IsBlank(s) {
		return []validation.Error{{
			Code:    ErrMissingPresenter,
			Field:   "presented_by",
			Message: "복위(伏爲)는 필수로 입력해야 합니다.",
		}}
	}
	return nil
}

// ValidateBusinessName checks that the 사업등 business name is non-blank.
func ValidateBusinessName(s string) []validation.Error {
	if xlsx.IsBlank(s) {
		return []validation.Error{{
			Code:    ErrMissingBusinessName,
			Field:   "business_name",
			Message: "사업체명은 필수로 입력해야 합니다.",
		}}
	}
	return nil
}
