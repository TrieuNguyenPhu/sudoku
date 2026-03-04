package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"sudoku/internal/models"
)

// SolverService defines the interface for Sudoku solving operations
type SolverService interface {
	Solve(puzzle *models.SudokuPuzzle) (*models.SudokuSolution, error)
	SolveWithTimeout(puzzle *models.SudokuPuzzle, timeout time.Duration) (*models.SudokuSolution, error)
	Validate(puzzle *models.SudokuPuzzle) error
	ValidateForSolving(request *models.SolverRequest) error
	IsValidMove(grid *models.Grid, row, col, num int) bool
}

// solverService implements the SolverService interface
type solverService struct{}

// NewSolverService creates a new instance of SolverService
func NewSolverService() SolverService {
	return &solverService{}
}

// Solve solves a Sudoku puzzle using backtracking algorithm (requirement 8.4)
func (s *solverService) Solve(puzzle *models.SudokuPuzzle) (*models.SudokuSolution, error) {
	// Default timeout of 3000ms as per requirement 4.5
	return s.SolveWithTimeout(puzzle, 3000*time.Millisecond)
}

// SolveWithTimeout solves a Sudoku puzzle with a specified timeout
func (s *solverService) SolveWithTimeout(puzzle *models.SudokuPuzzle, timeout time.Duration) (*models.SudokuSolution, error) {
	if puzzle == nil {
		return nil, errors.New("puzzle cannot be nil")
	}

	// Validate the input puzzle first
	if err := s.Validate(puzzle); err != nil {
		return nil, err
	}

	// Create a copy of the grid to solve
	solutionGrid := puzzle.Grid

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Solve using backtracking with timeout
	solved := s.solveWithContext(ctx, &solutionGrid)
	
	// Check if context was cancelled (timeout)
	if ctx.Err() == context.DeadlineExceeded {
		return nil, errors.New("solving timeout exceeded")
	}

	if !solved {
		return nil, errors.New("Puzzle has no solution")
	}

	// Create and return the solution
	solution := &models.SudokuSolution{
		Grid:     solutionGrid,
		PuzzleID: puzzle.ID,
		SolvedAt: time.Now(),
	}

	// Validate the solution before returning
	if err := solution.Validate(); err != nil {
		return nil, errors.New("generated solution is invalid: " + err.Error())
	}

	return solution, nil
}

// Validate validates a Sudoku puzzle for solving (requirement 4.6, 8.6)
func (s *solverService) Validate(puzzle *models.SudokuPuzzle) error {
	if puzzle == nil {
		return errors.New("puzzle cannot be nil")
	}

	// For solver validation, we only check grid validity, not difficulty constraints
	// since users can submit any puzzle to be solved
	if err := puzzle.Grid.ValidateGrid(); err != nil {
		return err
	}

	if puzzle.ID == "" {
		return errors.New("puzzle ID cannot be empty")
	}

	return nil
}

// IsValidMove checks if placing a number at a specific position is valid
func (s *solverService) IsValidMove(grid *models.Grid, row, col, num int) bool {
	if grid == nil {
		return false
	}

	// Check bounds
	if row < 0 || row >= models.GridSize || col < 0 || col >= models.GridSize {
		return false
	}

	// Check number range
	if num < 1 || num > 9 {
		return false
	}

	// Check if cell is already filled
	if grid[row][col] != models.EmptyCell {
		return false
	}

	// Check row constraint (requirement 8.1)
	for j := 0; j < models.GridSize; j++ {
		if grid[row][j] == num {
			return false
		}
	}

	// Check column constraint (requirement 8.2)
	for i := 0; i < models.GridSize; i++ {
		if grid[i][col] == num {
			return false
		}
	}

	// Check 3x3 box constraint (requirement 8.3)
	boxRow := (row / models.BoxSize) * models.BoxSize
	boxCol := (col / models.BoxSize) * models.BoxSize

	for i := boxRow; i < boxRow+models.BoxSize; i++ {
		for j := boxCol; j < boxCol+models.BoxSize; j++ {
			if grid[i][j] == num {
				return false
			}
		}
	}

	return true
}

// solveWithContext implements the backtracking algorithm with context for timeout handling
func (s *solverService) solveWithContext(ctx context.Context, grid *models.Grid) bool {
	// Check for timeout
	select {
	case <-ctx.Done():
		return false
	default:
	}

	// Find the next empty cell
	row, col := s.findEmptyCell(grid)
	if row == -1 {
		// No empty cells found, puzzle is solved
		return true
	}

	// Try numbers 1-9 in the empty cell
	for num := 1; num <= 9; num++ {
		if s.IsValidMove(grid, row, col, num) {
			// Place the number
			grid[row][col] = num

			// Recursively solve the rest
			if s.solveWithContext(ctx, grid) {
				return true
			}

			// Backtrack: remove the number if it doesn't lead to a solution
			grid[row][col] = models.EmptyCell
		}
	}

	// No valid number found for this cell
	return false
}

// findEmptyCell finds the next empty cell in the grid
func (s *solverService) findEmptyCell(grid *models.Grid) (int, int) {
	for i := 0; i < models.GridSize; i++ {
		for j := 0; j < models.GridSize; j++ {
			if grid[i][j] == models.EmptyCell {
				return i, j
			}
		}
	}
	return -1, -1 // No empty cell found
}

// ValidateForSolving performs comprehensive validation for solver input (requirement 4.6)
func (s *solverService) ValidateForSolving(request *models.SolverRequest) error {
	if request == nil {
		return errors.New("request cannot be nil")
	}

	// Validate the request structure
	if err := request.Validate(); err != nil {
		return err
	}

	// Additional validation specific to solving
	grid := models.Grid(request.Grid)
	
	// Check if the puzzle has at least some empty cells to solve
	if grid.IsComplete() {
		return errors.New("puzzle is already complete")
	}

	// Check if puzzle has reasonable number of clues (at least 17 is mathematically required)
	filledCells := grid.CountFilledCells()
	if filledCells < 17 {
		return errors.New("puzzle has insufficient clues (minimum 17 required)")
	}

	return nil
}
// SolveFromRequest is a helper function to solve from a SolverRequest
func SolveFromRequest(service SolverService, request *models.SolverRequest) (*models.SudokuSolution, error) {
	if request == nil {
		return nil, errors.New("request cannot be nil")
	}

	// Validate the request with enhanced validation
	if err := service.ValidateForSolving(request); err != nil {
		return nil, err
	}

	// Create a temporary puzzle from the request
	puzzle := &models.SudokuPuzzle{
		Grid:      models.Grid(request.Grid),
		ID:        uuid.New().String(), // Generate temporary ID
		CreatedAt: time.Now(),
	}

	// Solve the puzzle
	return service.Solve(puzzle)
}