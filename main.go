package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	// "strconv"
	// "errors"
	"sync"
	"time"
)

// Any live cell with fewer than two live neighbours dies, as if by underpopulation.
// Any live cell with two or three live neighbours lives on to the next generation.
// Any live cell with more than three live neighbours dies, as if by overpopulation.
// Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.

type Pos struct {
	x int
	y int
}

type Cell struct {
	Alive          bool
	Pos            Pos
	Channel        chan int
	AliveNeighbors int
}

type Event struct {
	EventType EventType
	Pos       Pos
	Val       int
}

type World struct {
	Rows int
	Cols int
	Map  [][]*Cell
}

type EventType byte

const (
	CellComeAlive EventType = 0
	CellDied      EventType = 1
)

func NewWorld(rows int, cols int) *World {
	Map := make([][]*Cell, rows)

	for y := 0; y < rows; y++ {
		Map[y] = make([]*Cell, cols)
		for x := 0; x < cols; x++ {
			Map[y][x] = NewCell(x, y)
		}
	}

	w := World{
		Rows: rows,
		Cols: cols,
		Map:  Map,
	}

	return &w
}

func (e *Event) ExecuteEvent(world *World, affectedNeighbours chan *Cell) error {
	neighbours := getNeighbours(e.Pos, world)

	// switch e.EventType {
	// case CellComeAlive:
	//
	// case CellDied:
	//
	// default:
	// 	return errors.New("Wrong event type.")
	// }

	for _, neighbour := range neighbours {
		affectedNeighbours <- neighbour
	}

	return nil
}

func (c *Cell) getPos() Pos {
	return c.Pos
}

func (c *Cell) getChan() chan int {
	return c.Channel
}

func (c *Cell) updateCell(ch chan *Event) {
	if c.Alive == true {
		if c.AliveNeighbors < 2 || c.AliveNeighbors > 3 {
			c.Alive = false
			newEvent := &Event{
				Pos:       c.Pos,
				EventType: CellDied,
				Val:       -1,
			}

			// fmt.Printf("Cell dies at x: %v y: %v\n", c.Pos.x, c.Pos.y)

			ch <- newEvent
		}
	} else {
		if c.AliveNeighbors == 3 {
			c.Alive = true

			newEvent := &Event{
				Pos:       c.Pos,
				EventType: CellComeAlive,
				Val:       1,
			}

			// fmt.Printf("Cell comes alive at x: %v y: %v\n", c.Pos.x, c.Pos.y)

			ch <- newEvent
		}
	}
}

func NewCell(x int, y int) *Cell {
	c := Cell{
		Alive: false,
		Pos: Pos{
			x: x,
			y: y,
		},
		Channel: make(chan int, 1),
	}

	return &c
}

func GetPos() *Pos {
	p := Pos{
		x: 4,
		y: 10,
	}

	return &p
}

func getNeighbours(pos Pos, world *World) []*Cell {
	neighbours := []*Cell{}
	minx := max(0, pos.x-1)
	maxx := min(world.Cols, pos.x+1)

	miny := max(0, pos.y-1)
	maxy := min(world.Rows, pos.y+1)

	for y := miny; y <= maxy; y++ {
		for x := minx; x <= maxx; x++ {
			if x == pos.x && y == pos.y {
				continue
			}
			// fmt.Printf("neighbour identified: x: %v, y: %v\n", x, y)

			neighbours = append(neighbours, world.Map[y][x])
		}
	}

	return neighbours
}

// func getNeighbours(neighbours chan *Cell, pos Pos, world [][]*Cell) {
// 	minx := max(0, pos.x-1)
// 	maxx := min(99, pos.x+1)
//
// 	miny := max(0, pos.y-1)
// 	maxy := min(39, pos.y+1)
//
// 	for y := miny; y <= maxy; y++ {
// 		for x := minx; x <= maxx; x++ {
// 			neighbours <- world[y][x]
// 		}
// 	}
// }

func updateNeighbors(val int, neighbours []*Cell) {
	for _, cell := range neighbours {
		cell.AliveNeighbors += val
	}
}

func (c *Cell) updateAlives(val int) {
	c.AliveNeighbors += val
}

func isObjInSlice[K comparable](object *K, slice []*K) bool {
	for _, objectInSlice := range slice {
		if object == objectInSlice {
			return true
		}
	}

	return false
}

func isPosInEventSlice(object *Pos, slice []*Event, valueCheck func(*Pos, *Pos) bool) bool {
	for _, objectInSlice := range slice {
		if valueCheck(object, &objectInSlice.Pos) {
			return true
		}
	}

	return false
}

