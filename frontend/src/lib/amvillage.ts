import { writable } from "svelte/store"

export type State = {
	config: Config
	balances: number[][]
	locks: (Lock | null)[]
	notice: string

	team: number
	username: string
}

export type Config = {
	currencies: string[]
	gems: string[]
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
	config: { currencies: [], gems: [], teams: [] },
	balances: [],
	locks: [],
	notice: "",
	team: 0,
	username: "",
})

export const connected = writable(false)
export const error = writable("")

export const ws = writable<WebSocket>()

export const connect = (username: string, password: string) => {
	const url = new URL("./ws", location.href)
	url.protocol = location.protocol === "https:" ? "wss:" : "ws:"
	const websocket = new WebSocket(url)
	websocket.onopen = () => {
		connected.set(true)
		websocket.send(`login ${btoa(password)} ${username}`)
	}
	websocket.onclose = () => {
		connected.set(false)
		console.warn("Websocket close")
		setTimeout(() => connect(username, password), 1000)
	}
	websocket.onerror = err => {
		connected.set(false)
		console.warn("Websocket error:", err)
		setTimeout(() => connect(username, password), 1000)
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
			console.warn(content)
			error.set(content)
			break
		default:
			console.warn("unknown message", data)
		}
	}
	ws.set(websocket)
}
