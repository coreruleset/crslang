package types

import "fmt"

type Variable struct {
	Name     VariableName `yaml:"name"`
	Excluded bool         `yaml:"excluded,omitempty"`
}

type VariableName string

const (
	ARGS_COMBINED_SIZE               VariableName = "ARGS_COMBINED_SIZE"
	AUTH_TYPE                        VariableName = "AUTH_TYPE"
	DURATION                         VariableName = "DURATION"
	FILES_COMBINED_SIZE              VariableName = "FILES_COMBINED_SIZE"
	FILES_NAMES                      VariableName = "FILES_NAMES"
	FILES_SIZES                      VariableName = "FILES_SIZES"
	FILES_TMP_CONTENT                VariableName = "FILES_TMP_CONTENT"
	FILES_TMPNAMES                   VariableName = "FILES_TMPNAMES"
	FULL_REQUEST                     VariableName = "FULL_REQUEST"
	FULL_REQUEST_LENGTH              VariableName = "FULL_REQUEST_LENGTH"
	HIGHEST_SEVERITY                 VariableName = "HIGHEST_SEVERITY"
	INBOUND_DATA_ERROR               VariableName = "INBOUND_DATA_ERROR"
	MATCHED_VAR                      VariableName = "MATCHED_VAR"
	MATCHED_VAR_NAME                 VariableName = "MATCHED_VAR_NAME"
	MODSEC_BUILD                     VariableName = "MODSEC_BUILD"
	MSC_PCRE_LIMITS_EXCEEDED         VariableName = "MSC_PCRE_LIMITS_EXCEEDED"
	MULTIPART_CRLF_LF_LINES          VariableName = "MULTIPART_CRLF_LF_LINES"
	MULTIPART_FILENAME               VariableName = "MULTIPART_FILENAME"
	MULTIPART_NAME                   VariableName = "MULTIPART_NAME"
	MULTIPART_STRICT_ERROR           VariableName = "MULTIPART_STRICT_ERROR"
	MULTIPART_UNMATCHED_BOUNDARY     VariableName = "MULTIPART_UNMATCHED_BOUNDARY"
	OUTBOUND_DATA_ERROR              VariableName = "OUTBOUND_DATA_ERROR"
	PATH_INFO                        VariableName = "PATH_INFO"
	PERF_ALL                         VariableName = "PERF_ALL"
	PERF_COMBINED                    VariableName = "PERF_COMBINED"
	PERF_GC                          VariableName = "PERF_GC"
	PERF_LOGGING                     VariableName = "PERF_LOGGING"
	PERF_PHASE1                      VariableName = "PERF_PHASE1"
	PERF_PHASE2                      VariableName = "PERF_PHASE2"
	PERF_PHASE3                      VariableName = "PERF_PHASE3"
	PERF_PHASE4                      VariableName = "PERF_PHASE4"
	PERF_PHASE5                      VariableName = "PERF_PHASE5"
	PERF_SREAD                       VariableName = "PERF_SREAD"
	PERF_SWRITE                      VariableName = "PERF_SWRITE"
	QUERY_STRING                     VariableName = "QUERY_STRING"
	REMOTE_ADDR                      VariableName = "REMOTE_ADDR"
	REMOTE_HOST                      VariableName = "REMOTE_HOST"
	REMOTE_PORT                      VariableName = "REMOTE_PORT"
	REMOTE_USER                      VariableName = "REMOTE_USER"
	REQBODY_ERROR                    VariableName = "REQBODY_ERROR"
	REQBODY_ERROR_MSG                VariableName = "REQBODY_ERROR_MSG"
	REQBODY_PROCESSOR                VariableName = "REQBODY_PROCESSOR"
	REQUEST_BASENAME                 VariableName = "REQUEST_BASENAME"
	REQUEST_BODY                     VariableName = "REQUEST_BODY"
	REQUEST_BODY_LENGTH              VariableName = "REQUEST_BODY_LENGTH"
	REQUEST_FILENAME                 VariableName = "REQUEST_FILENAME"
	REQUEST_LINE                     VariableName = "REQUEST_LINE"
	REQUEST_METHOD                   VariableName = "REQUEST_METHOD"
	REQUEST_PROTOCOL                 VariableName = "REQUEST_PROTOCOL"
	REQUEST_URI                      VariableName = "REQUEST_URI"
	REQUEST_URI_RAW                  VariableName = "REQUEST_URI_RAW"
	RESPONSE_BODY                    VariableName = "RESPONSE_BODY"
	RESPONSE_CONTENT_LENGTH          VariableName = "RESPONSE_CONTENT_LENGTH"
	RESPONSE_CONTENT_TYPE            VariableName = "RESPONSE_CONTENT_TYPE"
	RESPONSE_PROTOCOL                VariableName = "RESPONSE_PROTOCOL"
	RESPONSE_STATUS                  VariableName = "RESPONSE_STATUS"
	RESOURCE                         VariableName = "RESOURCE"
	SCRIPT_BASENAME                  VariableName = "SCRIPT_BASENAME"
	SCRIPT_FILENAME                  VariableName = "SCRIPT_FILENAME"
	SCRIPT_GID                       VariableName = "SCRIPT_GID"
	SCRIPT_GROUPNAME                 VariableName = "SCRIPT_GROUPNAME"
	SCRIPT_MODE                      VariableName = "SCRIPT_MODE"
	SCRIPT_UID                       VariableName = "SCRIPT_UID"
	SCRIPT_USERNAME                  VariableName = "SCRIPT_USERNAME"
	SDBM_DELETE_ERROR                VariableName = "SDBM_DELETE_ERROR"
	SERVER_ADDR                      VariableName = "SERVER_ADDR"
	SERVER_NAME                      VariableName = "SERVER_NAME"
	SERVER_PORT                      VariableName = "SERVER_PORT"
	SESSIONID                        VariableName = "SESSIONID"
	STATUS_LINE                      VariableName = "STATUS_LINE"
	STREAM_INPUT_BODY                VariableName = "STREAM_INPUT_BODY"
	STREAM_OUTPUT_BODY               VariableName = "STREAM_OUTPUT_BODY"
	TIME                             VariableName = "TIME"
	TIME_DAY                         VariableName = "TIME_DAY"
	TIME_EPOCH                       VariableName = "TIME_EPOCH"
	TIME_HOUR                        VariableName = "TIME_HOUR"
	TIME_MIN                         VariableName = "TIME_MIN"
	TIME_MON                         VariableName = "TIME_MON"
	TIME_SEC                         VariableName = "TIME_SEC"
	TIME_WDAY                        VariableName = "TIME_WDAY"
	TIME_YEAR                        VariableName = "TIME_YEAR"
	UNIQUE_ID                        VariableName = "UNIQUE_ID"
	URLENCODED_ERROR                 VariableName = "URLENCODED_ERROR"
	USER                             VariableName = "USER"
	USERAGENT_IP                     VariableName = "USERAGENT_IP"
	USERID                           VariableName = "USERID"
	WEBAPPID                         VariableName = "WEBAPPID"
	WEBSERVER_ERROR_LOG              VariableName = "WEBSERVER_ERROR_LOG"
	MSC_PCRE_ERROR                   VariableName = "MSC_PCRE_ERROR"
	MULTIPART_BOUNDARY_QUOTED        VariableName = "MULTIPART_BOUNDARY_QUOTED"
	MULTIPART_BOUNDARY_WHITESPACE    VariableName = "MULTIPART_BOUNDARY_WHITESPACE"
	MULTIPART_DATA_AFTER             VariableName = "MULTIPART_DATA_AFTER"
	MULTIPART_DATA_BEFORE            VariableName = "MULTIPART_DATA_BEFORE"
	MULTIPART_FILE_LIMIT_EXCEEDED    VariableName = "MULTIPART_FILE_LIMIT_EXCEEDED"
	MULTIPART_HEADER_FOLDING         VariableName = "MULTIPART_HEADER_FOLDING"
	MULTIPART_INVALID_HEADER_FOLDING VariableName = "MULTIPART_INVALID_HEADER_FOLDING"
	MULTIPART_INVALID_PART           VariableName = "MULTIPART_INVALID_PART"
	MULTIPART_INVALID_QUOTING        VariableName = "MULTIPART_INVALID_QUOTING"
	MULTIPART_LF_LINE                VariableName = "MULTIPART_LF_LINE"
	MULTIPART_MISSING_SEMICOLON      VariableName = "MULTIPART_MISSING_SEMICOLON"
	MULTIPART_SEMICOLON_MISSING      VariableName = "MULTIPART_SEMICOLON_MISSING"
	REQBODY_PROCESSOR_ERROR          VariableName = "REQBODY_PROCESSOR_ERROR"
	REQBODY_PROCESSOR_ERROR_MSG      VariableName = "REQBODY_PROCESSOR_ERROR_MSG"
	STATUS                           VariableName = "STATUS"
)

