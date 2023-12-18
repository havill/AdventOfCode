package main

import (
	"container/heap"
	"fmt"
	"image"
	"os"
	"strings"
)

type hqi[T any] struct {
	v T
	p int
}

type HeapQ[T any] []hqi[T]

func (q HeapQ[_]) Len() int           { return len(q) }
func (q HeapQ[_]) Less(i, j int) bool { return q[i].p < q[j].p }
func (q HeapQ[_]) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }
func (q *HeapQ[T]) Push(x any)        { *q = append(*q, x.(hqi[T])) }
func (q *HeapQ[_]) Pop() (x any)      { x, *q = (*q)[len(*q)-1], (*q)[:len(*q)-1]; return x }
func (q *HeapQ[T]) GPush(v T, p int)  { heap.Push(q, hqi[T]{v, p}) }
func (q *HeapQ[T]) GPop() (T, int)    { x := heap.Pop(q).(hqi[T]); return x.v, x.p }

type State struct {
	Pos image.Point
	Dir image.Point
}

func main() {
	input, _ := os.ReadFile("input.txt")
	split := strings.Fields(string(input))

	grid, end := map[image.Point]int{}, image.Point{0, 0}
	for y, s := range split {
		for x, r := range s {
			grid[image.Point{x, y}] = int(r - '0')
			end = image.Point{x, y}
		}
	}

	recurseMinimax := func(min, max int) int {
		queue, visited := HeapQ[State]{}, map[State]struct{}{}
		queue.GPush(State{image.Point{0, 0}, image.Point{1, 0}}, 0)
		queue.GPush(State{image.Point{0, 0}, image.Point{0, 1}}, 0)

		for len(queue) > 0 {
			node, heat := queue.GPop()

			if node.Pos == end {
				return heat
			}
			if _, ok := visited[node]; ok {
				continue
			}
			visited[node] = struct{}{}

			for _, d := range []image.Point{
				{node.Dir.Y, node.Dir.X}, {-node.Dir.Y, -node.Dir.X},
			} {
				for i := min; i <= max; i++ {
					n := node.Pos.Add(d.Mul(i))
					if _, ok := grid[n]; ok {
						h := 0
						for j := 1; j <= i; j++ {
							h += grid[node.Pos.Add(d.Mul(j))]
						}
						queue.GPush(State{n, d}, heat+h)
					}
				}
			}
		}
		return -1
	}

	fmt.Println(recurseMinimax(1, 3))  // min of 1 block, max of 3 blocks forward
	fmt.Println(recurseMinimax(4, 10)) // part 2: ultra crucibles: 4 min, 10 max
}
