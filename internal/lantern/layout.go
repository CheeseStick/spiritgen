package lantern

import "spiritgen/internal/pdf"

// Page-level constants shared across every lantern kind. These are fixed by
// the A4 sheet and the embedded background image registration; per-kind
// variation happens inside the tablet, not at the page level.
const (
	a4PaperMarginX = 12.0
	a4PaperMarginY = 16.0

	// Used to identify the background image registered in the PDF.
	backgroundImageKey = "lantern_bg"

	// Used when the caller-provided title is empty.
	defaultTitle = "뉴질랜드 남국선사"
)

// PersonFontConfig holds the font sizes/styles used for one person's row in
// a member or name list. Each renderer file picks the per-row config it
// wants via a local fontSizeForCount(n) function.
type PersonFontConfig struct {
	nameFontSize        float64
	nameFontStyle       string
	dharmaNameFontSize  float64
	dharmaNameFontStyle string
}

// Layout holds only the basic tablet geometry that each lantern kind needs to
// place itself on the page and draw its background. All title/address/person
// styling (font sizes, gaps, single-vs-double column thresholds, etc.) lives
// as local constants in the corresponding render_<kind>.go file.
type Layout struct {
	// ───────── Tablet placement (mm) ─────────
	TabletSize    float64 // square edge length
	TabletMarginX float64 // gap between tablets horizontally
	TabletMarginY float64 // gap between rows of tablets

	// ───────── Inside-tablet padding (mm) ─────────
	PaddingX      float64
	PaddingTop    float64
	PaddingBottom float64

	// ───────── Background image area (mm) ─────────
	BackgroundSize float64

	// ───────── Cut guideline circle ─────────
	CircleLineWidth float64

	// ───────── Font family ─────────
	LabelFontName string // shared font family across this layout

	// Title config
	titleFontStyle    string
	titleFontSize     float64
	titleBottomMargin float64
}

// ───────── Layout helpers — geometry only ─────────

// lineHeightFor approximates mm of vertical line height for a font of the
// given point size, using a 0.5 mm/pt heuristic.
func (l *Layout) lineHeightFor(fontSize float64) float64 {
	return fontSize * 0.5
}

// contentX returns the left edge of the inside-tablet content area.
func (l *Layout) contentX(originX float64) float64 {
	return originX + l.PaddingX
}

// contentW returns the width of the inside-tablet content area.
func (l *Layout) contentW() float64 {
	return l.TabletSize - 2*l.PaddingX
}

// contentBottom returns the y-coordinate of the bottom of the content area
// for a tablet whose top is at originY.
func (l *Layout) contentBottom(originY float64) float64 {
	return originY + l.TabletSize - l.PaddingBottom
}

// tabletBottom returns the y-coordinate of the very bottom edge of the tablet.
func (l *Layout) tabletBottom(originY float64) float64 {
	return originY + l.TabletSize
}

// placeAndRender runs the standard "left-to-right, wrap row, wrap page"
// cursor loop using this layout's TabletSize / Margin values. The renderOne
// callback draws a single tablet at the given (originX, originY).
func (l *Layout) placeAndRender(doc *pdf.Doc, count int, renderOne func(originX, originY float64, i int)) {
	pageWidth, pageHeight := doc.GetPageSize()
	leftMargin, topMargin, rightMargin, bottomMargin := doc.GetMargins()
	pageRight := pageWidth - rightMargin
	pageBottom := pageHeight - bottomMargin

	cursorX := leftMargin
	cursorY := topMargin

	for i := 0; i < count; i++ {
		if cursorX+l.TabletSize > pageRight {
			cursorX = leftMargin
			cursorY += l.TabletSize + l.TabletMarginY
		}
		if cursorY+l.TabletSize > pageBottom {
			doc.AddPage()
			cursorX = leftMargin
			cursorY = topMargin
		}
		renderOne(cursorX, cursorY, i)
		cursorX += l.TabletSize + l.TabletMarginX
	}
}

// drawTabletBackground paints the aspect-fit background image then the
// circular cut-guideline outline on top of it.
func (l *Layout) drawTabletBackground(doc *pdf.Doc, originX, originY float64) {
	bgOffset := (l.TabletSize - l.BackgroundSize) / 2
	doc.DrawImageAspectFit(backgroundImageKey,
		originX+bgOffset, originY+bgOffset, l.BackgroundSize, l.BackgroundSize)

	prevLW := doc.GetLineWidth()
	doc.SetLineWidth(l.CircleLineWidth)
	doc.Circle(originX+l.TabletSize/2, originY+l.TabletSize/2, l.TabletSize/2, "D")
	doc.SetLineWidth(prevLW)
}

// drawCenteredLine sets the font then draws text horizontally centered
// within the tablet's content area. Returns the Y position immediately
// below the line (y + line-height). Callers manage any extra bottom margin.
func (l *Layout) drawCenteredLine(doc *pdf.Doc, originX, y float64, text, fontStyle string, fontSize float64) float64 {
	doc.SetFont(l.LabelFontName, fontStyle, fontSize)
	lh := l.lineHeightFor(fontSize)
	doc.DrawTextCentered(text, l.contentX(originX), l.contentW(), y+lh)
	return y + lh
}

// renderCenteredNameList draws each name centered horizontally on its own
// line. Used by 영가등 / 사업등 where each entry is just a name (no relation
// or dharma). Lines past tabletBottom are skipped.
func (l *Layout) renderCenteredNameList(doc *pdf.Doc, names []string, cfg PersonFontConfig, boxX, y, boxWidth, tabletBottom float64) {
	doc.SetFont(l.LabelFontName, cfg.nameFontStyle, cfg.nameFontSize)
	lh := l.lineHeightFor(cfg.nameFontSize)
	for _, name := range names {
		if y+lh > tabletBottom {
			break
		}
		doc.DrawTextCentered(name, boxX, boxWidth, y+lh)
		y += lh
	}
}
