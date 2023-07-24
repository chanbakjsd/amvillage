<script lang="ts">
	import { fly } from "svelte/transition"
	import Button from "../components/Button.svelte"
	import NumberInput from "../components/NumberInput.svelte"
	import { state, ws } from "../lib/amvillage"
	import { status } from "../lib/state"

	$: isAdmin = $state.config.teams[$state.team].is_admin
	$: lock = $state.locks[$state.team]
	$: target = $status.status === "trade" ? $status.target : 0
	const value = Array($state.config.currencies.length).fill(0)
	const returnToMenu = () => {
		$status = {
			status: "mainMenu",
		}
	}
	const cancel = () => {
		returnToMenu()
		$ws.send("cancel")
	}
	const trade = () => {
		$ws.send(`trade ${target} ${value.map(a => `${a}`).join(" ")}`)
		returnToMenu()
	}
</script>

<main transition:fly={{ y: 500 }}>
	<div class="content">
		{#if $state.locks[$state.team] === null}
			<div class="error">交易失败，请重新尝试</div>
		{:else if lock.member !== $state.username}
			<div class="error">
				<span><b>{$state.locks[$state.team].member}</b>正在交易中</span>
			</div>
		{:else}
			<div class="trade">
				<div>
					<span class="explainer">转账至{$state.config.teams[target].name}</span>
					<span class="lock">您剩下<b>{Math.floor(lock.millis_left / 1000)}</b>秒</span>
					<hr />
				</div>
				<div class="control">
					{#each $state.config.currencies as currency, i}
						<div>
							<span>{currency}</span>
							<NumberInput max={isAdmin ? undefined : $state.balances[$state.team][i]} bind:value={value[i]} />
						</div>
					{/each}
				</div>
			</div>
		{/if}
	</div>
	<div class="buttons">
		{#if lock?.member === $state.username}
			<button class="confirm" on:click={trade} disabled={value.every(a => !a)}>成交</button>
		{/if}
		<Button on:click={cancel} classes="w-full rounded-full">回到主页面</Button>
	</div>
</main>

<style lang="postcss">
	main {
		@apply p-4 flex flex-col gap-8 h-full;
	}
	.content {
		@apply flex-grow;
	}
	.error {
		@apply grid justify-center items-center text-red-800 font-bold text-3xl h-full;
	}
	b {
		@apply font-semibold;
	}
	.trade {
		@apply flex flex-col gap-8;
	}
	.trade > div {
		@apply flex flex-col items-center w-full;
	}
	.explainer {
		@apply text-3xl font-bold text-center;
	}
	hr {
		@apply w-full border-black;
	}
	.trade .control {
		@apply flex flex-wrap gap-x-3 gap-y-2 items-center justify-center;
	}
	.control > div {
		@apply flex gap-x-3 items-center;
	}
	.buttons {
		@apply sticky bottom-0 flex flex-col gap-2;
	}
	.confirm {
		@apply bg-green-400 text-2xl w-full p-4 rounded-xl transition-all;
	}
	.confirm:hover {
		@apply bg-green-600;
	}
	.confirm:disabled {
		@apply bg-gray-300 text-gray-700;
	}
</style>
