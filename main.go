package main

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

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

		fmt.Println("connected")
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

	// InteractionContext
	GUILD  = 0
	BOT_DM = 1

	// Interaction Callback Type
	CHANNEL_MESSAGE_WITH_SOURCE = 4
	PONG                        = 1


	userAgent = "DiscordBot (https://github.com/ekefan/discord-bot, 1.0.0)"
)

func interactionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if payload["type"].(float64) == PING {
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

		if payload["type"].(float64) == APPLICATION_COMMMAND {
			resp := Resp{
				Type: CHANNEL_MESSAGE_WITH_SOURCE,
				Data: RespData{
					Content: "received",
				},
			}
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("User-Agent", userAgent)
			err := json.NewEncoder(w).Encode(resp)
			if err != nil {
				http.Error(w, "Server Error", http.StatusInternalServerError)
				slog.Error("error encoding application command response", "details", err.Error())
				return
			}
			return
		}
	}
}
