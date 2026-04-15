# ADR-0011: First-Class Scoring Model

- **Status:** Proposed
- **Date:** 2026-04-15
- **Phase:** 3 (part of structured actions redesign)

## Context

CRS uses an anomaly scoring model where rules accumulate scores rather than blocking
immediately. Each attack detection rule adds to a total anomaly score and a
category-specific score. A separate evaluation rule in a later phase checks whether the
accumulated score exceeds a threshold and takes action.

Today, this is implemented as manual TX variable manipulation in every rule:

```
SecRule REQUEST_ARGS "@detectSQLi" \
    "id:942100,\
    phase:2,\
    block,\
    capture,\
    t:none,t:urlDecodeUni,t:lowercase,\
    msg:'SQL Injection Attack Detected',\
    severity:'CRITICAL',\
    setvar:'tx.sql_injection_score=+%{tx.critical_anomaly_score}',\
    setvar:'tx.inbound_anomaly_score_pl1=+%{tx.critical_anomaly_score}'"
```

The score values are indirectly derived from severity — `tx.critical_anomaly_score` is
set to `5` in `crs-setup.conf`, `tx.warning_anomaly_score` to `3`, etc. But this
indirection is manual wiring: every rule must include the correct `setvar` actions with
the correct TX variable reference for its severity level.

### Problems

1. **Boilerplate** — every attack detection rule has 2-3 `setvar` actions for scoring.
   This is the single largest source of repetition in CRS rules.
2. **Error-prone** — using the wrong score variable (`tx.warning_anomaly_score` instead
   of `tx.critical_anomaly_score`) for a CRITICAL severity rule is a silent bug.
3. **Indirection** — the severity-to-score mapping is defined in `crs-setup.conf` as TX
   variables, then referenced via `%{tx.critical_anomaly_score}` interpolation in every
   rule. The relationship between severity and score is not explicit.
4. **Category scores are manual** — whether a rule increments `tx.sql_injection_score`
   or `tx.xss_score` depends on the rule author remembering to add the right `setvar`.
5. **Not toolable** — scoring analysis (which rules contribute how much to the total
   score, what thresholds are configured) requires parsing `setvar` strings.

## Decision

Make anomaly scoring a **first-class language concept**. Scoring is derived from severity
and category, with explicit override capability.

### Severity-Derived Scoring

When a rule declares a severity, scoring happens automatically:

```
rule 942100 (severity: critical) {
  when request.args |> detect_sqli()
  then block
}
```

The compiler derives the score from the severity using a scoring table defined in
globals:

```
globals {
  scoring {
    critical = 5
    warning  = 3
    notice   = 2
  }
}
```

This rule automatically contributes `5` to the anomaly score because it has
`severity: critical` and the scoring table maps critical to 5.

### Category Scoring via Groups

Category-specific scores (SQL injection, XSS, RCE, etc.) are derived from the group
a rule belongs to:

```
group "sql_injection" {
  category = "sqli"

  rule 942100 (severity: critical) {
    when request.args |> detect_sqli()
    then block
  }

  rule 942110 (severity: critical) {
    when request.args |> detect_sqli()
    then block
  }
}
```

Rules in this group automatically increment both `anomaly_score` and `sqli_score` by the
severity-derived value. No `setvar` needed.

### Explicit Score Override

For rules that need non-standard scoring, an explicit `score()` function is available:

```
rule 920170 (severity: warning) {
  when request.method |> matches("^(?:GET|HEAD)$")
   and request.headers["Content-Length"] |> not(matches("^0?$"))
  then block {
    score(anomaly: 5, protocol: 3)  # override severity-derived scores
  }
}
```

### Paranoia Levels

CRS paranoia levels (PL1-PL4) control which rules are active and which score bucket
they contribute to (`tx.inbound_anomaly_score_pl1` through `pl4`). In the new model,
paranoia level is a rule attribute:

```
rule 942100 (severity: critical, paranoia: 1) {
  when request.args |> detect_sqli()
  then block
}

rule 942101 (severity: critical, paranoia: 2) {
  when request.args |> url_decode() |> detect_sqli()
  then block
}
```

The compiler uses the paranoia attribute to:
1. Determine which PL score bucket to increment (for SecLang output)
2. Determine whether the rule is active at a given paranoia level
3. Generate the appropriate skip/marker logic for SecLang

### Scoring Evaluation

The scoring evaluation rule (currently rule 949110/959100 in CRS) becomes a built-in
concept rather than a hand-written rule:

