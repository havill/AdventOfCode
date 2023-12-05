package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
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

func mapSrcToDest(src int, ranges MappingList) int {
	dest := src

	//fmt.Println("src = ", src)
	//fmt.Println("ranges = ", ranges)

	for _, i := range ranges {
		lo := i.Source
		hi := i.Source + i.Length
		// fmt.Println("lo = ", lo, ", hi = ", hi)
		if src >= lo && src < hi {
			// fmt.Println("src ", src, " is in between ", lo, " and ", hi-1)
			dest = i.Destination + src - lo
			break
		}
	}
	//fmt.Fprintln(os.Stderr, "src = ", src, ", dest = ", dest)
	return dest
}

func mapSeedToLocation(a int, s, f, w, l, t, h, z MappingList) int {
	soil := mapSrcToDest(a, s)
	fertilizer := mapSrcToDest(soil, f)
	water := mapSrcToDest(fertilizer, w)
	light := mapSrcToDest(water, l)
	temperature := mapSrcToDest(light, t)
	humidity := mapSrcToDest(temperature, h)
	location := mapSrcToDest(humidity, z)
	return location
}

func main() {
	var parserState State = Seeds

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
		if len(line) == 0 {
			continue // skip blank lines
		}
		// fmt.Fprintf(os.Stderr, "line: %s\n", line)
		left, right, error := splitByColon(line)
		// fmt.Fprintln(os.Stderr, "left:", left)
		// fmt.Fprintln(os.Stderr, "right:", right)
		// fmt.Fprintln(os.Stderr, "error:", error)
		if error == nil {
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
					//fmt.Fprintln(os.Stderr, "added seed: ", num)
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
			// fmt.Fprintln(os.Stderr, "parserState = ", parserState)
			continue
		}
		var categoryMap Category

		if parserState != Seeds {
			_, err := fmt.Sscanf(line, "%d %d %d", &categoryMap.Destination, &categoryMap.Source, &categoryMap.Length)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error scanning integers: ", err)
				os.Exit(1)
			}
			switch parserState {
			case Soil:
				soilMaps = append(soilMaps, categoryMap)
				//fmt.Fprintln(os.Stderr, "added Soil map: ", categoryMap)
			case Fertilizer:
				fertilizerMaps = append(fertilizerMaps, categoryMap)
				//fmt.Fprintln(os.Stderr, "added Fertilizer map: ", categoryMap)
			case Water:
				waterMaps = append(waterMaps, categoryMap)
				//fmt.Fprintln(os.Stderr, "added Water map: ", categoryMap)
			case Light:
				lightMaps = append(lightMaps, categoryMap)
				//fmt.Fprintln(os.Stderr, "added Light map: ", categoryMap)
			case Temperature:
				temperatureMaps = append(temperatureMaps, categoryMap)
				//fmt.Fprintln(os.Stderr, "added Temperature map: ", categoryMap)
			case Humidity:
				humidityMaps = append(humidityMaps, categoryMap)
				//fmt.Fprintln(os.Stderr, "added Humidity map: ", categoryMap)
			case Location:
				locationMaps = append(locationMaps, categoryMap)
				//fmt.Fprintln(os.Stderr, "added Location map: ", categoryMap)
			default:
				fmt.Fprintln(os.Stderr, "Unknown State")
				os.Exit(1)
			}
		}
	}

	//fmt.Println(mapSeedToLocation(79, soilMaps, fertilizerMaps, waterMaps, lightMaps, temperatureMaps, humidityMaps, locationMaps))
	//fmt.Println(mapSeedToLocation(14, soilMaps, fertilizerMaps, waterMaps, lightMaps, temperatureMaps, humidityMaps, locationMaps))
	//fmt.Println(mapSeedToLocation(55, soilMaps, fertilizerMaps, waterMaps, lightMaps, temperatureMaps, humidityMaps, locationMaps))
	//fmt.Println(mapSeedToLocation(13, soilMaps, fertilizerMaps, waterMaps, lightMaps, temperatureMaps, humidityMaps, locationMaps))

	var lowest int = math.MaxInt

	for _, num := range toBePlanted {

		//fmt.Println("Index:", i, "Number:", num)
		x := mapSeedToLocation(num, soilMaps, fertilizerMaps, waterMaps, lightMaps, temperatureMaps, humidityMaps, locationMaps)
		if x < lowest {
			lowest = x
		}
	}
	fmt.Println(lowest)

	lowest = math.MaxInt
	for i := 0; i < len(toBePlanted); i += 2 {
		lo := toBePlanted[i]
		hi := toBePlanted[i] + toBePlanted[i+1]
		// fmt.Println("lo = ", lo, ", hi = ", hi)
		for j := lo; j < hi; j++ {
			x := mapSeedToLocation(j, soilMaps, fertilizerMaps, waterMaps, lightMaps, temperatureMaps, humidityMaps, locationMaps)
			if x < lowest {
				lowest = x
			}
		}
	}
	fmt.Println(lowest)

}
