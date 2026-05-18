package lantern

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
	if state.householdCount != 0 {
		t.Error("expected householdCount to remain 0 after failed read")
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
