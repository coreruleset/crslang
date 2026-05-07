# ADR-0014: Imports and Modularity

- **Status:** Proposed
- **Date:** 2026-04-25
- **Phase:** 4 (text syntax era; affects how multi-file rulesets compose)

## Context

CRS today is a directory of SecLang `.conf` files loaded by the engine via implicit
filesystem ordering and `Include` directives:

```
modsecurity.conf
├── Include /etc/modsec/crs-setup.conf
├── Include /etc/modsec/rules/REQUEST-901-INITIALIZATION.conf
├── Include /etc/modsec/rules/REQUEST-905-COMMON-EXCEPTIONS.conf
├── Include /etc/modsec/rules/REQUEST-911-METHOD-ENFORCEMENT.conf
└── ...
```

Several proposals reference this multi-file structure without formalizing it:
- ADR-0012 mentions `setup.crs`, `main.crs`, `rules/REQUEST-920-...crs`
- ADR-0013 documents groups as the cohesion unit, replacing per-file banners
- ADR-0006 references rules across files via ID for `exclude rule` and `update rule`

But the proposals do not specify:
- How a file imports another
- Whether load order matters semantically
- How rule IDs / group names / macros are namespaced
- What happens on conflict (duplicate rule IDs across files)
- How a deployer overrides an upstream group from a separately-distributed package
- How a downstream "custom rules" file extends the upstream CRS distribution

Without an explicit modularity model, every multi-file feature in the other ADRs is
under-specified.

### What multi-file rulesets need

1. **Explicit composition** — a deployer's `main.crs` lists what it imports, in what
   order. No filesystem-glob magic, no hidden ordering.
2. **Single global namespace for rule IDs** — rule 942100 means the same thing
   everywhere; uniqueness is enforced at compile time across all imported files.
3. **Hierarchical namespacing for groups, macros, and data** — `crs.sql_injection` is
   distinct from `myorg.custom_sqli`; downstream consumers don't accidentally collide
   with upstream package contents.
4. **Path resolution that's predictable** — relative to the importing file, with a
   defined search path for distributed packages.
5. **Override mechanism that's explicit** — a downstream file can modify upstream rules
   only via the rule-management directives (ADR-0006); silent shadowing is not allowed.
6. **Idempotent imports** — importing the same file twice produces one effect, not two.

## Decision

CRSLang adopts an **explicit `import` directive** with **hierarchical namespaces** for
groups, macros, and external data, and a **single flat namespace** for rule IDs.

### Import Syntax

```
import "path/to/file.crs"
import "path/to/file.crs" as alias
import package "owasp_crs/4.18"
```

Three forms:

1. **Path import** — loads the named file relative to the current file's directory.
2. **Aliased import** — loads the file and binds its top-level namespace to a local alias.
   Group references inside the importing file can then use the alias as a prefix.
3. **Package import** — loads a named package from the configured search path. Used for
   distributed rulesets (CRS itself, third-party rule packs).

### Path Resolution

Path imports resolve relative to the importing file's directory. No `..` traversal
allowed past the ruleset root, defined as the directory containing the entry-point
file (typically `main.crs`).

Package imports resolve via a search path defined in `config {}` (ADR-0012):

```
config {
  package_path = [
    "./vendor",
    "/usr/share/crslang/packages",
  ]
}
```

The first matching path wins. Package versions are encoded in the package name
(`owasp_crs/4.18`), not as a separate field — packages are distributed as versioned
directory trees.

### Namespaces

**Rule IDs are flat and global.** Every rule across every imported file shares a single
ID space. Duplicate IDs are a compile-time error. This matches SecLang behavior and
preserves the property that `exclude rule 942100` is unambiguous across the entire
ruleset.

**Groups, macros, and external data are namespaced** by the file or package they're
declared in. The fully-qualified name uses dot notation:

```
# in owasp_crs/rules/sql_injection.crs
group "sql_injection" {
  macro detect() = request.args |> detect_sqli()
  rule 942100 { ... }
}

# in deployer's main.crs
import package "owasp_crs/4.18" as crs

# crs.sql_injection refers to the imported group
exclude group crs.sql_injection target request.args.post["query"]

# crs.sql_injection.detect references the macro
rule 999100 (severity: critical) {
  when crs.sql_injection.detect()
  then block
}
```

Within a file, names resolve in this order: local declarations → aliased imports →
unaliased imports. If two unaliased imports define the same name, that is a
compile-time error. Warnings apply only when a higher-precedence name (for example,
a local declaration or aliased import) shadows a name that would otherwise be
available from an unaliased import.

### Load Order

The order of `import` statements in a file determines load order. Compilation processes
imports depth-first: each imported file is fully loaded (including its own imports)
before the importing file's body is evaluated.

For SecLang output, load order matters because rule evaluation order in SecLang follows
file order. The compiler emits rules in the order they were declared (across imports),
preserving the deployer's intent.

For other targets (Cloud Armor, AWS WAF), order is preserved where the target supports
it. Where the target evaluates rules in priority order rather than declaration order,
the compiler maps declaration order to priority numbers.

### Idempotency

Importing the same file twice is a no-op on the second import. The compiler tracks
imported files by canonical path; re-imports do not duplicate rules, groups, or macros.

This means diamond-shape import graphs (A imports B and C; B and C both import D) work
correctly: D is loaded once.

### Conflict Resolution

| Conflict | Resolution |
|---|---|
| Duplicate rule ID across imports | Compile-time error |
| Duplicate group name in same namespace | Compile-time error |
| Duplicate macro name in same namespace | Compile-time error |
| Macro/group name collision in same namespace | Compile-time error |
| Local name shadows imported name | Compile-time warning; local wins |
| Two unaliased imports define the same name | Compile-time error |


