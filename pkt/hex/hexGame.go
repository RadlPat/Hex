package hex

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
)

type Game struct {
	PlayerOne      Player
	PlayerTwo      Player
	ActivePlayer   *Player
	InactivePlayer *Player
	Rounds         int
	Finished       bool
	Board          HexBoard
}

type Move struct {
	i          int
	j          int
	Neighbours []Move
	Visited    bool
}

// NewGame init new Game
func NewGame(size int) Game {
	newGame := Game{}
	// create game board
	newGame.Board = NewHexBoard(size)

	//set player numbers
	PlayerOne := Player{}
	PlayerOne.Number = 1
	newGame.PlayerOne = PlayerOne
	PlayerTwo := Player{}
	PlayerTwo.Number = 2
	newGame.PlayerTwo = PlayerTwo

	newGame.Rounds = 1
	newGame.Finished = false

	return newGame
}

// Run runs the Game
func (g *Game) Run() {
	//start game by choosing start player
	g.chooseStartPlayer()

	//first player kann make the first move
	firstMove := g.ActivePlayer.firstMove()
	g.Board.ReceiveMove(firstMove, *g.ActivePlayer)
	g.Board.PrintBoard()

	g.swapRule(firstMove)

	for g.Finished == false {
		//as long as not finished:
		// get move
		isMoveValid := false
		move := Move{}

		for isMoveValid == false {
			move = g.ActivePlayer.getMove(&g.Board)
			isMoveValid = g.Board.ReceiveMove(move, *g.ActivePlayer)

		}
		g.ActivePlayer.Moves = append(g.ActivePlayer.Moves, move)
		// check if move is neighbour of other moves
		g.ActivePlayer.addNeighbours(move)

		g.Finished = g.ActivePlayer.isFinisherMove(g.Board.Size)
		if g.Finished {
			fmt.Print("Finished")
		}

		// print board
		g.Board.PrintBoard()

		// increase Round
		g.Rounds = g.Rounds + 1

		// check if finished

		// change player
		g.swapPlayers()
	}

}

// chooseStartPlayer sets randomly the Player who starts
func (g *Game) chooseStartPlayer() {
	//choose randomly start player
	randomNumber := rand.Intn(100)
	if randomNumber > 50 {
		g.ActivePlayer = &g.PlayerOne
		g.InactivePlayer = &g.PlayerTwo
	} else {
		g.ActivePlayer = &g.PlayerTwo
		g.InactivePlayer = &g.PlayerOne
	}

}

// swapRule Swap rule: on their first move the second player may move normally,
// or choose to swap their piece with that placed by the first player.
// [This encourages the first player to only choose a moderately strong first
// move and so reduces any advantage of going first. Ignore the swap rule for the first few games.]
func (g *Game) swapRule(move Move) {
	// ask if inactive player wants to swap pieces
	fmt.Println("Do you wanne swap pieces? yes/no")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	if input == "yes" || input == "y" {
		// add move to other player, change board
		g.InactivePlayer.Moves = append(g.InactivePlayer.Moves, move)
		g.Board.Board[move.i][move.j] = g.InactivePlayer.Number
	} else {
		g.ActivePlayer.Moves = append(g.ActivePlayer.Moves, move)
		g.swapPlayers()
	}

}

// swapPlayers swaps activ and inactiv players
func (g *Game) swapPlayers() {
	tmpAct := g.ActivePlayer
	g.ActivePlayer = g.InactivePlayer
	g.InactivePlayer = tmpAct
}

type Player struct {
	Number int
	Moves  []Move
}

