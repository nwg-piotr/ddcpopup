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
var busNum = flag.Int("b", -1, "Bus number for /dev/i2c-<number> (Required; check 'ddcutil 'detect')")
var debug = flag.Bool("d", false, "turn on Debug messages")
var darkIcons = flag.Bool("k", false, "use darK icons")
var displayVersion = flag.Bool("v", false, "display Version information")
var hpos = flag.String("hpos", "", "window Horizontal POSition: 'left', 'right' or none for center (Wayland)")
var vpos = flag.String("vpos", "", "window Vertical POSition: 'top', 'bottom' or none for center (Wayland)")
var hm = flag.Int("hm", 0, "window Horizontal Margin (Wayland)")
var vm = flag.Int("vm", 0, "window Vertical Margin (Wayland)")

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

	var iconsPath string
	if *darkIcons {
		iconsPath = path.Join("/usr/share/ddcpopup/icons/dark")
	} else {
		iconsPath = path.Join("/usr/share/ddcpopup/icons/light")
	}

	log.Debugf("Icons path: %s", iconsPath)

	if *executor {
		bri := getBrightness()
		if bri >= 0 {
			if *label == "" {
				// 2 lines (image path / value) for nwg-panel or Tint2
				var iconName string
				if bri > 80 {
					iconName = "brightness-full"
				} else if bri >= 50 {
					iconName = "brightness-medium"
				} else if bri >= 20 {
					iconName = "brightness-low"
				} else {
					iconName = "brightness-off"
				}
				fmt.Printf("%s.svg\n", path.Join(iconsPath, iconName))
				fmt.Printf("%v%%\n", bri)
			} else {
				// One-liner for textual panels
				fmt.Printf("%s %v%%\n", *label, bri)
			}
			os.Exit(0)
		} else {
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
		layershell.SetLayer(win, layershell.LAYER_SHELL_LAYER_TOP)

		if *hpos == "left" {
			layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_LEFT, true)
		} else if *hpos == "right" {
			layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_RIGHT, true)
		}

		if *vpos == "top" {
			layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_TOP, true)
		} else if *vpos == "bottom" {
			layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_BOTTOM, true)
		}

		if *hm != 0 {
			layershell.SetMargin(win, layershell.LAYER_SHELL_EDGE_LEFT, *hm)
			layershell.SetMargin(win, layershell.LAYER_SHELL_EDGE_RIGHT, *hm)
		}

		if *vm != 0 {
			layershell.SetMargin(win, layershell.LAYER_SHELL_EDGE_TOP, *vm)
			layershell.SetMargin(win, layershell.LAYER_SHELL_EDGE_BOTTOM, *vm)
		}
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
		bri := getBrightness()
		if bri > -1 {
			briSlider.SetValue(float64(bri))
			briValueChanged = false
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
