package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	combinations "github.com/mxschmitt/golang-combinations"
)

type Graph struct {
	edge map[string][]string
}

type EdgePair struct {
	a string
	b string
}

func NewGraph() *Graph {
	return &Graph{
		edge: make(map[string][]string),
	}
}

func CloneGraph(g *Graph) *Graph {
	newGraph := NewGraph()
	for node, edges := range g.edge {
		newEdges := make([]string, len(edges))
		copy(newEdges, edges)
		newGraph.edge[node] = newEdges
	}
	return newGraph
}

func (g *Graph) AddEdge(node1, node2 string, directed bool) {
	// Add edge from node1 to node2
	g.edge[node1] = append(g.edge[node1], node2)

	if !directed {
		// Add edge from node2 to node1
		g.edge[node2] = append(g.edge[node2], node1)
	}
}

func (g *Graph) DeleteEdge(node1, node2 string, directed bool) {
	g.edge[node1] = remove(g.edge[node1], node2)
	if !directed {
		g.edge[node2] = remove(g.edge[node2], node1)
	}
}

func remove(slice []string, val string) []string {
	for i, item := range slice {
		if item == val {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func (g *Graph) SimpleDFS(node string, visited map[string]bool) {
	if _, ok := visited[node]; ok {
		return
	}
	visited[node] = true
	for _, neighbor := range g.edge[node] {
		g.SimpleDFS(neighbor, visited)
	}
}

func (g *Graph) DFS(v string, visited map[string]bool, parent string) bool {
	visited[v] = true

	for _, i := range g.edge[v] {
		if !visited[i] {
			if g.DFS(i, visited, v) {
				return true
			}
		} else if i != parent {
			return true
		}
	}

	return false
}

func (g *Graph) DetectCycle() bool {
	visited := make(map[string]bool)

	for v := range g.edge {
		if !visited[v] {
			if g.DFS(v, visited, "") {
				return true
			}
		}
	}

	return false
}

func (g *Graph) CountEdges() int {
	edgeMap := make(map[EdgePair]bool)
	for node, neighbors := range g.edge {
		for _, neighbor := range neighbors {
			// Ensure the edge is always represented in the same way,
			// regardless of the order of the nodes
			edge := EdgePair{node, neighbor}
			if node > neighbor {
				edge = EdgePair{neighbor, node}
			}
			edgeMap[edge] = true
		}
	}
	return len(edgeMap)
}

func (g *Graph) CountNodes() int {
	visited := make(map[string]bool)
	for node := range g.edge {
		g.SimpleDFS(node, visited)
	}
	return len(visited)
}

func CountNodesInComponents(g *Graph) []int {
	visited := make(map[string]bool)
	counts := []int{}

	for node := range g.edge {
		if visited[node] {
			continue
		}

		count := CountDFS(g, node, visited)
		counts = append(counts, count)
	}

	return counts
}

func CountDFS(g *Graph, node string, visited map[string]bool) int {
	visited[node] = true
	count := 1

	for _, neighbor := range g.edge[node] {
		if !visited[neighbor] {
			count += CountDFS(g, neighbor, visited)
		}
	}

	return count
}

func (g *Graph) CountComponents() int {
	visited := make(map[string]bool)
	count := 0
	for node := range g.edge {
		if _, ok := visited[node]; !ok {
			g.SimpleDFS(node, visited)
			count++
		}
	}
	return count
}

func (g *Graph) EdgeExists(node1, node2 string) bool {
	for _, neighbor := range g.edge[node1] {
		if neighbor == node2 {
			return true
		}
	}
	return false
}

func (g *Graph) IterateEdges(directed bool) []EdgePair {
	var list []EdgePair

	for node, edges := range g.edge {
		for _, edge := range edges {
			if directed || node < edge {
				list = append(list, EdgePair{a: node, b: edge})
			}
		}
	}

	return list
}

func (g *Graph) String() string {
	var result strings.Builder
	for node, neighbors := range g.edge {
		result.WriteString(node + ": ")
		for _, neighbor := range neighbors {
			result.WriteString(neighbor + " ")
		}
		result.WriteString("\n")
	}
	return result.String()
}

func arrayProduct(nums []int) int {
	product := 1
	for _, num := range nums {
		product *= num
	}
	return product
}

// note that the Advent of Code problem uses the term "components" to refer to
// groups of snow producing components (iow, nodes / vertexes in CompSci terms)
// connected by "wires" (iow, edges in CompSci terms), not the traditional usage
func FindWiresToCut(g *Graph, groups, toDisconnect int) *Graph {
	// Find all edges in the graph
	wires := g.IterateEdges(false)
	combos := combinations.Combinations(wires, toDisconnect)

	// Iterate over all edges and remove them from the graph
	// If the graph contains a cycle after removing the edge,
	// then it is a critical edge
	for _, set := range combos {
		// Make a copy of the graph
		newGraph := CloneGraph(g)

		for _, edge := range set {
			newGraph.DeleteEdge(edge.a, edge.b, false)
		}

		if newGraph.CountComponents() == groups {
			return newGraph
		}
	}

	return nil
}

var debug bool

func main() {
	flag.BoolVar(&debug, "debug", false, "enable debug mode")
	flag.Parse()

	if debug {
		fmt.Println("Debug mode enabled")
	}

	g := NewGraph()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		re := regexp.MustCompile(`[:\s]+`)
		nodes := re.Split(line, -1)
		for i := 1; i < len(nodes); i++ {
			g.AddEdge(nodes[0], nodes[i], false)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	if debug {
		fmt.Println(g)
		fmt.Println(g.IterateEdges(false))
	}

	fmt.Println("BEFORE DELETION")
	fmt.Println("Does the graph contain cycle(s)?", g.DetectCycle())
	fmt.Println("Number of nodes in the graph   :", g.CountNodes())
	fmt.Println("Number of edges in the graph   :", g.CountEdges())
	fmt.Println("Number of components in graph  :", g.CountComponents())
	fmt.Println("Number of nodes in components  :", CountNodesInComponents(g))

	solution := FindWiresToCut(g, 2, 3)

	// https://adventofcode.com/2023/day/25
	// solution = g.CloneGraph(solution)
	// solution.DeleteEdge("hfx", "pzl", false)
	// solution.DeleteEdge("bvb", "cmg", false)
	// solution.DeleteEdge("nvd", "jqt", false)

	fmt.Println("AFTER DELETION")
	fmt.Println("Does the graph contain cycle(s)?", solution.DetectCycle())
	fmt.Println("Number of nodes in the graph   :", solution.CountNodes())
	fmt.Println("Number of edges in the graph   :", solution.CountEdges())
	fmt.Println("Number of components in graph  :", solution.CountComponents())
	fmt.Println("Number of nodes in components  :", CountNodesInComponents(solution))

	fmt.Println("Part 1: ", arrayProduct(CountNodesInComponents(solution)))
}
