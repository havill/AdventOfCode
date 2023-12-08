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

func main() {
	var instructions string
	var instructionsI int = 0
	var steps int = 0
	var network Network = make(Network)
	var youAreHere *node
	var ghosts []*node

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
			//fmt.Println("node: ", n)
			addNodeToNetwork(network, n)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(1)
	}

	youAreHere = network["AAA"]
	for youAreHere != nil && youAreHere.label != "ZZZ" {
		if instructionsI >= len(instructions) {
			instructionsI = 0
		}
		//fmt.Fprintln(os.Stdout, "You are at: ", youAreHere.label, " instructions: ", instructionsI)
		if instructions[instructionsI] == 'L' {
			left := network[youAreHere.left].label
			youAreHere = network[left]
			//fmt.Fprintln(os.Stdout, "You went left to: ", youAreHere.label)
		} else if instructions[instructionsI] == 'R' {
			right := network[youAreHere.right].label
			youAreHere = network[right]
			//fmt.Fprintln(os.Stdout, "You went right to: ", youAreHere.label)
		} else {
			//fmt.Fprintln(os.Stderr, "Invalid instruction: ", instructions[instructionsI])
			os.Exit(1)
		}
		steps += 1
		instructionsI += 1
	}
	if youAreHere == nil {
		fmt.Fprintln(os.Stderr, "You are lost!")
		os.Exit(1)
	}
	fmt.Println(steps)
}
