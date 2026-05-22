package lantern

import (
	"spiritgen/internal/pdf"
)

// renderSpiritSection lays out every 영가등 tablet, one per SpiritHousehold.
func renderSpiritSection(doc *pdf.Doc, households []SpiritHousehold, title string) {
	placeAndRender(doc, len(households), func(originX, originY float64, i int) {
		renderSpiritTablet(doc, households[i], originX, originY, title)
	})
}

// renderSpiritTablet is a first-pass layout for 영가등 — it currently just
// puts the address, the presenter (복위), and a vertically-centered list of
// names on the tablet. The actual visual design should be filled in here.
func renderSpiritTablet(doc *pdf.Doc, h SpiritHousehold, originX, originY float64, title string) {
	drawTabletBackground(doc, originX, originY)

	contentX := originX + tabletPaddingX
	contentW := tabletSize - 2*tabletPaddingX

	currentY := drawTitle(doc, originX, originY, title)
	currentY = drawAddress(doc, originX, currentY, h.Address)

	// 복위 line — same font family as address, slightly emphasized.
	// TODO: tune sizing/style for 영가등 specifically.
	doc.SetFont(labelFontName, addressFontStyle, addressFontSize)
	presenterLH := lineHeightFor(addressFontSize)
	doc.DrawTextCentered(h.PresentedBy+" 伏爲", contentX, contentW, currentY+presenterLH)
	currentY += presenterLH

	// Names — vertically centered list, clamped to not overlap the lines above.
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