```
# Defined in globals or as a special directive
scoring_threshold {
  inbound  = 5    # was: tx.inbound_anomaly_score_threshold
  outbound = 4    # was: tx.outbound_anomaly_score_threshold
  action   = deny(status: 403)
}
```

This replaces the complex phase-5 evaluation rules in CRS that manually compare
accumulated scores against thresholds.

### Compilation to SecLang

The compiler expands first-class scoring to SecLang `setvar` actions:

```
# CRSLang
rule 942100 (severity: critical, paranoia: 1) {
  when request.args |> detect_sqli()
  then block
}

# Compiled SecLang
SecRule REQUEST_ARGS "@detectSQLi" \
    "id:942100,\
    phase:2,\
    block,\
    t:none,\
    severity:'CRITICAL',\
    setvar:'tx.sql_injection_score=+%{tx.critical_anomaly_score}',\
    setvar:'tx.inbound_anomaly_score_pl1=+%{tx.critical_anomaly_score}'"
```

The `scoring {}` globals compile to TX variable initialization in `crs-setup.conf`.
The `scoring_threshold {}` compiles to the evaluation rules in phases 3-5.

### Compilation to Other Targets

For targets that don't have native anomaly scoring (Cloud Armor, AWS WAF, Cloudflare):
- Each rule becomes a standalone rule in the target format
- The scoring model is either:
  - Emulated via target-specific mechanisms (if available)
  - Simplified to direct block/allow (each rule acts independently)
  - Documented as a limitation of the target

### IR Representation

```go
type ScoringConfig struct {
    Levels    map[Severity]int     // critical: 5, warning: 3, notice: 2
    Threshold ScoringThreshold
}

type ScoringThreshold struct {
    Inbound  int
    Outbound int
    Action   DisruptiveAction
}

type RuleMetadata struct {
    ID       int
    Severity Severity        // critical, warning, notice
    Paranoia int             // 1-4
    Category string          // "sqli", "xss", "rce", etc. (from group)
    // ... other metadata
}

// Explicit score override in effects
type ScoreEffect struct {
    Scores map[string]int   // anomaly: 5, protocol: 3, etc.
}
```

## Alternatives Considered

### A: Keep scoring as TX variable manipulation

Leave `setvar` in effect blocks, just make the syntax nicer:

```
then block {
  tx.anomaly_score += 5
  tx.sqli_score += 5
}
```

**Rejected because:**
- Still boilerplate in every rule
- Still error-prone (wrong score value for severity)
- Still not toolable for score analysis
- ADR-0004 already has this syntax as the baseline; this ADR builds on it

### B: Score as a named effect only (no severity derivation)

```
then block {
  score(anomaly: 5, sqli: 5)
}
```

**Rejected as the only mechanism because:**
- Still requires every rule to declare its scores
- The severity-to-score mapping is the entire point — it's a policy, not a per-rule
  decision
- Explicit `score()` is kept as an override mechanism

### C: Score as rule metadata only (no explicit override)

```
rule 942100 (severity: critical, category: sqli) {
  when ...
  then block  # scoring is 100% automatic
}
```

**Rejected because:**
- Some rules legitimately need non-standard scoring
- No escape hatch for edge cases
- Too rigid for the full CRS ruleset

## Consequences

### Positive

- Eliminates 2-3 `setvar` actions from every attack detection rule (~80% of CRS rules)
- Severity-to-score mapping is a single configuration, not scattered across hundreds of
  rules
- Category scoring is automatic from group membership
- Paranoia levels become a typed attribute, not a naming convention
- Scoring analysis becomes trivial (query the IR for severity/category/paranoia)
- Score threshold configuration is a single declaration, not a complex evaluation rule
- Compile-time validation: severity is required on attack rules, preventing missing scores

### Negative

- Rules that deviate from the standard scoring model need the `score()` override, which
  mixes two scoring mechanisms (implicit + explicit)
- Paranoia levels as a language concept may not map to all target engines
- The scoring evaluation rule (949110/959100) becomes compiler-generated, which means
  debugging it requires understanding the compiler output

### Risks

- **Scoring model flexibility** — CRS's scoring model may evolve (new score categories,
  different threshold logic). The language must be flexible enough to accommodate changes
  without breaking existing rules.
- **Target compatibility** — anomaly scoring is fundamentally a ModSecurity/Coraza
  concept. Other targets may handle scoring differently or not at all. The multi-target
  compilation model (ADR-0010) must define how scoring degrades per target.
- **Migration** — converting existing CRS rules from manual `setvar` to severity-derived
  scoring requires validating that the derived scores match the original. An automated
  migration tool should compare outputs.
