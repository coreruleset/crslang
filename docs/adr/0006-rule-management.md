# ADR-0006: Rule Management Directives

- **Status:** Proposed
- **Date:** 2026-04-13
- **Phase:** 5

## Context

SecLang provides several directives for modifying rules after they are defined. These are
essential for CRS deployments where users need to customize rules without editing the
core ruleset:

- **`SecRuleRemoveById`** — disable a rule entirely by ID
- **`SecRuleRemoveByTag`** — disable all rules matching a tag
- **`SecRuleRemoveByMsg`** — disable all rules matching a message string
- **`SecRuleUpdateTargetById`** — add or remove targets from a rule
- **`SecRuleUpdateTargetByTag`** — add or remove targets from rules by tag
- **`SecRuleUpdateTargetByMsg`** — add or remove targets from rules by message
- **`SecRuleUpdateActionById`** — replace actions on a rule
- **`SecMarker`** — named label for `skipAfter` control flow

CRSLang v1 models these directly:

```yaml
- kind: remove
  metadata:
    comment: "Disable SQL injection check for admin API"
  selector:
    by_id: 942100

- kind: update_target
  selector:
    by_id: 920170
  collections:
    - name: ARGS
      arguments: [username]
      excluded: true
```

### Problems

1. **Selector inconsistency** — three selector types (by_id, by_tag, by_msg) with
   identical behavior but different fields.
2. **Target updates are complex** — adding and removing targets uses the same
   `collections` structure with an `excluded` flag, mixing positive and negative
   modifications.
3. **Action updates are opaque** — `SecRuleUpdateActionById` replaces the entire action
   list, which is error-prone and hard to review.
4. **No composition** — you cannot say "for all rules in phase 1 with tag X, remove
   target Y." Each operation is a standalone directive.
5. **SecMarker** — a control-flow label that exists only because `skipAfter` needs a
   target. If `skip`/`skipAfter` are deprecated (ADR-0004), markers may become
   unnecessary.

## Decision

Introduce first-class **rule management directives** with a unified selector system and
explicit modification operations.

### Exclude (Remove) Rules

```
# By ID
exclude rule 942100

# By tag
exclude rules where tag == "OWASP_CRS/SQL_INJECTION"

# By tag pattern
exclude rules where tag |> matches("^OWASP_CRS/SQL.*")

# By severity
exclude rules where severity == critical

# Combined
exclude rules where tag == "OWASP_CRS/SQL_INJECTION" and phase == request
```

### Update Targets

```
# Remove a specific target from a rule
update rule 920170 {
  remove target request.args["username"]
}

# Add a target
update rule 920170 {
  add target request.headers["X-Custom"]
}

# By tag — apply to all matching rules
update rules where tag == "OWASP_CRS/SQL_INJECTION" {
  remove target request.args["search_query"]
}
```

### Update Actions

Rather than replacing the entire action block, provide surgical modifications:

```
# Change disruptive action
update rule 920170 {
  set action pass
}

# Change severity
update rule 920170 {
  set severity warning
}

# Add an effect
update rule 920170 {
  add effect tx.custom_score += 10
}

# Remove an effect
update rule 920170 {
  remove effect capture()
}
```

### Rule Groups

Named groups replace `SecMarker` and provide a scope for batch operations:

```
group sql_injection_checks {
  rule 942100 (...) { ... }
  rule 942110 (...) { ... }
  rule 942120 (...) { ... }
}

# Exclude the entire group
exclude group sql_injection_checks

# Update all rules in a group
update group sql_injection_checks {
  remove target request.args["allowed_field"]
}
```

### Unified Selector System

All management directives use the same selector grammar:

```
selector       = "rule" INTEGER                              # single rule by ID
               | "rules" "where" selector_expr               # multiple rules by query
               | "group" IDENT                               # named group

selector_expr  = selector_expr "and" selector_expr
               | selector_expr "or" selector_expr
               | selector_field "==" value
               | selector_field "|>" func_call
               | "(" selector_expr ")"

selector_field = "tag" | "phase" | "severity" | "id" | "message"
```

### Ordering and Scope

Rule management directives are processed in file order, after all rules are loaded.
This matches SecLang's behavior where removal/update directives in a later file affect
rules from earlier files.

For CRS deployments, the convention is:

