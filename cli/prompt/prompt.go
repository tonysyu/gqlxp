package prompt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tonysyu/gqlxp/library"
)

// String prompts the user for a string input with an optional default value.
func String(prompt string, defaultValue string) (string, error) {
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", prompt, defaultValue)
	} else {
		fmt.Printf("%s: ", prompt)
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return "", fmt.Errorf("EOF reached")
		}
		return "", fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.TrimSpace(input)
	if input == "" && defaultValue != "" {
		return defaultValue, nil
	}

	return input, nil
}

// YesNo prompts the user for a yes/no confirmation.
func YesNo(prompt string) (bool, error) {
	fmt.Printf("%s (y/n): ", prompt)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return false, fmt.Errorf("EOF reached")
		}
		return false, fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.ToLower(strings.TrimSpace(input))
	return input == "y" || input == "yes", nil
}

// SchemaID prompts the user for a schema ID with validation.
func SchemaID(suggestedID string) (string, error) {
	for {
		p := "Enter schema ID (lowercase letters, numbers, hyphens)"
		id, err := String(p, suggestedID)
		if err != nil {
			return "", err
		}

		if err := library.ValidateSchemaID(id); err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		return id, nil
	}
}
