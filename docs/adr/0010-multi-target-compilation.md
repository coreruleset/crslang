# ADR-0010: Multi-Target Compilation Model

- **Status:** Proposed
- **Date:** 2026-04-15
- **Phase:** 0 (foundational — shapes language design constraints)

## Context

CRSLang's primary compilation target today is SecLang (ModSecurity/Coraza `.conf` files).
However, the long-term vision is to compile CRSLang rules to multiple WAF engines:

| Target | Expression language | Organization |
|--------|-------------------|--------------|
| ModSecurity / Coraza | SecLang directives | OWASP |
| Google Cloud Armor | CEL expressions | Google Cloud |
| AWS WAF | JSON rule statements | AWS |
| Cloudflare WAF | Wirefilter expressions | Cloudflare |

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
   - If there's a workaround (e.g., OR → flag-setting rules for SecLang, see below),
     the compiler applies it automatically.
   - If there's no workaround, the compiler emits a clear error: "Feature X is not
     supported by target Y."
   - Rules that use unsupported features are flagged, not silently dropped.
   - A rule can declare the targets it is meant for with the `backends` metadata
     property (see below). The compiler skips the rule for any other target and lists
     it in the compile report as excluded by `backends`, not as an error.

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

| Feature | SecLang | Cloud Armor | AWS WAF | Cloudflare |
|---------|---------|-------------|---------|------------|
| Boolean AND | chain (workaround) | native | native | native |
| Boolean OR | flag-setting rules (workaround) | native | native | native |
| Regex matching | native | native | native | native |
| String transforms | native (40+) | limited | very limited | limited |
| IP matching | native | native | native | native |
| Anomaly scoring | via TX variables | not native | not native | not native |
| `detect_sqli()` | native (libinjection) | not available | managed rules | managed rules |
| `detect_xss()` | native (libinjection) | not available | managed rules | managed rules |
| Response body inspection | native | not available | not available | limited |
| `ctl:` runtime overrides | native | not available | not available | not available |
| File-based pattern lists | native (`@pmFromFile`) | not available | not available | not available |
| Rule exclusions | native | custom | custom | custom |

### SecLang-Specific Considerations

Because SecLang is the primary target and has the most constraints, the compiler must
handle several CRSLang features that SecLang cannot express directly:

**Boolean OR** — CRSLang's `A or B` has no SecLang equivalent. Emitting one rule per
branch, each with the rule's effects, would apply the effects twice when both branches
match. The compiler instead emits one flag-setting rule per branch and keeps the
original rule ID on a single rule that tests the flag and applies the action and effects
once:

```
# CRSLang
rule 100001 {
  when condition_a or condition_b
  then block { tx.custom_score += 5 }
}

# Compiled SecLang: branches only set a flag; the effects run once
SecRule <condition_a> "id:9000101,phase:2,pass,nolog,setvar:'tx.or_100001=1'"
SecRule <condition_b> "id:9000102,phase:2,pass,nolog,setvar:'tx.or_100001=1'"
SecRule TX:or_100001 "@eq 1" \
    "id:100001,phase:2,block,setvar:'tx.custom_score=+5'"
```

Branch rules take IDs from the reserved range for compiler-generated rules (ADR-0012), so
the authored ID still identifies the rule in logs and exclusions. Branch rules copy
`MATCHED_VAR` and `MATCHED_VAR_NAME` into TX so the acting rule can log what matched.
Mixed expressions such as `x and (y or z)` compile the same way: the `or` sub-expression
becomes a flag, and the acting rule chains `x` with the flag test.

When the branches share the same operator and transformations and differ only in the
field, the compiler can skip the flags and emit one rule with the targets joined by `|`
(SecLang's target separator), as long as this does not change how many times the
rule's effects run. SecLang runs non-disruptive actions for each matching target value,
so the merge is safe for a rule with no counted effects (like the one below) and not for
one that adds to a score:

```
# CRSLang
rule 100003 {
  when request.args |> matches("evil") or request.cookies |> matches("evil")
  then block
}

# Compiled SecLang: one rule, no flag
SecRule ARGS|REQUEST_COOKIES "@rx evil" "id:100003,phase:2,block"
```

**Macros / named compositions** — expanded inline during compilation. The compiled
SecLang contains the full transformation chain.

**Phase inference** — resolved to a numeric phase value in the SecLang output.

**Typed fields** — mapped back to SecLang variable/collection names.

**First-class scoring** — expanded to `setvar` actions in SecLang.

### Restricting Rules to Targets: `backends`

By default the compiler tries to compile every rule for the selected target. A rule
that only makes sense for some targets can say so in its metadata:

```
rule 9100200 {
  metadata {
    backends = ["seclang"]
  }
  when response.body |> contains("ORA-00933")
  then block
}
```

For any target not in the list, the compiler skips the rule and lists it in the compile
report, so users can see which rules their target will not run. A rule without
`backends` is compiled for every target, and unsupported features in it are errors as
described above. CRS rules must not exclude `seclang`, since the whole ruleset must
compile to SecLang (principle 2).

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

- **SecLang workaround correctness** — decomposing OR into flag-setting rules changes
  evaluation semantics (every branch rule runs; nothing is short-circuited). Must validate
  that the decomposed output is semantically equivalent.
- **Target drift** — cloud WAF targets change their expression languages over time.
  Compilation backends must be maintained per-target.
- **CRS adoption** — if CRS rules use features that only compile to SecLang (e.g.,
  `detect_sqli()`), the "write once" promise is weakened. Document which features are
  portable vs target-specific.
- **Performance** — compiler workarounds may produce less optimized output than native
  rules. For SecLang, measure the impact of OR decomposition on rule evaluation
  performance.
