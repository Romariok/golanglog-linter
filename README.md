# golanglog-linter — Go-линтер для валидации лог-сообщений

Статический анализатор для Go, проверяющий вызовы функций логирования на соответствие правилам стиля и безопасности. Работает как самостоятельный инструмент через `go vet` или как модульный плагин `golangci-lint`.

Четыре правила: строчная первая буква, только ASCII, никаких спецсимволов и эмодзи, никаких чувствительных данных. Для правила 1 есть автоисправление (`golangci-lint run --fix`). Ключевые слова и regexp-паттерны настраиваются через `.golangci.yml` или флаги CLI.

Детектирование лог-вызовов работает через `pass.TypesInfo` — разрешается фактический тип получателя, а не просто имя. `mypackage.Info("msg")` не будет принят за `slog.Info`. Цепочки `Sugar()` тоже поддерживаются.

---

## Поддерживаемые логгеры

| Пакет | Методы |
|-------|--------|
| `log/slog` | `Debug`, `Info`, `Warn`, `Error`, `Log`, `DebugContext`, `InfoContext`, `WarnContext`, `ErrorContext` |
| `go.uber.org/zap` · Logger | `Debug`, `Info`, `Warn`, `Error`, `DPanic`, `Panic`, `Fatal` |
| `go.uber.org/zap` · SugaredLogger | `Debug`, `Info`, `Warn`, `Error`, `Debugf`, `Infof`, `Warnf`, `Errorf` — через цепочку `Sugar()` |

---

## Быстрый старт

> Другие варианты запуска (go vet, CLI флаги) — в [`docs/QUICKSTART.md`](docs/QUICKSTART.md).

### Требования

- Go 1.22+
- golangci-lint v2.x

### Module-плагин для `golangci-lint`

Создать `.custom-gcl.yml` в корне проекта:

```yaml
version: v2.11.3
plugins:
  - module: 'github.com/romariok/golanglog-linter'
    import: 'github.com/romariok/golanglog-linter/plugin'
    version: latest
```

Собрать кастомный бинарник:

```bash
golangci-lint custom           # создаёт ./custom-gcl
```

Добавить в `.golangci.yml`:

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

Запустить:

```bash
./custom-gcl run               # проверка
./custom-gcl run --fix         # проверка + автоисправление правила 1
```

---

## Правила

| # | Правило | Диагностическое сообщение | Автоисправление |
|:-:|---------|--------------------------|:---------------:|
| 1 | Первая буква должна быть строчной | `log message should start with a lowercase letter` | Да |
| 2 | Только английский (ASCII) | `log message must be in English only` | Нет |
| 3 | Запрет спецсимволов и эмодзи | `log message must not contain special characters or emojis` | Нет |
| 4 | Запрет чувствительных данных | `log message may contain sensitive data: found keyword "%s"` | Нет |

---

## Архитектура

```
analysis.Pass (AST + TypesInfo)
         │
         ▼
   analyzer.go          ← точка входа, подключает правила, обходит узлы CallExpr
         │
   logcall.go           ← IsLogCall(): определение вызовов по типу через pass.TypesInfo
         │
   ┌─────┴──────────────────────────┐
   │                                │
rules/lowercase.go     rules/english.go
rules/special_chars.go rules/sensitive.go
```

`logcall.go` разрешает тип получателя и путь пакета — не просто имя функции — чтобы отличить `slog.Info` от `mypackage.Info`. Цепочки `Sugar()` разрешаются обходом выражения-селектора и проверкой возвращаемого типа через `pass.TypesInfo`.

Каждое правило получает `*analysis.Pass`, `*ast.CallExpr` и уже извлечённую строку сообщения. Правила — чистые функции без общего состояния и побочных эффектов.

Конфигурация обрабатывается через `analysis.Analyzer.Flags` в режиме CLI и проецируется на тот же struct `Config`, что и в режиме плагина golangci-lint.

---

## Структура проекта

