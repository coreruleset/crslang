# ADR-0005: Parser Strategy

- **Status:** Proposed (see also [ADR-0009](0009-language-base-evaluation.md) for broader evaluation including HCL, CEL, and Expr)
- **Date:** 2026-04-13
- **Phase:** 4 (implementation of ADR-0009 decision)

## Context

CRSLang currently has two parsers:

1. **SecLang parser** — ANTLR4-based, external dependency
   (`github.com/coreruleset/seclang_parser`). Parses `.conf` files into a parse tree,
   which a listener converts to the Go AST.
2. **YAML parser** — `go.yaml.in/yaml/v4`. Loads YAML directly into typed structs.

Phase 4 introduces a third format: the native CRSLang text syntax. This ADR decides
how to build its parser.

## Decision

Build a **hand-written recursive-descent parser** in Go with no external dependencies.

### Rationale

The CRSLang grammar (see below) is small and regular. It does not need:
- Ambiguous grammar resolution (no ambiguity by design)
- Left-recursive rules (the grammar is LL(1) with minor lookahead)
- Grammar-driven code generation (the AST types already exist)

A hand-written parser provides:
- **Zero dependencies** — no ANTLR runtime, no build-time code generation
- **Clear error messages** — custom error recovery with context-aware messages
- **Full control** — easy to add syntax extensions without regenerating code
- **Performance** — no generic parser overhead
- **Debuggability** — step through parsing in a standard Go debugger

### Grammar

```ebnf
(* Top-level *)
file           = { rule | globals | defaults | comment } ;

(* Rules *)
rule           = "rule" INTEGER metadata? "{" when_clause then_clause "}" ;
metadata       = "(" key_value { "," key_value } ")" ;
key_value      = IDENT ":" value ;

(* Conditions *)
when_clause    = "when" expr ;
expr           = and_expr { "or" and_expr } ;
and_expr       = unary_expr { "and" unary_expr } ;
unary_expr     = "not" unary_expr
               | "(" expr ")"
               | pipeline ;

(* Pipelines *)
pipeline       = pipeline_source { "|>" func_call } ;
pipeline_source= field_ref | func_call | literal | "(" expr ")" ;
field_ref      = IDENT { "." IDENT } [ "[" selector "]" ] ;
selector       = STRING | "!" STRING ;
func_call      = IDENT "(" [ arg { "," arg } ] ")" ;
literal        = STRING | INTEGER | "true" | "false" ;

(* Actions *)
then_clause    = "then" disruptive [ "(" named_args ")" ] [ effect_block ] ;
disruptive     = "pass" | "block" | "deny" | "drop" | "allow" | "redirect" ;
effect_block   = "{" { effect_stmt } "}" ;
effect_stmt    = assignment | func_call | configure_block ;
assignment     = field_ref assign_op value ;
assign_op      = "=" | "+=" | "-=" ;
configure_block= "configure" "{" { IDENT "=" value } "}" ;

(* Globals and defaults *)
globals        = "globals" "{" { key_value } "}" ;
defaults       = "defaults" "{" { key_value } "}" ;

(* Rule management *)
exclude_rule   = "exclude" "rule" INTEGER ;
update_rule    = "update" "rule" INTEGER "{" update_body "}" ;

(* Literals *)
value          = STRING | INTEGER | FLOAT | BOOLEAN | list ;
list           = "[" [ value { "," value } ] "]" ;

(* Tokens *)
STRING         = '"' { char } '"' ;
INTEGER        = digit { digit } ;
IDENT          = letter { letter | digit | "_" } ;
comment        = "#" { any } newline ;
```

### Parser Architecture

```
Source Text
    │
    ▼
┌──────────┐     ┌──────────┐     ┌──────────┐
│  Lexer   │────▶│  Parser  │────▶│   AST    │
│ (tokens) │     │ (recdes) │     │ (Go IR)  │
└──────────┘     └──────────┘     └──────────┘
                                       │
                              ┌────────┴────────┐
                              ▼                  ▼
                         ┌─────────┐      ┌──────────┐
                         │  YAML   │      │ SecLang  │
                         │ export  │      │ export   │
                         └─────────┘      └──────────┘
```

**Components:**

1. **Lexer** (`parser/lexer.go`) — hand-written scanner producing tokens:
   - Keywords: `rule`, `when`, `then`, `and`, `or`, `not`, `globals`, `defaults`, etc.
   - Operators: `|>`, `=`, `+=`, `-=`
   - Delimiters: `(`, `)`, `{`, `}`, `[`, `]`, `,`, `.`
   - Literals: strings, integers, identifiers
   - Comments: `#` to end of line

