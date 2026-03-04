package handlers

import (
	"net/http"

	"sudoku/internal/models"
	"sudoku/internal/services"

	"github.com/gin-gonic/gin"
)

// PuzzleHandler handles HTTP requests for puzzle generation
type PuzzleHandler struct {
	puzzleService services.PuzzleService
}

// NewPuzzleHandler creates a new PuzzleHandler
func NewPuzzleHandler(puzzleService services.PuzzleService) *PuzzleHandler {
	return &PuzzleHandler{
		puzzleService: puzzleService,
	}
}

// GenerateEasy handles GET /easy requests (requirement 1.1, 1.2, 1.4, 1.5)
func (h *PuzzleHandler) GenerateEasy(c *gin.Context) {
	puzzle, err := h.puzzleService.GenerateEasy()
	if err != nil {
		// Return HTTP 500 for generation failures (requirement 1.5)
		errorResponse := models.NewErrorResponse(
			http.StatusInternalServerError,
			"Generation Failed",
			err.Error(),
		)
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	// Return puzzle in JSON format (requirement 1.2)
	response := models.NewPuzzleResponse(puzzle, true, "Easy puzzle generated successfully")
	c.JSON(http.StatusOK, response)
}

// GenerateMedium handles GET /medium requests (requirement 2.1, 2.2, 2.4, 2.5)
func (h *PuzzleHandler) GenerateMedium(c *gin.Context) {
	puzzle, err := h.puzzleService.GenerateMedium()
	if err != nil {
		// Return HTTP 500 for generation failures (requirement 2.5)
		errorResponse := models.NewErrorResponse(
			http.StatusInternalServerError,
			"Generation Failed",
			err.Error(),
		)
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	// Return puzzle in JSON format (requirement 2.2)
	response := models.NewPuzzleResponse(puzzle, true, "Medium puzzle generated successfully")
	c.JSON(http.StatusOK, response)
}

// GenerateHard handles GET /hard requests (requirement 3.1, 3.2, 3.4, 3.5)
func (h *PuzzleHandler) GenerateHard(c *gin.Context) {
	puzzle, err := h.puzzleService.GenerateHard()
	if err != nil {
		// Return HTTP 500 for generation failures (requirement 3.5)
		errorResponse := models.NewErrorResponse(
			http.StatusInternalServerError,
			"Generation Failed",
			err.Error(),
		)
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	// Return puzzle in JSON format (requirement 3.2)
	response := models.NewPuzzleResponse(puzzle, true, "Hard puzzle generated successfully")
	c.JSON(http.StatusOK, response)
}