package terminal

import "github.com/charmbracelet/glamour"

// Renderer provides markdown rendering functionality
type Renderer interface {
	Render(in string) (string, error)
}

// NewMarkdownRenderer creates a new markdown renderer with auto-styled output
func NewMarkdownRenderer() (Renderer, error) {
	return glamour.NewTermRenderer(glamour.WithAutoStyle())
}

// RenderMarkdownOrPlain attempts to render markdown, falling back to plain content on error
func RenderMarkdownOrPlain(renderer Renderer, content string) string {
	if renderer != nil {
		rendered, err := renderer.Render(content)
		if err == nil {
			return rendered
		}
	}
	return content
}
