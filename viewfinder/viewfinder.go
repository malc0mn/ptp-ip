package viewfinder

import (
	"github.com/malc0mn/ptp-ip/ptp"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
)

// WidgetDrawer defines the signature of the drawer function of a widget.
type WidgetDrawer func(*Widget, int64)

// Viewfinder holds a list of pointers to Widgets mapped to their ptp.DevicePropCode.
type Viewfinder struct {
	Widgets map[ptp.DevicePropCode]*Widget
}

// DrawWidget draws the widget mapped to the given device property code on the given image with the given value.
func (vf *Viewfinder) DrawWidget(img *image.RGBA, code ptp.DevicePropCode, val int64) {
	if w, ok := vf.Widgets[code]; ok {
		w.Dst = img
		w.Draw(w, val)
	}
}

// NewViewfinder creates a vendor specific viewfinder using the image passed in to allow each Widget to calibrate its
// starting position.
// When the vendor has no viewfinder defined, nothing will happen.
func NewViewfinder(img *image.RGBA, v ptp.VendorExtension) *Viewfinder {
	switch v {
	case ptp.VE_FujiPhotoFilmCoLtd:
		return NewFujiXT1Viewfinder(img)
	}

	return nil
}

// DrawViewfinder draws all the viewfinder widgets when present in the given ptp.DevicePropDesc list.
func DrawViewfinder(vf *Viewfinder, img *image.RGBA, s []*ptp.DevicePropDesc) {
	for _, p := range s {
		vf.DrawWidget(img, p.DevicePropertyCode, p.CurrentValueAsInt64())
	}
}

// Widget defines a viewfinder widget.
type Widget struct {
	*font.Drawer
	origin fixed.Point26_6
	face   font.Face
	colour *image.Uniform
	Draw   WidgetDrawer
}

// SetColour sets the font colour to the given red, green and blue values.
func (w *Widget) SetColour(r, g, b uint8) {
	w.Src = image.NewUniform(color.RGBA{R: r, G: g, B: b, A: 255})
}

// ResetColour resets the colour to the original one when the widget was first made.
func (w *Widget) ResetColour() {
	w.Src = w.colour
}

// ResetFace resets the font face to the original one when the widget was first made.
func (w *Widget) ResetFace() {
	w.Face = w.face
}

// ResetToOrigin resets the start drawing position to the original point calculated when the widget was first made.
func (w *Widget) ResetToOrigin() {
	w.Dot = w.origin
}

// NewWidget needs a colour to draw in and x/y coordinates to start drawing from.
// Important: the destination image is NOT set but can be set later using Widget.SetImage()!
func NewWidget(img *image.RGBA, r, g, b uint8, f *basicfont.Face, x, y int) *Widget {
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}
	col := image.NewUniform(color.RGBA{R: r, G: g, B: b, A: 255})

	return &Widget{
		Drawer: &font.Drawer{
			Dst:  img,
			Src:  col,
			Face: f,
			Dot:  point,
		},
		origin: point,
		face:   f,
		colour: col,
	}
}

// NewFontWidget returns a new Widget using basicfont.Face7x13 for its basicfont.Face.
func NewFontWidget(img *image.RGBA, r, g, b uint8, x, y int) *Widget {
	return NewWidget(img, r, g, b, basicfont.Face7x13, x, y)
}

// NewWhiteFontWidget returns a new Widget using basicfont.Face7x13 for its basicfont.Face and white (255, 255, 255) for
// its drawing colour.
func NewWhiteFontWidget(img *image.RGBA, x, y int) *Widget {
	return NewFontWidget(img, 255, 255, 255, x, y)
}

// NewGlyphWidget returns a new Widget using VFGlyphs6x13 for its basicfont.Face.
func NewGlyphWidget(img *image.RGBA, r, g, b uint8, x, y int) *Widget {
	return NewWidget(img, r, g, b, VFGlyphs6x13, x, y)
}

// NewWhiteGlyphWidget returns a new Widget using VFGlyphs6x13 for its basicfont.Face and white (255, 255, 255) for
// its drawing colour.
func NewWhiteGlyphWidget(img *image.RGBA, x, y int) *Widget {
	return NewGlyphWidget(img, 255, 255, 255, x, y)
}
