package lantern

import (
	"io"
	"spiritgen/internal/validation"
	"spiritgen/internal/xlsx"
	"strings"
)

// Recognized sheet names. Sheets not listed here are silently ignored.
const (
	SheetBig      = "큰등"
	SheetFamily   = "가족등"
	SheetSpirit   = "영가등"
	SheetBusiness = "사업등"
)

// ParseResult is the aggregated outcome of parsing a lantern XLSX workbook.
// Each lantern kind is parsed from its own sheet; errors are kept per-sheet
// so the GUI/CLI can surface them with the originating sheet name.
type ParseResult struct {
	Big      []FamilyHousehold
	Family   []FamilyHousehold
	Spirit   []SpiritHousehold
	Business []BusinessHousehold

	// Errors maps sheet name → row errors encountered on that sheet.
	// Entries are present even when there were no errors (empty slice) so
	// callers can detect "sheet existed but produced no households".
	Errors map[string][]validation.RowError
}

// TotalPersonCount returns the total number of individual persons / names
// across all lantern kinds.
func (r ParseResult) TotalPersonCount() int {
	total := 0
	for _, h := range r.Big {
		total += len(h.Members)
	}
	for _, h := range r.Family {
		total += len(h.Members)
	}
	for _, h := range r.Spirit {
		total += len(h.Names)
	}
	for _, h := range r.Business {
		total += len(h.Names)
	}
	return total
}

// TotalHouseholdCount returns the total number of households across all kinds.
func (r ParseResult) TotalHouseholdCount() int {
	return len(r.Big) + len(r.Family) + len(r.Spirit) + len(r.Business)
}

// TotalErrorCount returns the sum of row errors across all sheets.
func (r ParseResult) TotalErrorCount() int {
	total := 0
	for _, errs := range r.Errors {
		total += len(errs)
	}
	return total
}

// ParseXLSX reads an XLSX workbook and parses each recognized sheet into its
// corresponding household type. Unrecognized sheets are ignored.
func ParseXLSX(r io.Reader) (ParseResult, error) {
	sheets, err := xlsx.ReadAllSheets(r)
	if err != nil {
		return ParseResult{}, err
	}

	result := ParseResult{
		Errors: make(map[string][]validation.RowError),
	}

	if rows, ok := sheets[SheetBig]; ok {
		hs, errs := parseFamilyRows(rows)
		result.Big = hs
		result.Errors[SheetBig] = errs
	}
	if rows, ok := sheets[SheetFamily]; ok {
		hs, errs := parseFamilyRows(rows)
		result.Family = hs
		result.Errors[SheetFamily] = errs
	}
	if rows, ok := sheets[SheetSpirit]; ok {
		hs, errs := parseSpiritRows(rows)
		result.Spirit = hs
		result.Errors[SheetSpirit] = errs
	}
	if rows, ok := sheets[SheetBusiness]; ok {
		hs, errs := parseBusinessRows(rows)
		result.Business = hs
		result.Errors[SheetBusiness] = errs
	}

	return result, nil
}

// parseFamilyRows parses rows in the 가족등 / 큰등 schema:
//
//	[ address | name | relation (optional) | dharma_name (optional) ]
//
// An empty address column means the row belongs to the previous household.
func parseFamilyRows(rows [][]string) (households []FamilyHousehold, errs []validation.RowError) {
	var current *FamilyHousehold

	for i, row := range rows {
		if i == 0 {
			continue // header
		}
		if len(row) < 2 {
			if isEmptyRow(row) {
				continue
			}
			errs = append(errs, validation.RowError{
				RowIndex: i + 1,
				Errors: []validation.Error{{
					Code:    ErrRowColumnTooShort,
					Field:   "row",
					Message: "입력 칸이 부족합니다 (최소 2개 필요 - 주소|이름)",
				}},
			})
			continue
		}

		address := xlsx.NormalizeString(row[0])
		name := xlsx.NormalizeString(row[1])
		relation := ""
		dharma := ""
		if len(row) >= 3 {
			relation = xlsx.NormalizeString(row[2])
		}
		if len(row) >= 4 {
			dharma = xlsx.NormalizeString(row[3])
		}

		if address == "" && name == "" && relation == "" && dharma == "" {
			continue // empty row
		}

		var rowErrors []validation.Error
		if address != "" {
			rowErrors = append(rowErrors, ValidateAddress(address)...)
		}
		rowErrors = append(rowErrors, ValidatePerson(name)...)
		if len(rowErrors) > 0 {
			errs = append(errs, validation.RowError{RowIndex: i + 1, Errors: rowErrors})
			continue
		}

		member := FamilyMember{Name: name, Relation: relation, DharmaName: dharma}

		if address != "" {
			if current != nil {
				households = append(households, *current)
			}
			current = &FamilyHousehold{Address: address, Members: []FamilyMember{member}}
			continue
		}

		if current == nil {
			errs = append(errs, validation.RowError{
				RowIndex: i + 1,
				Errors: []validation.Error{{
					Code:    ErrMissingAddress,
					Field:   "address",
					Message: "주소 없이 사람 정보를 입력할 수 없습니다.",
				}},
			})
			continue
		}
		current.Members = append(current.Members, member)
	}

	if current != nil {
		households = append(households, *current)
	}
	return
}

