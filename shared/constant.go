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

type RabbitMQPayload struct {
	Command CommandType `json:"command"`
	ID      int64       `json:"id"`
}

const DownloadFolder = "./downloads"
const ConvertedFolder = "./converted"

type Status string

const (
	Initiated  Status = "initiated"
	Preparing  Status = "preparing"
	Processing Status = "processing"
	Completed  Status = "completed"
	Failed     Status = "failed"
)
