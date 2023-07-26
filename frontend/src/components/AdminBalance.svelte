<script lang="ts">
	import { state } from "../lib/amvillage"
</script>

<div>
	您是<span class="gold">管理员</span>，以下为组别余额。
</div>
<div class="table-container">
	<table>
		<tr>
			<th class="first" />
			{#each $state.config.teams as team}
				{#if team.is_admin}
					<td>{team.name}（给的数量）</td>
				{:else}
					<td>{team.name}</td>
				{/if}
			{/each}
		</tr>
		{#each $state.config.currencies as currency, i}
			<tr>
				<td class="first">{currency}</td>
				{#each $state.config.teams as team, j}
					{#if team.is_admin}
						<td>{-$state.balances[j][i]}</td>
					{:else}
						<td>{$state.balances[j][i]}</td>
					{/if}
				{/each}
			</tr>
		{/each}
	</table>
</div>

<style lang="postcss">
	.table-container {
		@apply relative max-w-full overflow-x-scroll border border-black;
	}
	table {
		@apply border-separate border-spacing-0;
	}
	th,
	td {
		@apply px-2 py-1 border border-black min-w-[3rem];
	}
	.first {
		@apply sticky left-0 bg-purple-700;
	}
	.gold {
		@apply font-semibold text-yellow-300;
	}
</style>
