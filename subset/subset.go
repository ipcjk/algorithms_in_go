package main

import (
	"fmt"
	"math"
)

func main() {
	generateSubsets([]int{}, []int{1, 2, 3}, "")
	generateSubsetsFromBitset([]string{"a", "b", "c"})
}

func generateSubsets(far, left []int, intend string) {
	fmt.Println(intend, "far:", far, "left", left)
	if len(left) == 0 {
		fmt.Println("finished:", far)
	} else {
		generateSubsets(append(far, left[0]), left[1:], intend+"   ")
		generateSubsets(far, left[1:], intend+"   ")
	}
}

func generateSubsetsFromBitset(subset []string) {
	y := math.Pow(2, float64(len(subset)))
	fmt.Println(len(subset))
	for i := 0; i < int(y); i++ {
		for j := uint32(0); j < uint32(len(subset)); j++ {
			if uint32(i)&(1<<j) > 0 {
				fmt.Printf("%s", subset[j])
			}
		}
		fmt.Println("")
	}

}
