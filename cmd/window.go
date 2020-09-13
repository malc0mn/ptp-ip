// Taken from github.com/fogleman/imview and slightly customized.

// +build with_lv

package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"image"
	_ "image/jpeg"
	"runtime"
)

func init() {
	runtime.LockOSThread()
}

type window struct {
	*glfw.Window
	image   image.Image
	texture *texture
}

func newWindow(im image.Image, title string) (*window, error) {
	const maxSize = 1200
	w := im.Bounds().Size().X
	h := im.Bounds().Size().Y
	a := float64(w) / float64(h)
	if a >= 1 {
		if w > maxSize {
			w = maxSize
			h = int(maxSize / a)
		}
	} else {
		if h > maxSize {
			h = maxSize
			w = int(maxSize * a)
		}
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	win, err := glfw.CreateWindow(w, h, title, nil, nil)
	if err != nil {
		return nil, err
	}

	win.MakeContextCurrent()
	glfw.SwapInterval(1)

	texture := newTexture()
	texture.setImage(im)
	result := &window{win, im, texture}
	result.SetRefreshCallback(result.onRefresh)

	return result, nil
}

func (window *window) setImage(im image.Image) {
	window.image = im
	window.texture.setImage(im)
	window.draw()
}

func (window *window) onRefresh(_ *glfw.Window) {
	window.draw()
}

func (window *window) draw() {
	window.MakeContextCurrent()
	gl.Clear(gl.COLOR_BUFFER_BIT)
	window.drawImage()
	window.SwapBuffers()
}

func (window *window) drawImage() {
	const padding = 0
	iw := window.image.Bounds().Size().X
	ih := window.image.Bounds().Size().Y
	w, h := window.GetFramebufferSize()
	s1 := float32(w) / float32(iw)
	s2 := float32(h) / float32(ih)
	f := float32(1 - padding)
	var x, y float32
	if s1 >= s2 {
		x = f * s2 / s1
		y = f
	} else {
		x = f
		y = f * s1 / s2
	}
	gl.Enable(gl.TEXTURE_2D)
	window.texture.bind()
	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0, 1)
	gl.Vertex2f(-x, -y)
	gl.TexCoord2f(1, 1)
	gl.Vertex2f(x, -y)
	gl.TexCoord2f(1, 0)
	gl.Vertex2f(x, y)
	gl.TexCoord2f(0, 0)
	gl.Vertex2f(-x, y)
	gl.End()
}
