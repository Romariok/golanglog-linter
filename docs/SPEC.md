# PRD: golanglog-linter — Go Linter for Log Message Validation

---

## 1. Introduction / Overview

**golanglog-linter** — линтер для языка Go, анализирующий вызовы функций логирования и проверяющий соответствие сообщений установленным правилам качества и безопасности.

### Проблема

В больших Go-проектах лог-сообщения часто написаны непоследовательно — с заглавных букв, на русском языке, с эмодзи и спецсимволами, а иногда содержат чувствительные данные (пароли, токены). Ручная ревью такого кода трудоёмка и ненадёжна.

### Цель

Автоматически обнаруживать нарушения стиля и безопасности в лог-вызовах на этапе CI/CD, не требуя изменений в рабочем процессе разработчиков — линтер интегрируется в `golangci-lint`.

---

## 2. Goals

1. Реализовать статический анализатор, проверяющий **4 обязательных правила** для лог-сообщений.
2. Поддержать два логгера: `log/slog` (stdlib) и `go.uber.org/zap`.
3. Интегрировать анализатор как **module plugin** для golangci-lint v2+.
4. Обеспечить покрытие **unit-тестами** каждого правила, включая граничные случаи.
5. Реализовать все **4 бонусных задания**: конфигурацию, авто-исправление, CI/CD и кастомные паттерны.
6. Предоставить документацию и примеры использования.

---

## 3. User Stories

> Ключевые сценарии использования линтера разными участниками команды.

- **Как разработчик** я хочу получать предупреждение при запуске `golangci-lint run`, если моё лог-сообщение начинается с заглавной буквы, чтобы поддерживать единый стиль в команде.
- **Как тимлид** я хочу, чтобы CI/CD падал при обнаружении чувствительных данных в логах (паролей, токенов), чтобы предотвратить утечки.
- **Как DevOps-инженер** я хочу настроить список запрещённых ключевых слов через `.golangci.yml`, не изменяя исходный код линтера.
- **Как разработчик** я хочу, чтобы линтер предложил авто-исправление для правила об uppercase, чтобы применить его одной командой.
- **Как контрибьютор** я хочу добавить кастомный regexp-паттерн для чувствительных данных через конфиг.

---

## 4. Functional Requirements

### 4.1 Правило 1 — Строчная буква в начале сообщения

1. Линтер **должен** обнаруживать вызовы лог-функций, где первый аргумент (строковый литерал) начинается с заглавной буквы (Unicode uppercase).
2. Диагностическое сообщение: `log message should start with a lowercase letter`.
3. Линтер **должен** предоставлять `analysis.SuggestedFix` — автоматически заменять первый символ на его lowercase-вариант.
4. Правило **не должно** срабатывать на сообщениях, начинающихся со строчной буквы, цифры, спецсимвола или пустой строке.
5. Правило **не должно** применяться к аргументам-переменным (только к строковым литералам).

---

### 4.2 Правило 2 — Только английский язык

1. Линтер **должен** обнаруживать строковые литералы, содержащие символы за пределами ASCII-диапазона, используя `unicode.RangeTable`.
2. Диагностическое сообщение: `log message must be in English only`.
3. Правило проверяет наличие хотя бы одного non-ASCII символа (`rune > 127`) в строке.
4. Цифры и ASCII-пунктуация (`. , : ; - _ ( ) [ ] /` и т.д.) **разрешены** — проверяются только non-ASCII символы.
5. Правило **не** предлагает авто-исправление (перевод невозможен автоматически).

---

### 4.3 Правило 3 — Запрет спецсимволов и эмодзи

1. Линтер **должен** обнаруживать строки, содержащие:
   - Эмодзи (Unicode range `\x{1F300}`–`\x{1FAFF}`, `\x{2600}`–`\x{26FF}`, `\x{2700}`–`\x{27BF}` и смежные блоки).
   - Повторяющиеся спецсимволы (`.`, `!`, `?` более одного раза подряд).
   - Отдельные спецсимволы в конце строки: `!`, `?`, `:`, `...`.
   - Одиночный `!` **в любой позиции** строки (включая середину).
2. Диагностическое сообщение: `log message must not contain special characters or emojis`.
3. Правило **не** предлагает авто-исправление.

---

### 4.4 Правило 4 — Запрет чувствительных данных

