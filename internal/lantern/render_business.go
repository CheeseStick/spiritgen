package lantern

import (
	"spiritgen/internal/pdf"
)

// renderBusinessSection lays out every 사업등 tablet, one per BusinessHousehold.
func renderBusinessSection(doc *pdf.Doc, households []BusinessHousehold, title string) {
	placeAndRender(doc, len(households), func(originX, originY float64, i int) {
		renderBusinessTablet(doc, households[i], originX, originY, title)
	})
}

// renderBusinessTablet is a first-pass layout for 사업등 — address + business
// name on top, then a vertically-centered list of associated names. The
// actual visual design should be filled in here.
func renderBusinessTablet(doc *pdf.Doc, h BusinessHousehold, originX, originY float64, title string) {
	drawTabletBackground(doc, originX, originY)

	contentX := originX + tabletPaddingX
	contentW := tabletSize - 2*tabletPaddingX

	currentY := drawTitle(doc, originX, originY, title)
	currentY = drawAddress(doc, originX, currentY, h.Address)

	// 사업체명 line.
	// TODO: tune sizing/style for 사업등 specifically.
	doc.SetFont(labelFontName, addressFontStyle, addressFontSize)
	bizLH := lineHeightFor(addressFontSize)
	doc.DrawTextCentered(h.BusinessName, contentX, contentW, currentY+bizLH)
	currentY += bizLH

	// Names — vertically centered list, clamped to not overlap above.
	cfg := fontSizeForCount(len(h.Names))
	contentBottom := originY + tabletSize - tabletPaddingBottom
	availableH := contentBottom - currentY
	namesH := float64(len(h.Names)) * lineHeightFor(cfg.nameFontSize)
	namesTop := currentY + (availableH-namesH)/2
	if namesTop < currentY {
		namesTop = currentY
	}

	tabletBottom := originY + tabletSize
	renderCenteredNameList(doc, h.Names, cfg, contentX, namesTop, contentW, tabletBottom)
}
