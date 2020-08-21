// +build with_lv

package main

import (
	"bytes"
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/malc0mn/ptp-ip/ip"
	"github.com/malc0mn/ptp-ip/ptp"
	"github.com/malc0mn/ptp-ip/viewfinder"
	"image"
	"image/draw"
	"time"
)

var (
	lvState   bool
	mainStack = make(chan func())
)

func init() {
	registerCommand(&liveview{})
}

type liveview struct{}

func (liveview) name() string {
	return "liveview"
}

func (liveview) alias() []string {
	return []string{}
}

func (liveview) execute(c *ip.Client, _ []string) string {
	errorFmt := "liveview error: %s\n"

	if lvState {
		return "already enabled!\n"
	}

	lvState = true

	if err := c.ToggleLiveView(lvState); err != nil {
		return fmt.Sprintf(errorFmt, err)
	}

	runOnMain(func() { liveViewUI(c) })

	return "enabled\n"
}

func (l liveview) help() string {
	return `"` + l.name() + `" opens a window and displays a live view through the camera lens. Not all vendors support this!` + "\n"
}

func (liveview) arguments() []string {
	return []string{}
}

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

// runOnMain executes f on the main thread but does not wait for it to finish.
func runOnMain(f func()) {
	mainStack <- f
}

func liveViewUI(c *ip.Client) error {
	s, err := c.GetDeviceState()
	if err != nil {
		s = []*ptp.DevicePropDesc{}
	}

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

	var vf *viewfinder.Viewfinder
	im, _, err := image.Decode(bytes.NewReader(img))
	if err == nil {
		vf = viewfinder.NewViewfinder(toRGBA(im), c.ResponderVendor())
	}

	ticker := time.NewTicker(1 * time.Second)

poller:
	for !window.ShouldClose() {
		select {
		case img := <-c.StreamChan:
			im, _, err := image.Decode(bytes.NewReader(img))
			if err == nil {
				rgba := toRGBA(im)
				data, ok := s.([]*ptp.DevicePropDesc)
				if vf != nil && ok {
					viewfinder.DrawViewfinder(vf, rgba, data)
				}
				window.setImage(rgba)
			}
		case <-ticker.C:
			s, _ = c.GetDeviceState()
		case <-quit:
			break poller
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

func toRGBA(img image.Image) *image.RGBA {
	rgba, ok := img.(*image.RGBA)
	if !ok {
		rgba = image.NewRGBA(img.Bounds())
		draw.Draw(rgba, rgba.Rect, img, image.Point{}, draw.Src)
	}

	return rgba
}

func preview(img []byte) string {
	// TODO: figure out how to cleanly have multiple windows open at the same time 'on the main thread' by introducing some
	//  sort of extremely simple window manager.
	if lvState {
		return "can currently not display preview while liveview is active"
	}

	runOnMain(func() { previewUI(img) })

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
	for !window.ShouldClose() {
		select {
		case <-quit:
			break poller
		default:
			glfw.WaitEvents()
		}
	}

	window.Destroy()

	return nil
}

func showImage(img []byte, title string) (*window, error) {
	im, _, err := image.Decode(bytes.NewReader(img))
	if err != nil {
		return nil, err
	}

	window, err := newWindow(im, title)
	if err != nil {
		return nil, err
	}
	window.draw()

	return window, nil
}
