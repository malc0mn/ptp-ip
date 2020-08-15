// +build with_lv

package main

import (
	ptpfmt "github.com/malc0mn/ptp-ip/fmt"
	"github.com/malc0mn/ptp-ip/ip"
	"github.com/malc0mn/ptp-ip/ptp"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"math"
	"strings"
)

func addViewfinder(img *image.RGBA, v ptp.VendorExtension, s []*ptp.DevicePropDesc) {
	switch v {
	case ptp.VE_FujiPhotoFilmCoLtd:
		fujiViewfinder(img, s)
	}
}

/*
use state here?

2020/08/13 17:01:26 [iShell] message received: 'state'
battery level: 2/3 || 3 - [3 0 0 0] - 0x03000000 - 0x3
image size: L 3:2 || 10 - [10 0 0 0] - 0x0a000000 - 0xa
white balance: automatic || 2 - [2 0 0 0] - 0x02000000 - 0x2
F-number: f/2.8 || 280 - [24 1 0 0] - 0x18010000 - 0x118
focus mode: single auto || 32769 - [1 128 0 0] - 0x01800000 - 0x8001
flash mode: enabled || 32778 - [10 128 0 0] - 0x0a800000 - 0x800a
shutter speed:  || 34464 - [160 134 1 128] - 0xa0860180 - 0x86a0
exposure program mode: automatic || 2 - [2 0 0 0] - 0x02000000 - 0x2
exposure bias compensation: 0 || 0 - [0 0 0 0] - 0x00000000 - 0x0
capture delay: off || 0 - [0 0 0 0] - 0x00000000 - 0x0
film simulation: PROVIA || 1 - [1 0 0 0] - 0x01000000 - 0x1
image quality: fine + RAW || 4 - [4 0 0 0] - 0x04000000 - 0x4
command dial mode: both || 0 - [0 0 0 0] - 0x00000000 - 0x0
ISO: 200 || 200 - [200 0 0 0] - 0xc8000000 - 0xc8
focus point: 4x4 || 1028 - [4 4 2 3] - 0x04040203 - 0x404
focus lock: off || 0 - [0 0 0 0] - 0x00000000 - 0x0
device error: none || 0 - [0 0 0 0] - 0x00000000 - 0x0
image space SD:  || 1454 - [174 5 0 0] - 0xae050000 - 0x5ae
movie remaining time:  || 1679 - [143 6 0 0] - 0x8f060000 - 0x68f
*/

// P [o]     F2.8     +/-  -3..-2..-1..|..1..2..3      iso200  bat3/3
func fujiViewfinder(img *image.RGBA, s []*ptp.DevicePropDesc) {
	for _, p := range s {
		switch p.DevicePropertyCode {
		case ip.DPC_Fuji_ExposureIndex:
			fujiIso(img, p.CurrentValueAsInt64())
		case ptp.DPC_BatteryLevel:
			fujiBattery3Bars(img, p.CurrentValueAsInt64())
		case ptp.DPC_ExposureBiasCompensation:
			fujiExposureBiasCompensation(img, p.CurrentValueAsInt64())
		case ptp.DPC_FNumber:
			fujiFNumber(img, p.CurrentValueAsInt64())
		case ptp.DPC_ExposureProgramMode:
			fujiExposureProgramMode(img, p.CurrentValueAsInt64())
		}
	}
}

func fujiExposureProgramMode(img *image.RGBA, epm int64) {
	col := color.RGBA{R: 255, G: 255, B: 255, A: 255} // white

	x := float64(img.Bounds().Min.X) + (float64(img.Bounds().Max.X) * 0.1)
	y := img.Bounds().Max.Y - 10
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: VFGlyphs6x13,
		Dot:  point,
	}

	icon := " "
	switch ptp.ExposureProgramMode(epm) {
	case ptp.EPM_Manual:
		icon = "Mm"
	case ptp.EPM_Automatic:
		icon = "Pp"
	case ptp.EPM_AperturePriority:
		icon = "Nn"
	case ptp.EPM_ShutterPriority:
		icon = "Ll"
	}

	d.DrawString(icon)
}

func fujiFNumber(img *image.RGBA, fn int64) {
	col := color.RGBA{R: 255, G: 255, B: 255, A: 255} // white

	x := float64(img.Bounds().Min.X) + (float64(img.Bounds().Max.X) * 0.25)
	y := img.Bounds().Max.Y - 10
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}

	d.DrawString(strings.Replace(ptpfmt.FNumberAsString(uint16(fn)), "f/", "F", 1))
}

