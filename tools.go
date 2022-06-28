package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

func waylandSession() bool {
	return os.Getenv("WAYLAND_DISPLAY") != "" || os.Getenv("XDG_SESSION_TYPE") == "waylandSession"
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
	log.Debugf("Command: %s", command)
	elements := strings.Split(command, " ")
	c, b := exec.Command(elements[0], elements[1:]...), new(strings.Builder)
	c.Stdout = b
	err := c.Run()

	return b.String(), err
}

func getBrightness() int {
	var command string
	if *busNum > -1 {
		command = fmt.Sprintf("ddcutil getvcp 10 --bus=%v", *busNum)
	} else {
		command = "ddcutil getvcp 10"
	}
	output, _ := getCommandOutput(command)
	lines := strings.Split(output, "\n")
	lineWithValue := ""
	for _, line := range lines {
		if strings.Contains(line, "Brightness") {
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

func getContrast() int {
	var command string
	if *busNum > -1 {
		command = fmt.Sprintf("ddcutil getvcp 12 --bus=%v", *busNum)
	} else {
		command = "ddcutil getvcp 12"
	}
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

func getActivePreset() string {
	var command string
	if *busNum > -1 {
		command = fmt.Sprintf("ddcutil getvcp 14 --bus=%v", *busNum)
	} else {
		command = "ddcutil getvcp 14"
	}
	output, _ := getCommandOutput(command)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Select color preset") {
			parts := strings.Split(line, " ")
			part := parts[len(parts)-1]
			if strings.Contains(part, "sl=0x") {
				val := strings.Split(part, "=")[1]
				val = val[:len(val)-1]
				return val
			}
		}
	}
	return ""
}

func getPresets() (name string, presets []string, e error) {
	var command string
	if *busNum > -1 {
		command = fmt.Sprintf("ddcutil capabilities --bus=%v", *busNum)
	} else {
		command = "ddcutil capabilities"
	}
	output, err := getCommandOutput(command)
	if err == nil {
		lines := strings.Split(output, "\n")

		name = strings.Split(lines[0], " ")[1]

		here := -1
		for i, line := range lines {
			if strings.Contains(line, "Feature: 14") {
				here = i
				log.Debugf("Feature 14 found in line %v", i)
				break
			}
		}

		for _, line := range lines[here+2:] {
			if !strings.Contains(line, "Feature") {
				presets = append(presets, strings.TrimSpace(line))
			} else {
				break
			}
		}

		return name, presets, nil
	}
	return "", nil, err
}
