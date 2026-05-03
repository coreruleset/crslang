# CRSLang Evolution: Open Topics and Future Work

This document tracks structural concepts that have been identified but not yet
addressed by an ADR. Items are grouped by priority. Each item notes the affected
existing ADRs, the design questions to answer, and a suggested ADR number when one
is appropriate to write.

## Status Legend

- **Open** — identified, no ADR yet
- **Partially covered** — touched by existing ADRs but lacks a dedicated decision
- **Deferred** — intentionally postponed; adequate workaround exists

## Important Second-Tier (next batch when there's appetite)

### 1. Time and Duration Types — Open
Proposed: ADR-0017

Need a typed literal for durations (`5m`, `1h`, `30s`) used by:
- Variable expiry (`expire(tx.var, 3600)` from ADR-0004 — currently raw seconds)
- Rate-limit windows (when added)
- Interval-based reload policies (`reload interval: 5m` from ADR-0016)
- Session/IP collection TTLs

Design questions:
- Literal syntax: `5m`, `5 minutes`, `Duration("5m")`?
- Type coercion to/from integer seconds for legacy interop
- Calendar-aware durations (months, years) — needed or out of scope?
- Time literals (`2026-04-25`) — separate concern or part of this?

Cross-references: ADR-0004 (effects/expire), ADR-0016 (reload policies).

@theseion: depending on the complexity, ISO 8601 offers a well established syntax.

### 2. Versioning and Lifecycle Annotations — Open
Proposed: ADR-0018

Rules and groups have lifecycle (added, modified, deprecated, removed). The language
needs to express this without losing information when rules are renamed or split.

Design questions:
- Annotation syntax: `@since("4.0")`, `@deprecated("4.5", reason: "...", replacement: rule X)`,
  `@experimental`?
- Compile-time warnings on use of deprecated rules?
- Schema/language-version compatibility — how does a v2 ruleset declare it requires
  CRSLang spec version N?
- Migration trail: when rule X is split into X+Y, how is the lineage tracked?

Cross-references: ADR-0006 (rule management), ADR-0013 (documentation), ADR-0014
(package versioning).

@thesesion:
- deprecation might depend on something, e.g., a platform or framemwork; we might need a way to declare such dependencies
- migration trail: we could take a hint from git and give rules a parent - child relationship. We migth also want to consider
  sibling relationships (similar rule at different PL) and family relationships (same group, e.g. RCE) (though families will
  probably already be covered by another mechanism, e.g., grouping)

### 3. Function Signatures and Type System — Partially covered
Proposed: ADR-0019

ADR-0001 types the field namespace. ADR-0015 introduces typed macro signatures. But
the broader function library (transformations, predicates, custom operators) has no
documented signature scheme. Without this, "compile-time validation" is partial.

Design questions:
- Signature declaration syntax for built-in functions?
- Generic types? (`count<T>(collection<T>) -> int`)
- Type inference vs explicit annotations
- Engine-specific function variants (some functions exist on Coraza but not Cloud Armor)
- How does the type system interact with macros (ADR-0015) and external data
  (ADR-0016) types?

Cross-references: ADR-0001 (typed fields), ADR-0010 (multi-target — typed function
availability per target), ADR-0015 (macros), ADR-0016 (data types).

### 4. Test Attachment — Open
Proposed: ADR-0020 (or rolled into ADR-0013)

Rules need a way to reference their FTW (or successor) test cases. Today FTW lives
entirely outside the rule language, with cross-references by string rule ID.

Design questions:
- Where to declare: rule attribute (`tests: "tests/942100.yaml"`), separate `test_file`
  field, or attached doc-comment section?
- Inline tests vs sidecar files (preference: sidecar to keep rule files clean)
- How does the compiler verify test file existence at compile time?
- Test case format: keep FTW, design a successor, or be format-agnostic?

Cross-references: ADR-0013 (documentation), ADR-0014 (package distribution — tests
ship with the package).

@theseion: we also have test overrides and platform overrides. At least when transpiling, the result should contain the tests with the platform
overrides applied for the selected platform. Where would platform overrides live in the future?

### 5. Persistent Collections (IP, GLOBAL, SESSION, USER) — Partially covered
Proposed: ADR-0021

ADR-0001 mentions `ip.*`, `global.*`, `session.*` field prefixes but only for *read*
access. The *write* side (initcol, setvar with TTL, persistent expiry across requests)
is not designed. CRS uses these heavily for rate limiting, repeat-offender tracking,
session-bound state.

Design questions:
- Initialization: `init_collection(ip: client.ip)` — at rule level or as a setup
  declaration?
- TTL semantics: when is a collection variable purged?
- Cross-engine portability: SecLang has explicit collections; cloud WAFs may model
  this entirely differently (or not at all)
- Relationship to TX: TX is per-transaction, collections are persistent — does the
  syntax distinguish?
- Collection size limits, eviction policies

Cross-references: ADR-0001 (field namespace), ADR-0003 (boolean algebra has a
901320 example using initcol), ADR-0010 (multi-target portability).

