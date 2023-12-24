package main

import (
	"bufio"
	"fmt"
	"os"
)

type Tile rune

const (
	Path   Tile = '.'
	Forest      = '#'
	North       = '^'
	East        = '>'
	South       = 'v'
	West        = '<'
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
}

func readTileMatrix() ([][]Tile, error) {
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

func main() {
	hikingTrails, err := readTileMatrix()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading tile matrix: %v\n", err)
		os.Exit(1)
	}

	stepped := make(Hiked, len(hikingTrails))
	for i := range hikingTrails {
		stepped[i] = make([]bool, len(hikingTrails[i]))
	}

	PrintMap(hikingTrails, stepped)
}
