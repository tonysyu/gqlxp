package prompt

import (
	"testing"

	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/library"
)

// These tests verify the validation logic used by prompt functions.
// Full interactive testing would require input mocking which is complex.

func TestSchemaIDValidation(t *testing.T) {
	is := is.New(t)

	t.Run("valid schema ID passes validation", func(t *testing.T) {
		err := library.ValidateSchemaID("github-api")
		is.NoErr(err)
	})

	t.Run("invalid schema ID fails validation", func(t *testing.T) {
		err := library.ValidateSchemaID("Invalid_ID")
		is.True(err != nil)
	})
}
