# golanglog-linter — Quick Start

## Requirements

- Go 1.22+
- golangci-lint v2.x (for plugin integration)

---

## Option 1 — Standalone tool (`go vet`)

Install the binary and run:

```bash
go install github.com/romariok/golanglog-linter/cmd/golanglog-linter@latest
go vet -vettool=$(which golanglog-linter) ./...
```

Or build locally from the repository:

```bash
go build -o golanglog-linter ./cmd/golanglog-linter/
go vet -vettool=./golanglog-linter ./...
```

With all flags explicitly set:

```bash
go vet -vettool=$(which golanglog-linter) \
  -golanglog-linter.rules.lowercase=true \
  -golanglog-linter.rules.english=true \
  -golanglog-linter.rules.special-chars=true \
  -golanglog-linter.rules.sensitive=true \
  -golanglog-linter.sensitive-keywords="password,token,secret,api_key" \
  -golanglog-linter.custom-patterns="credit.?card,ssn" \
  ./...
```

---

## Option 2 — Module plugin for `golangci-lint` (recommended)

Create `.custom-gcl.yml` in the project root:

```yaml
version: v2.11.3
plugins:
  - module: 'github.com/romariok/golanglog-linter'
    import: 'github.com/romariok/golanglog-linter/plugin'
    version: latest
```

Build the custom binary:

```bash
golangci-lint custom           # produces ./custom-gcl
```

Add to `.golangci.yml`:

```yaml
version: "2"

linters:
  enable:
    - golanglog
  settings:
    custom:
      golanglog:
        type: "module"
        description: "Validates log message style and security"
        settings:
          rules:
            lowercase: true
            english: true
            special-chars: true
            sensitive: true
          sensitive-keywords:
            - password
            - token
            - secret
            - api_key
          custom-patterns:
            - "credit.?card"
            - "ssn"
```

Run:

```bash
./custom-gcl run               # lint
./custom-gcl run --fix         # lint + auto-fix rule 1
```

---

## Option 3 — CLI flags (direct binary invocation)

```bash
# install
go install github.com/romariok/golanglog-linter/cmd/golanglog-linter@latest

# run with explicit flags
golanglog-linter \
  -rules.lowercase=true \
  -rules.english=true \
  -rules.special-chars=true \
  -rules.sensitive=true \
  -sensitive-keywords="password,token,secret,mytoken" \
  -custom-patterns="credit.?card,ssn" \
  ./...
```
