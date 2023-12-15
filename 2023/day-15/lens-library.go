package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type lens struct {
	label       string
	focalLength int
}

func addOrChangeLens(label string, focalLen int, box *[]lens) {
	for i, l := range *box {
		if l.label == label {
			(*box)[i].focalLength = focalLen
			return
		}
	}
	*box = append(*box, lens{label: label, focalLength: focalLen})
}

func removeLens(label string, box *[]lens) {
	for i, l := range *box {
		if l.label == label {
			// Remove the lens from the slice
			*box = append((*box)[:i], (*box)[i+1:]...)
			return
		}
	}
}

func parseStep(step string) (string, rune, int, error) {
	var label string
	var operation rune
	var focalLength int

	for _, char := range step {
		if unicode.IsLetter(char) {
			label += string(char)
		} else if char == '=' {
			operation = char
			break
		} else if char == '-' {
			operation = char
			return label, operation, 0, nil
		} else {
			return "", 0, 0, fmt.Errorf("invalid character: %c", char)
		}
	}

	if operation == 0 {
		return "", 0, 0, fmt.Errorf("no operation found in step")
	}

	focalLengthStr := strings.TrimPrefix(step, label+string(operation))
	focalLength, err := strconv.Atoi(focalLengthStr)
	if err != nil {
		return "", 0, 0, fmt.Errorf("invalid focal length: %v", err)
	}

	return label, operation, focalLength, nil
}

func calculateHash(s string) int {
	hash := 0
	for _, c := range s {
		hash = (hash + int(c)) * 17 % 256
	}
	return hash
}

func focusingPower(boxes [256][]lens) int {
	total := 0
	for j, box := range boxes {
		for i, lens := range box {
			total += (j + 1) * (i + 1) * lens.focalLength
		}
	}
	return total
}

func main() {
	var boxes [256][]lens

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	// Remove all whitespace
	input = strings.ReplaceAll(input, " ", "")
	input = strings.ReplaceAll(input, "\n", "")

	// Split the input into steps
	steps := strings.Split(input, ",")

	// Calculate the sum of the hash values
	sum := 0
	for _, step := range steps {
		sum += calculateHash(step)
		label, op, focalLen, err := parseStep(step)
		correctBox := calculateHash(label)
		//fmt.Fprintf(os.Stderr, "step=%q label=%q op='%c' focal=%d box=%d\n", step, label, op, focalLen, correctBox)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid step: %v\n", err)
			return
		}
		switch op {
		case '-':
			removeLens(label, &boxes[correctBox])
		case '=':
			addOrChangeLens(label, focalLen, &boxes[correctBox])
		default:
			fmt.Fprintf(os.Stderr, "invalid operation: %c\n", op)
			return
		}
		/*
			for i, box := range boxes {
				if len(box) == 0 {
					continue
				}
				fmt.Printf("\tBox %d: %v\n", i, box)
			}
			fmt.Println()
		*/
	}

	fmt.Println(sum)
	fmt.Println(focusingPower(boxes))
}