// parseSpiritRows parses rows in the 영가등 schema:
//
//	[ address | presenter (복위) | name ]
//
// First row of each household needs all three. Continuation rows (empty
// address column) only need name; column 2 (presenter) is ignored on
// continuation rows so callers can leave it blank.
func parseSpiritRows(rows [][]string) (households []SpiritHousehold, errs []validation.RowError) {
	var current *SpiritHousehold

	for i, row := range rows {
		if i == 0 {
			continue // header
		}
		if len(row) < 3 {
			if isEmptyRow(row) {
				continue
			}
			errs = append(errs, validation.RowError{
				RowIndex: i + 1,
				Errors: []validation.Error{{
					Code:    ErrRowColumnTooShort,
					Field:   "row",
					Message: "입력 칸이 부족합니다 (3개 열 필요 - 주소|복위|이름)",
				}},
			})
			continue
		}

		address := xlsx.NormalizeString(row[0])
		presenter := xlsx.NormalizeString(row[1])
		name := xlsx.NormalizeString(row[2])

		if address == "" && presenter == "" && name == "" {
			continue
		}

		var rowErrors []validation.Error
		if address != "" {
			rowErrors = append(rowErrors, ValidateAddress(address)...)
			rowErrors = append(rowErrors, ValidatePresenter(presenter)...)
		}
		rowErrors = append(rowErrors, ValidatePerson(name)...)
		if len(rowErrors) > 0 {
			errs = append(errs, validation.RowError{RowIndex: i + 1, Errors: rowErrors})
			continue
		}

		if !strings.HasSuffix(name, "영가") {
			name = name + " 영가" // auto-append "영가" for convenience, since it's implied in this sheet's context
		}

		if address != "" {
			if current != nil {
				households = append(households, *current)
			}
			current = &SpiritHousehold{Address: address, PresentedBy: presenter, Names: []string{name}}
			continue
		}

		if current == nil {
			errs = append(errs, validation.RowError{
				RowIndex: i + 1,
				Errors: []validation.Error{{
					Code:    ErrMissingAddress,
					Field:   "address",
					Message: "주소 없이 사람 정보를 입력할 수 없습니다.",
				}},
			})
			continue
		}

		current.Names = append(current.Names, name)
	}

	if current != nil {
		households = append(households, *current)
	}
	return
}

// parseBusinessRows parses rows in the 사업등 schema:
//
//	[ address | business_name (사업체명) | name ]
//
// Grouping rules mirror parseSpiritRows.
func parseBusinessRows(rows [][]string) (households []BusinessHousehold, errs []validation.RowError) {
	var current *BusinessHousehold

	for i, row := range rows {
		if i == 0 {
			continue // header
		}
		if len(row) < 3 {
			if isEmptyRow(row) {
				continue
			}
			errs = append(errs, validation.RowError{
				RowIndex: i + 1,
				Errors: []validation.Error{{
					Code:    ErrRowColumnTooShort,
					Field:   "row",
					Message: "입력 칸이 부족합니다 (3개 열 필요 - 주소|사업체명|이름)",
				}},
			})
			continue
		}

		address := xlsx.NormalizeString(row[0])
		bizName := xlsx.NormalizeString(row[1])
		name := xlsx.NormalizeString(row[2])

		if address == "" && bizName == "" && name == "" {
			continue
		}

		var rowErrors []validation.Error
		if address != "" {
			rowErrors = append(rowErrors, ValidateAddress(address)...)
			rowErrors = append(rowErrors, ValidateBusinessName(bizName)...)
		}
		rowErrors = append(rowErrors, ValidatePerson(name)...)
		if len(rowErrors) > 0 {
			errs = append(errs, validation.RowError{RowIndex: i + 1, Errors: rowErrors})
			continue
		}

		if address != "" {
			if current != nil {
				households = append(households, *current)
			}
			current = &BusinessHousehold{Address: address, BusinessName: bizName, Names: []string{name}}
			continue
		}

		if current == nil {
			errs = append(errs, validation.RowError{
				RowIndex: i + 1,
				Errors: []validation.Error{{
					Code:    ErrMissingAddress,
					Field:   "address",
					Message: "주소 없이 사람 정보를 입력할 수 없습니다.",
				}},
			})
			continue
		}
		current.Names = append(current.Names, name)
	}

	if current != nil {
		households = append(households, *current)
	}
	return
}

// isEmptyRow returns true if every cell in the row is blank after normalization.
func isEmptyRow(row []string) bool {
	for _, c := range row {
		if !xlsx.IsBlank(c) {
			return false
		}
	}
	return true
}
