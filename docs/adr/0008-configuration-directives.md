# ADR-0008: Separation of Configuration Directives from Rule Language

- **Status:** Proposed
- **Date:** 2026-04-13
- **Phase:** 0 (foundational — scoping decision for the language)

## Context

SecLang mixes engine configuration and rule definitions in the same file format. A
typical CRS deployment interleaves ~60 configuration directives with rules:

```
SecRuleEngine DetectionOnly
SecRequestBodyAccess On
SecRequestBodyLimit 13107200
SecRequestBodyNoFilesLimit 131072
SecRequestBodyJsonDepthLimit 512
SecPcreMatchLimit 500000
SecPcreMatchLimitRecursion 500000

SecRule REQUEST_HEADERS:Content-Type "@rx ^application/json" \
    "id:200001,phase:1,pass,t:none,ctl:requestBodyProcessor=JSON"

SecAuditEngine RelevantOnly
SecAuditLog /var/log/modsec_audit.log
SecAuditLogParts ABCFHZ
SecAuditLogType Serial
```

CRSLang v1 models these as `ConfigurationDirective` — a struct with a `Name` (string
enum of ~60 types) and a `Parameter` (opaque string):

```go
type ConfigurationDirective struct {
    Kind      Kind
    Metadata  *CommentMetadata
    Name      ConfigurationDirectiveType  // "SecRuleEngine", "SecRequestBodyLimit", etc.
    Parameter string                      // "DetectionOnly", "13107200", etc.
}
```

### Problems

1. **Mixed concerns** — engine tuning (body size limits, PCRE limits, log paths) and
   detection logic (rules) live in the same file and the same IR. A rule author should
   not need to think about `SecPcreMatchLimitRecursion`.

2. **Deployment-specific** — configuration values change between environments (dev vs
   staging vs prod, nginx vs Apache vs Envoy). Rules do not. Mixing them means either
   duplicating rules across environments or templating config values out of rule files.

3. **Engine-specific** — many directives are ModSecurity-specific (`SecHashEngine`,
   `SecGuardianLog`, `SecStreamInBodyInspection`) or even version-specific (ModSec v2 vs
   v3). A language targeting multiple engines (Coraza, cloud WAFs) cannot standardize
   these.

4. **Opaque parameters** — the `Parameter` field is an untyped string. There is no
   validation that `SecRequestBodyLimit` receives an integer or that `SecRuleEngine`
   receives one of `On|Off|DetectionOnly`.

5. **No composition** — configuration directives cannot reference each other or be
   conditional. They are flat key-value pairs scattered across files.

### Directive Inventory

The ~60 configuration directives in CRSLang v1 break down as follows:

**Engine behavior (32 directives)** — how the WAF processes requests:

| Category | Directives |
|----------|-----------|
| Rule engine | `SecRuleEngine`, `SecConnEngine`, `SecStatusEngine` |
| Request body | `SecRequestBodyAccess`, `SecRequestBodyLimit`, `SecRequestBodyNoFilesLimit`, `SecRequestBodyInMemoryLimit`, `SecRequestBodyLimitAction`, `SecRequestBodyJsonDepthLimit` |
| Response body | `SecResponseBodyAccess`, `SecResponseBodyLimit`, `SecResponseBodyLimitAction`, `SecResponseBodyMimeType`, `SecResponseBodyMimeTypesClear` |
| PCRE tuning | `SecPcreMatchLimit`, `SecPcreMatchLimitRecursion` |
| Paths | `SecDataDir`, `SecTmpDir`, `SecChrootDir` |
| Connections | `SecConnReadStateLimit`, `SecConnWriteStateLimit`, `SecCollectionTimeout` |
| Hashing | `SecHashEngine`, `SecHashKey`, `SecHashParam`, `SecHashMethodRx`, `SecHashMethodPm` |
| Streaming | `SecStreamInBodyInspection`, `SecStreamOutBodyInspection` |
| Misc engine | `SecContentInjection`, `SecDisableBackendCompression`, `SecInterceptOnError`, `SecXmlExternalEntity`, `SecCacheTransformations`, `SecArgumentSeparator`, `SecArgumentsLimit`, `SecCookieFormat`, `SecCookieV0Separator` |

**Logging (12 directives)** — how and where to log:

| Directives |
|-----------|
| `SecAuditEngine`, `SecAuditLog`, `SecAuditLog2`, `SecAuditLogFormat`, `SecAuditLogParts`, `SecAuditLogRelevantStatus`, `SecAuditLogType`, `SecAuditLogStorageDir`, `SecAuditLogDirMode`, `SecAuditLogFileMode` |
| `SecDebugLog`, `SecDebugLogLevel` |
| `SecGuardianLog` |

