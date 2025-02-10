<script lang="ts">
    import { onMount } from 'svelte';



    const rows = 6;
    const cols = 7;
    
    let board = Array.from({length: rows}, () => Array(cols).fill(0));
    let player = 1; // 1 for red, -1 for yellow
    let winner = 0;
    let isProcessing = false;


    async function resetGame() {
        board = Array.from({length: rows}, () => Array(cols).fill(0));
        player = 1;
        winner = 0;
    }

    async function userMoveResponse(player: number, col: number) {
        const response = await fetch('http://localhost:8080/user-move', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                board: board,
                player: player,
                column: col
            })
        });

        return await response.json();
    }

    async function computerMoveResponse() {
        const response = await fetch('http://localhost:8080/computer-move', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                board: board
            })
        });

        return await response.json();
    }

    async function validMovesResponse() {
        const response = await fetch('http://localhost:8080/valid-moves', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                board: board
            })
        });

        return await response.json();
    }

    async function checkWinnerResponse() {
        const response = await fetch('http://localhost:8080/check-win', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                board: board
            })
        });

        return await response.json();
    }


    async function userMove(mouseEvent: MouseEvent) {
        if (isProcessing || winner !== 0) return;
        
        const cell = mouseEvent.target as HTMLElement;
        const col = cell.parentElement ? 
            Array.from(cell.parentElement.children).indexOf(cell) : -1;

        // Validate column index
        if (col < 0 || col >= cols) return;

        isProcessing = true;
        
        try {
            const movesResponse = await validMovesResponse();
            const validMoves = movesResponse.valid_moves;

            if (validMoves.includes(col)) {
                // Human move
                const moveResponse = await userMoveResponse(player, col);
                board = moveResponse.board;
                
                // Check for win
                const winCheck = await checkWinnerResponse();
                winner = winCheck.winner;
                
                if (winner) {
                    console.log(`Player ${winner === 1 ? 'Red' : 'Yellow'} wins!`);
                } else {
                    // Switch to AI turn
                    player = -1;
                    await computerMove();
                }
            }
        } catch (error) {
            console.error('Move error:', error);
        } finally {
            isProcessing = false;
        }
    }

    async function computerMove() {
        if (winner !== 0) return;
        
        isProcessing = true;
        try {
            // AI move
            const moveResponse = await computerMoveResponse();
            board = moveResponse.board;
            
            // Check for win
            const winCheck = await checkWinnerResponse();
            winner = winCheck.winner;
            
            if (winner) {
                console.log(`Player ${winner === 1 ? 'Red' : 'Yellow'} wins!`);
            } else {
                // Switch back to human turn
                player = 1;
            }
        } catch (error) {
            console.error('AI move error:', error);
        } finally {
            isProcessing = false;
        }
    }

    onMount(() => {
        window.addEventListener('click', userMove);
    });   

</script>

<style lang="postcss">
    .board {
        @apply grid grid-cols-7 border-blue-400 border-2;
    } 

    .cell {
        @apply border-blue-400 border-2 aspect-square flex justify-center items-center;
    }

    .red-token {
        @apply rounded-full bg-red-500;
    }

    .yellow-token {
        @apply rounded-full bg-yellow-500;
    }

</style>

<div>
    <div class="board">
        {#each Array(rows).fill(null).map((_, i) => rows - 1 - i) as rowIndex}
            {#each Array(cols) as _, colIndex}
                <div class="cell">
                    {#if board[rowIndex][colIndex] === -1}
                        <div class="yellow-token h-5/6 w-5/6"></div>
                    {:else if board[rowIndex][colIndex] === 1}
                        <div class="red-token h-5/6 w-5/6"></div>
                    {/if}
                </div>
            {/each}
        {/each}
    </div>
    <button on:click={resetGame}>Reset</button>
</div>
