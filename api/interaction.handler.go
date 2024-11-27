package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ekefan/discord-bot/domain/command"
	"github.com/ekefan/discord-bot/domain/interaction"
)

const (
	// InteractionTypes
	PING                 = 1
	APPLICATION_COMMMAND = 2
	MESSAGE_COMPONENT    = 3

	// Interaction Callback Type
	CHANNEL_MESSAGE_WITH_SOURCE = 4
	PONG                        = 1

	userAgent = "DiscordBot (https://github.com/ekefan/discord-bot, 1.0.0)"
)

func (bs *BotServer) InteractionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset-UTF-8")
	w.Header().Set("User-Agent", userAgent)

	var reqPayload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&reqPayload); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if reqPayload["type"].(float64) == PING {
		bs.HandleDiscordPing(w)
		return
	}

	if reqPayload["type"].(float64) == APPLICATION_COMMMAND {
		var cmdInteraction interaction.SlashCommandInteraction
		dataBytes, _ := json.Marshal(reqPayload)
		err := json.Unmarshal(dataBytes, &cmdInteraction)
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			slog.Error("could not assert the type of the data from discord interaction")
			return
		}
		if cmdInteraction.Data.Name == command.TestCommand {
			bs.HandleTestCmd(w)
			return
		}

		if cmdInteraction.Data.Name == command.ChallengeCommand {
			bs.HandleChanllengeCmd(w, cmdInteraction)
		} else {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			slog.Error("received bad request interaction from discord", "details", "requested command doesn't exist on server")
			return
		}
		return
	}
	if reqPayload["type"].(float64) == MESSAGE_COMPONENT {
		var cmpInteraction interaction.ComponentInteraction
		dataBytes, _ := json.Marshal(reqPayload)
		err := json.Unmarshal(dataBytes, &cmpInteraction)
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			slog.Error("could not assert the type of the data from discord interaction")
			return
		}
		if strings.HasPrefix(cmpInteraction.Data.CustomId, "accept_button_") {
			bs.HandleAcceptComponentInteraction(w, cmpInteraction)
			return
		}
		if strings.HasPrefix(cmpInteraction.Data.CustomId, "select_choice_") {
			bs.HandleChoiceSelectionInteraction(w, cmpInteraction)
			return
		}
	} else {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		slog.Error("received bad request interaction from discord", "details", "interaction type not supported on this server")
		return
	}
}
