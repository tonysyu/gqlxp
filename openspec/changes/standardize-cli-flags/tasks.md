## 1. Implementation

- [x] 1.1 Update `app` command to use `--schema/-s` flag (library selector when omitted)
- [x] 1.2 Update `search` command to use `--schema/-s` flag (default schema when omitted)
- [x] 1.3 Update `show` command to use `--schema/-s` flag (default schema when omitted)
- [x] 1.4 Remove `--default/-d` flag (use default when `--schema` is omitted)
- [x] 1.5 Update command help text and examples

## 2. Testing

- [x] 2.1 Test app command with `--schema` flag
- [x] 2.2 Test app command without flags (library selector)
- [x] 2.3 Test search command with `--schema` flag
- [x] 2.4 Test search command without flags (uses default)
- [x] 2.5 Test show command with `--schema` flag
- [x] 2.6 Test show command without flags (uses default)
- [x] 2.7 Verify all commands work with file paths and library IDs

## 3. Documentation

- [x] 3.1 Update README examples
- [x] 3.2 Update documentation to reflect simplified flag usage
