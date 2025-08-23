package parser

func IsBlank(s string) bool {
	return NormalizeString(s) == ""
}

func ValidatePresenter(s string) []ValidationError {
	var errors []ValidationError

	if IsBlank(s) {
		errors = append(errors, ValidationError{
			Code:    ValidationErrorMissingPresenterName,
			Field:   "presented_by",
			Message: "복위(伏爲)의 이름은 필수로 입력해야합니다.",
		})
	}

	return errors
}

func ValidateDeceasedInput(name, dharma, clan, relation string) []ValidationError {
	var errors []ValidationError

	if IsBlank(name) && IsBlank(dharma) {
		errors = append(errors, ValidationError{
			Code:    ValidationErrorMissingDeceasedName,
			Field:   "name/dharma_name",
			Message: "망자의 이름 또는 법명은 필수로 입력해야합니다.",
		})
	}

	if IsBlank(relation) {
		errors = append(errors, ValidationError{
			Code:    ValidationErrorMissingDeceasedRelation,
			Field:   "relation",
			Message: "망자와 복위의 관계는 필수로 입력해야합니다.",
		})
	}

	return errors
}
