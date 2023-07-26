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

{#if $state.notice === ""}
	<svelte:component this={views[$status.status]} />
{:else}
	<Notice />
{/if}

{#if !$connected && $status.status !== "login"}
	<div class="warning">⚠️</div>
{/if}

<style lang="postcss">
	.warning {
		@apply absolute bottom-2 right-2 text-5xl animate-bounce;
	}
</style>
