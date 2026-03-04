package sudoku

import (
	"math/rand"
	"time"
)

type Difficulty int

const (
	Easy Difficulty = iota
	Medium
	Hard
)

func (d Difficulty) String() string {
	switch d {
	case Easy:
		return "easy"
	case Medium:
		return "medium"
	case Hard:
		return "hard"
	default:
		return "unknown"
	}
}

func GeneratePuzzle(d Difficulty) [9][9]int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	base := seededSolvedBoard(r)

	var holes int
	switch d {
	case Easy:
		holes = 36
	case Medium:
		holes = 46
	case Hard:
		holes = 54
	default:
		holes = 46
	}

	cells := r.Perm(81)
	for i := 0; i < holes && i < len(cells); i++ {
		pos := cells[i]
		row := pos / 9
		col := pos % 9
		base[row][col] = 0
	}

	return base
}

func seededSolvedBoard(r *rand.Rand) [9][9]int {
	var b [9][9]int
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			b[row][col] = (row*3+row/3+col)%9 + 1
		}
	}

	shuffleRowsWithinBands(&b, r)
	shuffleColsWithinStacks(&b, r)
	shuffleBands(&b, r)
	shuffleStacks(&b, r)
	permuteDigits(&b, r)

	return b
}

func shuffleRowsWithinBands(b *[9][9]int, r *rand.Rand) {
	for band := 0; band < 3; band++ {
		start := band * 3
		order := r.Perm(3)
		var temp [3][9]int
		for i := 0; i < 3; i++ {
			temp[i] = b[start+order[i]]
		}
		for i := 0; i < 3; i++ {
			b[start+i] = temp[i]
		}
	}
}

func shuffleColsWithinStacks(b *[9][9]int, r *rand.Rand) {
	for stack := 0; stack < 3; stack++ {
		start := stack * 3
		order := r.Perm(3)
		for row := 0; row < 9; row++ {
			var temp [3]int
			for i := 0; i < 3; i++ {
				temp[i] = b[row][start+order[i]]
			}
			for i := 0; i < 3; i++ {
				b[row][start+i] = temp[i]
			}
		}
	}
}

func shuffleBands(b *[9][9]int, r *rand.Rand) {
	order := r.Perm(3)
	var temp [9][9]int
	for band := 0; band < 3; band++ {
		sourceBand := order[band]
		for i := 0; i < 3; i++ {
			temp[band*3+i] = b[sourceBand*3+i]
		}
	}
	*b = temp
}

func shuffleStacks(b *[9][9]int, r *rand.Rand) {
	order := r.Perm(3)
	var temp [9][9]int
	for row := 0; row < 9; row++ {
		for stack := 0; stack < 3; stack++ {
			sourceStack := order[stack]
			for i := 0; i < 3; i++ {
				temp[row][stack*3+i] = b[row][sourceStack*3+i]
			}
		}
	}
	*b = temp
}

func permuteDigits(b *[9][9]int, r *rand.Rand) {
	perm := r.Perm(9)
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			b[row][col] = perm[b[row][col]-1] + 1
		}
	}
}

func IsValidBoard(board [9][9]int) bool {
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			val := board[row][col]
			if val < 0 || val > 9 {
				return false
			}
			if val == 0 {
				continue
			}

			board[row][col] = 0
			if !isSafe(board, row, col, val) {
				return false
			}
			board[row][col] = val
		}
	}
	return true
}

func Solve(board *[9][9]int) bool {
	row, col, ok := findEmpty(*board)
	if !ok {
		return true
	}

	for n := 1; n <= 9; n++ {
		if isSafe(*board, row, col, n) {
			board[row][col] = n
			if Solve(board) {
				return true
			}
			board[row][col] = 0
		}
	}
	return false
}

func findEmpty(board [9][9]int) (row int, col int, ok bool) {
	for row = 0; row < 9; row++ {
		for col = 0; col < 9; col++ {
			if board[row][col] == 0 {
				return row, col, true
			}
		}
	}
	return 0, 0, false
}

func isSafe(board [9][9]int, row int, col int, n int) bool {
	for i := 0; i < 9; i++ {
		if board[row][i] == n || board[i][col] == n {
			return false
		}
	}

	startRow := (row / 3) * 3
	startCol := (col / 3) * 3
	for r := startRow; r < startRow+3; r++ {
		for c := startCol; c < startCol+3; c++ {
			if board[r][c] == n {
				return false
			}
		}
	}
	return true
}
