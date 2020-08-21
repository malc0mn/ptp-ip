package main

import (
	"fmt"
	"github.com/malc0mn/ptp-ip/ip"
	"io/ioutil"
	"strconv"
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
	amount := 1
	hasIntArg := false
	if len(f) >= 1 {
		if val, err := strconv.Atoi(f[0]); err == nil {
			amount = val
			hasIntArg = true
		}
	}

	var img []byte
	for i := 0; i < amount; i++ {
		var err error
		img, err = c.InitiateCapture()
		if err != nil {
			return err.Error()
		}
		// TODO: add support to save/view ALL captured images!
	}

	if !hasIntArg && len(f) >= 1 {
		if cap.isView(f[0]) {
			return preview(img) + "\n"
		}

		if err := ioutil.WriteFile(f[0], img, 0644); err != nil {
			return err.Error() + "\n"
		}

		return fmt.Sprintf("Image preview saved to %s\n", f[0])
	}

	plural := ""
	if amount > 1 {
		plural = "s"
	}

	return fmt.Sprintf("Image%s captured, check the camera\n", plural)
}

func (cap capture) help() string {
	help := `"` + cap.name() + `" will make the responder capture a single image.` + "\n"
	help += helpAddAliases(cap.alias())

	if args := cap.arguments(); len(args) > 0 {
		help += helpAddArgumentsTitle()
		for i, arg := range args {
			switch i {
			case 0:
				help += "\t- " + arg + ": an integer value to indicate the amount of captures to make\n\tOR\n"
			case 1:
				help += "\t- " + `"` + arg + `" opens a window to display the capture preview if the camera returns it` + "\n\tOR\n"
			case 2:
				help += "\t- a " + arg + " to save the capture preview to\n"
			}
		}
	}

	return help
}

func (capture) arguments() []string {
	return []string{"amount", "view", "filepath"}
}

func (cap capture) isView(param string) bool {
	return param == cap.arguments()[1]
}
