package spirittablet

import (
	"spiritgen/internal/validation"
	"spiritgen/internal/xlsx"
)

// ValidatePresenter checks that the presenter (복위) name is non-blank.
func ValidatePresenter(s string) []validation.Error {
	var errors []validation.Error
	if xlsx.IsBlank(s) {
		errors = append(errors, validation.Error{
			Code:    ErrMissingPresenterName,
			Field:   "presented_by",
			Message: "복위(伏爲)의 이름은 필수로 입력해야합니다.",
		})
	}
	return errors
}

// ValidateDeceasedInput checks that a deceased entry has either a name or a
// dharma name, and that the relation field is present.
func ValidateDeceasedInput(name, dharma, clan, relation string) []validation.Error {
	var errors []validation.Error
	if xlsx.IsBlank(name) && xlsx.IsBlank(dharma) {
		errors = append(errors, validation.Error{
			Code:    ErrMissingDeceasedName,
			Field:   "name/dharma_name",
			Message: "망자의 이름 또는 법명은 필수로 입력해야합니다.",
		})
	}
	if xlsx.IsBlank(relation) {
		errors = append(errors, validation.Error{
			Code:    ErrMissingDeceasedRelation,
			Field:   "relation",
			Message: "망자와 복위의 관계는 필수로 입력해야합니다.",
		})
	}
	return errors
}
