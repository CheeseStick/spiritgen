package lantern

import (
	"fmt"

	"spiritgen/assets"
	"spiritgen/internal/pdf"
)

// RenderPDF writes the given households as a multi-page lantern tablet PDF at
// outputPath, using bgImage (PNG bytes) as the per-tablet background. title is
// rendered at the top center of every tablet — if empty, the title line is
// skipped.
func RenderPDF(households []Household, outputPath string, bgImage []byte, title string) error {
	if title == "" {
		title = defaultTitle
	}

	doc := pdf.NewA4Landscape(a4PaperMarginX, a4PaperMarginY)
	// We only ship one TTF (NotoSerifKR-Black). Register it under both "B" and
	// "" so SetFont calls with either style resolve to the same bytes.
	doc.RegisterUTF8Font(labelFontName, "B", assets.FontNotoSerifKR)
	doc.RegisterUTF8Font(labelFontName, "M", assets.FontNotoSerifKRMedium)
	doc.RegisterUTF8Font(labelFontName, "L", assets.FontNotoSerifKRLight)
	doc.RegisterPNGImage(backgroundImageKey, bgImage)
	doc.AddPage()

	pageWidth, pageHeight := doc.GetPageSize()
	leftMargin, topMargin, rightMargin, bottomMargin := doc.GetMargins()
	pageRight := pageWidth - rightMargin
	pageBottom := pageHeight - bottomMargin

	cursorX := leftMargin
	cursorY := topMargin

	for _, h := range households {
		if cursorX+tabletSize > pageRight {
			cursorX = leftMargin
			cursorY += tabletSize + tabletMarginY
		}
		if cursorY+tabletSize > pageBottom {
			doc.AddPage()
			cursorX = leftMargin
			cursorY = topMargin
		}

		renderTablet(doc, h, cursorX, cursorY, title)
		cursorX += tabletSize + tabletMarginX
	}

	return doc.OutputFileAndClose(outputPath)
}

func renderTablet(doc *pdf.Doc, h Household, originX, originY float64, title string) {
	// 1) Background image
	bgOffset := (tabletSize - backgroundSize) / 2
	doc.DrawImageAspectFit(backgroundImageKey,
		originX+bgOffset, originY+bgOffset, backgroundSize, backgroundSize)

	// 2) Circle cut-guideline outline (drawn on top of background image)
	prevLW := doc.GetLineWidth()
	doc.SetLineWidth(circleLineWidth)
	doc.Circle(originX+tabletSize/2, originY+tabletSize/2, tabletSize/2, "D")
	doc.SetLineWidth(prevLW)

	// Content area (inside the horizontal padding).
	contentX := originX + tabletPaddingX
	contentW := tabletSize - 2*tabletPaddingX

	// 3) Title
	currentY := originY + tabletPaddingTop
	if title != "" {
		doc.SetFont(labelFontName, titleFontStyle, titleFontSize)
		titleLH := lineHeightFor(titleFontSize)
		doc.DrawTextCentered(title, contentX, contentW, currentY+titleLH)
		currentY += titleLH + titleBottomMargin
	}

	// 4) Address
	doc.SetFont(labelFontName, addressFontStyle, addressFontSize)
	addressLH := lineHeightFor(addressFontSize)
	doc.DrawTextCentered(h.Address, contentX, contentW, currentY+addressLH)
	currentY += addressLH

	// 5) Persons
	personFontConfig := fontSizeForCount(len(h.Persons))
	contentBottom := originY + tabletSize - tabletPaddingBottom
	availableH := contentBottom - currentY
	personsH := personsContentHeight(h, personFontConfig)
	personsTop := currentY + (availableH-personsH)/2
	if personsTop < currentY {
		personsTop = currentY
	}

	tabletBottom := originY + tabletSize
	switch {
	case len(h.Persons) == 1:
		renderSinglePerson(doc, h.Persons[0], personFontConfig, contentX, contentW, personsTop)
	case len(h.Persons) <= namesSingleColumnMax:
		// Single column has plenty of room — render dharma names too.
		renderAlignedColumn(doc, h.Persons, personFontConfig, contentX, personsTop, contentW, tabletBottom, true)
	default:
		// Two columns: narrow per-column width — skip dharma to keep relation/
		// name comfortably spaced.
		half := (len(h.Persons) + 1) / 2 // odd → left column carries the extra
		colW := (contentW - namesColumnGap) / 2
		renderAlignedColumn(doc, h.Persons[:half], personFontConfig, contentX, personsTop, colW, tabletBottom, false)
		renderAlignedColumn(doc, h.Persons[half:], personFontConfig, contentX+colW+namesColumnGap, personsTop, colW, tabletBottom, false)
	}
}

