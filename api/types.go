package api

// Configuration and environment constants
const (
	envKey            = "ENV"
	envProduction     = "production"
	defaultTokenValue = "gowsay"
)

// Form field names
const (
	fieldToken = "token"
	fieldText  = "text"
)

// Command keywords
const (
	commandHelp     = "help"
	commandList     = "list"
	commandSurprise = "surprise"
	commandRandom   = "random"
)

// Slack response types
const (
	responseEphemeral = "ephemeral"
	responseInChannel = "in_channel"
)

// Default values
const (
	defaultCow = "default"
)

// Module holds handler dependencies
type Module struct {
	token   string
	columns int
}

// SlackResponse represents a Slack-compatible response
type SlackResponse struct {
	ResponseType string       `json:"response_type,omitempty"`
	Text         string       `json:"text,omitempty"`
	Attachments  []Attachment `json:"attachments,omitempty"`
}

// Attachment represents a Slack attachment
type Attachment struct {
	Text string `json:"text,omitempty"`
}
