package main

// import "fmt"
import "math"
import "math/rand"

type Random struct {
}

// runs in O(n!), so expect this to blow up, lol
func (r Random) SquareSumSort(n int) [][]int {
	list := make([]int, n)
	for i := 1; i <= n; i++ {
		list[i-1] = i
	}

	// randomly swap numbers until we found a sequence that works
	for {
		index1 := rand.Intn(n)
		index2 := rand.Intn(n)
		tmp := list[index1]
		list[index1] = list[index2]
		list[index2] = tmp
		// fmt.Println("trying", list)
		if r.verify(list) {
			return [][]int{list}
		}
	}
}
func (r Random) verify(orderedList []int) bool {
	for i, v := range orderedList {
		if i == 0 {
			continue
		}
		if !r.isSquareNumber((float64)(v + orderedList[i-1])) {
			return false
		}
	}
	return true
}

func (r Random) isSquareNumber(n float64) bool {
	sqrt := math.Sqrt(n)
	return math.Floor(sqrt) == sqrt
}

func factorial(n int) int {
	if n == 1 {
		return 1
	}
	return factorial(n-1) * n
}
