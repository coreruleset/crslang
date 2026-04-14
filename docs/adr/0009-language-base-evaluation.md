# ADR-0009: Language Base Evaluation — HCL, CEL, Expr, or Custom Parser

- **Status:** Proposed
- **Date:** 2026-04-13
- **Phase:** 0 (foundational — must decide before IR design in Phases 1-3)

## Context

ADR-0005 proposed a hand-written recursive-descent parser for CRSLang's text syntax.
Before committing to building a parser from scratch, we should evaluate whether an
existing language can serve as CRSLang's foundation — particularly given new requirements
that emerged from design discussions:

### New Requirements

**1. Reusable transform composition (macros/functions)**

CRS rules repeat the same transformation chains across many rules. Instead of spelling
out 4-5 transforms per rule, authors should define named compositions:

```
# Instead of repeating this in every SQL injection rule:
request.args |> url_decode() |> lowercase() |> remove_whitespace() |> detect_sqli()

# Define once, use everywhere:
detect_sqli(normalize(request.args))
# Or with pipeline:
request.args |> normalize() |> detect_sqli()
```

This means the pipeline operator (`|>`) is a per-rule convenience, not a language-level
necessity. If transforms are encapsulated in named functions, nesting depth stays at 1-2
levels, and nested function call syntax (`detect_sqli(normalize(field))`) is readable.

**2. Global defaults (version, common tags)**

CRS rules repeat `version: OWASP_CRS/4.18.0-dev` and `tag: OWASP_CRS` on every rule.
CRSLang v1 already extracts these (`ExtractDefaultValues` in `configuration.go`). The
new syntax should make this a first-class concept:

```
# Applied to all rules in all files
globals {
  version = "OWASP_CRS/4.18.0-dev"
  tags    = ["OWASP_CRS"]
}
```

**3. File-level / group-level defaults**

Some tags and metadata repeat within a file or logical group but not globally:

```
# Applied to all rules in this file/group
group "sql_injection" {
  tags = ["OWASP_CRS/SQL_INJECTION", "attack-sqli"]

  rule 942100 { ... }
  rule 942110 { ... }
}
```

**4. Condition expression is the key design question**

The structural parts (rule blocks, metadata, groups, globals) are straightforward in
any block-based language. The critical differentiator is how conditions are expressed —
the boolean expressions with function calls that form the matching logic.

## Evaluation Criteria

| Criterion | Weight | Description |
|-----------|--------|-------------|
| Condition expressiveness | High | Can it express boolean combinations of field tests with transforms? |
| Function composition | High | Can users define reusable transform chains? |
| Block structure | Medium | Does it naturally support rules, groups, globals? |
| Go ecosystem | Medium | Library quality, license, WASM support, maintenance |
| Extensibility | Medium | Can we add CRSLang-specific constructs without forking? |
| Learning curve | Medium | Familiarity for security engineers and DevOps |
| Pipeline support | Low | Nice-to-have if macro functions reduce nesting (see above) |

## Candidates

### Option A: HCL2

**Library:** `github.com/hashicorp/hcl/v2` (MPL-2.0)

**Structure:**

```hcl
globals {
  version = "OWASP_CRS/4.18.0-dev"
  tags    = ["OWASP_CRS"]
}

group "sql_injection" {
  tags = ["OWASP_CRS/SQL_INJECTION", "attack-sqli"]

  rule "942100" {
    severity = "critical"
    message  = "SQL Injection Attack Detected"

    condition = detect_sqli(normalize(request.args))
    action    = "block"

    effects {
      tx_sqli_score    = "+=5"
      tx_anomaly_score = "+=5"
      capture          = true
      log_audit        = true
    }
  }
}

exclude "rule" {
  id = 942100
}

update "rule" {
  id = 920170
  remove_target = "request.args[\"username\"]"
}
```

**Conditions:**

```hcl
# Simple match
condition = matches(request.method, "^(?:GET|HEAD)$")

# Boolean composition
condition = (
  matches(request.method, "^(?:GET|HEAD)$")
  && !matches(request.headers["Content-Length"], "^0?$")
)

# With named transforms (custom functions registered in Go)
condition = detect_xss(normalize_input(request.args))

# Complex boolean
condition = (
  (eq(request.method, "POST") || eq(request.method, "PUT"))
  && !contains(request.headers["Content-Type"], "application/json")
  && gt(length(request.body), 0)
)
```

**Custom functions:** Registered in Go via `function.New()`. CRSLang defines a function
registry that includes both primitive transforms (`lowercase`, `url_decode`) and
composite macros (`normalize_input` = url_decode + lowercase + remove_whitespace).
Users can extend the registry with their own macros.

