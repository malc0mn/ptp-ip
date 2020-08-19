// Taken from github.com/fogleman/imview and slightly customized.

// +build with_lv

package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"image"
	"image/draw"
)

type texture struct {
	handle uint32
}

func newTexture() *texture {
	var handle uint32
	gl.GenTextures(1, &handle)
	t := &texture{handle}
	t.setMinFilter(gl.LINEAR)
	t.setMagFilter(gl.NEAREST)
	t.setWrapS(gl.CLAMP_TO_EDGE)
	t.setWrapT(gl.CLAMP_TO_EDGE)
	return t
}

func (t *texture) bind() {
	gl.BindTexture(gl.TEXTURE_2D, t.handle)
}

func (t *texture) setImage(im image.Image) {
	rgba, ok := im.(*image.RGBA)
	if !ok {
		rgba = image.NewRGBA(im.Bounds())
		draw.Draw(rgba, rgba.Rect, im, image.Point{}, draw.Src)
	}
	t.setRGBA(rgba)
}

func (t *texture) setRGBA(im *image.RGBA) {
	t.bind()
	size := im.Rect.Size()
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA, int32(size.X), int32(size.Y),
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(im.Pix),
	)
}

func (t *texture) setMinFilter(x int32) {
	t.bind()
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, x)
}

func (t *texture) setMagFilter(x int32) {
	t.bind()
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, x)
}

func (t *texture) setWrapS(x int32) {
	t.bind()
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, x)
}

func (t *texture) setWrapT(x int32) {
	t.bind()
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, x)
}