**Uploads (4 directives):**

| Directives |
|-----------|
| `SecUploadDir`, `SecUploadFileLimit`, `SecUploadFileMode`, `SecUploadKeepFiles`, `SecTmpSaveUploadedFiles` |

**External data (3 directives):**

| Directives |
|-----------|
| `SecGeoLookupDb`, `SecGsbLookupDb`, `SecUnicodeMapFile` |

**Rule-adjacent metadata (5 directives)** — related to rules but not rules themselves:

| Directive | Purpose |
|-----------|---------|
| `SecComponentSignature` | Declares the application/version protected |
| `SecWebAppId` | Names the application context |
| `SecDefaultAction` | Sets default phase/action for subsequent rules |
| `SecMarker` | Label for `skipAfter` control flow |
| `SecServerSignature` | Overrides the server header |

**Operational (3 directives):**

| Directives |
|-----------|
| `SecRemoteRulesFailAction`, `SecRuleInheritance`, `SecRulePerfTime`, `SecSensorId`, `SecHttpBlKey` |

## Decision

**CRSLang does not model engine configuration directives.** The rule language describes
*what to detect*, not *how to run the engine*. Configuration belongs in a separate,
engine-specific configuration layer.

### What Stays in CRSLang

Only **rule-adjacent metadata** that is part of the rule semantics:

**1. Component signature → `component` declaration**

```
# SecLang: SecComponentSignature "OWASP_CRS/4.18.0-dev"
component "OWASP_CRS" version "4.18.0-dev"
```

This identifies the ruleset, not the engine. It belongs in the rule file as a header.

**2. Default actions → rule-level defaults block**

```
# SecLang: SecDefaultAction "phase:1,log,auditlog,pass"
defaults (phase: request_headers) {
  action: pass
  log(audit: true)
}
```

Or, if phase inference (ADR-0007) eliminates most explicit phase declarations, defaults
may simplify to just the action and effects:

```
defaults {
  action: pass
  log(audit: true)
}
```

Default actions apply to all subsequent rules in the file that do not override them.

**3. Markers → labels or groups (ADR-0006)**

```
# SecLang: SecMarker END_SQL_CHECKS
label END_SQL_CHECKS

# Or, preferably:
group sql_checks { ... }
```

**4. Web app ID → file-level metadata**

```
# SecLang: SecWebAppId "myapp"
app "myapp"
```

### What Leaves CRSLang

All engine behavior, logging, upload, path, PCRE, connection, hashing, and streaming
directives. These are **not part of the rule language**.

For CRS deployments, these settings move to an engine-specific configuration file:

```yaml
# engine.yaml (example — format is engine-specific, not standardized by CRSLang)
engine:
  mode: detection_only   # SecRuleEngine DetectionOnly

request_body:
  access: true           # SecRequestBodyAccess On
  limit: 13107200        # SecRequestBodyLimit
  no_files_limit: 131072 # SecRequestBodyNoFilesLimit
  json_depth_limit: 512  # SecRequestBodyJsonDepthLimit

response_body:
  access: true           # SecResponseBodyAccess On
  limit: 524288          # SecResponseBodyLimit
  mime_types:            # SecResponseBodyMimeType
    - text/html
    - application/json

pcre:
  match_limit: 500000          # SecPcreMatchLimit
  match_limit_recursion: 500000 # SecPcreMatchLimitRecursion

audit:
  engine: relevant_only  # SecAuditEngine RelevantOnly
  log: /var/log/modsec_audit.log
  format: JSON           # SecAuditLogFormat
  parts: ABCFHZ          # SecAuditLogParts

data:
  geo_db: /usr/share/GeoIP/GeoLite2-Country.mmdb  # SecGeoLookupDb
  unicode_map: /etc/modsec/unicode.mapping          # SecUnicodeMapFile
```

CRSLang does **not** standardize this format. Each engine (Coraza, ModSecurity, cloud
WAFs) defines its own configuration schema. CRSLang only defines the rule language.

### Runtime Overrides via `ctl:` Actions

Some configuration can be changed per-rule via `ctl:` actions in SecLang. These are
already handled by ADR-0004 as `configure {}` blocks inside rules:

```
rule 200001 {
  when request.headers["Content-Type"] |> matches("^application/json")
  then pass {
    configure(request_body_processor: json)
  }
}
```

