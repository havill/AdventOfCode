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

func parseBinaryOp(expr string) (left string, op Comparison, right string, err error) {
	if strings.Contains(expr, "<") {
		parts := strings.Split(expr, "<")
		return strings.TrimSpace(parts[0]), LessThan, strings.TrimSpace(parts[1]), nil
	} else if strings.Contains(expr, "=") {
		parts := strings.Split(expr, "=")
		return strings.TrimSpace(parts[0]), EqualTo, strings.TrimSpace(parts[1]), nil
	} else if strings.Contains(expr, ">") {
		parts := strings.Split(expr, ">")
		return strings.TrimSpace(parts[0]), GreaterThan, strings.TrimSpace(parts[1]), nil
	} else {
		return "", 0, "", fmt.Errorf("invalid expression: %s", expr)
	}
}

func parseTernaryOp(expr string) (left string, op Comparison, action string, err error) {
	parts := strings.Split(expr, ":")
	if len(parts) != 2 {
		return "", 0, "", fmt.Errorf("invalid ternary expression: %s", expr)
	}

	condition := strings.TrimSpace(parts[0])
	action = strings.TrimSpace(parts[1])

	left, op, _, err = parseBinaryOp(condition)
	if err != nil {
		return "", 0, "", fmt.Errorf("invalid condition in ternary expression: %s", condition)
	}

	return left, op, action, nil
}

func main() {
	//var w Workflows = make(Workflows, 0)
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
				id, _, val, err := parseTernaryOp(item)
				if err != nil {
					fmt.Println("Error parsing binary expression:", err)
					continue
				}
				i, err := strconv.Atoi(val)
				if err != nil {
					fmt.Println("Error converting string to integer:", err)
					continue
				}
				switch id {
				case "x":
					p = append(p, Ratings{Cool: i})
				case "m":
					p = append(p, Ratings{Musical: i})
				case "a":
					p = append(p, Ratings{Aerodynamic: i})
				case "s":
					p = append(p, Ratings{Shiny: i})
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
