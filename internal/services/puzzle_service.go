package services

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"time"

	"github.com/google/uuid"
	"sudoku/internal/models"
)

// PuzzleService defines the interface for Sudoku puzzle generation operations
type PuzzleService interface {
	GenerateEasy() (*models.SudokuPuzzle, error)
	GenerateMedium() (*models.SudokuPuzzle, error)
	GenerateHard() (*models.SudokuPuzzle, error)
	GenerateWithDifficulty(difficulty models.Difficulty) (*models.SudokuPuzzle, error)
	GenerateWithTimeout(difficulty models.Difficulty, timeout time.Duration) (*models.SudokuPuzzle, error)
}

// puzzleService implements the PuzzleService interface
type puzzleService struct {
	solverService SolverService
}

// NewPuzzleService creates a new instance of PuzzleService
func NewPuzzleService(solverService SolverService) PuzzleService {
	return &puzzleService{
		solverService: solverService,
	}
}

// GenerateEasy generates an easy Sudoku puzzle (40-45 pre-filled cells, 500ms timeout)
func (p *puzzleService) GenerateEasy() (*models.SudokuPuzzle, error) {
	return p.GenerateWithTimeout(models.Easy, 500*time.Millisecond)
}

// GenerateMedium generates a medium Sudoku puzzle (30-35 pre-filled cells, 1000ms timeout)
func (p *puzzleService) GenerateMedium() (*models.SudokuPuzzle, error) {
	return p.GenerateWithTimeout(models.Medium, 1000*time.Millisecond)
}

// GenerateHard generates a hard Sudoku puzzle (20-25 pre-filled cells, 2000ms timeout)
func (p *puzzleService) GenerateHard() (*models.SudokuPuzzle, error) {
	return p.GenerateWithTimeout(models.Hard, 2000*time.Millisecond)
}

// GenerateWithDifficulty generates a puzzle with the specified difficulty using default timeouts
func (p *puzzleService) GenerateWithDifficulty(difficulty models.Difficulty) (*models.SudokuPuzzle, error) {
	var timeout time.Duration
	
	switch difficulty {
	case models.Easy:
		timeout = 500 * time.Millisecond
	case models.Medium:
		timeout = 1000 * time.Millisecond
	case models.Hard:
		timeout = 2000 * time.Millisecond
	default:
		return nil, errors.New("invalid difficulty level")
	}
	
	return p.GenerateWithTimeout(difficulty, timeout)
}

// GenerateWithTimeout generates a puzzle with the specified difficulty and timeout
func (p *puzzleService) GenerateWithTimeout(difficulty models.Difficulty, timeout time.Duration) (*models.SudokuPuzzle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Retry logic for failed generations (requirement 5.4)
	maxRetries := 3
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return nil, errors.New("puzzle generation timeout exceeded")
		default:
		}

		puzzle, err := p.generatePuzzleWithContext(ctx, difficulty)
		if err == nil {
			return puzzle, nil
		}
		
		lastErr = err
		
		// If we're running out of time, don't retry
		if ctx.Err() == context.DeadlineExceeded {
			break
		}
	}

	// All retries failed
	if ctx.Err() == context.DeadlineExceeded {
		return nil, errors.New("puzzle generation timeout exceeded")
	}
	
	return nil, errors.New("failed to generate puzzle after retries: " + lastErr.Error())
}

// generatePuzzleWithContext generates a single puzzle attempt with context for timeout handling
func (p *puzzleService) generatePuzzleWithContext(ctx context.Context, difficulty models.Difficulty) (*models.SudokuPuzzle, error) {
	// Step 1: Generate a complete valid Sudoku solution using backtracking
	completeGrid, err := p.generateCompleteSolution(ctx)
	if err != nil {
		return nil, err
	}

	// Step 2: Remove cells based on difficulty to create the puzzle
	puzzleGrid, err := p.createPuzzleFromSolution(ctx, completeGrid, difficulty)
	if err != nil {
		return nil, err
	}

	// Step 3: Verify the puzzle has exactly one unique solution
	if !p.hasUniqueSolution(ctx, puzzleGrid) {
		return nil, errors.New("generated puzzle does not have unique solution")
	}

	// Step 4: Create the puzzle object
	puzzle := &models.SudokuPuzzle{
		Grid:       *puzzleGrid,
		Difficulty: difficulty,
		ID:         uuid.New().String(),
		CreatedAt:  time.Now(),
	}

	// Step 5: Final validation
	if err := puzzle.Validate(); err != nil {
		return nil, errors.New("generated puzzle failed validation: " + err.Error())
	}

	return puzzle, nil
}

// generateCompleteSolution generates a complete valid Sudoku solution using backtracking with constraint propagation
func (p *puzzleService) generateCompleteSolution(ctx context.Context) (*models.Grid, error) {
	var grid models.Grid
	
	// Initialize empty grid
	for i := 0; i < models.GridSize; i++ {
		for j := 0; j < models.GridSize; j++ {
			grid[i][j] = models.EmptyCell
		}
	}

	// Fill the grid using backtracking with randomization for variety
	if !p.fillGridRandomly(ctx, &grid) {
		return nil, errors.New("failed to generate complete solution")
	}

	return &grid, nil
}

// fillGridRandomly fills the grid using backtracking with randomized number selection
func (p *puzzleService) fillGridRandomly(ctx context.Context, grid *models.Grid) bool {
	// Check for timeout
	select {
	case <-ctx.Done():
		return false
	default:
	}

	// Find the next empty cell
	row, col := p.findEmptyCell(grid)
	if row == -1 {
		// No empty cells found, grid is complete
		return true
	}

	// Create randomized list of numbers 1-9
	numbers := p.getRandomizedNumbers()

	// Try each number in random order
	for _, num := range numbers {
		if p.solverService.IsValidMove(grid, row, col, num) {
			// Place the number
			grid[row][col] = num

			// Recursively fill the rest
			if p.fillGridRandomly(ctx, grid) {
				return true
			}

			// Backtrack: remove the number if it doesn't lead to a solution
			grid[row][col] = models.EmptyCell
		}
	}

	// No valid number found for this cell
	return false
}