Functions like Sprig (`github.com/Masterminds/sprig`) can be registered directly,
giving access to 100+ utility functions (string manipulation, date formatting, crypto,
etc.) for free.

**Strengths:**
- Block syntax is an excellent fit for rules, groups, globals
- Widely known in DevOps/infrastructure community (Terraform)
- Custom functions are trivial to add from Go
- `&&`, `||`, `!` for boolean expressions work naturally
- Sprig integration gives a rich function library for free
- HCL's `for` expressions could enable iteration over collections
- Mature parser with good error messages and recovery
- Compiles to WASM (pure Go, no CGo)

**Weaknesses:**
- No pipeline operator (`|>`) — grammar is fixed, cannot be extended
- Conditions are attribute values, not first-class syntax — they lack visual prominence
- Function nesting reads inside-out (mitigated by named compositions)
- `effects` block uses string values for operations (`"+=5"`) because HCL attributes
  are assignments, not arbitrary statements
- No `let` bindings in expressions for data-flow dependencies (ADR-0003 Category 3)
- HCL variable references use `var.name` prefix which would need custom handling

**Pipeline workaround with Go templates:**

HCL2 supports Go template syntax in string interpolation, and Go templates have `|`
(pipe) as a first-class operator. However, this only works inside strings:

```hcl
# This is string interpolation, not expression evaluation
condition = "${request.args | normalize | detect_sqli}"
```

This is not viable — it would mean conditions are opaque strings, losing type checking
and IDE support. The pipe in Go templates is for text rendering, not AST construction.

### Option B: CEL (Common Expression Language)

**Library:** `github.com/google/cel-go` (Apache-2.0)

**Structure:** CEL is an expression-only language — it has no block syntax. It would need
to be embedded inside another structural format (HCL, YAML, or custom blocks).

**Hybrid approach — HCL for structure, CEL for conditions:**

```hcl
globals {
  version = "OWASP_CRS/4.18.0-dev"
  tags    = ["OWASP_CRS"]
}

rule "920170" {
  severity = "warning"

  # CEL expression as a string attribute, parsed by cel-go
  condition = "request.method.matches('^(?:GET|HEAD)$') && !request.headers['Content-Length'].matches('^0?$')"

  action = "block"
  effects {
    tx_anomaly_score = "+=5"
  }
}
```

**CEL expression examples:**

```cel
// Method chaining (CEL's strength)
request.method.matches("^(?:GET|HEAD)$")

// Boolean composition
request.method.matches("^POST$") && request.body.size() > 0

// With custom functions
normalize(request.args).detect_sqli()

// With macros (CEL supports parse-level macros)
request.args.all(arg, !arg.detect_sqli())

// Exists macro
request.headers.exists(h, h.key == "X-Forwarded-For")
```

**Custom functions:** Defined in Go via `cel.Function()`. CEL also supports **custom
macros** that rewrite the AST at parse time — more powerful than runtime functions.

**Strengths:**
- Purpose-built for policy and rule evaluation (Google IAM, Kubernetes, Envoy)
- Method-style chaining: `field.transform().predicate()` — not a pipeline but close
- CEL macros enable reusable compositions at the AST level
- Strong type system with protocol buffer integration
- Partial evaluation support (useful for optimization)
- `has()` macro for optional field checks
- Battle-tested in production at scale

**Weaknesses:**
- Expression-only — needs a structural wrapper (HCL, YAML, or custom)
- Two parsers: one for structure, one for conditions
- CEL expressions as strings inside HCL lose syntax highlighting and editor support
- Method-style `field.transform()` creates ambiguity with field access (`request.headers`
  vs `request.method.matches()` — is `.matches` a field or a method?)
- Heavier dependency than HCL alone
- WASM support: `cel-go` is pure Go and compiles to WASM, but the binary size is
  significant

### Option C: Expr

**Library:** `github.com/expr-lang/expr` (MIT)

**Structure:** Like CEL, Expr is expression-only and needs a structural wrapper.

**Hybrid approach — HCL for structure, Expr for conditions:**

```hcl
rule "920170" {
  condition = "matches(request.method, '^(?:GET|HEAD)$') && !matches(request.headers['Content-Length'], '^0?$')"
  action    = "block"
}
```

**Expr expression examples:**

```
// Function calls
matches(request.method, "^(?:GET|HEAD)$")

// Boolean
matches(request.method, "POST") && length(request.body) > 0

// Member access
request.headers["Content-Type"] contains "json"

// Built-in operators
request.method in ["GET", "HEAD", "OPTIONS"]

// With custom functions
detect_sqli(normalize(request.args))
```

**Custom functions:** Defined in Go via `expr.Function()`. Expr also supports operator
overloading.