2. **Parser** (`parser/parser.go`) — recursive-descent functions:
   - One function per grammar production
   - Returns typed AST nodes (shared with YAML and SecLang importers)
   - Produces clear errors with line/column information

3. **Error recovery** — on syntax error:
   - Report the error with context (expected vs found, surrounding tokens)
   - Attempt to synchronize at the next `rule` keyword or `}`
   - Continue parsing to find multiple errors in one pass

### Package Structure

```
parser/
├── lexer.go          # Token scanner
├── lexer_test.go
├── token.go          # Token types
├── parser.go         # Recursive-descent parser
├── parser_test.go
├── error.go          # Error types and formatting
└── printer.go        # AST -> CRSLang text (formatter/pretty-printer)
```

### Three-Importer Architecture

After Phase 4, the system has three importers and two exporters sharing a single IR:

```
SecLang (.conf) ──► ANTLR parser ──► SecLang AST ──► Normalize ──┐
                                                                   │
YAML (.yaml)   ──► yaml.Unmarshal ──► YAML structs ──► Normalize ─┤──► Unified IR
                                                                   │
CRSLang (.crs) ──► Hand-written ──► AST nodes ────────────────────┘
                    parser                                │
                                                          ├──► YAML exporter
                                                          ├──► CRSLang text exporter
                                                          └──► SecLang exporter (existing)
```

### WASM/Playground Integration

The hand-written parser compiles to WASM without issues (no CGo, no external deps).
The playground gains a third pane/mode for the native syntax:

```javascript
// Existing
seclangToCRSLang(input)    // SecLang -> YAML
crslangToSeclang(input)    // YAML -> SecLang

// New
parseCRSLang(input)        // CRSLang text -> IR (validates)
formatCRSLang(ir)          // IR -> CRSLang text (pretty-print)
crslangTextToYaml(input)   // CRSLang text -> YAML
crslangTextToSeclang(input)// CRSLang text -> SecLang
```

### Formatting Conventions

The pretty-printer enforces consistent style:

```
# One condition per line, aligned
rule 920170 (phase: request, severity: warning) {
  when request.method |> matches("^(?:GET|HEAD)$")
   and request.headers["Content-Length"] |> not(matches("^0?$"))
  then block {
    tx.anomaly_score += 5
    log(audit: true)
  }
}

# Short rules on fewer lines
rule 900100 (phase: request) {
  when true
  then pass { tx.paranoia_level = 1 }
}
```

Indentation: 2 spaces. Line width guide: 100 characters. `and`/`or` align under `when`.

## Alternatives Considered

### A: ANTLR4 for the new syntax

Use ANTLR4 (same as SecLang parser) with a new `.g4` grammar.

**Rejected because:**
- Adds build-time code generation step
- ANTLR runtime is a heavy dependency for a small grammar
- Generated code is hard to customize for error messages
- The SecLang grammar is complex (legacy syntax); the CRSLang grammar is intentionally
  simple and does not need ANTLR's power

### B: PEG parser (pigeon, peg)

Use a PEG parser generator for Go.

**Considered viable but:**
- Still requires code generation
- PEG's ordered alternation can produce surprising behavior
- Error messages from PEG parsers are typically poor
- For this grammar size, hand-writing is faster than debugging a PEG grammar

### C: Parser combinator library (participle)

Use a Go parser combinator library that derives the parser from struct tags.

**Considered viable but:**
- Adds a runtime dependency
- Grammar changes require changing struct definitions, coupling syntax to IR
- Less control over error messages and recovery
- The CRSLang IR already exists; forcing it to match a combinator library's
  conventions would be counterproductive

### D: Tree-sitter grammar

Write a Tree-sitter grammar for editor integration, use it as the parser.

**Rejected as primary parser** — Tree-sitter is designed for incremental editor parsing,
not batch compilation. However, a Tree-sitter grammar **should be written alongside** the
hand-written parser for editor support (syntax highlighting, code folding, etc.).

## Consequences

### Positive

- Zero new dependencies
- Compiles to WASM without issues
- Full control over error messages and recovery
- Easy to extend as the language grows
- Debuggable with standard Go tools
- Fast compilation (no code generation step)

### Negative

- More initial code to write than using a parser generator
- No grammar file that doubles as documentation (mitigated by the EBNF in this ADR
  and a separate grammar reference document)
- Must be kept in sync with any Tree-sitter grammar manually

### Risks

- **Grammar evolution** — as the language grows, the hand-written parser must be updated
  in lockstep. Mitigated by comprehensive parser tests and the small grammar size.
- **Edge cases** — hand-written parsers can have subtle bugs in lookahead and error
  recovery. Mitigated by fuzzing and property-based testing.