1. Линтер **должен** срабатывать в двух случаях:

   - **Строковый литерал-маркер:** первый аргумент лог-вызова содержит одно из ключевых слов: `password`, `passwd`, `secret`, `token`, `api_key`, `apikey`, `auth`, `credential`, `private_key` — в любом регистре и с разделителями (`:`, `=`, пробел). Совпадение проверяется по **границам слова** (`\b`), чтобы избежать ложных срабатываний на `passthrough`, `authorize` и т.п.

     _Пример:_ `log.Info("user password: " + password)` → маркер `"user password: "`.

   - **Переменная-аргумент:** любой аргумент лог-вызова называется именем из списка чувствительных слов. Совпадение также проверяется по **границам слова** (`\b`).

     _Пример:_ `log.Debug("api_key=" + apiKey)` → имя переменной `apiKey` содержит `key`.

2. Диагностическое сообщение: `log message may contain sensitive data: found keyword "%s"`.
3. Список ключевых слов **должен** быть расширяем через конфигурацию.

---

### 4.5 Поддерживаемые логгеры

| Пакет             | Методы                                                                                                                                          |
| ----------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- |
| `log/slog`        | `Debug`, `Info`, `Warn`, `Error`, `Log`, `DebugContext`, `InfoContext`, `WarnContext`, `ErrorContext`                                            |
| `go.uber.org/zap` | `Debug`, `Info`, `Warn`, `Error`, `DPanic`, `Panic`, `Fatal`                                                                                    |
| `go.uber.org/zap` _(Sugar)_ | `Sugar().Debug`, `Sugar().Info`, `Sugar().Warn`, `Sugar().Error`, `Sugar().Infof`, `Sugar().Debugf`, `Sugar().Warnf`, `Sugar().Errorf` |

> **Важно:** Линтер **не должен** анализировать методы других пакетов (чтобы избежать ложных срабатываний). Первый аргумент (message string) — всегда первый позиционный аргумент метода.
>
> Для zap Sugar поддерживается анализ **цепочки вызовов** `logger.Sugar().Method(...)` с разрешением типа через `pass.TypesInfo`.

---

### 4.6 Структура проекта

```
golanglog-linter/
├── cmd/
│   └── golanglog-linter/   # standalone binary (go vet -vettool)
│       └── main.go
├── pkg/
│   └── golanglog/
│       ├── analyzer.go     # *analysis.Analyzer entry point
│       ├── rules/
│       │   ├── lowercase.go        # Правило 1
│       │   ├── english.go          # Правило 2
│       │   ├── special_chars.go    # Правило 3
│       │   ├── sensitive.go        # Правило 4
│       │   └── logcall.go          # Общая логика определения лог-вызовов
│       └── config/
│           └── config.go   # Структура конфигурации
├── plugin/
│   └── plugin.go           # golangci-lint module plugin entry
├── testdata/
│   └── src/
│       ├── lowercase/      # тестовые файлы для правила 1
│       ├── english/        # тестовые файлы для правила 2
│       ├── special/        # тестовые файлы для правила 3
│       └── sensitive/      # тестовые файлы для правила 4
├── .golangci.yml
├── .github/
│   └── workflows/
│       └── ci.yml
├── go.mod
├── go.sum
└── README.md
```

---

## 5. Non-Goals _(Out of Scope)_

| Что не входит в scope | Причина |
| --- | --- |
| Анализ аргументов-переменных (не литералов) | Слишком сложно без data-flow анализа |
| Поддержка logrus, zerolog и других логгеров | Только slog и zap |
| Анализ `fmt.Sprintf`/`fmt.Fprintf` вне контекста лога | Высокий риск ложных срабатываний |
| Авто-перевод сообщений на английский | Технически невозможно статически |
| Изменение поведения программы | Только диагностика и `SuggestedFix` |

---

## 6. Design Considerations

### Архитектура анализатора

- Использовать `golang.org/x/tools/go/analysis` и паттерн `*analysis.Analyzer`.
- Центральный анализатор (`analyzer.go`) использует `analysis.Pass` для обхода AST.
- Для определения лог-вызовов использовать `go/types` — проверять тип получателя вызова (`CallExpr`) и сопоставлять с известными пакетами.
- Каждое правило реализуется как отдельная функция `checkXxx(pass *analysis.Pass, call *ast.CallExpr, msg string)`.

### Конфигурация через golangci-lint

```yaml
linters-settings:
  custom:
    golanglog:
      type: module
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
          - custom_pattern
        custom-patterns:
          - "credit.?card"
          - "ssn"
```

