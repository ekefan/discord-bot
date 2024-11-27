package interaction

// InteractionResponse defines the response payload of an interaction
type InteractionResponse struct {
	Type int          `json:"type"`
	Data ResponseData `json:"data"`
}

// ResponseData is a sub field holding the data of the Interaction Response
type ResponseData struct {
	Content    string                  `json:"content"`
	Flags      int                     `json:"flags,omitempty"` //optional
	Components []ResponseDataComponent `json:"components,omitempty"` //optional
}

// ResponseDataComponent is a sub field of the Response Data of an Interaction Response
//
// It holds the component object of the response.
// Where components can either be an slice of Button Componets or String Select Components.
// Support for other types of components are not necessary in this version of the bot.
type ResponseDataComponent struct {
	Type       int   `json:"type"` // create Type for this
	Components interface{} `json:"components"`
}

type BtnComponent struct {
	Type     int    `json:"type"` // create Type for this
	Label    string `json:"label"`
	Style    int    `json:"style"`
	CustomId string `json:"custom_id"`
}

type StringSelectComponent struct {
	Type     int               `json:"type"` // create type for this
	CustomId string            `json:"custom_id"`
	Options  []StrSelectOption `json:"options"`
}

type StrSelectOption struct {
	Label       string `json:"label"`
	Value       string `json:"value"`
	Description string `json:"description"`
}
