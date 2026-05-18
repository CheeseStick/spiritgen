package lantern_test

import (
	"bytes"
	"testing"

	"spiritgen/assets"
	"spiritgen/internal/lantern"
)

// TestParseXLSX_Template smoke-tests that the embedded template parses without
// error. The template may include sample rows; the assertion only checks that
// the parser returns without an error. A richer parser test should be added
// once testdata/lantern_testdata.xlsx is available.
func TestParseXLSX_Template(t *testing.T) {
	r := bytes.NewReader(assets.LanternTabletTemplate)
	if _, err := lantern.ParseXLSX(r); err != nil {
		t.Fatalf("ParseXLSX(template) failed: %v", err)
	}
}
