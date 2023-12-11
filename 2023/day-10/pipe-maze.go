package main

import (
	"bufio"
	"fmt"
	"os"
)

type Tile int

const (
	unknown Tile = iota
	gt
	st
	ns
	ew
	ne
	nw
	sw
	se
)

func charToTile(c rune) Tile {
	switch c {
	case '.':
		return gt
	case 'S':
		return st
	case '|':
		return ns
	case '-':
		return ew
	case 'L':
		return ne
	case 'J':
		return nw
	case '7':
		return sw
	case 'F':
		return se
		/*
			default:
				return unknown // default case if none of the above match
		*/
	}
	fmt.Fprintf(os.Stderr, "charToTile: unknown character %c\n", c)
	os.Exit(1)
	return unknown
}

func stringToTileRow(s string) []Tile {
	tiles := make([]Tile, 0, len(s))
	for _, c := range s {
		tile := charToTile(c)
		tiles = append(tiles, tile)
	}
	return tiles
}

func main() {
	var field [][]Tile

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		row := stringToTileRow(line)
		field = append(field, row)
		//fmt.Println(row)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
