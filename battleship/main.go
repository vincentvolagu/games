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

	playComputer(ships)
	// playHuman(ships)
}

func playComputer(ships []Ship) {
	// formation := MakeRandomFormation()
	formation := MakeRandomCoordinator()
	summary := []int{}
	for i := 0; i < 1; i++ {
		board := &Board{}
		board.Init(10)
		formation.PlaceShips(board, ships)

		board.Print()

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
		board.Print()
		fmt.Println("takes", steps, "steps")
		summary = append(summary, steps)
		fmt.Println("================================================")
	}

	sum := 0
	for _, v := range summary {
		sum = sum + v
	}
	avg := sum / len(summary)
	fmt.Println("Summary: avg = ", avg, "steps")
}

func playHuman(ships []Ship) {
	// initialise the computer board and place ships
	computerBoard := &Board{}
	computerBoard.Init(10)
	formation := MakeRandomFormation()
	formation.PlaceShips(computerBoard, ships)
	computerBoard.Print()
	computerBoard.PrintForHumanPlayer()

	humanBoard := &Board{}
	humanBoard.Init(10)
	gunner := NewClusterGunner(humanBoard, []int{5, 4, 3, 2})

	// initialise the human player board and ask player to place all ships
	totalPoints := 0
	for _, ship := range ships {
		totalPoints += ship.length
	}
	for i := 0; i < totalPoints; i++ {
		// take human player's next target
		fmt.Println("================================================")
		fmt.Println("You have the following ships to put on the board")
		for _, ship := range ships {
			fmt.Println(ship.name, "length", ship.length)
		}
		fmt.Println("================================================")
		fmt.Println("enter row number")
		var row int
		fmt.Scanln(&row)
		fmt.Println("enter column")
		var col string
		fmt.Scanln(&col)

		p := computerBoard.TransformHumanPointInput(row, col)
		humanBoard.PlaceShipAt(p, "S")
		fmt.Println("================================================")
		humanBoard.Print()
	}

	for {
		// take human player's next target
		fmt.Println("================================================")
		fmt.Println("enter row number")
		var row int
		fmt.Scanln(&row)

		fmt.Println("enter column")
		var col string
		fmt.Scanln(&col)

		target := computerBoard.TransformHumanPointInput(row, col)

		if computerBoard.IsOutOfBound(target) {
			fmt.Println("you target", row, col, "is a invalid point, please try again")
			continue
		}

		isHit := computerBoard.Hit(target)
		fmt.Println("you target", row, col, "is", isHit)
		fmt.Println("============== Computer's Board =============================")
		computerBoard.PrintForHumanPlayer()

		if computerBoard.IsGameOver() {
			fmt.Println("Congrats, you win, thanks for playing")
			break
		}

		// computer's turn
		target = gunner.Target()
		if humanBoard.IsOutOfBound(target) {
			fmt.Println("I'm out of guesses, you win, thanks for playing")
			break
		}

		hit := humanBoard.Hit(target)
		if hit {
			gunner.Hit(target)
		} else {
			gunner.Miss(target)
		}
		fmt.Println("============== Player's Board =============================")
		humanBoard.PrintForHumanPlayer()

		if humanBoard.IsGameOver() {
			fmt.Println("Oops, you lose, thanks for playing")
			break
		}
	}
}
