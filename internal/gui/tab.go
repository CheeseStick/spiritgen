package gui

import "fyne.io/fyne/v2"

// Tab is implemented by product packages that contribute a tab to the main window.
type Tab interface {
	Title() string
	Build(w fyne.Window) fyne.CanvasObject
}
