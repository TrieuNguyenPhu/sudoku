package models

import (
	"errors"
	"time"
)

// Grid represents a 9x9 Sudoku board where 0 indicates empty cells
type Grid [9][9]int

// Difficulty represents the difficulty level of a Sudoku puzzle
type Difficulty string

const (
	Easy   Difficulty = "easy"   // 40-45 pre-filled cells
	Medium Difficulty = "medium" // 30-35 pre-filled cells
	Hard   Difficulty = "hard"   // 20-25 pre-filled cells
)

// Validation constants
const (
	GridSize    = 9
	BoxSize     = 3
	MinValue    = 0 // 0 represents empty cell
	MaxValue    = 9
	EmptyCell   = 0
	
	// Pre-filled cell counts for each difficulty
	EasyMinCells   = 40
	EasyMaxCells   = 45
	MediumMinCells = 30
	MediumMaxCells = 35
	HardMinCells   = 20
	HardMaxCells   = 25
)

// SudokuPuzzle represents a Sudoku puzzle with metadata
type SudokuPuzzle struct {
	Grid       Grid       `json:"grid"`
	Difficulty Difficulty `json:"difficulty"`
	ID         string     `json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
}

// SudokuSolution represents a solved Sudoku puzzle
type SudokuSolution struct {
	Grid     Grid      `json:"grid"`
	PuzzleID string    `json:"puzzle_id"`
	SolvedAt time.Time `json:"solved_at"`
}

// ValidateGrid validates that the grid follows basic Sudoku rules
func (g *Grid) ValidateGrid() error {
	// Check grid dimensions (already guaranteed by type, but good to be explicit)
	if len(g) != GridSize {
		return errors.New("grid must be 9x9")
	}
	
	for i := 0; i < GridSize; i++ {
		if len(g[i]) != GridSize {
			return errors.New("grid must be 9x9")
		}
		
		// Check cell values are in valid range (0-9)
		for j := 0; j < GridSize; j++ {
			if g[i][j] < MinValue || g[i][j] > MaxValue {
				return errors.New("cell values must be between 0 and 9")
			}
		}
	}
	
	// Validate Sudoku constraints
	if err := g.validateRows(); err != nil {
		return err
	}
	
	if err := g.validateColumns(); err != nil {
		return err
	}
	
	if err := g.validateBoxes(); err != nil {
		return err
	}
	
	return nil
}

// validateRows validates that each row contains unique numbers 1-9 (requirement 8.1)
func (g *Grid) validateRows() error {
	for i := 0; i < GridSize; i++ {
		seen := make(map[int]bool)
		for j := 0; j < GridSize; j++ {
			value := g[i][j]
			if value != EmptyCell {
				if seen[value] {
					return errors.New("duplicate number found in row")
				}
				seen[value] = true
			}
		}
	}
	return nil
}

// validateColumns validates that each column contains unique numbers 1-9 (requirement 8.2)
func (g *Grid) validateColumns() error {
	for j := 0; j < GridSize; j++ {
		seen := make(map[int]bool)
		for i := 0; i < GridSize; i++ {
			value := g[i][j]
			if value != EmptyCell {
				if seen[value] {
					return errors.New("duplicate number found in column")
				}
				seen[value] = true
			}
		}
	}
	return nil
}

// validateBoxes validates that each 3x3 box contains unique numbers 1-9 (requirement 8.3)
func (g *Grid) validateBoxes() error {
	for boxRow := 0; boxRow < BoxSize; boxRow++ {
		for boxCol := 0; boxCol < BoxSize; boxCol++ {
			seen := make(map[int]bool)
			
			// Check each cell in the 3x3 box
			for i := boxRow * BoxSize; i < (boxRow+1)*BoxSize; i++ {
				for j := boxCol * BoxSize; j < (boxCol+1)*BoxSize; j++ {
					value := g[i][j]
					if value != EmptyCell {
						if seen[value] {
							return errors.New("duplicate number found in 3x3 box")
						}
						seen[value] = true
					}
				}
			}
		}
	}
	return nil
}

// CountFilledCells returns the number of pre-filled cells in the grid
func (g *Grid) CountFilledCells() int {
	count := 0
	for i := 0; i < GridSize; i++ {
		for j := 0; j < GridSize; j++ {
			if g[i][j] != EmptyCell {
				count++
			}
		}
	}
	return count
}

// IsComplete checks if the grid is completely filled with valid numbers
func (g *Grid) IsComplete() bool {
	for i := 0; i < GridSize; i++ {
		for j := 0; j < GridSize; j++ {
			if g[i][j] == EmptyCell {
				return false
			}
		}
	}
	return true
}

// ValidateDifficulty validates that the puzzle has the correct number of pre-filled cells for its difficulty
func (p *SudokuPuzzle) ValidateDifficulty() error {
	filledCells := p.Grid.CountFilledCells()
	
	switch p.Difficulty {
	case Easy:
		if filledCells < EasyMinCells || filledCells > EasyMaxCells {
			return errors.New("easy puzzle must have 40-45 pre-filled cells")
		}
	case Medium:
		if filledCells < MediumMinCells || filledCells > MediumMaxCells {
			return errors.New("medium puzzle must have 30-35 pre-filled cells")
		}
	case Hard:
		if filledCells < HardMinCells || filledCells > HardMaxCells {
			return errors.New("hard puzzle must have 20-25 pre-filled cells")
		}
	default:
		return errors.New("invalid difficulty level")
	}
	
	return nil
}

// Validate validates the entire puzzle structure and constraints
func (p *SudokuPuzzle) Validate() error {
	if err := p.Grid.ValidateGrid(); err != nil {
		return err
	}
	
	if err := p.ValidateDifficulty(); err != nil {
		return err
	}
	
	if p.ID == "" {
		return errors.New("puzzle ID cannot be empty")
	}
	
	return nil
}

// Validate validates the solution structure
func (s *SudokuSolution) Validate() error {
	if err := s.Grid.ValidateGrid(); err != nil {
		return err
	}
	
	if !s.Grid.IsComplete() {
		return errors.New("solution must be complete (no empty cells)")
	}
	
	if s.PuzzleID == "" {
		return errors.New("puzzle ID cannot be empty")
	}
	
	return nil
}