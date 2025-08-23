package parser

import "fmt"

type ValidationErrorCode string

const (
	ValidationErrorRowColumnTooShort       ValidationErrorCode = "ROW_COLUMN_TOO_SHORT"      // Row의 Column이 모자랄 때 (5개 필요)
	ValidationErrorMissingPresenterName    ValidationErrorCode = "MISSING_PRESENTER"         // 복위 이름이 없을 때
	ValidationErrorMissingDeceasedName     ValidationErrorCode = "MISSING_DECEASED_NAME"     // 망자의 이름이 없을 때
	ValidationErrorMissingDeceasedRelation ValidationErrorCode = "MISSING_DECEASED_RELATION" // 망자와 복위의 관계가 없을 때
)

type ValidationError struct {
	Code    ValidationErrorCode
	Field   string
	Message string
}

type RowValidationError struct {
	RowIndex int
	Errors   []ValidationError
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e RowValidationError) Error() string {
	codes := make([]string, 0, len(e.Errors))

	for _, err := range e.Errors {
		codes = append(codes, string(err.Code))
	}

	return fmt.Sprintf("[%d] %s", e.RowIndex, codes)
}