### SuggestedFix для Правила 1

```go
analysis.SuggestedFix{
    Message: "Convert first letter to lowercase",
    TextEdits: []analysis.TextEdit{{
        Pos:     firstCharPos,
        End:     firstCharPos + token.Pos(utf8.RuneLen(firstRune)),
        NewText: []byte(string(unicode.ToLower(firstRune))),
    }},
}
```

---

## 7. Technical Considerations

### Зависимости

| Пакет | Версия | Назначение |
| --- | --- | --- |
| `golang.org/x/tools` | `v0.20.0+` | `go/analysis` framework |
| `go.uber.org/zap` | latest | тестовые данные и type-checking |
| `golangci-lint` | `v2.x` | module plugin интеграция |

### Версия Go

> Требуется **Go 1.22+** (требование задания).

### golangci-lint Module Plugin

Файл `plugin/plugin.go` должен экспортировать интерфейс `AnalyzerPlugin`:

```go
package main

import (
    "github.com/user/golanglog-linter/pkg/golanglog"
    "golang.org/x/tools/go/analysis"
)

type analyzerPlugin struct{}

func (*analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
    return []*analysis.Analyzer{golanglog.Analyzer}
}

var AnalyzerPlugin analyzerPlugin
```

### Определение лог-вызовов

> Для корректного определения вызовов необходимо использовать `pass.TypesInfo` для разрешения типов, а **не** просто имён функций, чтобы не ловить `mypackage.Info("...")` как лог-вызов.

---

## 8. Testing Requirements

### 8.1 Общие требования к тестам

- Все unit-тесты писать с использованием `analysistest.Run` из `golang.org/x/tools/go/analysis/analysistest`.
- Тестовые файлы помещать в `testdata/src/<rule>/`.
- В тестовых файлах использовать `// want` комментарии для указания ожидаемых диагностик.
- Достичь покрытия **не менее 80%** для каждого модуля правил.

---

### 8.2 Тест-кейсы — Правило 1 (Lowercase)

| #   | Сценарий                   | Вход                                 | Ожидание              |
| --- | -------------------------- | ------------------------------------ | --------------------- |
| 1   | Заглавная буква ASCII      | `slog.Info("Starting server")`       | **FAIL**              |
| 2   | Строчная буква ASCII       | `slog.Info("starting server")`       | PASS                  |
| 3   | Первый символ — цифра      | `slog.Info("3 retries left")`        | PASS                  |
| 4   | Первый символ — спецсимвол | `slog.Info("!error")`                | PASS _(rule 3 catches)_ |
| 5   | Пустая строка              | `slog.Info("")`                      | PASS                  |
| 6   | Переменная-аргумент        | `slog.Info(msg)`                     | PASS _(не литерал)_   |
| 7   | Заглавная unicode-буква    | `slog.Info("Ñew server")`            | **FAIL**              |
| 8   | Конкатенация с заглавной   | `slog.Info("Error: " + err.Error())` | **FAIL**              |
| 9   | zap.Info с заглавной       | `logger.Info("Starting")`            | **FAIL**              |
| 10  | SuggestedFix корректен     | `slog.Info("Hello")` → `"hello"`     | Fix верен             |

---

### 8.3 Тест-кейсы — Правило 2 (English only)

| #   | Сценарий                   | Вход                          | Ожидание  |
| --- | -------------------------- | ----------------------------- | --------- |
| 1   | Кириллица                  | `log.Info("запуск сервера")`  | **FAIL**  |
| 2   | Китайские иероглифы        | `log.Info("服务器启动")`       | **FAIL**  |
| 3   | Только ASCII               | `log.Info("server started")`  | PASS      |
| 4   | Смешанный (кирилл. + лат.) | `log.Info("server запущен")`  | **FAIL**  |
| 5   | Латиница с акцентами       | `log.Info("café connection")` | **FAIL**  |
| 6   | Числа и ASCII знаки        | `log.Info("port 8080 ready")` | PASS      |
| 7   | Пустая строка              | `log.Info("")`                | PASS      |
| 8   | Переменная                 | `log.Info(msg)`               | PASS      |
| 9   | Арабские символы           | `log.Info("خطأ")`             | **FAIL**  |
| 10  | Японские символы           | `log.Info("エラー")`          | **FAIL**  |

---

### 8.4 Тест-кейсы — Правило 3 (Special chars & emoji)

