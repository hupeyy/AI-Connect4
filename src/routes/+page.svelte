<script lang="ts">
    import { onMount } from 'svelte';
    import Board from '$lib/components/Board.svelte';

</script>

<div class="align-center text-center">
    <h1>Connect 4 AI</h1>
    <Board>
        <slot></slot>
    </Board>
</div>