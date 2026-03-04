package handlers

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all application routes
func RegisterRoutes(router *gin.Engine) {
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

	// API routes (placeholder for now)
	api := router.Group("/api")
	{
		api.GET("/easy", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Easy puzzle endpoint - to be implemented"})
		})
		api.GET("/medium", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Medium puzzle endpoint - to be implemented"})
		})
		api.GET("/hard", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Hard puzzle endpoint - to be implemented"})
		})
		api.POST("/sudoku-solver", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Solver endpoint - to be implemented"})
		})
	}
}