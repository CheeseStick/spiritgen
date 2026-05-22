package lantern

import (
	"io"

	"spiritgen/internal/validation"
	"spiritgen/internal/xlsx"
)

// ParseResult is the aggregated outcome of parsing a lantern XLSX file.
type ParseResult struct {
	Success []Household
	Errors  []validation.RowError
}

// ParseXLSX reads the first sheet of an XLSX document and groups rows into
// Households. The expected schema is:
//
//	[ address (주소) | name (이름) | relation (관계, optional) | dharma_name (법명, optional) ]
//
// Required: address (only on the first row of each household), name.
// Optional: relation, dharma_name. A row with an empty address column belongs
// to the previous household. Unlike spirit tablets, households are not split —
// every person of a household must fit on a single lantern tablet.
func ParseXLSX(r io.Reader) (ParseResult, error) {
	rows, err := xlsx.ReadFirstSheet(r)
	if err != nil {
		return ParseResult{}, err
	}

	var (
		result  ParseResult
		current *Household
	)

	for i, row := range rows {
		if i == 0 {
			continue // header row
		}
		if len(row) < 2 {
			result.Errors = append(result.Errors, validation.RowError{
				RowIndex: i + 1,
				Errors: []validation.Error{{
					Code:    ErrRowColumnTooShort,
					Field:   "row",
					Message: "입력 칸이 부족합니다 (4개 열 필요 - 주소|이름|관계|법명)",
				}},
			})
			continue
		}

		address := xlsx.NormalizeString(row[0])
		name := xlsx.NormalizeString(row[1])
		relation := ""
		dharma := ""

		if 3 <= len(row) {
			relation = xlsx.NormalizeString(row[2])
		}

		if 4 <= len(row) {
			dharma = xlsx.NormalizeString(row[3])
		}

		var rowErrors []validation.Error
		if address != "" {
			rowErrors = append(rowErrors, ValidateAddress(address)...)
		}
		rowErrors = append(rowErrors, ValidatePerson(name)...)

		if len(rowErrors) > 0 {
			result.Errors = append(result.Errors, validation.RowError{
				RowIndex: i + 1,
				Errors:   rowErrors,
			})
			continue
		}

		person := Person{Name: name, Relation: relation, DharmaName: dharma}

		if address != "" {
			if current != nil {
				result.Success = append(result.Success, *current)
			}
			current = &Household{Address: address, Persons: []Person{person}}
			continue
		}

		if current == nil {
			result.Errors = append(result.Errors, validation.RowError{
				RowIndex: i + 1,
				Errors: []validation.Error{{
					Code:    ErrMissingAddress,
					Field:   "address",
					Message: "주소 없이 사람 정보를 입력할 수 없습니다.",
				}},
			})
			continue
		}

		current.Persons = append(current.Persons, person)
	}

	if current != nil {
		result.Success = append(result.Success, *current)
	}
	return result, nil
}
