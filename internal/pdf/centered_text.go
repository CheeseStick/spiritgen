package pdf

// DrawTextCentered horizontally centers text within the box that starts at
// boxX and has width boxWidth. baselineY follows fpdf's Text() convention —
// it is the y-coordinate of the text's baseline (the bottom edge of most
// glyphs); callers typically pass `top + lineHeight` to place the text just
// below `top`.
//
// A font must already be active via SetFont before calling — both
// GetStringWidth and Text rely on it. For UTF-8 fonts, that font must have
// been registered via RegisterUTF8Font.
func (d *Doc) DrawTextCentered(text string, boxX, boxWidth, baselineY float64) {
	w := d.GetStringWidth(text)
	x := boxX + (boxWidth-w)/2
	d.Text(x, baselineY, text)
}
