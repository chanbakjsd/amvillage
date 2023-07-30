<script lang="ts">
	export let value: number
	export let max: number
	export let isMin = false

	$: displayMax = Math.max(max, 1)
	$: width = value < 0 ? 100 : (value / displayMax) * 100
</script>

<div class="bar" class:isMin class:negative={value < 0}>
	<div class="fill" style="width:{width}%" />
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
	.bar.negative {
		@apply border-red-300;
	}
	.bar.negative .fill {
		@apply border-red-300;
	}
</style>
