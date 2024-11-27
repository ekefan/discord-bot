package command

import (
	"errors"
	"fmt"
)

var (
	ErrCreateSlashCommand  = errors.New("can not create slash command")
	ErrInvalidSlashCommand = errors.New("can not configure nil slash command")
)

type CmdType int
type CmdOptionType int
type CmdIntegrationType int
type CmdContext int


// Bot Command
const (
	TestCommand      = "test"
	ChallengeCommand = "challenge"
)
const (
	// Command Types
	CHAT_INPUT CmdType = 1
	USER       CmdType = 2
	MESSAGE    CmdType = 3
	PRIMARY    CmdType = 4

	// Command Option Types
	STRING      CmdOptionType = 3
	INTEGER     CmdOptionType = 4
	SUB_COMMAND CmdOptionType = 1
	BOOLEAN     CmdOptionType = 5

	// Command IntegrationTypes
	GUILD_INSTALL CmdIntegrationType = 0
	USER_INSTALL  CmdIntegrationType = 1

	// CommandInteractionContext
	GUILD           CmdContext = 0
	BOT_DM          CmdContext = 1
	PRIVATE_CHANNEL CmdContext = 2
)

// SlashCommand is a discord model for slash commands
type SlashCommand struct {
	Name              string               `json:"name"`
	Description       string               `json:"description"`
	Type              CmdType              `json:"type"`
	IntergrationTypes []CmdIntegrationType `json:"integration_types"`
	Contexts          []CmdContext         `json:"context"`
	Options           []CommandOption      `json:"options,omitempty"` // Options can be of different types
}

// CommandOptions is a discord sub-model of Slash Command model
type CommandOption struct {
	Type        CmdOptionType     `json:"type"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Required    bool              `json:"required"`
	Choices     []CmdOptionChoice `json:"options"`
}

type CmdOptionChoice struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Slash command Configuration
type SlashCmdConfiguration func(slashCmd *SlashCommand) error

// NewSlashCOmmand creates a SlashCommand based on the slash command configuration
func NewSlashCommand(configureCmd SlashCmdConfiguration) (*SlashCommand, error) {
	var slashCommand SlashCommand
	if err := configureCmd(&slashCommand); err != nil {
		return nil, fmt.Errorf("%v:%v", ErrCreateSlashCommand, err)
	}
	return &slashCommand, nil
}

// TestCommandConfiguration implements
// a slash command configuration to configure a test command
func WithTestCommandConfiguration(slashCmd *SlashCommand) error {
	if slashCmd == nil {
		return ErrInvalidSlashCommand
	}
	slashCmd.Name = "test"
	slashCmd.Description = "Basic Command"
	slashCmd.Type = CHAT_INPUT
	slashCmd.IntergrationTypes = []CmdIntegrationType{
		GUILD_INSTALL, USER_INSTALL,
	}
	slashCmd.Contexts = []CmdContext{
		GUILD, BOT_DM, PRIVATE_CHANNEL,
	}
	slashCmd.Options = nil
	return nil
}

// ChallengeCOmmandConfiguration implements
// a slash command configuration to configure a challenge command
func WithChallengeCommandConfiguration(slashCmd *SlashCommand) error {
	if slashCmd == nil {
		return ErrInvalidSlashCommand
	}
	slashCmd.Name = "challenge"
	slashCmd.Description = "Challenge to a match of rock paper scissors"
	slashCmd.Type = CHAT_INPUT
	slashCmd.IntergrationTypes = []CmdIntegrationType{
		GUILD_INSTALL, USER_INSTALL,
	}
	slashCmd.Contexts = []CmdContext{
		GUILD, PRIVATE_CHANNEL,
	}
	slashCmd.Options = []CommandOption{
		{
			Type:        STRING,
			Name:        "object",
			Description: "Pick your object",
			Required:    true,
			Choices: []CmdOptionChoice{
				{
					Name:  "Rock",
					Value: "rock",
				}, {
					Name:  "Paper",
					Value: "paper",
				}, {
					Name:  "Scissor",
					Value: "scissor",
				},
			},
		},
	}
	return nil
}