**Strengths:**
- Very lightweight (MIT license, small dependency)
- Fast evaluation (compiles to bytecode)
- Simple API: `expr.Compile()` + `expr.Run()`
- Built-in operators: `in`, `contains`, `matches`, `startsWith`, `endsWith`
- Familiar syntax (similar to JavaScript/Python expressions)
- Compiles to WASM easily

**Weaknesses:**
- No method chaining or pipeline
- Less mature than CEL for policy use cases
- No macro system — compositions must be runtime functions
- Same two-parser problem as CEL
- Expressions as strings lose editor support

### Option D: Custom Parser (ADR-0005)

**Library:** None — hand-written lexer + recursive-descent parser in Go.

**Structure and conditions unified:**

```
globals {
  version = "OWASP_CRS/4.18.0-dev"
  tags    = ["OWASP_CRS"]
}

macro normalize(field) {
  field |> url_decode() |> lowercase() |> compress_whitespace()
}

group "sql_injection" {
  tags = ["OWASP_CRS/SQL_INJECTION", "attack-sqli"]

  rule 942100 (severity: critical) {
    when detect_sqli(normalize(request.args))
      or detect_sqli(normalize(request.body))
    then block {
      tx.sqli_score += 5
      tx.anomaly_score += 5
      capture()
      log(audit: true)
    }
  }
}
```

**Strengths:**
- Pipeline operator (`|>`) — native, first-class
- `macro` / user-defined functions — in the language, not just in Go
- Unified parser: structure and conditions in one grammar
- Full control over syntax, error messages, and evolution
- `let` bindings (ADR-0003) trivially addable
- `when`/`then` keywords give conditions visual prominence
- Assignment operators (`+=`, `-=`) native in effect blocks
- Zero dependencies
- WASM: trivial

**Weaknesses:**
- Must write and maintain the parser (~1500-2500 lines of Go)
- No pre-existing editor support (must write Tree-sitter grammar, LSP, etc.)
- No ecosystem of existing functions (Sprig, etc.) — must build the function library
- Novel syntax — no existing community familiarity
- Risk of language design mistakes without real-world iteration

### Option E: HCL for Structure + Custom Expression Micro-Language for Conditions

A hybrid that uses HCL for the block structure but replaces HCL's expression language
with a custom parser for conditions only.

```hcl
globals {
  version = "OWASP_CRS/4.18.0-dev"
  tags    = ["OWASP_CRS"]
}

macro "normalize" {
  steps = ["url_decode", "lowercase", "compress_whitespace"]
}

group "sql_injection" {
  tags = ["OWASP_CRS/SQL_INJECTION", "attack-sqli"]

  rule "942100" {
    severity = "critical"

    # Custom expression language parsed separately
    when = <<-EXPR
      detect_sqli(normalize(request.args))
      or detect_sqli(normalize(request.body))
    EXPR

    action = "block"

    effects {
      tx_sqli_score    = "+=5"
      tx_anomaly_score = "+=5"
      capture          = true
    }
  }
}
```

**Strengths:**
- HCL handles the hard parts: block parsing, heredocs, interpolation, error recovery
- Custom expression parser is small (only boolean + function calls + field refs)
- HCL heredoc syntax (`<<-EXPR ... EXPR`) cleanly embeds multi-line expressions
- Macros can be HCL blocks that the compiler processes
- Editor support: HCL highlighting works for structure; expression highlighting
  could be layered

**Weaknesses:**
- Two parsers, two grammars, two mental models
- Conditions as heredoc strings lose HCL-level type checking
- Macros are HCL blocks (declarative list of steps) rather than composable expressions
- More complex build: HCL dependency + custom parser code

## Comparison Matrix

| | HCL (A) | HCL+CEL (B) | HCL+Expr (C) | Custom (D) | HCL+Custom (E) |
|---|---|---|---|---|---|
| **Block structure** | Native | HCL | HCL | Must build | HCL |
| **Condition syntax** | `f(g(x)) && h(y)` | `x.g().f() && y.h()` | `f(g(x)) && h(y)` | `x \|> g() \|> f() and y \|> h()` | Custom in heredoc |
| **Pipeline** | No | Method chain | No | Native `\|>` | Possible in custom |
| **User macros** | Go functions only | CEL macros (AST) | Go functions only | Language-level | HCL blocks + Go |
| **Globals/groups** | Native blocks | HCL blocks | HCL blocks | Must build | HCL blocks |
| **Go dependency** | `hcl/v2` (MPL) | `hcl/v2` + `cel-go` (Apache) | `hcl/v2` + `expr` (MIT) | None | `hcl/v2` (MPL) |
| **WASM size** | Medium | Large | Medium | Small | Medium |
| **Editor support** | Terraform ecosystem | Partial | Partial | Must build | Partial |
| **Learning curve** | Low (Terraform users) | Medium (two syntaxes) | Low-Medium | Medium (new syntax) | Medium (two layers) |
| **Sprig/function libs** | Direct integration | Via CEL functions | Via Expr functions | Must wrap | Direct for HCL parts |
| **Effects syntax** | String values (`"+=5"`) | String values | String values | Native (`+= 5`) | String values |
| **Condition prominence** | Attribute value | String attribute | String attribute | `when`/`then` keywords | Heredoc block |
| **Type checking** | HCL schema | CEL type system | Expr type system | Custom | Split |
| **Maintenance cost** | Low | Medium | Low-Medium | Medium-High | Medium |

