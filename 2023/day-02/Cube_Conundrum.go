package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func splitByColon(s string) (string, string, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return "", "", errors.New("input string does not contain a colon")
	}
	return parts[0], parts[1], nil
}

func extractInt(s string) (int, error) {
	re := regexp.MustCompile(`\d+`)
	match := re.FindStringSubmatch(s)
	if match == nil {
		return 0, errors.New("no integer found in string")
	}
	return strconv.Atoi(match[0])
}

func extractIntBeforeColor(set string, color string) (int, error) {
	colorIndex := strings.Index(set, color)
	if colorIndex == -1 {
		return 0, errors.New("color not found in set")
	}
	commaIndex := strings.LastIndex(set[:colorIndex], ",")
	if commaIndex == -1 {
		return extractInt(set)
	}
	substring := set[commaIndex+1 : colorIndex]
	return extractInt(substring)
}

func main() {
	sum := 0
	power_sum := 0
	max_red := 12
	max_green := 13
	max_blue := 14

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		impossible := false
		min_red := 0
		min_green := 0
		min_blue := 0

		game, revealed, _ := splitByColon(line)
		id, _ := extractInt(game)
		tokens := strings.Split(revealed, ";")
		for _, set := range tokens {
			red, _ := extractIntBeforeColor(set, "red")
			green, _ := extractIntBeforeColor(set, "green")
			blue, _ := extractIntBeforeColor(set, "blue")
			if red > max_red || green > max_green || blue > max_blue {
				impossible = true
			}
			if red > min_red {
				min_red = red
			}
			if green > min_green {
				min_green = green
			}
			if blue > min_blue {
				min_blue = blue
			}
		}
		power := min_red * min_green * min_blue
		if !impossible {
			sum += id
		}
		power_sum += power
	}
	fmt.Println(sum)
	fmt.Println(power_sum)
}
