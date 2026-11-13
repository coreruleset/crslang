# CRSLang Evolution: From YAML/SecLang to Composable Expressions

## Vision

CRSLang aims to become a modern, composable rule language for web application firewalls.
Today it is a YAML serialization of ModSecurity's SecLang AST. The goal is to evolve it
into a standalone expression language where conditions are composed from typed fields,
transformed through composable functions, and combined with boolean algebra — similar in
spirit to Cloudflare's Wirefilter or Google's CEL, but purpose-built for the CRS
ecosystem.

## Status

This document and the linked ADRs describe a design that is mostly still **proposed**.
Two Phase 0 decisions are accepted: CRSLang uses a **custom parser** as its language base
([ADR-0009](adr/0009-language-base-evaluation.md), Option D), and engine configuration is
out of the language's scope ([ADR-0008](adr/0008-configuration-directives.md)). Every
other ADR is in `Status: Proposed`. For those, the framing is "the design proposes X"
rather than "X is decided" — even where the proposed direction reads declaratively for
clarity.

The remaining Phase 0 ADRs (multi-target compilation, ruleset initialization) are
proposals with reasonable consensus but are still subject to revision.

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
rule 949110 {
  metadata {
    phase = logging
    severity = critical
  }
  when tx.inbound_anomaly_score |> gt(tx.inbound_anomaly_score_threshold)
  then deny(status: 403)
}
```

Key properties:
- **Typed fields** with dot-notation and bracket access for maps
- **Phase inference** from field types — no redundant phase metadata
- **A `metadata {}` block** for rule properties (`phase`, `severity`, `paranoia`,
  `multi_match`, `backends`, ...) instead of arguments after the rule ID
- **Composable functions** for chaining transformations and terminal predicates, joined
  with the pipeline operator `|>` and reusable as named macros
- **Boolean algebra** (`and`, `or`, `not`, parentheses) replacing chains, with `let` to
  name conditions
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

6. **Engine independence** — CRSLang describes *what* to match, not *how* a specific
   engine implements it. Compilation targets (Coraza, ModSecurity, cloud WAFs) are
   separate concerns. The proposal in [ADR-0008](adr/0008-configuration-directives.md)
   places engine configuration (body limits, PCRE tuning, log paths) outside the rule
   language, in engine-specific config files.

## Evolution Phases

### Phase 0: Foundational Decisions

Four foundational topics need to be settled before IR or syntax work begins, because
they shape everything downstream. The first two are decided; the other two have
proposed directions that are still subject to revision.

**Language base (accepted):** CRSLang builds a custom parser (ADR-0009, Option D) instead
of adopting HCL, CEL or Expr. It gets the pipeline operator, macros defined in the
language, native `when`/`then` structure with imperative effects, and full control over
how the language evolves. The cost is maintaining a small parser and building editor
support and the function library from scratch.

**Scope (accepted):** CRSLang describes *what to detect*, not *how to run the engine*.
The ~60 SecLang configuration directives (body limits, PCRE tuning, log paths, etc.)
move out of scope; only rule-adjacent metadata (component signature, default actions,
markers, app ID) stays in the language.

**Multi-target compilation (proposed):** CRSLang would not be constrained by any single
compilation target. Today it compiles to SecLang (ModSecurity/Coraza); the proposed
architecture leaves room for future backends targeting Google Cloud Armor (CEL), AWS WAF,
Cloudflare (Wirefilter), and others. SecLang generation must remain lossless for the CRS
ruleset. Features that a target cannot express are handled by compiler workarounds or
clear error messages.

**Ruleset initialization (proposed):** The `crs-setup.conf` + 901 rules two-layer init
chain would be replaced by a `config {}` block that holds user-tunable deployment policy
(paranoia level, allowed methods, score thresholds, argument limits, etc.). The compiler
would generate all SecLang initialization output — no hand-authored 901 rules, no
two-file split. The deployer would copy `setup.crs.example` and edit it, replacing the
current `crs-setup.conf.example` workflow.

See [ADR-0009: Language Base Evaluation](adr/0009-language-base-evaluation.md),
[ADR-0008: Separation of Configuration](adr/0008-configuration-directives.md),
[ADR-0010: Multi-Target Compilation](adr/0010-multi-target-compilation.md), and
[ADR-0012: Ruleset Initialization and Deployment Configuration](adr/0012-ruleset-initialization.md).

### Phase 1: Typed Field System

Replace the split `variables` + `collections` model with a unified typed field namespace.

| Current | New |
|---|---|
| `variables: [REQUEST_METHOD]` | `request.method` |
| `collections: [{name: REQUEST_HEADERS, arguments: [Content-Type]}]` | `request.headers["Content-Type"]` |
| `collections: [{name: TX, arguments: [score], count: true}]` | `count(tx.score)` |
| `variables: [REMOTE_ADDR]` | `client.ip` |

Work:
- Define a typed field registry mapping dotted names to types (String, Int, IP, Bytes, Map)
- Build bidirectional mappings: SecLang variables/collections <-> field names
- Update the IR to use `Field` instead of separate `Variable`/`Collection`
- YAML backward compat: old format loads and normalizes to new

A key benefit of typed fields: **phase becomes inferable**. Since each field belongs to a
known processing phase (`request.headers` → phase 1, `response.body` → phase 4), the
compiler can derive the phase automatically: a rule runs in the latest phase among the
fields it reads, since earlier data stays available later. Declaring a phase earlier than
a referenced field is a compile-time error. Phase is only required as explicit metadata
for rules that use exclusively cross-phase fields (e.g., `tx.*`).

Phase 1 also introduces a structured **documentation model** spanning four tiers
(ruleset, group, rule, rule-management directive), with doc-comments for narrative
prose and typed metadata fields for machine-readable content (references, FP notes,
CWE/OWASP tags). Groups nest for conceptual sub-grouping; cross-cutting concepts like
paranoia and severity get their own documentation blocks.

See [ADR-0001: Typed Field Namespace](adr/0001-typed-field-namespace.md),
[ADR-0007: Phase Inference](adr/0007-phase-inference.md), and
[ADR-0013: Documentation Model](adr/0013-documentation-model.md).

### Phase 2: Boolean Algebra

Replace `chain` actions and implicit logic with explicit boolean expressions.

| Current | New |
|---|---|
| Rule + `chain` + chained rule | `condition1 and condition2` |
| Multiple conditions (implicit OR) | `condition1 or condition2` |
| Negated operator | `not(condition)` |

Work:
- Define `Expression` as recursive: `And(Expr, Expr)`, `Or(Expr, Expr)`, `Not(Expr)`,
  `Predicate(Pipeline)`
- Eliminate `chain` as a concept
- Remove `ChainedRule` from the IR
- Parenthesized grouping for precedence

This is the second core IR change. It does not depend on how functions or actions are
surfaced syntactically — boolean composition is a structural change to the condition
model. When compiling to SecLang, which has no OR, each `or` branch becomes a rule that
sets a flag, and one rule with the authored ID tests the flag and applies the action and
effects once ([ADR-0010](adr/0010-multi-target-compilation.md)).

See [ADR-0003: Boolean Algebra Replacing Chains](adr/0003-boolean-algebra.md).

### Phase 3: Function Composition and Structured Actions

**Function composition** — transformations and operators become composable functions,
chained left-to-right with the pipeline operator (`|>`). Reusable transform chains are
defined once as macros and invoked by name.

| Current | New |
|---|---|
| `transformations: [lowercase, urlDecode]` + `operator: {name: rx, ...}` | `field \|> url_decode() \|> lowercase() \|> matches("...")` |

**Structured actions and scoring** — the action bag is replaced with a structured model.
Effects use native assignment operators, as in `then block { tx.score += 5 }`.

Additionally, **anomaly scoring becomes first-class**: severity-derived scoring
eliminates the `setvar` boilerplate from every attack detection rule. A rule just declares
its severity, and the scoring model (defined in globals) handles the rest. Adding to a
score variable by hand in a rule that declares a severity is a compile-time error, since
it would count the score twice.

| Current | New |
|---|---|
| `disruptive: block` + `setvar: "tx.score=+10"` | `severity: critical` (score auto-derived) |
| Manual `setvar` per category | Category derived from group membership (`category = "sql_injection"` → `tx.sql_injection_score`) |
| Non-standard score `setvar`s | `score(anomaly: 5, protocol: 3)` replaces the severity-derived score |
| Complex phase-5 evaluation rules | `config { scoring { inbound_threshold = 5 } }` |

Work:
- Define a `Function` type with signature: name, args, return type
- Implement the pipeline operator and language-level macros
- Redesign actions as structured model (disruptive + effects)
- `ctl:` directives become `configure {}` blocks

Phase 3 also introduces **macros** for reusable expression compositions (eliminating
copy-pasted normalization chains and predicate patterns) and **external data sources**
as first-class typed declarations (IP lists, regex sets, geo databases, regex assemblies).

See [ADR-0002: Pipeline Operator](adr/0002-pipeline-operator.md),
[ADR-0004: Structured Action Model](adr/0004-structured-actions.md),
[ADR-0011: First-Class Scoring](adr/0011-first-class-scoring.md),
[ADR-0015: Macros and Named Compositions](adr/0015-macros-and-compositions.md),
[ADR-0016: External Data Sources](adr/0016-external-data-sources.md), and
[ADR-0009: Language Base Evaluation](adr/0009-language-base-evaluation.md).

### Phase 4: Text Syntax and Parser

Implement the custom parser chosen in Phase 0 (ADR-0009). The syntax supports the core
structural requirements:
- **Globals block** for version, common tags — inherited by all rules
- **Group blocks** for file-level metadata — tags shared within a logical grouping
- **Named function compositions** (macros) — reusable transform chains
- **Boolean conditions** with custom functions for matching

Work:
- Write the lexer and recursive-descent parser, build the function library
- YAML becomes one serialization of the IR, not the canonical form
- SecLang import stays via existing ANTLR parser (read-only)
- Update WASM/playground for the new syntax

Phase 4 also formalizes the **import and modularity** model: explicit `import`
statements, hierarchical namespaces for groups/macros/data, and a single flat namespace
for rule IDs. Multi-file rulesets gain predictable composition semantics, and packages
(CRS itself, third-party rule packs) become first-class distribution units.

See [ADR-0005: Parser Strategy](adr/0005-parser-strategy.md),
[ADR-0014: Imports and Modularity](adr/0014-imports-and-modularity.md), and
[ADR-0009: Language Base Evaluation](adr/0009-language-base-evaluation.md).

### Phase 5: Rule Management Directives

Handle meta-operations (exclusions, target updates, action updates) in the new syntax.

```
exclude rule 920170

