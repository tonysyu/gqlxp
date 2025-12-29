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