// firstMove takes the first move from activ player
func (p *Player) firstMove() Move {
	// ask if inactive player wants to swap pieces
	fmt.Printf("Player %d make your first Move \n", p.Number)
	fmt.Print("Enter i: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	i := scanner.Text()
	fmt.Print("Enter j: ")
	scanner.Scan()
	j := scanner.Text()

	newMove := Move{}
	newMove.i, _ = strconv.Atoi(i)
	newMove.j, _ = strconv.Atoi(j)
	return newMove
}

// getMove ask for new move from active player
func (p *Player) getMove(g *HexBoard) Move {
	// ask if inactive player wants to swap pieces

	newMove := Move{}

	fmt.Printf("Player %d Make your Move \n", p.Number)
	fmt.Print("Enter i: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	i := scanner.Text()
	fmt.Print("Enter j: ")
	scanner.Scan()
	j := scanner.Text()

	newMove.i, _ = strconv.Atoi(i)
	newMove.j, _ = strconv.Atoi(j)

	return newMove

}

// addNeighbours checks if move is neighbour of older moves and
// adds it
func (p *Player) addNeighbours(newMove Move) {
	for i, _ := range p.Moves {
		//0 0 0 0 0
		//  0 0 0 0 0
		//   0 0 0 0 0
		//    0 0 0 0 0

		// 0 1 2 3
		//  1 2 3 4
		//   2 3 4 5
		//    3 4 5 6
		// 1-2, 1-2, 1-2 3-2 2-3 3-2
		newMoveSum := newMove.i + newMove.j
		moveSum := p.Moves[i].i + p.Moves[i].j

		if newMoveSum == moveSum {
			if p.Moves[i].i != newMove.i && p.Moves[i].j != newMove.j {
				// is neighbour
				p.Moves[len(p.Moves)-1].Neighbours = append(p.Moves[len(p.Moves)-1].Neighbours, p.Moves[i])
				p.Moves[i].Neighbours = append(p.Moves[i].Neighbours, newMove)
			}
		} else if newMoveSum+1 == moveSum {
			// is neighbour
			p.Moves[len(p.Moves)-1].Neighbours = append(p.Moves[len(p.Moves)-1].Neighbours, p.Moves[i])
			p.Moves[i].Neighbours = append(p.Moves[i].Neighbours, newMove)
		} else if newMoveSum == moveSum+1 {
			// is neighbour
			p.Moves[len(p.Moves)-1].Neighbours = append(p.Moves[len(p.Moves)-1].Neighbours, p.Moves[i])
			p.Moves[i].Neighbours = append(p.Moves[i].Neighbours, newMove)
		}

	}

}

// Stack implementation in go
type Stack struct {
	items  []Move
	rwLock sync.RWMutex
}

// Push pushs adds an element to the Stack
func (s *Stack) Push(t Move) {
	if s.items == nil {
		s.items = []Move{}
	}
	s.rwLock.Lock()
	s.items = append(s.items, t)
	s.rwLock.Unlock()
}

// Pop gets the last item and deletes it from Stack
func (s *Stack) Pop() *Move {
	if len(s.items) == 0 {
		return nil
	}
	s.rwLock.Lock()
	item := s.items[len(s.items)-1]
	s.items = s.items[0 : len(s.items)-1]
	s.rwLock.Unlock()
	return &item

}

// Connects checks if there is a connection between opposite boarders
func (s *Stack) Connects(size int) bool {
	top := false
	bottom := false
	right := false
	left := false

	for _, item := range s.items {
		if item.i == 0 {
			top = true
		}
		if item.i == size-1 {
			bottom = true
		}
		if item.j == 0 {
			left = true
		}
		if item.j == size-1 {
			right = true
		}
	}

	if top && bottom {
		return true
	} else if right && left {
		return true
	}
	return false
}

// is FinisherMove checks of there is a connection between two moves
// from the oposit boarders
func (p *Player) isFinisherMove(boardSize int) bool {

	var stack Stack
	var visited Stack
	move := p.Moves[len(p.Moves)-1]
	//init visited = false for all moves
	for i, _ := range p.Moves {
		p.Moves[i].Visited = false
	}

	stack.Push(move)

	//go trough all moves which are conected
	for len(stack.items) > 0 {
		moveToCheck := stack.Pop()
		if moveToCheck.Visited == false {
			moveToCheck.Visited = true
			visited.Push(*moveToCheck)

			//add all nightbours to the stack to get subgraph in visited
			for _, neigbourMoves := range moveToCheck.Neighbours {
				stack.Push(neigbourMoves)
			}

		}
		fmt.Println("------------------------------------")
	}

	// todo board size
	finished := visited.Connects(boardSize)
	return finished

}
