package initcmd

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"
)

//go:embed templates/gqlxp_skill.template.md
var skillTemplateContent string

type skillTemplateData struct {
	SchemaFlag       string
	RuntimeSelection bool
}

// renderSkillTemplate renders the skill template with the given schema configuration.
func renderSkillTemplate(runtimeSelection bool, selectedSchemaID string) (string, error) {
	data := skillTemplateData{
		SchemaFlag:       fmt.Sprintf("-s %s", selectedSchemaID),
		RuntimeSelection: runtimeSelection,
	}
	if runtimeSelection {
		data.SchemaFlag = "-s <schema-id>"
	}

	tmpl, err := template.New("skill").Parse(skillTemplateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse skill template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to render skill template: %w", err)
	}

	return buf.String(), nil
}

