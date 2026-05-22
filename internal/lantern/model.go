// Package lantern contains the self-contained 연등/인등 (Buddhist lantern)
// product. It supports four lantern kinds, each parsed from a separate sheet
// of the input XLSX file:
//
//	큰등    — same shape as 가족등 but rendered at a different size
//	가족등  — household with multiple family members (name + optional relation/dharma)
//	영가등  — household with one presenter (복위) and multiple deceased names
//	사업등  — household with one business name (사업체명) and multiple names
//
// Each kind is parsed by a dedicated function in parser.go and rendered by a
// dedicated function in render_<kind>.go. The dispatch lives in render.go.
package lantern

// FamilyMember is one person in a 가족등 / 큰등 household.
// Relation and DharmaName are both optional.
type FamilyMember struct {
	Name       string `json:"name"`
	Relation   string `json:"relation"`
	DharmaName string `json:"dharma_name"`
}

// FamilyHousehold is one entry of 가족등 or 큰등. The same struct is used for
// both kinds; the only difference is the rendering pass.
type FamilyHousehold struct {
	Address string         `json:"address"`
	Members []FamilyMember `json:"members"`
}

// SpiritHousehold is one entry of 영가등 — single presenter (복위) with a list
// of deceased names.
type SpiritHousehold struct {
	Address     string   `json:"address"`
	PresentedBy string   `json:"presented_by"` // 복위
	Names       []string `json:"names"`        // 망자 이름들
}

// BusinessHousehold is one entry of 사업등 — single business name with a list
// of associated names.
type BusinessHousehold struct {
	Address      string   `json:"address"`
	BusinessName string   `json:"business_name"` // 사업체명
	Names        []string `json:"names"`         // 관련자 이름들
}
