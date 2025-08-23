package model

import "fmt"

// Deceased - 망자 데이터
type Deceased struct {
	Name       string `json:"name"`        // 망자의 이름
	DharmaName string `json:"dharma_name"` // 망자의 법명 (Optional)
	ClanOrigin string `json:"clan_origin"` // 망자의 본관
	Relation   string `json:"relation"`    // 복위와의 관계
}

func (d Deceased) String() string {
	clanOrigin := ""
	dharmaName := ""

	if 0 < len(d.ClanOrigin) {
		clanOrigin = fmt.Sprintf("%s ", d.ClanOrigin)
	}

	if 0 < len(d.DharmaName) {
		clanOrigin = fmt.Sprintf("%s ", d.DharmaName)
	}

	return fmt.Sprintf("선망%s%s %s 영가", clanOrigin, dharmaName, d.Name)
}
