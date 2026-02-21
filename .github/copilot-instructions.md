# CRSLang – Copilot Agent Onboarding Instructions

## Project Overview

**CRSLang** is a Go library and CLI tool that translates [OWASP Core Rule Set (CRS)](https://coreruleset.org/) rules between their native [SecLang/ModSecurity](https://github.com/coreruleset/seclang_parser) configuration format (`.conf` files) and a new, language-independent YAML representation called **CRSLang**.

The project also compiles to WebAssembly (WASM) so the translation logic can run in a browser without any server-side component.

---

## Repository Layout

```
crslang/
├── main.go                    # CLI entry point (seclang ↔ CRSLang conversion)
├── mage.go / magefile.go      # Mage build targets (generate, build, test)
├── Makefile                   # Short-cut targets: build, test, wasm, clean
├── go.mod / go.sum            # Go module definition (module: github.com/coreruleset/crslang)
│
├── listener/                  # ANTLR listener that walks the parse tree and
│   │                          # produces internal types.* structs
│   └── extended_seclang_parser_listener.go
│       (+ sec_rule.go, sec_action.go, variables.go, actions.go, …)
│
├── translator/                # Orchestration layer
│   ├── seclang.go             # LoadSeclang / LoadSeclangFromString / PrintSeclang
│   ├── crslang.go             # ToCRSLang (seclang → CRSLang conversion)
│   ├── file.go                # writeToFile helper
│   └── translator_test.go     # Round-trip and unit tests
│
├── types/                     # Core CRSLang type definitions + YAML serialisation
│   ├── configuration.go       # ConfigurationList, DirectiveList, ExtractDefaultValues
│   ├── secrule.go             # SecRule (variables, operator, actions, chaining)
│   ├── secaction.go           # SecAction
│   ├── metadata.go            # SecRuleMetadata, CommentMetadata, etc.
│   ├── operators.go           # Operator type + all operator constants
│   ├── variables.go           # Variable type + all variable-name constants
│   ├── collections.go         # Collection type + all collection-name constants
│   ├── actions.go             # Action type + all action constants
│   ├── transformations.go     # Transformation type + all transformation constants
│   ├── condition_directives.go# RuleWithCondition (CRSLang rule representation)
│   └── *_test.go              # Per-file unit tests
│
├── wasm/
│   ├── main.go                # WASM build (js/wasm build tag)
│   │                          # Exports: seclangToCRSLang(), crslangToSeclang()
│   └── demo.html              # Browser demo page
│
├── base_test.go               # Lexer/parser smoke tests (TestSecLang, TestCRSLang, …)
├── listener_test.go           # Full listener round-trip tests (TestLoadSecLang, …)
│
└── testdata/                  # Test fixture .conf files
    ├── crs/                   # Subset of OWASP CRS rule files
    ├── plugins/               # Sample plugin configuration files
    └── test*.conf             # Individual unit-test fixtures
```

---

## Key Concepts

### Data Flow

```
.conf file(s)
    └─► ANTLR SecLang lexer/parser  (github.com/coreruleset/seclang_parser)
            └─► ExtendedSeclangParserListener  (listener/)
                    └─► types.ConfigurationList  (raw Seclang structs)
                            └─► translator.ToCRSLang()
                                    └─► types.ConfigurationList  (CRSLang structs)
                                            └─► YAML marshal  →  crslang.yaml
```

Reverse path:

```
crslang.yaml
    └─► types.LoadDirectivesWithConditionsFromFile()
            └─► types.ConfigurationList  (CRSLang structs)
                    └─► types.FromCRSLangToUnformattedDirectives()
                            └─► DirectiveList.ToSeclang()  →  .conf files
```

### Important Types

| Type | Package | Description |
|---|---|---|
| `ConfigurationList` | `types` | Top-level container; holds `Global` defaults + `[]DirectiveList` |
| `DirectiveList` | `types` | A group of directives (typically one `.conf` file); has an `Id` and optional `Marker` |
| `SecRule` | `types` | A parsed `SecRule` directive; Variables/Collections, Operator, Transformations, Actions, optional ChainedRule |
| `SecAction` | `types` | A parsed `SecAction` directive |
| `RuleWithCondition` | `types` | CRSLang representation of a rule with structured conditions |
| `Operator` | `types` | Operator name + value + optional negation flag |
| `SeclangActions` | `types` | Disruptive, Non-disruptive and Flow actions |

---

## Build & Run

### Prerequisites

- Go 1.22+ (`go.mod` declares `go 1.22.2`)
- No code-generation step is required; the ANTLR-generated parser is vendored via `github.com/coreruleset/seclang_parser`.

### Build the CLI

```bash
go build .
# Produces: ./crslang
```

### Run the CLI

```bash
# Convert a directory of CRS .conf files to CRSLang YAML
./crslang -o output_name path/to/rules/

# Convert CRSLang YAML back to SecLang .conf files
./crslang -s -o output_dir/ crslang.yaml
```

### Build the WASM module

```bash
make wasm
# Produces: wasm/crslang.wasm  and  wasm/wasm_exec.js
```

### Serve the browser demo

```bash
cd wasm
python3 -m http.server 8080
# Open http://localhost:8080/demo.html
```

---

## Testing

### Run all tests

```bash
go test ./...
```

### Run tests with verbose output

```bash
go test -v ./...
```

### Run a specific test or package

```bash
go test -v -run TestLoadSecLang .
go test -v ./translator/
go test -v ./types/
```

### Test structure

- **`base_test.go`** — lexer/parser smoke tests over all files in `testdata/`
- **`listener_test.go`** — `TestLoadSecLang`: table-driven tests for the ANTLR listener. Each case provides a raw SecLang payload string and an expected `types.ConfigurationList` value.
- **`translator/translator_test.go`** — round-trip tests (Seclang → CRSLang YAML → reload → compare)
- **`types/*_test.go`** — unit tests for individual type helpers

### Adding new tests

Follow the table-driven pattern used in `listener_test.go`:

```go
{
    name:    "MyNewCase",
    payload: `SecRule ... "id:12345,phase:1,pass"`,
    expected: types.ConfigurationList{ /* ... */ },
},
```

For `types` tests use `testify/require`:

```go
require.Equal(t, expected, actual)
```

---

## CI / Workflows

### `go.yml` — Main test workflow

Triggers on push/PR to `main`. Runs `go test -v ./...`.

### `check_format.yml` — End-to-end format check

1. Clones the OWASP CRS repository at a pinned commit.
2. Converts all CRS rules via `go run . -o output crs/rules/`.
3. Converts the CRSLang output back to SecLang: `go run . -s -o ... output.yaml`.
4. Runs the CRS Python rules-check script to validate output conformance.

### Known CI Issues

- Both workflow files currently specify **`go-version: '1.25.6'`** which does not exist (likely a typo; the project requires `go 1.22+`). If these workflows fail with a "version not found" error, the fix is to change the `go-version` value to a valid Go release such as `'1.22.x'` or `'1.23.x'`.
- The `check_format.yml` workflow uses **`python-version: "3.14"`** (pre-release at time of writing). If that version is unavailable, switch to `"3.12"` or `"3.13"`.

---

## Code Conventions

- **Go module path**: `github.com/coreruleset/crslang`
- All packages live under the module root; import paths look like `github.com/coreruleset/crslang/types`.
- YAML serialization uses `go.yaml.in/yaml/v4` (not the more common `gopkg.in/yaml.v3`). Use this library whenever marshaling/unmarshaling CRSLang structures.
- The WASM entry point (`wasm/main.go`) is gated with the `//go:build js && wasm` build tag and should never be imported by non-WASM code.
- Test files in the root package (`package main`) can use helpers defined in `listener_test.go` (e.g., `mustNewActionOnly`, `mustNewActionWithParam`, `mustNewSetvarAction`).
- Copyright header style: `// Copyright <year> <author>\n// SPDX-License-Identifier: Apache-2.0`

---

## Common Errors & Workarounds

| Error | Cause | Workaround |
|---|---|---|
| `go: toolchain go1.25.6 not found` (CI) | Invalid Go version in workflow file | Change `go-version` in `.github/workflows/go.yml` and `check_format.yml` to a real Go release (e.g., `'1.22.x'`) |
| `python: 3.14 not available` (CI) | Pre-release Python specified in workflow | Change `python-version` in `check_format.yml` to `'3.12'` or `'3.13'` |
| `Invalid variable name: <name>` | `listener` received an unknown SecLang variable | Add the new constant to `types/variables.go` and wire it up in `stringToVariableName()` |
| `Invalid collection name: <name>` | Unknown SecLang collection | Add to `types/collections.go` and `stringToCollectionName()` |
| WASM build fails on non-WASM machine | `wasm/main.go` has `js && wasm` build constraint | Always build WASM with `GOOS=js GOARCH=wasm go build -o wasm/crslang.wasm ./wasm/` |
