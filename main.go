package main

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ekefan/discord-bot/handlers"
	"github.com/ekefan/discord-bot/util"
)

func main() {
	config := util.LoadConfig()
	http.HandleFunc("/interactions", verifyDiscordSignature(interactionsHandler, config))
	http.ListenAndServe(":8080", nil)
}

func verifyDiscordSignature(f http.HandlerFunc, c *util.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// read discords security headers
		signature := r.Header.Get("X-Signature-Ed25519")
		timestamp := r.Header.Get("X-Signature-Timestamp")

		if signature == "" || timestamp == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// read request body to be able to get request message
		body, err := io.ReadAll((r.Body))
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
		// restore request body to the stream
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		// get message
		message := append([]byte(timestamp), body...)

		// Decode the keys and verify the signature
		pubKeyBytes, err := hex.DecodeString(c.PublicKey)
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			slog.Error("error decoding public key", "details", err.Error())
			return
		}
		sigBytes, err := hex.DecodeString(signature)
		if err != nil {
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			slog.Error("error decoding discord signature", "details", err.Error())
			return
		}

		if !ed25519.Verify(pubKeyBytes, message, sigBytes) {
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}
		f(w, r)
	}
}

type Resp struct {
	Type int      `json:"type"`
	Data RespData `json:"data,omitempty"`
}

type RespData struct {
	Content string `json:"content"`
}

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

func interactionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset-UTF-8")
	w.Header().Set("User-Agent", userAgent)
	var reqPayload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&reqPayload); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if reqPayload["type"].(float64) == PING {
		resp := Resp{
			Type: PONG,
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			slog.Error("error encoding ping response", "details", err.Error())
			return
		}
		return
	}

	if reqPayload["type"].(float64) == APPLICATION_COMMMAND {
		var cmdInteraction util.SlashCommandPayload
		dataBytes, _ := json.Marshal(reqPayload)
		err := json.Unmarshal(dataBytes, &cmdInteraction)
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			slog.Error("could not assert the type of the data from discord interaction")
			return
		}
		if cmdInteraction.Data.Name == util.TestCmd {
			handlers.HandleTestCmd(w)
		}

		if cmdInteraction.Data.Name == util.ChallengeCmd {
			handlers.HandleChanllengeCmd(w, cmdInteraction)
		}
		return
	}
	if reqPayload["type"].(float64) == MESSAGE_COMPONENT {
		var cmpInteraction util.ComponentInteractionPayload
		dataBytes, _ := json.Marshal(reqPayload)
		err := json.Unmarshal(dataBytes, &cmpInteraction)
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			slog.Error("could not assert the type of the data from discord interaction")
			return
		}
		if strings.HasPrefix(cmpInteraction.Data.CustomId, "accept_button_") {
			handlers.HandleAcceptComponentInteraction(w, cmpInteraction)
			return
		}
		if strings.HasPrefix(cmpInteraction.Data.CustomId, "select_choice_") {
			handlers.HandleChoiceSelectionInteraction(w, cmpInteraction)
			return
		}
	}
}
