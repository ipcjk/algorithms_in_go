package main

import (
	"fmt"
	"math"
)

func main() {
	generateSubsets([]int{}, []int{1,2})
	generateSubsetsFromBitset([]int{1,2,3})
}

func generateSubsets(far, left []int) {
	if len(left) == 0 {
		fmt.Println(far)
	} else {
		generateSubsets(append(far, left[0]), left[1:])
		generateSubsets(far, left[1:])
	}
}

func generateSubsetsFromBitset(subset []int) {
	y := math.Pow(2, float64(len(subset)))

	for i := 0; i < int(y); i++ {
		for j := uint32(0); j < uint32(len(subset)); j++ {
			if uint32(i) & (1<<j) > 0 {
				fmt.Printf("%d ", subset[j])
			}
		}
		fmt.Println("")
	}

}
