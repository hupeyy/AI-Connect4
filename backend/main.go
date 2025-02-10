package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"math"
)

// GameState represents the board state and current player
type GameState struct {
	Board  [][]int	`json:"board"`
	Player int 		`json:"player"`
	Column int      `json:"column"`
}

func main() {
	http.HandleFunc("/check-win", checkWinHandler)
	http.HandleFunc("/user-move", userMoveHandler)
	http.HandleFunc("/computer-move", computerMoveHandler)
	http.HandleFunc("/valid-moves", validMovesHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			setHeaders(w)
			w.WriteHeader(http.StatusOK)
			return
		}

		http.NotFound(w, r)
	})

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func checkWinHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var state GameState
	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	winner := check_win(state.Board)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"winner": winner,
	})
}

func userMoveHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var state GameState
	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	newBoard := make_move(state.Board, state.Player, state.Column)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"board":   newBoard,
		"message": "Move executed at column " + fmt.Sprint(state.Column),
	})
}

func computerMoveHandler(w http.ResponseWriter, r *http.Request) {
    setHeaders(w)
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

    var state GameState
    err := json.NewDecoder(r.Body).Decode(&state)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    _, col := minimax(state.Board, MAX_DEPTH, math.MinInt32, math.MaxInt32, true)
    newBoard := make_move(state.Board, -1, col)

    json.NewEncoder(w).Encode(map[string]interface{}{
        "board":   newBoard,
        "column":  col,
        "message": fmt.Sprintf("Computer moved to column %d", col),
    })
}

func validMovesHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var state GameState
	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validColumns := valid_moves(state.Board)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid_moves": validColumns,
	})
}

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func check_win(board [][]int) int {
    rows := len(board)
    cols := len(board[0])

    // 0 = no winner, 1 = red, -1 = yellow

    // Check rows
    for r := 0; r < rows; r++ {
        for c := 0; c < cols-3; c++ {
            if board[r][c] == 1 && board[r][c+1] == 1 && 
               board[r][c+2] == 1 && board[r][c+3] == 1 {
                return 1
            }
            if board[r][c] == -1 && board[r][c+1] == -1 && 
               board[r][c+2] == -1 && board[r][c+3] == -1 {
                return -1
            }
        }
    }

    // Check columns
    for r := 0; r < rows-3; r++ {
        for c := 0; c < cols; c++ {
            if board[r][c] == 1 && board[r+1][c] == 1 && 
               board[r+2][c] == 1 && board[r+3][c] == 1 {
                return 1
            } 
            if board[r][c] == -1 && board[r+1][c] == -1 && 
               board[r+2][c] == -1 && board[r+3][c] == -1 {
                return -1
            }
        }
    }

    // Check positive diagonals
    for r := 0; r < rows-3; r++ {
        for c := 0; c < cols-3; c++ {
            if board[r][c] == 1 && board[r+1][c+1] == 1 && 
               board[r+2][c+2] == 1 && board[r+3][c+3] == 1 {
                return 1
            }
            if board[r][c] == -1 && board[r+1][c+1] == -1 && 
               board[r+2][c+2] == -1 && board[r+3][c+3] == -1 {
                return -1
            }
        }
    }

    // Check negative diagonals
    for c := 0; c < cols-3; c++ {
		for r := 3; r < rows; r++ {
			if board[r][c] == 1 && board[r-1][c+1] == 1 && 
			board[r-2][c+2] == 1 && board[r-3][c+3] == 1 {
				return 1
			}
			if board[r][c] == -1 && board[r-1][c+1] == -1 && 
			board[r-2][c+2] == -1 && board[r-3][c+3] == -1 {
				return -1
			}
		}
	}

    return 0
}

func valid_moves(board [][]int) []int {

	if board == nil || len(board) == 0 {
		return []int{}
	}
	rows := len(board)
	cols := len(board[0])
	moves := []int{}

	for c := 0; c < cols; c++ {
		if board[rows-1][c] == 0 {
			moves = append(moves, c)
		}
	}

	return moves
}

func make_move(board [][]int , player int, col int) [][]int {
	rows := len(board)

	for r := 0; r < rows; r++ {
		if board[r][col] == 0 {
			board[r][col] = player
			break
		}
	}

	return board
}

func print_board(board [][]int) {
	rows := len(board)
	cols := len(board[0])

	for r := rows-1; r >= 0; r-- {
		for c := 0; c < cols; c++ {
			if board[r][c] == 0 {
				print(".")
			} else {
				print(string(board[r][c]))
			}
		}
		println()
	}
}

