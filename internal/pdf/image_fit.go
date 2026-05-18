package pdf

// DrawImageAspectFit draws the previously-registered image identified by key
// into the box (boxX, boxY, boxW, boxH), centered and aspect-ratio preserved.
// The image must have been registered via RegisterPNGImage (or equivalent)
// before calling this method. If the image is not found or has invalid
// dimensions, the call is a no-op.
func (d *Doc) DrawImageAspectFit(key string, boxX, boxY, boxW, boxH float64) {
	info := d.GetImageInfo(key)
	if info == nil {
		return
	}
	imgW, imgH := info.Extent()
	if imgW <= 0 || imgH <= 0 {
		return
	}

	ratio := imgW / imgH
	boxRatio := boxW / boxH

	var drawW, drawH float64
	if ratio > boxRatio {
		drawW = boxW
		drawH = boxW / ratio
	} else {
		drawH = boxH
		drawW = boxH * ratio
	}

	drawX := boxX + (boxW-drawW)/2
	drawY := boxY + (boxH-drawH)/2
	d.Image(key, drawX, drawY, drawW, drawH, false, "", -1, "")
}
