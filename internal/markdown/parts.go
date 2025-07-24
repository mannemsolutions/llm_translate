/*
Package markdown is a simple package to split markdown files in parts and collect info on it
*/
package markdown

import (
	"regexp"
	"strings"
)

var (
	containsText    = regexp.MustCompile("[a-zA-Z]")
	justText        = regexp.MustCompile(`[a-zA-Z :.,/\0-9%?!\n-]+`)
	isURL           = regexp.MustCompile(`^https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)$`)
	fileChars       = `a-zA-Z0-9 _`
	doubleDot       = `\.{1,2}`
	fileNameWithDot = `[` + fileChars + `][` + fileChars + `.]*`
	dotWithFileName = `[` + fileChars + `.]*[` + fileChars + `]`
	homedir         = `(~[` + fileChars + `]*/|/)?`
	file            = `(` + doubleDot + `|` + fileNameWithDot + `|` + dotWithFileName + `)?`
	dir             = `((` + doubleDot + `|` + fileNameWithDot + `|` + dotWithFileName + `)/)*`
	isPath          = regexp.MustCompile(`^` + homedir + dir + file + `$`)
	isWord          = regexp.MustCompile(`[a-zA-Z0-9]+`)
	remoteAiNotes   = regexp.MustCompile(`\(Note: .*?\.\)`)
	headerRegexp    = regexp.MustCompile(`^(?P<prefix>#+ *)(?P<text>.*)`)
)

// Parts represents all parts of an md file
type Parts struct {
	document Part
	parts    []Part
}

// NewParts splits an md file into chunks and returns a Parts with every chunk being a new Part
func NewParts(md string) Parts {
	parts := Parts{document: Part(md)}

	for _, node := range strings.Split(md, "\n\n") {
		parts.parts = append(parts.parts, Part(node))
	}
	return parts
}

// Part is a part of an md file
type Part string

func (p Part) String() string {
	return string(p)
}

// Cleansed removes Ai notes if they are in there
func (p Part) Cleansed() Part {
	return Part(remoteAiNotes.ReplaceAllString(p.String(), ""))
}

// ContainsText checks if the node has text values
func (p Part) ContainsText() bool {
	return containsText.MatchString(p.String())
}

// IsHeader checks if a chunk is a header
func (p Part) IsHeader() bool {
	str := p.String()
	return strings.HasPrefix(str, "#") && !strings.Contains(str, "\n")
}

// HeaderText returns the prefix and text of a header
func (p Part) HeaderText() (string, string) {
	if !p.IsHeader() {
		return "", p.String()
	}
	match := headerRegexp.FindStringSubmatch(p.String())

	paramsMap := map[string]string{}
	for i, name := range headerRegexp.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return paramsMap["prefix"], paramsMap["text"]
}

// IsURL checks if a chunk is url (we want to skip them)
func (p Part) IsURL() bool {
	return isURL.MatchString(p.String())
}

// IsPath checks if a Part is a path (we want to skip them)
func (p Part) IsPath() bool {
	return isPath.MatchString(p.String())
}

// WordCount returns the number of words
func (p Part) WordCount() int {
	return len(isWord.FindAllString(p.String(), -1))
}

// WordRatio returns the ratio between number of words in pre and post Part
func (p Part) WordRatio(other Part) float32 {
	return float32(p.WordCount()) / float32(other.WordCount())
}

// CharCount returns the number of words
func (p Part) CharCount() int {
	return len(p.String())
}

// CharRatio returns the ratio between number of words in pre and post Part
func (p Part) CharRatio(other Part) float32 {
	return float32(p.CharCount()) / float32(other.CharCount())
}
