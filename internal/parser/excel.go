package parser

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"spiritgen/internal/model"
)

type ParseResult struct {
	Success []model.SpiritTablet
	Errors  []RowValidationError
}

func ParseFromXLSX(r io.Reader) (ParseResult, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return ParseResult{}, fmt.Errorf("excel open failed: %w", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return ParseResult{}, fmt.Errorf("read rows failed: %w", err)
	}

	var (
		result        ParseResult
		currentTablet *model.SpiritTablet
	)

	for i, row := range rows {
		if i == 0 {
			continue // Skips header
		}
		if len(row) < 5 {
			result.Errors = append(result.Errors, RowValidationError{
				RowIndex: i + 1,
				Errors: []ValidationError{{
					Code:    ValidationErrorRowColumnTooShort,
					Field:   "row",
					Message: "입력 칸이 부족합니다 (5개 열 필요 - 복위|망자이름|망자법명|망자본관|관계)",
				}},
			})
			continue
		}

		// Normalize inputs
		presenter := NormalizeString(row[0])
		name := NormalizeString(row[1])
		dharma := NormalizeString(row[2])
		clan := NormalizeString(row[3])
		relation := NormalizeString(row[4])

		// Validate
		var rowErrors []ValidationError
		if presenter != "" { // Conditionally check presenter as it can be omitted
			rowErrors = append(rowErrors, ValidatePresenter(presenter)...)
		}
		rowErrors = append(rowErrors, ValidateDeceasedInput(name, dharma, clan, relation)...)

		if len(rowErrors) > 0 {
			result.Errors = append(result.Errors, RowValidationError{
				RowIndex: i + 1,
				Errors:   rowErrors,
			})
			continue
		}

		deceased := model.Deceased{
			Name:       name,
			DharmaName: dharma,
			ClanOrigin: clan,
			Relation:   relation,
		}

		// Handle grouping by presenter
		if presenter != "" {
			// Flush previous tablet if exists
			if currentTablet != nil {
				result.Success = append(result.Success, *currentTablet)
			}
			currentTablet = &model.SpiritTablet{
				PresentedBy:  presenter,
				DeceasedList: []model.Deceased{deceased},
			}
		} else if currentTablet != nil {
			currentTablet.DeceasedList = append(currentTablet.DeceasedList, deceased)
		} else {
			// First row missing presenter
			result.Errors = append(result.Errors, RowValidationError{
				RowIndex: i + 1,
				Errors: []ValidationError{{
					Code:    ValidationErrorMissingPresenterName,
					Field:   "presented_by",
					Message: "복위자 정보가 없는 상태에서 망자 정보를 입력할 수 없습니다.",
				}},
			})
			continue
		}
	}

	if currentTablet != nil {
		result.Success = append(result.Success, *currentTablet)
	}

	return result, nil
}
