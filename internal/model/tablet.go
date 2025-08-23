package model

import "fmt"

type SpiritTablet struct {
	PresentedBy  string     `json:"presented_by"`  // 복위
	DeceasedList []Deceased `json:"deceased_list"` // 망자 리스트
}

func (s SpiritTablet) String() string {
	return fmt.Sprintf("%s 복위\n %s", s.PresentedBy, s.DeceasedList)
}
