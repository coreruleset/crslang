# ADR-0003: Boolean Algebra Replacing Chains

- **Status:** Proposed
- **Date:** 2026-04-13
- **Phase:** 2

## Context

SecLang uses the `chain` action to create AND logic between rules. When a rule has a
`chain` action, the next rule in sequence becomes a chained condition — the second rule
only evaluates if the first matches. This is how SecLang expresses multi-condition logic:

```
SecRule REQUEST_METHOD "@rx ^(?:GET|HEAD)$" \
    "id:920170,phase:1,block,chain"
    SecRule REQUEST_HEADERS:Content-Length "!@rx ^0?$" \
        "setvar:'tx.anomaly_score=+5'"
```

CRSLang v1 models this directly:

```yaml
kind: rule
metadata:
  id: 920170
conditions:
  - variables: [REQUEST_METHOD]
    operator: {name: rx, value: "^(?:GET|HEAD)$"}
actions:
  disruptive: {action: block}
  flow: [chain]
chainedRule:
  kind: rule
  conditions:
    - collections:
        - name: REQUEST_HEADERS
          arguments: [Content-Length]
      operator: {negate: true, name: rx, value: "^0?$"}
  actions:
    non-disruptive:
      - action: setvar
        param: "tx.anomaly_score=+5"
```

### Problems with chains

1. **Only AND** — chains can only express AND. There is no native OR at the condition
   level (only multiple variables in a single `SecRule`, which is a limited form of OR).
2. **Sequential nesting** — deeply chained rules are hard to read and maintain. A 4-way
   AND requires 3 levels of nesting.
3. **Side-effects at intermediate links** — chained rules can have their own
   transformations, actions (setvar, log), and even different targets. This conflates
   "condition" with "intermediate action."
4. **Fragile ordering** — the chain is positional. Inserting a rule between chain links
   breaks the logic.
5. **No grouping** — you cannot express `(A AND B) OR (C AND D)` in SecLang. CRS works
   around this with separate rules and anomaly scoring.

## Decision

Replace chains with **explicit boolean algebra**: `and`, `or`, `not`, with parenthesized
grouping.

### Syntax

```
rule 920170 (phase: request) {
  when request.method |> matches("^(?:GET|HEAD)$")
   and request.headers["Content-Length"] |> not(matches("^0?$"))
  then block {
    tx.anomaly_score += 5
  }
}
```

### Operator Precedence

From highest to lowest:
1. `not` (unary)
2. `and`
3. `or`

Parentheses override precedence:

```
when (A or B) and (C or D)
```

### Expression Tree IR

```go
type Expr interface {
    exprNode()
}

type AndExpr struct {
    Left  Expr
    Right Expr
}

type OrExpr struct {
    Left  Expr
    Right Expr
}

type NotExpr struct {
    Inner Expr
}

type PredicateExpr struct {
    Pipeline Pipeline  // from ADR-0002
}
```

A rule has exactly one `Expr` as its condition:

```go
type Rule struct {
    ID       int
    Metadata RuleMetadata
    When     Expr           // single expression tree
    Then     ActionBlock    // from ADR-0004
}
```

### Handling Chain Side-Effects

#### Current Semantics in CRSLang v1

The existing normalizer (`types/condition_directives.go:117-128`) reveals how CRSLang
already handles chains:

- **Chain links without non-disruptive actions** are **flattened** — their conditions are
  merged into a single AND group with the next link.
- **Chain links with non-disruptive actions** (setvar, capture, ctl, etc.) are **kept as
  separate `ChainedRule` nodes** — preserving the side-effect boundary.

This means the IR already distinguishes between "pure condition" chain links and "chain
links with side-effects." The boolean algebra model must account for both.

#### Analysis of CRS Chain Patterns

An audit of real CRS rules reveals three categories of intermediate chain side-effects:

**Category 1: `capture` on the first link (used by the final link)**

CRS rule 920190 (`REQUEST-920-PROTOCOL-ENFORCEMENT.conf`) uses `capture` on the first
link to extract regex groups, which are then referenced by the chained link's actions:

```
SecRule REQUEST_HEADERS:Range "@rx (\d+)-(\d+)" \
    "id:920190,phase:1,block,capture,chain"
    SecRule REQUEST_METHOD "@streq GET" \
        "setvar:'tx.inbound_anomaly_score_pl1=+%{tx.critical_anomaly_score}'"
```

Here `capture` must execute when the first condition matches (not when the full chain
matches), because captured groups are consumed downstream. However, the captured values
are only *useful* if the full chain matches. In the boolean model, `capture()` becomes a
pipeline function that extracts groups as a side-effect of matching:

