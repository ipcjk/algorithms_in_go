package main

import "fmt"

const (
	ROWS    = 8
	COLUMNS = 8
)

var calls = 0

func main() {
	var b Board

	eightqueens(&b, 8, 0)

}

func eightqueens(b *Board, queen int, colplaced int) bool {
	calls++
	if queen == 0 {
		fmt.Println(calls)
		b.print()
		return true
	} else {

		for i := 0; i < ROWS; i++ {
			if b.is_safe(i, colplaced) {
				b.place(i, colplaced)
				finished := eightqueens(b, queen-1, colplaced+1)
				if finished {
					return true
				}
				b.unplace(i, colplaced)
			}
		}

	}
	return false
}
