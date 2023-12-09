package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type node struct {
	label string
	left  string
	right string
}

type Network map[string]*node

func addNodeToNetwork(network Network, n *node) {
	network[n.label] = n
}

func parseInput(line string) *node {
	// Remove all punctuation
	re := regexp.MustCompile(`[^\w\s]`)
	line = re.ReplaceAllString(line, "")

	// Split the line into words
	words := strings.Fields(line)

	// Convert words to upper case
	for i, word := range words {
		words[i] = strings.ToUpper(word)
	}

	// Create a new node with the words as fields
	n := &node{
		label: words[0],
		left:  words[1],
		right: words[2],
	}

	return n
}

func isStartingNode(label string) bool {
	return strings.HasSuffix(label, "A")
}

func isEndingNode(label string) bool {
	return strings.HasSuffix(label, "Z")
}

func ghostsAreHome(ghostIsHere []*node) bool {
	for _, ghost := range ghostIsHere {
		if !isEndingNode(ghost.label) {
			return false
		}
	}
	return true
}

func debugGhosts(ghostIsHere []*node, count int) bool {
	arrived := 0

	for i, ghost := range ghostIsHere {
		fmt.Fprint(os.Stdout, "Ghost #", i, ":", ghost.label)
		if isEndingNode(ghost.label) {
			fmt.Fprint(os.Stdout, " (home)")
			arrived += 1
		}
		fmt.Fprintln(os.Stdout, "")
	}
	fmt.Fprintln(os.Stdout, "---------------------")
	return arrived > count
}

func main() {
	var instructions string
	var instructionsI int = 0
	var steps int = 0
	var network Network = make(Network)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		caps := strings.ToUpper(trimmed)
		if len(caps) == 0 {
			continue
		}
		if len(instructions) == 0 {
			instructions = caps
		} else {
			n := parseInput(line)
			addNodeToNetwork(network, n)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(1)
	}
	var youAreHere *node

	youAreHere = network["AAA"]
	for youAreHere != nil && youAreHere.label != "ZZZ" {
		if instructionsI >= len(instructions) {
			instructionsI = 0
		}
		if instructions[instructionsI] == 'L' {
			left := network[youAreHere.left].label
			youAreHere = network[left]
		} else if instructions[instructionsI] == 'R' {
			right := network[youAreHere.right].label
			youAreHere = network[right]
		} else {
			fmt.Fprintln(os.Stderr, "Invalid instruction: ", instructions[instructionsI])
			os.Exit(1)
		}
		steps += 1
		instructionsI += 1
	}
	fmt.Println(steps)
	steps = 0

	var ghostIsHere []*node

	for key, value := range network {
		if isStartingNode(key) {
			ghostIsHere = append(ghostIsHere, value)
		}
	}
	instructionsI = 0
	for ghostIsHere != nil && !ghostsAreHome(ghostIsHere) {
		/*
			if debugGhosts(ghostIsHere, math.MaxInt32) {
				os.Exit(1)
			}
		*/
		if instructions[instructionsI] == 'L' {
			//fmt.Fprint(os.Stdout, "GO LEFT (#", instructionsI, ") ")
			for i, ghost := range ghostIsHere {
				left := network[ghost.left].label
				ghostIsHere[i] = network[left]
			}
		} else if instructions[instructionsI] == 'R' {
			//fmt.Fprint(os.Stdout, "GO RIGHT (#", instructionsI, ") ")
			for i, ghost := range ghostIsHere {
				right := network[ghost.right].label
				ghostIsHere[i] = network[right]
			}
		} else {
			fmt.Fprintln(os.Stderr, "Invalid instruction: ", instructions[instructionsI])
			os.Exit(1)
		}
		steps += 1
		//fmt.Fprint(os.Stdout, "[STEP ", steps, "]\n\n")
		instructionsI += 1
		if instructionsI >= len(instructions) {
			instructionsI = 0
		}
	}
	if ghostIsHere == nil {
		fmt.Fprintln(os.Stderr, "Ghosts are lost!")
		os.Exit(1)
	}
	fmt.Println(steps)

}
