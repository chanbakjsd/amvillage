<script lang="ts">
	import { connected, state } from "./lib/amvillage"
	import { status } from "./lib/state"
	import Login from "./views/Login.svelte"
	import MainMenu from "./views/MainMenu.svelte"
	import Notice from "./views/Notice.svelte"
	import Trade from "./views/Trade.svelte"

	const views = {
		login: Login,
		mainMenu: MainMenu,
		trade: Trade,
	}
</script>

<div class="popup" class:show={$state.popup && $state.notice === ""}><div>{$state.popup}</div></div>
<div class="main-container">
	{#if $state.notice === ""}
		<svelte:component this={views[$status.status]} />
	{:else}
		<Notice />
	{/if}
</div>

{#if !$connected && $status.status !== "login"}
	<div class="warning">⚠️</div>
{/if}

<style lang="postcss">
	.popup {
		@apply w-full bg-red-300 max-h-0 overflow-hidden transition-all;
	}
	.popup.show {
		@apply max-h-full;
	}
	.popup > div {
		@apply p-2;
	}
	.main-container {
		@apply relative block flex-grow;
	}
	.warning {
		@apply absolute bottom-2 right-2 text-5xl animate-bounce;
	}
</style>
