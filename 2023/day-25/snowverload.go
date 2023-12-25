package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Graph struct {
	edge map[string][]string
}

func NewGraph() *Graph {
	return &Graph{
		edge: make(map[string][]string),
	}
}

func (g *Graph) AddEdge(node1, node2 string, undirected bool) {
	// Add edge from node1 to node2
	g.edge[node1] = append(g.edge[node1], node2)

	if undirected {
		// Add edge from node2 to node1
		g.edge[node2] = append(g.edge[node2], node1)
	}
}

func (g *Graph) DeleteEdge(node1, node2 string, undirected bool) {
	g.edge[node1] = remove(g.edge[node1], node2)
	if undirected {
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
	edgeMap := make(map[string]bool)
	for node, neighbors := range g.edge {
		for _, neighbor := range neighbors {
			// Ensure the edge is always represented in the same way,
			// regardless of the order of the nodes
			edge := node + "-" + neighbor
			if node > neighbor {
				edge = neighbor + "-" + node
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

func main() {
	g := NewGraph()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		re := regexp.MustCompile(`[:\s]+`)
		nodes := re.Split(line, -1)
		for i := 1; i < len(nodes); i++ {
			g.AddEdge(nodes[0], nodes[i], true)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	fmt.Println(g)

	fmt.Println("Does the graph contain a cycle?", g.DetectCycle())
	fmt.Println("Number of nodes in the graph:", g.CountNodes())
	fmt.Println("Number of edges in the graph:", g.CountEdges())
	fmt.Println("Number of connected components in the graph:", g.CountComponents())
}
