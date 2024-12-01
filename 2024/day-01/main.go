package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func sortSlice(slice []int) {
	sort.Ints(slice)
}

func distanceList(leftList, rightList []int) []int {
	if len(leftList) != len(rightList) {
		fmt.Println("Error: Lists are not of the same length")
		return nil
	}

	var distances []int
	for i := 0; i < len(leftList); i++ {
		distances = append(distances, int(math.Abs(float64(leftList[i]-rightList[i]))))
	}
	return distances
}

func sumSlice(slice []int) int {
	sum := 0
	for _, value := range slice {
		sum += value
	}
	return sum
}

func similarityScore(leftList, rightList []int) []int {
	var similarityScore []int
	for _, leftValue := range leftList {
		appearances := 0
		for _, rightValue := range rightList {
			if leftValue == rightValue {
				appearances++
			}
		}
		similarityScore = append(similarityScore, leftValue*appearances)
	}
	return similarityScore
}

func main() {
	var leftList []int
	var rightList []int

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		numbers := strings.Fields(line)
		if len(numbers) != 2 {
			fmt.Println("Each line must contain exactly two integers")
			continue
		}

		left, err1 := strconv.Atoi(numbers[0])
		right, err2 := strconv.Atoi(numbers[1])
		if err1 != nil || err2 != nil {
			fmt.Println("Invalid input, please enter integers")
			continue
		}

		leftList = append(leftList, left)
		rightList = append(rightList, right)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}

	sortSlice(leftList)
	sortSlice(rightList)

	//fmt.Println("Left List:", leftList)
	//fmt.Println("Right List:", rightList)

	distances := distanceList(leftList, rightList)
	//fmt.Println("Distances:", distances)

	sum := sumSlice(distances)
	fmt.Println("Sum of Distances:", sum)

	similarity := similarityScore(leftList, rightList)
	//fmt.Println("Similarity Score:", similarity)

	totalSimilarityScore := sumSlice(similarity)
	fmt.Println("Total Similarity Score:", totalSimilarityScore)
}
