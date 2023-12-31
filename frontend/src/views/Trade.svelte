<script lang="ts">
	import { fly } from "svelte/transition"
	import { _ } from "svelte-i18n"
	import Button from "../components/Button.svelte"
	import I18n from "../components/I18n.svelte"
	import NumberInput from "../components/NumberInput.svelte"
	import { state, ws } from "../lib/amvillage"
	import { status } from "../lib/state"

	const value = Array($state.config.currencies.length).fill(0)
	const gemValue = Array($state.config.gems.length).fill(0)
	$: currencyCount = $state.config.currencies.length
	$: isAdmin = $state.config.teams[$state.team].is_admin
	$: lock = $state.locks[$state.team]
	$: target = $status.status === "trade" ? $status.target : 0
	$: hasBalance =
		value.every((num, i) => num <= $state.balances[$state.team][i] || num === 0) &&
		gemValue.every((num, i) => num <= $state.balances[$state.team][currencyCount + i] || num === 0)
	$: tradeOK = (value.some(num => num !== 0) || gemValue.some(num => num !== 0)) && (isAdmin || hasBalance)
	const returnToMenu = () => {
		$status = {
			status: "mainMenu",
		}
	}
	const cancel = () => {
		returnToMenu()
		if (!isAdmin) $ws.send("cancel")
	}
	const trade = () => {
		$ws.send(`trade ${target} ${[...value, ...gemValue].map(a => `${a}`).join(" ")}`)
		returnToMenu()
	}
</script>

<main transition:fly={{ y: 500 }}>
	<div class="content">
		{#if $state.locks[$state.team] === null && !isAdmin}
			<div class="error">{$_("trade.error.timeout")}</div>
		{:else if lock?.member !== $state.username && !isAdmin}
			<div class="error">
				<I18n id="trade.error.othersTrading" values={{ username: lock?.member }} />
			</div>
		{:else}
			<div class="trade">
				<div>
					<span class="explainer">
						{$_("trade.title", { values: { target: $state.config.teams[target].name } })}
					</span>
					{#if !isAdmin}
						<span class="lock">
							<I18n id="trade.timeRemaining" values={{ sec: Math.floor(lock.millis_left / 1000) }} />
						</span>
					{/if}
					<hr />
				</div>
				<div class="control">
					{#each $state.config.currencies as currency, i}
						<span>{currency}</span>
						<NumberInput
							min={isAdmin ? -$state.balances[target][i] : 0}
							max={isAdmin ? undefined : $state.balances[$state.team][i]}
							bind:value={value[i]}
						/>
					{/each}
				</div>
				<div class="control">
					{#each $state.config.gems as gem, i}
						<span>{gem}</span>
						<NumberInput
							min={isAdmin ? -$state.balances[$state.team][currencyCount + i] : 0}
							max={isAdmin ? undefined : $state.balances[$state.team][currencyCount + i]}
							bind:value={gemValue[i]}
						/>
					{/each}
				</div>
			</div>
		{/if}
	</div>
	<div class="buttons">
		{#if lock?.member === $state.username || isAdmin}
			<button class="confirm" on:click={trade} disabled={!tradeOK}>{$_("trade.button.confirm")}</button>
		{/if}
		<Button on:click={cancel} classes="w-full rounded-full">{$_("trade.button.cancel")}</Button>
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
		@apply grid gap-x-3 gap-y-2 items-center justify-center;
		grid-template-columns: auto auto;
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
