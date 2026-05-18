// Package lantern contains the self-contained 연등/인등 (Buddhist lantern)
// product: domain model, XLSX schema, validation rules, PDF layout/rendering,
// and GUI tab. It depends only on horizontal infrastructure packages
// (xlsx, pdf, validation, gui) and the embedded assets.
package lantern

// Person represents a single member of a household on a lantern tablet.
type Person struct {
	Name       string `json:"name"`        // 이름
	Relation   string `json:"relation"`    // 관계
	DharmaName string `json:"dharma_name"` // 법명 (optional)
}

// Household groups a set of Persons sharing the same address.
// All persons of a household are rendered onto a single lantern tablet —
// there is no splitting like SpiritTablet.
type Household struct {
	Address string   `json:"address"`
	Persons []Person `json:"persons"`
}