update rule 920170 {
  remove target request.args["foo"]
}

# Only for some requests
exclude rule 942100 target request.body.args["query"]
  when request.uri |> begins_with("/admin/query-builder")
```

`where` selects rules by their metadata (`exclude rules where tag == "..."`); `when` limits
an exclusion to matching requests and compiles to SecLang `ctl:` runtime exclusions.

Work:
- First-class exclusion syntax
- Rule override/extension mechanism
- Replace `SecRuleRemoveById`, `SecRuleUpdateTargetById`, etc.

See [ADR-0006: Rule Management Directives](adr/0006-rule-management.md).

## Migration Strategy

```
Phase 0      Settle the four foundational topics: language base (decided:
             custom parser), scope separation (decided: engine config out of
             the language), multi-target compilation model, and ruleset
             initialization (config {} block replacing crs-setup.conf + 901
             rules). These decisions gate all subsequent work.

Phase 1      Internal IR changes: typed fields, phase inference.
             YAML syntax updated but backward-compatible.
             SecLang import still works.

Phase 2      Semantic model change: boolean expressions replace chains.
             New YAML schema version (v2).
             Automated migration tool: v1 YAML -> v2 YAML.

Phase 3      Function composition and action model redesign.
             Pipeline operator, macros, structured effects.

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

