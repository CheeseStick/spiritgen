package lantern

type PersonFontConfig struct {
	nameFontSize        float64
	nameFontStyle       string
	dharmaNameFontSize  float64
	dharmaNameFontStyle string
}

// Page-level margins (mm) for the A4 landscape sheet that hosts the tablets.
const (
	a4PaperMarginX = 12.0
	a4PaperMarginY = 16.0
)

// Per-tablet layout constants (mm).
const (
	tabletSize    = 110.0 // 11 cm square area per tablet
	tabletMarginX = 2.0
	tabletMarginY = 2.0

	tabletPaddingX      = 16.0 // horizontal padding inside the tablet square for text content
	tabletPaddingTop    = 16.0 // vertical padding inside the tablet square for text content
	tabletPaddingBottom = 34.0

	// Cut guideline: 1 px hairline circle filling the tablet square.
	circleLineWidth = 0.1

	// Background image area: nearly fills the tablet square (1mm inset)
	// and is rendered with aspect-ratio preserved.
	backgroundSize     = 108.0
	backgroundImageKey = "lantern_bg"
)

// Text layout.
const (
	labelFontName = "Noto Serif KR Black"

	// Title (centered at top; defaults to defaultTitle when caller passes "")
	titleFontSize     = 16.0
	titleFontStyle    = "B"
	titleBottomMargin = 4.0 // mm gap between title and body text
	defaultTitle      = "뉴질랜드 남국선사"

	// Address
	addressFontSize  = 10.0
	addressFontStyle = "M"

	// Persons
	namesSingleColumnMax = 5   // up to 5 persons → 1 column; 6+ → 2 columns
	namesColumnGap       = 3.0 // padding between the two columns
	relationNameGap      = 3.0 // horizontal gap between the relation column and the name column inside a row
	dharmaGap            = 2.0 // gap between "relation name" and dharma name
	dharmaGapY           = 4.0 // vertical gap between relation name and dharma name when both are present
)

// fontSizeForCount returns the PersonFontConfig for a given number of persons, which determines the font sizes and styles
func fontSizeForCount(n int) PersonFontConfig {
	switch {
	case n == 1:
		return PersonFontConfig{nameFontSize: 32.0, nameFontStyle: "M", dharmaNameFontSize: 24.0, dharmaNameFontStyle: "L"} // 1명
	case n <= namesSingleColumnMax:
		return PersonFontConfig{nameFontSize: 16.0, nameFontStyle: "M", dharmaNameFontSize: 14.0, dharmaNameFontStyle: "L"} // 2 - 5명
	default:
		return PersonFontConfig{nameFontSize: 12.0, nameFontStyle: "M", dharmaNameFontSize: 10.0, dharmaNameFontStyle: "L"} // 6명 이상
	}
}

// lineHeightFor approximates the vertical line height (mm) for a given font
// size in points using a 0.5 mm/pt heuristic that matches the existing layout.
func lineHeightFor(fontSize float64) float64 {
	return fontSize * 0.5
}
