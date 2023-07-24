<script lang="ts">
	import { onMount } from "svelte"
	import { connect, connected } from "./lib/amvillage"
	import { status } from "./lib/state"
	import Login from "./views/Login.svelte"
	import MainMenu from "./views/MainMenu.svelte"
	import Trade from "./views/Trade.svelte"

	onMount(() => {
		connect()
	})

	const views = {
		login: Login,
		mainMenu: MainMenu,
		trade: Trade,
	}
</script>

<svelte:component this={views[$status.status]} />

{#if !$connected}
	<div class="warning">⚠️</div>
{/if}

<style lang="postcss">
	.warning {
		@apply absolute bottom-2 right-2 text-5xl animate-bounce;
	}
</style>
