package main

import (
	"fmt"
	"github.com/malc0mn/ptp-ip/ip"
	"io/ioutil"
)

func init() {
	registerCommand(&capture{})
}

type capture struct{}

func (capture) name() string {
	return "capture"
}

func (capture) alias() []string {
	return []string{"shoot", "shutter", "snap"}
}

func (cap capture) execute(c *ip.Client, f []string) string {
	img, err := c.InitiateCapture()
	if err != nil {
		return err.Error()
	}
	if len(f) == 1 {
		if cap.isView(f[0]) {
			return preview(img) + "\n"
		}

		if err := ioutil.WriteFile(f[0], img, 0644); err != nil {
			return err.Error() + "\n"
		}

		return fmt.Sprintf("Image preview saved to %s\n", f[0])
	}

	return "Image captured, check the camera\n"
}

func (cap capture) help() string {
	help := `"` + cap.name() + `" will make the responder capture a single image.` + "\n"
	help += helpAddAliases(cap.alias())

	if args := cap.arguments(); len(args) > 0 {
		help += helpAddArgumentsTitle()
		for i, arg := range args {
			switch i {
			case 0:
				help += "\t- " + `"` + arg + `" opens a window to display the capture preview if the camera returns it` + "\n\tOR\n"
			case 1:
				help += "\t- a " + arg + " to save the capture preview to"
			}
		}
	}

	return help
}

func (capture) arguments() []string {
	return []string{"view", "filepath"}
}

func (cap capture) isView(param string) bool {
	return param == cap.arguments()[0]
}
