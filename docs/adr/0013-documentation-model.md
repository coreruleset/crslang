# ADR-0013: Documentation Model for Rules, Groups, and Rulesets

- **Status:** Proposed
- **Date:** 2026-04-25
- **Phase:** 1 (IR-level concept; surface syntax depends on Phase 0 / Phase 4)

## Context

CRSLang v1's metadata model captures only what SecLang carries: `id`, `phase`, `msg`,
`severity`, `tags`, `version`, `revision`, `maturity`, plus a free-form `comment` field.
The IR struct in `types/metadata.go` reflects this directly.

This is enough to round-trip SecLang `.conf` files but not enough to describe a rule's
*intent*. Substantive per-rule documentation today lives in three disconnected places:

1. **SecLang `#` comments** above each rule — narrative, references, FP notes
2. **File-level banner comments** at the top of each `.conf` file — describes what the
   file contains, paranoia coverage, attack family
3. **External markdown** — the CRS website, `CHANGES.md`, individual ruleset docs

None of this is structured. None of it is queryable by tooling. None of it round-trips
through the IR. Migration tools that parse SecLang lose every comment except the immediate
pre-rule comment, and even that is opaque text.

### What this proposal must address

- **Per-rule narrative** — what the rule detects, why, references (CWE, OWASP, CVE),
  known false positives, related rules
- **Per-group narrative** — replaces the file-level banner comment. Groups (ADR-0006,
  ADR-0011) are now the cohesion unit, not files. A `group "sql_injection"` deserves the
  same kind of overview that `REQUEST-942-APPLICATION-ATTACK-SQLI.conf` has today.
- **Per-ruleset narrative** — what the distribution is, version policy, supported
  targets, paranoia model, scoring model. Today scattered across README/CHANGES/setup-conf
  comments.
- **Cross-cutting attribute documentation** — paranoia levels, severity levels, scoring
  categories: shared concepts that don't belong inside any one group but need
  authoritative documentation.
- **Documentation on rule-management directives** — `exclude rule` and `update rule`
  (ADR-0006) deserve a "why" comment, just as rules do.

### What's out of scope

- Documentation versioning beyond what's already tracked in git
- Localized/translated documentation (English-only for this proposal)
- Auto-generated documentation from rule structure (separate tooling concern)
- The CRS website's full content model (the website consumes this ADR's output, but
  its layout and styling are not the language's concern)

## Decision

CRSLang adopts a **layered documentation model** with two complementary mechanisms:

1. **Doc-comments** (Rust/Go style) — `///` comment lines that attach to the next AST
   node. Used for narrative prose: descriptions, rationale, examples, FP notes.
2. **Structured documentation fields** — typed metadata for machine-readable content:
   `references`, `cwe`, `owasp`, `false_positives`, `see_also`. Used by tooling
   (crs-toolchain, IDE hovers, audit log enrichment, `--explain` mode).

Documentation attaches at four tiers: **ruleset → group (nestable) → rule → directive**.

### Tier 1: Ruleset

The top-of-file narrative, replacing scattered README/CHANGES/comment content.

```
/// OWASP Core Rule Set — generic attack detection for web applications.
///
/// This ruleset covers OWASP Top 10 attack categories using anomaly scoring with
/// paranoia gating. Default deployment uses paranoia level 1; higher levels
/// trade false-positive rate for broader detection.
///
/// References:
///   - https://coreruleset.org/
///   - OWASP Top 10:2021
ruleset "OWASP_CRS" version "4.18.0" {
  description {
    paranoia_model = "1-4, see paranoia { } block below"
    scoring_model  = "anomaly accumulation, threshold in config { scoring {} }"
    targets        = ["seclang", "coraza"]
  }
}
```

The optional `description {}` block carries structured fields that don't fit naturally
in a doc-comment. This block is open-ended — fields are documented in the CRSLang spec
but adding new ones is a non-breaking change.

### Tier 2: Group (nestable)

Groups carry the file-level documentation that today lives in `.conf` banner comments.
**Groups can nest** for conceptual sub-grouping, but nesting is reserved for genuine
hierarchical cohesion — not for crossing orthogonal axes (paranoia, severity, phase).

