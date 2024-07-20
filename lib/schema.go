package lib

type MessageReq struct {
	Content string `json:"content"`
}

type Metrics struct {
	MessagesTotal     int
	MessagesLastDay   int
	MessagesLastHour  int
	ProcessedTotal    int
	ProcessedLastDay  int
	ProcessedLastHout int
	// MessagesReceived int
}
