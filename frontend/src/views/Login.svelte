<script lang="ts">
	import { fly } from "svelte/transition"
	import Button from "../components/Button.svelte"
	import { error, state, ws } from "../lib/amvillage"
	import { status } from "../lib/state"

	let username = ""
	const password = location.pathname.split("/").at(-1)
	const logon = () => {
		$error = ""
		$ws.send(`login ${btoa(password)} ${username}`)
	}

	$: {
		if ($state.username) {
			$status = {
				status: "mainMenu"
			}
		}
	}

</script>

<main transition:fly={{ y: 500 }}>
	<div>
		<label for="username">显示名：</label>
		<input id="username" type="text" autocomplete="off" bind:value={username} placeholder="你的名字" />
	</div>
	<Button on:click={logon} disabled={!username}>开始</Button>
		<div class="error">{$error}&nbsp;</div>
</main>

<style lang="postcss">
	main {
		@apply flex flex-col items-center justify-center gap-4 w-full h-full;
	}
	div {
		@apply flex gap-2 items-center justify-center;
	}
	input {
		@apply border border-black p-2;
	}
	.error {
		@apply text-red-600 font-semibold;
	}
</style>
