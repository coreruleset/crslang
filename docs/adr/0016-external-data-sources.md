# ADR-0016: External Data Sources

- **Status:** Proposed
- **Date:** 2026-04-25
- **Phase:** 3 (function composition era; data sources are referenced from pipelines)

## Context

CRS rules depend heavily on external data — lists, sets, databases — that the rule
references but does not contain inline:

- **IP allow/blocklists** — `tx.high_risk_country_codes`, customer-specific blocklists,
  threat intel feeds
- **Geo databases** — MaxMind GeoLite2 country/city/ASN lookups
- **Regex assemblies** — CRS's `regex-assembly` files compile multi-pattern regexes from
  structured `.ra` files
- **String sets** — allowed methods, allowed content types, restricted extensions
  (currently flattened into space-separated TX variables)
- **Key-value maps** — sometimes used for header→action mappings
- **Pre-computed indices** — pattern-matching tables for `@pm`-style multi-pattern
  matching

In SecLang, these are accessed via specialized operators:
- `@ipMatchFromFile /path/to/ips.txt`
- `@geoLookup` (uses `SecGeoLookupDb` from ADR-0008's engine config)
- `@pmFromFile /path/to/patterns.txt`
- `@rxFromFile /path/to/regex.txt`

CRSLang has no story for any of this. Rules in current ADRs implicitly assume data is
available (`client.ip in blocked_ips`, `request.uri |> matches_any(restricted_paths)`)
without specifying where `blocked_ips` or `restricted_paths` come from.

### What external data needs

1. **Declaration** — a typed reference to an external file with a known format
2. **Compile-time validation** — confirm the file exists and parses as the declared type
3. **Predictable lookup** — same syntax for "is X in this set?" regardless of whether
   the set is inline or external
4. **Hot reload semantics** — some data (threat intel) updates frequently; the language
   should express whether a data source is static or hot-reloadable
5. **Multi-target compilation** — SecLang has `@ipMatchFromFile`; cloud WAFs have their
   own list constructs (Cloud Armor's IP lists, AWS WAF IP sets); the compiler must
   map appropriately

## Decision

CRSLang adopts a **typed `data` declaration** that names an external source with a
declared type and format. Data sources are referenced from rules and macros as
typed values, with operations dispatched based on the data type.

### Declaration Syntax

```
data <name>
  (from <source> | values [...])
  [checksum  <"sha256:hex">]
  [signed_by <"key-file">]
  type <type>
  [format <format>]
  [reload <policy>]
  [managed_by_engine]
```

URL sources require at least one of `checksum` or `signed_by` (both may be present).
File and inline sources accept `checksum` as an optional integrity check; `signed_by`
is meaningful only for URL sources.

The `managed_by_engine` clause marks the data as provided and lifecycle-managed by
the engine (see [Engine-Managed vs Compiled-In](#engine-managed-vs-compiled-in)).
When present, the compiler validates the declaration but does not embed or fetch
the data; the engine resolves it at runtime. `managed_by_engine` is mutually
exclusive with `reload` (the engine owns the reload policy).

Examples:

```
# IP list from a file, hot-reloadable
data blocked_ips
  from "data/blocked_ips.txt"
  type ip_list
  format line_separated
  reload on_change

# Regex set from a CRS regex-assembly file
data sqli_patterns
  from "data/sqli.ra"
  type regex_set
  format regex_assembly

# Country code allowlist (inline)
data allowed_countries
  type string_set
  values ["US", "CA", "GB", "DE", "FR"]

# Geo database (managed by engine, declared for reference)
data geo_db
  from "data/GeoLite2-Country.mmdb"
  type geo_db
  format mmdb
  managed_by_engine

# Multi-pattern matching set (compiled from individual patterns)
data malware_signatures
  from "data/malware.txt"
  type string_set
  format line_separated
  reload static
```

### Data Types

| Type | Lookup operations | Example file format |
|---|---|---|
| `ip_list` | `in`, `not_in`, `contains_ip()` | `192.168.0.0/16\n10.0.0.0/8\n...` |
| `string_set` | `in`, `not_in`, `contains()` | `application/json\ntext/xml\n...` |
| `regex_set` | `matches_any()`, `matches_assembly()` | `.ra` file or `pattern\n...` |
| `kv_map` | `lookup(key)`, `has_key(key)` | `key=value\n...` (or JSON) |
| `geo_db` | `geo_country(ip)`, `geo_asn(ip)`, `geo_city(ip)` | MaxMind `.mmdb` |
| `pattern_table` | `pattern_match()` (Aho-Corasick) | `pattern\n...` |

The standard library's matching functions are typed: `matches_any` accepts a `regex_set`
and a `string`, `contains_ip` accepts an `ip_list` and an `ip`. Wrong-type arguments
are compile errors.

### Source Forms

Three source forms:

1. **`from "path"`** — load from a file. Path resolves relative to the declaring file
   (per ADR-0014's path resolution rules).
2. **`values [...]`** — inline literal values; no file involved. Useful for short
   allowlists.
3. **`from <url>`** — load from a URL (HTTPS only). URL sources require integrity
   verification via one of two mechanisms (or both):

   - **`checksum "sha256:<hex>"`** — pin the payload to a specific content hash.
     Suitable for *content-pinned* URLs where the bytes at the URL never change
     (e.g., versioned releases like `crs-rules/4.18.0/sqli.ra`). Compatible with
     `reload static` and `reload on_startup`. **Not** compatible with `reload
     interval` or `reload on_change`, because a live-updating feed will fail
     verification on the first reload — the checksum pins specific bytes, and
     updates necessarily change those bytes.

   - **`signed_by "<key-file>"`** — verify each fetch against a detached signature
     served alongside the payload (conventionally at `<url>.sig`). The key file is a
     path within the package (vendored as part of the deployment). Suitable for
     live feeds because the signer is pinned, not the content — each new payload
     can carry a fresh signature.

   At least one of `checksum` or `signed_by` is required for URL sources. Both may
   be declared together (the compiler verifies both).

   ```
   # Pinned: a versioned release URL where bytes never change at this path
   data sqli_assembly
     from     "https://example.com/crs-rules/4.18.0/sqli.ra"
     checksum "sha256:9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
     type     regex_set
     format   regex_assembly
     reload   on_startup

   # Live feed: signer is pinned, content updates between reloads
   data threat_feed
     from      "https://intel.example.com/feeds/blocklist.txt"
     signed_by "data/threat_feed.pubkey.pem"
     type      ip_list
     format    line_separated
     reload    interval: 1h
   ```

   For `checksum`-verified sources, the compiler verifies at compile time (and
   again on each reload, where the policy allows) and fails the build or reload on
   mismatch. For `signed_by` sources, the compiler verifies the detached signature
   at every fetch; a failed verification keeps the previously loaded copy in place
   and emits a runtime error log.

   @theseion: tool support for computing the SHA-256 of a new URL would be nice

### Format Hints

The `format` keyword tells the compiler how to parse the source. Defaults are inferred
from the file extension when possible:

| Format | Description |
|---|---|
| `line_separated` | One value per line, `#` comments and blanks ignored |
| `csv` | Comma-separated; columns described by `columns: [...]` |
| `json` | JSON array or object |
| `regex_assembly` | CRS `.ra` format (compiled to a single regex) |
| `mmdb` | MaxMind binary database |
| `yaml` | YAML structure |

Custom formats can be registered by ruleset authors via a parser plugin (out of scope
for this ADR; deferred to a future "extension points" ADR).

### Reload Policies

- `static` (default) — loaded once at compile time. Becomes part of the compiled output.
  Most performant; suitable for slow-changing data (allowed methods, file extensions).
- `on_startup` — loaded once at engine startup. Suitable for medium-volatility data
  (geo databases) that should not require recompilation.
- `on_change` — engine watches the file and reloads on change. Suitable for high-
  volatility data (threat intel feeds).
- `interval: <duration>` — engine reloads on a fixed schedule.

Not all targets support all reload policies. SecLang/Coraza supports static and
on_startup. Cloud WAFs typically require on_startup or external propagation. The
compiler reports which policies are supported per target.

### Lookup Syntax

External data is referenced as a value in expressions:

```
data blocked_ips from "data/ips.txt" type ip_list

rule 902100 (severity: warning) {
  when client.ip in blocked_ips
  then block
}
```

Or via typed operators:

```
data restricted_extensions from "data/restricted.txt" type string_set

rule 920360 (severity: warning) {
  when request.uri.basename |> ends_with_any(restricted_extensions)
  then block
}
```

For multi-pattern data:

```
data sqli_patterns from "data/sqli.ra" type regex_set format regex_assembly

rule 942100 (severity: critical) {
  when request.args |> matches_any(sqli_patterns)
  then block
}
```

### Compile-Time Validation

For each `data` declaration, the compiler:

1. Resolves the source path (or fetches the URL with checksum check)
2. Parses the file according to the declared format
3. Validates that contents match the declared type (e.g., `ip_list` must contain valid
   IPs/CIDRs; `regex_set` must contain syntactically valid regexes)
4. For static data, embeds the parsed contents in the compiled output
5. For dynamic data, embeds a reference (path + checksum) and the engine loads at runtime

Parse errors fail the compilation with file:line context.

### Compilation to SecLang

| CRSLang | SecLang output |
|---|---|
| `client.ip in <ip_list>` | `@ipMatchFromFile path` (file emitted alongside `.conf`) |
| `field \|> matches_any(<regex_set>)` | `@rxFromFile path` |
| `field in <string_set>` | `@pmFromFile path` |
| `geo_country(client.ip)` | `@geoLookup` + reference to `GEO:COUNTRY_CODE` |
| Inline `values [...]` (small) | Expanded to a regex literal |
| Inline `values [...]` (large) | Emitted as a generated file |

Static data with small inline values may compile to inline regex/pattern literals for
performance; the compiler picks based on size thresholds.

@theseion: `@rxFromFile` is not a SecLang operator (at least not in ModSecurity)

### Compilation to Other Targets

For Cloud Armor, AWS WAF, Cloudflare, and other cloud WAFs:

- **Cloud Armor**: IP lists become Address Groups; string/regex sets become custom
  rules with combined patterns; managed lists (GeoLite, threat intel) map to Cloud
  Armor's preconfigured WAF rules where available.
- **AWS WAF**: IP lists map to IPSets; regex sets map to RegexPatternSets; string
  sets map to rule logic.
- **Cloudflare**: IP lists map to IP Lists; regex sets compile inline (Wirefilter has
  no separate list type).

For each target, the compiler emits a separate manifest file describing what data
sources to upload/configure outside the rule logic itself. ADR-0010's multi-target
model handles the manifest format.

### Engine-Managed vs Compiled-In

Some data sources are managed by the engine, not the compiler:

- `geo_db` is configured at the engine level (MaxMind path is engine config per ADR-0008)
- Hot-reloadable data is loaded by the engine, not embedded by the compiler
- The compiler emits references; the engine resolves them at runtime

For these, `data` declarations are essentially type-checked references that document
what the engine must provide. The compiler verifies the file exists at compile time
but does not embed it.

### IR Representation

```go
type DataDecl struct {
    Doc       *Documentation        // ADR-0013
    Name      string
    Namespace string                // per ADR-0014
    Type      DataType
    Source    DataSource            // file, url, or inline
    Format    DataFormat
    Reload    ReloadPolicy
    Static    bool                  // true if compiled-in
}

type DataType interface { /* ip_list, string_set, regex_set, etc. */ }

type DataSource interface { /* FileSource, URLSource, InlineSource */ }
```

A `DataDecl` is referenced by name in the AST; the type system uses the declared
`DataType` to validate operations applied to it.

## Alternatives Considered

### A: Inline-only data (no external files)

Require all data to be inlined in the ruleset.

**Rejected because:**
- CRS regex assemblies and IP lists routinely have thousands of entries; inlining
  becomes unmaintainable
- Threat intel feeds must update faster than ruleset deployments
- Geo databases are gigabytes; cannot be inlined

### B: Untyped data references (string-based, like SecLang)

Reference files by string path with no type:

```
when client.ip ipMatchFromFile "data/ips.txt" then block
```

**Rejected because:**
- Loses compile-time validation (typos in paths, wrong file format)
- Loses type-checking on operations (passing a regex file to an IP operator)
- Discards the language's "typed everything" property

### C: Data as fields in the namespace

Treat external data as virtual fields under `data.*`:

```
when client.ip in data.blocked_ips
```

**Considered viable but:**
- Conflates the field namespace (request data, response data, TX state) with reference
  data (pre-loaded sets and lookups)
- ADR-0001's namespace is for runtime-changing data; data sources are static or
  engine-managed
- Cleaner to have a separate `data` declaration and reference

### D: Auto-discovery of data files

Look for `*.txt`, `*.ra`, `*.mmdb` files in conventional directories and auto-import.

**Rejected because:**
- Same problems as ADR-0014's rejected filesystem-glob discovery: silent ordering,
  brittle file naming, no compile-time signal when a file is missing
- Type cannot be inferred from filename reliably

## Consequences

### Positive

- IP lists, regex sets, string sets, and geo databases become first-class concepts
- Compile-time validation of file existence, format, and type
- Hot-reload semantics are explicit, not engine-implementation-dependent
- Multi-target compilation has a clear model for what each target supports
- Lookup syntax is uniform: same `in`, `matches_any`, etc. regardless of source
- The CRS regex-assembly subsystem becomes integrated into the language

### Negative

- Authors must declare data sources explicitly (vs SecLang's "just point to a file")
- Data type catalog must be extensible; standard library grows
- Format support requires parsers in the compiler (line-separated, regex_assembly,
  mmdb, JSON, YAML)

### Risks

- **Format proliferation** — every team has a slightly different data format. The
  language must support extension (custom format parsers) without becoming a parser
  zoo. Deferred extension-points ADR will address this.
- **Hot-reload semantics across targets** — `reload on_change` works for Coraza
  (file watch) but not Cloud Armor (requires API push). Compiler reports unsupported
  policies per target; deployer handles target-specific provisioning.
- **Source authenticity** — pulling data from URLs requires integrity verification
  to prevent supply-chain attacks. The language enforces this by requiring at least
  one of `checksum` (for content-pinned URLs) or `signed_by` (for live feeds) on
  every URL source; without one, the compile fails. The choice between the two is
  governed by the reload policy — pinned checksums are incompatible with live
  reloads, since content updates necessarily change the hash.
- **Large static data** — a 10MB IP list compiled into a SecLang `@ipMatchFromFile`
  is fine, but if the compiler tries to inline it as a regex literal, performance
  craters. Size thresholds for inline-vs-file decisions are needed.


@theseion: I can also imagine URL + checksum URL + reload interval / on_change. I don't know whether we can assume that every URL that needs refreshing also supplies
a signature. They may update their hash list though. Same goes for signatures: the signature may live at a URL different to the convention.
