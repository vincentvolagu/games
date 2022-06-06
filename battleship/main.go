package main

import "fmt"

func main() {
	ships := []Ship{
		Ship{"C", 5},
		Ship{"B", 4},
		Ship{"S", 3},
		Ship{"D", 3},
		Ship{"T", 2},
	}

	st := &stats{name: "lucky draw vs lucky gunner"}
	playComputerPair(ships, MakeLuckyDraw(), LukcyGunnerFactory{}, st)
	st.print()
	st = &stats{name: "lucky draw vs linear gunner"}
	playComputerPair(ships, MakeLuckyDraw(), LinearGunnerFactory{}, st)
	st.print()
	st = &stats{name: "lucky draw vs cluster gunner"}
	playComputerPair(ships, MakeLuckyDraw(), ClusterGunnerFactory{}, st)
	st.print()
	st = &stats{name: "edge lover vs lucky gunner"}
	playComputerPair(ships, EdgeLover{}, LukcyGunnerFactory{}, st)
	st.print()
	st = &stats{name: "edge lover vs linear gunner"}
	playComputerPair(ships, EdgeLover{}, LinearGunnerFactory{}, st)
	st.print()
	st = &stats{name: "edge lover vs cluster gunner"}
	playComputerPair(ships, EdgeLover{}, ClusterGunnerFactory{}, st)
	st.print()
	// playHuman(ships)
}

type stats struct {
	name          string // name of the pair
	numberOfGames int
	totalSteps    int
}

func (s stats) avgSteps() int {
	return s.totalSteps / s.numberOfGames
}
func (s stats) print() {
	fmt.Println("Summary:", s.name)
	fmt.Println("total number of games:", s.numberOfGames)
	fmt.Println("total steps:", s.totalSteps)
	fmt.Println("avg steps per game:", s.avgSteps(), "steps")
	fmt.Println("=======================================")
}

func playComputerPair(
	ships []Ship,
	coordinator ShipCoordinator,
	gunnerFactory GunnerFactory,
	stats *stats,
) {
	for i := 0; i < 20; i++ {
		board := &Board{}
		board.Init(10)
		coordinator.PlaceShips(board, ships)

		// board.Print()

		// TODO: extract the hard coded ship size slice here
		gunner := gunnerFactory.MakeGunner(board, []int{5, 4, 3, 2})
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
		stats.numberOfGames = stats.numberOfGames + 1
		stats.totalSteps += steps
		// board.Print()
		// fmt.Println("================================================")
	}
}

func playHuman(ships []Ship) {
	// initialise the computer board and place ships
	computerBoard := &Board{}
	computerBoard.Init(10)
	formation := MakeLuckyDraw()
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
