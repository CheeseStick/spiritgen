package lantern

// bigLayout holds the basic geometry for 큰등 tablets. 큰등 has the same data
// shape as 가족등 (FamilyHousehold) but a different physical size. Tune the
// geometry here; styling (fonts, gaps) lives in render_big.go.
var bigLayout = Layout{
	TabletSize:    110.0, // TODO: bump for 큰등's larger physical size
	TabletMarginX: 2.0,
	TabletMarginY: 2.0,

	PaddingX:      16.0,
	PaddingTop:    16.0,
	PaddingBottom: 34.0,

	BackgroundSize:  108.0,
	CircleLineWidth: 0.1,

	LabelFontName: "Noto Serif KR Black",

	titleFontStyle:    "B",
	titleFontSize:     16.0,
	titleBottomMargin: 4.0,
}
