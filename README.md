# Sudoku (Go + Gin + SSR)

Website game Sudoku vi?t b?ng Golang, dùng `gin-gonic/gin` và Server-Side Rendering.

## Tính nang hi?n có

- Trang SSR t?i `/`
- API t?o puzzle theo d? khó:
  - `GET /easy`
  - `GET /medium`
  - `GET /hard`
- API gi?i Sudoku:
  - `POST /sudoku-solver`

## Yêu c?u môi tru?ng

- Go `1.26.2` (theo `go.mod`)

## Cài d?t và ch?y

```bash
go mod tidy
go run ./cmd/server
```

M? trình duy?t t?i: `http://localhost:8080`

## API chi ti?t

### 1) L?y puzzle theo d? khó

```http
GET /easy
GET /medium
GET /hard
```

Response m?u:

```json
{
  "difficulty": "easy",
  "board": [
    [5,3,0,0,7,0,0,0,0],
    [6,0,0,1,9,5,0,0,0],
    [0,9,8,0,0,0,0,6,0],
    [8,0,0,0,6,0,0,0,3],
    [4,0,0,8,0,3,0,0,1],
    [7,0,0,0,2,0,0,0,6],
    [0,6,0,0,0,0,2,8,0],
    [0,0,0,4,1,9,0,0,5],
    [0,0,0,0,8,0,0,7,9]
  ]
}
```

### 2) Gi?i Sudoku

```http
POST /sudoku-solver
Content-Type: application/json
```

Body m?u:

```json
{
  "board": [
    [5,3,0,0,7,0,0,0,0],
    [6,0,0,1,9,5,0,0,0],
    [0,9,8,0,0,0,0,6,0],
    [8,0,0,0,6,0,0,0,3],
    [4,0,0,8,0,3,0,0,1],
    [7,0,0,0,2,0,0,0,6],
    [0,6,0,0,0,0,2,8,0],
    [0,0,0,4,1,9,0,0,5],
    [0,0,0,0,8,0,0,7,9]
  ]
}
```

Response thành công:

```json
{
  "solution": [[...9 s?...], ...]
}
```

Response l?i thu?ng g?p:

- `400`: payload không dúng format ho?c board có conflict
- `422`: board không th? gi?i

## C?u trúc thu m?c

```text
cmd/server/main.go          # Entry point ch?y HTTP server
internal/http/router.go     # Route SSR + REST API
internal/sudoku/engine.go   # Generator + solver
templates/index.html        # SSR template
static/css/style.css        # CSS
static/js/app.js            # Frontend logic g?i API
```

## G?i ý phát tri?n ti?p

- Cho phép nh?p s? tr?c ti?p trên grid
- Thêm validate realtime khi ngu?i choi nh?p sai
- Thêm timer, undo, và luu game theo session
- Vi?t unit test cho generator/solver
