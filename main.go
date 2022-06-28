package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/dlasky/gotk3-layershell/layershell"
	"github.com/gotk3/gotk3/gtk"
	log "github.com/sirupsen/logrus"
)

const ver = "0.0.1"

var (
	wayland   bool
	briSlider *gtk.Scale
	conSlider *gtk.Scale
	combo     *gtk.ComboBoxText
)

var executor = flag.Bool("e", false, "print brightness Executor data")
var busNum = flag.Int("b", -1, "Bus number for /dev/i2c-<bus number>")
var debug = flag.Bool("d", false, "turn on Debug messages")
var iconSet = flag.String("i", "light", "Icon set to use")
var displayVersion = flag.Bool("v", false, "display Version information")

func main() {
	flag.Parse()
	if *debug {
		log.SetLevel(log.DebugLevel)
	}
	if *displayVersion {
		fmt.Printf("ddcpopup version %s\n", ver)
		os.Exit(0)
	}

	iconsPath := path.Join(configDir(), fmt.Sprintf("nwg-panel/icons_%v", *iconSet))
	log.Debugf("Icons path: %s", iconsPath)

	if *executor {
		bri := getBrightness()
		iconName := "display-brightness-low-symbolic"
		if bri > 70 {
			iconName = "display-brightness-high-symbolic"
		} else if bri >= 30 {
			iconName = "display-brightness-medium-symbolic"
		}
		fmt.Printf("%s.svg\n", path.Join(iconsPath, iconName))
		fmt.Printf("%v%%\n", bri)
		os.Exit(0)
	}

	wayland = waylandSession()

	// brightness := getBrightness()
	// contrast := getContrast()
	// activePreset := getActivePreset()
	name, presets, err := getPresets()
	if err == nil {
		fmt.Println(name)
		for _, preset := range presets {
			fmt.Printf("'%s'\n", preset)
		}
	}

	gtk.Init(nil)
	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	if wayland {
		layershell.InitForWindow(win)
	}

	var displayName string
	if *busNum > -1 {
		displayName = fmt.Sprintf(" %s (%v) ", name, *busNum)
	} else {
		displayName = name
	}
	frame, _ := gtk.FrameNew(displayName)
	_ = frame.SetProperty("margin", 6)
	frame.SetLabelAlign(0.5, 0.5)
	win.Add(frame)

	grid, _ := gtk.GridNew()
	grid.SetColumnSpacing(6)
	grid.SetRowSpacing(12)
	// grid.SetColumnHomogeneous(true)
	_ = grid.SetProperty("margin", 6)
	frame.Add(grid)

	lbl, _ := gtk.LabelNew("")
	_ = lbl.SetProperty("halign", gtk.ALIGN_END)
	lbl.SetMarkup("<tt>Brightness:</tt>")
	grid.Attach(lbl, 0, 0, 1, 1)

	briSlider, _ = gtk.ScaleNewWithRange(gtk.ORIENTATION_HORIZONTAL, 0, 100, 1)
	grid.Attach(briSlider, 1, 0, 2, 1)
	// briSlider.SetValue(float64(brightness))

	lbl, _ = gtk.LabelNew("")
	_ = lbl.SetProperty("halign", gtk.ALIGN_END)
	lbl.SetMarkup("<tt>Contrast:</tt>")
	grid.Attach(lbl, 0, 1, 1, 1)

	conSlider, _ = gtk.ScaleNewWithRange(gtk.ORIENTATION_HORIZONTAL, 0, 100, 1)
	grid.Attach(conSlider, 1, 1, 2, 1)
	// conSlider.SetValue(float64(contrast))

	lbl, _ = gtk.LabelNew("")
	_ = lbl.SetProperty("halign", gtk.ALIGN_END)
	lbl.SetMarkup("<tt>Preset:</tt>")
	grid.Attach(lbl, 0, 2, 1, 1)

	if presets != nil {
		combo, _ = gtk.ComboBoxTextNew()
		for _, preset := range presets {
			vals := strings.Split(preset, ": ")
			id := fmt.Sprintf("0x%s", vals[0])
			fmt.Println("id = ", id)
			text := vals[1]
			combo.Append(id, text)
		}
		// fmt.Println("activePreset = ", activePreset)
		// combo.SetActiveID(activePreset)
		grid.Attach(combo, 1, 2, 1, 1)
	}

	btn, _ := gtk.ButtonNew()
	btn.SetLabel("Close")
	_ = btn.SetProperty("halign", gtk.ALIGN_END)
	grid.Attach(btn, 2, 2, 1, 1)
	btn.Connect("clicked", func() {
		gtk.MainQuit()
	})

	win.ShowAll()

	go func() {
		briSlider.SetValue(float64(getBrightness()))
		conSlider.SetValue(float64(getContrast()))
		combo.SetActiveID(getActivePreset())
	}()

	// time.Sleep(300 * time.Millisecond)

	// go func(msg string) {
	// 	fmt.Println(msg)
	// 	conSlider.SetValue(float64(getContrast()))
	// }("contrast")

	// time.Sleep(300 * time.Millisecond)

	// go func(msg string) {
	// 	fmt.Println(msg)
	// 	combo.SetActiveID(getActivePreset())
	// }("preset")

	gtk.Main()
}
