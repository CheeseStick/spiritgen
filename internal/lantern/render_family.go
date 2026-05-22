package lantern

import (
	"fmt"

	"spiritgen/internal/pdf"
)

// ─────────────── 가족등 per-page constants ───────────────
// All title / address / member-row styling for 가족등 lives here.

const (
	// Address.
	familyAddressFontStyle = "M"
	familyAddressFontSize  = 12.0

	// Member rows.
	familyNamesSingleColumnMax = 5   // ≤N → 1 column; otherwise 2 columns
	familyNamesColumnGap       = 3.0 // padding between the two columns
	familyRelationNameGap      = 3.0 // horizontal gap between relation and name within a row
	familyDharmaGap            = 2.0 // gap between name and dharma name
	familyDharmaGapY           = 4.0 // vertical gap when single-person two-line layout is used
)

func familyFontSizeForCount(n int) PersonFontConfig {
	switch {
	case n == 1:
		return PersonFontConfig{nameFontSize: 32.0, nameFontStyle: "M", dharmaNameFontSize: 24.0, dharmaNameFontStyle: "L"} // 1명
	case n <= familyNamesSingleColumnMax:
		return PersonFontConfig{nameFontSize: 16.0, nameFontStyle: "M", dharmaNameFontSize: 15.0, dharmaNameFontStyle: "L"} // 2~5명
	default:
		return PersonFontConfig{nameFontSize: 12.0, nameFontStyle: "M", dharmaNameFontSize: 12.0, dharmaNameFontStyle: "L"} // 6명+
	}
}

// ─────────────── Renderers ───────────────

// renderFamilySection lays out every 가족등 tablet using familyLayout.
func renderFamilySection(doc *pdf.Doc, households []FamilyHousehold, title string) {
	familyLayout.placeAndRender(doc, len(households), func(originX, originY float64, i int) {
		renderFamilyTablet(doc, &familyLayout, households[i], originX, originY, title)
	})
}

// renderFamilyTablet draws one FamilyHousehold tablet using the supplied
// Layout (geometry) and the family-page constants above (styling).
func renderFamilyTablet(doc *pdf.Doc, l *Layout, h FamilyHousehold, originX, originY float64, title string) {
	l.drawTabletBackground(doc, originX, originY)

	contentX := l.contentX(originX)
	contentW := l.contentW()

	// Title — top-anchored, skipped when empty.
	currentY := originY + l.PaddingTop
	if title != "" {
		currentY = l.drawCenteredLine(doc, originX, currentY, title, l.titleFontStyle, l.titleFontSize) + l.titleBottomMargin
	}

	// Address — top-anchored, right under the title.
	currentY = l.drawCenteredLine(doc, originX, currentY, h.Address, familyAddressFontStyle, familyAddressFontSize)

	// Members — vertically centered, clamped so it never overlaps the address.
	cfg := familyFontSizeForCount(len(h.Members))
	availableH := l.contentBottom(originY) - currentY
	membersH := familyMembersHeight(l, h.Members, cfg)
	membersTop := currentY + (availableH-membersH)/2
	if membersTop < currentY {
		membersTop = currentY
	}

	tabletBottom := l.tabletBottom(originY)
	switch {
	case len(h.Members) == 1:
		renderSingleMember(doc, l, h.Members[0], cfg, contentX, contentW, membersTop)
	case len(h.Members) <= familyNamesSingleColumnMax:
		renderAlignedMemberColumn(doc, l, h.Members, cfg, contentX, membersTop, contentW, tabletBottom, true)
	default:
		half := (len(h.Members) + 1) / 2 // odd → left column carries the extra
		colW := (contentW - familyNamesColumnGap) / 2
		renderAlignedMemberColumn(doc, l, h.Members[:half], cfg, contentX, membersTop, colW, tabletBottom, false)
		renderAlignedMemberColumn(doc, l, h.Members[half:], cfg, contentX+colW+familyNamesColumnGap, membersTop, colW, tabletBottom, false)
	}
}

