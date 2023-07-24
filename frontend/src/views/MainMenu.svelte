<script lang="ts">
	import { onMount } from "svelte"
	import { fly } from "svelte/transition"
	import Bar from "../components/Bar.svelte"
	import Button from "../components/Button.svelte"
	import { connect, state, ws } from "../lib/amvillage"
	import { status } from "../lib/state"

	onMount(() => {
		connect()
	})

	$: teamName = $state.config.teams[$state.team]?.name ?? ""
	$: balance = $state.balances[$state.team] ?? Array($state.config.currencies.length)
	$: sum = Math.min(...balance)
	$: max = Math.max(...balance)

	let loading = false
	const trade = (i: number) => () => {
		$ws.send("lock")
		loading = true
		let interval: number
		interval = setInterval(() => {
			if ($state.locks[$state.team] !== null) {
				$status = {
					status: "trade",
					target: i,
				}
				clearInterval(interval)
				loading = false
			}
		}, 100)
	}
</script>

<main transition:fly={{ y: 500 }}>
	<div class="welcome">
		欢迎回来，<b>{$state.username}</b>！
	</div>
	<div class="balance">
		<p class="label"><b>{teamName}</b>拥有</p>
		<span class="summary">{sum} 庄</span>
		<hr />
		<p class="label">1庄 = 全部资源的组合</p>
		<table>
			{#each balance as v, i}
				<tr class="entry" class:isMin={v === sum}>
					<td><Bar value={v} {max} isMin={v === sum} /></td>
					<td>{v} {$state.config.currencies[i]}</td>
				</tr>
			{/each}
		</table>
	</div>
	<div class="trade">
		{#each $state.config.teams as team, i}
			{#if i !== $state.team}
				<Button on:click={trade(i)} disabled={loading}>{team.name}</Button>
			{/if}
		{/each}
	</div>
</main>

<style lang="postcss">
	main {
		@apply flex flex-col gap-4 py-4;
	}
	.welcome {
		@apply px-4;
	}
	b {
		@apply font-semibold;
	}
	.balance {
		@apply flex flex-col bg-purple-800 text-white text-center gap-2 p-4;
	}
	.summary {
		@apply text-4xl;
	}
	hr {
		@apply border-gray-400;
	}
	table {
		@apply mx-auto border-separate border-spacing-x-2;
	}
	.isMin td {
		@apply text-red-300;
	}
	.trade {
		@apply flex flex-wrap gap-2 justify-center;
	}
</style>
