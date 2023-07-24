import { writable } from "svelte/store"

type Login = { status: "login" }
type MainMenu = { status: "mainMenu" }
type Trade = { status: "trade", target: number }

type Status = Login | MainMenu | Trade

export const status = writable<Status>({
	"status": "login"
})
