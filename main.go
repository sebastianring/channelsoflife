package main

import (
	"fmt"
	"math/rand"
	// "os"
	// "os/exec"
	"strconv"
	"sync"
	// "time"
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

func (c *Cell) getPos() Pos {
	return c.Pos
}

func (c *Cell) getChan() chan int {
	return c.Channel
}

func (c *Cell) monitorChan() {
	for {
		select {
		case change := <-c.Channel:
			c.AliveNeighbors += change
			c.updateCell()
		}
	}
}

func (c *Cell) updateCell() {
	if c.Alive == true {
		if c.AliveNeighbors < 2 || c.AliveNeighbors > 3 {
			c.Alive = false
		}
	} else {
		if c.AliveNeighbors == 3 {
			c.Alive = true
		}
	}
}

func NewCell(alive bool, x int, y int, alives *[]*Cell) *Cell {
	c := Cell{
		Alive: alive,
		Pos: Pos{
			x: x,
			y: y,
		},
		Channel: make(chan int, 1),
	}

	if alive {
		*alives = append(*alives, &c)
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

func getNeighbours(neighbours chan *Cell, pos Pos, world [][]*Cell) {
	// neighbours := make([]*Cell, 3)

	minx := max(0, pos.x-1)
	maxx := min(99, pos.x+1)

	miny := max(0, pos.y-1)
	maxy := min(39, pos.y+1)

	for y := miny; y <= maxy; y++ {
		for x := minx; x <= maxx; x++ {
			neighbours <- world[y][x]
		}
	}
}

func updateNeighbors(val int, neighbours []*Cell) {
	for _, cell := range neighbours {
		cell.AliveNeighbors += val
	}
}

func isCellInSlice(cell *Cell, slice []*Cell) bool {
	for _, cellInSlice := range slice {
		if cell.Pos == cellInSlice.Pos {
			return true
		}
	}

	return false
}

type world struct {
	width  uint
	height uint
	alives []*Cell
	world  [][]*Cell
}

func main() {
	width := 100
	height := 40
	alives := []*Cell{}

	world := make([][]*Cell, height)

	// updateChannel := make(chan int, 1)

	for y := 0; y < height; y++ {
		world[y] = make([]*Cell, width)

		for x := 0; x < width; x++ {
			var alive bool
			chance := rand.Intn(500)

			if chance == 1 {
				alive = true
			} else {
				alive = false
			}

			world[y][x] = NewCell(alive, x, y, &alives)
		}
	}

	printWorld(&world)
	fmt.Println("number of alives: " + strconv.Itoa(len(alives)))

	var wg sync.WaitGroup
	neighboursChan := make(chan *Cell, 3)
	neighbours := []*Cell{}

	for _, cell := range alives {
		wg.Add(1)
		go func(cell *Cell) {
			defer wg.Done()
			getNeighbours(neighboursChan, cell.Pos, world)
		}(cell)
	}

	go func() {
		wg.Wait()
		close(neighboursChan)
	}()

	for data := range neighboursChan {
		if !isCellInSlice(data, neighbours) {
			neighbours = append(neighbours, data)
		}
	}

	fmt.Println(len(neighbours))

	// for _, cell := range neighbours {
	// 	fmt.Println(cell.Pos)
	// }
	//
	// 	select {
	// 	case <-updateChannel:
	//
	// 	}
	// 	cmd := exec.Command("clear")
	// 	cmd.Stdout = os.Stdout
	// 	cmd.Run()
	//
	// 	printWorld(&world)
	//
	// 	time.Sleep(1000 * time.Millisecond)
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
