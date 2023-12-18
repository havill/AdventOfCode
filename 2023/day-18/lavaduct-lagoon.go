package main

import (
	"fmt"
	"strconv"
	"strings"
)

type rgba struct {
	red, green, blue, alpha int64
}

type Ground struct {
	Hole  bool
	color rgba
}

func extractRGB(hex24 string) (rgba, error) {
	var color rgba
	var err error

	hex24 = strings.TrimSpace(hex24)
	if strings.HasPrefix(hex24, "#") {
		hex24 = hex24[1:]
	}
	if len(hex24) != 6 {
		return color, fmt.Errorf("invalid color length")
	}

	color.red, err = strconv.ParseInt(hex24[0:2], 16, 64)
	if err != nil {
		return color, err
	}

	color.green, err = strconv.ParseInt(hex24[2:4], 16, 64)
	if err != nil {
		return color, err
	}

	color.blue, err = strconv.ParseInt(hex24[4:6], 16, 64)
	if err != nil {
		return color, err
	}

	return color, nil
}

func parseDirection(s string) (horizontalChange, verticalChange int) {
	switch strings.ToUpper(s) {
	case "U":
		return 0, -1
	case "D":
		return 0, 1
	case "L":
		return -1, 0
	case "R":
		return 1, 0
	}
	return 0, 0
}

func parseDigPlan(line string) (direction string, meters int, color string, err error) {
	fields := strings.Fields(line)
	if len(fields) != 3 {
		return "", 0, "", fmt.Errorf("invalid dig plan format")
	}

	direction = fields[0]

	meters, err = strconv.Atoi(fields[1])
	if err != nil {
		return "", 0, "", fmt.Errorf("invalid meters value: %v", err)
	}

	color = fields[2]
	if !strings.HasPrefix(color, "#") {
		return "", 0, "", fmt.Errorf("invalid color format")
	}

	return direction, meters, color, nil
}

func main() {

}
