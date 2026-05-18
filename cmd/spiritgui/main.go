package main

import (
	"spiritgen/internal/gui"
	"spiritgen/internal/spirittablet"

	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := a.NewWindow("위패/인등/연등 생성 프로그램")
	w.SetMaster()

	gui.BuildUI(w, spirittablet.NewTab())

	w.ShowAndRun()
}
