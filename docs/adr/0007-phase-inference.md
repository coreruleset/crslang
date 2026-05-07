# ADR-0007: Phase Inference from Field Types

- **Status:** Proposed
- **Date:** 2026-04-13
- **Phase:** 1 (builds on ADR-0001)

## Context

SecLang requires every rule to declare a numeric phase (1-5) that controls when the rule
evaluates during request processing:

| Phase | Name | Available data |
|-------|------|----------------|
| 1 | Request Headers | Method, URI, headers, cookies, query string args |
| 2 | Request Body | POST body, file uploads, multipart data |
| 3 | Response Headers | Response status, response headers |
| 4 | Response Body | Response body content |
| 5 | Logging | Post-processing (rarely used in rules) |

In CRSLang v1 phase is a string in metadata:

```yaml
metadata:
  phase: "1"
```

This is redundant information in the vast majority of rules. If a rule inspects
`REQUEST_HEADERS:Content-Type`, it must be phase 1. If it inspects `RESPONSE_BODY`, it
must be phase 4. The phase is already implied by the fields the rule references.

### How Phase Maps to Fields

With the typed field namespace from ADR-0001:

| Field prefix | Implied phase | SecLang phase |
|---|---|---|
| `request.method`, `request.uri`, `request.protocol`, `request.line` | request_headers | 1 |
| `request.headers[*]`, `request.cookies[*]` | request_headers | 1 |
| `request.args.get[*]`, `request.args.get.names` | request_headers | 1 |
| `request.body`, `request.args.post[*]`, `request.args[*]` | request_body | 2 |
| `multipart.*`, `files.*` | request_body | 2 |
| `response.status`, `response.protocol` | response_headers | 3 |
| `response.headers[*]` | response_headers | 3 |
| `response.body` | response_body | 4 |
| `tx.*`, `ip.*`, `global.*`, `session.*` | any (cross-phase) | context-dependent |
| `matched.*`, `rule.*` | any (cross-phase) | context-dependent |
| `client.ip`, `server.*` | any (cross-phase) | context-dependent |
| `time.*`, `env.*` | any (cross-phase) | context-dependent |

### Phase-Ambiguous Fields

Some fields are available in all phases. When a rule uses **only** cross-phase fields
(e.g., only `tx.*` variables), the phase cannot be inferred and must be declared
explicitly.

In CRS today, this pattern appears in:
- Initialization rules (phase 1, TX-only) — e.g., rule 901001 checking
  `count(tx.crs_setup_version)`
- Paranoia level skip rules (phase 1-4, TX-only) — e.g., rules 911011/911012 checking
  `tx.detection_paranoia_level`
- Scoring evaluation rules (phase 5, TX-only) — evaluating accumulated anomaly scores

Most of these specific patterns are eliminated by other proposals: initialization rules
become compiler-generated (ADR-0012), paranoia skip rules are replaced by guarded groups
(ADR-0006), and scoring evaluation is driven by `config { scoring {} }` (ADR-0012). The
inference rules below still matter for the remaining cases — custom rules that read
or write TX state across rules, and any cross-phase logic the user authors directly.

## Decision

**Phase is inferred from field references when unambiguous, and explicitly declared only
when needed.**

### Inference Rules

1. If a rule references any **request-headers-phase** field → phase is `request_headers`
2. If a rule references any **request-body-phase** field → phase is `request_body`
3. If a rule references any **response-headers-phase** field → phase is `response_headers`
4. If a rule references any **response-body-phase** field → phase is `response_body`
5. If a rule references only **cross-phase** fields → phase must be declared explicitly

### Conflict Detection

If a rule references fields from multiple phases, that is a **compile-time error**:

```
# ERROR: request.method is phase 1, response.body is phase 4
rule 999999 {
  when request.method |> eq("GET")
   and response.body |> contains("error")
  then block
}
```

This is also an error in SecLang (you cannot inspect response body in phase 1), but
SecLang does not detect it — it silently produces wrong results. CRSLang catches it
at compile time.

@theseion: I don't understand. In SecLang, this rule would be evaluated in phase 4. Why would that be an error? And from the logic above, the example rule would
also run in phase 4. The way this section is phrased, a rule would never be able to inspect headers and body at the same time because they belong to different phases.

### Mixed Phase-Specific and Cross-Phase Fields

When a rule mixes phase-specific and cross-phase fields, the phase-specific field wins:

```
# Phase is inferred as request_headers (from request.method)
# tx.anomaly_score is cross-phase, so it does not affect inference
rule 920170 {
  when request.method |> matches("^(?:GET|HEAD)$")
   and tx.some_flag |> eq(1)
  then block
}
```

### Explicit Phase Declaration

When inference is not possible (TX-only rules) or when the author wants to override for
clarity, phase is declared in metadata. The example below uses rule 901001 as a familiar
illustration of the inference behavior; in practice 901001 itself is compiler-generated
per ADR-0012, and explicit phase declaration applies to user-authored TX-only rules:

```
# Must declare phase — only TX fields used
rule 901001 (phase: request_headers) {
  when count(tx.crs_setup_version) |> eq(0)
  then deny(status: 500)
}

# Paranoia gating — replaced by guarded groups (ADR-0006)
# Instead of skip rules duplicated across phases:
group "method_enforcement_pl1" (requires: paranoia >= 1) {
  rule 911100 (phase: request_headers) { ... }
  rule 911101 (phase: request_body) { ... }
}
```

