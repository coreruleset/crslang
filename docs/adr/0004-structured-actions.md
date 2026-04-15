# ADR-0004: Structured Action Model

- **Status:** Proposed
- **Date:** 2026-04-13
- **Phase:** 3

## Context

CRSLang v1 inherits SecLang's action model, which groups actions into four categories
in a flat bag:

```yaml
actions:
  disruptive:
    action: deny
  non-disruptive:
    - action: log
    - action: auditlog
    - action: setvar
      param: "tx.anomaly_score=+5"
    - action: setvar
      param: "tx.inbound_anomaly_score=+5"
    - action: capture
  flow:
    - chain
  data:
    - action: status
      param: "403"
```

### Problems

1. **Category confusion** — the four categories (disruptive, non-disruptive, flow, data)
   are SecLang implementation details. Rule authors think in terms of "what happens when
   this matches," not which category an action belongs to.
2. **setvar string parsing** — `setvar: "tx.anomaly_score=+5"` encodes a field, an
   operation, and a value in a string that must be parsed separately.
3. **flow actions in the action bag** — `chain` is a structural concept (ADR-0003
   replaces it), not an action. `skip`/`skipAfter` are control flow, not side-effects.
4. **ctl overloading** — `ctl:forceRequestBodyVariable=On` mixes engine configuration
   with rule actions.
5. **status as a data action** — the HTTP response status is a parameter of the
   disruptive action, not a separate action.

## Decision

Replace the action bag with a **structured action block** consisting of a disruptive
action (with parameters) and an optional block of side-effects.

### Syntax

```
then <disruptive-action>(<params>) {
  <side-effects>
}
```

### Disruptive Actions

Each disruptive action can take named parameters:

```
then deny(status: 403)
then block(status: 403)
then redirect(url: "https://example.com/blocked")
then drop
then pass
then allow
```

`status` as a parameter of the disruptive action, not a separate data action.

### Side-Effect Statements

The block body contains zero or more side-effect statements:

#### Variable Assignment

```
tx.anomaly_score += 5          # was: setvar:'tx.anomaly_score=+5'
tx.sql_injection_score += 5    # was: setvar:'tx.sql_injection_score=+5'
tx.blocking_early = 1          # was: setvar:'tx.blocking_early=1'
ip.reput_block_flag = 1        # was: setvar:'ip.reput_block_flag=1'
```

Assignment operators:
- `=`  — set (was: `setvar:'collection.var=value'`)
- `+=` — increment (was: `setvar:'collection.var=+value'`)
- `-=` — decrement (was: `setvar:'collection.var=-value'`)

Variable expiry:

```
tx.block_duration = 3600 (expires: 3600)   # was: expirevar:'tx.block_duration=3600'
```

#### Logging

```
log()                    # was: action: log
audit_log()              # was: action: auditlog
log(data: "%{MATCHED_VAR}")  # was: logdata:'%{MATCHED_VAR}'
```

Or combined (since most rules use both):

```
log(audit: true)         # log + auditlog
log(audit: true, data: matched.var)
```

#### Capture

```
capture()                # was: action: capture
```

#### Engine Configuration

```
configure {
  force_request_body = true    # was: ctl:forceRequestBodyVariable=On
  request_body_processor = XML # was: ctl:requestBodyProcessor=XML
  rule_engine = detect_only    # was: ctl:ruleEngine=DetectOnly
  audit_engine = off           # was: ctl:auditEngine=Off
}
```

Or inline for single settings:

```
configure(request_body_processor: XML)
```

### Complete Examples

**Simple block with scoring:**
```
rule 941100 (phase: request, severity: critical) {
  when request.args |> js_decode() |> detect_xss()
  then block {
    tx.xss_score += 5
    tx.anomaly_score += 5
    capture()
    log(audit: true, data: matched.var)
  }
}
```

**Pass with configuration:**
```
rule 900100 (phase: request) {
  when true
  then pass {
    tx.paranoia_level = 1
    tx.anomaly_score_threshold = 5
  }
}
```

