package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Model              string    `json:"model"`
	CreatedAt          time.Time `json:"created_at"`
	Message            Message   `json:"message"`
	Done               bool      `json:"done"`
	TotalDuration      int64     `json:"total_duration"`
	LoadDuration       int       `json:"load_duration"`
	PromptEvalCount    int       `json:"prompt_eval_count"`
	PromptEvalDuration int       `json:"prompt_eval_duration"`
	EvalCount          int       `json:"eval_count"`
	EvalDuration       int64     `json:"eval_duration"`
}

const (
	llmModel = "phi4"
	llmURL   = "http://localhost:11434/api/chat"
)

var (
	prompts = []string{
		"You are a translation AI.",
		"Your sole purpose is to translate the provided text into English.",
		"Do not add any extra comments, notes, or meta-information about the source text, its language, or the translation process.",
		"Provide only the English translation.",
		"For parts you fail to translate, keep the original text",
	}
	prompt        = strings.Join(prompts, " ")
	containsText  = regexp.MustCompile("[a-zA-Z]")
	justText      = regexp.MustCompile(`[a-zA-Z :.,/\0-9%?!]+`)
	isURL         = regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)
	isPath        = regexp.MustCompile(`^(.+)\/([^\/]+)$`)
	remoteAiNotes = regexp.MustCompile(`\n?\(Note: .*?\.\)`)
)

func cleanse(in string) (cleansed string) {
	out := []byte(in)
	out = remoteAiNotes.ReplaceAll([]byte(out), []byte{})
	return string(out)
}

func main() {
	var (
		readError error
		line      string
	)
	start := time.Now()
	logger := log.New(os.Stderr, "", 0)
	reader := bufio.NewReader(os.Stdin)

	for {
		if readError != nil {
			break
		}
		line, readError = reader.ReadString('\n')
		orgLine := line
		allParts := justText.FindAllString(line, -1)
		if allParts == nil {
			fmt.Println(line)
			continue
		}
		for _, part := range allParts {
			if !containsText.Match([]byte(part)) {
				continue
			}
			if isURL.MatchString(part) {
				continue
			}
			if isPath.MatchString(part) {
				continue
			}
			translated, err := translate(part)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			translated = cleanse(translated)
			line = strings.ReplaceAll(line, part, translated)
			/*
				if strings.HasPrefix(resp.Message.Content, "I'm unable to translate the text") {
					fmt.Println(line)
					continue
				}
			*/
		}
		fmt.Println(line)
		logger.Printf("%s > %s", orgLine, line)
	}
	logger.Printf("Completed in %v", time.Since(start))
}

func translate(toTranslate string) (translated string, err error) {
	req := Request{
		Model:  llmModel,
		Stream: false,
		Messages: []Message{
			{Role: "system", Content: prompt},
			{Role: "user", Content: toTranslate},
		},
	}
	js, err := json.Marshal(&req)
	if err != nil {
		return "", err
	}
	client := http.Client{}
	httpReq, err := http.NewRequest(http.MethodPost, llmURL, bytes.NewReader(js))
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
