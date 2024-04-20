package shared

type CommandType string

const (
	StartCommand CommandType = "start"
	HelpCommand  CommandType = "help"
	PhotoCommand CommandType = "photo"
)

type UpdateType string

const (
	Message  UpdateType = "message"
	Callback UpdateType = "callback"
)
