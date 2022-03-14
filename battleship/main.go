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

	playHuman(ships)
	// formation := MakeRandomFormation()
	// summary := []int{}
	// for i := 0; i < 1; i++ {
	// board := &Board{}
	// board.Init(10)
	// formation.PlaceShips(board, ships)

	// board.PrintForHumanPlayer()

	// // fmt.Println("================================================")
	// gunner := NewClusterGunner(board, []int{5, 4, 3, 2})
	// // gunner := NewLinearGunner(board)
	// // gunner := NewDiagonalGunner(board, []int{5, 2})
	// // gunner := NewRandomGunner(board)
	// steps := 0
	// for {
	// steps = steps + 1
	// target := gunner.Target()
	// if board.IsOutOfBound(target) {
	// // fmt.Println("gunner finished, game over")
	// break
	// }
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
	// fmt.Println("takes", steps, "steps")
	// summary = append(summary, steps)
	// // board.Print()
	// // fmt.Println("================================================")
	// }

	// sum := 0
	// for _, v := range summary {
	// sum = sum + v
	// }
	// avg := sum / len(summary)
	// fmt.Println("Summary: avg = ", avg, "steps")
}

func playHuman(ships []Ship) {
	// initialise the computer board and place ships
	myBoard := &Board{}
	myBoard.Init(10)
	formation := MakeRandomFormation()
	formation.PlaceShips(myBoard, ships)
	myBoard.Print()
	myBoard.PrintForHumanPlayer()

	humanBoard := &Board{}
	humanBoard.Init(10)
	humanBoard.SetShips(ships)
	gunner := NewClusterGunner(humanBoard, []int{5, 4, 3, 2})

	for {
		// take human player's next target
		fmt.Println("================================================")
		fmt.Println("enter row number")
		var row int
		fmt.Scanln(&row)

		fmt.Println("enter column")
		var col string
		fmt.Scanln(&col)

		p := myBoard.TransformHumanPointInput(row, col)

		if myBoard.IsOutOfBound(p) {
			fmt.Println("you target", row, col, "is a invalid point, please try again")
			continue
		}

		isHit := myBoard.Hit(p)
		fmt.Println("you target", row, col, "is", isHit)
		fmt.Println("============== Computer's Board =============================")
		myBoard.PrintForHumanPlayer()

		if myBoard.IsGameOver() {
			fmt.Println("Congrats, you win, thanks for playing")
			break
		}

		// computer's turn
		target := gunner.Target()
		if humanBoard.IsOutOfBound(target) {
			fmt.Println("I'm out of guesses, you win, thanks for playing")
			break
		}
		targetRow, targetCol := humanBoard.TransformPointForHuman(target)
		fmt.Println("Computer picks", target, targetRow, targetCol, "please enter if it's a hit [y/n]")
		var hit string
		fmt.Scanln(&hit)
		if hit == "y" {
			humanBoard.RecordHit(target)
			gunner.Hit(target)
		} else {
			humanBoard.RecordMiss(target)
			gunner.Miss(target)
		}
		fmt.Println("============== Player's Board =============================")
		humanBoard.PrintForHumanPlayer()

		if humanBoard.HasHitAllShips() {
			fmt.Println("Oops, you lose, thanks for playing")
			break
		}
	}
}