```
when request.headers["Range"] |> capture("(\\d+)-(\\d+)")
 and request.method |> eq("GET")
then block { tx.inbound_anomaly_score_pl1 += tx.critical_anomaly_score }
```

The `capture()` function both tests the regex and stores the match groups. It is a
predicate with a side-effect — similar to how `matches()` works, but additionally
populating the match context.

**Category 2: `ctl` on intermediate links (engine configuration)**

CRS rule 905111 (`testdata/test_36_chain.conf`) uses `ctl:ruleRemoveByTag` on an
intermediate link in a 3-level chain:

```
SecRule REMOTE_ADDR "@ipMatch 127.0.0.1,::1" \
    "id:905111,phase:1,pass,chain"
    SecRule REQUEST_HEADERS:User-Agent "@endsWith (internal dummy connection)" \
        "ctl:ruleRemoveByTag=TEST_CHAIN,chain"
        SecRule REQUEST_LINE "@rx ^(?:GET /|OPTIONS \*) HTTP/[12]\.[01]$" \
            "ctl:ruleRemoveByTag=OWASP_CRS,ctl:auditEngine=Off"
```

The intermediate `ctl` modifies engine state only if that specific link matches. In
ModSecurity, `ctl` actions on chain links execute per-link, not deferred to full-chain
match. This is the hardest pattern to model in a pure boolean expression.

**Category 3: `setvar` on the first link (consumed by a later link)**

CRS rule 901320 (`REQUEST-901-INITIALIZATION.conf`) sets a TX variable on the first link,
then the final link uses that variable in a transformation:

```
SecRule &TX:ENABLE_DEFAULT_COLLECTIONS "@eq 1" \
    "id:901320,phase:1,pass,nolog,\
    setvar:'tx.ua_hash=%{REQUEST_HEADERS.User-Agent}',chain"
    SecRule TX:ENABLE_DEFAULT_COLLECTIONS "@eq 1" \
        "chain"
        SecRule TX:ua_hash "@unconditionalMatch" \
            "t:none,t:sha1,t:hexEncode,\
            initcol:global=global,\
            initcol:ip=%{remote_addr}_%{MATCHED_VAR}"
```

Here, `setvar` on link 1 creates `tx.ua_hash`, and link 3 reads it. The side-effect
must execute before the downstream link evaluates. This is fundamentally a **data-flow
dependency between conditions**, not just a side-effect.

#### Design Options

**Option A: Side-effects only on the rule (move to `then` block)**

All setvar/log/ctl actions move to the `then` block and execute only when the full
expression matches.

```
rule 920190 (phase: request) {
  when request.headers["Range"] |> matches("(\\d+)-(\\d+)")
   and request.method |> eq("GET")
  then block {
    capture()
    tx.inbound_anomaly_score_pl1 += tx.critical_anomaly_score
  }
}
```

- Works for **most** CRS rules (the majority of chains have side-effects only on the
  final link).
- **Breaks Category 2** (intermediate `ctl`): if `ctl` is deferred, the engine
  configuration change happens too late.
- **Breaks Category 3** (data-flow `setvar`): if `setvar` moves to `then`, the
  downstream condition cannot read the variable.

**Option B: Conditional side-effects on sub-expressions**

Allow inline effect blocks on individual predicates:

```
rule 905111 (phase: request) {
  when client.ip |> ip_in_range("127.0.0.1", "::1")
   and request.headers["User-Agent"] |> ends_with("(internal dummy connection)") {
         configure(rule_remove_by_tag: "TEST_CHAIN")
       }
   and request.line |> matches("^(?:GET /|OPTIONS \\*) HTTP/[12]\\.[01]$") {
         configure(rule_remove_by_tag: "OWASP_CRS", audit_engine: off)
       }
  then pass
}
```

- Handles all three categories.
- Adds complexity: the reader must understand that `{ ... }` after a predicate is a
  conditional side-effect, not part of the boolean expression.
- Requires defining execution semantics precisely (see below).

**Option C: `let` bindings for data-flow dependencies**

Introduce variable bindings to handle Category 3 without inline side-effects:

```
rule 901320 (phase: request) {
  let ua_hash = request.headers["User-Agent"] |> sha1() |> hex_encode()
  when count(tx.enable_default_collections) |> eq(1)
   and tx.enable_default_collections |> eq(1)
  then pass {
    init_collection(global: "global")
    init_collection(ip: client.ip + "_" + ua_hash)
  }
}
```

- Separates data flow from conditions cleanly.
- `let` bindings are evaluated eagerly, before the `when` clause.
- Handles Category 3 without inline side-effects.
- Does not address Category 2 (`ctl` on intermediate links).

#### Recommendation

