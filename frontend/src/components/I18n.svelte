<script lang="ts">
	import { getMessageFormatter, json } from "svelte-i18n"

	export let id: string
	export let values: Record<string, any> = {}

	type BoldedText = {
		type: "bold"
		content: string
	}
	type GoldText = {
		type: "gold"
		content: string
	}

	type CustomText = BoldedText | GoldText
	type I18nChunk = {
		value: string | CustomText
	}

	let chunks: I18nChunk[] = []
	$: {
		try {
			let message = $json(id)
			if (typeof message !== "string") {
				console.warn(`Unexpected type for message ${id}: `, message)
				message = id
			}
			chunks = getMessageFormatter(message as string).formatToParts<CustomText>({
				...values,
				b: (chunks: CustomText[]): BoldedText => {
					// Bold everything inside.
					const content = chunks
						.map(chunk => {
							if (typeof chunk === "string") return chunk
							return chunk.content
						})
						.join("")
					return { type: "bold", content }
				},
				gold: (chunks: CustomText[]): GoldText => {
					// Gold everything inside.
					const content = chunks
						.map(chunk => {
							if (typeof chunk === "string") return chunk
							return chunk.content
						})
						.join("")
					return { type: "gold", content }
				},
			})
		} catch (err: any) {
			console.warn("Error formatting:", err.message)
			chunks = [{ value: err.message }]
		}
	}
</script>

{#each chunks as c}
	{#if typeof c.value === "string"}
		{c.value}
	{:else if c.value.type === "bold"}
		<b>{c.value.content}</b>
	{:else}
		<span class="gold">{c.value.content}</span>
	{/if}
{/each}

<style lang="postcss">
	b {
		@apply font-semibold;
	}
	.gold {
		@apply font-semibold text-yellow-300;
	}
</style>
