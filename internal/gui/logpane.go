package gui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// LogPane is a scrollable, read-only multi-line text widget used by product
// tabs to surface row-level errors and progress messages. Each call to Append
// also writes the line to the standard logger (stderr) so the same info shows
// up in the console.
//
// The pane is always visible; an empty pane displays the configured
// placeholder.
type LogPane struct {
	entry  *widget.Entry
	scroll *container.Scroll
}

// NewLogPane constructs a LogPane that is visible from the start.
func NewLogPane() *LogPane {
	entry := widget.NewMultiLineEntry()
	entry.Disable()
	entry.SetPlaceHolder("읽기 결과가 여기에 표시됩니다…")
	entry.Wrapping = fyne.TextWrapWord

	scroll := container.NewVScroll(entry)
	scroll.SetMinSize(fyne.NewSize(0, 140))

	return &LogPane{entry: entry, scroll: scroll}
}

// Widget returns the underlying fyne.CanvasObject for placement in a layout.
func (p *LogPane) Widget() fyne.CanvasObject {
	return p.scroll
}

// Append writes a line to the pane and to the console log.
func (p *LogPane) Append(line string) {
	log.Println(line)
	if p == nil || p.entry == nil {
		return
	}
	if p.entry.Text == "" {
		p.entry.SetText(line)
	} else {
		p.entry.SetText(p.entry.Text + "\n" + line)
	}
}

// Clear empties the pane (it remains visible).
func (p *LogPane) Clear() {
	if p == nil || p.entry == nil {
		return
	}
	p.entry.SetText("")
}
