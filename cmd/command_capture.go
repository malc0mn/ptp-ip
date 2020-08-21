package main

import (
	"fmt"
	"github.com/malc0mn/ptp-ip/ip"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
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
	if len(f) >= 1 {
		if val, err := strconv.Atoi(f[0]); err == nil {
			amount = val
			f = f[1:] // drop processed amount argument
		}
	}

	var (
		imgs chan []byte
		output string
		wg sync.WaitGroup
	)
	if len(f) >= 1 {
		imgs = make(chan []byte, 10)
		var path string
		if !cap.isView(f[0]) {
			ext := filepath.Ext(f[0])
			path = strings.TrimSuffix(f[0], ext) + "-%d" + ext
		}

		wg.Add(1)
		go func() {
			i := 1
			for img := range imgs {
				if path != "" {
					file := fmt.Sprintf(path, i)
					if err := ioutil.WriteFile(file, img, 0644); err != nil {
						output += err.Error() + "\n"
						continue
					}
					output += fmt.Sprintf("Image preview saved to %s\n", file)
					i++
				} else {
					output += preview(img) + "\n"
				}
			}
			wg.Done()
		}()
	}

	// TODO: send fmt.Sprintf("capturing %d images", amount) to a channel for instant output to the end user.
	for i := 0; i < amount; i++ {
		// TODO: send fmt.Sprintf("capturing image %d", i) to a channel for instant output to the end user.
		var err error
		img, err := c.InitiateCapture()
		if err != nil {
			return err.Error()
		}
		if imgs != nil {
			imgs <- img
		}
	}
	if imgs != nil {
		close(imgs)
		wg.Wait()
		return output
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
				help += "\t- " + arg + ": an integer value to indicate the amount of captures to make\n"
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
