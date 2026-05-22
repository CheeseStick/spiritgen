package lantern

import (
	"fmt"

	"spiritgen/internal/pdf"
)

// renderFamilySection lays out every 가족등 tablet, one per FamilyHousehold.
func renderFamilySection(doc *pdf.Doc, households []FamilyHousehold, title string) {
	placeAndRender(doc, len(households), func(originX, originY float64, i int) {
		renderFamilyTablet(doc, households[i], originX, originY, title)
	})
}

func renderFamilyTablet(doc *pdf.Doc, h FamilyHousehold, originX, originY float64, title string) {
	drawTabletBackground(doc, originX, originY)

	contentX := originX + tabletPaddingX
	contentW := tabletSize - 2*tabletPaddingX

	currentY := drawTitle(doc, originX, originY, title)
	currentY = drawAddress(doc, originX, currentY, h.Address)

	// Members block — vertically centered, clamped so it never overlaps the
	// address line above.
	cfg := fontSizeForCount(len(h.Members))
	contentBottom := originY + tabletSize - tabletPaddingBottom
	availableH := contentBottom - currentY
	membersH := familyMembersHeight(h.Members, cfg)
	membersTop := currentY + (availableH-membersH)/2
	if membersTop < currentY {
		membersTop = currentY
	}

	tabletBottom := originY + tabletSize
	switch {
	case len(h.Members) == 1:
		renderSingleMember(doc, h.Members[0], cfg, contentX, contentW, membersTop)
	case len(h.Members) <= namesSingleColumnMax:
		// Single column: enough room to show dharma names too.
		renderAlignedMemberColumn(doc, h.Members, cfg, contentX, membersTop, contentW, tabletBottom, true)
	default:
		// Two columns: narrow per-column width — skip dharma to keep things readable.
		half := (len(h.Members) + 1) / 2 // odd → left column carries the extra
		colW := (contentW - namesColumnGap) / 2
		renderAlignedMemberColumn(doc, h.Members[:half], cfg, contentX, membersTop, colW, tabletBottom, false)
		renderAlignedMemberColumn(doc, h.Members[half:], cfg, contentX+colW+namesColumnGap, membersTop, colW, tabletBottom, false)
	}
}

// familyMembersHeight returns the vertical space the members block will
// occupy. Address is laid out separately above the members.
func familyMembersHeight(members []FamilyMember, cfg PersonFontConfig) float64 {
	memberLH := lineHeightFor(cfg.nameFontSize)
	switch {
	case len(members) == 1:
		total := memberLH
		if members[0].DharmaName != "" {
			total += lineHeightFor(cfg.dharmaNameFontSize)
		}
		return total
	case len(members) <= namesSingleColumnMax:
		return float64(len(members)) * memberLH
	default:
		rows := (len(members) + 1) / 2
		return float64(rows) * memberLH
	}
}

// renderSingleMember draws "relation name" (or just "name" when relation is
// empty) centered on one line and, if present, the dharma name centered on
// the next.
func renderSingleMember(doc *pdf.Doc, m FamilyMember, cfg PersonFontConfig, contentX, contentW, y float64) {
	doc.SetFont(labelFontName, cfg.nameFontStyle, cfg.nameFontSize)
	nameLH := lineHeightFor(cfg.nameFontSize)
	text := m.Name
	if m.Relation != "" {
		text = fmt.Sprintf("%s %s", m.Relation, m.Name)
	}
	doc.DrawTextCentered(text, contentX, contentW, y+nameLH)
	y += nameLH + dharmaGapY

	if m.DharmaName != "" {
		doc.SetFont(labelFontName, cfg.dharmaNameFontStyle, cfg.dharmaNameFontSize)
		dharmaLH := lineHeightFor(cfg.dharmaNameFontSize)
		doc.DrawTextCentered(m.DharmaName, contentX, contentW, y+dharmaLH)
	}
}

// renderAlignedMemberColumn draws members as a "relation | name [dharma]"
// table within the column box (boxX, boxX+boxWidth). All rows share the same
// gutter so the relation→name gap is constant across rows. The widest row
// determines the gutter position, and the assembly is horizontally centered
// in the column box. Rows past tabletBottom are skipped.
//
// When showDharma is false, dharma names are omitted entirely.
// When no member in the column has a relation, the relation column is
// dropped (no gutter, no gap) and names are simply centered.
func renderAlignedMemberColumn(doc *pdf.Doc, members []FamilyMember, cfg PersonFontConfig, boxX, y, boxWidth, tabletBottom float64, showDharma bool) {
	memberLH := lineHeightFor(cfg.nameFontSize)

	// Pre-pass: measure each row's relation and name in the name font.
	relW := make([]float64, len(members))
	nameW := make([]float64, len(members))
	doc.SetFont(labelFontName, cfg.nameFontStyle, cfg.nameFontSize)
	var maxRelW, maxNameW float64
	for i, m := range members {
		relW[i] = doc.GetStringWidth(m.Relation)
		nameW[i] = doc.GetStringWidth(m.Name)
		if relW[i] > maxRelW {
			maxRelW = relW[i]
		}
		if nameW[i] > maxNameW {
			maxNameW = nameW[i]
		}
	}
	showRelation := maxRelW > 0

	// Optionally measure dharma names in the dharma font.
	dharmaW := make([]float64, len(members))
	var maxDharmaW float64
	if showDharma {
		doc.SetFont(labelFontName, cfg.dharmaNameFontStyle, cfg.dharmaNameFontSize)
		for i, m := range members {
			if m.DharmaName == "" {
				continue
			}
			dharmaW[i] = doc.GetStringWidth(m.DharmaName)
			if dharmaW[i] > maxDharmaW {
				maxDharmaW = dharmaW[i]
			}
		}
	}

	rightSideW := maxNameW
	if maxDharmaW > 0 {
		rightSideW += dharmaGap + maxDharmaW
	}
	totalW := maxRelW + rightSideW
	if showRelation {
		totalW += relationNameGap
	}
	startX := boxX + (boxWidth-totalW)/2
	if startX < boxX {
		startX = boxX
	}

	var gutterX, nameStartX float64
	if showRelation {
		gutterX = startX + maxRelW
		nameStartX = gutterX + relationNameGap
	} else {
		nameStartX = startX
	}

	for i, m := range members {
		if y+memberLH > tabletBottom {
			break
		}
		if showRelation && m.Relation != "" {
			doc.SetFont(labelFontName, cfg.nameFontStyle, cfg.nameFontSize)
			doc.Text(gutterX-relW[i], y+memberLH, m.Relation)
		}
		doc.SetFont(labelFontName, cfg.nameFontStyle, cfg.nameFontSize)
		doc.Text(nameStartX, y+memberLH, m.Name)
		if showDharma && m.DharmaName != "" {
			dharmaX := nameStartX + nameW[i] + dharmaGap
			if dharmaX+dharmaW[i] <= boxX+boxWidth {
				doc.SetFont(labelFontName, cfg.dharmaNameFontStyle, cfg.dharmaNameFontSize)
				doc.Text(dharmaX, y+memberLH, m.DharmaName)
			}
		}
		y += memberLH
	}
}
