<script lang="ts">
	import { fly } from "svelte/transition"
	import AdminBalance from "../components/AdminBalance.svelte"
	import Bar from "../components/Bar.svelte"
	import Button from "../components/Button.svelte"
	import { state, ws } from "../lib/amvillage"
	import { status } from "../lib/state"

	$: teamName = $state.config.teams[$state.team]?.name ?? ""
	$: isAdmin = $state.config.teams[$state.team].is_admin
	$: balance = $state.balances[$state.team] ?? Array($state.config.currencies.length)
	$: sum = Math.min(...balance)
	$: max = Math.max(...balance)
	$: lock = $state.locks[$state.team]

	let loading = false
	let notice = ""
	const trade = (i: number) => () => {
		if (!isAdmin) $ws.send("lock")
		loading = true
		let interval: number
		interval = setInterval(() => {
			if ($state.locks[$state.team] !== null || isAdmin) {
				$status = {
					status: "trade",
					target: i,
				}
				clearInterval(interval)
				loading = false
			}
		}, 100)
	}
	const sendNotice = () => {
		$ws.send("notice " + notice)
	}
</script>

<main transition:fly={{ y: 500 }}>
	<div class="welcome">
		欢迎回来，<b>{$state.username}</b>！
	</div>
	<div class="balance">
		{#if isAdmin}
			<AdminBalance />
		{:else}
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
		{/if}
	</div>
	<div class="trade">
		{#if $state.config.teams[$state.team].is_admin || lock === null}
			{#each $state.config.teams as team, i}
				{#if i !== $state.team}
					<Button on:click={trade(i)} disabled={loading}>{team.name}</Button>
				{/if}
			{/each}
		{:else}
			<span class="error">
				{lock.member}正在交易中（{Math.floor(lock.millis_left / 1000)}秒后失效）
			</span>
		{/if}
	</div>
	{#if isAdmin}
		<hr />
		<div class="adminButtons">
			<div class="notice">
				<span>突发事件</span>
				<input bind:value={notice} />
				<Button on:click={sendNotice}>发送</Button>
			</div>
		</div>
	{/if}
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
		@apply flex flex-col items-center bg-purple-800 text-white text-center gap-2 p-4;
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
	.error {
		@apply text-red-800;
	}
	.adminButtons {
		@apply flex flex-col items-center;
	}
	.notice {
		@apply flex gap-2 items-center;
	}
	input {
		@apply border border-black;
	}
</style>
