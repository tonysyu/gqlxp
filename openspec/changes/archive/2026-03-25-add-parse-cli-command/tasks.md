## 1. Implement parse command

- [x] 1.1 Create `cli/parse.go` with `parseCommand()` that accepts an optional file path argument and `--schema/-s` flag
- [x] 1.2 Implement input reading: read from the file path argument when provided, or from stdin when omitted
- [x] 1.3 Add operation validation using `gqlparser/v2/validator.Validate()` against the resolved schema
- [x] 1.4 Format and print validation errors as `<source>:<line>:<col>: <message>` (one per line to stdout), exiting with code 1
- [x] 1.5 Register `parseCommand()` in `cli/main.go` alongside `searchCommand()` and `showCommand()`

## 2. Tests

- [x] 2.1 Test valid operation from file exits with code 0 and prints no output
- [x] 2.2 Test syntax error in operation prints error with correct format and exits with code 1
- [x] 2.3 Test unknown field reference prints error identifying the field and exits with code 1
- [x] 2.4 Test multiple errors are all reported (one per line)
- [x] 2.5 Test stdin input uses `<stdin>` as the source in error messages
