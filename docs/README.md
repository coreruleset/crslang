# CRSLang Evolution: From YAML/SecLang to Composable Expressions

## Vision

CRSLang aims to become a modern, composable rule language for web application firewalls.
Today it is a YAML serialization of ModSecurity's SecLang AST. The goal is to evolve it
into a standalone expression language where conditions are composed from typed fields,
piped through transformation functions, and combined with boolean algebra — similar in
spirit to Cloudflare's Wirefilter or Google's CEL, but purpose-built for the CRS
ecosystem.

## Current State

CRSLang v1 is a bidirectional translator between SecLang `.conf` files and structured
YAML. Internally it models SecLang concepts directly:

- **Variables and Collections** — two separate concepts (`REQUEST_METHOD` vs
  `REQUEST_HEADERS:Content-Type`) mapped from SecLang's target system
- **Transformations** — a flat ordered list (`t:lowercase,t:urlDecode`) applied before
  the operator
- **Operators** — a single operator per condition (`@rx`, `@pm`, `@eq`, etc.)
- **Actions** — a bag of disruptive, non-disruptive, flow, and data actions
- **Chains** — sequential rule linking via the `chain` flow action, simulating AND logic
- **Phases** — string metadata (`"1"` through `"5"`) controlling evaluation order,
  redundant with the variables used in the rule in most cases

A rule today looks like:

```yaml
kind: rule
metadata:
  id: 920170
  phase: "1"
  message: "GET/HEAD with body"
  severity: WARNING
conditions:
  - collections:
      - name: REQUEST_HEADERS
        arguments:
          - Content-Type
    operator:
      name: rx
      value: "^application/json"
    transformations:
      - lowercase
actions:
  disruptive:
    action: block
  non-disruptive:
    - action: log
    - action: setvar
      param: "tx.anomaly_score=+5"
```

### Strengths

- Lossless round-trip between SecLang and YAML
- Familiar to CRS maintainers
- Toolable (YAML parsers are everywhere)
- Playground with WASM support

### Limitations

- Verbose: a simple rule requires deep nesting
- SecLang concepts leak through (chains, `count: true` flags, string phases)
- No boolean composition: `chain` is implicit AND with no OR support at the
  condition level
- Transformations are a flat list, not composable with the operator
- No type safety: everything is a string
- Cannot express new patterns without adding more YAML structure

## Target State

A rule in the new CRSLang should look like:

```
rule 920170 {
  metadata {
    severity = warning
  }
  when request.method |> matches("^(?:GET|HEAD)$")
   and request.headers["Content-Length"] |> not(matches("^0?$"))
  then block {
    log()
  }
}
```

The phase is not declared — it is inferred as `request_headers` (SecLang phase 1) from
the `request.method` and `request.headers` field references. When exported to SecLang,
the compiler emits `phase:1` automatically.

Rules that use only cross-phase fields (e.g., `tx.*`) must declare the phase explicitly:

```
rule 901001 {
  metadata {
    phase = request_headers
    severity = critical
  }
  when count(tx.crs_setup_version) |> eq(0)
  then deny(status: 500)
}
```

Key properties:

- **Typed fields** with dot-notation and bracket access for maps
- **Phase inference** from field types — no redundant phase metadata
- **Pipeline operator** (`|>`) for chaining transformations and terminal predicates
- **Boolean algebra** (`and`, `or`, `not`, parentheses) replacing chains
- **Structured action blocks** replacing the action bag
- **Compile-time validation** — phase/field mismatches caught before deployment
- **Concise** but readable by security engineers

## Design Principles

1. **Readability first** — CRS rules are read far more than they are written. The syntax
   must be scannable by security engineers who are not language designers.

2. **Composability** — any field expression should be combinable with any other via
   boolean operators. No artificial limits from the underlying engine.

3. **Type safety** — fields have types; functions have signatures. Catch errors at parse
   time, not at runtime.

