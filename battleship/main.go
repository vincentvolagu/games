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

	for i := 0; i < 3; i++ {
		board := &Board{}
		board.Init(10)
		formation.PlaceShips(board, ships)

		board.Print()

		fmt.Println("================================================")
		gunner := NewClusterGunner(board)
		// gunner := NewLinearGunner(board)
		// gunner := NewDiagonalGunner(board, []int{5, 2})
		// gunner := NewRandomGunner(board)
		steps := 0
		for {
			steps = steps + 1
			target := gunner.Target()
			if board.IsOutOfBound(target) {
				fmt.Println("gunner finished, game over")
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
		board.Print()
		fmt.Println("================================================")
	}
}
