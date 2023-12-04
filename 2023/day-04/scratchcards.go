package main

type Card struct {
	Winning []int
	Possess []int
}

func (c *Card) AddToWinning(num int) {
	c.Winning = append(c.Winning, num)
}

func (c *Card) AddToPossess(num int) {
	c.Possess = append(c.Winning, num)
}

func (c *Card) InWinning(Possessed int) bool {
	for i, num := range c.Winning {
		if num == Possessed {
			return true
		}
	}
	return false
}
