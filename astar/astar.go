package main

import (
	"fmt"
	"sort"
)

const (
	L = iota
	H
	START
	END
)

type vector struct {
	x, y, h, g, id, object, parent, g_hint, f int
}

type sortByF []*vector

func (a sortByF) Len() int           { return len(a) }
func (a sortByF) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByF) Less(i, j int) bool { return a[i].f < a[j].f }

/* playfield in pure integers */
var playfield = [][]int{
	{L, START, H, H, L, L},
	{L, L, L, H, H, L},
	{L, H, L, H, H, L},
	{L, L, H, L, L, L},
	{L, H, H, L, L, L},
	{L, L, END, L, L, L},
}

var playfieldVectoren [][]vector

func main() {
	var openlist []*vector
	var on_openlist = make(map[int]bool)
	var closedlist = make(map[int]bool)
	var success bool
	var endNode vector

	/* Convert our ascii playground to a type "vector"-structure-based playground, assign ids and
	h-values for each vector */
	if err := generateField(playfield); err != nil {
		fmt.Printf("Cant create field: %s\n", err)
	}

	/* test for end - node */
	_, _, err := findSpecialField(playfield, END)
	if err != nil {
		fmt.Println("Cant find end-position")
		return
	}

	/* get start node */
	x, y, err := findSpecialField(playfield, START)
	if err != nil {
		fmt.Println("Cant find start-position")
		return
	}

	node := playfieldVectoren[x][y]
	openlist = append(openlist, &node)

	for node_i := 0; node_i < len(openlist); node_i++ {
		node = *openlist[node_i]
		closedlist[node.id] = true

		fmt.Println("Checking node:", node.id, "parent: ", node.parent)
		printplayfield(playfieldVectoren, closedlist, on_openlist)

		if node.object == END {
			endNode = node
			success = true
			break
		}

		neighbors := node.getNeighbors()

		for i, n := range neighbors {
			if closedlist[n.id] {
				continue
			}
			if n.object == L || n.object == END {
				/* Already on openlist? Dont add again */
				if !on_openlist[n.id] {
					openlist = append(openlist, neighbors[i])
				}
				/* Update neighbors g and parent only if path lucks better  */
				if neighbors[i].g == -1 || (node.g+neighbors[i].g_hint < neighbors[i].g) {
					if neighbors[i].parent != 0 {
						fmt.Println("Re parent", node.id, "for", n.id, "old", neighbors[i].parent)
					}
					neighbors[i].setParent(node.id)
					neighbors[i].updateG(node.g)
					neighbors[i].calculateF()
					on_openlist[n.id] = true
				}
			}
		}

		/* remove top element and rewind node_i */
		openlist = openlist[1:]
		node_i--

		/* resort by value F */
		sort.Sort(sortByF(openlist))

	}

	if success {
		fmt.Println("")
		fmt.Println("solution in reversed steps from ending node")
		for node := endNode; ; {
			fmt.Println(node.id)
			node, err = findVectorByID(playfieldVectoren, node.parent)
			if err != nil {
				break
			}
		}
	} else {
		fmt.Println("Did not find the ending ;/")
	}

}

func printplayfield(playfieldVectoren [][]vector, closedlist map[int]bool, on_openlist map[int]bool) {
	for i, v := range playfieldVectoren {
		for j := range v {
			if playfieldVectoren[i][j].object == H {
				fmt.Printf("H (%02d: g:%02d,p:%02d,f:%02d,h:%02d) ", playfieldVectoren[i][j].id, playfieldVectoren[i][j].g, playfieldVectoren[i][j].parent, playfieldVectoren[i][j].f, playfieldVectoren[i][j].h)
				continue
			}
			if closedlist[playfieldVectoren[i][j].id] {
				fmt.Printf("x (%02d: g:%02d,p:%02d,f:%02d,h:%02d) ", playfieldVectoren[i][j].id, playfieldVectoren[i][j].g, playfieldVectoren[i][j].parent, playfieldVectoren[i][j].f, playfieldVectoren[i][j].h)
			} else if on_openlist[playfieldVectoren[i][j].id] {
				fmt.Printf("o (%02d: g:%02d,p:%02d,f:%02d,h:%02d) ", playfieldVectoren[i][j].id, playfieldVectoren[i][j].g, playfieldVectoren[i][j].parent, playfieldVectoren[i][j].f, playfieldVectoren[i][j].h)
			} else {
				fmt.Printf("- (%02d: g:%02d,p:%02d,f:%02d,h:%02d) ", playfieldVectoren[i][j].id, playfieldVectoren[i][j].g, playfieldVectoren[i][j].parent, playfieldVectoren[i][j].f, playfieldVectoren[i][j].h)
			}
		}
		fmt.Println("")
	}
}