Use a **layered approach**:

1. **Phase 2a** — implement Option A (side-effects in `then` only). This handles the
   vast majority of CRS rules. Chain links without intermediate side-effects are already
   flattened by the existing normalizer, so this is the natural starting point.

2. **Phase 2b** — add `let` bindings (Option C) for data-flow dependencies. This
   cleanly handles Category 3 without complicating the boolean expression model.

3. **Phase 2c** — if Category 2 (`ctl` on intermediate links) proves common enough to
   warrant language support, add conditional side-effects (Option B) as an extension.
   Before doing so, audit whether these `ctl` patterns can be restructured as separate
   rules instead.

#### Execution Semantics (if Option B is adopted)

If conditional side-effects are added:
- Side-effects on a predicate execute **when that predicate evaluates to true** during
  expression evaluation.
- Evaluation order is **left-to-right** with **short-circuit**: `A and B` does not
  evaluate `B` if `A` is false; `A or B` does not evaluate `B` if `A` is true.
- Side-effects from predicates that are not evaluated (due to short-circuit) do **not**
  execute.
- This matches ModSecurity's chain behavior where a chain link's actions only fire if
  that link matches.

### Migration from Chains

The conversion is mechanical for most cases:

```
# Chain of depth N (no intermediate side-effects):
rule + chain -> chained1 + chain -> chained2
# Becomes:
when condition1 and condition2 and condition3

# Negated chain link:
rule + chain -> (negated operator in chained rule)
# Becomes:
when condition1 and not(condition2)
```

For chains with intermediate side-effects, by category:

**Category 1 (`capture`):** Model as a pipeline function. `capture()` replaces both the
`@rx` operator and the `capture` action — it matches and extracts in one step:

```
# SecLang:
SecRule REQUEST_HEADERS:Range "@rx (\d+)-(\d+)" "capture,chain"
# CRSLang:
when request.headers["Range"] |> capture("(\\d+)-(\\d+)") and ...
```

**Category 2 (`ctl` on intermediate links):** Restructure as separate rules where
possible. If not possible and Option B is adopted, use conditional side-effects:

```
# SecLang:
SecRule REMOTE_ADDR "@ipMatch 127.0.0.1" "chain"
    SecRule REQUEST_HEADERS:User-Agent "@endsWith foo" "ctl:ruleRemoveByTag=X,chain"
# CRSLang (Option B):
when client.ip |> ip_in_range("127.0.0.1")
 and request.headers["User-Agent"] |> ends_with("foo") {
       configure(rule_remove_by_tag: "X")
     }
 and ...
```

**Category 3 (`setvar` consumed downstream):** Use `let` bindings:

```
# SecLang:
SecRule &TX:ENABLE "@eq 1" "setvar:'tx.ua_hash=%{REQUEST_HEADERS.User-Agent}',chain"
    SecRule TX:ua_hash "@unconditionalMatch" "t:sha1,t:hexEncode,initcol:ip=..."
# CRSLang:
let ua_hash = request.headers["User-Agent"] |> sha1() |> hex_encode()
when count(tx.enable) |> eq(1) ...
then pass { init_collection(ip: client.ip + "_" + ua_hash) }
```

### Collection Quantifier: `each()`

SecLang's `multiMatch` action changes how a collection-targeting condition evaluates —
instead of stopping at the first match, it iterates all values and fires side-effects
per match. This is currently modeled as a non-disruptive action, but it is semantically
a condition quantifier.

**Recommendation: `each()` as a condition-level quantifier (Option A).**

```
# Without each(): first match wins, effects fire once
when request.args |> detect_sqli()

# With each(): all values tested, effects fire per match
when each(request.args) |> detect_sqli()
then block {
  tx.sqli_score += 5     # incremented per matching argument
  log(data: matched.var)  # logged per matching argument
}
```

`each()` wraps a map-typed field and signals "iterate all values." Without it, the
default is first-match semantics.

**Alternatives documented:**

- **Option B: Effect-level modifier** — `then block (per_match: true) { ... }`. Simpler
  to parse but misleading: the reader assumes first-match from the condition until
  they notice the modifier.
- **Option C: Separate iteration block** — `for each match { per-match effects } then
  block { once-only effects }`. Most expressive (supports both per-match and once-only
  effects) but adds a new block type.
- **Option D: Drop it** — if scoring becomes first-class (ADR-0011), per-match scoring
  may be handled at the scoring level rather than as a language construct.

`multiMatch` is rarely used in CRS, so Option A is sufficient for the foreseeable
future. Options B/C can be revisited if use cases emerge.

### String Interpolation