@theseion: duplicate aliases should raise a compile time error, right?

Conflicts are never silent. Authors must either rename, alias, or use the rule-management
directives (ADR-0006) to express override intent explicitly.

### Override Mechanism

Downstream consumers do not modify upstream files. They modify upstream behavior via
ADR-0006 directives:

```
# main.crs
import package "owasp_crs/4.18" as crs

# Disable a specific rule
exclude rule 942100

# Change a rule's severity
update rule 942100 {
  severity = critical_payment
}

# Disable an entire group
exclude group crs.sql_injection where request.uri |> starts_with("/admin")

# Add custom rules in a deployer-owned namespace
group "myorg.custom_payment_protection" {
  rule 9100100 (severity: critical) {
    when request.uri |> starts_with("/payment")
     and request.body |> contains("sql_keyword")
    then block
  }
}
```

The `exclude` and `update` directives in `main.crs` reference imported names by
fully-qualified path. The upstream package files are never modified.

@theseion: I don't understand this sentence. Where is the "fully-qualified path" in the example above?

### Distribution Model

A CRSLang package is a directory tree:

```
owasp_crs/4.18/
├── package.crs              # entry point — declares the package
├── setup.crs.example        # deployer-editable defaults
├── rules/
│   ├── sql_injection.crs
│   ├── xss.crs
│   └── ...
└── tests/                   # FTW test cases (per ADR-0013)
```

`package.crs` declares the package and re-exports the rule files:

```
package "owasp_crs" version "4.18.0" {
  description = "OWASP Core Rule Set"
  targets     = ["seclang", "coraza"]
}

import "rules/sql_injection.crs"
import "rules/xss.crs"
# ... etc.
```

A deployer's `main.crs` imports the package:

```
import package "owasp_crs/4.18" as crs
import "setup.crs"          # deployer-edited config

# Optional deployer-specific rules
import "custom/payment.crs"

```

@theseion: shouldn't we then ship a `main.crs.example` as well?
That being said, we could generate both `setup.crs` and `main.crs` directly as transpiler output instead.

I think I might be confused by how this package concept relates to the current structure of CRS. If I understand correctly, `main.crs` and `setup.crs` would be
inputs for the transpiler, correct? Based on the imports, the generated SecLang representations may, for example, only include some of the `9xx-REQUEST-xx.conf` files.
I think an example that explains how the package relates to the generated SecLang representation would be very helpful here.

### Compilation Output

For SecLang export, all imports are flattened into the output `.conf` file with banner
comments delineating the source:

```
# ===== imported from owasp_crs/4.18/rules/sql_injection.crs =====
SecRule ARGS "@detectSQLi" \
    "id:942100,..."
# ===== end owasp_crs/4.18/rules/sql_injection.crs =====
```

For YAML export, the import structure is preserved as nested documents or a manifest
file pointing to per-import YAML files (deployer's choice).

## Alternatives Considered

### A: Filesystem-glob discovery (current SecLang behavior)

Load all `*.crs` files in a directory, alphabetically.

**Rejected because:**
- Order depends on filenames; renaming a file changes evaluation order silently
- No way to express "use these files but not those"
- Conflicts with package distribution (which file in which directory is "the" package?)
- Hidden coupling: removing a file changes behavior elsewhere with no compile-time signal

### B: Single-file rulesets only

Require all rules in one file. No imports.

**Rejected because:**
- CRS today has ~50 files for good reason; 50,000-line files are unmaintainable
- Loses the cohesion benefit of grouping by attack family
- No package distribution model possible

### C: Implicit namespacing by file path

Every file's contents are automatically namespaced by its path; no explicit `as alias`.

**Considered viable but:**
- Long namespaces (`rules/protocol_enforcement/methods.crs.allowed_methods`) are awkward
- Renaming a file silently changes references
- Aliases give the deployer naming control without forcing path-based names

### D: First-class merge directive

`merge group crs.sql_injection { ... add more rules ... }` for downstream additions.

**Rejected because:**
- ADR-0006's `update` and `exclude` already cover modification
- Adding new rules in a deployer namespace (`myorg.custom_*`) is cleaner than merging
  into upstream namespaces
- Merging blurs the "upstream/downstream" boundary; explicit namespaces preserve it

## Consequences

### Positive

- Multi-file rulesets have a clear, predictable composition model
- Package distribution becomes a first-class concept (CRS itself, third-party packs)
- Conflicts are caught at compile time, never silent
- Override semantics (ADR-0006) work cleanly across package boundaries
- Diamond imports work correctly via idempotency
- Deployer customizations live in a separate namespace from upstream content

### Negative

- Authors must write `import` statements explicitly (vs SecLang's implicit `Include`)
- Namespace resolution rules (local → aliased → unaliased) require learning
- Package distribution requires a directory layout convention; not just "put files anywhere"

### Risks

- **Migration friction** — existing CRS deployments use a flat directory of `.conf`
  files with implicit ordering. The migration tool must generate a `main.crs` with
  the right import order to preserve behavior. This is automatable but needs care.
- **Package versioning conflicts** — if a deployer imports `owasp_crs/4.18` and a
  third-party pack imports `owasp_crs/4.17`, there's a conflict. The compiler must
  detect this and require the deployer to pin a single version.
- **Deep import graphs** — large rulesets with many cross-imports may be hard to
  reason about. Tooling (visualization, dependency analysis) helps; the language
  itself doesn't restrict graph shape beyond enforcing acyclicity.
