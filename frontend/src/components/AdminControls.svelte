<script lang="ts">
	import { _ } from "svelte-i18n"
	import Button from "../components/Button.svelte"
	import { ws } from "../lib/amvillage"

	let notice = ""
	let popup = ""
	const sendNotice = () => {
		$ws.send("notice " + notice)
	}
	const sendPopup = () => {
		$ws.send("popup " + popup)
	}
	const hidePopup = () => {
		$ws.send("popup")
	}
</script>

<div class="controls">
	<span>{$_("admin.label.notice")}</span>
	<input bind:value={notice} />
	<Button on:click={sendNotice}>{$_("admin.button.sendNotice")}</Button>
	<span>{$_("admin.label.banner")}</span>
	<input bind:value={popup} />
	<div class="flex flex-wrap justify-center gap-1">
		<Button on:click={sendPopup}>{$_("admin.button.sendBanner")}</Button>
		<Button on:click={hidePopup}>{$_("admin.button.hideBanner")}</Button>
	</div>
</div>

<style lang="postcss">
	.controls {
		@apply grid justify-center items-center p-2 gap-2;
		grid-template-columns: auto 10rem auto;
	}
	input {
		@apply border border-black;
	}
</style>