## Key Trade-Off Analysis

### The Nesting Question (With Macros)

With named composition functions, the nesting concern largely dissolves:

```hcl
# HCL: 1-2 levels of nesting — readable
condition = detect_sqli(normalize(request.args))
condition = detect_xss(normalize(request.headers["User-Agent"]))

# vs. Custom: pipeline — also readable
when request.args |> normalize() |> detect_sqli()
when request.headers["User-Agent"] |> normalize() |> detect_xss()
```

Both are acceptable. The pipeline is more readable for ad-hoc transform chains that
are not yet captured in a macro, but macros are the steady state for mature rulesets.

### The Effects Block Problem

HCL's attribute syntax does not naturally express side-effect operations:

```hcl
# Awkward: operations encoded as strings
effects {
  tx_anomaly_score = "+=5"     # Is this assignment or string?
  capture          = true       # Ok
  log_audit        = true       # Ok
}

# Or: structured but verbose
effect "setvar" {
  target = "tx.anomaly_score"
  op     = "+="
  value  = 5
}
```

A custom parser handles this natively:

```
then block {
  tx.anomaly_score += 5
  capture()
  log(audit: true)
}
```

This is a genuine advantage of the custom parser — effect blocks are imperative
statements, and HCL is not designed for imperative code.

### The Two-Parser Problem (Options B, C, E)

Hybrid approaches add complexity:
- Two grammars to document and maintain
- Editor support must understand both layers
- Error messages cross parser boundaries ("HCL parsed fine, but the CEL expression
  on line 15 has a type error")
- Macro definitions live in one syntax but are invoked in another

This is manageable but adds friction. A single-parser approach (A or D) is simpler.

### The Maintenance Argument

HCL's parser is maintained by HashiCorp. A custom parser must be maintained by the
CRSLang team. For a small team, this is a real cost — but the grammar is small (~30
productions) and well-understood. The total parser code is estimated at 1500-2500 lines
of Go, comparable to a medium-sized package.

## Decision

**Deferred.** This ADR presents the analysis; the decision depends on the team's
priorities:

**Choose HCL (Option A) if:**
- Minimizing parser maintenance is the top priority
- The team values Terraform-ecosystem familiarity
- Macros defined in Go (not in the language) are acceptable
- The effects block can use structured sub-blocks instead of imperative syntax
- Sprig / existing function library integration is important

**Choose Custom Parser (Option D) if:**
- The pipeline operator is valued even with macros available
- Language-level macro definitions are desired (users define macros in `.crs` files)
- Native `when`/`then` structure and imperative effects matter
- The team is comfortable maintaining a small parser
- Full control over language evolution is important

**Avoid the hybrid options (B, C, E)** unless there is a compelling reason to split
the parser. The added complexity of two grammars outweighs the benefits for a language
this small.

### Regardless of Choice

Some design elements are constant across all options:

1. **Globals block** — inherited by all rules (version, common tags)
2. **Group blocks** — file-level or logical grouping with inherited tags/metadata
3. **Named function compositions** — whether Go-registered (HCL) or language-level
   (custom), reusable transform chains are essential
4. **SecLang import** — the ANTLR-based importer remains unchanged
5. **YAML serialization** — the IR can always be exported to YAML
6. **WASM/playground** — all options compile to WASM

## Consequences

### If HCL Is Chosen

- ADR-0005 is superseded; no custom parser needed
- Pipeline operator (ADR-0002) is dropped from the text syntax; macros replace it
- Effects (ADR-0004) must be redesigned for HCL's declarative attribute model
- The function registry becomes the primary extension point
- Sprig and other Go function libraries can be integrated directly
- Editor support comes largely for free (HCL plugins exist)

### If Custom Parser Is Chosen

- ADR-0005 stands as designed
- Pipeline, macros, effects all work as proposed in earlier ADRs
- Must build editor support (Tree-sitter grammar, LSP)
- Function library must be built from scratch (no Sprig shortcut)
- More initial investment, more long-term control
