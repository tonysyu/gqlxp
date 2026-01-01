---
name: Repository Test Review
description: Review all tests in the repository for coverage, readability, and best practices.
category: Repository
tags: [review, testing, quality, coverage]
---

Review all test code in the repository for coverage, readability, and adherence to testing best practices.

**Scope**
Reviews all `*_test.go` files and test directories (`tests/`). Analyzes unit tests, acceptance tests, and fitness tests.

**Standards Checked**

1. **Code Coverage**
   - Run `go test -cover ./...` to check coverage metrics
   - Identify packages with low coverage (<70%)
   - Find untested critical paths and edge cases
   - Suggest additional test cases for uncovered code

2. **Test Readability**
   - Clear test names that describe behavior (not implementation)
   - Proper use of test tables for similar scenarios
   - Minimal setup/teardown complexity
   - Self-documenting assertions with helpful failure messages
   - Appropriate use of test helpers vs inline code
   - See `docs/coding.md` for readability guidelines

3. **Test Level Appropriateness**
   - Identify functionality tested at unit level that should use acceptance tests
   - Check if acceptance tests verify complete user workflows
   - Ensure unit tests focus on single components and edge cases
   - Validate fitness tests enforce architectural constraints
   - Flag over-complicated unit tests that mock excessively

4. **General Testing Best Practices**
   - Proper package naming (`package foo_test` vs `package foo`)
   - Use of `github.com/matryer/is` assertion library
   - Deterministic tests (no race conditions, time dependencies)
   - Independent tests (no shared state between tests)
   - Fast test execution (identify slow tests)
   - Clear arrange-act-assert structure

**Output Format**
For each issue found:
- File location (file:line)
- Category of issue
- Current code/pattern
- Suggested improvement
- Rationale based on testing principles

**Steps**
1. Run coverage analysis: `go test -cover ./...`
2. Find all test files in repository
3. Analyze each test file against the four standard categories
4. Identify patterns across tests (good and bad)
5. Provide coverage summary and prioritized recommendations
6. Highlight exemplary tests that follow best practices
