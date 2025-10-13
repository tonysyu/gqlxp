package text

import (
	"strings"
)

func JoinParagraphs(lines ...string) string {
	return strings.Join(lines, "\n\n")
}

func JoinLines(lines ...string) string {
	return strings.Join(lines, "\n")
}

func GqlCode(code string) string {
	return JoinLines("```graphqls", code, "```")
}

func H1(title string) string {
	return "# " + title
}
