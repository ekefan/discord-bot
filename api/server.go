package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/ekefan/discord-bot/domain/command"
	"github.com/ekefan/discord-bot/memory"
	"github.com/ekefan/discord-bot/util"
)

type BotServer struct {
	Config *util.EnvConfig
	Store  memory.ChallangeRespository
}

func NewBotServer(config *util.EnvConfig, store memory.ChallangeRespository) *BotServer {
	return &BotServer{
		Config: config,
		Store:  store,
	}
}

type ReqMethod string

var (
	ErrInvalidReqMethod     = errors.New("request method not supported")
	ErrEncodingRequestBody  = errors.New("could not marshall request body")
	ErrCreateDiscordRequest = errors.New("could not create a discord request")
)

const (
	POST   ReqMethod = "POST"
	GET    ReqMethod = "GET"
	PUT    ReqMethod = "PUT"
	PATCH  ReqMethod = "PATCH"
	DELETE ReqMethod = "DELETE"
)

func (reqMethod ReqMethod) Valid() bool {
	switch reqMethod {
	case POST:
		return true
	case GET:
		return true
	case PUT:
		return true
	case PATCH:
		return true
	case DELETE:
		return true
	default:
		return false
	}
}

type DiscordRequestOption struct {
	Method ReqMethod
	Body   interface{}
}

func (bs *BotServer) DiscordRequest(ctx context.Context, endpoint string, options DiscordRequestOption) (*http.Response, error) {
	if !options.Method.Valid() {
		return nil, ErrInvalidReqMethod
	}

	var reqBodyIOStream io.Reader
	if options.Body != nil {
		reqBodyBytes, err := json.Marshal(options.Body)
		if err != nil {
			return nil, ErrEncodingRequestBody
		}
		reqBodyIOStream = bytes.NewBuffer(reqBodyBytes)
	}

	url := fmt.Sprintf("%v/%v", bs.Config.DiscordBaseUrl, endpoint)

	request, err := http.NewRequestWithContext(ctx, string(options.Method), url, reqBodyIOStream)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("%v: %v", ErrCreateDiscordRequest, context.DeadlineExceeded)
		}
		if ctx.Err() == context.Canceled {
			return nil, fmt.Errorf("%v: %v", ErrCreateDiscordRequest, ctx.Err())
		}
		return nil, fmt.Errorf("%v: %v", ErrCreateDiscordRequest, err)
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bot %v", bs.Config.DiscordToken))
	request.Header.Add("Content-Type", "application/json; charset=UTF-8")
	request.Header.Add("User-Agent", "DiscordBot (https://github.com/ekefan/discord-bot, 1.0.0)")

	return http.DefaultClient.Do(request)
}

func retryRequest(client *http.Client, request *http.Request, retries int) (*http.Response, error) {
	for i := 0; i < retries; i++ {
		resp, err := client.Do(request)
		if err == nil {
			return resp, nil
		}
		if i < retries-1 {
			time.Sleep(time.Duration(i+1) * time.Second)
		}
	}
	return nil, fmt.Errorf("request failed after %d retries", retries)
}

func (bs *BotServer) InstallGlobalCommands(ctx context.Context, appId, botToken string, commands []command.SlashCommand) error {
	url := fmt.Sprintf("applications/%v/commands", appId)
	options := DiscordRequestOption{
		Method: POST,
		Body:   commands,
	}
	_, err := bs.DiscordRequest(ctx, url, options)
	if err != nil {
		slog.Error("error installing global commands", "details", err)
		return err
	}
	return nil
}
