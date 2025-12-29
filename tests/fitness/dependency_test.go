package fitness_test

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
	"slices"
)

// ALLOWED_IMPORTERS maps external dependencies to the packages allowed to import them.
var ALLOWED_IMPORTERS = map[string][]string{
	"github.com/blevesearch/bleve":       {"search"},
	"github.com/charmbracelet/bubbles":   {"tui"},
	"github.com/charmbracelet/bubbletea": {"tui"},
	// Glamour used by both cli and tui to pretty-print markdown
	"github.com/charmbracelet/glamour": {"cli", "tui"},
	// Lipgloss used by cli for colored output and tui for styling
	"github.com/charmbracelet/lipgloss": {"cli", "tui"},
	"github.com/muesli/reflow":          {"utils"},
	"github.com/urfave/cli/v3":          {"cli"},
	"github.com/vektah/gqlparser/v2":    {"gql"},
	"golang.org/x/term":                 {"cli"},
}

// TestDependencyRestrictions enforces that certain external dependencies
// can only be imported by specific packages.
func TestDependencyRestrictions(t *testing.T) {
	is := is.New(t)

	violations := findDependencyViolations(is)

	if len(violations) > 0 {
		reportDependencyViolations(t, violations)
	}
}

func findDependencyViolations(is *is.I) []string {
	projectRoot := filepath.Join("..", "..")
	violations := []string{}

	for _, pkg := range PACKAGE_HIERARCHY {
		pkgPath := filepath.Join(projectRoot, pkg)

		// Skip if package doesn't exist
		if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
			continue
		}

		files, err := collectGoFiles(pkgPath)
		is.NoErr(err) // should collect Go files without error

		// Check each file for violations
		for _, file := range files {
			fileViolations := checkDependencyImports(file, pkg, projectRoot)
			violations = append(violations, fileViolations...)
		}
	}

	return violations
}

func checkDependencyImports(filePath, pkg string, projectRoot string) []string {
	violations := []string{}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ImportsOnly)
	if err != nil {
		return violations
	}

	for _, imp := range node.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)

		// Check if this import is restricted
		for restrictedDep, allowedPkgs := range ALLOWED_IMPORTERS {
			if strings.HasPrefix(importPath, restrictedDep) {
				// Check if current package is allowed to import this dependency
				if !slices.Contains(allowedPkgs, pkg) {
					relPath, _ := filepath.Rel(projectRoot, filePath)
					violations = append(violations, fmt.Sprintf(
						"Package '%s' imports restricted dependency '%s' in %s (only %v can import it)",
						pkg, restrictedDep, relPath, allowedPkgs,
					))
				}
			}
		}
	}

	return violations
}

func reportDependencyViolations(t *testing.T, violations []string) {
	t.Error("Dependency restriction violations found:")
	for _, v := range violations {
		t.Error("  " + v)
	}
	t.Error("\nRestricted dependencies:")
	for dep, allowedPkgs := range ALLOWED_IMPORTERS {
		t.Errorf("  %s -> %v", dep, allowedPkgs)
	}
}