```
/// SQL Injection detection.
///
/// Covers SQLi via libinjection (detectSQLi operator) and pattern-based detection
/// across query string, POST body, headers, and cookies. PL1 covers high-confidence
/// patterns; PL2-4 progressively widen detection at the cost of FP rate.
///
/// Common false-positive families:
///   - Admin UIs that legitimately accept SQL fragments
///   - Search forms with SQL keywords in the query
///   - JSON APIs that include SQL-like strings as data values
group "sql_injection" {
  category = "sqli"

  references     = ["CWE-89", "OWASP A1:2021"]
  false_positives = [
    "Admin UIs accepting SQL fragments — apply targeted exclusions",
    "Search forms with SQL keywords — exclude on the search field",
  ]

  /// libinjection-based detection. Uses detectSQLi() on URL-decoded args.
  group "libinjection_based" {
    rule 942100 (paranoia: 1, severity: critical) { ... }
    rule 942110 (paranoia: 2, severity: critical) { ... }
  }

  /// Pattern-based detection. Catches patterns libinjection misses (UNION-based,
  /// time-based blind, comment injection).
  group "pattern_based" {
    rule 942120 (paranoia: 2) { ... }
  }
}
```

#### Nesting Semantics

- Inner groups **inherit** outer group's structured fields (`category`, guards from
  `requires:`, default action). Inner can add or override.
- Doc-comments do **not** inherit — each group has its own narrative.
- Guards **compose** via conjunction: outer `(requires: paranoia >= 1)` + inner
  `(requires: paranoia >= 2)` = effective `paranoia >= 2`.
- Nesting depth is unlimited but expected to stay shallow (1-2 levels in practice).

### Tier 3: Rule

Per-rule documentation with both narrative and structured fields:

```
/// SQL injection via libinjection on URL-decoded query string and POST args.
///
/// References:
///   - CWE-89
///   - https://github.com/libinjection/libinjection
///
/// Known false positives:
///   - Search queries containing SQL keywords (use exclusion on the search field)
rule 942100 (
  severity:   critical,
  paranoia:   1,
  references: ["CWE-89"],
) {
  when request.args |> detect_sqli()
  then block
}
```

Doc-comments carry prose; the `references:` metadata field carries a machine-readable
list. Tooling can render the prose, link the references, and emit both in audit logs.

### Tier 4: Rule-Management Directives

`exclude rule` and `update rule` (ADR-0006) carry their own documentation — an
exclusion without an explanation is a maintenance hazard:

```
/// Admin UI legitimately accepts SQL fragments in the "query" field of the
/// query builder at /admin/query-builder. Reviewed by security team 2026-04-12.
exclude rule 942100 target request.args.post["query"]
  where request.uri |> starts_with("/admin/query-builder")

/// Bumped severity for our payment endpoints — SQLi here is an immediate
/// PCI-DSS incident, not just a CRS event.
update rule 942100 {
  severity = critical_payment
}
```

### Cross-Cutting Attribute Documentation

Concepts like paranoia levels and severity levels span every group and rule. They get
their own documentation construct, parallel to `globals { scoring {} }` from ADR-0011:

```
paranoia {
  /// High-confidence detections. Minimal false positives expected. Default for
  /// most deployments. Suitable for production traffic without exclusions.
  level 1 {
    expected_fp_rate = "<1%"
    coverage         = "common attack patterns, well-known signatures"
  }

  /// Broader coverage. Some false positives expected; targeted exclusions
  /// usually needed for application-specific traffic.
  level 2 {
    expected_fp_rate = "1-5%"
    coverage         = "obfuscated variants, less-common attack patterns"
  }

  level 3 { ... }
  level 4 { ... }
}
```

This block is metadata-only — it does not affect rule evaluation. Tooling and the
ruleset documentation generator consume it to explain what each paranoia level means
for deployers.

### IR Representation

