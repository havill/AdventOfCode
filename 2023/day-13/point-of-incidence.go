package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type imgs struct {
	img []img
}

func newImages() imgs {
	return imgs{img: make([]img, 0)}
}

type img struct {
	raw []string
}

func newImg() img {
	return img{raw: make([]string, 0)}
}

func (i img) transposed() []string {
	transposed := []string{}

	for x := 0; x < len(i.raw[0]); x++ {
		buf := ""
		for j := len(i.raw) - 1; j >= 0; j-- {
			buf += string(i.raw[j][x])
		}
		transposed = append(transposed, buf)
	}

	return transposed
}

var images imgs

func findReflection(image []string) (int, bool) {
	walkToEdge := func(lower, upper int) bool {
		for lower >= 0 && upper < len(image) && strings.Compare(image[lower], image[upper]) == 0 {
			lower--
			upper++
		}
		lower++
		return (lower == 0 || upper == len(image))
	}

	for i := 0; i < len(image)-1; i++ {
		if strings.Compare(image[i], image[i+1]) == 0 {
			fmt.Fprintln(os.Stderr, "start point: ", image[i], " vs ", image[i+1])
			if walkToEdge(i, i+1) {
				return i, true
			}
		}
	}
	return 0, false
}

func horizontalReflection(image img, fn func(image []string) (int, bool)) (int, bool) {
	return fn(image.raw)
}

func verticalReflection(image img, fn func(image []string) (int, bool)) (int, bool) {
	return fn(image.transposed())
}

func findReflectionWithDifference(image []string) (int, bool) {
	walkToEdge := func(lower, upper int) bool {
		diffs := 0
		for lower >= 0 && upper < len(image) {
			n := numDiffs(image[lower], image[upper])
			fmt.Fprintln(os.Stderr, n, " difference(s) found at line ", lower)
			diffs += n
			lower--
			upper++
		}
		lower++
		return diffs == 1
	}

	for i := 0; i < len(image)-1; i++ {
		fmt.Fprintln(os.Stderr, "start point: ", image[i], " vs ", image[i+1])
		if walkToEdge(i, i+1) {
			return i, true
		}
	}
	return 0, false
}

func numDiffs(a, b string) int {
	diff := 0
	for i := range a {
		if a[i] != b[i] {
			diff++
		}
	}
	return diff
}

func findSmudge(lines []string) int {
	sum := 0
	for _, image := range images.img {
		fmt.Fprintln(os.Stderr, "checking for horizontal reflection")
		n, ok := horizontalReflection(image, findReflectionWithDifference)
		if ok {
			sum += 100 * (n + 1)
		}
		fmt.Fprintln(os.Stderr, "checking for vertical reflection")
		n, ok = verticalReflection(image, findReflectionWithDifference)
		if ok {
			sum += (n + 1)
		}
	}

	return sum
}

func reflectionLines(lines []string) int {
	parselines(lines)

	sum := 0
	for _, image := range images.img {
		fmt.Fprintln(os.Stderr, "checking for horizontal reflection")
		n, ok := horizontalReflection(image, findReflection)
		if ok {
			sum += 100 * (n + 1)
		}
		fmt.Fprintln(os.Stderr, "checking for vertical reflection")
		n, ok = verticalReflection(image, findReflection)
		if ok {
			sum += (n + 1)
		}
	}
	return sum
}

func parselines(lines []string) {
	addImage := func(i int) int {
		image := newImg()
		for ; i < len(lines); i++ {
			if len(lines[i]) == 0 {
				break
			}
			image.raw = append(image.raw, lines[i])
		}
		images.img = append(images.img, image)
		return i
	}

	images = newImages()

	for i := 0; i < len(lines); i++ {
		next := addImage(i)
		i = next
	}
}

func readlines() []string {
	lines := []string{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}

	return lines
}

func main() {
	lines := readlines()

	noteSummary := reflectionLines(lines)
	fmt.Printf("%d\n", noteSummary)

	noteSummary = findSmudge(lines)
	fmt.Printf("%d\n", noteSummary)
}