func fujiExposureBiasCompensation(img *image.RGBA, ex int64) {
	zero := 9 // don't forget: zero indexed!
	stops := 3 // bias dial is per 3 stops
	stop := 0 // default stop is '0' meaning no fractional stop
	bias := []rune("6..5..4..0..1..2..3")
	marker := []rune("                   ")

	// Make sure the center point of our bias widget is in the center of the image.
	glyphwidth := 6
	offset := glyphwidth * len(bias) / 2

	x := float64(img.Bounds().Max.X) - (float64(img.Bounds().Max.X) * 0.5) - float64(offset)
	y := img.Bounds().Max.Y - 10
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}
	col := color.RGBA{R: 255, G: 255, B: 255, A: 255} // white

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: VFGlyphs6x13,
		Dot:  point,
	}

	// Draw the leading +/- icon
	d.Dot.X -= fixed.Int26_6(18 * 64) // offset icon 3 glyphs to the left
	d.DrawString("+-")
	d.Dot = point // reset drawing position

	// Calculate marker position.
	i, f := math.Modf(float64(int16(ex)) / float64(1000))
	onZero := i == 0 && f == 0
	if f != 0 {
		stop = 1
		if math.Abs(f) > 0.4 {
			stop = 2
		}
		if math.Signbit(f) {
			stop = -stop
		}
	}
	pos := zero + stop + int(i * float64(stops))

	// When we're not on a fractional number, replace the number with an 'empty' marker.
	if f == 0 {
		bias[pos] = '"'
	}

	// When the marker is on 0, the widget must be drawn in grey.
	if onZero {
		col = color.RGBA{R: 100, G: 100, B: 100, A: 255} // grey
		d.Src = image.NewUniform(col)
	}

	// Now draw the basic exposure bias compensation widget.
	d.DrawString(string(bias))

	// When the marker is on 0, the the marker and '0' position must be drawn in white.
	if onZero {
		col = color.RGBA{R: 255, G: 255, B: 255, A: 255} // white
		d.Src = image.NewUniform(col)
		for _, r := range []rune{'"', '!'} {
			d.Dot = point // reset drawing position
			marker[pos] = r
			d.DrawString(string(marker))
		}

		return
	}

	// Draw the marker on the the calculated position in yellow!
	marker[pos] = '!'
	col = color.RGBA{R: 255, G: 185, B: 10, A: 255} // yellow
	d.Src = image.NewUniform(col)
	d.Dot = point // reset drawing position
	d.DrawString(string(marker))
}

func fujiIso(img *image.RGBA, ex int64) {
	col := color.RGBA{R: 255, G: 255, B: 255, A: 255} // white

	x := float64(img.Bounds().Max.X) - (float64(img.Bounds().Max.X) * 0.2)
	y := img.Bounds().Max.Y - 10
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: VFGlyphs6x13,
		Dot:  point,
	}

	iso := ptpfmt.FujiExposureIndexAsString(ip.FujiExposureIndex(ex))

	d.DrawString("is") // iso icon

	if strings.HasPrefix(iso, "S") {
		d.Dot.X -= fixed.Int26_6(18 * 64)
		d.Dot.Y -= fixed.Int26_6(8 * 64)
		d.DrawString("ISO") // auto icon
		d.Dot.Y += fixed.Int26_6(8 * 64) // reset Y axis
		iso = string([]rune(iso)[1:]) // drop the leading S
	}

	d.Face = basicfont.Face7x13
	d.Dot.X += fixed.Int26_6(6 * 64)
	d.Dot.Y += fixed.Int26_6(2 * 64)

	// actual value
	d.DrawString(iso)
}

func fujiBattery3Bars(img *image.RGBA, bl int64) {
	col := color.RGBA{R: 255, G: 255, B: 255, A: 255} // white
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}

	var lvl string
	switch ip.FujiBatteryLevel(bl) {
	case ip.BAT_Fuji_3bOne:
		col = red
		lvl = "baU"
	case ip.BAT_Fuji_3bTwo:
		lvl = "bCT"
	case ip.BAT_Fuji_3bFull:
		lvl = "BAT"
	}

	x := float64(img.Bounds().Max.X) - (float64(img.Bounds().Max.X) * 0.1)
	y := img.Bounds().Max.Y - 8
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: VFGlyphs6x13,
		Dot:  point,
	}
	d.DrawString(lvl)
}
