<script lang="ts">
	export let classes = ""
	export let disabled = false
</script>

<button on:click class={classes} {disabled}><slot /></button>

<style lang="postcss">
	button {
		@apply border border-black p-2 rounded-xl cursor-pointer transition-all;
		@apply hover:bg-black hover:text-white;
		@apply disabled:bg-gray-300 disabled:text-gray-700;
	}
</style>
