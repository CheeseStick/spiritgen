package model

import "fmt"

type SpiritTablet struct {
	PresentedBy  string     `json:"presented_by"`  // 복위
	DeceasedList []Deceased `json:"deceased_list"` // 망자 리스트
}

func (s SpiritTablet) String() string {
	return fmt.Sprintf("%s 복위\n %s", s.PresentedBy, s.DeceasedList)
}

func (s SpiritTablet) Split(maximum int) []SpiritTablet {
	if maximum <= 0 {
		maximum = 3 // fallback
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
