package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"sudoku/internal/models"
	"sudoku/internal/services"
)

func TestSolverHandler_SolvePuzzle(t *testing.T) {
	// Set up test environment
	gin.SetMode(gin.TestMode)
	
	// Create services
	solverService := services.NewSolverService()
	handler := NewSolverHandler(solverService)
	
	// Set up router
	router := gin.New()
	router.POST("/sudoku-solver", handler.SolvePuzzle)

	// Test valid puzzle
	validRequest := models.SolverRequest{
		Grid: [9][9]int{
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
	}

	requestBody, _ := json.Marshal(validRequest)
	req, _ := http.NewRequest("POST", "/sudoku-solver", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response models.SolverResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !response.Success {
		t.Error("Expected successful response")
	}

	if response.Solution == nil {
		t.Error("Expected solution in response")
	}
}

func TestSolverHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	solverService := services.NewSolverService()
	handler := NewSolverHandler(solverService)
	
	router := gin.New()
	router.POST("/sudoku-solver", handler.SolvePuzzle)

	// Test invalid JSON
	req, _ := http.NewRequest("POST", "/sudoku-solver", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestSolverHandler_ValidationErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	solverService := services.NewSolverService()
	handler := NewSolverHandler(solverService)
	
	router := gin.New()
	router.POST("/sudoku-solver", handler.SolvePuzzle)

	// Test puzzle with insufficient clues
	insufficientRequest := models.SolverRequest{
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

	requestBody, _ := json.Marshal(insufficientRequest)
	req, _ := http.NewRequest("POST", "/sudoku-solver", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	if response.Error != "Validation Error" {
		t.Errorf("Expected 'Validation Error', got %s", response.Error)
	}
}