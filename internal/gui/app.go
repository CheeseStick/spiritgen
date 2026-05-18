// Package gui provides the generic Fyne shell that hosts one or more product
// tabs. It exposes a Tab interface (see tab.go) and BuildUI, which is
// agnostic about which products are registered.
package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// BuildUI populates w with the given product tabs.
func BuildUI(w fyne.Window, tabs ...Tab) {
	items := make([]*container.TabItem, 0, len(tabs))
	for _, t := range tabs {
		items = append(items, container.NewTabItem(t.Title(), t.Build(w)))
	}
	w.SetContent(container.NewAppTabs(items...))
	w.Resize(fyne.NewSize(600, 280))
}