// createPuzzleFromSolution creates a puzzle by removing cells from a complete solution
func (p *puzzleService) createPuzzleFromSolution(ctx context.Context, solution *models.Grid, difficulty models.Difficulty) (*models.Grid, error) {
	// Copy the solution to create the puzzle
	puzzle := *solution

	// Determine target cell count based on difficulty
	var minCells, maxCells int
	switch difficulty {
	case models.Easy:
		minCells, maxCells = models.EasyMinCells, models.EasyMaxCells
	case models.Medium:
		minCells, maxCells = models.MediumMinCells, models.MediumMaxCells
	case models.Hard:
		minCells, maxCells = models.HardMinCells, models.HardMaxCells
	default:
		return nil, errors.New("invalid difficulty level")
	}

	// Calculate target number of cells to remove
	targetCells := p.getRandomInRange(minCells, maxCells)
	cellsToRemove := 81 - targetCells

	// Create list of all cell positions
	positions := make([][2]int, 0, 81)
	for i := 0; i < models.GridSize; i++ {
		for j := 0; j < models.GridSize; j++ {
			positions = append(positions, [2]int{i, j})
		}
	}

	// Shuffle positions for random removal
	p.shufflePositions(positions)

	// Remove cells while maintaining unique solution
	removed := 0
	for _, pos := range positions {
		if removed >= cellsToRemove {
			break
		}

		select {
		case <-ctx.Done():
			return nil, errors.New("timeout during cell removal")
		default:
		}

		row, col := pos[0], pos[1]
		originalValue := puzzle[row][col]

		// Temporarily remove the cell
		puzzle[row][col] = models.EmptyCell

		// Check if puzzle still has unique solution
		if p.hasUniqueSolution(ctx, &puzzle) {
			// Keep the cell removed
			removed++
		} else {
			// Restore the cell
			puzzle[row][col] = originalValue
		}
	}

	// Verify we have the right number of cells
	filledCells := puzzle.CountFilledCells()
	if filledCells < minCells || filledCells > maxCells {
		return nil, errors.New("failed to achieve target difficulty cell count")
	}

	return &puzzle, nil
}

// hasUniqueSolution checks if a puzzle has exactly one unique solution
func (p *puzzleService) hasUniqueSolution(ctx context.Context, puzzle *models.Grid) bool {
	// Count the number of solutions using a modified solver
	solutionCount := p.countSolutions(ctx, puzzle, 0, 2) // Stop counting after 2 solutions
	return solutionCount == 1
}

// countSolutions counts the number of possible solutions (up to maxCount)
func (p *puzzleService) countSolutions(ctx context.Context, grid *models.Grid, currentCount, maxCount int) int {
	// Check for timeout or max count reached
	select {
	case <-ctx.Done():
		return currentCount
	default:
	}

	if currentCount >= maxCount {
		return currentCount
	}

	// Find the next empty cell
	row, col := p.findEmptyCell(grid)
	if row == -1 {
		// No empty cells found, this is a complete solution
		return currentCount + 1
	}

	// Try numbers 1-9 in the empty cell
	for num := 1; num <= 9; num++ {
		if p.solverService.IsValidMove(grid, row, col, num) {
			// Place the number
			grid[row][col] = num

			// Recursively count solutions
			currentCount = p.countSolutions(ctx, grid, currentCount, maxCount)
			
			// Early exit if we've found enough solutions
			if currentCount >= maxCount {
				grid[row][col] = models.EmptyCell
				return currentCount
			}

			// Backtrack: remove the number
			grid[row][col] = models.EmptyCell
		}
	}

	return currentCount
}

// Helper functions

// findEmptyCell finds the next empty cell in the grid
func (p *puzzleService) findEmptyCell(grid *models.Grid) (int, int) {
	for i := 0; i < models.GridSize; i++ {
		for j := 0; j < models.GridSize; j++ {
			if grid[i][j] == models.EmptyCell {
				return i, j
			}
		}
	}
	return -1, -1 // No empty cell found
}

// getRandomizedNumbers returns a shuffled slice of numbers 1-9
func (p *puzzleService) getRandomizedNumbers() []int {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	
	// Fisher-Yates shuffle
	for i := len(numbers) - 1; i > 0; i-- {
		j := p.getRandomInt(i + 1)
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
	
	return numbers
}

// shufflePositions shuffles an array of positions
func (p *puzzleService) shufflePositions(positions [][2]int) {
	// Fisher-Yates shuffle
	for i := len(positions) - 1; i > 0; i-- {
		j := p.getRandomInt(i + 1)
		positions[i], positions[j] = positions[j], positions[i]
	}
}

// getRandomInRange returns a random integer between min and max (inclusive)
func (p *puzzleService) getRandomInRange(min, max int) int {
	if min == max {
		return min
	}
	return min + p.getRandomInt(max-min+1)
}

// getRandomInt returns a cryptographically secure random integer between 0 and n-1
func (p *puzzleService) getRandomInt(n int) int {
	if n <= 0 {
		return 0
	}
	
	nBig := big.NewInt(int64(n))
	result, err := rand.Int(rand.Reader, nBig)
	if err != nil {
		// Fallback to a simple method if crypto/rand fails
		return int(time.Now().UnixNano()) % n
	}
	
	return int(result.Int64())
}