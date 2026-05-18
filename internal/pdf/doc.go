// Package pdf provides a thin presentation-agnostic wrapper over fpdf,
// exposing A4 document creation and font/image registration primitives.
package pdf

import (
	"bytes"

	"codeberg.org/go-pdf/fpdf"
)

// Doc wraps *fpdf.Fpdf so domain packages can call lower-level methods directly
// while still using the helpers defined here.
type Doc struct {
	*fpdf.Fpdf
}

// NewA4Landscape returns a new A4-landscape PDF document with the given margins.
// SetAutoPageBreak is configured with marginY as the bottom margin.
func NewA4Landscape(marginX, marginY float64) *Doc {
	pdf := fpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(marginX, marginY, -1)
	pdf.SetAutoPageBreak(true, marginY)
	return &Doc{Fpdf: pdf}
}

// RegisterUTF8Font registers a UTF-8 TrueType font from in-memory bytes.
func (d *Doc) RegisterUTF8Font(name, style string, ttf []byte) {
	d.AddUTF8FontFromBytes(name, style, ttf)
}

// RegisterPNGImage registers a PNG image from bytes under the given key,
// which can later be referenced by Image(key, ...).
func (d *Doc) RegisterPNGImage(key string, png []byte) {
	d.RegisterImageOptionsReader(key, fpdf.ImageOptions{ImageType: "PNG"}, bytes.NewReader(png))
}