| #   | Сценарий                    | Вход                                  | Ожидание     |
| --- | --------------------------- | ------------------------------------- | ------------ |
| 1   | Эмодзи в конце              | `log.Info("server started 🚀")`       | **FAIL**     |
| 2   | Множество `!`               | `log.Error("connection failed!!!")`   | **FAIL**     |
| 3   | Многоточие                  | `log.Warn("something went wrong...")` | **FAIL**     |
| 4   | Один `!` в конце            | `log.Error("connection failed!")`     | **FAIL**     |
| 5   | Чистое сообщение            | `log.Info("server started")`          | PASS         |
| 6   | Двоеточие в середине        | `log.Info("key: value set")`          | PASS         |
| 7   | Вопросительный знак в конце | `log.Warn("retry failed?")`           | **FAIL**     |
| 8   | Пунктуация в числах         | `log.Info("v1.2.3 deployed")`         | PASS         |
| 9   | Эмодзи в начале             | `log.Info("🔥 hot reload")`           | **FAIL**     |
| 10  | Несколько эмодзи            | `log.Info("done ✅🎉")`               | **FAIL**     |
| 11  | Спецсимвол `\n` в строке    | `log.Info("line1\nline2")`            | **FAIL**     |

---

### 8.5 Тест-кейсы — Правило 4 (Sensitive data)

| #   | Сценарий                        | Вход                                           | Ожидание                      |
| --- | ------------------------------- | ---------------------------------------------- | ----------------------------- |
| 1   | Ключевое слово в литерале       | `log.Info("user password: " + p)`              | **FAIL** (`password`)         |
| 2   | Имя переменной `password`       | `log.Debug("auth", zap.String("password", p))` | **FAIL**                      |
| 3   | Имя переменной `apiKey`         | `log.Debug("api_key=" + apiKey)`               | **FAIL**                      |
| 4   | Токен в строке                  | `log.Info("token: " + token)`                  | **FAIL**                      |
| 5   | Нейтральное сообщение           | `log.Info("user authenticated")`               | PASS                          |
| 6   | `SECRET` в верхнем регистре     | `log.Info("SECRET=" + val)`                    | **FAIL**                      |
| 7   | Частичное совпадение            | `log.Info("password_reset complete")`          | **FAIL**                      |
| 8   | Кастомный паттерн из конфига    | `log.Info("credit_card: " + cc)`               | **FAIL** _(если настроен)_    |
| 9   | Поле структуры `req.Password`   | `log.Info("req", req.Password)`                | **FAIL**                      |
| 10  | Слово `passthrough`             | `log.Info("passthrough enabled")`              | PASS _(не `password`)_        |
| 11  | Слово в `slog.String` key       | `slog.Info("auth", slog.String("token", t))`   | **FAIL**                      |
| 12  | `credential` в строке           | `log.Info("credential stored")`                | **FAIL**                      |

---

### 8.6 Граничные и негативные сценарии (edge cases)

| #   | Сценарий                                | Ожидание                                                                             |
| --- | --------------------------------------- | ------------------------------------------------------------------------------------ |
| 1   | Не-лог вызов с таким же именем метода   | `myCache.Info("Starting")` — **PASS** (чужой пакет)                                  |
| 2   | Вложенный лог-вызов                     | `slog.Info(fmt.Sprintf("val: %s", v))` — PASS для правила 1                          |
| 3   | Многострочный литерал                   | `slog.Info("line1\n" + "line2")` — конкатенация анализируется                        |
| 4   | Нулевой аргумент                        | `slog.Info()` — линтер **не паникует**                                               |
| 5   | Слишком длинное сообщение (>1000 chars) | Корректная обработка без OOM                                                         |
| 6   | Файл без импорта slog/zap               | Линтер не анализирует, нет ложных срабатываний                                       |
| 7   | Все правила одновременно                | `slog.Info("Пароль: 🔑" + password)` — **4 диагностики**                             |
| 8   | `//nolint:golanglog` директива          | Нарушение игнорируется                                                               |
| 9   | Сообщение только из пробелов            | `slog.Info(" ")` — PASS (нет буквы)                                                  |
| 10  | Unicode uppercase İ (Turkish)           | `slog.Info("İstanbul")` — **FAIL** (заглавная)                                       |

---

### 8.7 Требования к тестовым файлам в `testdata/`

