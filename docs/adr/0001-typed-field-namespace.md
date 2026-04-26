# ADR-0001: Typed Field Namespace

- **Status:** Proposed
- **Date:** 2026-04-13
- **Phase:** 1

## Context

CRSLang currently models SecLang's target system directly, splitting targets into two
separate concepts:

- **Variables** — standalone values like `REQUEST_METHOD`, `REMOTE_ADDR`,
  `RESPONSE_STATUS`. Represented as a flat enum with 120+ entries.
- **Collections** — map-like structures like `REQUEST_HEADERS`, `ARGS`, `TX` that
  require argument selectors (e.g., `REQUEST_HEADERS:Content-Type`). Represented with a
  `name`, `arguments[]`, `excluded[]`, and a `count` boolean flag.

This dual model creates several problems:

1. **Cognitive overhead** — rule authors must know which targets are variables vs
   collections. The distinction is a SecLang implementation detail, not a conceptual one.
2. **Verbose YAML** — accessing `REQUEST_HEADERS:Content-Type` requires 4 lines of YAML
   nesting.
3. **Weak typing** — all values are strings. There is no way to express that
   `REMOTE_ADDR` is an IP address or that `RESPONSE_STATUS` is an integer.
4. **Count as a flag** — `&COLLECTION` (count) is a boolean on the collection struct
   rather than a function applied to it.
5. **No validation** — nothing prevents combining invalid variable/collection names with
   incompatible operators.

## Decision

Replace variables and collections with a **unified typed field namespace** using
dot-notation for hierarchy and bracket notation for map access.

### Field Naming Convention

Fields use a hierarchical dot-separated namespace:

```
request.method          # was: REQUEST_METHOD
request.uri             # was: REQUEST_URI (full URI including query string)
request.uri_raw         # was: REQUEST_URI_RAW (uri before normalization)
request.filename        # was: REQUEST_FILENAME (path without query string)
request.basename        # was: REQUEST_BASENAME (last path component)
request.line            # was: REQUEST_LINE (raw request line)
request.protocol        # was: REQUEST_PROTOCOL
request.body            # was: REQUEST_BODY
request.body.length     # was: REQUEST_BODY_LENGTH
request.headers         # was: REQUEST_HEADERS (entire map)
request.headers["Host"] # was: REQUEST_HEADERS:Host
request.cookies         # was: REQUEST_COOKIES (entire map)
request.cookies["sid"]  # was: REQUEST_COOKIES:sid
request.args            # was: ARGS (all arguments)
request.args.get        # was: ARGS_GET
request.args.post       # was: ARGS_POST
request.args["id"]      # was: ARGS:id

response.status         # was: RESPONSE_STATUS
response.protocol       # was: RESPONSE_PROTOCOL
response.body           # was: RESPONSE_BODY
response.headers        # was: RESPONSE_HEADERS
response.content_type   # was: RESPONSE_CONTENT_TYPE

client.ip               # was: REMOTE_ADDR
client.port             # was: REMOTE_PORT
server.ip               # was: SERVER_ADDR
server.port             # was: SERVER_PORT

tx.anomaly_score                 # was: TX:anomaly_score
tx.inbound_anomaly_score_pl1     # was: TX:inbound_anomaly_score_pl1
tx.inbound_anomaly_score_pl2     # was: TX:inbound_anomaly_score_pl2
tx.inbound_anomaly_score_pl3     # was: TX:inbound_anomaly_score_pl3
tx.inbound_anomaly_score_pl4     # was: TX:inbound_anomaly_score_pl4
# Note: CRS stores scores as TX variables. ADR-0011 proposes replacing
# direct TX access with a first-class scoring model.

matched.var             # was: MATCHED_VAR
matched.var_name        # was: MATCHED_VAR_NAME
matched.vars            # was: MATCHED_VARS (map)

files.names             # was: FILES_NAMES
files.sizes             # was: FILES_SIZES
files.tmpnames          # was: FILES_TMPNAMES

multipart.filename      # was: MULTIPART_FILENAME
multipart.name          # was: MULTIPART_NAME
multipart.part_headers  # was: MULTIPART_PART_HEADERS

time.epoch              # was: TIME_EPOCH
time.year               # was: TIME_YEAR

rule.id                 # was: RULE:id (inside RULE collection)
```