4. **Backward compatibility** — SecLang import is permanent. YAML remains a supported
   serialization. The new syntax adds a better authoring format, it does not remove
   existing ones.

5. **Incremental adoption** — each phase delivers standalone value. The project does not
   need to complete all phases to be useful.

6. **Engine independence** — CRSLang describes _what_ to match, not _how_ a specific
   engine implements it. Compilation targets (Coraza, ModSecurity, cloud WAFs) are
   separate concerns. Engine configuration (body limits, PCRE tuning, log paths) is
   explicitly out of scope — it belongs in engine-specific config files, not in the rule
   language (see [ADR-0008](adr/0008-configuration-directives.md)).

## Evolution Phases

### Phase 0: Foundational Decisions

Two decisions must be made before any IR or syntax work begins, because they shape
everything downstream.

**Language base:** Should CRSLang adopt an existing language (HCL, with its Terraform
ecosystem, Sprig function library, and zero parser maintenance) or build a custom parser
(with pipeline operator, language-level macros, and full control)? Both approaches
support globals, groups, named function compositions, and boolean conditions. The key
trade-off is maintenance cost vs syntactic control — and the answer determines how
Phases 2-4 are designed.

**Scope:** CRSLang describes _what to detect_, not _how to run the engine_. The ~60
SecLang configuration directives (body limits, PCRE tuning, log paths, etc.) are
explicitly out of scope. Only rule-adjacent metadata (component signature, default
actions, markers, app ID) stays in the language.

**Multi-target compilation:** CRSLang is not constrained by any single compilation
target. Today it compiles to SecLang (ModSecurity/Coraza), but the architecture supports
future backends for Google Cloud Armor (CEL), AWS WAF, Cloudflare (Wirefilter), and
others. SecLang generation must be lossless for the CRS ruleset. Features that a target
cannot express are handled by compiler workarounds or clear error messages.

**Ruleset initialization:** The `crs-setup.conf` + 901 rules two-layer init chain is
replaced by a `config {}` block that holds user-tunable deployment policy (paranoia
level, allowed methods, score thresholds, argument limits, etc.). The compiler generates
all SecLang initialization output — no hand-authored 901 rules, no two-file split. The
deployer copies `setup.crs.example` and edits it, replacing the current
`crs-setup.conf.example` workflow.

See [ADR-0009: Language Base Evaluation](adr/0009-language-base-evaluation.md),
[ADR-0008: Separation of Configuration](adr/0008-configuration-directives.md),
[ADR-0010: Multi-Target Compilation](adr/0010-multi-target-compilation.md), and
[ADR-0012: Ruleset Initialization and Deployment Configuration](adr/0012-ruleset-initialization.md).

### Phase 1: Typed Field System

Replace the split `variables` + `collections` model with a unified typed field namespace.

| Current                                                             | New                               |
| ------------------------------------------------------------------- | --------------------------------- |
| `variables: [REQUEST_METHOD]`                                       | `request.method`                  |
| `collections: [{name: REQUEST_HEADERS, arguments: [Content-Type]}]` | `request.headers["Content-Type"]` |
| `collections: [{name: TX, arguments: [score], count: true}]`        | `count(tx.score)`                 |
| `variables: [REMOTE_ADDR]`                                          | `client.ip`                       |

Work:

- Define a typed field registry mapping dotted names to types (String, Int, IP, Bytes, Map)
- Build bidirectional mappings: SecLang variables/collections <-> field names
- Update the IR to use `Field` instead of separate `Variable`/`Collection`
- YAML backward compat: old format loads and normalizes to new

A key benefit of typed fields: **phase becomes inferable**. Since each field belongs to a
known processing phase (`request.headers` → phase 1, `response.body` → phase 4), the
compiler can derive the phase automatically and catch phase/field mismatches at compile
time. Phase is only required as explicit metadata for rules that use exclusively
cross-phase fields (e.g., `tx.*`).

