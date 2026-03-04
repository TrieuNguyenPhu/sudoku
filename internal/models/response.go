package models

import (
	"time"
)

// ErrorResponse represents a standardized error response structure
type ErrorResponse struct {
	Error     string    `json:"error"`
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	RequestID string    `json:"request_id,omitempty"`
}

// PuzzleRequest represents a request to generate a puzzle
type PuzzleRequest struct {
	Difficulty Difficulty `json:"difficulty" binding:"required"`
}

// SolverRequest represents a request to solve a Sudoku puzzle
type SolverRequest struct {
	Grid [9][9]int `json:"grid" binding:"required"`
}

// PuzzleResponse represents the response for puzzle generation
type PuzzleResponse struct {
	Puzzle    *SudokuPuzzle `json:"puzzle"`
	Success   bool          `json:"success"`
	Message   string        `json:"message,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
}

// SolverResponse represents the response for puzzle solving
type SolverResponse struct {
	Solution  *SudokuSolution `json:"solution,omitempty"`
	Success   bool            `json:"success"`
	Message   string          `json:"message,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
}

// NewErrorResponse creates a new error response with timestamp
func NewErrorResponse(code int, error string, message string) *ErrorResponse {
	return &ErrorResponse{
		Error:     error,
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
	}
}

// NewErrorResponseWithRequestID creates a new error response with request ID
func NewErrorResponseWithRequestID(code int, error string, message string, requestID string) *ErrorResponse {
	return &ErrorResponse{
		Error:     error,
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		RequestID: requestID,
	}
}

// NewPuzzleResponse creates a new puzzle response
func NewPuzzleResponse(puzzle *SudokuPuzzle, success bool, message string) *PuzzleResponse {
	return &PuzzleResponse{
		Puzzle:    puzzle,
		Success:   success,
		Message:   message,
		Timestamp: time.Now(),
	}
}

// NewSolverResponse creates a new solver response
func NewSolverResponse(solution *SudokuSolution, success bool, message string) *SolverResponse {
	return &SolverResponse{
		Solution:  solution,
		Success:   success,
		Message:   message,
		Timestamp: time.Now(),
	}
}

// ValidateSolverRequest validates the solver request structure
func (r *SolverRequest) Validate() error {
	// Create a Grid from the request data for validation
	var grid Grid
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			grid[i][j] = r.Grid[i][j]
		}
	}
	
	return grid.ValidateGrid()
}