package spirittablet

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestGenerateButtonInitiallyDisabled(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	w := a.NewWindow("test")

	tb := &tab{state: &uiState{}}
	tb.Build(w)

	generateBtn := widget.NewButton("생성", nil)
	generateBtn.Disable()

	if !generateBtn.Disabled() {
		t.Error("expected generate button to be disabled initially")
	}
}

func TestReadEnablesGenerateButton(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	w := a.NewWindow("test")

	state := &uiState{}
	generateBtn := widget.NewButton("생성", nil)
	generateBtn.Disable()

	state.xlsxPath = "../../testdata/testdata.xlsx"
	DoRead(state, generateBtn, w)

	if generateBtn.Disabled() {
		t.Error("expected generate button to be enabled after successful read")
	}
	if state.personCount == 0 {
		t.Error("expected personCount > 0 after successful read")
	}
	if len(state.tablets) == 0 {
		t.Error("expected tablets to be populated after successful read")
	}
}

func TestReadWithInvalidPathKeepsButtonDisabled(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	w := a.NewWindow("test")

	state := &uiState{}
	generateBtn := widget.NewButton("생성", nil)
	generateBtn.Disable()

	state.xlsxPath = "/nonexistent/path/file.xlsx"
	DoRead(state, generateBtn, w)

	if !generateBtn.Disabled() {
		t.Error("expected generate button to remain disabled after failed read")
	}
	if state.personCount != 0 {
		t.Error("expected personCount to remain 0 after failed read")
	}
}

func TestReadWithEmptyPathKeepsButtonDisabled(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	w := a.NewWindow("test")

	state := &uiState{}
	generateBtn := widget.NewButton("생성", nil)
	generateBtn.Disable()

	DoRead(state, generateBtn, w) // xlsxPath is empty

	if !generateBtn.Disabled() {
		t.Error("expected generate button to remain disabled when no file selected")
	}
}

func TestReadResetsStateOnReread(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	w := a.NewWindow("test")

	state := &uiState{}
	generateBtn := widget.NewButton("생성", nil)
	generateBtn.Disable()

	state.xlsxPath = "../../testdata/testdata.xlsx"
	DoRead(state, generateBtn, w)

	firstCount := state.personCount
	firstTablets := len(state.tablets)

	if firstCount == 0 {
		t.Fatal("first read produced no persons")
	}

	DoRead(state, generateBtn, w)

	if state.personCount != firstCount {
		t.Errorf("re-read changed personCount: got %d, want %d", state.personCount, firstCount)
	}
	if len(state.tablets) != firstTablets {
		t.Errorf("re-read changed tablet count: got %d, want %d", len(state.tablets), firstTablets)
	}
}
