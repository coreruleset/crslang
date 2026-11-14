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

Examples here use the pipeline (ADR-0002) and action block (ADR-0004) syntax for
readability. The boolean structure this ADR adds to the IR does not depend on either.

```
rule 920170 {
  when request.method |> matches("^(?:GET|HEAD)$")
   and request.headers["Content-Length"] |> not(matches("^0?$"))
  then block {
    tx.anomaly_score += 5
  }
}
```

### Naming Conditions with `let`

Long pipelines inside a boolean expression are hard to read: the reader has to find
where each pipeline ends and the next operator begins. `let` names a predicate so the
`when` clause shows only the boolean structure:

```
rule 920170 {
  let is_get_or_head  = request.method |> matches("^(?:GET|HEAD)$")
  let has_body_length = request.headers["Content-Length"] |> not(matches("^0?$"))
  when is_get_or_head and has_body_length
  then block {
    tx.anomaly_score += 5
  }
}
```

`let` is syntactic sugar. The parser inlines each binding into the `Expr` tree where it
is used, so it needs no IR node and no YAML v2 form, and short-circuit evaluation is
unchanged. Pipelines are pure, so inlining a binding gives the same result as evaluating
it first. A binding can also name a value rather than a predicate (`let referer_host =
request.headers["Referer"] |> extract(...)`), which is how a rule tests part of a field
(see Conditions Are Pure).

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

### Conditions Are Pure

A `when` expression has no side-effects. Effects (`setvar`, `ctl`, scoring, logging)
live only in the `then` block and run once, when the whole expression is true. Because
nothing in a condition can be observed, evaluation order and short-circuiting are not
part of the language: the compiler may reorder predicates (for example, cheap ones
first, or the OR lowering in ADR-0010) without changing what a rule means.

This replaces chain links that carry actions. In CRS those actions only move data from
one link to the next, so pure value functions can express all of them.

#### What CRS Chains Do Between Links

An audit of CRS v4.30.0 (`rules/*.conf`) and 11 CRS plugins (`plugins/*.conf`, mostly
rule-exclusion plugins) found:

- 54 chained rules in CRS and 68 in the plugins.
- In CRS, 24 rules carry actions on a link before the last one. Those actions are only
  `capture` and `setvar`.
- No `ctl` on a link before the last one, in CRS or in the plugins. The usual example,
  905111, is test data (`testdata/test_36_chain.conf`), not a CRS rule. The plugins
  add 2 `capture`, 1 `setvar`, and 1 `skip` on a chain starter, which applies to the
  whole chain.

The actions follow three patterns:

| Pattern | CRS rules | What the earlier link does |
|---|---|---|
| Filter matched values | 932200, 932206, 932240, 933150, 941310, 942200, 942440, 954130 | A later link tests `MATCHED_VARS`/`MATCHED_VAR`: "of the values that matched A, one also matches B" |
| Extract, then test | 901320, 920190, 920420, 920440, 920450, 920451, 920480, 922100, 931130, 932205, 932207, 942130, 942131, 942521, 943110 | `capture` (and often `setvar tx.x=%{tx.1}`) pulls out a part of a value (content type, charset, extension, header name, Referer host), and a later link tests `TX:1` or `TX:x` |
| Keep a value for the log | 932200, 932205, 932206, 932207, 932240, 942130, 942131, 942521, 943110, 943120 | `setvar tx.<id>_matched_var_name=%{matched_var_name}` on the first link, because `logdata` on the last link would otherwise see the last link's match |

#### Collection Semantics

Two rules, matching SecLang, make the patterns expressible:

- A transformation applied to a collection (`request.args`, `request.headers`) applies
  to each value. A pipeline step `not(p)` also applies to each value.
- A predicate applied to a collection is true when **at least one** value satisfies it.
  `request.args |> matches(A) and request.args |> matches(B)` is therefore two
  independent tests; the two values need not be the same.

#### Filter Matched Values: `filter()`

`filter(p)` keeps the entries of a collection whose value satisfies `p`. A later
predicate on the result plays the role of `MATCHED_VARS`:

