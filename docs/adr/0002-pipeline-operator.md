# ADR-0002: Pipeline Operator for Composition

- **Status:** Proposed
- **Date:** 2026-04-13
- **Phase:** 3 (conditional — applies only if ADR-0009 chooses custom parser)

## Context

CRSLang currently separates transformations from operators. A condition applies an
ordered list of transformations to the target value, then tests the result with a single
operator:

```yaml
conditions:
  - collections:
      - name: REQUEST_HEADERS
        arguments: [User-Agent]
    transformations:
      - urlDecode
      - lowercase
      - removeWhitespace
    operator:
      name: rx
      value: "malicious-pattern"
```

This model has limitations:

1. **Transformations and operators are disconnected** — they sit in separate YAML
   sections even though they form a logical pipeline.
2. **No intermediate branching** — you cannot apply different operators to different
   transformation stages of the same field.
3. **Operators are special** — they have different syntax from transformations despite
   being conceptually similar (functions that take input and produce output).
4. **No composition** — you cannot nest or combine transformation results.

## Decision

Introduce a **pipeline operator** (`|>`) that chains functions left-to-right. Both
transformations and operators become functions connected by the pipeline.

### Syntax

```
field |> transform1() |> transform2() |> predicate("pattern")
```

The pipeline reads left-to-right:
1. Start with a field reference (from ADR-0001)
2. Each `|>` passes the output of the left side as the first argument to the right side
3. Intermediate functions (transformations) return the transformed value
4. Terminal functions (operators/predicates) return `bool`

### Examples

**Simple match:**
```
request.uri |> matches("^/admin")
```

**Transformation chain:**
```
request.headers["User-Agent"]
  |> url_decode()
  |> lowercase()
  |> remove_whitespace()
  |> matches("malicious-pattern")
```

**Negation:**
```
not(request.headers["Content-Length"] |> matches("^0?$"))
```

**Count:**
```
count(tx.crs_setup_version) |> eq(0)
```

**IP matching:**
```
client.ip |> ip_in_range("10.0.0.0/8", "172.16.0.0/12")
```

### Function Categories

#### Transformation Functions (return transformed value)

| Function | Input -> Output | SecLang equivalent |
|----------|-----------------|-------------------|
| `lowercase()` | String -> String | `t:lowercase` |
| `uppercase()` | String -> String | `t:uppercase` |
| `url_decode()` | String -> String | `t:urlDecode` |
| `url_decode_uni()` | String -> String | `t:urlDecodeUni` |
| `html_entity_decode()` | String -> String | `t:htmlEntityDecode` |
| `js_decode()` | String -> String | `t:jsDecode` |
| `css_decode()` | String -> String | `t:cssDecode` |
| `base64_decode()` | String -> String | `t:base64Decode` |
| `hex_decode()` | String -> String | `t:hexDecode` |
| `compress_whitespace()` | String -> String | `t:compressWhitespace` |
| `remove_whitespace()` | String -> String | `t:removeWhitespace` |
| `remove_nulls()` | String -> String | `t:removeNulls` |
| `remove_comments()` | String -> String | `t:removeComments` |
| `normalize_path()` | String -> String | `t:normalisePath` |
| `normalize_path_win()` | String -> String | `t:normalisePathWin` |
| `cmd_line()` | String -> String | `t:cmdLine` |
| `sql_hex_decode()` | String -> String | `t:sqlHexDecode` |
| `escape_seq_decode()` | String -> String | `t:escapeSeqDecode` |
| `trim()` | String -> String | `t:trim` |
| `length()` | String -> Int | `t:length` |
| `md5()` | String -> String | `t:md5` |
| `sha1()` | String -> String | `t:sha1` |

#### Predicate Functions (return bool)

| Function | Input Type | SecLang equivalent |
|----------|------------|-------------------|
| `matches(pattern)` | String -> Bool | `@rx` |
| `matches_global(pattern)` | String -> Bool | `@rxGlobal` |
| `eq(value)` | String/Int -> Bool | `@eq` / `@streq` |
| `gt(value)` | Int -> Bool | `@gt` |
| `ge(value)` | Int -> Bool | `@ge` |
| `lt(value)` | Int -> Bool | `@lt` |
| `le(value)` | Int -> Bool | `@le` |
| `contains(value)` | String -> Bool | `@contains` |
| `contains_word(values...)` | String -> Bool | `@pm` |
| `contains_word_from_file(path)` | String -> Bool | `@pmFromFile` (external word lists) |
| `begins_with(value)` | String -> Bool | `@beginsWith` |
| `ends_with(value)` | String -> Bool | `@endsWith` |
| `within(value)` | String -> Bool | `@within` |
| `ip_in_range(ranges...)` | IP -> Bool | `@ipMatch` |
| `ip_in_range_from_file(path)` | IP -> Bool | `@ipMatchFromFile` (external IP lists) |
| `detect_sqli()` | String -> Bool | `@detectSQLi` |
| `detect_xss()` | String -> Bool | `@detectXSS` |
| `validate_byte_range(range)` | Bytes -> Bool | `@validateByteRange` |
| `validate_url_encoding()` | String -> Bool | `@validateUrlEncoding` |
| `validate_utf8()` | String -> Bool | `@validateUtf8Encoding` |
| `verify_cc(pattern)` | String -> Bool | `@verifyCC` |
| `rbl(server)` | IP -> Bool | `@rbl` |
| `geo_lookup()` | IP -> Bool | `@geoLookup` |

