package main

import (
	"justn0w-bot/configs"
	"justn0w-bot/internal/router"
)

func main() {
	configs.Init()
	r := router.Router()
	r.Run(":8080")
}
