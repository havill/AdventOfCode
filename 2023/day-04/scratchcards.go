package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Card struct {
	Number    int
	Winning   []int
	Possessed []int
}

type Deck struct {
	Cards []Card
}

func (c *Card) AddToWinning(num int) {
	c.Winning = append(c.Winning, num)
}

func (c *Card) AddToPossessed(num int) {
	c.Possessed = append(c.Winning, num)
}

func (c *Card) InWinning(Possessed int) bool {
	for _, num := range c.Winning {
		if num == Possessed {
			return true
		}
	}
	return false
}

func TotalCardPointsAndMatches(c *Card) (int, int) {
	total := 0
	matches := 0

	for _, num := range c.Possessed {
		if c.InWinning(num) {
			if total == 0 {
				total = 1
			} else {
				total *= 2
			}
			matches++
		}
	}
	return total, matches
}

func splitByColon(s string) (string, string, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return "", "", errors.New("input string does not contain a colon")
	}
	return parts[0], parts[1], nil
}

func splitByVerticalBar(s string) (string, string, error) {
	parts := strings.SplitN(s, "|", 2)
	if len(parts) != 2 {
		return "", "", errors.New("input string does not contain a vertical bar")
	}
	return parts[0], parts[1], nil
}

func extractInt(s string) (int, error) {
	re := regexp.MustCompile(`\d+`)
	match := re.FindStringSubmatch(s)
	if match == nil {
		return 0, errors.New("no integer found in string")
	}
	return strconv.Atoi(match[0])
}

func (d *Deck) AddCard(c Card) {
	d.Cards = append(d.Cards, c)
}

func (d *Deck) CopyCardsToEndofDeck(Number int, copies int) {
	start := 0
	for start < len(d.Cards) {
		if d.Cards[start].Number == Number {
			break
		}
		start++
	}
	for i := start + 1; copies > 0; i++ {
		d.Cards = append(d.Cards, d.Cards[i])
		copies--
		// fmt.Println("Copying card ", d.Cards[i].Number, " to end of deck")
	}
}

func ExtractIntsFromString(s string) []int {
	re := regexp.MustCompile(`\d+`)
	match := re.FindAllString(s, -1)
	if match == nil {
		return nil
	}
	ints := make([]int, len(match))
	for i, num := range match {
		ints[i], _ = strconv.Atoi(num)
	}
	return ints
}

func main() {
	var d Deck
	var c Card
	total := 0

	d.Cards = make([]Card, 0)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		left, right, _ := splitByColon(line)
		winning, possessed, _ := splitByVerticalBar(right)

		c.Number, _ = extractInt(left)
		c.Winning = ExtractIntsFromString(winning)
		c.Possessed = ExtractIntsFromString(possessed)
		d.AddCard(c)
	}
	for _, card := range d.Cards {
		points, _ := TotalCardPointsAndMatches(&card)
		total += points

	}
	fmt.Println(total)
	i := 0
	for i < len(d.Cards) {
		card := d.Cards[i]
		_, matches := TotalCardPointsAndMatches(&card)
		// fmt.Println("Card ", card.Number, " has ", matches, " matches")
		d.CopyCardsToEndofDeck(card.Number, matches)
		i++
		// fmt.Println("we now have ", len(d.Cards), " cards")
	}
	fmt.Println(len(d.Cards))
}
