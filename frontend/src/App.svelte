<script lang="ts">
	import { onMount } from "svelte"
	import Bar from "./components/Bar.svelte"
	import { connect, connected, state } from "./lib/amvillage"

	onMount(() => {
		connect()
	})

	$: teamName = $state.config.teams[$state.team]?.name ?? ""
	$: balance = $state.balances[$state.team] ?? Array($state.config.currencies.length)
	$: sum = Math.min(...balance)
	$: max = Math.max(...balance)
</script>

<main>
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
</main>

{#if !$connected}
	<div class="warning">⚠️</div>
{/if}

<style lang="postcss">
	.welcome {
		@apply p-4;
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
	.warning {
		@apply absolute bottom-2 right-2 text-5xl animate-bounce;
	}
</style>
