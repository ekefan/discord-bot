package util

// Interaction Context Types
const (
	GUILD int = iota
	BOT_DM
	PRIVATE_CHANNEL
)

// Discord Command Types
const (
	_          int = iota
	CHAT_INPUT     //The default
	USER
	MESSAGE
	PRIMARY_ENTRY_POINT
)

// Discord Integration Types
const (
	GUILD_INSTALL int = iota
	USER_INSTALL
)

// CommandOptions Types
const (
	STRING  = 3
	INTEGER = 4
	BOOLEAN = 5
)

// Bot Command
const (
	TestCmd      = "test"
	ChallengeCmd = "challenge"
)

// Command Options structure for options that can be added
// to a discord slash command
type CommandOptions struct {
	Type        int      `json:"type"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Required    bool     `json:"required"`
	Choices     []string `json:"choices"`
}

// DiscordCommand structure for defining a discord command
type DiscordCommand struct {
	Name              string           `json:"name"`
	Description       string           `json:"description"`
	Type              int              `json:"type"`
	IntergrationTypes []int            `json:"integration_types"`
	Contexts          []int            `json:"contexts"`
	Options           []CommandOptions `json:"options"`
}

// InteractionOptions structure for options associated
// with a discord command interaction request
type InteractionOptions struct {
	Type  int    `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

// CommandInteraction structure for a discord slash
// command interaction
type CommandInteractionData struct {
	ID      string               `json:"id"`
	Name    string               `json:"name"`
	Type    int                  `json:"type"`
	Options []InteractionOptions `json:"options"`
}

type SlashCommandPayload struct {
	Type   int                    `json:"type"`
	Token  string                 `json:"token"`
	Member SlashCommandMember     `json:"member"`
	ID     string                 `json:"id"`
	Data   CommandInteractionData `json:"data"`
}

type SlashCommandMember struct {
	User  MemberUser `json:"user"`
	Roles []string   `json:"roles"`
	// still contains more fields but not required
}

type MemberUser struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Avatar       string `json:"avatar"`
	Dicriminator string `json:"discriminator"`
	PublicFlags  int    `json:"public_flags"`
}