```go
type Documentation struct {
    DocComment string                  // narrative prose from /// comments
    Fields     map[string]Value        // structured fields (references, cwe, etc.)
}

// Attached to each documentable AST node:
type Ruleset   struct { Doc *Documentation; ... }
type Group     struct { Doc *Documentation; Children []GroupChild; ... }
type Rule      struct { Doc *Documentation; Metadata Metadata; ... }
type Exclusion struct { Doc *Documentation; ... }
type Update    struct { Doc *Documentation; ... }

// Cross-cutting attribute docs (separate top-level blocks):
type ParanoiaDocs struct { Levels map[int]*Documentation }
type SeverityDocs struct { Levels map[Severity]*Documentation }
```

Doc-comments are stored as a single string (the `///` lines concatenated, leading marker
stripped). Tooling that wants to parse internal structure (e.g., extract a "References:"
section) does so over the prose; the language doesn't impose a sub-structure on the
narrative.

### Surface Syntax

Doc-comments shown above use `///` as the marker. The exact marker depends on the
Phase 0 language base decision (see ADR-0009):

- **Custom parser:** `///` doc-comments that attach to the next AST node, plus
  structured fields in the metadata clause. Native, lightweight.
- **HCL:** HCL does not have node-attached doc-comments. The same content goes into
  multi-line `description` and structured-field attributes:
  ```hcl
  rule "942100" {
    severity = "critical"
    description = <<-EOT
      SQL injection via libinjection...
    EOT
    references = ["CWE-89"]
  }
  ```

The IR is identical regardless of surface syntax. The compiler/parser populates the
`Documentation` struct from whichever surface form it sees.

### YAML Round-Trip

YAML serialization captures the documentation natively:

```yaml
kind: rule
metadata:
  id: 942100
  severity: critical
  paranoia: 1
documentation:
  description: |
    SQL injection via libinjection on URL-decoded query string and POST args.
  references:
    - CWE-89
    - https://github.com/libinjection/libinjection
  false_positives:
    - "Search queries containing SQL keywords"
when: ...
then: ...
```

The `description` field carries the doc-comment prose (multi-line YAML string). All
other structured fields carry through directly.

### Compilation to SecLang

SecLang's documentation model is limited to `msg:`, `tag:`, and `#` comments. The
compiler maps the documentation IR back as follows:

| CRSLang | SecLang output |
|---|---|
| Rule doc-comment first line | `msg:` |
| Rule doc-comment full text | `#` comment block above the rule |
| Group doc-comment | `#` banner above the group's compiled rules |
| Ruleset doc-comment | `#` banner at the top of the generated file |
| `references` field | `tag:cwe/89`, `tag:owasp/a1-2021`, etc. |
| `false_positives` field | `#` comment block above the rule |
| `paranoia { level N {} }` docs | dropped (no SecLang equivalent); preserved in YAML |

Round-trip is **not lossless** in the SecLang direction — SecLang cannot represent
nested groups or structured documentation. The CRSLang→SecLang export is deliberately
treated as a downgrade. Authors edit CRSLang; SecLang is the deployment target.

### Compilation to Other Targets

For Cloud Armor, AWS WAF, and Cloudflare (ADR-0010), documentation maps to whatever
labels, descriptions, or annotations the target supports. Most cloud WAFs accept a
`description` field per rule; structured fields become tags or labels where supported,
or are dropped where not. The compiler emits a target-specific manifest noting what
was preserved vs dropped.

### Tooling Contract

Documentation is consumed by:

- **crs-toolchain** — generates the CRS website from the IR's documentation nodes
- **LSP / IDE** — hover, go-to-definition for rule references, autocomplete for tag values
- **`--explain` CLI mode** — given a rule ID or matched event, prints the rule's
  description, references, and FP notes
- **Audit log enrichment** — embed rule description and references into audit log
  entries for downstream SOC tooling
- **Test report generation** — FP notes inform test case generation

Each consumer reads from the same `Documentation` IR struct. Adding a new tooling
consumer doesn't require schema changes.

## Alternatives Considered

### A: Inline metadata extension only

Push everything into the existing `metadata = "(" metadata_kv ")"` clause with multi-line
string support.

**Rejected because:**
- Rules with substantial documentation become unreadable inline
- Forces the metadata clause to handle multi-paragraph prose, which fights its parens-
  delimited design