See [ADR-0001: Typed Field Namespace](adr/0001-typed-field-namespace.md) and
[ADR-0007: Phase Inference](adr/0007-phase-inference.md).

### Phase 2: Boolean Algebra

Replace `chain` actions and implicit logic with explicit boolean expressions.

| Current                           | New                         |
| --------------------------------- | --------------------------- |
| Rule + `chain` + chained rule     | `condition1 and condition2` |
| Multiple conditions (implicit OR) | `condition1 or condition2`  |
| Negated operator                  | `not(condition)`            |

Work:

- Define `Expression` as recursive: `And(Expr, Expr)`, `Or(Expr, Expr)`, `Not(Expr)`,
  `Predicate(Pipeline)`
- Eliminate `chain` as a concept
- Remove `ChainedRule` from the IR
- Parenthesized grouping for precedence

This is the second core IR change. It does not depend on how functions or actions are
surfaced syntactically — boolean composition is a structural change to the condition
model.

See [ADR-0003: Boolean Algebra Replacing Chains](adr/0003-boolean-algebra.md).

### Phase 3: Function Composition and Structured Actions

How functions and actions are expressed depends on the Phase 0 language base decision.

**Function composition** — transformations and operators become composable functions.
If the custom parser is chosen, a pipeline operator (`|>`) chains them left-to-right.
If HCL is chosen, named composition functions (macros) keep nesting shallow. Either way,
reusable transform chains are defined once and invoked by name.

| Current                                                                 | New (custom)                                                | New (HCL)                          |
| ----------------------------------------------------------------------- | ----------------------------------------------------------- | ---------------------------------- |
| `transformations: [lowercase, urlDecode]` + `operator: {name: rx, ...}` | `field \|> url_decode() \|> lowercase() \|> matches("...")` | `matches(normalize(field), "...")` |

**Structured actions and scoring** — the action bag is replaced with a structured model.
If the custom parser is chosen, `then block { tx.score += 5 }` uses native assignment
operators. If HCL is chosen, effects use structured sub-blocks or string-encoded
operations.

Additionally, **anomaly scoring becomes first-class**: severity-derived scoring
eliminates the `setvar` boilerplate from every attack detection rule. A rule just declares
its severity, and the scoring model (defined in globals) handles the rest.

| Current                                        | New                                       |
| ---------------------------------------------- | ----------------------------------------- |
| `disruptive: block` + `setvar: "tx.score=+10"` | `severity: critical` (score auto-derived) |
| Manual `setvar` per category                   | Category derived from group membership    |
| Complex phase-5 evaluation rules               | `scoring_threshold { inbound = 5 }`       |

Work:

- Define a `Function` type with signature: name, args, return type
- If custom: implement pipeline operator, language-level macros
- If HCL: register function library (including Sprig), define macro blocks
- Redesign actions as structured model (disruptive + effects)
- `ctl:` directives become `configure {}` blocks

See [ADR-0002: Pipeline Operator](adr/0002-pipeline-operator.md) (custom parser path),
[ADR-0004: Structured Action Model](adr/0004-structured-actions.md),
[ADR-0011: First-Class Scoring](adr/0011-first-class-scoring.md), and
[ADR-0009: Language Base Evaluation](adr/0009-language-base-evaluation.md).

### Phase 4: Text Syntax and Parser

Implement the text syntax based on the Phase 0 decision.

Both approaches support the core structural requirements:

- **Globals block** for version, common tags — inherited by all rules
- **Group blocks** for file-level metadata — tags shared within a logical grouping
- **Named function compositions** (macros) — reusable transform chains
- **Boolean conditions** with custom functions for matching

Work:

- If HCL: define schema, register function library, design effects model
- If custom: write lexer + recursive-descent parser, build function library
- YAML becomes one serialization of the IR, not the canonical form
- SecLang import stays via existing ANTLR parser (read-only)
- Update WASM/playground for the new syntax

