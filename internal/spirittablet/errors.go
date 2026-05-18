package spirittablet

import "spiritgen/internal/validation"

const (
	ErrRowColumnTooShort       validation.ErrorCode = "ROW_COLUMN_TOO_SHORT"
	ErrMissingPresenterName    validation.ErrorCode = "MISSING_PRESENTER"
	ErrMissingDeceasedName     validation.ErrorCode = "MISSING_DECEASED_NAME"
	ErrMissingDeceasedRelation validation.ErrorCode = "MISSING_DECEASED_RELATION"
)
