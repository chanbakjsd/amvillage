<script lang="ts">
	export let min = 0
	export let max: number | undefined
	export let value = min

	let timeout: number
	const cap = (num: number): number => {
		num = Math.max(min, num)
		if (typeof max === "number") num = Math.min(max, num)
		return num
	}
	const hold = (num: number) => () => {
		let speed = 0.1
		let internalValue = value - 0
		timeout = setInterval(() => {
			const prev = internalValue
			internalValue += speed * num
			internalValue = cap(internalValue)
			if (internalValue === prev) {
				clearInterval(timeout)
				return
			}
			value = Math.ceil(internalValue)
			speed *= 1.1
		}, 50)
	}
	const release = () => {
		clearInterval(timeout)
	}

	const update = () => {
		value = cap(value)
	}
</script>

<div>
	<button on:mousedown={hold(-1)} on:mouseup={release} disabled={value <= min}>−</button>
	<input type="number" inputmode="numeric" {min} {max} bind:value on:change={update} />
	<button on:mousedown={hold(1)} on:mouseup={release} disabled={value >= max}>+</button>
</div>

<style lang="postcss">
	div {
		@apply flex;
	}
	button {
		@apply inline-flex w-8 h-8 bg-black text-white items-center justify-center;
	}
	button:disabled {
		@apply bg-gray-300 text-gray-700;
	}
	input {
		@apply border border-black p-1 w-16 h-8 text-center;
	}

	/* Reset. */
	input::-webkit-outer-spin-button,
	input::-webkit-inner-spin-button {
		@apply m-0;
		-webkit-appearance: none;
	}

	input[type="number"] {
		-moz-appearance: textfield;
	}
</style>
