# ADR-0010: Multi-Target Compilation Model

- **Status:** Proposed
- **Date:** 2026-04-15
- **Phase:** 0 (foundational — shapes language design constraints)

## Context

CRSLang's primary compilation target today is SecLang (ModSecurity/Coraza `.conf` files).
However, the long-term vision is to compile CRSLang rules to multiple WAF engines:

| Target               | Expression language    | Organization |
| -------------------- | ---------------------- | ------------ |
| ModSecurity / Coraza | SecLang directives     | OWASP        |
| Google Cloud Armor   | CEL expressions        | Google Cloud |
| AWS WAF              | JSON rule statements   | AWS          |
| Cloudflare WAF       | Wirefilter expressions | Cloudflare   |

Each target has a different expression language, different capabilities, and different
limitations. CRSLang must be expressive enough to author rules once and compile them to
any supported target.

### The Constraint Question

Should CRSLang be constrained to the intersection of what all targets support? Or should
it be a superset, with per-target compilation producing errors or warnings for
unsupported features?

CRSLang v1 was constrained by SecLang — it modeled SecLang concepts directly (chains,
phases, transformation lists, action bags). This made SecLang round-tripping lossless but
prevented the language from expressing anything SecLang couldn't.

## Decision

**CRSLang is not constrained by any single compilation target.** The language expresses
the full range of WAF detection logic. Each compilation backend maps as much as it can
and reports clear errors for unsupported features.

### Design Principles

1. **Language over target** — CRSLang's syntax and semantics are defined by what makes
   rules readable and composable, not by what a specific engine supports.

2. **SecLang generation must be lossless for the CRS ruleset** — the entire OWASP CRS
   must compile to SecLang without loss. This is a practical constraint (CRS deploys on
   ModSecurity/Coraza today), not a language constraint.

3. **Graceful degradation** — when a CRSLang feature has no equivalent in a target:
   - By default, the compiler tries to compile all rules for the target backend, however if a rule has the properties `backends` specified, the compiler will only attempt to compile it if the target is included in the list. This allows CRS consumers to know which rules are expected in their target and which are not.
   - If there's a workaround (e.g., OR → multiple rules with scoring for SecLang),
     the compiler applies it automatically.
   - If there's no workaround, the compiler emits a clear error: "Feature X is not
     supported by target Y."
   - Rules that use unsupported features are flagged, not silently dropped.

4. **Target capability profiles** — each backend declares what it supports. Tooling can
   validate a ruleset against a target before compilation.

### Compilation Architecture

```
                            ┌──► SecLang (.conf)
                            │
CRSLang (.crs) ──► IR ─────┼──► Cloud Armor (CEL)
                            │
                            ├──► AWS WAF (JSON)
                            │
                            └──► Cloudflare (Wirefilter)
```

The IR (intermediate representation) is the shared substrate. It is richer than any
single target. Compilation backends are independent — adding a new target does not affect
existing backends or the language itself.

### Feature Support Matrix (Expected)

| Feature                  | SecLang                     | Cloud Armor   | AWS WAF       | Cloudflare    |
| ------------------------ | --------------------------- | ------------- | ------------- | ------------- |
| Boolean AND              | chain (workaround)          | native        | native        | native        |
| Boolean OR               | separate rules (workaround) | native        | native        | native        |
| Regex matching           | native                      | native        | native        | native        |
| String transforms        | native (40+)                | limited       | very limited  | limited       |
| IP matching              | native                      | native        | native        | native        |
| Anomaly scoring          | via TX variables            | not native    | not native    | not native    |
| `detect_sqli()`          | native (libinjection)       | not available | managed rules | managed rules |
| `detect_xss()`           | native (libinjection)       | not available | managed rules | managed rules |
| Response body inspection | native                      | not available | not available | limited       |
| `ctl:` runtime overrides | native                      | not available | not available | not available |
| File-based pattern lists | native (`@pmFromFile`)      | not available | not available | not available |
| Rule exclusions          | native                      | custom        | custom        | custom        |

### SecLang-Specific Considerations

Because SecLang is the primary target and has the most constraints, the compiler must
handle several CRSLang features that SecLang cannot express directly:

**Boolean OR** — CRSLang's `A or B` has no SecLang equivalent. The compiler decomposes
it into separate rules, each condition rule will set a variable if it matches, and a final rule checks the variable to apply the disruptive action. The final rule will use the CRSLang id:

