package render_test

import (
	"spiritgen/internal/model"
	"spiritgen/internal/render"
	"testing"
)

func TestFromSpiritTablet(t *testing.T) {
	tablet := model.SpiritTablet{
		PresentedBy: "김누구",
		DeceasedList: []model.Deceased{
			{Name: "누군가", DharmaName: "법명", ClanOrigin: "본관 후인", Relation: "관계"},
			{Name: "누군나", DharmaName: "법명", ClanOrigin: "본관 유인", Relation: "관계"},
			{Name: "누군다", DharmaName: "", ClanOrigin: "본관 후인", Relation: "관계"},
			{Name: "누군라", DharmaName: "법명", ClanOrigin: "", Relation: "관계"},
			{Name: "", DharmaName: "법명", ClanOrigin: "본관 유인", Relation: "관계"},
		},
	}

	labels := render.FromSpiritTablet(tablet)

	if len(labels) != 2 {
		t.Errorf("Expected 2 pages, got %d", len(labels))
	}

	ExpectedPresenterLabel := "김누구 伏爲"
	if labels[0].PresenterLabel != ExpectedPresenterLabel || labels[1].PresenterLabel != ExpectedPresenterLabel {
		t.Errorf("Expected presenter label %s but got: %s", ExpectedPresenterLabel, labels[0].PresenterLabel)
	}

	if labels[0].DeceasedLabels[0] != "망 관계 본관 후인 법명 누군가 靈駕" {
		t.Errorf("Expected DeceasedLabel differs: %s", labels[0].DeceasedLabels[0])
	}

	if labels[0].DeceasedLabels[2] != "망 관계 본관 후인 누군다 靈駕" {
		t.Errorf("Expected DeceasedLabel differs: %s", labels[0].DeceasedLabels[0])
	}

	if labels[1].DeceasedLabels[0] != "망 관계 법명 누군라 靈駕" {
		t.Errorf("Expected DeceasedLabel differs: %s", labels[0].DeceasedLabels[0])
	}

	if labels[1].DeceasedLabels[1] != "망 관계 본관 유인 법명 靈駕" {
		t.Errorf("Expected DeceasedLabel differs: %s", labels[0].DeceasedLabels[0])
	}
}
