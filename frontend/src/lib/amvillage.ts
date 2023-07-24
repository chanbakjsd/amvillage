import { writable } from "svelte/store"

export type State = {
	config: Config
	balances: number[][]
	locks: (Lock | null)[]

	team: number
	username: string
}

export type Config = {
	currencies: string[]
	teams: {
		name: string
		is_admin: boolean
	}[]
}

type Lock = {
	member: string
	millis_left: number
}

export const state = writable<State>({
	config: { currencies: [], teams: [] },
	balances: [],
	locks: [],
	team: 0,
	username: "",
})

export const connected = writable(false)

export const ws = writable<WebSocket>()

export const connect = () => {
	const url = new URL("ws://localhost:8080/ws") // "./ws", location.href)
	url.protocol = location.protocol === "https:" ? "wss:" : "ws:"
	const websocket = new WebSocket(url)
	websocket.onopen = () => {
		connected.set(true)
		websocket.send("login " + btoa("group1") + " Wen Xu")
	}
	websocket.onclose = () => {
		connected.set(false)
		console.warn("Websocket close")
		setTimeout(connect, 1000)
	}
	websocket.onerror = err => {
		connected.set(false)
		console.warn("Websocket error:", err)
		setTimeout(connect, 1000)
	}
	websocket.onmessage = (msg: MessageEvent) => {
		const data = msg.data as string
		if (data.indexOf(" ") === -1) {
			console.warn("expected space in event", data)
			return
		}
		const cmd = data.slice(0, data.indexOf(" "))
		const content = data.slice(data.indexOf(" "))
		switch (cmd) {
		case "state":
			console.log(JSON.parse(content))
			state.set(JSON.parse(content))
			break
		case "error":
		default:
			console.warn("unknown message", data)
		}
	}
	ws.set(websocket)
}
