package lantern

import "spiritgen/internal/pdf"

// ─────────────── 영가등 per-page constants ───────────────

const (
	spiritAddressFontStyle = "M"
	spiritAddressFontSize  = 12.0

	spiritPresenterFontStyle = "B"
	spiritPresenterFontSize  = 12.0
)

func spiritFontSizeForCount(n int) PersonFontConfig {
	switch {
	case n == 1:
		return PersonFontConfig{nameFontSize: 16.0, nameFontStyle: "M", dharmaNameFontSize: 16.0, dharmaNameFontStyle: "L"}
	case n <= 4:
		return PersonFontConfig{nameFontSize: 14.0, nameFontStyle: "M", dharmaNameFontSize: 14.0, dharmaNameFontStyle: "L"}
	default:
		return PersonFontConfig{nameFontSize: 12.0, nameFontStyle: "M", dharmaNameFontSize: 12.0, dharmaNameFontStyle: "L"}
	}
}

// ─────────────── Renderers ───────────────

// renderSpiritSection lays out every 영가등 tablet using spiritLayout.
func renderSpiritSection(doc *pdf.Doc, households []SpiritHousehold, title string) {
	spiritLayout.placeAndRender(doc, len(households), func(originX, originY float64, i int) {
		renderSpiritTablet(doc, &spiritLayout, households[i], originX, originY, title)
	})
}

// renderSpiritTablet is a first-pass layout for 영가등 — address + presenter
// (복위) on top, then a vertically-centered list of names.
func renderSpiritTablet(doc *pdf.Doc, l *Layout, h SpiritHousehold, originX, originY float64, title string) {
	l.drawTabletBackground(doc, originX, originY)

	contentX := l.contentX(originX)
	contentW := l.contentW()

	currentY := originY + l.PaddingTop
	if title != "" {
		currentY = l.drawCenteredLine(doc, originX, currentY, title, l.titleFontStyle, l.titleFontSize) + l.titleBottomMargin
	}

	// Address.
	currentY = l.drawCenteredLine(doc, originX, currentY, h.Address, spiritAddressFontStyle, spiritAddressFontSize)

	// Presenter (복위).
	currentY = l.drawCenteredLine(doc, originX, currentY, h.PresentedBy+" 복위", spiritPresenterFontStyle, spiritPresenterFontSize)

	// Names — vertically centered, clamped so it never overlaps the lines above.
	cfg := spiritFontSizeForCount(len(h.Names))
	availableH := l.contentBottom(originY) - currentY
	namesH := float64(len(h.Names)) * l.lineHeightFor(cfg.nameFontSize)
	namesTop := currentY + (availableH-namesH)/2
	if namesTop < currentY {
		namesTop = currentY
	}

	l.renderCenteredNameList(doc, h.Names, cfg, contentX, namesTop, contentW, l.tabletBottom(originY))
}
