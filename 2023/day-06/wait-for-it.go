package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Race struct {
	Time     []int
	Distance []int
}

type RaceList []Race

func splitByColon(s string) (string, string, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return "", "", errors.New("input string does not contain a colon")
	}
	return parts[0], parts[1], nil
}

func main() {
	var RaceList Race

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue // skip blank lines
		}
		left, right, error := splitByColon(line)

		if error == nil {
			if strings.EqualFold(left, "time") {
				str := right
				fields := strings.Fields(str)

				for _, field := range fields {
					num, err := strconv.Atoi(field)
					if err != nil {
						fmt.Fprintln(os.Stderr, "Error converting string to integer: ", err)
						os.Exit(1)
					}
					RaceList.Time = append(RaceList.Time, num)
					fmt.Fprintln(os.Stderr, "added time: ", num)
				}
			} else if strings.EqualFold(left, "distance") {
				str := right
				fields := strings.Fields(str)

				for _, field := range fields {
					num, err := strconv.Atoi(field)
					if err != nil {
						fmt.Fprintln(os.Stderr, "Error converting string to integer: ", err)
						os.Exit(1)
					}
					RaceList.Distance = append(RaceList.Distance, num)
					fmt.Fprintln(os.Stderr, "added distance: ", num)
				}
			} else {
				fmt.Fprintln(os.Stderr, "Unknown map")
				os.Exit(1)
			}
			continue
		}
	}
	if len(RaceList.Time) != len(RaceList.Distance) {
		fmt.Fprintln(os.Stderr, "Error: time and distance lists are not the same length")
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "time list: ", RaceList.Time)
	fmt.Fprintln(os.Stderr, "distance list: ", RaceList.Distance)

}