See [ADR-0005: Parser Strategy](adr/0005-parser-strategy.md) and
[ADR-0009: Language Base Evaluation](adr/0009-language-base-evaluation.md).

### Phase 5: Rule Management Directives

Handle meta-operations (exclusions, target updates, action updates) in the new syntax.

```
exclude rule 920170

update rule 920170 {
  remove target request.args["foo"]
}
```

Work:

- First-class exclusion syntax
- Rule override/extension mechanism
- Replace `SecRuleRemoveById`, `SecRuleUpdateTargetById`, etc.

See [ADR-0006: Rule Management Directives](adr/0006-rule-management.md).

## Migration Strategy

```
Phase 0      Decide language base (HCL vs custom) and scope (no engine config).
             These decisions gate all subsequent work.

Phase 1      Internal IR changes: typed fields, phase inference.
             YAML syntax updated but backward-compatible.
             SecLang import still works.

Phase 2      Semantic model change: boolean expressions replace chains.
             New YAML schema version (v2).
             Automated migration tool: v1 YAML -> v2 YAML.

Phase 3      Function composition and action model redesign.
             Shape determined by Phase 0 decision.

Phase 4      Text syntax as primary authoring format.
             Three importers: SecLang -> IR, YAML -> IR, CRSLang text -> IR.
             Two exporters: IR -> YAML, IR -> CRSLang text.

Phase 5      Full language with rule management.
             SecLang becomes legacy import-only.
```

At every phase:

- Existing SecLang `.conf` files remain importable
- Existing v1 YAML remains loadable (with deprecation warnings from Phase 2+)
- The playground supports all active formats
- `crs-toolchain` integration is updated in lockstep

## Architecture Decision Records

| ADR                                          | Phase | Title                                                   | Status   |
| -------------------------------------------- | ----- | ------------------------------------------------------- | -------- |
| [0009](adr/0009-language-base-evaluation.md) | 0     | Language Base — HCL, CEL, Expr, or Custom               | Proposed |
| [0008](adr/0008-configuration-directives.md) | 0     | Separation of Configuration from Rule Language          | Proposed |
| [0010](adr/0010-multi-target-compilation.md) | 0     | Multi-Target Compilation Model                          | Proposed |
| [0012](adr/0012-ruleset-initialization.md)   | 0     | Ruleset Initialization and Deployment Configuration     | Proposed |
| [0001](adr/0001-typed-field-namespace.md)    | 1     | Typed Field Namespace                                   | Proposed |
| [0007](adr/0007-phase-inference.md)          | 1     | Phase Inference from Field Types                        | Proposed |
| [0003](adr/0003-boolean-algebra.md)          | 2     | Boolean Algebra Replacing Chains                        | Proposed |
| [0002](adr/0002-pipeline-operator.md)        | 3     | Pipeline Operator for Composition (conditional on 0009) | Proposed |
| [0004](adr/0004-structured-actions.md)       | 3     | Structured Action Model                                 | Proposed |
| [0011](adr/0011-first-class-scoring.md)      | 3     | First-Class Scoring Model                               | Proposed |
| [0005](adr/0005-parser-strategy.md)          | 4     | Parser Strategy (see 0009 for decision)                 | Proposed |
| [0006](adr/0006-rule-management.md)          | 5     | Rule Management Directives                              | Proposed |

## Reference Languages

- [Wirefilter](https://developers.cloudflare.com/ruleset-engine/rules-language/) —
  Cloudflare's flat-field boolean expression language for firewall rules
- [CEL](https://github.com/google/cel-spec) — Google's Common Expression Language,
  used in Envoy, Kubernetes, and IAM policies
- [OPA/Rego](https://www.openpolicyagent.org/docs/latest/policy-language/) —
  Open Policy Agent's logic-programming policy language
- [SecLang](https://github.com/owasp-modsecurity/ModSecurity/wiki/Reference-Manual-%28v3.x%29) —
  ModSecurity's directive-based rule language (current source format)
