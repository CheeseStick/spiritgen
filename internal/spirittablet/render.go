package spirittablet

import (
	"fmt"
	"slices"
	"strings"

	"spiritgen/assets"
	"spiritgen/internal/pdf"
)

// RenderPDF writes the given tablets as a multi-page PDF at outputPath using
// bgImage (PNG bytes) as the per-tablet background.
func RenderPDF(tablets []SpiritTablet, outputPath string, bgImage []byte) error {
	doc := pdf.NewA4Landscape(a4PaperMarginX, a4PaperMarginY)
	doc.RegisterUTF8Font(labelFontName, labelFontStyle, assets.FontNotoSerifKR)
	doc.RegisterPNGImage(backgroundImageKey, bgImage)
	doc.AddPage()

	pageWidth, _ := doc.GetPageSize()
	leftMargin, _, rightMargin, _ := doc.GetMargins()
	pageContentSize := pageWidth - (leftMargin + rightMargin)

	cursorX := a4PaperMarginX
	for _, s := range tablets {
		nextContentEndX := cursorX + labelWidth + labelMarginX
		if pageContentSize < nextContentEndX {
			doc.AddPage()
			cursorX = a4PaperMarginX
		}
		renderTablet(doc, s, cursorX, a4PaperMarginY)
		cursorX += labelWidth + labelMarginX
	}

	return doc.OutputFileAndClose(outputPath)
}

func renderTablet(doc *pdf.Doc, tablet SpiritTablet, startX, startY float64) {
	if len(tablet.DeceasedList) == 0 {
		return
	}

	doc.SetFont(labelFontName, labelFontStyle, labelFontSize)
	doc.Image(backgroundImageKey, startX, startY, labelWidth, labelHeight, false, "PNG", -1, "")

	endX := labelWidth - (labelPaddingX * 2)
	endY := startY + labelHeight - labelPaddingBottom
	textColumnWidth := (endX / 5) - labelTextHorizontalSpacing
	cursorX := startX + labelPaddingX

	if len(tablet.DeceasedList) < 3 {
		cursorX += (textColumnWidth + labelTextHorizontalSpacing) * float64(3-len(tablet.DeceasedList))
	}

	slices.Reverse(tablet.DeceasedList)

	for _, d := range tablet.DeceasedList {
		cursorY := startY + labelPaddingTop
		name := strings.Join([]string{d.DharmaName, d.Name}, " ")

		doc.DrawVerticalText(fmt.Sprintf("%s 靈駕", name), cursorX, endY, labelCharVerticalHeight, labelSpaceVerticalHeightRatio, true)

		doc.DrawVerticalText("망", cursorX, cursorY, labelCharVerticalHeight, labelSpaceVerticalHeightRatio, false)
		cursorY += labelCharVerticalHeight * 1.5

		doc.DrawVerticalText(d.Relation, cursorX, cursorY, labelCharVerticalHeight, labelSpaceVerticalHeightRatio, false)
		cursorY += labelCharVerticalHeight * 5.25

		if 0 < len(d.ClanOrigin) {
			doc.DrawVerticalText(d.ClanOrigin, cursorX, cursorY, labelCharVerticalHeight, labelSpaceVerticalHeightRatio, false)
		}

		cursorX += textColumnWidth + labelTextHorizontalSpacing
	}

	cursorX += textColumnWidth + labelTextHorizontalSpacing
	doc.DrawVerticalText(fmt.Sprintf("%s 伏爲", tablet.PresentedBy), cursorX, endY, labelCharVerticalHeight, labelSpaceVerticalHeightRatio, true)
}
