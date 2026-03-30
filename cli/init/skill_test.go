package initcmd

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestRenderSkillTemplate_HardCoded(t *testing.T) {
	is := is.New(t)

	result, err := renderSkillTemplate(false, "github")
	is.NoErr(err)

	is.True(strings.Contains(result, "-s github"))
	is.True(!strings.Contains(result, "Schema is selected at runtime"))
	is.True(!strings.Contains(result, "list available schemas and ask the user"))
}

func TestRenderSkillTemplate_RuntimeSelection(t *testing.T) {
	is := is.New(t)

	result, err := renderSkillTemplate(true, "github")
	is.NoErr(err)

	is.True(strings.Contains(result, "-s <schema-id>"))
	is.True(!strings.Contains(result, "-s github"))
	is.True(strings.Contains(result, "Schema is selected at runtime"))
	is.True(strings.Contains(result, "list available schemas and ask the user"))
	is.True(strings.Contains(result, "Schema Context Awareness"))
}

func TestRenderSkillTemplate_ContainsExpectedSections(t *testing.T) {
	is := is.New(t)

	result, err := renderSkillTemplate(false, "test")
	is.NoErr(err)

	expectedSections := []string{
		"Available gqlxp Commands",
		"Smart Field Selection",
		"Common Workflows",
		"Query Construction Guidelines",
		"Search Syntax Reference",
		"Response Format",
	}
	for _, section := range expectedSections {
		is.True(strings.Contains(result, section)) // missing section
	}
	is.True(!strings.Contains(result, "Schema Context Awareness"))
}
