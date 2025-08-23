package main

import "justn0w-bot/internal/router"

func main() {
	r := router.Router()
	r.Run(":8080")
}
