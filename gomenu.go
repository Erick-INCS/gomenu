package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/jroimartin/gocui"
)

type shortcut struct {
	Name     string
	cmd      string
	terminal bool
}

var regexName, rNameOk = regexp.Compile("Name=(.*)")
var regexCMD, rCMDOk = regexp.Compile("Exec=(.*)")
var regexTerm, rTermOk = regexp.Compile("Terminal=(.*)")
var regexType, rTypeOk = regexp.Compile("Type=(.*)")
var shortcuts = []shortcut{}
var pointer int = 0

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("main", 1, 1, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Wrap = true

		files, err := ioutil.ReadDir("/usr/share/applications")
		if err != nil {
			log.Fatal(err)
		}

		var arr = make([]string, 2)

		for _, file := range files {
			if file.Size() > 0 {
				content, err := ioutil.ReadFile("/usr/share/applications/" + file.Name())
				if err != nil {
					log.Fatal(err)
				}

				arr = regexType.FindStringSubmatch(fmt.Sprintf("%s", content))
				if len(arr) == 2 && strings.ToUpper(arr[1]) == "APPLICATION" {

					var s shortcut

					arr = regexName.FindStringSubmatch(fmt.Sprintf("%s", content))
					if len(arr) == 2 {
						s.Name = arr[1]
					}

					arr = regexCMD.FindStringSubmatch(fmt.Sprintf("%s", content))
					if len(arr) == 2 {
						s.cmd = arr[1]
					}
					arr = regexTerm.FindStringSubmatch(fmt.Sprintf("%s", content))
					if len(arr) == 2 {
						s.terminal = strings.ToUpper(arr[1]) == "TERMINAL"
					}

					shortcuts = append(shortcuts, s)
				}

				// fmt.Fprintln(v, fmt.Sprintf("%s", content))
			}
		}

		// _, lines := v.Size()

		for _, b := range shortcuts {
			fmt.Fprintln(v, b.Name)
		}

		//line := strings.Repeat("This is a long line -- ", 10)
		// fmt.Fprintf(v, "%s\n\n", line)
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func moveDownV(g *gocui.Gui, v *gocui.View) error {
	v2, e := g.View("main")
	if e != nil {
		log.Fatal(e)
	}

	_, sy := v2.Size()

	if pointer+1+sy < len(shortcuts) {
		pointer += 1
		v2.SetOrigin(0, pointer)
	}
	return nil
}

func moveUpV(g *gocui.Gui, v *gocui.View) error {
	v2, e := g.View("main")
	if e != nil {
		log.Fatal(e)
	}

	if pointer-1 >= 0 {
		pointer -= 1
		v2.SetOrigin(0, pointer)
	}
	return nil
}

func main() {
	if rNameOk != nil {
		log.Panicln(rNameOk)
	}
	if rCMDOk != nil {
		log.Panicln(rCMDOk)
	}
	if rTermOk != nil {
		log.Panicln(rTermOk)
	}
	if rTypeOk != nil {
		log.Panicln(rTypeOk)
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, moveDownV); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, moveUpV); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
