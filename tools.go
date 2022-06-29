package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

func waylandSession() bool {
	return os.Getenv("WAYLAND_DISPLAY") != "" || strings.Contains(os.Getenv("XDG_SESSION_TYPE"), "wayland")
}

func configDir() string {
	var dir string
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		dir = os.Getenv("XDG_CONFIG_HOME")
	} else if os.Getenv("HOME") != "" {
		dir = path.Join(os.Getenv("HOME"), ".config")
	}

	return dir
}

func getCommandOutput(command string) (string, error) {
	log.Debugf("CMD: %s", command)
	elements := strings.Split(command, " ")
	c, b := exec.Command(elements[0], elements[1:]...), new(strings.Builder)
	c.Stdout = b
	err := c.Run()
	output := b.String()
	log.Debugf("OUT: %s", output)

	return output, err
}

func getBrightness() (int, error) {
	command := fmt.Sprintf("ddcutil getvcp 10 --bus=%v", *busNum)
	output, _ := getCommandOutput(command)
	lines := strings.Split(output, "\n")
	lineWithValue := ""
	for _, line := range lines {
		if strings.Contains(line, "Brightness") {
			lineWithValue = strings.Split(line, ",")[0]
		}
	}
	var e error
	if lineWithValue != "" {
		parts := strings.Split(lineWithValue, " ")
		strVal := parts[len(parts)-1]
		intVal, e := strconv.Atoi(strVal)
		if e == nil {
			return intVal, nil
		}
	}
	return 0, e
}

func getContrast() int {
	command := fmt.Sprintf("ddcutil getvcp 12 --bus=%v", *busNum)
	output, _ := getCommandOutput(command)
	lines := strings.Split(output, "\n")
	lineWithValue := ""
	for _, line := range lines {
		if strings.Contains(line, "Contrast") {
			lineWithValue = strings.Split(line, ",")[0]
		}
	}
	if lineWithValue != "" {
		parts := strings.Split(lineWithValue, " ")
		strVal := parts[len(parts)-1]
		intVal, err := strconv.Atoi(strVal)
		if err == nil {
			return intVal
		}
	}
	return 0
}

func getActivePreset() (string, error) {
	command := fmt.Sprintf("ddcutil getvcp 14 --bus=%v", *busNum)
	output, err := getCommandOutput(command)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Select color preset") {
			values := strings.Split(line, ":")[1]
			preset := strings.Split(values, ",")[0]
			xPos := strings.Index(preset, "x")
			if len(preset) >= xPos+3 {
				return preset[xPos-1 : xPos+3], nil
			}
		}
	}
	return "", err
}

func getPresets() (name string, presets []string, e error) {
	command := fmt.Sprintf("ddcutil capabilities --bus=%v", *busNum)
	output, err := getCommandOutput(command)
	if err == nil {
		lines := strings.Split(output, "\n")

		// name = strings.Split(lines[0], " ")[1]

		here := -1
		for i, line := range lines {
			if strings.HasPrefix(line, "Model:") {
				name = strings.Split(line, " ")[1]
			}
			if strings.Contains(line, "Feature: 14") {
				here = i
				log.Debugf("Feature 14 found in line %v", i)
				break
			}
		}

		if here != -1 {
			for _, line := range lines[here+2:] {
				if !strings.Contains(line, "Feature") {
					presets = append(presets, strings.TrimSpace(line))
				} else {
					break
				}
			}
			return name, presets, nil
		} else {
			return "Unrecognized", presets, errors.New("error parsing capabilities")
		}

	}
	return "unknown", nil, err
}

func launch(command string) {
	log.Debugf("Executing: %s", command)
	tStart := time.Now()
	parts := strings.Split(command, " ")

	cmd := exec.Command(parts[0], parts[1:]...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}
	if cmd.Start() != nil {
		log.Warnf("Couldn't execute: %s", command)
	} else {
		go func() {
			_ = cmd.Wait()
			t := time.Now()
			log.Debugf("It took %v ms.", t.Sub(tStart).Milliseconds())
		}()
	}
}
