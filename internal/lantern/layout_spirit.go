package lantern

// spiritLayout holds the basic geometry for 영가등 tablets. Styling lives in
// render_spirit.go.
var spiritLayout = Layout{
	TabletSize:    110.0,
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