```
# 932200: SecLang links  ARGS @rx A  ->  MATCHED_VARS @rx /  ->  MATCHED_VARS @rx \s
rule 932200 (severity: critical) {
  let candidates = request.args |> lowercase() |> url_decode_uni() |> filter(matches(A))
  when candidates |> matches("/")
   and candidates |> matches("\\s")
  then block
}
```

This is exact. ModSecurity v2 and Coraza do not clear `MATCHED_VARS` between links of a
chain; each link adds its matches. The values a later link adds come from the same
filtered set, so every `MATCHED_VARS` link tests that set, and the two tests above may
match different values, just as in SecLang. The compiler lowers a filtered collection
used by later predicates to a first link on the original targets and later links on
`MATCHED_VARS`.

#### Extract, Then Test: `extract()`

`extract(pattern, group)` returns a capture group of a value (`group` defaults to 0,
the whole match). A value that does not match yields nothing, and any predicate on
nothing is false, so `extract(...) |> not(...)` is false when the pattern does not
match, as the SecLang chain is.

```
# 920420: REQUEST_HEADERS:Content-Type @rx ^[^;\s]+  (capture, setvar tx.content_type)
#         -> TX:content_type !@within %{tx.allowed_request_content_type}
# (simplified: CRS also wraps the value in |...| for @within)
rule 920420 (severity: critical) {
  when request.headers["Content-Type"]
       |> extract("^[^;\\s]+")
       |> lowercase()
       |> not(within(tx.allowed_request_content_type))
  then block
}

# 943110: the Referer host must not end with the Host header
rule 943110 (severity: critical) {
  let referer_host = request.headers["Referer"] |> extract("^(?:ht|f)tps?://(.*?)/", 1)
  when request.args.keys |> lowercase() |> matches(SESSION_PARAM_NAMES)
   and referer_host |> not(ends_with(request.headers["Host"]))
  then block
}
```

`extract()` replaces the `capture()` predicate: capture groups are values, not
transaction state, so a binding that uses them can be referenced any number of times.
The compiler emits `capture` plus `TX:n` (or a `setvar` to a TX variable) when lowering
to SecLang.

On a collection, `extract()` applies to each value, so the language checks every value.
SecLang's `TX:n` holds only the groups of the **last** value that matched. 942130 and
942131 compare two groups of the same value (`TX:1 @streq %{TX.2}`). In CRSLang this
means "some argument has equal groups 1 and 2", while the compiled SecLang chain, like
the CRS rule it came from, sees the last matching argument only. That is the one place
where SecLang output detects less than the CRSLang rule means. Round-tripping CRS through
CRSLang still produces the same SecLang (ADR-0010).

#### Keep a Value for the Log: `matched.var`

In `then`, `msg` and `log(data: ...)`, `matched.var` and `matched.var_name` refer to the
value matched by the **first predicate** of the condition, in source order. That is the
value CRS keeps in `tx.<id>_matched_var_name` today. The compiler generates that
`setvar` when a rule has more than one predicate, so authors do not write it.

#### Chains Outside These Patterns

The SecLang importer recognizes the three patterns. A chain with any other action on a
link before the last one (for example `ctl`, as in 905111) is not representable: the
importer reports it and either splits it into separate rules that pass state through a
TX variable, or keeps it as raw SecLang. No CRS or plugin rule needs this.

#### Open Questions

- **Testing two parts of one value.** 942130/942131 compare two capture groups of the
  same argument. That needs a way to refer to the value inside a per-value test, such as
  `request.args |> any(v -> v |> extract(R, 1) |> eq(v |> extract(R, 2)))` or a
  dedicated `groups_equal(R, 1, 2)`. The syntax is not decided.
- **Choosing the logged value.** `matched.var` follows the first predicate. Whether a
  rule can name another predicate's value (for example through a `let` binding) is not
  decided.

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

For chains with actions on earlier links, by pattern (see Conditions Are Pure):

