# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## MANDATORY DOCUMENTATION RULE

**ALL documentation relevant to developers and users MUST be placed in either:**
- `README.md` for project overview and getting started information
- `docs/*` for detailed technical documentation

**CLAUDE.md is exclusively for AI instructions and should NOT contain user/developer documentation.**

**When creating new documentation files:**
- Add entry to `docs/index.md` under "Available Documentation"
- Add entry to `README.md` under "Developer Documentation"

## DOCUMENTATION CONSISTENCY RULE

**When discovering inconsistencies between code and documentation:**
- ALWAYS suggest updating the relevant documentation files (`README.md` or `docs/*`) to reflect current code state
- Proactively identify outdated information in documentation during code analysis
- Suggest specific documentation updates rather than just noting the inconsistency

## CODE VALIDATION RULE

**ALWAYS validate code changes by running tests:**
- After making any code changes, run `just test` to ensure all tests pass
- If tests fail, fix the issues before considering the task complete
- This helps catch compilation errors, type mismatches, and regression issues early

## BUILD SYSTEM RULE

**Use just commands instead of direct Go commands:**
- Use `just build` instead of `go build`
- Use `just run` instead of `go run`
- The project uses justfile for standardized build commands

## DOCUMENTATION STYLE RULE

**Keep documentation concise.** Prioritize brevity over completeness.

- Get to the point in 1-2 sentences when possible
- Eliminate examples, benefits lists, and redundant explanations
- Use bullet points, not paragraphs
- Cross-reference instead of duplicating content
- Avoid referencing specific code locations (file:line)

## Project Information

For project overview, development commands, and architecture details, refer to:
- [README.md](README.md) - Project overview and getting started
- [docs/development.md](docs/development.md) - Build, test, and development commands
- [docs/architecture.md](docs/architecture.md) - Package structure and technical details
- [docs/index.md](docs/index.md) - Complete documentation index
