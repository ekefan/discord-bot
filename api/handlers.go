package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/ekefan/discord-bot/domain"
	"github.com/ekefan/discord-bot/domain/challenge"
	"github.com/ekefan/discord-bot/domain/interaction"
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

// Message Flags
const (
	EPHEMERAL = 1 << 6
)

func (bs *BotServer) HandleDiscordPing(w http.ResponseWriter) {
	resp := interaction.InteractionResponse{
		Type: PONG,
	}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		slog.Error("error encoding ping response", "details", err.Error())
		return
	}
}

func (bs *BotServer) HandleTestCmd(w http.ResponseWriter) {
	resp := interaction.InteractionResponse{
		Type: CHANNEL_MESSAGE_WITH_SOURCE,
		Data: interaction.ResponseData{
			Content: "Servers Up🤗🙂",
		},
	}
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		slog.Error("error encoding ping response", "details", err.Error())
		return
	}
}

func (bs *BotServer) HandleChanllengeCmd(w http.ResponseWriter, reqData interaction.SlashCommandInteraction) {
	// get challenge and challenger details from request Data
	challengeId := reqData.ID
	challengerId := reqData.Member.User.ID
	choice := reqData.Data.Options[0].Value

	p1 := &domain.Player{
		ID:     challengerId,
		Choice: domain.RpsChoice(choice),
	}
	// create a new challenge
	challenge, err := challenge.NewChallenge(challengeId, p1)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		slog.Error("error encoding ping response", "details", err.Error())
		return
	}
	bs.Store.CreateChallenge(challenge) // support for another context is not provided

	// respond with a message component
	btnComponent := interaction.BtnComponent{
		Type:     BUTTON,
		Label:    "accept",
		Style:    PRIMARY,
		CustomId: fmt.Sprintf("accept_button_%s", challengeId),
	}
	var components interface{}
	components = []interaction.BtnComponent{
		btnComponent,
	}
	respCompnent := interaction.ResponseDataComponent{
		Type:       ACTION_ROW,
		Components: components.([]interaction.BtnComponent),
	}
	resp := interaction.InteractionResponse{
		Type: CHANNEL_MESSAGE_WITH_SOURCE,
		Data: interaction.ResponseData{
			Content: fmt.Sprintf("accept challenge from <@%s>", reqData.Member.User.ID),
			Components: []interaction.ResponseDataComponent{
				respCompnent,
			},
		},
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (bs *BotServer) HandleAcceptComponentInteraction(w http.ResponseWriter, cmpInteraction interaction.ComponentInteraction) {
	challengeId := strings.Replace(cmpInteraction.Data.CustomId, "accept_button_", "", -1)

	strSelect := interaction.StringSelectComponent{
		Type:     STRING_SELECT,
		CustomId: fmt.Sprintf("select_choice_%v", challengeId),
		Options: []interaction.StrSelectOption{
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
	var components interface{}
	components = []interaction.StringSelectComponent{
		strSelect,
	}
	respCompnent := interaction.ResponseDataComponent{
		Type:       ACTION_ROW,
		Components: components.([]interaction.StringSelectComponent),
	}
	cmpRespData := interaction.ResponseData{
		Content: "What is your object of choice?",
		Flags:   EPHEMERAL,
		Components: []interaction.ResponseDataComponent{
			respCompnent,
		},
	}
	resp := interaction.InteractionResponse{
		Type: CHANNEL_MESSAGE_WITH_SOURCE,
		Data: cmpRespData,
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to send interaction response", "error", err.Error())
		return
	}

	go func() {
		// delete the accept message so no other can accept again
		endpoint := fmt.Sprintf("webhooks/%v/%v/messages/%v", bs.Config.AppID, cmpInteraction.Token, cmpInteraction.Message.ID)
		options := DiscordRequestOption{
			Method: DELETE,
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		response, err := bs.DiscordRequest(ctx, endpoint, options)
		if err != nil {
			slog.Error("could not delete accept button compnent", "details", err.Error())
			return
		}
		if response.StatusCode != http.StatusNoContent {
			fmt.Printf("Failed to delete Discord message: unexpected status code %v\n", response.StatusCode)
			return
		}
	}()

}

func (bs *BotServer) HandleChoiceSelectionInteraction(w http.ResponseWriter, cmpInteraction interaction.ComponentInteraction) {
	challengeID := strings.Replace(cmpInteraction.Data.CustomId, "select_choice_", "", -1)

	challenge, err := bs.Store.GetChallenge(challengeID)
	if err != nil {
		// If game expands use a message interaction as the response
		http.Error(w, "Challenge not found", http.StatusNotFound)
		return
	}
	opponentId := cmpInteraction.Member.User.ID
	choice := cmpInteraction.Data.Values[0]
	opponent := &domain.Player{
		ID:     opponentId,
		Choice: domain.RpsChoice(choice),
	}
	challenge.SetOpponent(opponent)
	err = challenge.DetermineChallengeResult()
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		slog.Error("could not determin challenge result", "details", err.Error())
		return
	}
	resultStr, err := challenge.GetResultMsg()
	err = bs.Store.DeleteChallenge(challengeID)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		slog.Error("could not delete a challenge after getting it's result", "details", err.Error())
		return
	}
	resp := interaction.InteractionResponse{
		Type: CHANNEL_MESSAGE_WITH_SOURCE,
		Data: interaction.ResponseData{
			Content: resultStr,
		},
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

	go func() {
		endpoint := fmt.Sprintf("webhooks/%v/%v/messages/%v", bs.Config.AppID, cmpInteraction.Token, cmpInteraction.Message.ID)
		var body interface{}
		body = map[string]interface{}{
			"content":    fmt.Sprintf("Nice choice <@%v>", cmpInteraction.Member.User.ID),
			"components": nil,
		}
		options := DiscordRequestOption{
			Method: PATCH,
			Body:   body.(map[string]interface{}),
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		response, err := bs.DiscordRequest(ctx, endpoint, options)
		if err != nil {
			slog.Error("could not update the select choice message", "details", err.Error())
			return
		}
		if response.StatusCode != http.StatusOK {
			fmt.Printf("Failed to update Discord message: unexpected status code %v\n", response.StatusCode)
			return
		}
	}()
}
