package fitness_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
)

const PACKAGE_PREFIX = "github.com/tonysyu/gqlxp/"

// PACKAGE_HIERARCHY defines the allowed dependency order.
// Packages may only import from packages listed below them.
var PACKAGE_HIERARCHY = []string{
	"cmd", // highest level: Code can import from any package below
	"cli",
	"tui",
	"library",
	"search",
	"gqlfmt",
	"gql",
	"utils", // lowest level: Cannot import from any other internal package
}

// TestPackageHierarchy enforces that top-level packages can only import
// from packages below them in the hierarchy defined by PACKAGE_HIERARCHY.
func TestPackageHierarchy(t *testing.T) {
	is := is.New(t)

	hierarchyLevel := make(map[string]int)
	for i, pkg := range PACKAGE_HIERARCHY {
		hierarchyLevel[pkg] = i
	}

	violations := findAllViolations(hierarchyLevel, is)

	if len(violations) > 0 {
		reportViolations(t, violations)
	}
}

func findAllViolations(hierarchyLevel map[string]int, is *is.I) []string {
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
			fileViolations := checkFileImports(file, pkg, hierarchyLevel, projectRoot)
			violations = append(violations, fileViolations...)
		}
	}

	return violations
}

func collectGoFiles(pkgPath string) ([]string, error) {
	var files []string
	err := filepath.Walk(pkgPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Include only non-test Go files
		if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func checkFileImports(filePath, pkg string, hierarchyLevel map[string]int, projectRoot string) []string {
	violations := []string{}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ImportsOnly)
	if err != nil {
		return violations
	}

	currentLevel := hierarchyLevel[pkg]

	for _, imp := range node.Imports {
		if violation := validateImport(imp, pkg, currentLevel, hierarchyLevel, filePath, projectRoot); violation != "" {
			violations = append(violations, violation)
		}
	}

	return violations
}

func validateImport(imp *ast.ImportSpec, pkg string, currentLevel int, hierarchyLevel map[string]int, filePath, projectRoot string) string {
	importPath := strings.Trim(imp.Path.Value, `"`)

	// Only check imports from our own module
	if !strings.HasPrefix(importPath, PACKAGE_PREFIX) {
		return ""
	}

	importedPkg := extractTopLevelPackage(importPath)
	if importedPkg == "" || importedPkg == pkg {
		return ""
	}

	// Check if imported package is in our hierarchy
	importedLevel, inHierarchy := hierarchyLevel[importedPkg]
	if !inHierarchy {
		return fmt.Sprintf("Package '%s' imports '%s', which is not in hierarchy", pkg, importedPkg)
	}

	// Violation: importing from same level or higher
	if importedLevel <= currentLevel {
		relPath, _ := filepath.Rel(projectRoot, filePath)
		return fmt.Sprintf("Package '%s' (level %d) imports '%s' (level %d) in %s",
			pkg, currentLevel, importedPkg, importedLevel, relPath)
	}

	return ""
}

func extractTopLevelPackage(importPath string) string {
	parts := strings.Split(strings.TrimPrefix(importPath, PACKAGE_PREFIX), "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}

func reportViolations(t *testing.T, violations []string) {
	t.Error("Package hierarchy violations found:")
	for _, v := range violations {
		t.Error("  " + v)
	}
	t.Errorf("\nExpected hierarchy (top to bottom): %s", strings.Join(PACKAGE_HIERARCHY, " -> "))
}
