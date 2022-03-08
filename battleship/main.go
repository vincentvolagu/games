package main

import "fmt"

func main() {
	ships := []Ship{
		Ship{"T", 2},
		Ship{"D", 3},
		Ship{"S", 3},
		Ship{"B", 4},
		Ship{"C", 5},
	}

	formation := MakeRandomFormation()
	summary := []int{}
	for i := 0; i < 100; i++ {
		board := &Board{}
		board.Init(20)
		formation.PlaceShips(board, ships)

		// board.Print()

		// fmt.Println("================================================")
		gunner := NewClusterGunner(board, []int{5, 4, 3, 2})
		// gunner := NewLinearGunner(board)
		// gunner := NewDiagonalGunner(board, []int{5, 2})
		// gunner := NewRandomGunner(board)
		steps := 0
		for {
			steps = steps + 1
			target := gunner.Target()
			if board.IsOutOfBound(target) {
				// fmt.Println("gunner finished, game over")
				break
			}
			result := board.Hit(target)
			if result {
				gunner.Hit(target)
			} else {
				gunner.Miss(target)
			}
			if board.IsGameOver() {
				break
			}
		}
		fmt.Println("takes", steps, "steps")
		summary = append(summary, steps)
		// board.Print()
		// fmt.Println("================================================")
	}

	sum := 0
	for _, v := range summary {
		sum = sum + v
	}
	avg := sum / len(summary)
	fmt.Println("Summary: avg = ", avg, "steps")
}
