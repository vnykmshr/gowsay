package api

// constants
const (
	ActionSay      = "say"
	ActionThink    = "think"
	ActionList     = "list"
	ActionHelp     = "help"
	ActionSurprise = "surprise"
	ActionDefault  = ActionSay

	CommandMoo = "/moo"

	CowDefault = "default"
	CowRandom  = "random"

	MoodDefault = ""
	MoodRandom  = "random"

	ResponseEphemeral = "ephemeral"
	ResponseInChannel = "in_channel"

	FieldToken       = "token"
	FieldEnv         = "TKPENV"
	FieldText        = "text"
	FieldContentType = "Content-Type"

	ValueProduction      = "production"
	ValueApplicationJSON = "application/json"
	ValueDefaultToken    = "gowsay"
)

// ServerConfig the server config
type ServerConfig struct {
	Name string
}

// AppConfig the app config
type AppConfig struct {
	Token   string
	Columns int32
}

// Config the config struct
type Config struct {
	Server ServerConfig
	App    AppConfig
}

// Module the module
type Module struct {
	cfg *Config
}

// SlackResponse the slack response
type SlackResponse struct {
	ResponseType string       `json:"response_type,omitempty"`
	Text         string       `json:"text,omitempty"`
	Attachments  []Attachment `json:"attachments,omitempty"`
}

// Attachment attachment text
type Attachment struct {
	Text string `json:"text,omitempty"`
}
