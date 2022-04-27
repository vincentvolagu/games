package main

import "fmt"
import "math"

func main() {
	fmt.Println("vim-go, 16! = ", factorial(16))
	// bruteForce(16)

}

func bruteForce(n int) {
	list := make([]int, n)
	for i := 1; i <= n; i++ {
		list[i-1] = i
	}
	combo := factorialCombination(list)
	for _, row := range combo {
		// fmt.Println(row)
		if verify(row) {
			fmt.Println("found list that are sequence of adjacent square numbers", row)
			break
		}
	}
}

func factorial(n int) int {
	if n == 1 {
		return 1
	}
	return factorial(n-1) * n
}

func factorialCombination(list []int) [][]int {
	if len(list) == 1 {
		return [][]int{list}
	}
	// newCombo := make([][]int, factorial(len(list)))
	var newCombo [][]int
	for i, v := range list {
		head := v
		// cheap way to remove a value from slice, changing the ordering though
		tail := make([]int, len(list))
		copy(tail, list)
		tail[i] = tail[len(list)-1]
		tail = tail[:len(tail)-1]
		tailCombo := factorialCombination(tail)
		// fmt.Println("create combo for ", head, " by ", tail, " with tail ", tailCombo)
		for _, combo := range tailCombo {
			newCombo = append(newCombo, append([]int{head}, combo...))
			// fmt.Println("new combo ", newCombo[i])
		}
	}

	return newCombo
}

func verify(orderedList []int) bool {
	for i, v := range orderedList {
		if !isSquareNumber((float64)(v + orderedList[i+1])) {
			return false
		}
	}
	return true
}

func isSquareNumber(n float64) bool {
	return math.Floor(math.Sqrt(n)) == n
}