### Type Checking

The pipeline enables static type checking:

```
# Valid: String -> String -> String -> Bool
request.uri |> lowercase() |> url_decode() |> matches("pattern")

# Invalid: Int -> String (type error)
response.status |> lowercase()

# Valid: String -> Int -> Bool
request.headers["Content-Length"] |> length() |> gt(100)
```

Type errors are caught at parse time with clear messages:
```
Error: lowercase() expects String input, got Int from 'response.status'
```

### IR Representation

```go
type Pipeline struct {
    Source  FieldRef       // Starting field
    Steps   []FunctionCall // Ordered transformation/predicate calls
    Negated bool           // Wrapping not()
}

type FunctionCall struct {
    Name      string       // "lowercase", "matches", etc.
    Arguments []Value      // Function arguments (patterns, numbers, etc.)
    ReturnType FieldType   // Computed from function signature
}
```

### YAML v2 Representation (intermediate)

```yaml
when:
  pipeline:
      field: request.headers["User-Agent"]
      steps:
        - fn: url_decode
        - fn: lowercase
        - fn: matches
          args: ["malicious-pattern"]
```

## Alternatives Considered

### A: Nested Function Calls (Wirefilter-style)

```
matches(lowercase(url_decode(request.headers["User-Agent"])), "pattern")
```

**Viable with macros (see ADR-0009).** The original concern about deep nesting dissolves
when named composition functions keep calls at 1-2 levels:

```
# Without macros: 4 levels deep, hard to read
matches(lowercase(url_decode(js_decode(request.args))), "pattern")

# With macros: readable
detect_sqli(normalize(request.args))
```

If HCL is chosen as the language base (ADR-0009 Option A), nested function calls are the
only available syntax — HCL's grammar cannot be extended with `|>`. In that path, macros
are essential to keep conditions readable.

**Trade-offs vs pipeline:**
- Reads inside-out for ad-hoc chains (without macros), which is opposite to CRS authors'
  mental model from SecLang's left-to-right `t:` chains
- With macros as the steady state, the reading direction difference is minimal
- `&&`/`||` boolean operators (HCL) are more familiar to developers than `and`/`or`
  keywords, but less readable for security engineers

### B: Method Chaining (CEL-style)

```
request.headers["User-Agent"].urlDecode().lowercase().matches("pattern")
```

**Considered viable but:**
- Requires fields to be objects with methods, complicating the type system
- Dot already used for field hierarchy (ADR-0001), adding method calls on the same
  syntax creates ambiguity: is `request.headers.names` a field or a method?
- Less visually distinct between "accessing data" and "transforming data"

### C: Unix Pipe (`|`)

```
request.uri | lowercase | matches("pattern")
```

**Rejected because:**
- `|` is commonly used for bitwise OR or alternatives in regex context
- Could create parsing ambiguity with future boolean OR syntax
- `|>` is established in Elixir, F#, OCaml, and Gleam for this exact purpose

## Consequences

### Positive

- Unified function model: transformations and operators are both "just functions"
- Natural left-to-right reading matches CRS authors' mental model
- Static type checking becomes possible
- Easy to add new functions without syntax changes
- Pipeline can be serialized to any format (YAML, text, JSON)

### Negative

- New syntax to learn (`|>` is not universally known)
- Function naming must be carefully designed — once published, names are hard to change
- Some SecLang operators don't fit the function model cleanly (e.g., `@inspectFile`
  which is really a side-effect)

### Risks

- **Function explosion** — resist adding functions for every conceivable transformation.
  Start with the functions needed by CRS rules and grow from there.
- **Performance implications** — pipeline representation must compile efficiently to
  target engines. Ensure the IR preserves enough information for backends to optimize.

### Notes on File-Based Functions and Regex Assembly

**Regex patterns** from `crs-toolchain` regex assembly (`.ra` files) are **inlined** at
build time. The toolchain compiles them into optimized patterns before CRSLang sees them.
There is no `matches_from_file()` — the final regex is embedded in the rule.

**Pattern match word lists** (`@pmFromFile`) and **IP range lists** (`@ipMatchFromFile`)
remain as external file references via `contains_word_from_file(path)` and
`ip_in_range_from_file(path)`. These lists can be thousands of entries and are data, not
logic — inlining them would make rules unreadable.

### The `t:none` Convention

SecLang's `t:none` transformation (which resets default transformations from
`SecDefaultAction`) has no equivalent in CRSLang and is not needed. Each rule explicitly
states its transforms via the pipeline or named macros, and the `defaults {}` block
(ADR-0008) does not inject hidden transformation chains. The problem `t:none` solved
does not exist in the new model.
