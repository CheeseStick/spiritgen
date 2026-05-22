package lantern

import "spiritgen/internal/pdf"

const (
	businessAddressFontStyle = "M"
	businessAddressFontSize  = 12.0

	businessNameLineFontStyle    = "B" // for the 사업체명 line
	businessNameLineFontSize     = 24.0
	businessNameLineBottomMargin = 4.0
)

func businessFontSizeForCount(n int) PersonFontConfig {
	switch {
	case n == 1:
		return PersonFontConfig{nameFontSize: 14.0, nameFontStyle: "M", dharmaNameFontSize: 24.0, dharmaNameFontStyle: "L"}
	case n <= 5:
		return PersonFontConfig{nameFontSize: 12.0, nameFontStyle: "M", dharmaNameFontSize: 15.0, dharmaNameFontStyle: "L"}
	default:
		return PersonFontConfig{nameFontSize: 10.0, nameFontStyle: "M", dharmaNameFontSize: 12.0, dharmaNameFontStyle: "L"}
	}
}

// ─────────────── Renderers ───────────────

// renderBusinessSection lays out every 사업등 tablet using businessLayout.
func renderBusinessSection(doc *pdf.Doc, households []BusinessHousehold, title string) {
	businessLayout.placeAndRender(doc, len(households), func(originX, originY float64, i int) {
		renderBusinessTablet(doc, &businessLayout, households[i], originX, originY, title)
	})
}

// renderBusinessTablet is a first-pass layout for 사업등 — address + business
// name on top, then a vertically-centered list of associated names.
func renderBusinessTablet(doc *pdf.Doc, l *Layout, h BusinessHousehold, originX, originY float64, title string) {
	l.drawTabletBackground(doc, originX, originY)

	contentX := l.contentX(originX)
	contentW := l.contentW()

	currentY := originY + l.PaddingTop
	if title != "" {
		currentY = l.drawCenteredLine(doc, originX, currentY, title, l.titleFontStyle, l.titleFontSize) + l.titleBottomMargin
	}

	// Address.
	currentY = l.drawCenteredLine(doc, originX, currentY, h.Address, businessAddressFontStyle, businessAddressFontSize)

	// 사업체명
	currentY = l.drawCenteredLine(doc, originX, currentY, h.BusinessName, businessNameLineFontStyle, businessNameLineFontSize) + businessNameLineBottomMargin

	// Names — vertically centered, clamped so it never overlaps the lines above.
	cfg := businessFontSizeForCount(len(h.Names))
	availableH := l.contentBottom(originY) - currentY
	namesH := float64(len(h.Names)) * l.lineHeightFor(cfg.nameFontSize)
	namesTop := currentY + (availableH-namesH)/2
	if namesTop < currentY {
		namesTop = currentY
	}

	l.renderCenteredNameList(doc, h.Names, cfg, contentX, namesTop, contentW, l.tabletBottom(originY))
}
