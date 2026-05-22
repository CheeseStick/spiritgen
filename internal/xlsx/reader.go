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

// ReadAllSheets opens the XLSX document from r and returns a map of every
// sheet's name to its rows. Empty/non-existent sheets are still present in
// the map with an empty rows slice.
func ReadAllSheets(r io.Reader) (map[string][][]string, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, fmt.Errorf("excel open failed: %w", err)
	}
	defer f.Close()

	out := make(map[string][][]string)
	for _, name := range f.GetSheetList() {
		rows, err := f.GetRows(name)
		if err != nil {
			return nil, fmt.Errorf("read rows from sheet %q failed: %w", name, err)
		}
		out[name] = rows
	}
	return out, nil
}
