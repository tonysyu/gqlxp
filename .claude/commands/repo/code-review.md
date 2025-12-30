---
name: Repository Code Review
description: Review entire repository against project coding standards.
category: Repository
tags: [review, standards, quality, codebase]
---

Review all Go code in the repository against coding standards defined in `docs/coding.md`.

**Scope**
Reviews all `.go` files in the repository. Excludes vendor directories and generated code.

**Standards Checked**

1. **Public API Minimization**
   - Exported identifiers (capitalized) that could be private
   - Getter/setter patterns that violate field access guidelines
   - Public fields that should be private with validation

2. **Testing Conventions**
   - Test package naming (`package foo_test` for black-box vs `package foo` for white-box)
   - Use of `github.com/matryer/is` instead of `t.Error`/`t.Errorf`
   - Appropriate test categorization (unit/acceptance/fitness)
   - See `docs/coding.md` for test type guidelines

3. **Error Handling**
   - Panics that should be errors
   - Missing error context with `fmt.Errorf`
   - Late input validation
   - Unhandled error returns

4. **Line of Sight**
   - Deeply nested if statements (>2 levels)
   - Missing early returns for error cases
   - Happy path not left-aligned

**Output Format**
For each issue found:
- File location (file:line)
- Standard violated
- Current code
- Suggested fix
- Rationale

**Steps**
1. Find all Go files in the repository (excluding vendor and generated code)
2. Analyze each file against the four standard categories
3. For test files, apply testing-specific checks
4. Report findings with actionable suggestions
5. Summarize overall compliance and priority fixes
