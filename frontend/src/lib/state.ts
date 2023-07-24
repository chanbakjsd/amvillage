import { writable } from "svelte/store"

type MainMenu = { status: "mainMenu" }
type Trade = { status: "trade", target: number }

type Status = MainMenu | Trade

export const status = writable<Status>({
	"status": "mainMenu"
})
