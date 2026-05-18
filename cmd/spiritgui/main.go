package main

import (
	"spiritgen/internal/gui"

	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := a.NewWindow("SpiritGen")
	w.SetMaster()
	gui.BuildUI(w)
	w.ShowAndRun()
}
