package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/dlasky/gotk3-layershell/layershell"
	"github.com/gotk3/gotk3/gtk"
	log "github.com/sirupsen/logrus"
)

const ver = "0.0.1"

var (
	wayland         bool
	briValueChanged bool
	conValueChanged bool
	briSlider       *gtk.Scale
	conSlider       *gtk.Scale
	combo           *gtk.ComboBoxText
	timeStart       time.Time
)

var executor = flag.Bool("e", false, "print brightness Executor data")
var label = flag.String("l", "", "print this Label instead of image path")
var busNum = flag.Int("b", -1, "Bus number for /dev/i2c-<number>")
var debug = flag.Bool("d", false, "turn on Debug messages")
var iconSet = flag.String("i", "light", "Icon set to use")
var displayVersion = flag.Bool("v", false, "display Version information")

func main() {
	timeStart = time.Now()
	flag.Parse()
	if *debug {
		log.SetLevel(log.DebugLevel)
	}
	if *displayVersion {
		fmt.Printf("ddcpopup version %s\n", ver)
		os.Exit(0)
	}
	if *busNum == -1 {
		fmt.Println("Bus number required: -b <number>")
		os.Exit(1)
	}

	iconsPath := path.Join(configDir(), fmt.Sprintf("nwg-panel/icons_%v", *iconSet))
	log.Debugf("Icons path: %s", iconsPath)

	if *executor {
		bri, err := getBrightness()
		if err == nil {
			if *label == "" {
				// 2 lines (image path / value) for nwg-panel or Tint2
				iconName := "display-brightness-low-symbolic"
				if bri > 70 {
					iconName = "display-brightness-high-symbolic"
				} else if bri >= 30 {
					iconName = "display-brightness-medium-symbolic"
				}
				fmt.Printf("%s.svg\n", path.Join(iconsPath, iconName))
				fmt.Printf("%v%%\n", bri)
			} else {
				// One-liner for textual panels
				fmt.Printf("%s %v%%\n", *label, bri)
			}
			os.Exit(0)
		} else {
			log.Error(err)
			os.Exit(1)
		}
	}

	wayland = waylandSession()

	name, presets, err := getPresets()
	displayName := fmt.Sprintf(" %s (bus %v) ", name, *busNum)
	if err != nil {
		log.Error(err)
	}

	gtk.Init(nil)
	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	if wayland {
		layershell.InitForWindow(win)
	}

	frame, _ := gtk.FrameNew(displayName)
	_ = frame.SetProperty("margin", 6)
	frame.SetLabelAlign(0.5, 0.5)
	win.Add(frame)

	grid, _ := gtk.GridNew()
	grid.SetColumnSpacing(6)
	grid.SetRowSpacing(12)
	_ = grid.SetProperty("margin", 6)
	frame.Add(grid)

	lbl, _ := gtk.LabelNew("")
	_ = lbl.SetProperty("halign", gtk.ALIGN_END)
	lbl.SetMarkup("<tt>Brightness:</tt>")
	grid.Attach(lbl, 0, 0, 1, 1)

	briSlider, _ = gtk.ScaleNewWithRange(gtk.ORIENTATION_HORIZONTAL, 0, 100, 1)
	_ = briSlider.SetProperty("hexpand", true)
	_ = briSlider.Connect("value-changed", func() {
		briValueChanged = true
	})
	_ = briSlider.Connect("button-release-event", func() {
		if briValueChanged {
			launch(fmt.Sprintf("ddcutil setvcp 10 %v -b %v --noverify", int(briSlider.GetValue()), *busNum))
		}
	})
	grid.Attach(briSlider, 1, 0, 2, 1)

	lbl, _ = gtk.LabelNew("")
	_ = lbl.SetProperty("halign", gtk.ALIGN_END)
	lbl.SetMarkup("<tt>Contrast:</tt>")
	grid.Attach(lbl, 0, 1, 1, 1)

	conSlider, _ = gtk.ScaleNewWithRange(gtk.ORIENTATION_HORIZONTAL, 0, 100, 1)
	_ = conSlider.SetProperty("hexpand", true)
	_ = conSlider.Connect("value-changed", func() {
		conValueChanged = true
	})
	_ = conSlider.Connect("button-release-event", func() {
		if conValueChanged {
			launch(fmt.Sprintf("ddcutil setvcp 12 %v -b %v --noverify", int(conSlider.GetValue()), *busNum))
		}
	})
	grid.Attach(conSlider, 1, 1, 2, 1)

	lbl, _ = gtk.LabelNew("")
	_ = lbl.SetProperty("halign", gtk.ALIGN_END)
	lbl.SetMarkup("<tt>Preset:</tt>")
	grid.Attach(lbl, 0, 2, 1, 1)

	combo, _ = gtk.ComboBoxTextNew()
	grid.Attach(combo, 1, 2, 1, 1)
	if presets != nil {
		for _, preset := range presets {
			values := strings.Split(preset, ": ")
			id := fmt.Sprintf("0x%s", values[0])
			text := values[1]
			combo.Append(id, text)
		}
		combo.Connect("changed", func() {
			dec, err := strconv.ParseInt(strings.Split(combo.GetActiveID(), "x")[1], 16, 64)
			if err == nil {
				launch(fmt.Sprintf("ddcutil setvcp 14 %v -b %v --noverify", dec, *busNum))
			}
		})
	} else {
		combo.Append("unavailable", "unavailable")
		combo.SetActiveID("unavailable")
		combo.SetSensitive(false)
	}

	btn, _ := gtk.ButtonNew()
	btn.SetLabel("Close")
	_ = btn.SetProperty("halign", gtk.ALIGN_END)
	grid.Attach(btn, 2, 2, 1, 1)
	btn.Connect("clicked", func() {
		gtk.MainQuit()
	})

	win.ShowAll()
	win.SetSizeRequest(win.GetAllocatedHeight()*2, 0)

	go func() {
		bri, e := getBrightness()
		if e == nil {
			briSlider.SetValue(float64(bri))
			briValueChanged = false
		} else {
			log.Error(e)
		}

		conSlider.SetValue(float64(getContrast()))
		preset, e := getActivePreset()
		if e == nil {
			combo.SetActiveID(preset)
			conValueChanged = false
		} else {
			log.Error(e)
		}
		t := time.Now()
		log.Debugf("Startup time: %v ms", t.Sub(timeStart).Milliseconds())
	}()

	gtk.Main()
}
