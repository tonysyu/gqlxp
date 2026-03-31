package initcmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	clilib "github.com/tonysyu/gqlxp/cli/library"
	"github.com/tonysyu/gqlxp/cli/prompt"
	"github.com/tonysyu/gqlxp/library"
	"github.com/urfave/cli/v3"
)

// Command creates the init subcommand.
func Command() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Interactive setup wizard for gqlxp",
		Description: `Walks you through configuring a schema and optionally
installing the /gqlxp skill for AI assistants like Claude Code.`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return runInitWizard(ctx)
		},
	}
}

func runInitWizard(ctx context.Context) error {
	fmt.Println("Welcome to gqlxp!")
	fmt.Println()

	lib := library.NewLibrary()

	// Step 1: Schema setup (only if library is empty)
	schemas, err := lib.List()
	if err != nil {
		return fmt.Errorf("failed to list schemas: %w", err)
	}
	if len(schemas) == 0 {
		fmt.Println("--- Schema Setup ---")
		schemaID, err := addNewSchema(ctx, lib)
		if err != nil {
			return err
		}
		if err := setDefault(lib, schemaID); err != nil {
			return err
		}
	}

	// Step 2: AI skill installation
	fmt.Println("\n--- AI Skill Setup ---")
	if err := installSkill(lib); err != nil {
		return err
	}

	fmt.Println("\nSetup complete!")
	return nil
}

func addNewSchema(ctx context.Context, lib library.Library) (string, error) {
	source, err := prompt.String("Enter schema file path or URL", "")
	if err != nil {
		return "", err
	}
	if source == "" {
		return "", fmt.Errorf("schema source is required")
	}

	content, _, err := clilib.LoadSchemaContent(ctx, source, nil)
	if err != nil {
		return "", err
	}

	if err := clilib.ValidateSchema(content); err != nil {
		return "", err
	}

	suggested := clilib.ExtractSuggestedID(source)
	schemaID, err := prompt.SchemaID(suggested)
	if err != nil {
		return "", err
	}

	displayName, err := prompt.String("Enter display name", schemaID)
	if err != nil {
		return "", err
	}

	if err := lib.AddFromContent(schemaID, displayName, content, source); err != nil {
		return "", err
	}

	fmt.Printf("Added schema '%s' to library\n", schemaID)
	return schemaID, nil
}

func setDefault(lib library.Library, schemaID string) error {
	currentDefault, _ := lib.GetDefaultSchema()
	if currentDefault == schemaID {
		fmt.Printf("'%s' is already the default schema.\n", schemaID)
		return nil
	}

	setIt, err := prompt.YesNo(fmt.Sprintf("Set '%s' as default schema?", schemaID))
	if err != nil {
		return err
	}

	if setIt {
		if err := lib.SetDefaultSchema(schemaID); err != nil {
			return fmt.Errorf("failed to set default schema: %w", err)
		}
		fmt.Printf("Default schema set to '%s'\n", schemaID)
	}
	return nil
}

func installSkill(lib library.Library) error {
	install, err := prompt.YesNo("Install /gqlxp skill for AI assistants (e.g. Claude Code)?")
	if err != nil || !install {
		return err
	}

	fmt.Println("Schema mode:")
	fmt.Println("  1. Prompt at runtime (list schemas and ask each session)")
	fmt.Println("  2. Hard-code a specific schema")
	choice, err := prompt.String("Select an option", "1")
	if err != nil {
		return err
	}
	runtimeSelection := strings.TrimSpace(choice) != "2"

	selectedSchemaID := ""
	if !runtimeSelection {
		selectedSchemaID, err = selectSchema(lib)
		if err != nil {
			return err
		}
	}

	content, err := renderSkillTemplate(runtimeSelection, selectedSchemaID)
	if err != nil {
		return err
	}

	skillDir := filepath.Join(".claude", "skills", "gqlxp")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", skillDir, err)
	}

	skillPath := filepath.Join(skillDir, "SKILL.md")
	if err := os.WriteFile(skillPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("failed to write skill file: %w", err)
	}
	fmt.Printf("Created %s\n", skillPath)
	return nil
}

func selectSchema(lib library.Library) (string, error) {
	schemas, err := lib.List()
	if err != nil {
		return "", fmt.Errorf("failed to list schemas: %w", err)
	}

	fmt.Println("Available schemas:")
	for i, s := range schemas {
		fmt.Printf("  %d. %s (%s)\n", i+1, s.ID, s.DisplayName)
	}

	choice, err := prompt.String("Select a schema", "1")
	if err != nil {
		return "", err
	}

	num, err := strconv.Atoi(strings.TrimSpace(choice))
	if err != nil || num < 1 || num > len(schemas) {
		return "", fmt.Errorf("invalid selection: %s", choice)
	}

	return schemas[num-1].ID, nil
}
