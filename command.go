package s2y

const (
	CMD_TYPE_STR = 3
)

// Sent to the discord API as a created command.
type Command struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Options     []CommandOption `json:"options"`
}

func NewCommand(name, desc string, options ...CommandOption) Command {
	return Command{Name: name, Description: desc, Options: options}
}

// Options for a command.
type CommandOption struct {
	Type        int    `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

// Data received as part of an interaction.
type CommandData struct {
	ID      string
	Name    string
	Options []CommandDataOption
}

// Options for the received data.
type CommandDataOption struct {
	Name, Value string
}
