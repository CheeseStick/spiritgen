package lantern

import "spiritgen/internal/validation"

const (
	ErrRowColumnTooShort validation.ErrorCode = "ROW_COLUMN_TOO_SHORT"
	ErrMissingAddress    validation.ErrorCode = "MISSING_ADDRESS"
	ErrMissingName       validation.ErrorCode = "MISSING_NAME"
	ErrMissingRelation   validation.ErrorCode = "MISSING_RELATION"
)