**Filter matched values:** a later link on `MATCHED_VARS` becomes a predicate on a
`filter()`ed collection:

```
# SecLang:
SecRule ARGS "@rx A" "capture,chain"
    SecRule MATCHED_VARS "@rx /"
# CRSLang:
when request.args |> filter(matches(A)) |> matches("/")
```

**Extract, then test:** `capture` plus a later `TX:n` (or a `setvar` copy of it)
becomes `extract()`:

```
# SecLang:
SecRule REQUEST_HEADERS:Content-Type "@rx ^[^;\s]+" "capture,setvar:'tx.content_type=%{tx.0}',chain"
    SecRule TX:content_type "!@within %{tx.allowed_request_content_type}"
# CRSLang:
when request.headers["Content-Type"] |> extract("^[^;\\s]+")
     |> not(within(tx.allowed_request_content_type))
```

**Keep a value for the log:** `setvar tx.<id>_matched_var_name` is dropped; `matched.var`
in `then` refers to the first predicate's value.

### Multi-Stage Matching: `multi_match`

SecLang's `multiMatch` action changes how the operator sees a target that has
transformations. Without it, the operator runs once, on the fully transformed value.
With it, the operator runs on the original value and again after each transformation
that changes the value, and every stage that matches counts as a match. This is
currently modeled as a non-disruptive action, but it changes how the condition
evaluates, not what happens after it matches.

**Recommendation: `multi_match` as rule metadata.**

```
# Without multi_match: the operator tests only the final transformed value
rule 942100 {
  when request.args |> utf8_to_unicode() |> url_decode_uni() |> detect_sqli()
  then block
}

# With multi_match: the operator tests the original value and the value after
# each transformation
rule 942100 (multi_match: true) {
  when request.args |> utf8_to_unicode() |> url_decode_uni() |> detect_sqli()
  then block
}
```

The flag applies to every predicate in the rule. In SecLang, `multiMatch` is set per
chain link, so a rule-level flag is lossy only for a chain where some links use it and
others transform their input without it. CRS has no such chain: in the one chained rule
that uses `multiMatch` (942131), the other link has no transformations, and multi-stage
matching has no effect without them. The CRS ruleset therefore still compiles to SecLang
losslessly (ADR-0010).

**Alternatives documented:**

- **Option B: Predicate-level wrapper** — `multi_match(request.args |> utf8_to_unicode()
  |> detect_sqli())`. Exact per-predicate control and lossless for any chain, but it adds
  a construct to the pipeline syntax for a rarely used feature.
- **Option C: Drop it** — always test only the final value. Changes the behavior of
  the eight CRS rules that use `multiMatch`.

`multiMatch` is rarely used in CRS, so the metadata flag is sufficient for the
foreseeable future. Option B can be revisited if a rule needs different settings per
predicate.

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
rule 100001 {
  when request.uri |> matches("/admin")
   and (client.ip |> ip_in_range("10.0.0.0/8")
        or request.headers["X-Internal"] |> eq("true"))
  then pass
}

# Complex grouping (phase inferred as request_body from request.body, ADR-0007)
rule 100002 {
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
  phase: request_headers
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
- Conditions are pure, so evaluation order is not part of a rule's meaning
- Closer to how security engineers think about conditions

### Negative

- `filter()`, `extract()` and `let` add concepts SecLang does not have; authors have to
  learn that a predicate on a collection means "at least one value"
- For 942130/942131, SecLang output checks fewer values than the CRSLang rule means
  (see Extract, Then Test)
- SecLang chains with other actions on earlier links (such as `ctl`) cannot be imported
  as single rules
- Some deeply chained rules may become long single expressions (mitigated by
  line breaks and formatting conventions)

### Risks

- **Pattern coverage** — the importer's three patterns cover CRS v4.30.0 and the audited
  plugins. User rules may use other actions on earlier links; those are reported and
  split or kept as raw SecLang, not supported by the language.
- **Migration correctness** — automated chain-to-boolean conversion must be validated
  against the full CRS test suite. The three patterns need separate migration paths in the
  converter tool.
