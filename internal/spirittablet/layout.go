package spirittablet

// Page-level margins (mm) for the A4 landscape sheet that hosts the tablets.
const (
	a4PaperMarginX = 12.0
	a4PaperMarginY = 16.0
)

// Per-tablet layout constants (mm) — dimensions and paddings of a single label.
const (
	labelWidth         = 58.0
	labelHeight        = 184.0
	labelMarginX       = 4.0
	labelPaddingX      = 14.0
	labelPaddingTop    = 32.0
	labelPaddingBottom = 28.0

	backgroundImageKey = "spirittablet_bg"
)

// Font configuration for vertical Korean text.
const (
	labelFontName  = "Noto Serif KR Black"
	labelFontStyle = "B"
	labelFontSize  = 14

	labelCharVerticalHeight       = 5.0
	labelSpaceVerticalHeightRatio = 1.5
	labelTextHorizontalSpacing    = 2.0
)
