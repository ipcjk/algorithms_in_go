package main

import "fmt"

func main() {
	generateSubsets([]int{}, []int{1,2,3})
}

func generateSubsets(far, left []int) {
	if len(left) == 0 {
		fmt.Println(far)
	} else {
		generateSubsets(append(far, left[0]), left[1:])
		generateSubsets(far, left[1:])
	}
}