```
# CRSLang
rule 100001 {
  when condition_a or condition_b
  then block { tx.anomaly_score += 5 }
}

# Compiled SecLang: two rules, same score target
# Rule IDs must be integers; the compiler allocates adjacent IDs deterministically.
SecRule ... "id:1000011,...,setvar:'tx.internal_rule_1000011=1'"
SecRule ... "id:1000012,...,setvar:'tx.internal_rule_1000012=1'"
SecRule TX:internal_rule_1000011||TX:internal_rule_1000012 "@eq 1" "id:100001,...,deny"
```

**Macros / named compositions** — expanded inline during compilation. The compiled
SecLang contains the full transformation chain.

**Phase inference** — resolved to a numeric phase value in the SecLang output.

**Typed fields** — mapped back to SecLang variable/collection names.

**First-class scoring** — expanded to `setvar` actions in SecLang.

### File Organization Metadata

CRS rules are organized into specific `.conf` files with a defined loading order
(REQUEST-901-INITIALIZATION.conf, REQUEST-920-PROTOCOL-ENFORCEMENT.conf, etc.). This
file structure matters for SecLang deployment because ModSecurity loads files in
order.

CRSLang preserves file organization as metadata in the IR. When compiling to SecLang,
the compiler emits rules into the correct files. When compiling to other targets that
don't have file-order dependencies, this metadata is ignored.

### Testing Strategy

Two testing layers support the multi-target model:

1. **Language-level tests** — validate CRSLang rules directly against the IR, without
   compilation. Check conditions, scoring, field access, type safety. Fast, no engine
   needed.

2. **E2E / integration tests** — compile to a target format (`.conf`, CEL, etc.),
   deploy to the engine, run `ftw` or equivalent test framework. Validates that
   compilation output works in the real engine. Essential for SecLang, where the
   compilation involves workarounds (OR decomposition, chain reconstruction).

Both layers are necessary. Language-level tests catch rule logic errors. E2E tests
catch compilation and engine-specific issues.

## Alternatives Considered

### A: Constrain to SecLang's capabilities

Limit CRSLang to what SecLang can express. No OR, no macros, no features beyond what
compiles 1:1.

**Rejected because:**

- Defeats the purpose of a new language
- Cloud WAF targets natively support OR, grouping, and richer expressions
- CRS maintainers cannot express patterns they can think of
- The language becomes a thin skin over SecLang rather than an improvement

### B: Constrain to the intersection of all targets

Limit CRSLang to features supported by every target.

**Rejected because:**

- The intersection is too small (string transforms, regex, IP matching, basic AND)
- Cloud WAF targets have very different capabilities
- Would prevent CRS from using features like `detect_sqli()` which only some engines
  support
- New targets could further shrink the intersection

### C: Per-target language profiles

Define subsets of CRSLang for each target (e.g., "CRSLang-SecLang", "CRSLang-CEL").

**Rejected because:**

- Fragments the language and the community
- Rule authors must know which profile they're writing for
- Tooling must handle multiple profiles
- Defeats the "write once" goal

## Consequences

### Positive

- CRSLang can express any WAF detection pattern, not just SecLang patterns
- New targets can be added without changing the language
- Rules that work on all targets are portable by default
- Target capability profiles enable early validation
- CRS can eventually deploy to cloud WAFs without manual translation

### Negative

- Compiler complexity increases with each target
- Feature support matrix must be maintained and documented
- Rule authors who care about portability must check target compatibility
- SecLang workarounds (OR decomposition) may produce less efficient output than
  hand-written SecLang

### Risks

- **SecLang workaround correctness** — decomposing OR into multiple rules changes
  evaluation semantics (each rule is independent, not short-circuited). Must validate
  that the decomposed output is semantically equivalent.
- **Target drift** — cloud WAF targets change their expression languages over time.
  Compilation backends must be maintained per-target.
- **CRS adoption** — if CRS rules use features that only compile to SecLang (e.g.,
  `detect_sqli()`), the "write once" promise is weakened. Document which features are
  portable vs target-specific.
- **Performance** — compiler workarounds may produce less optimized output than native
  rules. For SecLang, measure the impact of OR decomposition on rule evaluation
  performance.
