package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ekefan/discord-bot/util"
)

type RespComponents struct {
	Type       int                  `json:"type"`
	Components []BtnComponentObject `json:"components"`
}
type RespData struct {
	Content    string           `json:"content"`
	Flags      int              `json:"flags,omitempty"` //optional
	Components []RespComponents `json:"components"`
}

type Resp struct {
	Type int      `json:"type"`
	Data RespData `json:"data,omitempty"`
}

type ComponentResponse struct {
	Type int                   `json:"type"`
	Data ComponentResponseData `json:"data"`
}

type ComponentResponseData struct {
	Content    string                        `json:"content"`
	Flags      int                           `json:"flags"`
	Components []StrSelectMsgActionRowObject `json:"components"`
}

type StrSelectComponentObject struct {
	Type     int               `json:"type"`
	CustomId string            `json:"custom_id"`
	Options  []StrSelectOption `json:"options"`
}

type StrSelectOption struct {
	Label       string `json:"label"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type StrSelectMsgActionRowObject struct {
	Type       int                        `json:"type"`
	Components []StrSelectComponentObject `json:"components"`
}
type BtnComponentObject struct {
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

// Message Flags
const (
	EPHEMERAL = 1 << 6
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
	btnComponent := BtnComponentObject{
		Type:     BUTTON,
		Label:    "accept",
		Style:    PRIMARY,
		CustomId: fmt.Sprintf("accept_button_%s", reqData.ID),
	}
	respCompnent := RespComponents{
		Type: ACTION_ROW,
		Components: []BtnComponentObject{
			btnComponent,
		},
	}
	resp := Resp{
		Type: CHANNEL_MESSAGE_WITH_SOURCE,
		Data: RespData{ // TODO: should ComponentMsgData
			Content: fmt.Sprintf("accept challenge from <@%s>", reqData.Member.User.ID),
			Components: []RespComponents{
				respCompnent,
			},
		},
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func HandleAcceptComponentInteraction(w http.ResponseWriter, cmpInteraction util.ComponentInteractionPayload) {
	gameID := strings.Replace(cmpInteraction.Data.CustomId, "accept_button_", "", -1)
	fmt.Println(gameID)
	strSelect := StrSelectComponentObject{
		Type:     STRING_SELECT,
		CustomId: fmt.Sprintf("select_choice_%v", gameID),
		Options: []StrSelectOption{
			{
				Label:       "Rock",
				Value:       "rock",
				Description: "sedimentary, igneous, or perphaps even metamorphic",
			}, {
				Label:       "Scissors",
				Value:       "scissors",
				Description: "careful ! sharp ! edges !!",
			}, {
				Label:       "Paper",
				Value:       "paper",
				Description: "versatile and iconic",
			},
		},
	}
	actionRow := StrSelectMsgActionRowObject{
		Type: ACTION_ROW,
		Components: []StrSelectComponentObject{
			strSelect,
		},
	}
	cmpRespData := ComponentResponseData{
		Content: "What is your object of choice?",
		Flags:   EPHEMERAL,
		Components: []StrSelectMsgActionRowObject{
			actionRow,
		},
	}
	resp := ComponentResponse{
		Type: CHANNEL_MESSAGE_WITH_SOURCE,
		Data: cmpRespData,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func HandleChoiceSelectionInteraction(w http.ResponseWriter, cmpInteraction util.ComponentInteractionPayload) {
	gameID := strings.Replace(cmpInteraction.Data.CustomId, "select_choice_", "", -1)
	fmt.Println(gameID)

	context := cmpInteraction.Context
	fmt.Println(context)
	fmt.Println(cmpInteraction.Member)
	fmt.Println(cmpInteraction.ID)
	resp := Resp{
		Type: CHANNEL_MESSAGE_WITH_SOURCE,
		Data: RespData{
			Content: fmt.Sprintf("nice move <@%v>", cmpInteraction.Member.User.ID),
		},
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
