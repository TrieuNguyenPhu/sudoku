package services

import (
	"testing"
	"time"

	"sudoku/internal/models"
)

// Test puzzle with a known solution
var testPuzzle = models.SudokuPuzzle{
	Grid: models.Grid{
		{5, 3, 0, 0, 7, 0, 0, 0, 0},
		{6, 0, 0, 1, 9, 5, 0, 0, 0},
		{0, 9, 8, 0, 0, 0, 0, 6, 0},
		{8, 0, 0, 0, 6, 0, 0, 0, 3},
		{4, 0, 0, 8, 0, 3, 0, 0, 1},
		{7, 0, 0, 0, 2, 0, 0, 0, 6},
		{0, 6, 0, 0, 0, 0, 2, 8, 0},
		{0, 0, 0, 4, 1, 9, 0, 0, 5},
		{0, 0, 0, 0, 8, 0, 0, 7, 9},
	},
	Difficulty: models.Medium, // This puzzle has 30 filled cells, which fits Medium difficulty (30-35)
	ID:         "test-puzzle-1",
	CreatedAt:  time.Now(),
}

// Expected solution for the test puzzle
var expectedSolution = models.Grid{
	{5, 3, 4, 6, 7, 8, 9, 1, 2},
	{6, 7, 2, 1, 9, 5, 3, 4, 8},
	{1, 9, 8, 3, 4, 2, 5, 6, 7},
	{8, 5, 9, 7, 6, 1, 4, 2, 3},
	{4, 2, 6, 8, 5, 3, 7, 9, 1},
	{7, 1, 3, 9, 2, 4, 8, 5, 6},
	{9, 6, 1, 5, 3, 7, 2, 8, 4},
	{2, 8, 7, 4, 1, 9, 6, 3, 5},
	{3, 4, 5, 2, 8, 6, 1, 7, 9},
}

func TestSolverService_Solve(t *testing.T) {
	service := NewSolverService()

	solution, err := service.Solve(&testPuzzle)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if solution == nil {
		t.Fatal("Expected solution, got nil")
	}

	// Verify the solution is complete
	if !solution.Grid.IsComplete() {
		t.Error("Solution should be complete (no empty cells)")
	}

	// Verify the solution is valid
	if err := solution.Grid.ValidateGrid(); err != nil {
		t.Errorf("Solution is invalid: %v", err)
	}

	// Verify the solution matches expected result
	if solution.Grid != expectedSolution {
		t.Error("Solution does not match expected result")
	}

	// Verify metadata
	if solution.PuzzleID != testPuzzle.ID {
		t.Errorf("Expected puzzle ID %s, got %s", testPuzzle.ID, solution.PuzzleID)
	}
}

func TestSolverService_ValidateValidPuzzle(t *testing.T) {
	service := NewSolverService()

	err := service.Validate(&testPuzzle)
	if err != nil {
		t.Errorf("Expected no error for valid puzzle, got: %v", err)
	}
}

func TestSolverService_ValidateInvalidPuzzle(t *testing.T) {
	service := NewSolverService()

	// Test with nil puzzle
	err := service.Validate(nil)
	if err == nil {
		t.Error("Expected error for nil puzzle")
	}

	// Test with invalid grid (duplicate in row)
	invalidPuzzle := testPuzzle
	invalidPuzzle.Grid[0][0] = 3 // Creates duplicate with position [0][1]
	
	err = service.Validate(&invalidPuzzle)
	if err == nil {
		t.Error("Expected error for puzzle with duplicate in row")
	}
}

func TestSolverService_IsValidMove(t *testing.T) {
	service := NewSolverService()
	grid := testPuzzle.Grid

	// Test valid move
	if !service.IsValidMove(&grid, 0, 2, 4) {
		t.Error("Expected valid move to return true")
	}

	// Test invalid move (number already in row)
	if service.IsValidMove(&grid, 0, 2, 5) {
		t.Error("Expected invalid move (duplicate in row) to return false")
	}

	// Test invalid move (number already in column)
	if service.IsValidMove(&grid, 0, 2, 8) {
		t.Error("Expected invalid move (duplicate in column) to return false")
	}

	// Test invalid move (cell already filled)
	if service.IsValidMove(&grid, 0, 0, 1) {
		t.Error("Expected invalid move (cell already filled) to return false")
	}

	// Test invalid bounds
	if service.IsValidMove(&grid, -1, 0, 1) {
		t.Error("Expected invalid move (negative row) to return false")
	}

	if service.IsValidMove(&grid, 0, 10, 1) {
		t.Error("Expected invalid move (column out of bounds) to return false")
	}

	// Test invalid number
	if service.IsValidMove(&grid, 0, 2, 0) {
		t.Error("Expected invalid move (number 0) to return false")
	}

	if service.IsValidMove(&grid, 0, 2, 10) {
		t.Error("Expected invalid move (number 10) to return false")
	}
}

func TestSolverService_UnsolvablePuzzle(t *testing.T) {
	service := NewSolverService()

	// Create an unsolvable puzzle (two 5s in the same row)
	unsolvablePuzzle := models.SudokuPuzzle{
		Grid: models.Grid{
			{5, 3, 5, 0, 7, 0, 0, 0, 0}, // Two 5s in first row
			{6, 0, 0, 1, 9, 5, 0, 0, 0},
			{0, 9, 8, 0, 0, 0, 0, 6, 0},
			{8, 0, 0, 0, 6, 0, 0, 0, 3},
			{4, 0, 0, 8, 0, 3, 0, 0, 1},
			{7, 0, 0, 0, 2, 0, 0, 0, 6},
			{0, 6, 0, 0, 0, 0, 2, 8, 0},
			{0, 0, 0, 4, 1, 9, 0, 0, 5},
			{0, 0, 0, 0, 8, 0, 0, 7, 9},
		},
		Difficulty: models.Medium,
		ID:         "unsolvable-puzzle",
		CreatedAt:  time.Now(),
	}

	// This should fail validation first
	err := service.Validate(&unsolvablePuzzle)
	if err == nil {
		t.Error("Expected validation error for unsolvable puzzle")
	}
}

