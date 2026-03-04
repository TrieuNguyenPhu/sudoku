// Sudoku Game JavaScript
class SudokuGame {
    constructor() {
        this.grid = Array(9).fill().map(() => Array(9).fill(0));
        this.originalGrid = Array(9).fill().map(() => Array(9).fill(0));
        this.init();
    }

    init() {
        this.bindEvents();
        this.updateStatusMessage('Select a difficulty to start playing!');
    }

    bindEvents() {
        // Difficulty buttons
        document.getElementById('easy-btn').addEventListener('click', () => this.generatePuzzle('easy'));
        document.getElementById('medium-btn').addEventListener('click', () => this.generatePuzzle('medium'));
        document.getElementById('hard-btn').addEventListener('click', () => this.generatePuzzle('hard'));
        
        // Action buttons
        document.getElementById('solve-btn').addEventListener('click', () => this.solvePuzzle());
        
        // Cell input handling
        document.querySelectorAll('.sudoku-cell').forEach(cell => {
            cell.addEventListener('input', (e) => this.handleCellInput(e));
            cell.addEventListener('keydown', (e) => this.handleKeyDown(e));
        });
    }

    async generatePuzzle(difficulty) {
        try {
            this.updateStatusMessage('Generating puzzle...');
            const response = await fetch(`/api/${difficulty}`);
            
            if (!response.ok) {
                throw new Error(`Failed to generate ${difficulty} puzzle`);
            }
            
            const data = await response.json();
            
            // For now, just show the API response since puzzle generation isn't implemented yet
            this.updateStatusMessage(`${difficulty.charAt(0).toUpperCase() + difficulty.slice(1)} puzzle loaded! (API: ${data.message})`);
            
            // TODO: Load actual puzzle data when puzzle generation is implemented
            // this.loadPuzzle(data.grid);
            
        } catch (error) {
            this.updateStatusMessage(`Error: ${error.message}`, 'error');
        }
    }

    async solvePuzzle() {
        try {
            this.updateStatusMessage('Solving puzzle...');
            
            const response = await fetch('/api/sudoku-solver', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ grid: this.grid })
            });
            
            if (!response.ok) {
                throw new Error('Failed to solve puzzle');
            }
            
            const data = await response.json();
            
            // For now, just show the API response since solver isn't implemented yet
            this.updateStatusMessage(`Solver called! (API: ${data.message})`);
            
            // TODO: Load solution when solver is implemented
            // this.loadSolution(data.grid);
            
        } catch (error) {
            this.updateStatusMessage(`Error: ${error.message}`, 'error');
        }
    }

    handleCellInput(event) {
        const cell = event.target;
        const value = cell.value;
        
        // Only allow numbers 1-9 or empty
        if (value && (!/^[1-9]$/.test(value))) {
            cell.value = '';
            return;
        }
        
        const row = parseInt(cell.dataset.row);
        const col = parseInt(cell.dataset.col);
        
        this.grid[row][col] = value ? parseInt(value) : 0;
        
        // Basic validation (will be enhanced when validation logic is implemented)
        this.validateCell(cell, row, col);
    }

    handleKeyDown(event) {
        // Allow navigation with arrow keys
        if (['ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight'].includes(event.key)) {
            event.preventDefault();
            this.navigateGrid(event.target, event.key);
        }
    }

    navigateGrid(currentCell, direction) {
        const row = parseInt(currentCell.dataset.row);
        const col = parseInt(currentCell.dataset.col);
        let newRow = row;
        let newCol = col;
        
        switch (direction) {
            case 'ArrowUp':
                newRow = Math.max(0, row - 1);
                break;
            case 'ArrowDown':
                newRow = Math.min(8, row + 1);
                break;
            case 'ArrowLeft':
                newCol = Math.max(0, col - 1);
                break;
            case 'ArrowRight':
                newCol = Math.min(8, col + 1);
                break;
        }
        
        const nextCell = document.querySelector(`[data-row="${newRow}"][data-col="${newCol}"]`);
        if (nextCell) {
            nextCell.focus();
        }
    }

    validateCell(cell, row, col) {
        // Basic validation - will be enhanced when full validation is implemented
        const value = this.grid[row][col];
        
        if (value === 0) {
            cell.classList.remove('invalid', 'valid');
            return;
        }
        
        // Simple duplicate check in row and column
        let isValid = true;
        
        // Check row
        for (let c = 0; c < 9; c++) {
            if (c !== col && this.grid[row][c] === value) {
                isValid = false;
                break;
            }
        }
        
        // Check column
        if (isValid) {
            for (let r = 0; r < 9; r++) {
                if (r !== row && this.grid[r][col] === value) {
                    isValid = false;
                    break;
                }
            }
        }
        
        cell.classList.toggle('invalid', !isValid);
        cell.classList.toggle('valid', isValid);
    }

    loadPuzzle(puzzleGrid) {
        // Will be implemented when puzzle generation is ready
        this.grid = puzzleGrid.map(row => [...row]);
        this.originalGrid = puzzleGrid.map(row => [...row]);
        this.renderGrid();
    }

    loadSolution(solutionGrid) {
        // Will be implemented when solver is ready
        this.grid = solutionGrid.map(row => [...row]);
        this.renderGrid();
        this.updateStatusMessage('Puzzle solved!', 'success');
    }

    renderGrid() {
        document.querySelectorAll('.sudoku-cell').forEach(cell => {
            const row = parseInt(cell.dataset.row);
            const col = parseInt(cell.dataset.col);
            const value = this.grid[row][col];
            
            cell.value = value === 0 ? '' : value.toString();
            
            // Mark pre-filled cells
            if (this.originalGrid[row][col] !== 0) {
                cell.classList.add('pre-filled');
                cell.readOnly = true;
            } else {
                cell.classList.remove('pre-filled');
                cell.readOnly = false;
            }
        });
    }

    updateStatusMessage(message, type = 'info') {
        const statusElement = document.getElementById('status-message');
        const statusContainer = statusElement.parentElement;
        
        statusElement.textContent = message;
        
        // Remove existing status classes
        statusContainer.classList.remove('success', 'error', 'info');
        statusContainer.classList.add(type);
    }
}

// Initialize the game when the page loads
document.addEventListener('DOMContentLoaded', () => {
    new SudokuGame();
});