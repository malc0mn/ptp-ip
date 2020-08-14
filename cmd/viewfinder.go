// +build with_lv

package main

import (
	"github.com/malc0mn/ptp-ip/ip"
	"github.com/malc0mn/ptp-ip/ptp"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
)

func addViewfinder(img *image.RGBA, c *ip.Client) {
	switch c.ResponderVendor() {
	case ptp.VE_FujiPhotoFilmCoLtd:
		fujiViewfinder(img)
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
func fujiViewfinder(img *image.RGBA) {
	x := 10
	y := 20
	col := color.RGBA{255, 255, 255, 255}
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: VFGlyphs6x13,
		Dot:  point,
	}

	// Battery test: renders 3/3 2/3 1/3 critical.
	d.DrawString("BAT bCT baU bct")
}
