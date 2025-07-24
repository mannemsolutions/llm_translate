package main

import (
	"fmt"
	"llm_translate/internal/markdown"
	"llm_translate/internal/ollama"
	"log"
	"os"
	"strings"
	"time"
)

const (
	llmModel     = "phi4"
	llmURL       = "http://localhost:11434/api/chat"
	llmCharRatio = 2.0
	llmWordRatio = 2.0
)

var (
	prompts = []string{
		"You are a translation AI.",
		"Your sole purpose is to translate the provided text from Dutch into English.",
		"The provided text is in markdown format and your output must be markdown too.",
		"Do not add any extra comments, notes, or meta-information about the source text, its language, or the translation process.",
		"Provide only the English translation in markdown format.",
		"For parts you fail to translate, keep the original.",
	}
	prompt        = strings.Join(prompts, " ")
	llmConnection = ollama.Connection{Model: llmModel, Prompt: prompt, URL: llmURL}
)

func main() {
	start := time.Now()
	logger := log.New(os.Stderr, "", 0)

	mdReader := markdown.NewFromStdin()

	logger.Printf("started")
	for {
		part, readErr := mdReader.Read()
		if part.String() == "" {
			logger.Printf("read nothing")
			break
		}
		part = part.TrimSpace()
		if !part.ContainsText() {
			logger.Printf("skipping non-text part: %s", part)
			fmt.Println(part + "\n")
			continue
		}
		if part.IsURL() {
			logger.Printf("skipping url: %s", part)
			fmt.Println(part + "\n")
			continue
		}
		if part.IsPath() {
			logger.Printf("skipping path: %s", part)
			fmt.Println(part + "\n")
			continue
		}

		prefix, text := part.HeaderText()
		translated, err := translate(markdown.Part(text))
		if err != nil {
			logger.Fatal(err)
		}
		fmt.Println(prefix + string(translated.TrimSpace()) + "\n")

		logger.Printf("translated: %s", part)
		logger.Printf("into: %s%s", prefix, translated)
		if readErr != nil {
			break
		}
	}
	logger.Printf("Completed in %v", time.Since(start))
}

func translate(toTranslate markdown.Part) (translated markdown.Part, err error) {
	logger := log.New(os.Stderr, "", 0)
	tl, err := llmConnection.Translate(string(toTranslate))
	if err != nil {
		logger.Fatalln(err)
	}
	translated = markdown.Part(tl).Cleansed()

	charRatio := translated.CharRatio(toTranslate)
	logger.Printf("char ratio: %f", charRatio)
	if charRatio > llmCharRatio {
		logger.Printf("skipping due to char ratio (%d,%d,%f)",
			toTranslate.CharCount(),
			translated.CharCount(),
			translated.CharRatio(toTranslate),
		)
		logger.Printf("part: %s", toTranslate)
		logger.Printf("llm translation: %s", translated)
		return toTranslate, nil
	}

	wordRatio := translated.WordRatio(toTranslate)
	logger.Printf("word ratio: %f", wordRatio)
	if translated.WordRatio(toTranslate) > llmCharRatio {
		logger.Printf("skipping due to char ratio (%d,%d,%f)",
			toTranslate.WordCount(),
			translated.WordCount(),
			translated.WordRatio(toTranslate),
		)
		logger.Printf("part: %s", toTranslate)
		logger.Printf("llm translation: %s", translated)
		return toTranslate, nil
	}
	return translated, nil
}
