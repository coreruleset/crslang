# ADR-0012: Ruleset Initialization and Deployment Configuration

- **Status:** Proposed
- **Date:** 2026-04-25
- **Phase:** 0 (resolves the open risk in ADR-0008; precedes Phase 3 scoring work in ADR-0011)

## Context

CRS uses a two-layer initialization chain that predates any notion of a composable rule
language:

1. **`crs-setup.conf`** — a SecLang file the deployer edits. It uses `SecAction` with
   `setvar` to push deployment policy into TX variables: paranoia level, anomaly score
   thresholds, allowed HTTP methods, allowed content types, argument length limits, and
   similar detection parameters.

2. **901 rules** — hand-authored SecLang rules that run at request time (phase 1) to:
   - Verify `crs-setup.conf` was loaded (rule 901001 checks `&TX:crs_setup_version`)
   - Initialize scoring TX variables (`tx.critical_anomaly_score = 5`, etc.)
   - Propagate the paranoia level from TX into the skip/marker machinery
   - Initialize defaults for any config variable the deployer did not set

This design is a consequence of SecLang having no compile-time constants. Everything is a
runtime TX variable. CRSLang v2 has no such constraint.

### Problems

1. **Runtime init for compile-time policy** — paranoia level, scoring thresholds, and
   allowed methods do not change between requests. They are deployment policy. Expressing
   them as runtime `setvar` instructions is an artifact of the SecLang model, not a
   design choice.

2. **Boilerplate rules** — the 901 block exists solely to wire configuration into the rule
   engine. None of these rules contain detection logic. They are infrastructure, and
   infrastructure should be compiler-generated.

3. **Sanity check as a rule** — rule 901001 checks that `crs-setup.conf` was loaded. This
   is a workaround for SecLang's lack of include-time validation. A compiler or loader can
   enforce this guarantee statically.

4. **Split across two files** — deployment policy (crs-setup.conf) and runtime init (901
   rules) are separate files that must stay in sync. Adding a new config parameter
   requires touching both.

5. **TX variables as a configuration API** — parameters like `tx.allowed_methods` are TX
   variables that rules read as if they were static data. This works but loses type
   information, validates nothing at load time, and makes the configuration API implicit
   rather than declared.

### Scope Relative to ADR-0008

ADR-0008 draws the line between engine configuration (out of scope: body size limits,
PCRE limits, log paths) and rule-adjacent concerns (in scope: defaults, markers,
signatures). It notes in its Risks section that `crs-setup.conf` "is a mix of
configuration directives and `SecAction` rules that initialize TX variables" and defers
the question of what replaces the TX-initializing part.

This ADR resolves that deferred question. It defines the `config {}` block as the
replacement for crs-setup.conf's detection policy parameters and the 901 init rules.

## Decision

CRSLang introduces a **`config {}` block** that holds user-tunable deployment policy.
The compiler generates all necessary SecLang initialization from this block — no
hand-authored 901 rules, no two-file split.

### Separation of Concerns

Four distinct concepts that were previously conflated:

| Block | Defined by | Changed by | Examples |
|-------|-----------|------------|---------|
| `globals {}` | Ruleset authors | Nobody — read-only | Scoring severity table, component version |
| `config {}` | Ruleset authors (defaults) | Deployers (overrides) | Paranoia level, allowed methods, thresholds |
| `defaults {}` | Ruleset authors | Deployers (rarely) | Default action, default log behavior |
| Rule metadata | Rule authors | Nobody after authoring | `severity:`, `paranoia:`, `phase:` |

`globals {}` is ruleset policy — it defines the scoring model and is the same for every
deployment. `config {}` is deployment policy — it is the tunable surface the deployer
is expected to edit.

### The `config {}` Block

