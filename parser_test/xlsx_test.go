package parser

import (
	"bytes"
	"os"
	"path/filepath"
	"spiritgen/internal/parser"
	"testing"
)

func TestParseFromXLSX(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}

	t.Logf("Current working dir: %s", wd)

	abs := filepath.Join(wd, "..", "testdata", "template.xlsx")
	data, err := os.ReadFile(abs)
	if err != nil {
		t.Fatal(err)
	}
	r := bytes.NewReader(data)

	result, err := parser.ParseFromXLSX(r)
	if err != nil {
		t.Fatalf("ParseFromXLSX failed: %v", err)
	}

	if len(result.Success) == 0 {
		t.Error("expected at least one sprint tablet")
	}

	if len(result.Success) != 6 {
		t.Error("expected 6 sprint tablets, got ", len(result.Success))
	}

	if result.Success[0].PresentedBy != "1" || len(result.Success[0].DeceasedList) != 3 {
		t.Error("expected 3 Sprint tablets presented by 1, got ", result.Success[0].PresentedBy, len(result.Success))
	}

	if result.Success[1].PresentedBy != "2" || result.Success[1].DeceasedList[0].Name != "2-망자-1" {
		t.Error("expected sprint tablets presented by 2, got ", result.Success[1].PresentedBy, result.Success[1].DeceasedList[0].Name)
	}

	if result.Success[2].PresentedBy != "2" || result.Success[2].DeceasedList[0].Name != "2-망자-2" {
		t.Error("expected sprint tablets presented by 2, got ", result.Success[1].PresentedBy, result.Success[1].DeceasedList[0].Name)
	}
}
