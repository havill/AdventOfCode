package main

import (
	"bufio"
	"fmt"
	"os"
)

type Tile rune

const (
	Path   Tile = '.'
	Forest Tile = '#'
	North  Tile = '^'
	East   Tile = '>'
	South  Tile = 'v'
	West   Tile = '<'
)

type Hiked [][]bool

func PrintMap(hikingTrails [][]Tile, stepped Hiked) {
	for y, row := range hikingTrails {
		for x, tile := range row {
			if stepped[y][x] {
				fmt.Print("O")
			} else {
				fmt.Print(string(tile))
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func ReadTileMatix() ([][]Tile, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var matrix [][]Tile

	for scanner.Scan() {
		line := scanner.Text()
		row := make([]Tile, len(line))
		for i, char := range line {
			row[i] = Tile(char)
		}
		matrix = append(matrix, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading standard input: %v", err)
	}

	return matrix, nil
}

func cloneHiked(original Hiked) Hiked {
	clone := make(Hiked, len(original))
	for i := range original {
		clone[i] = make([]bool, len(original[i]))
		copy(clone[i], original[i])
	}
	return clone
}

func CanGoUp(slippery bool, hikingTrails [][]Tile, stepped Hiked, x, y int) bool {
	if y-1 < 0 {
		return false
	}
	if slippery && hikingTrails[y][x] != Path && hikingTrails[y][x] != North {
		return false
	}
	return hikingTrails[y-1][x] != Forest && !stepped[y-1][x]
}

func CanGoRight(slippery bool, hikingTrails [][]Tile, stepped Hiked, x, y int) bool {
	width := len(hikingTrails[0])
	if x+1 >= width {
		return false
	}
	if slippery && hikingTrails[y][x] != Path && hikingTrails[y][x] != East {
		return false
	}
	return hikingTrails[y][x+1] != Forest && !stepped[y][x+1]
}

func CanGoDown(slippery bool, hikingTrails [][]Tile, stepped Hiked, x, y int) bool {
	height := len(hikingTrails)
	if y+1 >= height {
		return false
	}
	if slippery && hikingTrails[y][x] != Path && hikingTrails[y][x] != South {
		return false
	}
	return hikingTrails[y+1][x] != Forest && !stepped[y+1][x]
}

func CanGoLeft(slippery bool, hikingTrails [][]Tile, stepped Hiked, x, y int) bool {
	if x-1 < 0 {
		return false
	}
	if slippery && hikingTrails[y][x] != Path && hikingTrails[y][x] != West {
		return false
	}
	return hikingTrails[y][x-1] != Forest && !stepped[y][x-1]
}

func AtGoal(hikingTrails [][]Tile, x, y int) bool {
	height := len(hikingTrails)
	return y+1 >= height
}

func WalkToBottom(solutions *[]int, slippery bool, hikingTrails [][]Tile, stepped Hiked, steps, x, y int) int {
	width := len(hikingTrails[0])
	height := len(hikingTrails)

	north, east, south, west := 0, 0, 0, 0

	if x < 0 || x >= width || y < 0 || y >= height {
		return 0
	}
	stepped[y][x] = true
	steps++
	//PrintMap(hikingTrails, stepped)
	if CanGoUp(slippery, hikingTrails, stepped, x, y) {
		//fmt.Println("Going up from ", x, y)
		newMap := cloneHiked(stepped)
		north = WalkToBottom(solutions, slippery, hikingTrails, newMap, steps, x, y-1)
	}
	if CanGoRight(slippery, hikingTrails, stepped, x, y) {
		//fmt.Println("Going right from ", x, y)
		newMap := cloneHiked(stepped)
		east = WalkToBottom(solutions, slippery, hikingTrails, newMap, steps, x+1, y)
	}
	if CanGoDown(slippery, hikingTrails, stepped, x, y) {
		//fmt.Println("Going down from ", x, y)
		newMap := cloneHiked(stepped)
		south = WalkToBottom(solutions, slippery, hikingTrails, newMap, steps, x, y+1)
	}
	if CanGoLeft(slippery, hikingTrails, stepped, x, y) {
		//fmt.Println("Going left from ", x, y)
		newMap := cloneHiked(stepped)
		west = WalkToBottom(solutions, slippery, hikingTrails, newMap, steps, x-1, y)
	}
	if north >= east && north >= south && north >= west {
		steps += north
	} else if east >= north && east >= south && east >= west {
		steps += east
	} else if south >= north && south >= east && south >= west {
		steps += south
	} else if west >= north && west >= east && west >= south {
		steps += west
	}
	if AtGoal(hikingTrails, x, y) {
		fmt.Println(steps)
		*solutions = append(*solutions, steps)
		return steps
	}
	return 0
}

func FindStart(hikingTrails [][]Tile) (int, int) {
	for x, tile := range hikingTrails[0] {
		if tile != Forest {
			return x, 0
		}
	}
	return -1, -1
}

func maxInt(arr []int) (max int, err error) {
	if len(arr) == 0 {
		return 0, fmt.Errorf("array is empty")
	}

	max = arr[0]
	for _, value := range arr {
		if value > max {
			max = value
		}
	}
	return max, nil
}

func main() {
	solutions1, solutions2 := []int{}, []int{}
	longestHike := 0

	hikingTrails, err := ReadTileMatix()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading tile matrix: %v\n", err)
		os.Exit(1)
	}

	stepped := make(Hiked, len(hikingTrails))
	for i := range hikingTrails {
		stepped[i] = make([]bool, len(hikingTrails[i]))
	}
	//PrintMap(hikingTrails, stepped)

	x, y := FindStart(hikingTrails)

	WalkToBottom(&solutions1, true, hikingTrails, stepped, -1, x, y)
	longestHike, _ = maxInt(solutions1)
	fmt.Printf("Part 1: %d\n", longestHike)

	WalkToBottom(&solutions2, false, hikingTrails, stepped, -1, x, y)
	longestHike, _ = maxInt(solutions2)
	fmt.Printf("Part 2: %d\n", longestHike)
}
