package main

import (
	"bufio"
	"fmt"
	"os"
)

type Tile int

const (
	ground      Tile = 0
	north       Tile = 1
	south       Tile = 2
	east        Tile = 4
	west        Tile = 8
	footprint   Tile = 16
	north_south Tile = north | south
	east_west   Tile = east | west
	north_east  Tile = north | east
	north_west  Tile = north | west
	south_east  Tile = south | east
	south_west  Tile = south | west
	starting    Tile = north | south | east | west | footprint
)

type Point struct {
	x int
	y int
}

type Animal struct {
	previous Point
	current  Point
	next     Point
}

func availableDirections(field [][]Tile, p Point) []Point {
	directions := make([]Point, 0, 4)
	if field[p.y][p.x]&north == north && p.y > 0 && field[p.y-1][p.x]&south == south {
		directions = append(directions, Point{p.x, p.y - 1})
	}
	if field[p.y][p.x]&south == south && p.y < len(field)-1 && field[p.y+1][p.x]&north == north {
		directions = append(directions, Point{p.x, p.y + 1})
	}
	if field[p.y][p.x]&east == east && p.x < len(field[p.y])-1 && field[p.y][p.x+1]&west == west {
		directions = append(directions, Point{p.x + 1, p.y})
	}
	if field[p.y][p.x]&west == west && p.x > 0 && field[p.y][p.x-1]&east == east {
		directions = append(directions, Point{p.x - 1, p.y})
	}
	return directions
}

func removePreviousDirection(directions []Point, p Animal) []Point {
	for i, direction := range directions {
		if direction == p.previous {
			// Remove the direction from the slice
			directions = append(directions[:i], directions[i+1:]...)
			break
		}
	}
	return directions
}

func pointsEqual(p1, p2 Point) bool {
	return p1.x == p2.x && p1.y == p2.y
}

func allTogether(animals []Animal) bool {
	for i := 0; i < len(animals)-1; i++ {
		if !pointsEqual(animals[i].current, animals[i+1].current) {
			return false
		}
	}
	return true
}

func moveAnimal(a Animal) Animal {
	a.previous = a.current
	a.current = a.next
	a.next = Point{-1, -1}
	return a
}

func charToTile(c rune) Tile {
	switch c {
	case '.':
		return ground
	case 'S':
		return starting
	case '|':
		return north_south
	case '-':
		return east_west
	case 'L':
		return north_east
	case 'J':
		return north_west
	case '7':
		return south_west
	case 'F':
		return south_east
	}
	fmt.Fprintf(os.Stderr, "charToTile: unknown character %c\n", c)
	return 0
}

func stringToTileRow(s string) []Tile {
	tiles := make([]Tile, 0, len(s))
	for _, c := range s {
		tile := charToTile(c)
		tiles = append(tiles, tile)
	}
	return tiles
}

func findStartingTile(field [][]Tile) Point {
	for y, row := range field {
		for x, tile := range row {
			if tile == starting {
				return Point{x, y}
			}
		}
	}
	return Point{-1, -1} // return Point{-1, -1} if "starting" is not found
}

func main() {
	var field [][]Tile
	var animals []Animal
	var start Point
	var distance int

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		row := stringToTileRow(line)
		field = append(field, row)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	start = findStartingTile(field)
	distance = 0

	choices := availableDirections(field, start)
	for i := 0; i < len(choices); i++ {
		animals = append(animals, Animal{start, choices[i], Point{-1, -1}})
		field[animals[i].current.y][animals[i].current.x] |= footprint
	}
	distance++

	for !allTogether(animals) {
		for i := 0; i < len(animals); i++ {
			animals[i].next = animals[i].current
			choices := availableDirections(field, animals[i].current)
			choices = removePreviousDirection(choices, animals[i])
			if len(choices) > 0 {
				animals[i].next = choices[0]
			}
			animals[i] = moveAnimal(animals[i])
			field[animals[i].current.y][animals[i].current.x] |= footprint
		}
		distance++
	}
	fmt.Println(distance)
}