```
config {
  # Which rules are active (ADR-0011: paranoia: N attribute on rules uses this)
  paranoia_level = 1

  # Anomaly mode: accumulate scores and evaluate at end of phase.
  # immediate: each rule blocks independently (no scoring evaluation needed).
  blocking_mode  = anomaly

  # Detection policy — parameters used directly by detection rules
  allowed_methods       = ["GET", "HEAD", "POST", "OPTIONS"]
  allowed_content_types = [
    "application/x-www-form-urlencoded",
    "multipart/form-data",
    "text/xml",
    "application/xml",
    "application/json",
  ]
  allowed_http_versions  = ["HTTP/1.0", "HTTP/1.1", "HTTP/2", "HTTP/3"]
  restricted_extensions  = [".asa", ".asax", ".ascx", ".axd", ".backup",
                             ".bak", ".bat", ".cdx", ".cer", ".cfg", ".cmd"]
  restricted_headers     = ["Proxy-Authorization", "Lock-Token"]

  # Argument limits (rules that enforce sizes reference these)
  arg_limits {
    max_num_args      = 255
    max_arg_name_len  = 100
    max_arg_value_len = 400
    max_total_len     = 64000
  }

  # Upload limits (rules that enforce upload sizes reference these)
  upload_limits {
    max_file_size     = 10485760   # 10 MB
    max_combined_size = 10485760
  }

  # Anomaly score thresholds — when to block (see ADR-0011 for score model)
  scoring {
    inbound_threshold  = 5
    outbound_threshold = 4
    action             = deny(status: 403)
  }
}
```
@theseion: HTTP versions should use an enum, not strings. Otherwise, we'll have issues with things like `HTTP/2` vs `HTTP/2.0`
@theseion: would be nice to be able to use units for bytes, like in K8s resources. Good for reading and writing, the transpiler can take care of the target value.

The fields inside `config {}` are defined by the CRS ruleset, not by the CRSLang
language. The language specifies the block syntax and semantics; CRS defines which
parameters it supports. A ruleset can declare a schema for its `config {}` parameters
(types, defaults, documentation), which the compiler uses for validation.

#### Relationship to ADR-0011

ADR-0011 introduced `scoring_threshold {}` as a top-level block. That block is deployment
policy (different deployments set different thresholds) and belongs in `config {}`. The
`globals { scoring {} }` block from ADR-0011 remains as-is — it is a ruleset constant,
not tunable per deployment.

```
# Ruleset constant — belongs in globals (ADR-0011)
globals {
  scoring {
    critical = 5   # severity to score mapping
    warning  = 3
    notice   = 2
  }
}

# Deployment policy — belongs in config (this ADR)
config {
  scoring {
    inbound_threshold  = 5
    outbound_threshold = 4
    action             = deny(status: 403)
  }
}
```

### 901 Rules Replaced by the Compiler

The hand-authored 901 rules disappear. Each category of init work is handled differently:

| 901 purpose | CRSLang v2 replacement |
|---|---|
| Scoring variable init (`tx.critical_anomaly_score = 5`) | Compiler-generated from `globals { scoring {} }` |
| Score threshold init (`tx.inbound_anomaly_score_threshold = 5`) | Compiler-generated from `config { scoring {} }` |
| Detection policy init (`tx.allowed_methods = [...]`) | Compiler-generated from `config {}` |
| Sanity check (`count(tx.crs_setup_version) eq 0 → deny`) | Loader/compiler guarantee (see below) |
| Paranoia skip propagation | Compiler-generated marker pairs from guarded groups (ADR-0006) |
| Default init for unset variables | Defaults declared in the config schema; compiler emits defaults |

### Sanity Check Replaced by Loader Guarantee

Rule 901001 detects that `crs-setup.conf` was not included. This check is needed in
SecLang because the rule engine cannot know at load time whether an include happened.

In CRSLang, the `config {}` block is part of the ruleset itself (or an explicitly
imported file). The compiler can enforce at compile time that a `config {}` block exists
and that all required parameters have values. If the deployer forgets to configure, the
compile fails — no runtime surprise.

For SecLang output, the compiler can emit a version marker in the generated init block
and omit the runtime sanity check entirely, since correctness is guaranteed before
deployment.

### Compilation to SecLang

The `config {}` block compiles to a generated initialization block in the SecLang output.
This generated block replaces both `crs-setup.conf` and the hand-authored 901 rules:

```
# CRSLang config {} block
config {
  paranoia_level = 1
  allowed_methods = ["GET", "HEAD", "POST", "OPTIONS"]
  scoring {
    inbound_threshold = 5
    action            = deny(status: 403)
  }
}

# Compiled SecLang (generated — not hand-authored)
# ---- CRSLang generated initialization ----
SecAction \
    "id:9000001,\
    phase:1,\
    nolog,\
    pass,\
    setvar:tx.detection_paranoia_level=1,\
    setvar:tx.allowed_methods=GET HEAD POST OPTIONS,\
    setvar:tx.critical_anomaly_score=5,\
    setvar:tx.warning_anomaly_score=3,\
    setvar:tx.inbound_anomaly_score_threshold=5"
# ---- end generated initialization ----
```