func generateField(playfield [][]int) error {
	var id = 1

	x, y, err := findSpecialField(playfield, END)
	if err != nil {
		return fmt.Errorf("Cant find position: %d\n", END)
	}

	target := &vector{x: x, y: y}

	playfieldVectoren = make([][]vector, 1)

	for i, v := range playfield {

		if i == len(playfieldVectoren) {
			playfieldVectoren = append(playfieldVectoren, []vector{})
		}

		playfieldVectoren[i] = make([]vector, 1)

		for j := range v {
			if j == len(playfieldVectoren[i]) {
				playfieldVectoren[i] = append(playfieldVectoren[i], vector{})
			}
			v := vector{x: i, y: j, object: playfield[i][j], id: id}
			v.generateH(*target)
			if v.object == L || v.object == END {
				v.updateG(-1)
			}
			playfieldVectoren[i][j] = v
			id++
		}
	}
	return nil
}

func (v *vector) setParent(id int) {
	v.parent = id
}

func (v *vector) updateG(parent int) {
	v.g = v.g_hint + parent
}

func (v *vector) calculateF() {
	v.f = v.g + v.h
}

func (v *vector) getNeighbors() []*vector {
	var neighbors = make([]*vector, 0)

	/* push neighbors up  */
	if v.x-1 >= 0 {
		playfieldVectoren[v.x-1][v.y].g_hint = 10
		neighbors = append(neighbors, &playfieldVectoren[v.x-1][v.y])
	}
	/* push neighbors left */
	if v.y-1 >= 0 && playfieldVectoren[v.x][v.y-1].object == L {
		playfieldVectoren[v.x][v.y-1].g_hint = 10
		neighbors = append(neighbors, &playfieldVectoren[v.x][v.y-1])
	}

	/* check neighbor left/up */
	if v.x-1 >= 0 && v.y-1 >= 0 {
		playfieldVectoren[v.x-1][v.y-1].g_hint = 14
		neighbors = append(neighbors, &playfieldVectoren[v.x-1][v.y-1])
	}

	/* check neighbor right/up */
	if v.x+1 < len(playfieldVectoren) && v.y-1 >= 0 {
		playfieldVectoren[v.x+1][v.y-1].g_hint = 14
		neighbors = append(neighbors, &playfieldVectoren[v.x+1][v.y-1])
	}

	/* check neighbor left/down */
	if v.x-1 >= 0 && v.y+1 < len(playfieldVectoren[v.x]) {
		playfieldVectoren[v.x-1][v.y+1].g_hint = 14
		neighbors = append(neighbors, &playfieldVectoren[v.x-1][v.y+1])
	}

	/* check neighbor right/down */
	if v.x+1 < len(playfieldVectoren) && v.y+1 < len(playfieldVectoren[v.x]) {
		playfieldVectoren[v.x+1][v.y+1].g_hint = 14
		neighbors = append(neighbors, &playfieldVectoren[v.x+1][v.y+1])
	}

	/* push neighbors right */
	if v.x+1 < len(playfieldVectoren) {
		playfieldVectoren[v.x+1][v.y].g_hint = 10
		neighbors = append(neighbors, &playfieldVectoren[v.x+1][v.y])
	}
	/* push neighbors down */
	if v.y+1 < len(playfieldVectoren[v.x]) {
		playfieldVectoren[v.x][v.y+1].g_hint = 10
		neighbors = append(neighbors, &playfieldVectoren[v.x][v.y+1])
	}
	return neighbors
}

func (v *vector) generateH(target vector) {
	if v.object != END {
		v.h = HeuManhattanDistance(target, *v)
	}
}

func HeuManhattanDistance(a vector, b vector) (h int) {
	return abs(a.x-b.x)*10 + abs(a.y-b.y)*10
}

func findVectorByID(playfieldVectoren [][]vector, id int) (vector, error) {
	for i, v := range playfieldVectoren {
		for j := range v {
			if playfieldVectoren[i][j].id == id {
				return playfieldVectoren[i][j], nil
			}
		}
	}
	return vector{}, fmt.Errorf("vector %d not found", id)
}

func findSpecialField(playfield [][]int, was int) (int, int, error) {
	for i, v := range playfield {
		for j := range v {
			if playfield[i][j] == was {
				return i, j, nil
			}
		}
	}
	return 0, 0, fmt.Errorf("%d not found", was)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
