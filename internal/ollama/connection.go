/*
Package ollama can be used to connect to ollama to request a translation
*/
package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Connection defines a llm which we can communicate with
type Connection struct {
	Model  string
	Prompt string
	URL    string
}

// Translate sends the tekst to the llm and returns the result
func (c Connection) Translate(toTranslate string) (translated string, err error) {
	req := Request{
		Model:  c.Model,
		Stream: false,
		Messages: []Message{
			{Role: "system", Content: c.Prompt},
			{Role: "user", Content: toTranslate},
		},
	}
	js, err := json.Marshal(&req)
	if err != nil {
		return "", err
	}
	client := http.Client{}
	httpReq, err := http.NewRequest(http.MethodPost, c.URL, bytes.NewReader(js))
	if err != nil {
		return "", err
	}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return "", err
	} else if httpResp.StatusCode != 200 {
		return "", fmt.Errorf("http response %v", httpResp)
	}
	defer httpResp.Body.Close()
	ollamaResp := Response{}
	err = json.NewDecoder(httpResp.Body).Decode(&ollamaResp)
	return ollamaResp.Message.Content, err
}
