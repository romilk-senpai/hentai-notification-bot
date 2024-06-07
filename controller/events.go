package events

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(e Event) error
}

type EventType int

const (
	Unknown EventType = iota
	Message
	Command
)

type Event struct {
	UserHash    string
	Type        EventType
	Text        string
	CommandInfo CommandInfo
	Meta        interface{}
}

type CommandInfo struct {
	Command   string
	Arguments string
}
