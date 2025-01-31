<script>
    import { onMount } from 'svelte';

    onMount(() => {
        console.log('mounted');
    });

    const rows = 6;
    const cols = 7;
</script>

<style lang="postcss">
    .board {
        @apply grid grid-cols-7 border-blue-400 border-2;
    } 

    .cell {
        @apply border-blue-400 border-2 aspect-square flex justify-center items-center;
    }

    .token {
        @apply rounded-full bg-red-500;
    }
</style>

<div class="board">
    {#each Array(rows) as _, row}
        {#each Array(cols) as _, col}
            <div class="cell">
                <div class="token h-5/6 w-5/6"></div> 
            </div>
        {/each}
    {/each}
</div>

