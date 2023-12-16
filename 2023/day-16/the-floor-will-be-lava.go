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

type beamList []beam

type headingType rune

const (
	north headingType = '^'
	south headingType = 'v'
	east  headingType = '>'
	west  headingType = '<'
)

type beam struct {
	yAdvance verticalDirection
	xAdvance horizontalDirection
	x        int
	y        int
}

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

func spawnBeam(beams beamList, x int, y int, going headingType) beamList {
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

	beams = append(beams, newBeam)
	return beams
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

func findBeamsAtPosition(beams beamList, x, y int) []beam {
	var matchingBeams []beam
	for _, b := range beams {
		if b.x == x && b.y == y {
			matchingBeams = append(matchingBeams, b)
		}
	}
	return matchingBeams
}

func debugDiagram(grid gridMatrix, beams beamList) {
	// let's animate!
	fmt.Printf("%s", "\033[H") // move cursor to top left
	for y, row := range grid {
		for x, space := range row {
			beamAtPosition := findBeamsAtPosition(beams, x, y)
			if beamAtPosition != nil {
				if len(beamAtPosition) > 1 {
					fmt.Printf("%s", "\033[48;2;255;0;0m") // set the background to red
					fmt.Printf("%s", "\033[38;2;0;0;0m")   // set the foreground to black
					fmt.Print("*")
					fmt.Printf("%s", "\033[0m") // restore the background and foreground colors to defaults
				} else {
					c := rune(toArrow(int(beamAtPosition[0].xAdvance), int(beamAtPosition[0].yAdvance)))
					fmt.Printf("%s", "\033[48;2;255;255;0m") // set the background to yellow
					fmt.Printf("%s", "\033[38;2;0;0;255m")   // set the foreground to blue
					fmt.Printf("%c", c)
					fmt.Printf("%s", "\033[0m") // restore the background and foreground colors to defaults
				}
			} else {
				if space.energized > 0 {
					fmt.Printf("%s", "\033[48;2;255;255;0m") // set the background to yellow
					fmt.Printf("%s", "\033[38;2;0;0;0m")     // set the foreground to black
				}
				fmt.Print(string(space.containing))
				fmt.Printf("%s", "\033[0m") // restore the background and foreground colors to defaults
			}
		}
		fmt.Println() // end of row
	}
	// pause until they press Enter
	//reader := bufio.NewReader(os.Stdin)
	//reader.ReadString('\n')
}

func heatTiles(grid gridMatrix, beams beamList) {
	for _, b := range beams {
		grid[b.y][b.x].energized++
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

func isBeamInHistory(history beamList, b beam) bool {
	for _, old := range history {
		if old == b {
			return true
		}
	}
	return false
}

func gcBeams(width, height int, beams, history beamList) beamList {
	var newBeams beamList
	for _, b := range beams {
		if b.x >= 0 && b.x < width && b.y >= 0 && b.y < height {
			if !isBeamInHistory(history, b) {
				newBeams = append(newBeams, b)
			}
		}
	}
	return newBeams
}

func advanceBeams(beams, history beamList) beamList {
	for i := range beams {
		history = append(history, beams[i])
		beams[i].x += int(beams[i].xAdvance)
		beams[i].y += int(beams[i].yAdvance)
	}
	return history
}

func deflectOrSplitBeams(grid gridMatrix, beams beamList) beamList {
	newBeams := make(beamList, 0, len(beams))

	for i := range beams {
		tile := grid[beams[i].y][beams[i].x]
		if tile.containing != emptySpace {
			if tile.containing == forwardMirror {
				if beams[i].xAdvance > 0 {
					beams[i].xAdvance = 0
					beams[i].yAdvance = up
				} else if beams[i].xAdvance < 0 {
					beams[i].xAdvance = 0
					beams[i].yAdvance = down
				} else if beams[i].yAdvance > 0 {
					beams[i].yAdvance = 0
					beams[i].xAdvance = left
				} else if beams[i].yAdvance < 0 {
					beams[i].yAdvance = 0
					beams[i].xAdvance = right
				}
			} else if tile.containing == backwardMirror {
				if beams[i].xAdvance > 0 {
					beams[i].xAdvance = 0
					beams[i].yAdvance = down
				} else if beams[i].xAdvance < 0 {
					beams[i].xAdvance = 0
					beams[i].yAdvance = up
				} else if beams[i].yAdvance < 0 {
					beams[i].yAdvance = 0
					beams[i].xAdvance = left
				} else if beams[i].yAdvance > 0 {
					beams[i].yAdvance = 0
					beams[i].xAdvance = right
				}
			} else if tile.containing == verticalSplitter {
				if beams[i].yAdvance == 0 {
					newBeams = spawnBeam(newBeams, beams[i].x, beams[i].y, north)
					beams[i].xAdvance = 0
					beams[i].yAdvance = down
				}
			} else if tile.containing == horizontalSplitter {
				if beams[i].xAdvance == 0 {
					newBeams = spawnBeam(newBeams, beams[i].x, beams[i].y, east)
					beams[i].xAdvance = left
					beams[i].yAdvance = 0
				}
			}
		}
	}
	return append(beams, newBeams...)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Please provide a file name")
	}

	fmt.Printf("%s", "\033[2J") // clear screen

	file := os.Args[1]
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer f.Close()

	grid, _ := loadGridFromFile(*f)
	x, y := gridDimensions(grid)

	var beams, history beamList

	beams = spawnBeam(beams, 0, 0, east)
	for len(beams) > 0 {
		debugDiagram(grid, beams) // debug

		beams = deflectOrSplitBeams(grid, beams)
		heatTiles(grid, beams)
		history = advanceBeams(beams, history)
		beams = gcBeams(x, y, beams, history)
		//fmt.Println(beams) // debug
	}

	count := energizedTiles(grid)
	fmt.Printf("Number of energized tiles: %d\n", count)

}
