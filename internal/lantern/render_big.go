package lantern

import "spiritgen/internal/pdf"

// renderBigSection lays out every 큰등 tablet using bigLayout (geometry) and
// — for now — the family-page styling constants from render_family.go.
//
// TODO: once 큰등's visual design diverges (larger fonts, single tablet per
// page, alternate ornamentation, etc.), copy the renderFamilyTablet body
// here under renderBigTablet and replace familyXxx constants with bigXxx
// equivalents defined locally in this file.
func renderBigSection(doc *pdf.Doc, households []FamilyHousehold, title string) {
	bigLayout.placeAndRender(doc, len(households), func(originX, originY float64, i int) {
		renderFamilyTablet(doc, &bigLayout, households[i], originX, originY, title)
	})
}
