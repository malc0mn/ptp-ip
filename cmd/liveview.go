// +build with_lv

package main

import (
	"bytes"
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/malc0mn/ptp-ip/ip"
	"image"
)

var lvState bool
var lvTerm chan bool

func init() {
	liveview = toggleLv
	lvTerm = make(chan bool, 1)
}

func toggleLv(c *ip.Client, f []string) string {
	errorFmt := "liveview error: %s\n"

	lvState = !lvState

	if err := c.ToggleLiveView(lvState); err != nil {
		return fmt.Sprintf(errorFmt, err)
	}

	if lvState {
		go liveviewUI(c.StreamChan, lvTerm)

		return "enabled\n"
	}

	lvTerm <- true

	return "disabled\n"
}

func liveviewUI(imgs chan []byte, term chan bool) error {
	if err := gl.Init(); err != nil {
		return err
	}

	if err := glfw.Init(); err != nil {
		return err
	}
	defer glfw.Terminate()

	img := <-imgs
	im, _, err := image.Decode(bytes.NewReader(img))
	if err != nil {
		return err
	}

	window, err := NewWindow(im, "Live view")
	if err != nil {
		return err
	}

	// TODO: properly handle window close!!!
	if window.ShouldClose() {
		window.Destroy()
	}
	window.Draw()
	glfw.PollEvents()

	for {
		select {
		case <-term:
			window.Destroy()

			return nil
		case img := <-imgs:
			im, _, err := image.Decode(bytes.NewReader(img))
			if err == nil {
				window.SetImage(im)
			}
		}
	}
}
