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

func init() {
	lvEnabled = true
	liveview = openLv
	preview = openPreview
}

func openLv(c *ip.Client, _ []string) string {
	errorFmt := "liveview error: %s\n"

	if lvState {
		return "already enabled!\n"
	}

	lvState = true

	if err := c.ToggleLiveView(lvState); err != nil {
		return fmt.Sprintf(errorFmt, err)
	}

	go liveViewUI(c)

	return "enabled\n"
}

// TODO: this is unstable for a yet unknown reason. The first window opened never has a problem, but any following
//  initialization of live view can crash with a segfault. Any pointers as to why would be great :/
func liveViewUI(c *ip.Client) error {
	if err := gl.Init(); err != nil {
		return err
	}

	if err := glfw.Init(); err != nil {
		return err
	}
	defer glfw.Terminate()

	img := <-c.StreamChan
	window, err := showImage(img, "Live view")
	if err != nil {
		return err
	}

	for !window.ShouldClose() {
		img := <-c.StreamChan
		im, _, err := image.Decode(bytes.NewReader(img))
		if err == nil {
			window.SetImage(im)
		}
		glfw.PollEvents()
	}

	window.Destroy()
	lvState = false
	if err := c.ToggleLiveView(lvState); err != nil {
		return err
	}

	return nil
}

func openPreview(img []byte) string {
	go previewUI(img)

	return "preview window opened"
}

// TODO: same as with the live view: this is unstable for an unknown reason.
func previewUI(img []byte) error {
	if err := gl.Init(); err != nil {
		return err
	}

	if err := glfw.Init(); err != nil {
		return err
	}
	defer glfw.Terminate()

	window, err := showImage(img, "Capture preview")
	if err != nil {
		return err
	}

	for !window.ShouldClose() {
		window.Draw()
		glfw.PollEvents()
	}

	window.Destroy()

	return nil
}

func showImage(img []byte, title string) (*Window, error) {
	im, _, err := image.Decode(bytes.NewReader(img))
	if err != nil {
		return nil, err
	}

	window, err := NewWindow(im, title)
	if err != nil {
		return nil, err
	}
	window.Draw()

	return window, nil
}
