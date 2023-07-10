<script lang="ts">
	export let value: number
	export let max: number
	export let isMin = false

	$: displayMax = Math.max(max, 1)
</script>

<div class="bar" class:isMin>
	<div class="fill" style="width:{(value / displayMax) * 100}%" />
</div>

<style lang="postcss">
	.bar {
		@apply w-40 h-5 bg-black rounded-lg overflow-hidden;
	}
	.fill {
		@apply h-full bg-gray-400;
	}
	.bar.isMin {
		@apply border-2 border-yellow-300;
	}
	.bar.isMin .fill {
		@apply bg-yellow-300;
	}
</style>
