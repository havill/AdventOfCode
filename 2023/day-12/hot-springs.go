package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func groupMatchesRecord(re *regexp.Regexp, conditionRecord string) bool {
	match := re.MatchString(conditionRecord)
	fmt.Println("groupMatchesRecord: ", conditionRecord, match)
	return match
}

func testCombosHelper(re *regexp.Regexp, brokenGroupsSum int, conditionRecord string, current string) int {
	total := 0

	if strings.Count(current, "#") > brokenGroupsSum {
		return total // short-circuit the recursion if we've already exceeded the number of broken groups
	}
	if len(conditionRecord) == 0 {
		if groupMatchesRecord(re, current) {
			//fmt.Println(current)
			total++
		}
		return total
	}
	if conditionRecord[0] == '?' {
		total += testCombosHelper(re, brokenGroupsSum, conditionRecord[1:], current+".")
		total += testCombosHelper(re, brokenGroupsSum, conditionRecord[1:], current+"#")
	} else {
		total += testCombosHelper(re, brokenGroupsSum, conditionRecord[1:], current+string(conditionRecord[0]))
	}
	return total
}

func testAllCombos(re *regexp.Regexp, brokenGroupsSum int, conditionRecord string) int {
	return testCombosHelper(re, brokenGroupsSum, conditionRecord, "")
}

func splitString(input string) (string, string) {
	parts := strings.Split(input, " ")
	left := ""
	right := ""

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if strings.ContainsAny(part, ".#?") {
			left += part
		} else {
			right += part
		}
	}

	return left, right
}

func stringToIntSlice(input string) ([]int, error) {
	parts := strings.Split(input, ",")
	result := make([]int, len(parts))

	for i, part := range parts {
		part = strings.TrimSpace(part)
		num, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		result[i] = num
	}

	return result, nil
}

func convertSliceToString(defectives []int) string {
	var result []string

	result = append(result, "^[^#]*")
	for i, defective := range defectives {
		result = append(result, fmt.Sprintf("[#?]{%d}", defective))
		if i < len(defectives)-1 {
			result = append(result, "[^#]+")
		}
	}
	result = append(result, "[^#]*$")
	return strings.Join(result, "")
}

func sumArray(numbers []int) int {
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	return sum
}

func repeatString(conditions string, folds int) string {
	repeatedConditions := strings.Repeat(conditions+"?", folds)
	return strings.TrimSuffix(repeatedConditions, "?") // Remove the trailing "?"
}

func repeatIntSlice(contiguousGroups []int, folds int) []int {
	result := make([]int, 0, len(contiguousGroups)*folds)
	for i := 0; i < folds; i++ {
		result = append(result, contiguousGroups...)
	}
	return result
}

func main() {
	total := 0
	unfoldedTotal := 0
	scanner := bufio.NewScanner(os.Stdin)
	progress := 0

	//fmt.Println("Please enter some lines of text. Press CTRL+D to end input.")
	for scanner.Scan() {
		var sum int
		var pattern string
		var re *regexp.Regexp

		line := scanner.Text()
		progress++
		fmt.Println("processing line #", progress, "at", time.Now())
		left, right := splitString(line)
		//fmt.Println("Left: ", left)
		//fmt.Println("Right: ", right)
		brokenGroups, err := stringToIntSlice(right)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		//fmt.Println("Broken groups: ", brokenGroups)
		sum = sumArray(brokenGroups)
		//fmt.Println("Sum: ", sum)
		pattern = convertSliceToString(brokenGroups)
		//fmt.Println("Uncompiled: ", pattern)
		re, err = regexp.Compile(pattern)
		if err != nil {
			fmt.Println("Error compiling regex:", err)
			return
		}
		total += testAllCombos(re, sum, left)

		// Part 2
		unfoldedConditions := repeatString(left, 5)
		unfoldedBrokenGroups := repeatIntSlice(brokenGroups, 5)
		//fmt.Println("Unfolded conditions: ", unfoldedConditions)
		//fmt.Println("Unfolded broken groups: ", unfoldedBrokenGroups)
		sum = sumArray(unfoldedBrokenGroups)
		//fmt.Println("Sum: ", sum)
		pattern = convertSliceToString(unfoldedBrokenGroups)
		//fmt.Println("Uncompiled: ", pattern)
		re, _ = regexp.Compile(pattern)
		matches := testAllCombos(re, sum, unfoldedConditions)
		//fmt.Println("Matches: ", matches)
		unfoldedTotal += matches
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	fmt.Println("Total matches (part 1): ", total)
	fmt.Println("Total matches (part 2): ", unfoldedTotal)
}