// personsContentHeight returns the vertical space the persons block will
// occupy for the given household. Used to position the block vertically.
// Address is not included — it is laid out separately above the persons.
func personsContentHeight(h Household, cfg PersonFontConfig) float64 {
	personLH := lineHeightFor(cfg.nameFontSize)
	switch {
	case len(h.Persons) == 1:
		total := personLH
		if h.Persons[0].DharmaName != "" {
			total += lineHeightFor(cfg.dharmaNameFontSize)
		}
		return total
	case len(h.Persons) <= namesSingleColumnMax:
		return float64(len(h.Persons)) * personLH
	default:
		rows := (len(h.Persons) + 1) / 2
		return float64(rows) * personLH
	}
}

// renderSinglePerson draws "relation name" centered on one line and, if
// present, the dharma name centered on the next, all within the content box
// (contentX, contentX+contentW).
func renderSinglePerson(doc *pdf.Doc, p Person, cfg PersonFontConfig, contentX, contentW, y float64) {
	doc.SetFont(labelFontName, cfg.nameFontStyle, cfg.nameFontSize)
	nameLH := lineHeightFor(cfg.nameFontSize)
	text := fmt.Sprintf("%s %s", p.Relation, p.Name)
	doc.DrawTextCentered(text, contentX, contentW, y+nameLH)
	y += nameLH + dharmaGapY

	if p.DharmaName != "" {
		doc.SetFont(labelFontName, cfg.dharmaNameFontStyle, cfg.dharmaNameFontSize)
		dharmaLH := lineHeightFor(cfg.dharmaNameFontSize)
		doc.DrawTextCentered(p.DharmaName, contentX, contentW, y+dharmaLH)
	}
}

// renderAlignedColumn draws persons as a "relation | name [dharma]" table
// within the column box (boxX, boxX+boxWidth). All rows share the same gutter
// so that the relation→name gap is constant across rows. The widest row
// determines the gutter position, and the assembly is horizontally centered
// in the column box. Rows past tabletBottom are skipped.
//
// When showDharma is false, dharma names are omitted entirely (no measurement,
// no rendering) — useful for narrow 2-column layouts where space is tight.
func renderAlignedColumn(doc *pdf.Doc, persons []Person, cfg PersonFontConfig, boxX, y, boxWidth, tabletBottom float64, showDharma bool) {
	personLH := lineHeightFor(cfg.nameFontSize)

	// Pre-pass: measure each row's relation and name in the name font.
	relW := make([]float64, len(persons))
	nameW := make([]float64, len(persons))
	doc.SetFont(labelFontName, cfg.nameFontStyle, cfg.nameFontSize)
	var maxRelW, maxNameW float64
	for i, p := range persons {
		relW[i] = doc.GetStringWidth(p.Relation)
		nameW[i] = doc.GetStringWidth(p.Name)
		if relW[i] > maxRelW {
			maxRelW = relW[i]
		}
		if nameW[i] > maxNameW {
			maxNameW = nameW[i]
		}
	}

	// Optionally measure dharma names in the dharma font.
	dharmaW := make([]float64, len(persons))
	var maxDharmaW float64
	if showDharma {
		doc.SetFont(labelFontName, cfg.dharmaNameFontStyle, cfg.dharmaNameFontSize)
		for i, p := range persons {
			if p.DharmaName == "" {
				continue
			}
			dharmaW[i] = doc.GetStringWidth(p.DharmaName)
			if dharmaW[i] > maxDharmaW {
				maxDharmaW = dharmaW[i]
			}
		}
	}

	// Center the widest-row assembly within the column box.
	rightSideW := maxNameW
	if maxDharmaW > 0 {
		rightSideW += dharmaGap + maxDharmaW
	}
	totalW := maxRelW + relationNameGap + rightSideW
	startX := boxX + (boxWidth-totalW)/2
	if startX < boxX {
		startX = boxX // clamp if totalW > boxWidth
	}
	gutterX := startX + maxRelW             // right edge of the relation column
	nameStartX := gutterX + relationNameGap // left edge of the name column

	for i, p := range persons {
		if y+personLH > tabletBottom {
			break
		}
		// Relation — right-aligned to gutterX
		doc.SetFont(labelFontName, cfg.nameFontStyle, cfg.nameFontSize)
		doc.Text(gutterX-relW[i], y+personLH, p.Relation)
		// Name — left-aligned at nameStartX
		doc.Text(nameStartX, y+personLH, p.Name)
		// Dharma — beside the name (only when enabled and it fits)
		if showDharma && p.DharmaName != "" {
			dharmaX := nameStartX + nameW[i] + dharmaGap
			if dharmaX+dharmaW[i] <= boxX+boxWidth {
				doc.SetFont(labelFontName, cfg.dharmaNameFontStyle, cfg.dharmaNameFontSize)
				doc.Text(dharmaX, y+personLH, p.DharmaName)
			}
		}
		y += personLH
	}
}
