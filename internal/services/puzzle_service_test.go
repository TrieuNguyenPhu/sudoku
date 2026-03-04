package services

import (
	"testing"
	"time"

	"sudoku/internal/models"
)

func TestPuzzleService_GenerateEasy(t *testing.T) {
	solverService := NewSolverService()
	puzzleService := NewPuzzleService(solverService)

	puzzle, err := puzzleService.GenerateEasy()
	if err != nil {
		t.Fatalf("Failed to generate easy puzzle: %v", err)
	}

	// Validate puzzle structure
	if err := puzzle.Validate(); err != nil {
		t.Errorf("Generated puzzle failed validation: %v", err)
	}

	// Check difficulty constraints
	filledCells := puzzle.Grid.CountFilledCells()
	if filledCells < models.EasyMinCells || filledCells > models.EasyMaxCells {
		t.Errorf("Easy puzzle has %d filled cells, expected %d-%d", filledCells, models.EasyMinCells, models.EasyMaxCells)
	}

	// Check that puzzle has correct difficulty
	if puzzle.Difficulty != models.Easy {
		t.Errorf("Expected difficulty %s, got %s", models.Easy, puzzle.Difficulty)
	}

	// Verify puzzle can be solved
	solution, err := solverService.Solve(puzzle)
	if err != nil {
		t.Errorf("Generated puzzle cannot be solved: %v", err)
	}

	// Verify solution is complete
	if !solution.Grid.IsComplete() {
		t.Error("Solution is not complete")
	}
}

func TestPuzzleService_GenerateMedium(t *testing.T) {
	solverService := NewSolverService()
	puzzleService := NewPuzzleService(solverService)

	puzzle, err := puzzleService.GenerateMedium()
	if err != nil {
		t.Fatalf("Failed to generate medium puzzle: %v", err)
	}

	// Validate puzzle structure
	if err := puzzle.Validate(); err != nil {
		t.Errorf("Generated puzzle failed validation: %v", err)
	}

	// Check difficulty constraints
	filledCells := puzzle.Grid.CountFilledCells()
	if filledCells < models.MediumMinCells || filledCells > models.MediumMaxCells {
		t.Errorf("Medium puzzle has %d filled cells, expected %d-%d", filledCells, models.MediumMinCells, models.MediumMaxCells)
	}

	// Check that puzzle has correct difficulty
	if puzzle.Difficulty != models.Medium {
		t.Errorf("Expected difficulty %s, got %s", models.Medium, puzzle.Difficulty)
	}
}

func TestPuzzleService_GenerateHard(t *testing.T) {
	solverService := NewSolverService()
	puzzleService := NewPuzzleService(solverService)

	puzzle, err := puzzleService.GenerateHard()
	if err != nil {
		t.Fatalf("Failed to generate hard puzzle: %v", err)
	}

	// Validate puzzle structure
	if err := puzzle.Validate(); err != nil {
		t.Errorf("Generated puzzle failed validation: %v", err)
	}

	// Check difficulty constraints
	filledCells := puzzle.Grid.CountFilledCells()
	if filledCells < models.HardMinCells || filledCells > models.HardMaxCells {
		t.Errorf("Hard puzzle has %d filled cells, expected %d-%d", filledCells, models.HardMinCells, models.HardMaxCells)
	}

	// Check that puzzle has correct difficulty
	if puzzle.Difficulty != models.Hard {
		t.Errorf("Expected difficulty %s, got %s", models.Hard, puzzle.Difficulty)
	}
}

func TestPuzzleService_GenerateWithTimeout(t *testing.T) {
	solverService := NewSolverService()
	puzzleService := NewPuzzleService(solverService)

	// Test with reasonable timeout that should succeed
	puzzle, err := puzzleService.GenerateWithTimeout(models.Easy, 2000*time.Millisecond)
	if err != nil {
		t.Errorf("Expected successful generation with reasonable timeout, got error: %v", err)
	}
	if puzzle == nil {
		t.Error("Expected puzzle to be generated")
	}

	// Test timeout behavior by checking that timeout is respected
	start := time.Now()
	_, err = puzzleService.GenerateWithTimeout(models.Hard, 50*time.Millisecond)
	elapsed := time.Since(start)
	
	// Should either succeed quickly or timeout within reasonable time
	if err != nil && err.Error() == "puzzle generation timeout exceeded" {
		// This is expected behavior for short timeout
		if elapsed > 200*time.Millisecond {
			t.Errorf("Timeout took too long: %v", elapsed)
		}
	} else if err == nil {
		// Generation succeeded quickly, which is also valid
		if elapsed > 100*time.Millisecond {
			t.Errorf("Generation took too long for success case: %v", elapsed)
		}
	} else {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestPuzzleService_InvalidDifficulty(t *testing.T) {
	solverService := NewSolverService()
	puzzleService := NewPuzzleService(solverService)

	_, err := puzzleService.GenerateWithDifficulty("invalid")
	if err == nil {
		t.Error("Expected error for invalid difficulty")
	}
	if err.Error() != "invalid difficulty level" {
		t.Errorf("Expected 'invalid difficulty level' error, got: %v", err)
	}
}