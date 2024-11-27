// interaction is a the object received from discord when a user interacts with the bot
package interaction

type SlashCommandInteraction struct {
	Type    int                `json:"type"` // Create Type for this
	Token   string             `json:"token"`
	Member  SlashCommandMember `json:"member"`
	ID      string             `json:"id"`
	Data    InteractionData    `json:"data"`
	Context int                `json:"context"`
}
type InteractionData struct {
	ID      string               `json:"id"`
	Name    string               `json:"name"`
	Type    int                  `json:"type"` // Create Type for this
	Options []InteractionOptions `json:"options"`
}

type InteractionOptions struct {
	Type  int    `json:"type"` // Create Type for this
	Name  string `json:"name"`
	Value string `json:"value"`
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
	PublicFlags  int    `json:"public_flags"` // Create Type for this
}
