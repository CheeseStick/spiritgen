package pdf

// DrawVerticalText writes the given text vertically starting at (startX, startY).
// When reversed is true the text is drawn upward from startY (i.e. the last rune
// is drawn at startY and earlier runes are placed above it). Otherwise, text is
// drawn downward.
//
// charHeight is the vertical spacing for a normal character; spaceRatio divides
// charHeight when the rune is a space, so spaces consume less vertical room.
func (d *Doc) DrawVerticalText(text string, startX, startY, charHeight, spaceRatio float64, reversed bool) {
	runes := []rune(text)
	length := len(runes)
	if length == 0 {
		return
	}

	if reversed {
		y := startY - charHeight
		for i := length; 0 < i; i-- {
			d.Text(startX, y, string(runes[i-1]))
			if runes[i-1] == ' ' {
				y -= charHeight / spaceRatio
			} else {
				y -= charHeight
			}
		}
		return
	}

	y := startY + charHeight
	for _, r := range runes {
		d.Text(startX, y, string(r))
		if r == ' ' {
			y += charHeight / spaceRatio
		} else {
			y += charHeight
		}
	}
}