- Каждый файл должен компилироваться (`go build ./testdata/...`).
- Аннотации вида `// want "log message should start"` для каждого ожидаемого FAIL.
- Структура файлов:
  - `testdata/src/<rule>/good.go` — только PASS-случаи
  - `testdata/src/<rule>/bad.go` — только FAIL-случаи

---

## 9. Bonus Tasks

### 9.1 Конфигурация через golangci-lint settings

- Структура конфига парсится из `pass.Analyzer.Flags` или через `golangci-lint` settings map.
- Поля: `rules` (`map[string]bool`), `sensitive-keywords` (`[]string`).
- Значения по умолчанию: все правила включены, стандартный список keywords.

### 9.2 Авто-исправление (SuggestedFixes)

- Реализовать для **Правила 1** (lowercase first letter).
- `analysis.SuggestedFix` с точным указанием позиции первого rune в AST.
- Применяется через `golangci-lint run --fix` или `go vet`.

### 9.3 Кастомные паттерны

- Пользователь может добавить regexp-паттерны в конфиг (`custom-patterns: ["credit.?card", "ssn"]`).
- Паттерны компилируются **один раз** при инициализации анализатора.
- При ошибке компиляции паттерна — вернуть понятное сообщение об ошибке.

### 9.4 CI/CD (GitHub Actions)

Файл `.github/workflows/ci.yml`:

```yaml
jobs:
  build:   # go build ./...
  test:    # go test ./... -race -cover
  lint:    # golangci-lint run на самом проекте
  release: # goreleaser при теге v*
```

---

## 10. Success Metrics

| Метрика | Критерий успеха |
| --- | --- |
| Все правила реализованы | `go test ./...` — зелёный |
| Покрытие тестами | ≥ **80%** для пакетов `rules/` |
| Интеграция с golangci-lint | Подключается через `type: module`, обнаруживает нарушения |
| Авто-исправление | `golangci-lint run --fix` применяет fix для Правила 1 |
| CI/CD | Pipeline проходит при push в `main` |
| Документация | README содержит рабочий пример установки и использования |

---

## 11. Decisions Log

| # | Тема | Решение |
| --- | --- | --- |
| 1 | **Правило 3:** одиночный `!` в середине сообщения | **Считать нарушением** — `!` запрещён в любой позиции строки |
| 2 | **Правило 4:** совпадение по границе слова vs вхождение | **Word boundary (`\b`)** — `passthrough` не является нарушением, `password` — является |
| 3 | **Правило 2:** разрешить цифры и ASCII-пунктуацию | **Разрешить** — запрещены только символы с `rune > 127` |
| 4 | **golangci-lint:** поддержка v1.x как fallback | **Да** — реализовать поддержку plugin binary для golangci-lint v1.x |
| 5 | **zap Sugar:** поддержка цепочки `Sugar().Info(...)` | **Да, в первом релизе** — анализировать через разрешение типа в `pass.TypesInfo` |

---

## 12. Implementation Phases

### Этап 1 — Базовая структура

- Инициализировать `go.mod`, настроить зависимости.
- Создать `pkg/golanglog/analyzer.go` с `*analysis.Analyzer`.
- Реализовать `logcall.go` — определение лог-вызовов по типу через `pass.TypesInfo`.

### Этап 2 — Реализация правил

- `rules/lowercase.go` — Правило 1 + SuggestedFix.
- `rules/english.go` — Правило 2 (unicode ranges).
- `rules/special_chars.go` — Правило 3 (regexp + rune ranges).
- `rules/sensitive.go` — Правило 4 (keyword matching).

### Этап 3 — Тестирование

- Написать unit-тесты через `analysistest.Run` для каждого правила.
- Создать файлы в `testdata/src/` с аннотациями `// want`.
- Покрыть все edge cases из раздела 8.

### Этап 4 — Интеграция с golangci-lint

- Реализовать `plugin/plugin.go` с `AnalyzerPlugin`.
- Добавить конфигурацию через `Flags` / settings map.
- Написать `.golangci.yml` с примером подключения.
- Настроить GitHub Actions CI/CD.
- Написать README с инструкцией по установке и примерами.

### Этап 5 — Проверка на реальных проектах

- Запустить линтер на 2–3 реальных open-source Go-проектах (например, из GitHub).
- Зафиксировать найденные нарушения и ложные срабатывания.
- Оценить производительность и стабильность анализатора на больших кодовых базах.
- Внести необходимые исправления по итогам проверки.
- Результаты (примеры найденных нарушений, метрики, выводы) опубликовать в README в разделе **Real-world Testing**.
