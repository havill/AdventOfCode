package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type State int

type Category struct {
	Destination int
	Source      int
	Length      int
}

type SeedList []int
type MappingList []Category

const (
	Seeds State = iota
	Soil
	Fertilizer
	Water
	Light
	Temperature
	Humidity
	Location
)

func splitByColon(s string) (string, string, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return "", "", errors.New("input string does not contain a colon")
	}
	return parts[0], parts[1], nil
}

func main() {
	var parserState State = Seeds
	fmt.Fprintln(os.Stderr, parserState)

	var toBePlanted SeedList
	var soilMaps []Category
	var fertilizerMaps []Category
	var waterMaps []Category
	var lightMaps []Category
	var temperatureMaps []Category
	var humidityMaps []Category
	var locationMaps []Category

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		left, right, error := splitByColon(line)
		if error != nil {
			fmt.Fprintln(os.Stderr, "Error: ", error)
			if strings.EqualFold(left, "seeds") {
				parserState = Seeds
				str := right
				fields := strings.Fields(str)

				for _, field := range fields {
					num, err := strconv.Atoi(field)
					if err != nil {
						fmt.Fprintln(os.Stderr, "Error converting string to integer: ", err)
						os.Exit(1)
					}
					toBePlanted = append(toBePlanted, num)
				}
			} else if strings.EqualFold(left, "seed-to-soil map") {
				parserState = Soil
			} else if strings.EqualFold(left, "soil-to-fertilizer map") {
				parserState = Fertilizer
			} else if strings.EqualFold(left, "fertilizer-to-water map") {
				parserState = Water
			} else if strings.EqualFold(left, "water-to-light map") {
				parserState = Light
			} else if strings.EqualFold(left, "light-to-temperature map") {
				parserState = Temperature
			} else if strings.EqualFold(left, "temperature-to-humidity map") {
				parserState = Humidity
			} else if strings.EqualFold(left, "humidity-to-location map") {
				parserState = Location
			} else {
				fmt.Fprintln(os.Stderr, "Unknown map")
				os.Exit(1)
			}

		}
		if len(line) == 0 {
			continue // skip blank lines
		}
		var categoryMap Category

		_, err := fmt.Sscanf(line, "%d %d %d", &categoryMap.Destination, &categoryMap.Source, &categoryMap.Length)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error scanning integers: ", err)
			os.Exit(1)
		}
		switch parserState {
		case Seeds:
			fmt.Println("State is Seeds")
		case Soil:
			soilMaps = append(soilMaps, categoryMap)
			fmt.Fprintln(os.Stderr, "added Soil map: ", categoryMap)
		case Fertilizer:
			fertilizerMaps = append(fertilizerMaps, categoryMap)
			fmt.Fprintln(os.Stderr, "added Fertilizer map: ", categoryMap)
		case Water:
			waterMaps = append(waterMaps, categoryMap)
			fmt.Fprintln(os.Stderr, "added Water map: ", categoryMap)
		case Light:
			lightMaps = append(lightMaps, categoryMap)
			fmt.Fprintln(os.Stderr, "added Light map: ", categoryMap)
		case Temperature:
			temperatureMaps = append(temperatureMaps, categoryMap)
			fmt.Fprintln(os.Stderr, "added Temperature map: ", categoryMap)
		case Humidity:
			humidityMaps = append(humidityMaps, categoryMap)
			fmt.Fprintln(os.Stderr, "added Humidity map: ", categoryMap)
		case Location:
			locationMaps = append(locationMaps, categoryMap)
			fmt.Fprintln(os.Stderr, "added Location map: ", categoryMap)
		default:
			fmt.Println("Unknown State")
		}
	}
}