func evaluate_board(board [][]int) int {
    score := 0
    rows := len(board)
    cols := len(board[0])

    // Center column preference (columns 3 and 4 in 0-6 index)
    centerCol1 := cols/2 - 1
    centerCol2 := cols/2
    for r := 0; r < rows; r++ {
        if board[r][centerCol1] == 1 {
            score += 2
        } else if board[r][centerCol1] == -1 {
            score -= 2
        }
        if board[r][centerCol2] == 1 {
            score += 2
        } else if board[r][centerCol2] == -1 {
            score -= 2
        }
    }

    // Evaluate all possible lines of 4 cells
    // Horizontal
    for r := 0; r < rows; r++ {
        for c := 0; c < cols-3; c++ {
            window := [4]int{
                board[r][c],
                board[r][c+1],
                board[r][c+2],
                board[r][c+3],
            }
            score += evaluateWindow(window)
        }
    }

    // Vertical
    for c := 0; c < cols; c++ {
        for r := 0; r < rows-3; r++ {
            window := [4]int{
                board[r][c],
                board[r+1][c],
                board[r+2][c],
                board[r+3][c],
            }
            score += evaluateWindow(window)
        }
    }

    // Positive slope diagonal
    for r := 0; r < rows-3; r++ {
        for c := 0; c < cols-3; c++ {
            window := [4]int{
                board[r][c],
                board[r+1][c+1],
                board[r+2][c+2],
                board[r+3][c+3],
            }
            score += evaluateWindow(window)
        }
    }

    // Negative slope diagonal
    for r := 3; r < rows; r++ {
        for c := 0; c < cols-3; c++ {
            window := [4]int{
                board[r][c],
                board[r-1][c+1],
                board[r-2][c+2],
                board[r-3][c+3],
            }
            score += evaluateWindow(window)
        }
    }

    return score
}

func evaluateWindow(window [4]int) int {
    score := 0
    redCount := 0
    yellowCount := 0

    for _, cell := range window {
        switch cell {
        case 1:
            redCount++
        case -1:
            yellowCount++
        }
    }

    // Evaluate red potential
    if redCount > 0 && yellowCount == 0 {
        switch redCount {
        case 3:
            score += 50
        case 2:
            score += 10
        case 1:
            score += 1
        }
    }

    // Evaluate yellow potential
    if yellowCount > 0 && redCount == 0 {
        switch yellowCount {
        case 3:
            score -= 50
        case 2:
            score -= 10
        case 1:
            score -= 1
        }
    }

    return score
}

const (
    MAX_DEPTH = 5
    WIN_SCORE = 100000
)

func minimax(board [][]int, depth int, alpha int, beta int, maximizingPlayer bool) (int, int) {
    validMoves := valid_moves(board)
    winner := check_win(board)
    
    // Terminal conditions
    if winner != 0 || depth == 0 || len(validMoves) == 0 {
        if winner != 0 {
            // Multiply by depth to prioritize faster wins
            return winner * (depth + 1) * WIN_SCORE, -1
        }
        if depth == 0 {
            return evaluate_board(board), -1
        }
        return 0, -1 // Draw
    }

    bestCol := validMoves[0]
    var bestVal int

    if maximizingPlayer {
        bestVal = math.MinInt32
        for _, col := range validMoves {
            // Create new board state for each move
            newBoard := makeMoveCopy(board, col, 1)
            val, _ := minimax(newBoard, depth-1, alpha, beta, false)
            
            if val > bestVal {
                bestVal = val
                bestCol = col
            }
            alpha = max(alpha, bestVal)
            if beta <= alpha {
                break // Beta cutoff
            }
        }
    } else {
        bestVal = math.MaxInt32
        for _, col := range validMoves {
            newBoard := makeMoveCopy(board, col, -1)
            val, _ := minimax(newBoard, depth-1, alpha, beta, true)
            
            if val < bestVal {
                bestVal = val
                bestCol = col
            }
            beta = min(beta, bestVal)
            if beta <= alpha {
                break // Alpha cutoff
            }
        }
    }
    
    return bestVal, bestCol
}

func makeMoveCopy(board [][]int, col int, player int) [][]int {
    // Create deep copy
    newBoard := make([][]int, len(board))
    for i := range board {
        newBoard[i] = make([]int, len(board[i]))
        copy(newBoard[i], board[i])
    }

    // Find first empty row in column
    for r := 0; r < len(newBoard); r++ {
        if newBoard[r][col] == 0 {
            newBoard[r][col] = player
            break
        }
    }
    return newBoard
}