| ADR | Phase | Title | Status |
|-----|-------|-------|--------|
| [0009](adr/0009-language-base-evaluation.md) | 0 | Language Base — HCL, CEL, Expr, or Custom | Accepted (custom parser) |
| [0008](adr/0008-configuration-directives.md) | 0 | Separation of Configuration from Rule Language | Accepted |
| [0010](adr/0010-multi-target-compilation.md) | 0 | Multi-Target Compilation Model | Proposed |
| [0012](adr/0012-ruleset-initialization.md) | 0 | Ruleset Initialization and Deployment Configuration | Proposed |
| [0001](adr/0001-typed-field-namespace.md) | 1 | Typed Field Namespace | Proposed |
| [0007](adr/0007-phase-inference.md) | 1 | Phase Inference from Field Types | Proposed |
| [0013](adr/0013-documentation-model.md) | 1 | Documentation Model for Rules, Groups, and Rulesets | Proposed |
| [0003](adr/0003-boolean-algebra.md) | 2 | Boolean Algebra Replacing Chains | Proposed |
| [0002](adr/0002-pipeline-operator.md) | 3 | Pipeline Operator for Composition | Proposed |
| [0004](adr/0004-structured-actions.md) | 3 | Structured Action Model | Proposed |
| [0011](adr/0011-first-class-scoring.md) | 3 | First-Class Scoring Model | Proposed |
| [0015](adr/0015-macros-and-compositions.md) | 3 | Macros and Named Compositions | Proposed |
| [0016](adr/0016-external-data-sources.md) | 3 | External Data Sources | Proposed |
| [0005](adr/0005-parser-strategy.md) | 4 | Parser Strategy | Proposed |
| [0014](adr/0014-imports-and-modularity.md) | 4 | Imports and Modularity | Proposed |
| [0006](adr/0006-rule-management.md) | 5 | Rule Management Directives | Proposed |

For structural concepts identified but not yet covered by an ADR (time/duration types,
versioning annotations, function signatures, test attachment, persistent collections,
capture groups, and others), see [TODO.md](TODO.md).

## Reference Languages

- [Wirefilter](https://developers.cloudflare.com/ruleset-engine/rules-language/) —
  Cloudflare's flat-field boolean expression language for firewall rules
- [CEL](https://github.com/google/cel-spec) — Google's Common Expression Language,
  used in Envoy, Kubernetes, and IAM policies
- [OPA/Rego](https://www.openpolicyagent.org/docs/latest/policy-language/) —
  Open Policy Agent's logic-programming policy language
- [SecLang](https://github.com/owasp-modsecurity/ModSecurity/wiki/Reference-Manual-%28v3.x%29) —
  ModSecurity's directive-based rule language (current source format)
