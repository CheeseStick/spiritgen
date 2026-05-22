package lantern

import "spiritgen/internal/pdf"

// renderBigSection lays out every 큰등 tablet. 큰등 has the same data shape as
// 가족등 (FamilyHousehold) but is intended to be drawn at a different size
// from 가족등. For now it just delegates to the family tablet renderer —
// resize / per-kind layout knobs should be added here when the design is
// finalized (e.g. larger tabletSize, larger fonts, single-tablet-per-page).
func renderBigSection(doc *pdf.Doc, households []FamilyHousehold, title string) {
	placeAndRender(doc, len(households), func(originX, originY float64, i int) {
		// TODO: customize for 큰등 — larger tablet, bigger font, single per page, etc.
		renderFamilyTablet(doc, households[i], originX, originY, title)
	})
}
