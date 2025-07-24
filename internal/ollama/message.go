package ollama

// Message is something we want to tell or ask the llm
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
