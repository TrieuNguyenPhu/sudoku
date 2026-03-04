const boardEl = document.getElementById("board");
const statusEl = document.getElementById("status");
let currentBoard = Array.from({ length: 9 }, () => Array(9).fill(0));

function renderBoard(board) {
  boardEl.innerHTML = "";
  board.forEach((row) => {
    row.forEach((n) => {
      const cell = document.createElement("div");
      cell.className = "cell";
      cell.textContent = n === 0 ? "" : String(n);
      boardEl.appendChild(cell);
    });
  });
}

async function fetchPuzzle(level) {
  const res = await fetch(`/${level}`);
  if (!res.ok) throw new Error("Không l?y du?c puzzle");
  const data = await res.json();
  currentBoard = data.board;
  renderBoard(currentBoard);
  statusEl.textContent = `Ğang choi m?c ${data.difficulty}`;
}

async function solveCurrent() {
  const res = await fetch("/sudoku-solver", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ board: currentBoard }),
  });

  const data = await res.json();
  if (!res.ok) {
    statusEl.textContent = data.error || "Solver l?i";
    return;
  }

  currentBoard = data.solution;
  renderBoard(currentBoard);
  statusEl.textContent = "Ğã gi?i xong.";
}

document.querySelectorAll("button[data-level]").forEach((btn) => {
  btn.addEventListener("click", () => fetchPuzzle(btn.dataset.level).catch((e) => {
    statusEl.textContent = e.message;
  }));
});

document.getElementById("solve-btn").addEventListener("click", () => {
  solveCurrent().catch((e) => {
    statusEl.textContent = e.message;
  });
});

fetchPuzzle("easy").catch((e) => {
  statusEl.textContent = e.message;
});
