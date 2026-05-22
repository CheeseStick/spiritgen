package lantern

import (
	"spiritgen/assets"
	"spiritgen/internal/pdf"
)

// RenderPDF writes every lantern kind in result as a single PDF at outputPath.
// Each kind that has at least one household starts on a fresh page. title is
// rendered at the top of every tablet; empty title falls back to defaultTitle.
//
// bgImage is the per-tablet background PNG; for now the same image is used
// for every kind. Per-kind backgrounds can be plumbed through later by
// extending Layout with its own image bytes / key.
func RenderPDF(result ParseResult, outputPath string, bgImage []byte, title string) error {
	if title == "" {
		title = defaultTitle
	}

	doc := pdf.NewA4Landscape(a4PaperMarginX, a4PaperMarginY)
	// We ship three weights of NotoSerifKR — register all three so any layout
	// can request them via SetFont(..., "B"|"M"|"L", ...).
	doc.RegisterUTF8Font(familyLayout.LabelFontName, "B", assets.FontNotoSerifKR)
	doc.RegisterUTF8Font(familyLayout.LabelFontName, "M", assets.FontNotoSerifKRMedium)
	doc.RegisterUTF8Font(familyLayout.LabelFontName, "L", assets.FontNotoSerifKRLight)
	doc.RegisterPNGImage(backgroundImageKey, bgImage)

	// Each non-empty section starts on its own page.
	if len(result.Big) > 0 {
		doc.AddPage()
		renderBigSection(doc, result.Big, title)
	}
	if len(result.Family) > 0 {
		doc.AddPage()
		renderFamilySection(doc, result.Family, title)
	}
	if len(result.Spirit) > 0 {
		doc.AddPage()
		renderSpiritSection(doc, result.Spirit, title)
	}
	if len(result.Business) > 0 {
		doc.AddPage()
		renderBusinessSection(doc, result.Business, title)
	}

	return doc.OutputFileAndClose(outputPath)
}