- No good way to distinguish narrative from machine-readable fields

### B: Sidecar markdown files only

Keep rule definitions lean; put all documentation in `rules/942100.md` files linked by
filename or rule ID.

**Rejected because:**
- Drift risk: a rule's logic changes without its docs being updated
- No IR-level documentation for tooling — the `--explain` mode would have to load and
  parse markdown separately
- Two-file workflow is friction for authors editing both
- Loses the natural co-location that doc-comments provide

(Note: external markdown is still useful for *very* long-form content like the
paranoia model deep-dive. The CRS website can pull narrative from CRSLang docs and
augment it with sidecar markdown for content that doesn't fit per-rule.)

### C: Dedicated `doc {}` block inside each rule

```
rule 942100 (severity: critical) {
  doc {
    description = "..."
    references  = [...]
  }
  when ...
  then ...
}
```

**Considered viable but:**
- More vertical noise on every rule (every rule gets a `doc { }` even for short docs)
- Description still wants multi-line strings, which works but feels heavier than `///`
- Doc-comments are a more familiar pattern from Rust, Go, TypeScript, Python (`"""..."""`)

The structured fields still appear (references, false_positives) but as top-level
metadata or alongside metadata, not inside a separate `doc {}` block.

### D: Comments as data (parse SecLang `#` comments and surface them)

Treat existing `#` comments as the documentation source; no new constructs.

**Rejected because:**
- Loses structure entirely — everything is opaque text
- Tooling can't reliably extract references, FP notes, etc.
- SecLang comments are the *current* state and the problem this ADR solves

### E: Tags as the only structured documentation

Reuse the existing `tags:` array for everything: `tag:cwe/89`, `tag:fp/admin-ui`, etc.

**Rejected because:**
- Tags are flat strings; nested or multi-field structures don't fit
- Tag namespace pollution — mixing taxonomic tags with documentation tags makes both
  harder to query
- No room for prose narrative

## Consequences

### Positive

- Per-rule documentation gains a structured home with both narrative (doc-comments)
  and machine-readable (`references`, `false_positives`) channels
- Group-level documentation replaces opaque file headers with structured, queryable content
- Ruleset-level documentation centralizes information today scattered across
  README/CHANGES/setup-conf comments
- Cross-cutting concepts (paranoia, severity) get authoritative documentation that
  doesn't have to live inside any one group
- Tooling has a single, consistent IR contract — `crs-toolchain`, IDE LSP, audit log
  enrichment, and `--explain` all read from the same `Documentation` nodes
- Exclusions and updates carry "why" comments, reducing drift in custom deployments
- YAML serialization preserves everything; SecLang export degrades gracefully

### Negative

- More to write per rule. Authors who don't fill in docs leave the rule silent in
  tooling. Mitigated by: doc-comments are optional, lint can warn on missing docs for
  rules above a severity threshold.
- Two mechanisms (doc-comments + structured fields) means authors must understand the
  split — narrative vs machine-readable. Mitigated by: a clear style guide.
- Group nesting adds grammar surface. Mitigated by: nesting is optional, most groups
  stay flat, and the inheritance rules are simple.

### Risks

- **Convention drift** — without enforcement, doc-comments may become stale. The CRS
  test pipeline can include doc-coverage checks (e.g., every rule above severity
  warning must have a description).
- **Cross-target degradation** — SecLang and cloud WAFs preserve only a subset of
  documentation. Authors editing in CRSLang must understand that exporting to SecLang
  loses nested group docs and structured fields beyond `tag:`. The migration tooling
  should print a clear summary of what was preserved.
- **HCL surface complexity** — if HCL is chosen as the language base, doc-comments
  become structured `description` attributes, which is verbose for small docs. The
  trade-off is part of the Phase 0 language base decision (ADR-0009).
- **Schema growth** — structured fields (`references`, `false_positives`, `cwe`,
  `owasp`, etc.) will accumulate over time. The IR's `Fields map[string]Value` keeps
  this open-ended; the spec documents which fields are recognized by core tooling.