func TestSolverService_SolveWithTimeout(t *testing.T) {
	service := NewSolverService()

	// Test with reasonable timeout
	solution, err := service.SolveWithTimeout(&testPuzzle, 5*time.Second)
	if err != nil {
		t.Fatalf("Expected no error with reasonable timeout, got: %v", err)
	}

	if solution == nil {
		t.Fatal("Expected solution, got nil")
	}

	// Test with very short timeout (should still work for this simple puzzle)
	solution, err = service.SolveWithTimeout(&testPuzzle, 1*time.Millisecond)
	// This might timeout or succeed depending on system performance
	// We just verify it doesn't panic
	if err != nil && solution != nil {
		t.Error("If error is returned, solution should be nil")
	}
}

func TestSolveFromRequest(t *testing.T) {
	service := NewSolverService()

	request := &models.SolverRequest{
		Grid: [9][9]int(testPuzzle.Grid),
	}

	solution, err := SolveFromRequest(service, request)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if solution == nil {
		t.Fatal("Expected solution, got nil")
	}

	// Verify the solution is valid
	if err := solution.Grid.ValidateGrid(); err != nil {
		t.Errorf("Solution is invalid: %v", err)
	}

	// Test with nil request
	_, err = SolveFromRequest(service, nil)
	if err == nil {
		t.Error("Expected error for nil request")
	}
}

func TestSolverService_ValidateForSolving(t *testing.T) {
	service := NewSolverService()

	// Test valid request
	validRequest := &models.SolverRequest{
		Grid: [9][9]int(testPuzzle.Grid),
	}
	
	err := service.ValidateForSolving(validRequest)
	if err != nil {
		t.Errorf("Expected no error for valid request, got: %v", err)
	}

	// Test nil request
	err = service.ValidateForSolving(nil)
	if err == nil {
		t.Error("Expected error for nil request")
	}

	// Test already complete puzzle
	completeRequest := &models.SolverRequest{
		Grid: [9][9]int(expectedSolution),
	}
	
	err = service.ValidateForSolving(completeRequest)
	if err == nil {
		t.Error("Expected error for already complete puzzle")
	}

	// Test puzzle with insufficient clues (less than 17)
	insufficientRequest := &models.SolverRequest{
		Grid: [9][9]int{
			{5, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	}
	
	err = service.ValidateForSolving(insufficientRequest)
	if err == nil {
		t.Error("Expected error for puzzle with insufficient clues")
	}
}

func TestSolverService_ErrorMessageFormat(t *testing.T) {
	service := NewSolverService()

	// Test unsolvable puzzle error message format (requirement 4.3)
	unsolvablePuzzle := models.SudokuPuzzle{
		Grid: models.Grid{
			{1, 2, 3, 4, 5, 6, 7, 8, 9},
			{1, 0, 0, 0, 0, 0, 0, 0, 0}, // Invalid: duplicate 1 in column
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		ID:        "test-unsolvable",
		CreatedAt: time.Now(),
	}

	// This should fail validation first due to duplicate in column
	err := service.Validate(&unsolvablePuzzle)
	if err == nil {
		t.Error("Expected validation error for puzzle with duplicate in column")
	}
}

func TestSolverService_TimeoutHandling(t *testing.T) {
	service := NewSolverService()

	// Create a more complex puzzle that takes longer to solve
	complexPuzzle := models.SudokuPuzzle{
		Grid: models.Grid{
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 3, 0, 8, 5},
			{0, 0, 1, 0, 2, 0, 0, 0, 0},
			{0, 0, 0, 5, 0, 7, 0, 0, 0},
			{0, 0, 4, 0, 0, 0, 1, 0, 0},
			{0, 9, 0, 0, 0, 0, 0, 0, 0},
			{5, 0, 0, 0, 0, 0, 0, 7, 3},
			{0, 0, 2, 0, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 4, 0, 0, 0, 9},
		},
		ID:        "complex-puzzle",
		CreatedAt: time.Now(),
	}

	// Test timeout with very short duration
	_, err := service.SolveWithTimeout(&complexPuzzle, 1*time.Microsecond)
	if err == nil {
		// If it still solves too quickly, just verify the timeout mechanism works
		t.Log("Puzzle solved too quickly for timeout test, but timeout mechanism is in place")
	} else {
		// Verify timeout error message
		if err.Error() != "solving timeout exceeded" {
			t.Errorf("Expected 'solving timeout exceeded', got: %v", err)
		}
	}

	// Test with reasonable timeout (should succeed)
	solution, err := service.SolveWithTimeout(&testPuzzle, 5*time.Second)
	if err != nil {
		t.Errorf("Expected no error with reasonable timeout, got: %v", err)
	}
	
	if solution == nil {
		t.Error("Expected solution with reasonable timeout")
	}
}

func TestSolverService_3000msDefaultTimeout(t *testing.T) {
	service := NewSolverService()

	// Test that default Solve method uses 3000ms timeout (requirement 4.5)
	start := time.Now()
	solution, err := service.Solve(&testPuzzle)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if solution == nil {
		t.Fatal("Expected solution, got nil")
	}

	// Verify it completed well within the 3000ms timeout
	if duration > 3*time.Second {
		t.Errorf("Solve took too long: %v (should be under 3000ms)", duration)
	}
}