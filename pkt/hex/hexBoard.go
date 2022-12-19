package hex

import "fmt"

type HexBoard struct {
	Board [][]int
	Size  int
}

func NewHexBoard(size int) HexBoard {
	newBord := HexBoard{}
	newBord.CreateBoard(size)
	newBord.Size = size
	return newBord
}

// createBoard creates the game board in n x n dimension
func (h *HexBoard) CreateBoard(n int) {
	h.Board = make([][]int, n)
	h.Size = n
	for i := 0; i < h.Size; i++ {
		h.Board[i] = make([]int, h.Size)

		for j := 0; j < h.Size; j++ {
			h.Board[i][j] = 0
		}
	}
}

// printBoard prints the game board in terminal in a readable view
func (h *HexBoard) PrintBoard() {
	for i := 0; i < h.Size; i++ {
		fmt.Println(h.Board[i])
	}

}

// receiveMove sets move on Board if possible, else it returns false
func (h *HexBoard) ReceiveMove(move Move, player Player) bool {
	if h.Board[move.i][move.j] == 0 {
		h.Board[move.i][move.j] = player.Number
		return true
	}

	return false

}
