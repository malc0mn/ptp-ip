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

var (
	lvState   bool
	mainStack = make(chan func())
)

// mainThread is used to execute on the main thread, which is what OpenGL requires.
func mainThread() {
	for {
		select {
		case f := <-mainStack:
			f()
		case <-quit:
			return
		}
	}
}

// do executes f on the main thread but does not wait for it to finish.
func do(f func()) {
	mainStack <- f
}

func liveview(c *ip.Client, _ []string) string {
	errorFmt := "liveview error: %s\n"

	if lvState {
		return "already enabled!\n"
	}

	lvState = true

	if err := c.ToggleLiveView(lvState); err != nil {
		return fmt.Sprintf(errorFmt, err)
	}

	do(func() { liveViewUI(c) })

	return "enabled\n"
}

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
		select {
		case img := <-c.StreamChan:
			im, _, err := image.Decode(bytes.NewReader(img))
			if err == nil {
				window.SetImage(im)
			}
			glfw.PollEvents()
		case <-quit:
			return nil
		}
	}

	window.Destroy()
	lvState = false
	if err := c.ToggleLiveView(lvState); err != nil {
		return err
	}

	return nil
}

func preview(img []byte) string {
	do(func() { previewUI(img) })

	return "preview window opened"
}

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

poller:
	for {
		select {
		case <-quit:
			break poller
		default:
			if window.ShouldClose() {
				break poller
			}
			glfw.PollEvents()
		}
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
