package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
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

type coordinate struct {
	x, y int
}

type graph struct {
	minX, minY int
	maxX, maxY int
	cube       map[coordinate]Ground
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
	case "U", "3":
		return 0, -1
	case "D", "1":
		return 0, 1
	case "L", "2":
		return -1, 0
	case "R", "0":
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

	color = strings.Trim(fields[2], "()")
	if !strings.HasPrefix(color, "#") {
		return "", 0, "", fmt.Errorf("invalid color format")
	}

	return direction, meters, color, nil
}

func dig(where coordinate, color rgba, lagoon graph) graph {
	// Update the min/max X and Y values.
	if where.x < lagoon.minX {
		lagoon.minX = where.x
	}
	if where.x > lagoon.maxX {
		lagoon.maxX = where.x
	}
	if where.y < lagoon.minY {
		lagoon.minY = where.y
	}
	if where.y > lagoon.maxY {
		lagoon.maxY = where.y
	}
	lagoon.cube[where] = Ground{true, color}

	return lagoon
}

func setTerminalForegroundColor(r, g, b int) {
	fmt.Printf("\033[38;2;%d;%d;%dm", r, g, b)
}
func setTerminalBackgroundColor(r, g, b int) {
	fmt.Printf("\033[48;2;%d;%d;%dm", r, g, b)
}
func resetTerminalColors() {
	fmt.Printf("\033[0m")
}
func resetCursorToTopLeft(clearScreen bool) {
	fmt.Printf("%s", "\033[H") // move cursor to Ln 1, Col 1
	if clearScreen {
		fmt.Printf("%s", "\033[2J") // clear screen
	}
}
func debugPrintLagoon(lagoon graph) {
	for y := lagoon.minY; y <= lagoon.maxY; y++ {
		for x := lagoon.minX; x <= lagoon.maxX; x++ {
			ground := lagoon.cube[coordinate{x, y}]
			setTerminalBackgroundColor(0, 0, 0)
			if ground.Hole {
				setTerminalForegroundColor(int(ground.color.red), int(ground.color.green), int(ground.color.blue))
				fmt.Printf("#")
			} else {
				setTerminalForegroundColor(255, 255, 255)
				fmt.Printf(".")
			}
			resetTerminalColors()
		}
		fmt.Println()
	}
	fmt.Println()
}

func fillPolygonStack(area graph, color rgba, x int, y int) {
	stack := []coordinate{{x, y}}

	for len(stack) > 0 {
		// Pop a coordinate from the stack.
		where := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// Check if the coordinate is out of bounds or if it's an edge pixel.
		pixel := area.cube[where]
		if where.x < area.minX || where.y < area.minY || where.y > area.maxY || where.x > area.maxX || pixel.Hole {
			continue
		}

		// Fill this pixel.
		area.cube[where] = Ground{true, color}

		// Push the neighboring coordinates to the stack.
		stack = append(stack, coordinate{where.x - 1, where.y}) // Left
		stack = append(stack, coordinate{where.x + 1, where.y}) // Right
		stack = append(stack, coordinate{where.x, where.y - 1}) // Up
		stack = append(stack, coordinate{where.x, where.y + 1}) // Down
	}
}

func fillPolygon(area graph, color rgba, x int, y int) {
	where := coordinate{x, y}
	pixel := area.cube[where]
	if x < area.minX || y < area.minY || y > area.maxY || x > area.maxX || pixel.Hole {
		// We're out of bounds or this is an edge pixel, so we return without filling.
		return
	}

	// Fill this pixel.
	area.cube[where] = Ground{true, color}

	// Recursively fill the neighboring pixels.
	fillPolygon(area, color, x-1, y) // Left
	fillPolygon(area, color, x+1, y) // Right
	fillPolygon(area, color, x, y-1) // Up
	fillPolygon(area, color, x, y+1) // Down
}

func countHoles(lagoon graph) int {
	count := 0
	for x := lagoon.minX; x <= lagoon.maxX; x++ {
		for y := lagoon.minY; y <= lagoon.maxY; y++ {
			if lagoon.cube[coordinate{x, y}].Hole {
				count++
			}
		}
	}
	return count
}

func decodeHexadecimal(hex string) (int, string, error) {
	left := hex[1:6]
	right := hex[6:]
	result, err := strconv.ParseInt(left, 16, 64)
	if err != nil {
		return 0, "", err
	}
	return int(result), right, nil
}

func main() {
	var lagoon graph
	x, y := 0, 0
	lagoon.cube = make(map[coordinate]Ground)

	part2 := flag.Bool("part2", false, "--- Part Two ---")
	debug := flag.Bool("debug", false, "show the before and after maps")
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		direction, meters, color, err := parseDigPlan(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing dig plan: %v\n", err)
			continue
		}
		xDelta, yDelta := parseDirection(direction)
		rgb, err := extractRGB(color)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error extracting RGB: %v\n", err)
			continue
		}

		if *part2 {
			meters, direction, _ = decodeHexadecimal(color)
			rgb = rgba{255, 255, 255, 0}
		}

		for meters > 0 {
			where := coordinate{x, y}
			lagoon = dig(where, rgb, lagoon)
			y += yDelta
			x += xDelta
			meters--
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
	}

	if *debug {
		resetCursorToTopLeft(true)
		debugPrintLagoon(lagoon)
	}
	//fillPolygon(lagoon, rgba{255, 0, 0, 0}, 1, 1)
	fillPolygonStack(lagoon, rgba{255, 0, 0, 0}, 1, 1)
	if *debug {
		fmt.Println("After filling:")
		debugPrintLagoon(lagoon)
	}
	fmt.Printf("Lava Area: %dm\u00B3\n", countHoles(lagoon))
}
