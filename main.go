package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

// Any live cell with fewer than two live neighbours dies, as if by underpopulation.
// Any live cell with two or three live neighbours lives on to the next generation.
// Any live cell with more than three live neighbours dies, as if by overpopulation.
// Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.

type Pos struct {
	x uint
	y uint
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

	// updateNeighbours(c.Pos, 1)
}

func NewCell(alive bool, x uint, y uint, alives *[]*Cell) *Cell {
	c := Cell{
		Alive: alive,
		Pos: Pos{
			x: x,
			y: y,
		},
		Channel: make(chan int, 1),
	}

	*alives = append(*alives, &c)

	return &c
}

func GetPos() *Pos {
	p := Pos{
		x: 4,
		y: 10,
	}

	return &p
}

// func updateNeighbours(pos Pos, val int) {
// 	minx := max(0, pos.x-3)
// 	maxx := min()
//
// }

type world struct {
	width  uint
	height uint
	alives []*Cell
	world  [][]*Cell
}

func main() {
	width := uint(100)
	height := uint(40)
	alives := []*Cell{}

	world := make([][]*Cell, height)

	updateChannel := make(chan int, 1)

	for y := uint(0); y < height; y++ {
		world[y] = make([]*Cell, width)

		for x := uint(0); x < width; x++ {
			var alive bool
			chance := rand.Intn(100)

			if chance == 1 {
				alive = true
			} else {
				alive = false
			}

			world[y][x] = NewCell(alive, x, y, &alives)
		}
	}

	for {
		select {
		case <-updateChannel:

		}
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()

		printWorld(&world)

		time.Sleep(1000 * time.Millisecond)
	}
}

func printWorld(w *[][]*Cell) {
	for _, row := range *w {
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
