package ws

type Message struct {
    Username string `json:"username"`
    Text     string `json:"text"`
    ChatID   string `json:"chatid"`
}