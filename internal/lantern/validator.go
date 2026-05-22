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

// ValidatePerson checks that name and relation are both present.
// Dharma name is optional and not checked here.
func ValidatePerson(name string) []validation.Error {
	var errs []validation.Error
	if xlsx.IsBlank(name) {
		errs = append(errs, validation.Error{
			Code:    ErrMissingName,
			Field:   "name",
			Message: "이름은 필수로 입력해야 합니다.",
		})
	}
	//if xlsx.IsBlank(relation) {
	//	errs = append(errs, validation.Error{
	//		Code:    ErrMissingRelation,
	//		Field:   "relation",
	//		Message: "관계는 필수로 입력해야 합니다.",
	//	})
	//}
	return errs
}
