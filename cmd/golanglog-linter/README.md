# golanglog-linter — standalone binary

Standalone binary for running the linter via `go vet` without golangci-lint.

## Build

```bash
go build -o golanglog-linter ./cmd/golanglog-linter
```

## Usage

```bash
# via go vet (local build)
go vet -vettool=./golanglog-linter ./...

# via go install
go install github.com/romariok/golanglog-linter/cmd/golanglog-linter@latest
go vet -vettool=$(which golanglog-linter) ./...
```

## Flags

| Flag | Type | Default | Description |
|------|------|:-------:|-------------|
| `-rules.lowercase` | bool | `true` | Check for uppercase first letter |
| `-rules.english` | bool | `true` | Check for non-ASCII characters |
| `-rules.special-chars` | bool | `true` | Check for special characters and emoji |
| `-rules.sensitive` | bool | `true` | Check for sensitive data keywords |
| `-sensitive-keywords` | string | built-in list | Comma-separated list of sensitive keywords |
| `-custom-patterns` | string | — | Comma-separated list of custom regexp patterns |

## Example with flags

```bash
go vet -vettool=./golanglog-linter \
  -golanglog-linter.rules.sensitive=true \
  -golanglog-linter.sensitive-keywords="password,token,mykey" \
  ./...
```

> Full documentation — see [root README](../../README.md).