```
golanglog-linter/
├── cmd/golanglog-linter/
│   └── main.go              # самостоятельный бинарник (go vet -vettool)
├── pkg/golanglog/
│   ├── analyzer.go          # точка входа *analysis.Analyzer
│   ├── rules/
│   │   ├── logcall.go       # определение лог-вызовов через TypesInfo
│   │   ├── lowercase.go     # правило 1: заглавная буква + SuggestedFix
│   │   ├── english.go       # правило 2: обнаружение non-ASCII символов
│   │   ├── special_chars.go # правило 3: эмодзи, повторяющиеся/завершающие спецсимволы, \n
│   │   └── sensitive.go     # правило 4: совпадение по ключевым словам и именам переменных
│   └── config/
│       └── config.go        # struct Config, флаги, значения по умолчанию
├── plugin/
│   └── plugin.go            # точка входа модульного плагина golangci-lint
├── testdata/src/
│   ├── go.uber.org/zap/     # заглушка zap для testdata в режиме GOPATH
│   ├── lowercase/           # bad.go + good.go для правила 1
│   ├── english/             # bad.go + good.go для правила 2
│   ├── special/             # bad.go + good.go для правила 3
│   └── sensitive/           # bad.go + good.go для правила 4
├── .golangci.yml
├── .github/workflows/ci.yml
├── go.mod
└── README.md
```

---

## Тестирование

Для каждого правила — отдельный `*_test.go` и два файла с тестовыми данными.

```bash
go test ./...          # запустить все тесты
go test -cover ./...   # с отчётом о покрытии
go test -race ./...    # с детектором гонок
```

Файлы testdata используют аннотации `// want` для ожидаемых диагностик:

```go
// bad.go
slog.Info("Starting server") // want `log message should start with a lowercase letter`

// good.go
slog.Info("starting server")  // аннотация отсутствует — диагностика не ожидается
```

### Целевое покрытие

| Пакет | Цель |
|-------|:----:|
| `pkg/golanglog/rules` | ≥ 80% |

### Покрытые граничные случаи

| Сценарий | Результат |
|----------|:---------:|
| Не-лог вызов с таким же именем метода (`myCache.Info(...)`) | нет ложного срабатывания |
| Аргумент-переменная (`slog.Info(msg)`) | пропускается всеми правилами |
| Конкатенация литералов (`"Error: " + err.Error()`) | проверяется первый сегмент-литерал |
| Пустая строка `""` | нет диагностики |
| `slog.Info()` без аргументов | нет паники |
| Цепочка `logger.Sugar().Info(...)` | корректно определяется |
| Все четыре правила на одном вызове | четыре диагностики |

---

## CI/CD

GitHub Actions (`.github/workflows/ci.yml`):

| Задача | Команда | Триггер |
|--------|---------|:-------:|
| `build` | `go build ./...` | push |
| `test` | `go test ./... -race -cover` | push |
| `lint` | `golangci-lint run` | push |
| `release` | `goreleaser` | тег `v*` |

---

## Заметки по реализации

Первым написал `logcall.go` — до любых правил. `cache.Info("msg")` и `slog.Info("msg")` выглядят одинаково в AST, различить их можно только через `pass.TypesInfo`, и лучше решить это один раз централизованно.

Правила писались от простого к сложному. Lowercase — чтобы проверить весь пайплайн от обхода AST до `SuggestedFix`. English — тривиальная проверка `rune > 127`. Special chars потребовал аккуратной работы с regex: версионные строки (`v1.2.3`), паттерны `key: value` и Unix-пути нужно было явно исключить. Sensitive оказался самым нетривиальным: граничные слова `\b` в regex (чтобы `passthrough` не срабатывало на `password`), рекурсивный обход бинарных выражений, проверка полей селекторов для `req.Password`.

Несколько моментов, которые стоит знать:

- `SuggestedFix` для правила 1 использует `utf8.RuneLen(firstRune)` — иначе многобайтовый Unicode-символ в верхнем регистре заменится некорректно.
- `analysistest` работает в режиме GOPATH без разрешения модулей, поэтому для тестов zap нужна минимальная заглушка в `testdata/src/go.uber.org/zap/`.
- `!` запрещён в любой позиции, не только в конце — в структурированных логах для систем агрегации он неуместен везде.

---

## Документация

Подробная спецификация проекта — правила, тест-кейсы, архитектурные решения — находится в [`docs/SPEC.md`](docs/SPEC.md).

---

## Использование ИИ

ИИ использовался для ускорения написания кода и документации. Все решения принимались вручную, каждая часть логики проверялась.