The generated IDs occupy a reserved range (e.g., 9000001–9000999) that does not conflict
with authored rule IDs. Generated comments make the source of the init block clear.

### Workflow: The `setup.crs` Convention

For CRS deployments, the convention is:

```
crs/
├── setup.crs         # deployer edits config {} here
├── rules/
│   ├── REQUEST-901-INITIALIZATION.crs   # removed — no longer needed
│   ├── REQUEST-920-PROTOCOL-ENFORCEMENT.crs
│   └── ...
└── main.crs          # imports setup.crs and all rule files
```

`setup.crs` ships as `setup.crs.example` with commented-out defaults. The deployer
copies and edits it. If no `setup.crs` is provided, the compiler uses the defaults
declared in the config schema.
@theseion: In my mind, we wanted to avoid manual editing entirely and instead parameterize the transpilation, such that the transpiled output can be
deployed directly.

The analogy to the current `crs-setup.conf.example` is intentional — the workflow is
familiar, but the mechanism is now compile-time rather than runtime.

## Alternatives Considered

### A: Keep 901 rules as hand-authored CRSLang rules (status quo structure)

Translate 901 rules from SecLang to CRSLang syntax but keep them as explicit rules:

```
rule 901001 (phase: request_headers) {
  when count(tx.crs_setup_version) |> eq(0)
  then deny(status: 500)
}
```

**Rejected because:**
- The sanity check is a compiler concern, not a rule concern; no CRSLang user should
  author or read it
- All other 901 rules are pure boilerplate — translating them from SecLang to CRSLang
  syntax without conceptual improvement gains nothing
- The two-file split (setup.crs + 901-init.crs) recreates the crs-setup.conf problem
  in the new syntax

### B: Separate `setup.crs` file type with its own syntax

A dedicated file type that only allows `config {}` blocks, separate from rule files.

**Rejected because:**
- `config {}` inside a regular `.crs` file works fine and is simpler
- A separate file type requires a second parser or parser mode
- The distinction between "setup file" and "rules file" is a convention, not a language
  requirement — both use `.crs` extension, with `main.crs` importing `setup.crs`

### C: Runtime TX variables (current model, no change)

Keep `SecAction setvar` as the initialization mechanism, just with nicer CRSLang syntax.

**Rejected because:**
- Defers policy to runtime unnecessarily
- The compiler cannot validate configuration values before deployment
- Rules that read `tx.allowed_methods` are reading runtime state that could be missing
- The two-file coordination problem persists

### D: Configuration as rule metadata (inline)

Attach deployment policy to the group or ruleset that uses it:

```
group "protocol_enforcement" (allowed_methods: ["GET", "POST"]) {
  rule 920170 { ... }
}
```

**Rejected because:**
- Policy that applies across groups (e.g., paranoia level) would be duplicated
- The deployer cannot easily find and change deployment-wide settings
- Mixes authoring concerns (which rules exist) with deployment concerns (how they behave)

## Consequences

### Positive

- 901 rules as a file and a concept disappear from CRS authoring
- `crs-setup.conf` is replaced by a `config {}` block with typed, validated fields
- Compiler-time validation: missing or invalid config parameters are caught before deployment
- Deployer workflow stays familiar: copy `setup.crs.example`, edit values
- Generated SecLang init block is clearly labeled, not confused with hand-authored rules
- Detection policy parameters are co-located and typed — no more scattered `setvar` strings

### Negative

- Existing CRS users must migrate from `crs-setup.conf` to `setup.crs` (automated migration tool can handle most cases)
- The config parameter set is defined by the ruleset (CRS), not the language — if CRS adds a new parameter, it must be documented in the schema
- Generated rule IDs for the init block consume a range that must be reserved and documented

### Risks

- **Config schema versioning** — as CRS evolves, the set of supported `config {}` parameters may change. Deployers upgrading CRS must check whether their `setup.crs` uses deprecated parameters. The compiler should warn on unknown or removed parameters.
- **Default vs explicit** — if no `setup.crs` is provided, the compiler uses schema defaults. This is safe, but deployers may not realize the defaults are in effect. Clear documentation and a "warn on defaults used" compiler flag mitigate this.
- **Multi-target config semantics** — some config parameters (e.g., `scoring {}`) may not apply to all compilation targets (Cloud Armor, AWS WAF). The multi-target compiler (ADR-0010) must document which config parameters are target-specific.