var (
	allVariables = map[string]VariableName{
		"ARGS_COMBINED_SIZE":               ARGS_COMBINED_SIZE,
		"AUTH_TYPE":                        AUTH_TYPE,
		"DURATION":                         DURATION,
		"FILES_COMBINED_SIZE":              FILES_COMBINED_SIZE,
		"FILES_NAMES":                      FILES_NAMES,
		"FILES_SIZES":                      FILES_SIZES,
		"FILES_TMP_CONTENT":                FILES_TMP_CONTENT,
		"FILES_TMPNAMES":                   FILES_TMPNAMES,
		"FULL_REQUEST":                     FULL_REQUEST,
		"FULL_REQUEST_LENGTH":              FULL_REQUEST_LENGTH,
		"HIGHEST_SEVERITY":                 HIGHEST_SEVERITY,
		"INBOUND_DATA_ERROR":               INBOUND_DATA_ERROR,
		"MATCHED_VAR":                      MATCHED_VAR,
		"MATCHED_VAR_NAME":                 MATCHED_VAR_NAME,
		"MODSEC_BUILD":                     MODSEC_BUILD,
		"MSC_PCRE_LIMITS_EXCEEDED":         MSC_PCRE_LIMITS_EXCEEDED,
		"MULTIPART_CRLF_LF_LINES":          MULTIPART_CRLF_LF_LINES,
		"MULTIPART_FILENAME":               MULTIPART_FILENAME,
		"MULTIPART_NAME":                   MULTIPART_NAME,
		"MULTIPART_STRICT_ERROR":           MULTIPART_STRICT_ERROR,
		"MULTIPART_UNMATCHED_BOUNDARY":     MULTIPART_UNMATCHED_BOUNDARY,
		"OUTBOUND_DATA_ERROR":              OUTBOUND_DATA_ERROR,
		"PATH_INFO":                        PATH_INFO,
		"PERF_ALL":                         PERF_ALL,
		"PERF_COMBINED":                    PERF_COMBINED,
		"PERF_GC":                          PERF_GC,
		"PERF_LOGGING":                     PERF_LOGGING,
		"PERF_PHASE1":                      PERF_PHASE1,
		"PERF_PHASE2":                      PERF_PHASE2,
		"PERF_PHASE3":                      PERF_PHASE3,
		"PERF_PHASE4":                      PERF_PHASE4,
		"PERF_PHASE5":                      PERF_PHASE5,
		"PERF_SREAD":                       PERF_SREAD,
		"PERF_SWRITE":                      PERF_SWRITE,
		"QUERY_STRING":                     QUERY_STRING,
		"REMOTE_ADDR":                      REMOTE_ADDR,
		"REMOTE_HOST":                      REMOTE_HOST,
		"REMOTE_PORT":                      REMOTE_PORT,
		"REMOTE_USER":                      REMOTE_USER,
		"REQBODY_ERROR":                    REQBODY_ERROR,
		"REQBODY_ERROR_MSG":                REQBODY_ERROR_MSG,
		"REQBODY_PROCESSOR":                REQBODY_PROCESSOR,
		"REQUEST_BASENAME":                 REQUEST_BASENAME,
		"REQUEST_BODY":                     REQUEST_BODY,
		"REQUEST_BODY_LENGTH":              REQUEST_BODY_LENGTH,
		"REQUEST_FILENAME":                 REQUEST_FILENAME,
		"REQUEST_LINE":                     REQUEST_LINE,
		"REQUEST_METHOD":                   REQUEST_METHOD,
		"REQUEST_PROTOCOL":                 REQUEST_PROTOCOL,
		"REQUEST_URI":                      REQUEST_URI,
		"REQUEST_URI_RAW":                  REQUEST_URI_RAW,
		"RESPONSE_BODY":                    RESPONSE_BODY,
		"RESPONSE_CONTENT_LENGTH":          RESPONSE_CONTENT_LENGTH,
		"RESPONSE_CONTENT_TYPE":            RESPONSE_CONTENT_TYPE,
		"RESPONSE_PROTOCOL":                RESPONSE_PROTOCOL,
		"RESPONSE_STATUS":                  RESPONSE_STATUS,
		"RESOURCE":                         RESOURCE,
		"SCRIPT_BASENAME":                  SCRIPT_BASENAME,
		"SCRIPT_FILENAME":                  SCRIPT_FILENAME,
		"SCRIPT_GID":                       SCRIPT_GID,
		"SCRIPT_GROUPNAME":                 SCRIPT_GROUPNAME,
		"SCRIPT_MODE":                      SCRIPT_MODE,
		"SCRIPT_UID":                       SCRIPT_UID,
		"SCRIPT_USERNAME":                  SCRIPT_USERNAME,
		"SDBM_DELETE_ERROR":                SDBM_DELETE_ERROR,
		"SERVER_ADDR":                      SERVER_ADDR,
		"SERVER_NAME":                      SERVER_NAME,
		"SERVER_PORT":                      SERVER_PORT,
		"SESSIONID":                        SESSIONID,
		"STATUS_LINE":                      STATUS_LINE,
		"STREAM_INPUT_BODY":                STREAM_INPUT_BODY,
		"STREAM_OUTPUT_BODY":               STREAM_OUTPUT_BODY,
		"TIME":                             TIME,
		"TIME_DAY":                         TIME_DAY,
		"TIME_EPOCH":                       TIME_EPOCH,
		"TIME_HOUR":                        TIME_HOUR,
		"TIME_MIN":                         TIME_MIN,
		"TIME_MON":                         TIME_MON,
		"TIME_SEC":                         TIME_SEC,
		"TIME_WDAY":                        TIME_WDAY,
		"TIME_YEAR":                        TIME_YEAR,
		"UNIQUE_ID":                        UNIQUE_ID,
		"URLENCODED_ERROR":                 URLENCODED_ERROR,
		"USER":                             USER,
		"USERAGENT_IP":                     USERAGENT_IP,
		"USERID":                           USERID,
		"WEBAPPID":                         WEBAPPID,
		"WEBSERVER_ERROR_LOG":              WEBSERVER_ERROR_LOG,
		"MSC_PCRE_ERROR":                   MSC_PCRE_ERROR,
		"MULTIPART_BOUNDARY_QUOTED":        MULTIPART_BOUNDARY_QUOTED,
		"MULTIPART_BOUNDARY_WHITESPACE":    MULTIPART_BOUNDARY_WHITESPACE,
		"MULTIPART_DATA_AFTER":             MULTIPART_DATA_AFTER,
		"MULTIPART_DATA_BEFORE":            MULTIPART_DATA_BEFORE,
		"MULTIPART_FILE_LIMIT_EXCEEDED":    MULTIPART_FILE_LIMIT_EXCEEDED,
		"MULTIPART_HEADER_FOLDING":         MULTIPART_HEADER_FOLDING,
		"MULTIPART_INVALID_HEADER_FOLDING": MULTIPART_INVALID_HEADER_FOLDING,
		"MULTIPART_INVALID_PART":           MULTIPART_INVALID_PART,
		"MULTIPART_INVALID_QUOTING":        MULTIPART_INVALID_QUOTING,
		"MULTIPART_LF_LINE":                MULTIPART_LF_LINE,
		"MULTIPART_MISSING_SEMICOLON":      MULTIPART_MISSING_SEMICOLON,
		"MULTIPART_SEMICOLON_MISSING":      MULTIPART_SEMICOLON_MISSING,
		"REQBODY_PROCESSOR_ERROR":          REQBODY_PROCESSOR_ERROR,
		"REQBODY_PROCESSOR_ERROR_MSG":      REQBODY_PROCESSOR_ERROR_MSG,
		"STATUS":                           STATUS,
	}
)

func VariablesToString(variables []Variable, separator string) string {
	result := ""
	for i, variable := range variables {
		if variable.Excluded {
			result += "!"
		}
		result += string(variable.Name)
		if i != len(variables)-1 {
			result += separator
		}
	}
	return result
}

func GetVariable(name string) (VariableName, error) {
	variable, exists := allVariables[name]
	if !exists {
		return "", fmt.Errorf("Invalid variable name: %s", name)
	}
	return variable, nil
}
