<script lang="ts">
	import { fly } from "svelte/transition"
	import { _ } from "svelte-i18n"
	import Button from "../components/Button.svelte"
	import { creds, error, state } from "../lib/amvillage"
	import { status } from "../lib/state"

	let username = ""
	const password = location.pathname.split("/").at(-1)
	const logon = () => {
		$error = ""
		$creds = { username, password }
	}

	$: {
		if ($state.username) {
			$status = {
				status: "mainMenu",
			}
		}
	}
</script>

<main transition:fly={{ y: 500 }}>
	<h1>{$_("login.text.welcome")}</h1>
	<div>
		<label for="username">{$_("login.label.username")}</label>
		<input
			id="username"
			type="text"
			autocomplete="off"
			bind:value={username}
			placeholder={$_("login.placeholder.username")}
		/>
	</div>
	<Button on:click={logon} disabled={!username}>{$_("login.button.start")}</Button>
	<div class="error">{$error}&nbsp;</div>
</main>

<style lang="postcss">
	main {
		@apply flex flex-col items-center justify-center gap-4 w-full h-full;
	}
	h1 {
		@apply font-semibold text-3xl text-center;
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
