package main

import (
	"fmt"
	"time"
	"os/exec"
	"os"
)


type Board struct {
	Feld [8][8]int
}

func (b *Board) print() {
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLUMNS; j++ {
			if b.Feld[i][j] == 1 {
				fmt.Printf(" Q") } else {
				fmt.Printf(" -")
			}
		}
		fmt.Println()
	}
	fmt.Println("")
}

func (b *Board) place(x, y int) bool {

	b.print()
	time.Sleep(time.Millisecond*500)
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	if x > ROWS || y > COLUMNS {
		return false
	}

	if !b.is_safe(x, y) {
		return false
	}

	if b.Feld[x][y] == 0 {
		b.Feld[x][y] = 1
		return true
	}

	return false

}

func (b *Board) unplace(x, y int) bool {

	if b.Feld[x][y] == 1 {
		b.Feld[x][y] = 0
		return true
	}
	return false

}

func (b *Board) is_safe(x, y int) bool {
	/* no queen allowed in same row or same column */

	/* check column first */
	for i := 0; i < ROWS; i++ {
		if b.Feld[i][y] == 1 {
			return false
		}
	}

	/* check rows second */
	for i := 0; i < COLUMNS; i++ {
		if b.Feld[x][i] == 1 {
			return false
		}
	}


	/* check diagonal crosses up, left */
	for j, i := y-1, x-1; i >= 0 && j >= 0 ; i-- {
		if b.Feld[i][j] == 1 {
			//fmt.Println("Up/Left", i, j, "cant place at", x, y)
			return false
		}
		j--
	}

	/* check diagonal crosses up, right */
	for j, i := y+1, x-1; i >= 0 && j < COLUMNS ; i-- {
		//fmt.Println("Up/Right","checking",i,j,"for",x,y)
		if b.Feld[i][j] == 1 {
			//fmt.Println("Up/Right", i, j, "cant place at", x, y)
			return false
		}
		j++
	}

	/* check diagonal crosses down, left */
	for j, i := y-1, x+1; i < ROWS && j >= 0 ; i++ {
		if b.Feld[i][j] == 1 {
			//fmt.Println("Down/Left", i, j, "cant place at", x, y)
			return false
		}
		j--
	}


	/* check diagonal crosses down, right */
	for j, i := y+1, x+1; i < ROWS && j < COLUMNS ; i++ {
		if b.Feld[i][j] == 1 {
			//fmt.Println("Down/Right", i, j, "cant place at", x, y)
			return false
		}
		j++
	}


	return true

}
