package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"sudoku/internal/sudoku"
)

type solveRequest struct {
	Board [9][9]int `json:"board"`
}

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/", renderHome)
	r.GET("/easy", getPuzzle(sudoku.Easy))
	r.GET("/medium", getPuzzle(sudoku.Medium))
	r.GET("/hard", getPuzzle(sudoku.Hard))
	r.POST("/sudoku-solver", solveSudoku)

	return r
}

func renderHome(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Sudoku SSR",
	})
}

func getPuzzle(level sudoku.Difficulty) gin.HandlerFunc {
	return func(c *gin.Context) {
		puzzle := sudoku.GeneratePuzzle(level)
		c.JSON(http.StatusOK, gin.H{
			"difficulty": level.String(),
			"board":      puzzle,
		})
	}
}

func solveSudoku(c *gin.Context) {
	var req solveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload, expected { board: [[...9], ...9] }"})
		return
	}

	board := req.Board
	if !sudoku.IsValidBoard(board) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "board contains invalid numbers or conflicts"})
		return
	}

	solved := board
	if !sudoku.Solve(&solved) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "board cannot be solved"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"solution": solved,
	})
}
