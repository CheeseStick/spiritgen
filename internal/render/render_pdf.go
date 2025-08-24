package render

import (
	"codeberg.org/go-pdf/fpdf"
	"spiritgen/internal/model"
)

const (
	// 결과물 출력 A4 용지 마진 (mm)
	a4PaperMarginX = 12.0
	a4PaperMarginY = 16.0

	// 위패 용지 사이즈 (mm)
	spiritTabletLabelWidth              = 58.0
	spiritTabletLabelHeight             = 184.0
	spiritTabletBackgroundImageFilePath = "assets/images/background.jpg"

	// 위패 용지 간 마진 (mm)
	spiritTabletLabelMarginX = 8.0

	// 위패 용지 내 패딩 (mm)
	spiritTabletLabelPaddingX = 14.0
	spiritTabletLabelPaddingY = 18.0

	// 위패 용지 폰트 설정
	spiritTabletLabelFontName     = "Noto Serif KR Black"
	spiritTabletLabelFontStyle    = "B"
	spiritTabletLabelFontSize     = 14
	spiritTabletLabelFontFilePath = "assets/fonts/NotoSerifKR-Black.ttf"

	// 위패 용지 글자 간격 설정 (mm)
	spiritTabletLabelCharVerticalHeight    = 5.0
	spiritTabletLabelTextHorizontalSpacing = 2.0
)

func RenderLabelsAsPDF(labels []model.RenderedSpiritTabletLabel, outputPath string) error {
	pdf := fpdf.New("L", "mm", "A4", "")

	// Register assets
	pdf.AddUTF8Font(spiritTabletLabelFontName, spiritTabletLabelFontStyle, spiritTabletLabelFontFilePath)
	pdf.RegisterImage(spiritTabletBackgroundImageFilePath, "jpg")

	// Setup page
	pdf.SetMargins(a4PaperMarginX, a4PaperMarginY, -1)
	pdf.SetAutoPageBreak(true, a4PaperMarginY)
	pdf.AddPage()

	// Calculate maximum content width
	pageWidth, _ := pdf.GetPageSize()
	leftMargin, _, rightMargin, _ := pdf.GetMargins()
	pageContentSize := pageWidth - (leftMargin + rightMargin)

	// Cursor
	cursorX := a4PaperMarginX

	for _, label := range labels {
		// Calculate next content ending X
		nextContentEndX := cursorX + spiritTabletLabelWidth + spiritTabletLabelMarginX

		// If new content exceeds the maximum content width
		if pageContentSize < nextContentEndX {
			// Add new page and reset cursor
			pdf.AddPage()
			cursorX = a4PaperMarginX
		}

		// Render Label
		renderLabel(pdf, label, cursorX, a4PaperMarginY)

		// Update cursor
		cursorX += spiritTabletLabelWidth + spiritTabletLabelMarginX
	}

	return pdf.OutputFileAndClose(outputPath)
}

func renderLabel(pdf *fpdf.Fpdf, label model.RenderedSpiritTabletLabel, startX, startY float64) {
	if len(label.DeceasedLabels) == 0 {
		return
	}

	pdf.SetFont(spiritTabletLabelFontName, spiritTabletLabelFontStyle, spiritTabletLabelFontSize)
	pdf.Image(spiritTabletBackgroundImageFilePath, startX, startY, spiritTabletLabelWidth, spiritTabletLabelHeight, false, "JPG", -1, "")

	endX := spiritTabletLabelWidth - (spiritTabletLabelPaddingX * 2)
	endY := startY + (spiritTabletLabelHeight - (spiritTabletLabelPaddingY * 2))
	textColumnWidth := (endX / 5) - spiritTabletLabelTextHorizontalSpacing
	cursorX := startX + spiritTabletLabelPaddingX

	if len(label.DeceasedLabels) < 3 {
		cursorX += (textColumnWidth + spiritTabletLabelTextHorizontalSpacing) * float64(3-len(label.DeceasedLabels))
	}

	for _, deceased := range label.DeceasedLabels {
		renderVerticalText(pdf, deceased, cursorX, endY, spiritTabletLabelCharVerticalHeight, true)
		cursorX += textColumnWidth + spiritTabletLabelTextHorizontalSpacing
	}

	cursorX += textColumnWidth + spiritTabletLabelTextHorizontalSpacing
	renderVerticalText(pdf, label.PresenterLabel, cursorX, endY, spiritTabletLabelCharVerticalHeight, true)
}

func renderVerticalText(pdf *fpdf.Fpdf, text string, startX, startY, charVerticalHeight float64, reversed bool) {
	runes := []rune(text)
	length := len(runes)

	if length == 0 {
		return
	}

	if reversed {
		y := startY - charVerticalHeight

		for i := length; 0 < i; i-- {
			pdf.Text(startX, y, string(runes[i-1]))
			y -= charVerticalHeight
		}
	} else {
		y := startY + charVerticalHeight

		for _, r := range runes {
			pdf.Text(startX, y, string(r))
			y += charVerticalHeight
		}
	}
}
