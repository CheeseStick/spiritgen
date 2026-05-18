// Package spirittablet contains the self-contained 영가위패 (Korean spirit tablet)
// product: domain model, XLSX schema, validation rules, PDF layout/rendering,
// and GUI tab. It depends only on horizontal infrastructure packages
// (xlsx, pdf, validation, gui) and the embedded assets.
package spirittablet

import "fmt"

// Deceased represents a 망자 entry on a spirit tablet.
type Deceased struct {
	Name       string `json:"name"`
	DharmaName string `json:"dharma_name"`
	ClanOrigin string `json:"clan_origin"`
	Relation   string `json:"relation"`
}

// SpiritTablet groups one or more Deceased entries under a single 복위 (PresentedBy).
type SpiritTablet struct {
	PresentedBy  string     `json:"presented_by"`
	DeceasedList []Deceased `json:"deceased_list"`
}

func (s SpiritTablet) String() string {
	return fmt.Sprintf("%s 복위\n %s", s.PresentedBy, s.DeceasedList)
}

// Split chunks a SpiritTablet whose DeceasedList exceeds maximum into multiple
// SpiritTablets sharing the same PresentedBy.
func (s SpiritTablet) Split(maximum int) []SpiritTablet {
	if maximum <= 0 {
		maximum = 3
	}

	var result []SpiritTablet
	total := len(s.DeceasedList)
	for i := 0; i < total; i += maximum {
		end := i + maximum
		if end > total {
			end = total
		}
		chunk := SpiritTablet{
			PresentedBy:  s.PresentedBy,
			DeceasedList: s.DeceasedList[i:end],
		}
		if len(chunk.DeceasedList) != 0 {
			result = append(result, chunk)
		}
	}
	return result
}
