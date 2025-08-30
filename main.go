package main

import (
	"justn0w-bot/config"
	"justn0w-bot/internal/router"
)

func main() {
	config.Init()
	r := router.Router()
	r.Run(":8080")
}
