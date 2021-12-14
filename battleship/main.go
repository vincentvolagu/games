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

	formation := makeRandomCoordinator()

	for i := 0; i < 3; i++ {
		board := &Board{}
		board.Init(10)
		formation.PlaceShips(board, ships)

		board.Print()

		fmt.Println("================================================")
		// gunner := NewLinearGunner(board)
		// for {
		// target := gunner.Target()
		// result := board.Hit(target)
		// if result {
		// gunner.Hit(target)
		// } else {
		// gunner.Miss(target)
		// }
		// if board.IsGameOver() {
		// break
		// }
		// }
		// board.Print()
		// fmt.Println("================================================")
	}

}
