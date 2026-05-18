package spirittablet

import (
	"io"

	"spiritgen/internal/validation"
	"spiritgen/internal/xlsx"
)

// ParseResult is the aggregated outcome of parsing an XLSX file:
// successful SpiritTablets plus per-row validation errors.
type ParseResult struct {
	Success []SpiritTablet
	Errors  []validation.RowError
}

// ParseXLSX reads the first sheet of an XLSX document and groups rows into
// SpiritTablets. The expected schema is a 5-column row:
//
//	[ presenter (복위) | name (망자이름) | dharma_name (망자법명) | clan_origin (망자본관) | relation (관계) ]
//
// The presenter column is optional after the first row of a group: empty
// presenter means the row belongs to the previous tablet. A tablet flushes
// once it accumulates 3 deceased entries, with subsequent rows starting a new
// tablet under the same presenter.
func ParseXLSX(r io.Reader) (ParseResult, error) {
	rows, err := xlsx.ReadFirstSheet(r)
	if err != nil {
		return ParseResult{}, err
	}

	var (
		result        ParseResult
		currentTablet *SpiritTablet
	)

	for i, row := range rows {
		if i == 0 {
			continue // header
		}
		if len(row) < 5 {
			result.Errors = append(result.Errors, validation.RowError{
				RowIndex: i + 1,
				Errors: []validation.Error{{
					Code:    ErrRowColumnTooShort,
					Field:   "row",
					Message: "입력 칸이 부족합니다 (5개 열 필요 - 복위|망자이름|망자법명|망자본관|관계)",
				}},
			})
			continue
		}

		presenter := xlsx.NormalizeString(row[0])
		name := xlsx.NormalizeString(row[1])
		dharma := xlsx.NormalizeString(row[2])
		clan := xlsx.NormalizeString(row[3])
		relation := xlsx.NormalizeString(row[4])

		var rowErrors []validation.Error
		if presenter != "" {
			rowErrors = append(rowErrors, ValidatePresenter(presenter)...)
		}
		rowErrors = append(rowErrors, ValidateDeceasedInput(name, dharma, clan, relation)...)

		if len(rowErrors) > 0 {
			result.Errors = append(result.Errors, validation.RowError{
				RowIndex: i + 1,
				Errors:   rowErrors,
			})
			continue
		}

		deceased := Deceased{
			Name:       name,
			DharmaName: dharma,
			ClanOrigin: clan,
			Relation:   relation,
		}

		if presenter != "" {
			if currentTablet != nil {
				result.Success = append(result.Success, *currentTablet)
			}
			currentTablet = &SpiritTablet{
				PresentedBy:  presenter,
				DeceasedList: []Deceased{deceased},
			}
			continue
		}

		if currentTablet == nil {
			result.Errors = append(result.Errors, validation.RowError{
				RowIndex: i + 1,
				Errors: []validation.Error{{
					Code:    ErrMissingPresenterName,
					Field:   "presented_by",
					Message: "복위자 정보가 없는 상태에서 망자 정보를 입력할 수 없습니다.",
				}},
			})
			continue
		}

		if 3 <= len(currentTablet.DeceasedList) {
			result.Success = append(result.Success, *currentTablet)
			currentTablet = &SpiritTablet{
				PresentedBy:  currentTablet.PresentedBy,
				DeceasedList: []Deceased{deceased},
			}
		} else {
			currentTablet.DeceasedList = append(currentTablet.DeceasedList, deceased)
		}
	}

	if currentTablet != nil {
		result.Success = append(result.Success, *currentTablet)
	}
	return result, nil
}
