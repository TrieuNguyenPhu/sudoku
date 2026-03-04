package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"sudoku/internal/models"
	"sudoku/internal/services"
)

// SolverHandler handles Sudoku solving requests
type SolverHandler struct {
	solverService services.SolverService
}

// NewSolverHandler creates a new solver handler
func NewSolverHandler(solverService services.SolverService) *SolverHandler {
	return &SolverHandler{
		solverService: solverService,
	}
}

// SolvePuzzle handles POST /sudoku-solver requests (requirement 4.1, 4.2)
func (h *SolverHandler) SolvePuzzle(c *gin.Context) {
	var request models.SolverRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&request); err != nil {
		response := models.NewErrorResponse(http.StatusBadRequest, "Invalid JSON", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Solve the puzzle using enhanced validation and error handling
	solution, err := services.SolveFromRequest(h.solverService, &request)
	if err != nil {
		// Handle specific error cases as per requirements
		if err.Error() == "Puzzle has no solution" {
			// Requirement 4.3: Return HTTP 400 for unsolvable puzzles
			response := models.NewErrorResponse(http.StatusBadRequest, "Unsolvable Puzzle", "Puzzle has no solution")
			c.JSON(http.StatusBadRequest, response)
			return
		}
		
		if err.Error() == "solving timeout exceeded" {
			// Requirement 4.5: Handle timeout errors
			response := models.NewErrorResponse(http.StatusRequestTimeout, "Timeout", "Solving timeout exceeded (3000ms)")
			c.JSON(http.StatusRequestTimeout, response)
			return
		}

		// Requirement 4.6: Handle validation errors
		response := models.NewErrorResponse(http.StatusBadRequest, "Validation Error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Return successful solution
	response := models.NewSolverResponse(solution, true, "Puzzle solved successfully")
	c.JSON(http.StatusOK, response)
}