**Deny with status:**
```
rule 901001 (phase: request, severity: critical) {
  when count(tx.crs_setup_version) |> eq(0)
  then deny(status: 500) {
    log(audit: true)
  }
}
```

### IR Representation

```go
type ActionBlock struct {
    Action     DisruptiveAction
    Parameters map[string]Value    // status, url, etc.
    Effects    []Effect
}

type Effect interface {
    effectNode()
}

type VarAssignment struct {
    Target   FieldRef       // tx.anomaly_score
    Operator AssignOp       // Set, Add, Subtract
    Value    Value
    Expiry   *int           // optional TTL in seconds
}

type LogEffect struct {
    Audit bool
    Data  *FieldRef         // optional logdata field
}

type CaptureEffect struct{}

type ConfigureEffect struct {
    Settings map[string]Value
}
```

### YAML v2 Representation

```yaml
then:
  action: block
  params:
    status: 403
  effects:
    - assign:
        target: tx.anomaly_score
        op: "+="
        value: 5
    - log:
        audit: true
        data: matched.var
    - capture: true
```

### Flow Actions

- **`chain`** — eliminated entirely (ADR-0003)
- **`skip`/`skipAfter`/`SecMarker`** — eliminated entirely. These were an
  implementation detail of ModSecurity's sequential evaluation model. The actual intent
  is **conditional rule activation** — "only run these rules if condition X holds."

  This is now expressed as **guarded groups** (ADR-0006):

  ```
  # was: SecRule TX:DETECTION_PARANOIA_LEVEL "@lt 2" "skipAfter:END-941"
  #      ... rules ...
  #      SecMarker "END-941"

  group "xss_pl2" (requires: paranoia >= 2) {
    rule 941120 (severity: critical) { ... }
    rule 941130 (severity: critical) { ... }
  }
  ```

  The compiler generates the appropriate `skipAfter`/`SecMarker` pairs when compiling
  to SecLang. For paranoia-level gating specifically, the `paranoia` attribute on rules
  (ADR-0011) allows the compiler to group and gate rules automatically.

  No `skip_to()`, `goto`, or `label` exists in CRSLang. The language expresses intent
  (conditional activation), not mechanism (skip/marker).

## Alternatives Considered

### A: Keep the four-category bag, just improve syntax

```yaml
actions:
  disruptive: deny
  status: 403
  setvar:
    - tx.anomaly_score += 5
  logging: [log, auditlog]
```

**Rejected because:**
- Still carries the category model that confuses authors
- Does not solve the fundamental problem of actions being a grab-bag

### B: All actions as function calls

```
then {
  deny(status: 403)
  set(tx.anomaly_score, "+", 5)
  log(audit: true)
}
```

**Rejected because:**
- The disruptive action is not a side-effect; it is the *outcome* of the rule.
  It deserves syntactic prominence, not burial in a block.
- `deny()` as a function call implies it could appear conditionally inside the block,
  which should not be possible.

### C: Separate `action` and `effects` keywords

```
rule 920170 {
  when ...
  action block(status: 403)
  effects {
    tx.anomaly_score += 5
    log()
  }
}
```

**Considered viable** but adds a keyword without clear benefit over `then action { effects }`.

## Consequences

### Positive

- Clear separation: the disruptive action is the rule's outcome; effects are
  side-effects that happen when it matches
- `setvar` string parsing eliminated; assignments use native syntax
- `ctl:` directives get their own `configure` block instead of being overloaded
  actions
- Status is a parameter, not a separate action type
- Familiar imperative syntax for the effect block

### Negative

- Two syntactic constructs for effects: assignment syntax (`tx.score += 5`) and
  function syntax (`log()`, `capture()`). This is intentional — assignments and
  function calls are conceptually different.
- `configure {}` may need engine-specific settings that vary across backends

### Risks

- **Engine mapping** — some effects may not have equivalents in all target engines.
  Define a core set that all backends must support and an extension mechanism for
  engine-specific effects.
- **Ordering semantics** — do effects execute in order? In parallel? Define this
  explicitly (recommendation: ordered, top-to-bottom, matching SecLang behavior).
