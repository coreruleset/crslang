/**
 * SecLang mode for CodeMirror 5
 * Provides syntax highlighting for ModSecurity/Coraza WAF configuration language
 */

(function(mod) {
  if (typeof exports == "object" && typeof module == "object") // CommonJS
    mod(require("../../lib/codemirror"));
  else if (typeof define == "function" && define.amd) // AMD
    define(["../../lib/codemirror"], mod);
  else // Plain browser env
    mod(CodeMirror);
})(function(CodeMirror) {
  "use strict";

  CodeMirror.defineMode("seclang", function(config, parserConfig) {
    var keywords = {
      // SecLang directives
      "SecRule": true,
      "SecAction": true,
      "SecMarker": true,
      "SecRuleEngine": true,
      "SecRequestBodyAccess": true,
      "SecResponseBodyAccess": true,
      "SecRequestBodyLimit": true,
      "SecRequestBodyLimitAction": true,
      "SecResponseBodyLimit": true,
      "SecResponseBodyLimitAction": true,
      "SecRuleInheritance": true,
      "SecRuleRemoveById": true,
      "SecRuleRemoveByMsg": true,
      "SecRuleRemoveByTag": true,
      "SecRuleUpdateActionById": true,
      "SecRuleUpdateTargetById": true,
      "SecRuleUpdateTargetByMsg": true,
      "SecRuleUpdateTargetByTag": true,
      "SecComponentSignature": true,
      "SecWebAppId": true,
      "SecServerSignature": true,
      "SecArgumentSeparator": true,
      "SecCookieFormat": true,
      "SecUnicodeMapFile": true,
      "SecStatusEngine": true,
      "SecInterceptOnError": true,
      "SecContentInjection": true,
      "SecStreamInBodyInspection": true,
      "SecStreamOutBodyInspection": true,
      "SecRulePerf": true,
      "SecConnEngine": true,
      "SecConnReadStateLimit": true,
      "SecConnWriteStateLimit": true,
      "SecSensorId": true,
      "SecRulePerfTime": true,
      "SecHashEngine": true,
      "SecHashKey": true,
      "SecHashParam": true,
      "SecHashMethodRx": true,
      "SecHashMethodPm": true,
      "SecGsbLookupDb": true,
      "SecGuardianLog": true,
      "SecUploadKeepFiles": true,
      "SecUploadFileMode": true,
      "SecUploadFileLimit": true,
      "SecUploadDir": true,
      "SecTmpDir": true,
      "SecDataDir": true,
      "SecTmpSaveUploadedFiles": true,
      "SecArgumentsLimit": true,
      "SecPcreMatchLimit": true,
      "SecPcreMatchLimitRecursion": true,
      "SecDebugLog": true,
      "SecDebugLogLevel": true,
      "SecAuditEngine": true,
      "SecAuditLog": true,
      "SecAuditLogParts": true,
      "SecAuditLogType": true,
      "SecAuditLogFormat": true,
      "SecAuditLogStorageDir": true,
      "SecAuditLogFileMode": true,
      "SecAuditLogDirMode": true,
      "SecAuditLogRelevantStatus": true,
      "SecRemoteRules": true,
      "SecRemoteRulesFailAction": true,
      "SecCollectionTimeout": true,
      "SecHttpBlkey": true,
      "Include": true,
      "IncludeOptional": true
    };

    var actions = {
      "accuracy": true, "allow": true, "append": true, "auditlog": true, "block": true,
      "capture": true, "chain": true, "ctl": true, "deny": true, "deprecatevar": true,
      "drop": true, "exec": true, "expirevar": true, "id": true, "initcol": true,
      "log": true, "logdata": true, "maturity": true, "msg": true, "multiMatch": true,
      "noauditlog": true, "nolog": true, "pass": true, "pause": true, "phase": true,
      "prepend": true, "proxy": true, "redirect": true, "rev": true, "sanitiseArg": true,
      "sanitiseMatched": true, "sanitiseMatchedBytes": true, "sanitiseRequestHeader": true,
      "sanitiseResponseHeader": true, "setenv": true, "setrsc": true, "setsid": true,
      "setuid": true, "setvar": true, "severity": true, "skip": true, "skipAfter": true,
      "status": true, "tag": true, "t": true, "ver": true, "xmlns": true
    };

    return {
      startState: function() {
        return {
          inString: false,
          stringDelim: null,
          inComment: false
        };
      },

      token: function(stream, state) {
        if (state.inComment) {
          stream.skipToEnd();
          state.inComment = false;
          return "comment";
        }

        if (state.inString) {
          while (!stream.eol()) {
            if (stream.next() == state.stringDelim) {
              if (stream.peek() != state.stringDelim) {
                state.inString = false;
                state.stringDelim = null;
                return "string";
              }
              stream.next();
            }
          }
          return "string";
        }

        // Comments
        if (stream.match(/^#.*/)) {
          return "comment";
        }

        // Strings
        if (stream.match(/^"/)) {
          state.inString = true;
          state.stringDelim = '"';
          return "string";
        }
        if (stream.match(/^'/)) {
          state.inString = true;
          state.stringDelim = "'";
          return "string";
        }

        // Variables (start with &)
        if (stream.match(/^&[A-Z_][A-Z0-9_:]*/)) {
          return "variable";
        }

        // Operators (start with @)
        if (stream.match(/^@[a-zA-Z][a-zA-Z0-9_]*/)) {
          return "operator";
        }

        // Numbers
        if (stream.match(/^[0-9]+/)) {
          return "number";
        }

        // Keywords (directives)
        var word = stream.match(/^[a-zA-Z][a-zA-Z0-9_]*/);
        if (word) {
          word = word[0];
          if (keywords.hasOwnProperty(word)) {
            return "keyword";
          }
          if (actions.hasOwnProperty(word)) {
            return "builtin";
          }
          return "atom";
        }

        // Skip whitespace
        if (stream.match(/^\s+/)) {
          return null;
        }

        stream.next();
        return null;
      },

      lineComment: "#"
    };
  });

  CodeMirror.defineMIME("text/x-seclang", "seclang");
  CodeMirror.defineMIME("text/x-modsecurity", "seclang");
}); 