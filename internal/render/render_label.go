package render

import (
	"fmt"
	"spiritgen/internal/model"
)

func FromSpiritTablet(t model.SpiritTablet) []model.RenderedSpiritTabletLabel {
	const maxPerPage = 3
	var labels []model.RenderedSpiritTabletLabel

	for i := 0; i < len(t.DeceasedList); i += maxPerPage {
		end := i + maxPerPage
		if end > len(t.DeceasedList) {
			end = len(t.DeceasedList)
		}

		var deceasedLabels []string
		for _, d := range t.DeceasedList[i:end] {
			clanOrigin := ""
			name := ""

			if 0 < len(d.ClanOrigin) {
				clanOrigin = fmt.Sprintf("%s ", d.ClanOrigin)
			}

			if 0 < len(d.DharmaName) {
				name = fmt.Sprintf("%s ", d.DharmaName)
			}

			if 0 < len(d.Name) {
				name = fmt.Sprintf("%s%s ", name, d.Name)
			}

			deceasedLabels = append(deceasedLabels, fmt.Sprintf("망 %s %s%s靈駕", d.Relation, clanOrigin, name))
		}

		labels = append(labels, model.RenderedSpiritTabletLabel{
			PresenterLabel: fmt.Sprintf("%s 伏爲", t.PresentedBy),
			DeceasedLabels: deceasedLabels,
		})
	}

	return labels
}

func FromSpiritTablets(tablets []model.SpiritTablet) []model.RenderedSpiritTabletLabel {
	var all []model.RenderedSpiritTabletLabel
	for _, t := range tablets {
		all = append(all, FromSpiritTablet(t)...)
	}
	return all
}