### Phase Names

CRSLang uses descriptive phase names instead of numbers:

| CRSLang name | SecLang number |
|---|---|
| `request_headers` | 1 |
| `request_body` | 2 |
| `response_headers` | 3 |
| `response_body` | 4 |
| `logging` | 5 |

@theseion nginx has more phases: https://nginx.org/en/docs/dev/development_guide.html#http_phases and https://nginx.org/en/docs/stream/stream_processing.html.
We might want to document the corresponding phases in nginx here, maybe even consider what this list could look like in the future.

The text syntax also accepts short forms for ergonomics:

| Short form | Full name |
|---|---|
| `request` | `request_headers` (most common request phase) |
| `response` | `response_headers` (most common response phase) |

@theseion: I feel like these shor forms will create more confusion than be helpful...

### Compilation to SecLang

When exporting to SecLang, the compiler maps inferred or declared phases back to numeric
values. This is lossless:

```
# CRSLang (phase inferred from request.headers)
rule 920170 {
  when request.headers["Content-Type"] |> matches("^application/json")
  then block
}

# Compiled SecLang
SecRule REQUEST_HEADERS:Content-Type "@rx ^application/json" \
    "id:920170,phase:1,block"
```

For rules with explicit phase:

```
# CRSLang
rule 901001 (phase: request_headers) {
  when count(tx.crs_setup_version) |> eq(0)
  then deny(status: 500)
}

# Compiled SecLang
SecRule &TX:crs_setup_version "@eq 0" \
    "id:901001,phase:1,deny,status:500"
```

### Sub-Phase Ambiguity: `request.args`

One field requires special handling: `request.args` (all arguments, GET + POST combined).
In SecLang, `ARGS` includes both query string and POST body parameters. It is technically
available in phase 1, but POST parameters are only populated after phase 2 processes the
body.

CRSLang resolves this:
- `request.args.get` → unambiguously phase 1 (query string only)
- `request.args.post` → unambiguously phase 2 (POST body only)
- `request.args` (combined) → inferred as phase 2 (the latest phase needed to have
  complete data)

@theseion: I would consider getting rid of `ARGS` entirely for CRSLang. I don't understand the value of having the combined collection.

This matches CRS best practice: rules that inspect `ARGS` should run in phase 2 to
ensure POST parameters are available.

**Trade-off:** This choice prevents inspecting `request.args` at phase 1 for GET-only
requests (where no POST body exists). Authors who need phase-1 inspection of query
string arguments should use `request.args.get` explicitly. The alternative — inferring
phase 1 and silently missing POST parameters — is a worse failure mode because it
would produce false negatives on POST-based attacks.

## Alternatives Considered

### A: Always require explicit phase

Keep phase as a required metadata field, even when it can be inferred.

**Rejected because:**
- Redundant in 80%+ of rules
- Creates a maintenance burden: if a rule's targets change, the author must remember
  to update the phase
- Mismatches between declared phase and actual field usage become possible (and are
  silent bugs)

### B: Infer phase, never allow explicit declaration

Always infer, error on TX-only rules that cannot be inferred.

**Rejected because:**
- TX-only rules are common in CRS (initialization, scoring, paranoia skipping)
- Some rules intentionally run in specific phases for ordering reasons
- Removes author control where it is legitimately needed

### C: Phase as a rule attribute, not metadata

```
@phase(request_headers)
rule 901001 { ... }
```

**Rejected because:**
- Adds syntax weight for something that is usually inferred
- Attributes/annotations are better reserved for cross-cutting concerns that affect
  many rules (e.g., `@deprecated`, `@disabled`)

## Consequences

### Positive

- Eliminates redundant metadata in the majority of rules
- **Catches phase/field mismatches at compile time** — a class of silent bugs in SecLang
  becomes a compile error in CRSLang
- Descriptive phase names (`request_headers`) are clearer than numeric phases (`1`)
- Lossless round-trip to SecLang: the compiler maps inferred phases to numbers
- Reduces cognitive load: authors think about *what* data to inspect, not *when* the
  engine should inspect it

### Negative

- Authors must understand inference rules to predict which phase their rule lands in
- TX-only rules require explicit phase, creating two code paths
- The `request.args` sub-phase resolution (→ phase 2) is a convention that must be
  documented and understood
- Short forms (`request` → `request_headers`) could cause confusion if someone means
  `request_body`

### Risks

- **Inference surprises** — an author adds a `tx.*` check to a rule that was previously
  inferred as phase 1. The phase does not change (TX is cross-phase), but the author
  might expect it to. Clear compiler output ("phase: request_headers (inferred from
  request.method)") mitigates this.
- **Future phases** — if a WAF engine adds new phases (e.g., a WebSocket phase),
  the field-to-phase mapping must be extended. The registry design from ADR-0001
  accommodates this.
- **Paranoia skip pattern** — the CRS pattern of duplicating TX-only skip rules across
  phases is eliminated entirely by guarded groups (ADR-0006). Instead of per-phase skip
  rules, a guarded group with `requires: paranoia >= N` handles conditional activation.
  The compiler generates the appropriate per-phase skip/marker pairs for SecLang output.
