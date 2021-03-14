package node

type DisconnectUserMessage struct {
	User string `json:"user"`
	Reason string `json:"reason"`
	Error bool `json:"error"`
}
