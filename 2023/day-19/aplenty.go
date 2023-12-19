package main

import (
	"bufio"
	"fmt"
	"os"
)

type Categories rune

const (
	Cool        Categories = 'x'
	Musical     Categories = 'm'
	Aerodynamic Categories = 'a'
	Shiny       Categories = 's'
)

type Ratings map[Categories]int

type Result rune

const (
	Reject Result = 'R'
	Accept Result = 'A'
)

type Comparison rune

const (
	LessThan    Comparison = '<'
	EqualTo     Comparison = '='
	GreaterThan Comparison = '>'
)

type Condition struct {
	Category string
	Operator Comparison
	Value    int
}

type Rule struct {
	If           Condition
	Then         Result
	WorkflowName string
}

type Workflow struct {
	Name  string
	Rules []Rule
}

type Workflows []Workflow
type Parts []Ratings

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
