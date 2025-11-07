package api

// constants
const (
	ActionList     = "list"
	ActionHelp     = "help"
	ActionSurprise = "surprise"

	CommandMoo = "/moo"

	CowDefault = "default"
	CowRandom  = "random"

	MoodDefault = ""
	MoodRandom  = "random"

	ResponseEphemeral = "ephemeral"
	ResponseInChannel = "in_channel"

	FieldToken       = "token"
	FieldEnv         = "ENV"
	FieldText        = "text"
	FieldContentType = "Content-Type"

	ValueProduction      = "production"
	ValueApplicationJSON = "application/json"
	ValueDefaultToken    = "gowsay"
)

// Module holds handler dependencies
type Module struct {
	token   string
	columns int
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