These are **rule-scoped overrides**, not global configuration. They remain in the rule
language because they are conditional on the rule matching.

The set of `ctl:` settings that CRSLang supports in `configure {}` blocks is a small,
standardized subset that is expected to be portable across engines:

| CRSLang | SecLang `ctl:` |
|---------|---------------|
| `request_body_processor` | `ctl:requestBodyProcessor` |
| `request_body_access` | `ctl:requestBodyAccess` |
| `response_body_access` | `ctl:responseBodyAccess` |
| `rule_engine` | `ctl:ruleEngine` |
| `force_request_body` | `ctl:forceRequestBodyVariable` |
| `audit_engine` | `ctl:auditEngine` |
| `audit_log_parts` | `ctl:auditLogParts` |
| `rule_remove_by_id` | `ctl:ruleRemoveById` |
| `rule_remove_by_tag` | `ctl:ruleRemoveByTag` |

Engine-specific `ctl:` actions that are not portable go into an `engine_configure {}`
escape hatch or are dropped on export to non-supporting engines.

### SecLang Import/Export

**Import:** When importing SecLang `.conf` files, configuration directives are:
1. Recognized and parsed (for validation and reporting)
2. Emitted as comments or a separate config section — **not** mixed into the rule IR
3. A migration report lists extracted configuration with suggested engine-config
   equivalents

**Export:** When compiling CRSLang to SecLang:
1. Rule-adjacent metadata (`component`, `defaults`, `app`) is emitted as the
   corresponding SecLang directives
2. Engine configuration is **not** emitted — the user provides it separately in their
   engine config
3. `configure {}` blocks compile to `ctl:` actions

## Alternatives Considered

### A: Model all configuration directives with typed fields

Define a typed configuration schema in CRSLang covering all ~60 directives:

```
configure {
  rule_engine = detection_only
  request_body.limit = 13107200
  pcre.match_limit = 500000
  audit.log = "/var/log/modsec_audit.log"
}
```

**Rejected because:**
- Most directives are engine-specific and would not apply to all backends
- CRSLang would need to track and update configuration options for every engine version
- Conflates the rule language with engine deployment, the exact problem we are solving
- Configuration drift: CRSLang's config schema would perpetually lag behind engine
  releases

### B: Configuration as a separate CRSLang file type

Define a `.crsconfig` format alongside `.crs` rule files:

```
# engine.crsconfig
rule_engine: detection_only
request_body_limit: 13107200
```

**Rejected because:**
- Still requires CRSLang to define and maintain a configuration schema
- Adds a second file format to the CRSLang specification
- Each engine already has its own config format; adding another creates more work
  without clear benefit

### C: Keep configuration directives as-is (passthrough)

Import them, store them in the IR, export them unchanged.

**Rejected because:**
- Perpetuates the mixed-concerns problem
- The IR carries engine-specific data that has no meaning in a multi-engine world
- Prevents CRSLang from being engine-independent

## Consequences

### Positive

- CRSLang becomes purely about detection logic — cleaner language, smaller spec
- Engine independence: rules work across Coraza, ModSecurity, cloud WAFs without
  carrying engine-specific config
- Deployment configuration is managed with deployment tools (Ansible, Terraform, Helm)
  rather than embedded in rule files
- Smaller IR: ~60 directive types removed from the type system
- The `ConfigurationDirective` struct and its ~60-entry enum can be removed from the
  Go types (or relegated to the SecLang importer only)

### Negative

- CRS users accustomed to `SecRuleEngine On` at the top of their rule files must learn
  to put it elsewhere
- The SecLang importer must handle configuration directives gracefully (extract, report,
  not error)
- CRS documentation must clearly explain where configuration goes in a CRSLang-based
  deployment
- `SecDefaultAction` semantics (applying defaults to subsequent rules) must be carefully
  preserved in the `defaults {}` block

### Risks

- **CRS setup file** — `crs-setup.conf` is a mix of configuration directives and
  `SecAction` rules that initialize TX variables. The TX-initializing rules stay in
  CRSLang; the engine config lines are extracted. The migration tooling must handle
  this file specially.
- **`ctl:` completeness** — if an engine adds new `ctl:` actions, the CRSLang
  `configure {}` standardized set may lag behind. The `engine_configure {}` escape
  hatch handles this, but rules using it are not portable.
- **Documentation gap** — users need clear guidance on "I used to have one `.conf` file,
  now I have CRSLang rules + engine config." Provide migration guides per engine.
