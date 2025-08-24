package render

import (
	"codeberg.org/go-pdf/fpdf"
	"fmt"
	"spiritgen/internal/model"
	"strings"
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
	spiritTabletLabelMarginX = 4.0

	// 위패 용지 내 패딩 (mm)
	spiritTabletLabelPaddingX      = 14.0
	spiritTabletLabelPaddingTop    = 34.0
	spiritTabletLabelPaddingBottom = 36.0

	// 위패 용지 폰트 설정
	spiritTabletLabelFontName     = "Noto Serif KR Black"
	spiritTabletLabelFontStyle    = "B"
	spiritTabletLabelFontSize     = 14
	spiritTabletLabelFontFilePath = "assets/fonts/NotoSerifKR-Black.ttf"

	// 위패 용지 글자 간격 설정 (mm)
	spiritTabletLabelCharVerticalHeight    = 5.0
	spiritTabletLabelTextHorizontalSpacing = 2.0
)

func FromSpiritTablets(tablets []model.SpiritTablet, outputPath string) error {
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

	for _, s := range tablets {
		// Calculate next content ending X
		nextContentEndX := cursorX + spiritTabletLabelWidth + spiritTabletLabelMarginX

		// If new content exceeds the maximum content width
		if pageContentSize < nextContentEndX {
			// Add new page and reset cursor
			pdf.AddPage()
			cursorX = a4PaperMarginX
		}

		// Render tablet
		renderTablet(pdf, s, cursorX, a4PaperMarginY)

		// Update cursor
		cursorX += spiritTabletLabelWidth + spiritTabletLabelMarginX
	}

	return pdf.OutputFileAndClose(outputPath)
}

func renderTablet(pdf *fpdf.Fpdf, tablet model.SpiritTablet, startX, startY float64) {
	if len(tablet.DeceasedList) == 0 {
		return
	}

	pdf.SetFont(spiritTabletLabelFontName, spiritTabletLabelFontStyle, spiritTabletLabelFontSize)
	pdf.Image(spiritTabletBackgroundImageFilePath, startX, startY, spiritTabletLabelWidth, spiritTabletLabelHeight, false, "JPG", -1, "")

	endX := spiritTabletLabelWidth - (spiritTabletLabelPaddingX * 2)
	endY := startY + spiritTabletLabelHeight - spiritTabletLabelPaddingBottom
	textColumnWidth := (endX / 5) - spiritTabletLabelTextHorizontalSpacing
	cursorX := startX + spiritTabletLabelPaddingX

	if len(tablet.DeceasedList) < 3 {
		cursorX += (textColumnWidth + spiritTabletLabelTextHorizontalSpacing) * float64(3-len(tablet.DeceasedList))
	}

	for _, d := range tablet.DeceasedList {
		cursorY := startY + spiritTabletLabelPaddingTop
		name := strings.Join([]string{d.DharmaName, d.Name}, " ")

		renderVerticalText(pdf, fmt.Sprintf("%s 靈駕", name), cursorX, endY, spiritTabletLabelCharVerticalHeight, true)

		renderVerticalText(pdf, "망", cursorX, cursorY, spiritTabletLabelCharVerticalHeight, false)
		cursorY += spiritTabletLabelCharVerticalHeight * 1.5

		renderVerticalText(pdf, d.Relation, cursorX, cursorY, spiritTabletLabelCharVerticalHeight, false)
		cursorY += spiritTabletLabelCharVerticalHeight * 4.5

		if 0 < len(d.ClanOrigin) {
			renderVerticalText(pdf, d.ClanOrigin, cursorX, cursorY, spiritTabletLabelCharVerticalHeight, false)
		}

		cursorX += textColumnWidth + spiritTabletLabelTextHorizontalSpacing
	}

	cursorX += textColumnWidth + spiritTabletLabelTextHorizontalSpacing
	renderVerticalText(pdf, fmt.Sprintf("%s 伏爲", tablet.PresentedBy), cursorX, endY, spiritTabletLabelCharVerticalHeight, true)
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
