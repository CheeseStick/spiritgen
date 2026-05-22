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
// for every kind. Per-kind backgrounds can be plumbed through later.
func RenderPDF(result ParseResult, outputPath string, bgImage []byte, title string) error {
	if title == "" {
		title = defaultTitle
	}

	doc := pdf.NewA4Landscape(a4PaperMarginX, a4PaperMarginY)
	doc.RegisterUTF8Font(labelFontName, "B", assets.FontNotoSerifKR)
	doc.RegisterUTF8Font(labelFontName, "M", assets.FontNotoSerifKRMedium)
	doc.RegisterUTF8Font(labelFontName, "L", assets.FontNotoSerifKRLight)
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

// placeAndRender runs the standard "place tablets left-to-right, wrap to next
// row, wrap to next page" cursor loop. The renderOne callback draws a single
// tablet at the given (x, y). Shared by every section so per-kind functions
// only need to supply their tablet-drawing closure.
func placeAndRender(doc *pdf.Doc, count int, renderOne func(originX, originY float64, i int)) {
	pageWidth, pageHeight := doc.GetPageSize()
	leftMargin, topMargin, rightMargin, bottomMargin := doc.GetMargins()
	pageRight := pageWidth - rightMargin
	pageBottom := pageHeight - bottomMargin

	cursorX := leftMargin
	cursorY := topMargin

	for i := 0; i < count; i++ {
		if cursorX+tabletSize > pageRight {
			cursorX = leftMargin
			cursorY += tabletSize + tabletMarginY
		}
		if cursorY+tabletSize > pageBottom {
			doc.AddPage()
			cursorX = leftMargin
			cursorY = topMargin
		}
		renderOne(cursorX, cursorY, i)
		cursorX += tabletSize + tabletMarginX
	}
}

// drawTabletBackground draws the background image (aspect-fit, centered) and
// then the circular cut-guideline on top. Shared by every per-kind renderer.
func drawTabletBackground(doc *pdf.Doc, originX, originY float64) {
	bgOffset := (tabletSize - backgroundSize) / 2
	doc.DrawImageAspectFit(backgroundImageKey,
		originX+bgOffset, originY+bgOffset, backgroundSize, backgroundSize)

	prevLW := doc.GetLineWidth()
	doc.SetLineWidth(circleLineWidth)
	doc.Circle(originX+tabletSize/2, originY+tabletSize/2, tabletSize/2, "D")
	doc.SetLineWidth(prevLW)
}

// drawTitle draws the (centered) title at the top of the tablet and returns
// the Y position where the next line of content should begin.
func drawTitle(doc *pdf.Doc, originX, originY float64, title string) float64 {
	y := originY + tabletPaddingTop
	if title == "" {
		return y
	}
	doc.SetFont(labelFontName, titleFontStyle, titleFontSize)
	titleLH := lineHeightFor(titleFontSize)
	doc.DrawTextCentered(title, originX+tabletPaddingX, tabletSize-2*tabletPaddingX, y+titleLH)
	return y + titleLH + titleBottomMargin
}

// drawAddress draws the (centered) address and returns the Y position where
// the next line of content should begin.
func drawAddress(doc *pdf.Doc, originX, y float64, address string) float64 {
	doc.SetFont(labelFontName, addressFontStyle, addressFontSize)
	addressLH := lineHeightFor(addressFontSize)
	doc.DrawTextCentered(address, originX+tabletPaddingX, tabletSize-2*tabletPaddingX, y+addressLH)
	return y + addressLH
}

// renderCenteredNameList draws each name centered horizontally within the
// box (boxX, boxX+boxWidth) on its own line. Used by 영가등 / 사업등 where
// each entry is just a name (no relation/dharma). Lines past tabletBottom
// are skipped.
func renderCenteredNameList(doc *pdf.Doc, names []string, cfg PersonFontConfig, boxX, y, boxWidth, tabletBottom float64) {
	doc.SetFont(labelFontName, cfg.nameFontStyle, cfg.nameFontSize)
	lh := lineHeightFor(cfg.nameFontSize)
	for _, name := range names {
		if y+lh > tabletBottom {
			break
		}
		doc.DrawTextCentered(name, boxX, boxWidth, y+lh)
		y += lh
	}
}
