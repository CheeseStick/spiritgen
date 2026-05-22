package lantern

import "spiritgen/internal/validation"

const (
	ErrRowColumnTooShort   validation.ErrorCode = "ROW_COLUMN_TOO_SHORT"
	ErrMissingAddress      validation.ErrorCode = "MISSING_ADDRESS"
	ErrMissingName         validation.ErrorCode = "MISSING_NAME"
	ErrMissingRelation     validation.ErrorCode = "MISSING_RELATION"
	ErrMissingPresenter    validation.ErrorCode = "MISSING_PRESENTER"     // 영가등 복위
	ErrMissingBusinessName validation.ErrorCode = "MISSING_BUSINESS_NAME" // 사업등 사업체명
)