@theseion: sounds like we need some form of abstraction for these. Even more so, as these things are engine-specific and shouldn't
be part of CRS (not in this form at least), if possible.

### 6. Capture Groups and Pattern Bindings — Partially covered
Proposed: ADR-0022

ADR-0002 mentions `capture()` as a predicate, but how regex captures are accessed
and composed is not designed. Today's SecLang `TX:0`–`TX:9` is the leak that needs
replacing.

Design questions:
- Access syntax: `matched.captures[1]`, `$1`, named captures?
- Scope: where are captures visible after the predicate that created them?
- Composition: if rule chain has multiple `matches()` predicates with captures, how
  do they coexist?
- Multiple matches: `find_all()` style operations that return a list of captures
- Backreferences within a regex: who handles them — the regex engine or the language?

Cross-references: ADR-0002 (pipeline operator with `capture()`), ADR-0003 (boolean
algebra — let bindings option mentioned this).

## Lower Priority (probably fine to defer)

### 7. Custom Operators and Extension Points — Deferred

Plugin system for adding detection algorithms beyond what ships with the engine.
Engines like Coraza already have plugin systems; the language could acknowledge
them or ignore them.

Likely resolution: ignore for now; document that engine-specific operators are
accessed via an `engine_call("operator_name", args)` escape hatch.

### 8. Conditional Compilation — Deferred

`@if(target == "seclang")` for target-specific tweaks. Mostly handled by the
multi-target compiler (ADR-0010) emitting target-specific output without explicit
language support. May become important if rule authors need to write target-specific
fallbacks.

### 9. Conditional Effects — Deferred

`then if (x) { block } else { pass }`. Current binary match-or-don't is probably
enough. CRS rules don't currently need branched effects.

### 10. Logging Structure — Adequately covered

ADR-0004 has `log()` as an effect. Rich structured logging (audit log enrichment,
custom fields, log levels per effect) is largely engine-side and out of scope for
the language.

### 11. Cross-Rule Communication / Signaling — Adequately covered

Beyond TX variables (which already serve this), there's no current need for an
event/signal model. CRS's existing TX-based communication works; a more formal
model can wait until a real use case emerges.

### 12. Rate Limiting Primitives — Open but engine-specific

Native syntax for "block if more than N requests from same IP in M minutes." Today
done via collection state and counters. A first-class rate limit construct would
be cleaner but is heavily engine-dependent (cloud WAFs have native rate limiting;
ModSecurity does not).

Likely resolution: model rate limiting as a pattern over collections (covered by
ADR-0021 once written), not as a dedicated language construct.

### 13. Severity and Category Customization — Open

Examples in earlier ADRs use custom severity values (`critical_payment`). Whether
authors can define custom severity/category values, and how those interact with
scoring (ADR-0011), is undecided.

Likely resolution: extend ADR-0011's scoring table to allow custom severities with
declared point values.

### 14. Rule reuse – Open

@theseion

When I thinking about conditional effects (9), I thought it might be interesting for longer cunjunctive logical chains. For example, imagine a chain of 4 conditions that increasingly narrow the scope of detection,
e.g.,:
- matches `wp-admin.php`
- has a suspect user agent string
- has a suspect IP
- uses suspect encoding

We could bail out after each condition and set an increasing scoring value. However, that is basically how CRS already works, only that rules stand for themselves.

Thinking on this further, a rule like "has suspect IP match" here could be reused in other contexts, even as a standalone rule. But we wouldn't want to duplicate the rule.
Instead, it would be nice if we could reference the rule and use it multiple times with a single definition. In the future, engines could cache the result of reused rules.
For the example above, if the first two conditions match, the IP rule would be run and the result would be cached (the reference could maybe use it's on scoring,
as it may differ for the reference case; the standalone rule might use +5, the chained rule might use +1 as part of the chain). The same IP rule later, standalone, would 
then run with the cached result (no second evaluation necessary).

Pros:
- rule reuse
- detection composition
- more precise detection (?)
- better "fuzzy" matching, i.e., composition of weak signals is easier because they can be grouped together
- Changing the rule definition affects all references

Cons:
- increased complexity in rule writing, reading, transpilation, and for engines
- the volume of rules to choose from makes it hard to choose a rule to reuse
- Changing the rule definition affects all references

## Cross-Cutting Notes

When writing ADRs from this list, check the following for consistency:

- **Documentation** — every new construct should declare how doc-comments and
  structured fields apply (ADR-0013)
- **Imports/Namespacing** — every new declarable entity (data, macro, etc.) lives
  in a namespace per ADR-0014
- **Multi-target** — every new feature must declare degradation behavior across
  ADR-0010's compilation targets
- **Phase placement** — IR-level concepts go in Phases 1-2; surface syntax in Phase 4

## Conventions

- ADR numbers are reserved when this TODO promises one (e.g., "Proposed: ADR-0017").
  Numbers are not used elsewhere until the ADR is written.
- An item moving from "Open" to written ADR is removed from this file and linked
  from the README ADR table.
- Items deferred indefinitely stay here as a record of explicit decisions not to
  pursue.
