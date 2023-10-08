<script lang="ts">
	import { onMount } from "svelte"
	import { _ } from "svelte-i18n"
	import { fly } from "svelte/transition"
	import Button from "../components/Button.svelte"
	import { state, ws } from "../lib/amvillage"

	$: isAdmin = $state.config.teams[$state.team].is_admin

	onMount(() => {
		if (window.navigator && window.navigator.vibrate) {
			window.navigator.vibrate([200, 100, 200])
		}
	})

	const cancel = () => {
		$ws.send("notice ")
	}
</script>

<main transition:fly={{ y: 500 }}>
	<div class="content">
		<div>
			{$state.notice}
		</div>
	</div>
	{#if isAdmin}
		<Button on:click={cancel} classes="w-full rounded-full">{$_("notice.button.stop")}</Button>
	{/if}
</main>

<style lang="postcss">
	main {
		@apply p-4 flex flex-col gap-8 h-full;
	}
	.content {
		@apply flex flex-grow items-center justify-center text-3xl break-words;
	}
</style>
