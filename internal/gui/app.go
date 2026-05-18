package gui

import (
	"spiritgen/assets"
	"spiritgen/internal/model"
	"spiritgen/internal/parser"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type AppState struct {
	xlsxPath       string
	parseResult    *parser.ParseResult
	tablets        []model.SpiritTablet
	personCount    int
	window         fyne.Window
	selectedDesign []byte
}

func BuildUI(w fyne.Window) {
	state := &AppState{window: w, selectedDesign: assets.DesignOne}
	tabs := container.NewAppTabs(
		container.NewTabItem("위패 생성", buildSpiritTabletTab(state, w)),
	)
	w.SetContent(tabs)
	w.Resize(fyne.NewSize(580, 220))
}