// familyMembersHeight returns the vertical space the members block will
// occupy. Address is laid out separately above the members.
func familyMembersHeight(l *Layout, members []FamilyMember, cfg PersonFontConfig) float64 {
	memberLH := l.lineHeightFor(cfg.nameFontSize)
	switch {
	case len(members) == 1:
		total := memberLH
		if members[0].DharmaName != "" {
			total += l.lineHeightFor(cfg.dharmaNameFontSize)
		}
		return total
	case len(members) <= familyNamesSingleColumnMax:
		return float64(len(members)) * memberLH
	default:
		rows := (len(members) + 1) / 2
		return float64(rows) * memberLH
	}
}

// renderSingleMember draws "relation name" (or just "name" when relation is
// empty) centered on one line and, if present, the dharma name centered on
// the next.
func renderSingleMember(doc *pdf.Doc, l *Layout, m FamilyMember, cfg PersonFontConfig, contentX, contentW, y float64) {
	doc.SetFont(l.LabelFontName, cfg.nameFontStyle, cfg.nameFontSize)
	nameLH := l.lineHeightFor(cfg.nameFontSize)
	text := m.Name
	if m.Relation != "" {
		text = fmt.Sprintf("%s %s", m.Relation, m.Name)
	}
	doc.DrawTextCentered(text, contentX, contentW, y+nameLH)
	y += nameLH + familyDharmaGapY

	if m.DharmaName != "" {
		doc.SetFont(l.LabelFontName, cfg.dharmaNameFontStyle, cfg.dharmaNameFontSize)
		dharmaLH := l.lineHeightFor(cfg.dharmaNameFontSize)
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
// When no member has a relation, the relation column is dropped entirely.
func renderAlignedMemberColumn(doc *pdf.Doc, l *Layout, members []FamilyMember, cfg PersonFontConfig, boxX, y, boxWidth, tabletBottom float64, showDharma bool) {
	memberLH := l.lineHeightFor(cfg.nameFontSize)

	// Pre-pass: measure each row's relation and name in the name font.
	relW := make([]float64, len(members))
	nameW := make([]float64, len(members))
	doc.SetFont(l.LabelFontName, cfg.nameFontStyle, cfg.nameFontSize)
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
		doc.SetFont(l.LabelFontName, cfg.dharmaNameFontStyle, cfg.dharmaNameFontSize)
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
		rightSideW += familyDharmaGap + maxDharmaW
	}
	totalW := maxRelW + rightSideW
	if showRelation {
		totalW += familyRelationNameGap
	}
	startX := boxX + (boxWidth-totalW)/2
	if startX < boxX {
		startX = boxX
	}

	var gutterX, nameStartX float64
	if showRelation {
		gutterX = startX + maxRelW
		nameStartX = gutterX + familyRelationNameGap
	} else {
		nameStartX = startX
	}

	for i, m := range members {
		if y+memberLH > tabletBottom {
			break
		}
		if showRelation && m.Relation != "" {
			doc.SetFont(l.LabelFontName, cfg.nameFontStyle, cfg.nameFontSize)
			doc.Text(gutterX-relW[i], y+memberLH, m.Relation)
		}
		doc.SetFont(l.LabelFontName, cfg.nameFontStyle, cfg.nameFontSize)
		doc.Text(nameStartX, y+memberLH, m.Name)
		if showDharma && m.DharmaName != "" {
			dharmaX := nameStartX + nameW[i] + familyDharmaGap
			if dharmaX+dharmaW[i] <= boxX+boxWidth {
				doc.SetFont(l.LabelFontName, cfg.dharmaNameFontStyle, cfg.dharmaNameFontSize)
				doc.Text(dharmaX, y+memberLH, m.DharmaName)
			}
		}
		y += memberLH
	}
}
