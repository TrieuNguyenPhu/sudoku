package handlers

import (
	"sudoku/internal/services"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all application routes with proper service initialization
func RegisterRoutes(router *gin.Engine) {
	// Initialize services
	solverService := services.NewSolverService()
	puzzleService := services.NewPuzzleService(solverService)
	
	// Initialize handlers
	puzzleHandler := NewPuzzleHandler(puzzleService)
	solverHandler := NewSolverHandler(solverService)

	// Web interface routes
	router.GET("/", func(c *gin.Context) {
		// Create arrays for template iteration
		rows := make([]int, 9)
		cols := make([]int, 9)
		for i := 0; i < 9; i++ {
			rows[i] = i
			cols[i] = i
		}
		
		c.HTML(200, "game.html", gin.H{
			"title": "Sudoku Game",
			"Rows":  rows,
			"Cols":  cols,
		})
	})

	// API routes for puzzle generation (requirements 1, 2, 3)
	router.GET("/easy", puzzleHandler.GenerateEasy)
	router.GET("/medium", puzzleHandler.GenerateMedium)
	router.GET("/hard", puzzleHandler.GenerateHard)
	
	// API route for puzzle solving (requirement 4)
	router.POST("/sudoku-solver", solverHandler.SolvePuzzle)

	// Legacy API routes (for backward compatibility)
	api := router.Group("/api")
	{
		api.GET("/easy", puzzleHandler.GenerateEasy)
		api.GET("/medium", puzzleHandler.GenerateMedium)
		api.GET("/hard", puzzleHandler.GenerateHard)
		api.POST("/sudoku-solver", solverHandler.SolvePuzzle)
	}
}