package tui

import (
	"strings"
)

func joinParagraphs(lines... string) string {
	return strings.Join(lines, "\n\n")
}

func joinLines(lines... string) string {
	return strings.Join(lines, "\n")
}

func gqlCode(code string) string {
	return joinLines("```graphqls", code, "```")
}

func h1(title string) string {
	return "# " + title
}