```
# Core rules (provided by CRS)
rules/
  REQUEST-901-INITIALIZATION.crs
  REQUEST-941-APPLICATION-ATTACK-XSS.crs
  REQUEST-942-APPLICATION-ATTACK-SQLI.crs
  ...

# User customizations (site-specific)
customizations/
  exclude-false-positives.crs      # exclude/update directives
  site-specific-rules.crs          # additional rules
```

### IR Representation

```go
type ExcludeDirective struct {
    Selector RuleSelector
    Comment  string
}

type UpdateDirective struct {
    Selector     RuleSelector
    Modifications []Modification
    Comment      string
}

type RuleSelector interface {
    selectorNode()
}

type ByIDSelector struct {
    ID int
}

type WhereSelector struct {
    Expr SelectorExpr  // boolean expression over rule metadata
}

type GroupSelector struct {
    Name string
}

type Modification interface {
    modNode()
}

type AddTarget struct {
    Target FieldRef
}

type RemoveTarget struct {
    Target FieldRef
}

type SetAction struct {
    Action DisruptiveAction
}

type SetMetadata struct {
    Field string
    Value Value
}

type AddEffect struct {
    Effect Effect
}

type RemoveEffect struct {
    Effect Effect
}
```

### YAML v2 Representation

```yaml
- kind: exclude
  comment: "Disable SQL injection check for admin API"
  selector:
    rule: 942100

- kind: exclude
  selector:
    where:
      tag: "OWASP_CRS/SQL_INJECTION"
      phase: request

- kind: update
  selector:
    rule: 920170
  modifications:
    - remove_target: request.args["username"]
    - add_target: request.headers["X-Custom"]

- kind: update
  selector:
    where:
      tag: "OWASP_CRS/SQL_INJECTION"
  modifications:
    - set_action: pass
```

### Migration from SecLang

| SecLang | CRSLang |
|---------|---------|
| `SecRuleRemoveById 942100` | `exclude rule 942100` |
| `SecRuleRemoveByTag "SQL_INJECTION"` | `exclude rules where tag == "SQL_INJECTION"` |
| `SecRuleUpdateTargetById 920170 "!ARGS:username"` | `update rule 920170 { remove target request.args["username"] }` |
| `SecRuleUpdateActionById 920170 "pass"` | `update rule 920170 { set action pass }` |
| `SecMarker END_SQL_CHECKS` | `group sql_checks { ... }` (or label if skip is needed) |

## Alternatives Considered

### A: Inheritance / Override Model

```
rule 920170 extends base:920170 {
  exclude target request.args["username"]
  override action pass
}
```

**Rejected because:**
- Implies a class hierarchy that does not exist
- Conflates "I want to modify this rule" with "I want to define a new rule based on it"
- The `extends` keyword suggests the original rule still runs, which is confusing

### B: Annotation-Based

```
@exclude(942100)
@remove_target(920170, request.args["username"])
```

**Rejected because:**
- Annotations are typically metadata, not imperative operations
- Does not support the `where` selector pattern
- Visually disconnected from the rules they affect

### C: Separate Configuration File Format

Use TOML, JSON, or a custom format for exclusions, keeping rule management separate
from rule definitions.

**Rejected because:**
- Forces users to learn a second file format
- The management directives are part of the rule language — they reference fields,
  targets, and actions using the same syntax
- Separate formats prevent tooling from validating references

## Consequences

### Positive

- Unified selector system replaces three separate selector types
- Surgical updates instead of full action replacement
- Rule groups provide natural scoping for batch operations
- Same expression syntax in selectors and conditions (reuse of `|>`, `matches()`, etc.)
- Clear, readable exclusion files for CRS deployments

### Negative

- The `where` selector adds query-like complexity to what is currently a simple
  ID/tag/msg lookup
- Rule groups are a new organizational concept that must be adopted by CRS maintainers
- Processing order must be well-defined to avoid surprises

### Risks

- **Selector performance** — `where` queries over large rulesets need efficient
  evaluation. For CRS-sized rulesets (~300 rules) this is not a concern; for very
  large custom rulesets, consider indexing by tag and phase.
- **Circular modifications** — an update that changes a tag could affect which rules
  another `where` selector matches. Define processing as single-pass: selectors
  evaluate against the original rule state, not the modified state.
- **Group adoption** — CRS rules are currently organized by file, not by explicit
  groups. Groups are optional and additive — existing file-based organization continues
  to work.
