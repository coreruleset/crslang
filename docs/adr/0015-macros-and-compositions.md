# ADR-0015: Macros and Named Compositions

- **Status:** Proposed
- **Date:** 2026-04-25
- **Phase:** 3 (function composition era; partially shapes Phase 0 language base decision)

## Context

Multiple ADRs reference reusable function compositions ("named macros", "named function
compositions") without specifying how they're defined or scoped:

- ADR-0002 (pipeline operator): "reusable transform chains are defined once and invoked
  by name"
- ADR-0009 (language base evaluation): HCL vs custom decision partially hinges on
  what macros look like — HCL doesn't have parameterized expression macros natively,
  custom can have whatever shape the parser supports
- README Phase 3: "named function compositions (macros) keep nesting shallow"
- ADR-0014 (imports): macros are namespaced like groups and external data

The actual macro design has never been written down.

### What needs reuse

CRS today has many transformation chains that repeat across rules:

```
t:none,t:utf8toUnicode,t:urlDecodeUni,t:htmlEntityDecode,t:jsDecode,t:cssDecode,t:removeNulls
```

This 7-step "deep normalization" chain appears in dozens of XSS and SQLi rules. Today
it's copy-pasted with the risk of one rule getting it slightly wrong.

CRS also has predicate-side reuse:
- "is this method GET or HEAD?" → `request.method |> matches("^(?:GET|HEAD)$")`
- "is this from a trusted IP?" → `client.ip |> ip_in_range(trusted_ips)`
- "is this a static asset?" → `request.uri |> matches("\\.(?:css|js|png|gif)$")`

These are repeated across many rules. They deserve names.

### What macros are not

- **Not functions**: macros do not introduce new operators or runtime behavior. They are
  named compositions of existing functions.
- **Not templates**: macros do not generate rules. A macro produces a value (an
  expression result), not a rule.
- **Not stateful**: macros cannot read or write TX state. They are pure expressions.

## Decision

CRSLang adopts **typed expression macros** — named, parameterized compositions of
existing functions that produce a value when invoked. Macros are pure, side-effect-free,
and resolved at compile time.

### Definition Syntax

```
macro <name>([<param>: <type>, ...]) [: <return_type>] = <expression>
```

Examples:

```
# No parameters — a fixed pipeline applied to a fixed input field
macro is_safe_method() = request.method |> matches("^(?:GET|HEAD|OPTIONS)$")

# One parameter — applies a transformation chain to any input
macro deep_normalize(input: string) = input
  |> utf8_to_unicode()
  |> url_decode_uni()
  |> html_entity_decode()
  |> js_decode()
  |> css_decode()
  |> remove_nulls()

# Multi-parameter — predicate over two inputs
macro request_from(field: string, allowlist: ip_list) =
  field |> ip_in_range(allowlist)

# Composition — macros calling macros
macro is_xss_attempt(input: string) = input
  |> deep_normalize()
  |> detect_xss()
```

### Type Signatures

Every macro has a typed signature, inferred or declared:

- **Parameter types** — `string`, `int`, `bool`, `ip_list`, `regex_set`, `field_ref`,
  field-namespace types from ADR-0001
- **Return type** — usually inferred from the body. Can be declared for clarity:
  `macro is_xss(...) : bool = ...`

Type checking happens at compile time. A macro called with the wrong types is a
compile error.

### Scoping

Macros live at three scope levels:

1. **File scope** — declared at the top of a file. Visible directly within the
   declaring file, and visible to other files in the same package only when those
   files explicitly `import` the declaring file (per ADR-0014).
2. **Group scope** — declared inside a group, visible only within that group and any
   nested groups.
3. **Package scope** — declared with `export` (see [Exports](#exports) below), making
   the macro visible to consumers that import the package as a whole, under the
   package's namespace (per ADR-0014).

```
# package.crs (CRS distribution)
package "owasp_crs" version "4.18.0"

# Exported — visible to importers as crs.deep_normalize
export macro deep_normalize(input: string) = input
  |> utf8_to_unicode()
  |> url_decode_uni()
  ...

# rules/sql_injection.crs
group "sql_injection" {
  # Group-scoped macro (visible only inside this group, not to importers)
  macro args_are_sqli() = request.args |> detect_sqli()

  # Exported group-scope macro — visible to importers as crs.sql_injection.detect
  export macro detect() = request.args |> detect_sqli()

  rule 942100 (severity: critical) {
    when args_are_sqli()           # local macro
       or deep_normalize(request.body) |> contains("union")  # package macro
    then block
  }
}

# deployer's main.crs
import package "owasp_crs/4.18" as crs

rule 9100100 {
  when crs.deep_normalize(request.headers["X-Custom"]) |> contains("attack")
  then block
}
```

The `export` keyword controls cross-package visibility — see [Exports](#exports) below.

### Exports

Macros default to **private**: visible only within their declaring scope. To make a
macro visible to importers, prefix the declaration with the `export` keyword:

```
# rules/_common.crs (inside the OWASP_CRS package)
export macro deep_normalize(input: string) = input
  |> utf8_to_unicode()
  |> url_decode_uni()
  |> remove_nulls()

# Not exported — file-scope helper, invisible outside _common.crs
macro _strip_quotes(input: string) = input |> replace("\"", "")
```

`export` is allowed at file scope and at group scope. The qualified path importers
use is determined by ADR-0014's namespacing rules:

| Declaration site | Importer-visible name (with alias `crs`) |
|---|---|
| File scope, anywhere in the package | `crs.<macro_name>` |
| Group scope inside a group named `bar` | `crs.bar.<macro_name>` |
| Group scope inside nested groups `outer.inner` | `crs.outer.inner.<macro_name>` |

`export` does not affect visibility *within* the package — a non-exported macro is
still usable by other files in the same package via an explicit `import` of the
declaring file (per ADR-0014). The keyword only controls visibility *across* the
package boundary, i.e., whether consumers that `import package "<name>"` can
reference the macro under the package namespace.

#### Default visibility rationale

CRSLang chooses **explicit opt-in for exports** rather than opt-out (e.g., "every
top-level macro is automatically exported"). Adding a new macro inside a package does
not accidentally widen the public API surface. The `export` keyword is grep-able,
making the package's external API enumerable with `rg '^export macro'`. This matches
the convention in Rust (`pub`), TypeScript (`export`), and Python (`__all__`).

#### Re-export from package.crs

A package author who wants to expose macros that live in other files declares
re-exports in `package.crs`:

```
# package.crs
package "owasp_crs" version "4.18.0"

import "rules/_common.crs"
import "rules/sql_injection.crs"

# Re-export individual names from imported files
export macro deep_normalize from "rules/_common.crs"
export macro is_safe_method from "rules/_common.crs"
```

Re-exports do not redefine the macro; they make an already-defined macro from another
file visible to importers under the package namespace. This is needed only when the
file authoring style keeps macros private at the file level but the package wants to
expose them; if the macro is declared with `export macro` directly in its file, no
re-export is needed.

#### Conflict with ADR-0014's name resolution

Re-exports go through the same conflict-resolution rules as ordinary declarations
(ADR-0014). Two `export` declarations producing the same fully-qualified name across
files of the same package are a compile-time error.

#### HCL surface

In the HCL surface (per ADR-0009), the equivalent is a `visibility` attribute on the
macro block:

```hcl
macro "deep_normalize" {
  visibility = "public"
  param "input" { type = string }
  body  = "..."
}
```

The IR is identical: both surfaces produce a `MacroDecl` with `Exported = true`.

### Pure Expression Constraint

Macro bodies are pure expressions. They cannot:

- Read or write `tx.*`, `ip.*`, `global.*`, `session.*` (these are runtime state)
- Have side-effects (logging, scoring, action invocation)
- Reference rule attributes (severity, paranoia)
- Contain `when`, `then`, or any rule-level keyword

They can:

- Reference any field from the typed namespace (ADR-0001)
- Call any function in the standard library
- Call other macros (acyclically — recursion is a compile error)
- Use boolean operators (`and`, `or`, `not`) per ADR-0003
- Use the pipeline operator (custom parser) or nested calls (HCL)

This constraint keeps macros tractable: a macro call is always replaceable by its
body. Tooling can inline macros without changing semantics.

### No Recursion

Macros cannot call themselves directly or indirectly. The compile-time call graph must
be acyclic. This guarantees macro expansion terminates.

### Compilation Strategy

Macros are **inlined at compile time**. A macro call expands into its body with
parameter substitution:

```
# Authored
rule 942100 {
  when args_are_sqli()
  then block
}

# After inlining (intermediate form)
rule 942100 {
  when request.args |> detect_sqli()
  then block
}

# Compiled SecLang (no trace of the macro)
SecRule ARGS "@detectSQLi" "id:942100,phase:2,block,..."
```

The compiled output has no notion of "macro" — only the expanded expressions.

### HCL vs Custom Surface

The macro concept is the same in both language bases (Phase 0 / ADR-0009), but the
surface differs:

**Custom parser:**

```
macro deep_normalize(input: string) = input
  |> utf8_to_unicode()
  |> url_decode_uni()
  |> remove_nulls()
```

**HCL:**

HCL doesn't have first-class parameterized macros. Two options:

- Use `locals` with no parameters (function-call equivalents):

  ```hcl
  locals {
    safe_methods_pattern = "^(?:GET|HEAD|OPTIONS)$"
  }
  ```

- Add a custom `macro` block via HCL extension (the project already extends HCL with
  custom block types):

  ```hcl
  macro "deep_normalize" {
    param "input" { type = string }
    body  = "url_decode(html_entity_decode(remove_nulls(input)))"
  }
  ```

  HCL's `templatefile()` and `function()` extension mechanisms make this feasible.

The IR is identical regardless of surface — both forms produce the same
`MacroDecl` and `MacroCall` AST nodes.

### Macro vs Function Distinction

| Aspect | Function | Macro |
|---|---|---|
| Defined by | Standard library or engine | Ruleset author |
| Implementation | Native code (Go) | CRSLang expression |
| Side effects | Possibly (rare) | None |
| Compilation | Stays as a call | Inlined |
| Examples | `matches()`, `count()`, `detect_sqli()` | `deep_normalize()`, `is_safe_method()` |

Functions are the language's primitives; macros are user-defined named compositions
of those primitives.

### IR Representation

```go
type MacroDecl struct {
    Doc        *Documentation        // from ADR-0013
    Name       string                // local name within scope
    Namespace  string                // package/file/group, per ADR-0014
    Exported   bool                  // visible to importers (set by `export` keyword)
    Params     []MacroParam
    ReturnType Type
    Body       Expr                  // the typed AST of the expression
}

type MacroParam struct {
    Name string
    Type Type
}

type MacroCall struct {
    Decl *MacroDecl
    Args []Expr
}
```

After macro inlining, `MacroCall` nodes are replaced by the substituted `Expr` from
the macro's body. The decl is retained in the IR for tooling (cross-reference,
documentation, refactoring).

## Alternatives Considered

### A: No macros — copy-paste pipelines

Require every rule to spell out its full pipeline.

**Rejected because:**
- CRS already has the deep-normalize chain repeated across many rules; copy-paste
  causes drift and inconsistency
- Higher cognitive load when reading rules
- Defeats one of the language's stated benefits (composability)

### B: Macros as engine functions (no user-defined)

Ship the standard library with `deep_normalize()` etc. baked in. Users cannot define
their own.

**Rejected because:**
- Different deployers have different normalization needs (e.g., different sequences
  for different content types)
- Custom rule packs (third-party) want to ship their own conventions
- Standard library would grow indefinitely to cover everyone's use cases

### C: Macros with side effects (full procedures)

Allow macros to log, set TX, etc.

**Rejected because:**
- Conflates expressions with statements
- Compile-time inlining no longer trivially correct (ordering of side effects matters)
- The action model (ADR-0004) already handles side effects in `then` blocks; macros
  would duplicate or compete with that

### D: Templating (rule-generating macros)

`for severity in [critical, warning] { rule N { ... } }` style metaprogramming.

**Rejected because:**
- Generated rules are hard to understand and debug
- Rule IDs become opaque (where did rule 942101 come from?)
- The current model (group inheritance, severity attribute) covers the common cases

### E: Recursive macros

Allow self-referential macros for tree-walking patterns.

**Rejected because:**
- Compilation can fail to terminate without careful guards
- No CRS use case currently requires recursion
- Can be added later if a real need emerges

## Consequences

### Positive

- Eliminates copy-paste of common pipelines (deep_normalize, safe_method check, etc.)
- Type-checked at compile time — wrong-type arguments caught early
- Macros + groups + imports + documentation form a coherent modularity story
- Tooling can show macro expansion (e.g., IDE inlay hints, `--explain` mode)
- Compiled output has no macro overhead — pure inlining

### Negative

- Authors learn a second naming layer (functions vs macros) and where each scope lives
- Group-scoped macros mean the same name can resolve to different things in different
  contexts; mitigated by IDE tooling and clear error messages
- HCL surface for macros is awkward without language extension; pushes weight onto the
  Phase 0 decision

### Risks

- **Macro proliferation** — without discipline, deployers may define macros for
  trivial expressions, hurting readability. Style guide and lint rules can encourage
  reuse only when a pipeline appears 3+ times.
- **Type-checking complexity** — full type inference across macro calls requires the
  type system from the second-tier "function signatures" gap. Until that's done,
  type checking may be partial. Acceptable as long as runtime errors remain rare.
- **Cross-target divergence** — non-SecLang targets may not support the same function
  set; macros that compile cleanly to SecLang may fail to compile to Cloud Armor.
  The multi-target compiler (ADR-0010) must report which macros are unsupported per
  target.
