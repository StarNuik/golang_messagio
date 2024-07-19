package internal

type MessageCodec = JsonCodec[Message]

type Message struct {
	Content string `json:"content"`
}

type Metrics struct {
	MessagesReceived int `json:"messagesReceived"`
}
