## 1. Implement parse command

- [ ] 1.1 Create `cli/parse.go` with `parseCommand()` that accepts an optional file path argument and `--schema/-s` flag
- [ ] 1.2 Implement input reading: read from the file path argument when provided, or from stdin when omitted
- [ ] 1.3 Add operation validation using `gqlparser/v2/validator.Validate()` against the resolved schema
- [ ] 1.4 Format and print validation errors as `<source>:<line>:<col>: <message>` (one per line to stdout), exiting with code 1
- [ ] 1.5 Register `parseCommand()` in `cli/main.go` alongside `searchCommand()` and `showCommand()`

## 2. Tests

- [ ] 2.1 Test valid operation from file exits with code 0 and prints no output
- [ ] 2.2 Test syntax error in operation prints error with correct format and exits with code 1
- [ ] 2.3 Test unknown field reference prints error identifying the field and exits with code 1
- [ ] 2.4 Test multiple errors are all reported (one per line)
- [ ] 2.5 Test stdin input uses `<stdin>` as the source in error messages
