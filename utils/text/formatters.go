package text

import (
	"fmt"
	"strings"

	"github.com/muesli/reflow/wordwrap"
)

const ellipsis = "…"

func JoinParagraphs(lines ...string) string {
	return strings.Join(lines, "\n\n")
}

func JoinLines(lines ...string) string {
	return strings.Join(lines, "\n")
}

func SplitLines(text string) []string {
	return strings.Split(text, "\n")
}

func WrapAndTruncate(rawString string, maxWidth, maxLines int) string {
	text := wordwrap.String(rawString, maxWidth)
	textLines := SplitLines(text)

	// Truncate to max height and add ellipsis if needed
	if len(textLines) <= maxLines {
		return text
	}

	textLines = textLines[:maxLines]
	lastLine := textLines[maxLines-1]
	if len(lastLine) > 0 {
		if len(lastLine) >= maxWidth {
			// Replace last character with ellipsis
			lastLine = lastLine[:len(lastLine)-1] + ellipsis
		} else {
			lastLine = lastLine + ellipsis
		}
		textLines[maxLines-1] = lastLine
	}
	return JoinLines(textLines...)
}

func GqlCode(code string) string {
	return JoinLines("```graphqls", code, "```")
}

func GqlDocString(description string) string {
	return fmt.Sprintf(`"""%s"""`, description)
}

func H1(title string) string {
	return "# " + title
}
