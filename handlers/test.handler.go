package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ekefan/discord-bot/util"
)

type RespComponents struct {
	Type       int               `json:"type"`
	Components []ComponentObject `json:"components"`
}
type RespData struct {
	Content    string           `json:"content"`
	Flags      int              `json:"flags"` //optional
	Components []RespComponents `json:"components"`
}

type Resp struct {
	Type int      `json:"type"`
	Data RespData `json:"data,omitempty"`
}

type ComponentObject struct {
	Type     int    `json:"type"`
	Label    string `json:"label"`
	Style    int    `json:"style"`
	CustomId string `json:"custom_id"`
}

// InteractionTypes
const (
	PING                 = 1
	APPLICATION_COMMMAND = 2
	MESSAGE_COMPONENT    = 3
)

// Interaction Callback(Response) Type
const (
	CHANNEL_MESSAGE_WITH_SOURCE = 4
	PONG                        = 1
)

// ComponentTypes
const (
	_ int = iota
	ACTION_ROW
	BUTTON
	STRING_SELECT
	TEXT_INPUT
)

// Button styles
const (
	_ int = iota
	PRIMARY
	SECONDARY
	SUCCESS
	DANGER
)
const userAgent = "DiscordBot (https://github.com/ekefan/discord-bot, 1.0.0)"

func HandleTestCmd(w http.ResponseWriter) {
	resp := Resp{
		Type: CHANNEL_MESSAGE_WITH_SOURCE,
		Data: RespData{
			Content: "Servers Up🤗🙂",
		},
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func HandleChanllengeCmd(w http.ResponseWriter, reqData util.SlashCommandPayload) {
	btnComponent := ComponentObject{
		Type:     BUTTON,
		Label:    "accept",
		Style:    PRIMARY,
		CustomId: fmt.Sprintf("accept_button_%s", reqData.ID),
	}
	respCompnent := RespComponents{
		Type: ACTION_ROW,
		Components: []ComponentObject{
			btnComponent,
		},
	}
	resp := Resp{
		Type: CHANNEL_MESSAGE_WITH_SOURCE,
		Data: RespData{
			Content: fmt.Sprintf("accept challenge from <@%s>", reqData.Member.User.ID),
			Components: []RespComponents{
				respCompnent,
			},
		},
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