func comparePos(pos1 *Pos, pos2 *Pos) bool {
	return pos1.x == pos2.x && pos1.y == pos2.y
}

type world struct {
	width  uint
	height uint
	alives []*Cell
	world  [][]*Cell
}

func main() {
	world := NewWorld(5, 15)
	initialSpawns := 10
	events := []*Event{}

	for initialSpawns > len(events) {
		for {
			randx := rand.Intn(world.Cols - 1)
			randy := rand.Intn(world.Rows - 1)

			pos := Pos{
				x: randx,
				y: randy,
			}

			if !isPosInEventSlice(&pos, events, comparePos) {
				event := Event{
					Pos:       pos,
					EventType: CellComeAlive,
					Val:       1,
				}

				events = append(events, &event)
				world.Map[pos.y][pos.x].Alive = true
				break
			}
		}
	}

	for {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()

		printWorld(&world.Map)

		for _, event := range events {
			fmt.Println(event)
		}

		eventChan := make(chan *Event, 3)
		var wg sync.WaitGroup
		var affectedNeighbours []*Cell

		for _, event := range events {
			wg.Add(1)
			go func(event *Event) {
				defer wg.Done()
				affectedNeighbours = getNeighbours(event.Pos, world)

				for _, cell := range affectedNeighbours {
					cell.AliveNeighbors += event.Val
					// fmt.Printf("Updated cell: %v, %v with val: %v  - now val is: %v \n", cell.Pos.x, cell.Pos.y, event.Val, cell.AliveNeighbors)
				}
			}(event)
		}

		wg.Wait()

		affectedCellsTotal := affectedNeighbours
		for _, event := range events {
			affectedCellsTotal = append(affectedCellsTotal, world.Map[event.Pos.y][event.Pos.x])
		}

		go func() {
			wg.Wait()
			close(eventChan)
		}()

		for _, cell := range affectedCellsTotal {
			wg.Add(1)
			go func(cell *Cell) {
				defer wg.Done()
				cell.updateCell(eventChan)
			}(cell)
		}

		printWorldWAlives(&world.Map)

		time.Sleep(1 * time.Second)

		events = []*Event{}
		counter := 0
		for data := range eventChan {
			counter++
			events = append(events, data)
		}

		for _, event := range events {
			fmt.Println(event)
		}

		fmt.Printf("counter: %v\n", counter)
		fmt.Scanln()
	}
}

func printWorld(w *[][]*Cell) {
	for i, row := range *w {
		if i < 10 {
			fmt.Printf("%v  ", i)
		} else {
			fmt.Printf("%v ", i)
		}
		for _, cell := range row {
			if !cell.Alive {
				fmt.Print(".")
			} else {
				fmt.Print("#")
			}
		}
		fmt.Print("\n")
	}
}

func printWorldWAlives(w *[][]*Cell) {
	for i, row := range *w {
		if i < 10 {
			fmt.Printf("%v  ", i)
		} else {
			fmt.Printf("%v ", i)
		}
		for _, cell := range row {
			fmt.Printf("%v", cell.AliveNeighbors)
		}
		fmt.Print("\n")
	}
}

func max(a int, b int) int {
	if a < b {
		return b
	}

	return a
}

func min(a int, b int) int {
	if a < b {
		return a
	}

	return b
}

//
// neighboursChan := make(chan *Cell, 3)
// 	neighbours := []*Cell{}
//
// 	for _, cell := range relevantCells {
// 		wg.Add(1)
// 		go func(cell *Cell) {
// 			defer wg.Done()
// 			getNeighbours(neighboursChan, cell.Pos, world)
// 		}(cell)
// 	}
//
// 	go func() {
// 		wg.Wait()
// 		close(neighboursChan)
// 	}()
//
// 	for data := range neighboursChan {
// 		if !isCellInSlice(data, neighbours) {
// 			neighbours = append(neighbours, data)
// 		}
// 	}
//
// 	fmt.Println(len(neighbours))
//
// 	for _, cell := range neighbours {
// 		wg.Add(1)
// 		go func(cell *Cell) {
// 			defer wg.Done()
// 			cell.AliveNeighbors++
// 		}(cell)
//

// func isCellInSlice(cell *Cell, slice []*Cell) bool {
// 	for _, cellInSlice := range slice {
// 		if cell.Pos == cellInSlice.Pos {
// 			return true
// 		}
// 	}
// 	return false
// }
//}