SecLang uses `%{TX:score}` and `%{MATCHED_VAR}` for string interpolation in actions
(`logdata`, `msg`, `setvar`). Most of these cases become direct field references or
expressions in CRSLang (e.g., `log(data: matched.var)`).

For cases that require composed strings (log messages, dynamic values), CRSLang needs
a string construction mechanism. The exact form — string interpolation
(`"Score: ${tx.anomaly_score}"`), concatenation (`"Score: " + string(tx.score)`), or
a format function (`format("Score: %d", tx.score)`) — is deferred to the effects model
design in Phase 3. The IR must support composed string values in effect arguments.

### New Capabilities

Boolean algebra enables patterns that are impossible or awkward in SecLang:

```
# OR conditions (currently requires separate rules + scoring)
rule 100001 (phase: request) {
  when request.uri |> matches("/admin")
   and (client.ip |> ip_in_range("10.0.0.0/8")
        or request.headers["X-Internal"] |> eq("true"))
  then pass
}

# Complex grouping
rule 100002 (phase: request) {
  when (request.method |> eq("POST") or request.method |> eq("PUT"))
   and not(request.headers["Content-Type"] |> contains("application/json"))
   and request.body |> length() |> gt(0)
  then block
}
```

### YAML v2 Representation

```yaml
kind: rule
metadata:
  id: 920170
  phase: request
when:
  and:
    - pipeline:
        field: request.method
        steps:
          - fn: matches
            args: ["^(?:GET|HEAD)$"]
    - not:
        pipeline:
          field: request.headers["Content-Length"]
          steps:
            - fn: matches
              args: ["^0?$"]
then:
  action: block
  effects:
    - set: tx.anomaly_score
      op: "+="
      value: 5
```

## Alternatives Considered

### A: Keep chains, add OR

Add `or_chain` alongside `chain` to keep the sequential model but add OR support.

**Rejected because:**
- Does not address nesting depth problem
- Creates a hybrid model that is harder to reason about than either pure chains or
  pure boolean algebra
- Still positional/fragile

### B: Implicit AND (Rego-style)

Multiple conditions in a block are implicitly ANDed:

```
rule 920170 {
  request.method |> matches("^(?:GET|HEAD)$")
  request.headers["Content-Length"] |> not(matches("^0?$"))
  -> block { ... }
}
```

**Rejected because:**
- Implicit AND with no explicit OR creates the same asymmetry as SecLang
- Harder to parse: where does one condition end and the next begin?
- Less familiar to the target audience

### C: Pattern Matching

```
match request {
  method: /^(?:GET|HEAD)$/,
  headers["Content-Length"]: not(/^0?$/)
} -> block
```

**Rejected because:**
- Only works for AND conditions on the same request
- Cannot express cross-category conditions (request + response + tx)
- Novel syntax with no established precedent

## Consequences

### Positive

- Full boolean expressiveness: AND, OR, NOT with arbitrary nesting
- Eliminates chain complexity and positional fragility
- Rules become self-contained: no hidden dependencies on adjacent rules
- Enables optimizations: expression tree can be reordered, short-circuited
- Closer to how security engineers think about conditions

### Negative

- Three categories of intermediate side-effects require a layered migration (Phase 2a/b/c)
  rather than a single clean cutover
- `let` bindings (Phase 2b) add a new language concept not present in SecLang
- Conditional side-effects (Phase 2c, if adopted) complicate the expression model and
  require precise execution semantics
- Some deeply chained rules may become long single expressions (mitigated by
  line breaks and formatting conventions)

### Risks

- **Side-effect ordering** — if Option B (conditional side-effects) is adopted, the
  semantics of short-circuit evaluation must be precisely defined: side-effects on
  unevaluated predicates do not fire. This matches ModSecurity chain behavior but must
  be documented and tested.
- **Engine mapping** — not all WAF engines support full boolean algebra natively. The
  compiler may need to decompose expressions back into chains for some targets
  (ModSecurity, Coraza). This is a compiler concern, not a language concern.
- **Category 2 prevalence** — intermediate `ctl` actions (Category 2) are rare in CRS
  (found primarily in initialization and internal-traffic rules like 905111). If they
  prove confined to a small set of rules, they can be handled by restructuring those
  rules rather than adding Option B to the language. A full CRS audit should quantify
  this before committing to Phase 2c.
- **Category 3 data flow** — `let` bindings change CRSLang from a purely declarative
  rule language to one with local variable scoping. This is a significant conceptual
  shift. The alternative is to require these patterns to be split into multiple rules
  with explicit TX variables, which is closer to how CRS already works.
- **Migration correctness** — automated chain-to-boolean conversion must be validated
  against the full CRS test suite. The three categories need separate migration paths
  in the converter tool.
