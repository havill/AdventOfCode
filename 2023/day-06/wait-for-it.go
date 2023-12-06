package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Races struct {
	TimeMS     []int
	DistanceMM []int
}

func splitByColon(s string) (string, string, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return "", "", errors.New("input string does not contain a colon")
	}
	return parts[0], parts[1], nil
}

func solve(r Races) int {
	answer := 1
	for i := 0; i < len(r.TimeMS); i++ {
		winners := 0
		for buttonMS := 1; buttonMS < r.TimeMS[i]-1; buttonMS++ {
			if buttonMS*(r.TimeMS[i]-buttonMS) > r.DistanceMM[i] {
				winners++
				//fmt.Fprintln(os.Stderr, "winner = ", buttonMS, "ms")
			}
		}
		answer *= winners
	}
	return answer
}

func main() {
	var r [2]Races

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
				nospace := strings.Join(strings.Fields(str), "")
				fmt.Fprintln(os.Stderr, "nospace = ", nospace)
				for _, field := range fields {
					num, err := strconv.Atoi(field)
					if err != nil {
						fmt.Fprintln(os.Stderr, "Error converting string to integer: ", err)
						os.Exit(1)
					}
					r[0].TimeMS = append(r[0].TimeMS, num)
					//fmt.Fprintln(os.Stderr, "added time: ", num)
				}
				num, _ := strconv.Atoi(nospace)
				r[1].TimeMS = append(r[1].TimeMS, num)
			} else if strings.EqualFold(left, "distance") {
				str := right
				fields := strings.Fields(str)
				nospace := strings.Join(strings.Fields(str), "")
				fmt.Fprintln(os.Stderr, "nospace = ", nospace)
				for _, field := range fields {
					num, err := strconv.Atoi(field)
					if err != nil {
						fmt.Fprintln(os.Stderr, "Error converting string to integer: ", err)
						os.Exit(1)
					}
					r[0].DistanceMM = append(r[0].DistanceMM, num)
					//fmt.Fprintln(os.Stderr, "added distance: ", num)
				}
				num, _ := strconv.Atoi(nospace)
				r[1].DistanceMM = append(r[1].DistanceMM, num)
			} else {
				fmt.Fprintln(os.Stderr, "Unknown map")
				os.Exit(1)
			}

		}
	}
	if len(r[0].TimeMS) != len(r[0].DistanceMM) {
		fmt.Fprintln(os.Stderr, "Error: time and distance lists are not the same length")
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "time list: ", r[0].TimeMS)
	fmt.Fprintln(os.Stderr, "distance list: ", r[0].DistanceMM)
	fmt.Fprintln(os.Stderr, "time list: ", r[1].TimeMS)
	fmt.Fprintln(os.Stderr, "distance list: ", r[1].DistanceMM)

	fmt.Println(solve(r[0]))
	fmt.Println(solve(r[1]))

}
