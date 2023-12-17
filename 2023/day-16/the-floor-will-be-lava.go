package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type verticalDirection int
type horizontalDirection int

const (
	up   verticalDirection = -1
	down verticalDirection = +1
)

const (
	left  horizontalDirection = -1
	right horizontalDirection = +1
)

type headingType rune

const (
	north headingType = '^'
	south headingType = 'v'
	east  headingType = '>'
	west  headingType = '<'
)

type beam struct {
	x        int // Col - 1
	y        int // Ln - 1
	xAdvance horizontalDirection
	yAdvance verticalDirection
}

type beamMap map[beam]int

type spaceType rune

const (
	emptySpace         spaceType = '.'
	forwardMirror      spaceType = '/'
	backwardMirror     spaceType = '\\'
	verticalSplitter   spaceType = '|'
	horizontalSplitter spaceType = '-'
)

type tile struct {
	containing spaceType
	energized  int
}

type gridMatrix [][]tile

func toArrow(x, y int) headingType {
	if y > 0 && x == 0 {
		return south
	} else if y < 0 && x == 0 {
		return north
	} else if x > 0 && y == 0 {
		return east
	} else if x < 0 && y == 0 {
		return west
	}
	return '?'
}

func spawnBeam(beams beamMap, x int, y int, going headingType) {
	newBeam := beam{x: x, y: y}

	switch going {
	case north:
		newBeam.yAdvance = up
	case south:
		newBeam.yAdvance = down
	case east:
		newBeam.xAdvance = right
	case west:
		newBeam.xAdvance = left
	}

	beams[newBeam] = 1
}

func loadGridFromFile(file os.File) (gridMatrix, error) {
	scanner := bufio.NewScanner(&file)
	var grid gridMatrix

	for scanner.Scan() {
		line := scanner.Text()
		row := make([]tile, len(line))
		for i, char := range line {
			switch spaceType(char) {
			case emptySpace, forwardMirror, backwardMirror, verticalSplitter, horizontalSplitter:
				row[i] = tile{containing: spaceType(char)}
			default:
				return nil, fmt.Errorf("invalid character %c", char)
			}
		}
		grid = append(grid, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return grid, nil
}

func energizedTiles(grid gridMatrix) int {
	count := 0
	for _, row := range grid {
		for _, space := range row {
			if space.energized > 0 {
				count++
			}
		}
	}
	return count
}

func findBeamsAtPosition(beams beamMap, x, y int) []beam {
	var matchingBeams []beam
	for b := range beams {
		if b.x == x && b.y == y {
			matchingBeams = append(matchingBeams, b)
		}
	}
	return matchingBeams
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

func debugDiagram(grid gridMatrix, beams beamMap, interactive bool) {
	if !interactive {
		resetCursorToTopLeft(interactive)
	}
	for y, row := range grid {
		for x, space := range row {
			beamAtPosition := findBeamsAtPosition(beams, x, y)
			if beamAtPosition != nil {
				if len(beamAtPosition) > 1 {
					setTerminalBackgroundColor(255, 0, 0)   // red
					setTerminalForegroundColor(255, 255, 0) // yellow
					fmt.Print("*")
					resetTerminalColors()
				} else {
					c := rune(toArrow(int(beamAtPosition[0].xAdvance), int(beamAtPosition[0].yAdvance)))
					setTerminalBackgroundColor(255, 255, 0) // yellow
					setTerminalForegroundColor(0, 0, 255)   // blue
					fmt.Printf("%c", c)
					resetTerminalColors()
				}
			} else {
				if space.energized > 0 {
					setTerminalBackgroundColor(255, 255, 0) // yellow
					setTerminalForegroundColor(0, 0, 0)     // black
				}
				fmt.Print(string(space.containing))
				resetTerminalColors()
			}
		}
		fmt.Println() // end of row
	}
	if interactive {
		fmt.Println("beams    =", beams)
		fmt.Printf("energized= %d\n", energizedTiles(grid))
		// pause until they press Enter
		var reader *bufio.Reader = bufio.NewReader(os.Stdin)
		reader.ReadString('\n')
	} else {
		fmt.Println()
	}
}

func heatTiles(grid gridMatrix, beams beamMap) {
	for k, heat := range beams {
		grid[k.y][k.x].energized += heat
	}
}

func gridDimensions(grid gridMatrix) (int, int) {
	height := len(grid)
	width := 0
	if height > 0 {
		width = len(grid[0])
	}
	return width, height
}

func gcBeams(beams, history beamMap, width, height int) beamMap {
	newBeams := make(beamMap)

	for k, v := range beams {
		_, found := history[k]
		if k.x >= 0 && k.x < width && k.y >= 0 && k.y < height && !found {
			newBeams[k] = v
		}
	}
	return newBeams
}

func advanceBeams(beams, history beamMap) beamMap {
	newBeams := make(beamMap)

	for k, v := range beams {
		history[k]++
		k.x += int(k.xAdvance)
		k.y += int(k.yAdvance)
		newBeams[k] = v
	}
	return newBeams
}

func deflectOrSplitBeams(beams beamMap, grid gridMatrix) beamMap {
	newBeams := make(beamMap)

	for k, v := range beams {
		tile := grid[k.y][k.x]
		switch tile.containing {
		case forwardMirror:
			if k.xAdvance > 0 {
				k.xAdvance = 0
				k.yAdvance = up
			} else if k.xAdvance < 0 {
				k.xAdvance = 0
				k.yAdvance = down
			} else if k.yAdvance > 0 {
				k.yAdvance = 0
				k.xAdvance = left
			} else if k.yAdvance < 0 {
				k.yAdvance = 0
				k.xAdvance = right
			}
		case backwardMirror:
			if k.xAdvance > 0 {
				k.xAdvance = 0
				k.yAdvance = down
			} else if k.xAdvance < 0 {
				k.xAdvance = 0
				k.yAdvance = up
			} else if k.yAdvance < 0 {
				k.yAdvance = 0
				k.xAdvance = left
			} else if k.yAdvance > 0 {
				k.yAdvance = 0
				k.xAdvance = right
			}
		case verticalSplitter:
			if k.yAdvance == 0 {
				spawnBeam(newBeams, k.x, k.y, north)
				k.xAdvance = 0
				k.yAdvance = down
			}
		case horizontalSplitter:
			if k.xAdvance == 0 {
				spawnBeam(newBeams, k.x, k.y, east)
				k.xAdvance = left
				k.yAdvance = 0
			}
		}
		newBeams[k] = v
	}
	return newBeams
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Please provide a file name")
	}

	resetCursorToTopLeft(true) // for interactive debug

	file := os.Args[1]
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer f.Close()

	grid, _ := loadGridFromFile(*f)
	x, y := gridDimensions(grid)

	beams := make(beamMap)
	history := make(beamMap)

	spawnBeam(beams, 0, 0, east)
	for len(beams) > 0 {
		beams = gcBeams(beams, history, x, y)

		debugDiagram(grid, beams, false) // debug

		heatTiles(grid, beams)
		beams = deflectOrSplitBeams(beams, grid)
		beams = advanceBeams(beams, history)

	}
	count := energizedTiles(grid)
	fmt.Printf("Part 1: Number of energized tiles = %d\n", count)
}
