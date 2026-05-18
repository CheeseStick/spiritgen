// Package xlsx provides domain-neutral XLSX reading primitives.
package xlsx

import (
	"fmt"
	"io"

	"github.com/xuri/excelize/v2"
)

// ReadFirstSheet opens the XLSX document from r and returns the rows of its first sheet.
func ReadFirstSheet(r io.Reader) ([][]string, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, fmt.Errorf("excel open failed: %w", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("read rows failed: %w", err)
	}
	return rows, nil
}