### Type System

Each field has a declared type:

| Type           | Description                 | Example fields                                |
| -------------- | --------------------------- | --------------------------------------------- |
| `String`       | UTF-8 text                  | `request.method`, `request.uri`               |
| `Int`          | Integer                     | `response.status`, `server.port`              |
| `IP`           | IP address (v4/v6)          | `client.ip`, `server.ip`                      |
| `Bytes`        | Raw byte sequence           | `request.body`, `response.body`               |
| `Map[String]`  | String-keyed map of strings | `request.headers`, `request.args`             |
| `List[String]` | Ordered list of strings     | `request.headers.names`, `request.args.names` |
| `Bool`         | Boolean                     | (future: computed fields)                     |

### Map Access

Map-typed fields support:

- **Full map** — `request.headers` (iterates all values)
- **Key access** — `request.headers["Content-Type"]` (single value, type `String`)
- **Names** — `request.headers.keys` (list of key names)

**Key exclusions** (SecLang's `!ARGS:foo` — "remove foo from previously defined iterations over ARGS") are not part of the field syntax. Instead, target exclusions are handled by the rule management system
(ADR-0006):

```
# The rule targets all args
rule 942100 (severity: critical) {
  when request.args |> detect_sqli()
  then block
}

# User excludes a specific key in their customization file
exclude rule 942100 target request.args["passwd"]
```

This separates the rule definition (what to detect) from deployment customization (which
fields to skip), matching how CRS already works in practice.

### Count as a Function

Instead of a boolean flag on the collection, `count()` is a function:

```
# Current YAML:
collections:
  - name: TX
    arguments: [crs_setup_version]
    count: true

# New:
count(tx.crs_setup_version)    # returns Int
```

### Implementation

1. **Field Registry** — a Go struct/map defining all known fields with their names,
   types, and SecLang equivalents:

   ```go
   type FieldDef struct {
       Name       string     // "request.headers"
       Type       FieldType  // MapString
       SecLangVar string     // "REQUEST_HEADERS" (for import/export)
       IsCollection bool     // true (for backward compat mapping)
   }
   ```

2. **IR change** — conditions use `Field` instead of `[]Variable` + `[]Collection`:

   ```go
   type Target struct {
       Field     FieldRef   // parsed field reference
       KeyAccess *string    // bracket key, if any
       Excluded  []string   // excluded keys
   }
   ```

3. **Bidirectional mapping** — for SecLang import and YAML v1 compat:
   - `REQUEST_HEADERS:Content-Type` -> `request.headers["Content-Type"]`
   - `&TX:score` -> `count(tx.score)`
   - `!ARGS:foo` -> target exclusion via ADR-0006: `exclude rule ... target request.args["foo"]`

4. **YAML v2 syntax** (intermediate, before Phase 5 text syntax):

   ```yaml
   when:
     pipeline:
       field: request.headers["Content-Type"]
       steps:
         - fn: matches
           args: ["^application/json"]
   ```

## Alternatives Considered

### A: Keep variables and collections separate, just rename them

Simpler migration but preserves the dual-model complexity and does not enable type
checking.

### B: Object traversal (CEL-style)

`request.headers.get("Content-Type")` — method calls on typed objects. More powerful
but requires a more complex runtime model. Could be revisited if CRSLang ever needs
computed fields or custom methods.

### C: Flat identifiers (Wirefilter-style)

`http.request.headers` with no real hierarchy — dots are part of the name, not
traversal. Simpler parser but loses the structural grouping that aids readability and
completion.

## Consequences

### Positive

- Unified mental model: everything is a field
- Type checking becomes possible
- Shorter, more readable conditions
- Natural path to function composition (Phase 2)
- Bracket notation handles map key access cleanly
- Target exclusions separated from rule definitions (ADR-0006)

### Negative

- Large mapping table to maintain (120+ SecLang variables -> field names)
- Community must learn new field names (mitigated by documentation and autocomplete)
- Potential naming conflicts with future WAF engines
- Serialization format must handle both old and new field references during transition

### Risks

- **Naming bikeshed** — field names will be debated. Propose a naming convention early
  and get community buy-in before implementation.
- **Incomplete coverage** — some obscure SecLang variables may not map cleanly.
  Maintain an `engine.*` escape hatch namespace for engine-specific fields.
