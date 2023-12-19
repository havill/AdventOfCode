package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var debug = flag.Bool("debug", false, "enable debug mode")

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

func parseLine(line string) (left string, right string) {
	line = strings.TrimSpace(line)
	index := strings.Index(line, "{")
	if index == -1 {
		return line, ""
	}

	left = strings.TrimSpace(line[:index])
	right = strings.TrimSpace(line[index+1 : len(line)-1])

	return left, right
}

func parseList(list string) []string {
	items := strings.Split(list, ",")
	for i, item := range items {
		items[i] = strings.TrimSpace(item)
	}
	return items
}

func parseTernaryOp(expr string) (subj string, operator Comparison, value int, obj string, err error) {
	parts := strings.Fields(expr)
	if len(parts) != 3 {
		return "", 0, 0, "", fmt.Errorf("invalid expression format")
	}

	subj = parts[0]

	switch parts[1] {
	case "<":
		operator = LessThan
	case "=":
		operator = EqualTo
	case ">":
		operator = GreaterThan
	default:
		return "", 0, 0, "", fmt.Errorf("invalid operator")
	}

	valueStrs := strings.Split(parts[2], ":")
	if len(valueStrs) == 2 {
		value, err = strconv.Atoi(valueStrs[0])
		if err != nil {
			return "", 0, 0, "", fmt.Errorf("invalid value: %v", err)
		}
		obj = valueStrs[1]
	} else {
		value, err = strconv.Atoi(parts[2])
		if err != nil {
			return "", 0, 0, "", fmt.Errorf("invalid value: %v", err)
		}
	}
	return subj, operator, value, obj, nil
}

func main() {
	var w Workflows = make(Workflows, 0)
	var p Parts = make(Parts, 0)

	flag.Parse()

	if *debug {
		fmt.Println("Debug mode enabled")
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		left, right := parseLine(line)
		fmt.Println("Left:", left, "Right:", right)
		if len(left) > 0 {
		} else { // no identifier to left of brace; we have a part
			list := parseList(right)
			for _, item := range list {
				fmt.Println("Item:", item)
				id, _, val, _, err := parseTernaryOp(item)
				if err != nil {
					fmt.Println("Error parsing binary expression:", err)
					continue
				}
				switch id {
				case "x":
					p = append(p, Ratings{Cool: val})
				case "m":
					p = append(p, Ratings{Musical: val})
				case "a":
					p = append(p, Ratings{Aerodynamic: val})
				case "s":
					p = append(p, Ratings{Shiny: val})
				default:
					fmt.Println("Unknown identifier:", id)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
