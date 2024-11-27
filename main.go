package main

import (
	"net/http"

	"github.com/ekefan/discord-bot/api"
	"github.com/ekefan/discord-bot/api/middleware"
	"github.com/ekefan/discord-bot/memory"
	"github.com/ekefan/discord-bot/util"
)

func main() {
	config := util.LoadConfig()
	storage := memory.NewInMemory()
	bs := api.NewBotServer(config, storage)
	http.HandleFunc("/interactions", middleware.VerifyDiscordSignature(bs.InteractionsHandler, config))
	http.ListenAndServe(":8080", nil)
}
