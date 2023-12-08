package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type HandType int

const (
	High_card HandType = iota
	One_pair
	Two_pair
	Three_of_a_kind
	Full_house
	Four_of_a_kind
	Five_of_a_kind
)

type Card int

const (
	Zero Card = iota // this is a sentinal value and should never be used
	Joker
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

type Parsed struct {
	Hand []Card
	Bid  int
	Type HandType
}

func charToCard(c rune, jokers bool) (Card, error) {
	switch c {
	case '0':
		return Zero, nil
	case '1':
		return Joker, nil
	case '2':
		return Two, nil
	case '3':
		return Three, nil
	case '4':
		return Four, nil
	case '5':
		return Five, nil
	case '6':
		return Six, nil
	case '7':
		return Seven, nil
	case '8':
		return Eight, nil
	case '9':
		return Nine, nil
	case 'T':
		return Ten, nil
	case 'J':
		if jokers {
			return Joker, nil
		}
		return Jack, nil
	case 'Q':
		return Queen, nil
	case 'K':
		return King, nil
	case 'A':
		return Ace, nil
	default:
		return Zero, fmt.Errorf("invalid character: %c", c)
	}
}

func removePunctuationAndWhitespace(s string) string {
	f := func(r rune) rune {
		if unicode.IsPunct(r) || unicode.IsSpace(r) {
			return -1
		}
		return r
	}
	return strings.Map(f, s)
}

func IsZOfAKind(hand []Card, z int) bool {
	counts := make(map[Card]int)
	var jokers int
	for _, card := range hand {
		if card == Joker {
			jokers++
		} else {
			counts[card]++
		}
	}
	if jokers >= z {
		return true
	}
	for _, count := range counts {
		if count+jokers >= z {
			return true
		}
	}
	return false
}

func IsFullHouse(hand []Card) bool {
	counts := make(map[Card]int)
	var jokers int
	for _, card := range hand {
		if card == Joker {
			jokers++
		} else {
			counts[card]++
		}
	}
	var pair, threeOfAKind bool
	var threeOfAKindCard Card
	for card, count := range counts {
		if count+jokers >= 3 {
			threeOfAKind = true
			threeOfAKindCard = card
			jokers -= 3 - count
			break
		}
	}
	for card, count := range counts {
		if card != threeOfAKindCard && count+jokers >= 2 {
			pair = true
			break
		}
	}
	return pair && threeOfAKind
}

func IsTwoPair(hand []Card) bool {
	counts := make(map[Card]int)
	var jokers int
	for _, card := range hand {
		if card == Joker {
			jokers++
		} else {
			counts[card]++
		}
	}
	pairs := 0
	for _, count := range counts {
		if count+jokers >= 2 {
			pairs++
			if jokers > 0 {
				jokers--
			}
		}
	}
	return pairs >= 2
}

func less(a, b []Card) bool {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return len(a) < len(b)
}

func totalWinnings(hands []Parsed) int {
	total := 0
	for rank, hand := range hands {
		total += hand.Bid * rank
	}
	return total
}

func appendCardToHand(c rune, x *Parsed, jokersWild bool) error {
	card, err := charToCard(c, jokersWild)
	if err != nil {
		return fmt.Errorf("Error converting char to a card: %v", c)
	}
	x.Hand = append(x.Hand, card)
	return nil
}

func determineHandTypeAndAppend(x *Parsed, hands []Parsed) []Parsed {
	if IsZOfAKind(x.Hand, 5) {
		x.Type = Five_of_a_kind
	} else if IsZOfAKind(x.Hand, 4) {
		x.Type = Four_of_a_kind
	} else if IsFullHouse(x.Hand) {
		x.Type = Full_house
	} else if IsZOfAKind(x.Hand, 3) {
		x.Type = Three_of_a_kind
	} else if IsTwoPair(x.Hand) {
		x.Type = Two_pair
	} else if IsZOfAKind(x.Hand, 2) {
		x.Type = One_pair
	} else {
		x.Type = High_card
	}
	return append(hands, *x)
}

func main() {
	var hands []Parsed = make([]Parsed, 1)
	// prefix with sentinel value so first real hand is rank 1
	hands[0] = Parsed{Hand: []Card{Zero, Zero, Zero, Zero, Zero}, Bid: 0, Type: -1}

	var jokerHands []Parsed = make([]Parsed, 1)
	// prefix with sentinel value so first real hand is rank 1
	jokerHands[0] = Parsed{Hand: []Card{Zero, Zero, Zero, Zero, Zero}, Bid: 0, Type: -1}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		var x, y Parsed

		line := scanner.Text()

		line = strings.TrimSpace(line)
		lastSpace := strings.LastIndex(line, " ")
		left := line[:lastSpace]
		right := line[lastSpace+1:]
		left = removePunctuationAndWhitespace(left)
		hand := strings.ToUpper(left)

		i, err := strconv.Atoi(right)
		if err != nil {
			fmt.Println("Error converting string to int:", err)
		}
		x.Bid = i
		y.Bid = i
		for _, c := range hand {
			appendCardToHand(c, &x, false)
			appendCardToHand(c, &y, true)
		}
		hands = determineHandTypeAndAppend(&x, hands)
		jokerHands = determineHandTypeAndAppend(&y, jokerHands)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	sort.Slice(hands, func(i, j int) bool {
		if hands[i].Type != hands[j].Type {
			return hands[i].Type < hands[j].Type
		}
		return less(hands[i].Hand, hands[j].Hand)
	})
	sort.Slice(jokerHands, func(i, j int) bool {
		if jokerHands[i].Type != jokerHands[j].Type {
			return jokerHands[i].Type < jokerHands[j].Type
		}
		return less(jokerHands[i].Hand, jokerHands[j].Hand)
	})

	fmt.Println(totalWinnings(hands))
	fmt.Println(totalWinnings(jokerHands))
}
