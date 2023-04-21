package app

// common props
const (
	TaskQueue       = "trade-capture-task-queue"
	IncomingTopic   = "dev.tradecapture.incoming"
	OutgoingTopic   = "dev.tradecapture.outgoing"
	KafkaAddress    = "localhost:9092"
	MaxConns        = 5
	MaxChannelDepth = 100
)
