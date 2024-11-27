package interaction

type ComponentInteraction struct {
	Type    int                         `json:"type"`
	Token   string                      `json:"token"`
	ID      string                      `json:"id"`
	Data    ComponentData               `json:"data"`
	Member  SlashCommandMember          `json:"member"`
	Message ComponentInteractionMessage `json:"message"`
	Context int                         `json:"context"`
}

type ComponentData struct {
	CustomId      string               `json:"custom_id"`
	ComponentType int                  `json:"component_type"`
	Values       []CmpInteractionValue `json:"values"`
}

type ComponentInteractionMessage struct {
	Type int    `json:"type"`
	ID   string `json:"id"`
}

type CmpInteractionValue string