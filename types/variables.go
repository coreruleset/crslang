package types

type Variables struct {
	Variables []string `yaml:"variables,omitempty"`
}

type Variable string

const (
	// Collections
	ARGS                   Variable = "ARGS"
	ARGS_GET               Variable = "ARGS_GET"
	ARGS_GET_NAMES         Variable = "ARGS_GET_NAMES"
	ARGS_NAMES             Variable = "ARGS_NAMES"
	ARGS_POST_NAMES        Variable = "ARGS_POST_NAMES"
	ARGS_POST              Variable = "ARGS_POST"
	ENV                    Variable = "ENV"
	FILES                  Variable = "FILES"
	GEO                    Variable = "GEO"
	GLOBAL                 Variable = "GLOBAL"
	IP                     Variable = "IP"
	MATCHED_VARS_NAMES     Variable = "MATCHED_VARS_NAMES"
	MATCHED_VARS           Variable = "MATCHED_VARS"
	MULTIPART_PART_HEADERS Variable = "MULTIPART_PART_HEADERS"
	PERF_RULES             Variable = "PERF_RULES"
	REQUEST_COOKIES_NAMES  Variable = "REQUEST_COOKIES_NAMES"
	REQUEST_COOKIES        Variable = "REQUEST_COOKIES"
	REQUEST_HEADERS_NAMES  Variable = "REQUEST_HEADERS_NAMES"
	REQUEST_HEADERS        Variable = "REQUEST_HEADERS"
	RESPONSE_HEADERS_NAMES Variable = "RESPONSE_HEADERS_NAMES"
	RESPONSE_HEADERS       Variable = "RESPONSE_HEADERS"
	RULE                   Variable = "RULE"
	SESSION                Variable = "SESSION"
	TX                     Variable = "TX"
	XML                    Variable = "XML"
	// Vars
	ARGS_COMBINED_SIZE           Variable = "ARGS_COMBINED_SIZE"
	AUTH_TYPE                    Variable = "AUTH_TYPE"
	DURATION                     Variable = "DURATION"
	FILES_COMBINED_SIZE          Variable = "FILES_COMBINED_SIZE"
	FILES_NAMES                  Variable = "FILES_NAMES"
	FILES_SIZES                  Variable = "FILES_SIZES"
	FILES_TMP_CONTENT            Variable = "FILES_TMP_CONTENT"
	FILES_TMPNAMES               Variable = "FILES_TMPNAMES"
	FULL_REQUEST                 Variable = "FULL_REQUEST"
	FULL_REQUEST_LENGTH          Variable = "FULL_REQUEST_LENGTH"
	HIGHEST_SEVERITY             Variable = "HIGHEST_SEVERITY"
	INBOUND_DATA_ERROR           Variable = "INBOUND_DATA_ERROR"
	MATCHED_VAR                  Variable = "MATCHED_VAR"
	MATCHED_VAR_NAME             Variable = "MATCHED_VAR_NAME"
	MODSEC_BUILD                 Variable = "MODSEC_BUILD"
	MSC_PCRE_LIMITS_EXCEEDED     Variable = "MSC_PCRE_LIMITS_EXCEEDED"
	MULTIPART_CRLF_LF_LINES      Variable = "MULTIPART_CRLF_LF_LINES"
	MULTIPART_FILENAME           Variable = "MULTIPART_FILENAME"
	MULTIPART_NAME               Variable = "MULTIPART_NAME"
	MULTIPART_STRICT_ERROR       Variable = "MULTIPART_STRICT_ERROR"
	MULTIPART_UNMATCHED_BOUNDARY Variable = "MULTIPART_UNMATCHED_BOUNDARY"
	OUTBOUND_DATA_ERROR          Variable = "OUTBOUND_DATA_ERROR"
	PATH_INFO                    Variable = "PATH_INFO"
	PERF_ALL                     Variable = "PERF_ALL"
	PERF_COMBINED                Variable = "PERF_COMBINED"
	PERF_GC                      Variable = "PERF_GC"
	PERF_LOGGING                 Variable = "PERF_LOGGING"
	PERF_PHASE1                  Variable = "PERF_PHASE1"
	PERF_PHASE2                  Variable = "PERF_PHASE2"
	PERF_PHASE3                  Variable = "PERF_PHASE3"
	PERF_PHASE4                  Variable = "PERF_PHASE4"
	PERF_PHASE5                  Variable = "PERF_PHASE5"
	PERF_SREAD                   Variable = "PERF_SREAD"
	PERF_SWRITE                  Variable = "PERF_SWRITE"
	QUERY_STRING                 Variable = "QUERY_STRING"
	REMOTE_ADDR                  Variable = "REMOTE_ADDR"
	REMOTE_HOST                  Variable = "REMOTE_HOST"
	REMOTE_PORT                  Variable = "REMOTE_PORT"
	REMOTE_USER                  Variable = "REMOTE_USER"
	REQBODY_ERROR                Variable = "REQBODY_ERROR"
	REQBODY_ERROR_MSG            Variable = "REQBODY_ERROR_MSG"
	REQBODY_PROCESSOR            Variable = "REQBODY_PROCESSOR"
	REQUEST_BASENAME             Variable = "REQUEST_BASENAME"
	REQUEST_BODY                 Variable = "REQUEST_BODY"
	REQUEST_BODY_LENGTH          Variable = "REQUEST_BODY_LENGTH"
	REQUEST_FILENAME             Variable = "REQUEST_FILENAME"
	REQUEST_LINE                 Variable = "REQUEST_LINE"
	REQUEST_METHOD               Variable = "REQUEST_METHOD"
	REQUEST_PROTOCOL             Variable = "REQUEST_PROTOCOL"
	REQUEST_URI                  Variable = "REQUEST_URI"
	REQUEST_URI_RAW              Variable = "REQUEST_URI_RAW"
	RESPONSE_BODY                Variable = "RESPONSE_BODY"
	RESPONSE_CONTENT_LENGTH      Variable = "RESPONSE_CONTENT_LENGTH"
	RESPONSE_CONTENT_TYPE        Variable = "RESPONSE_CONTENT_TYPE"
	RESPONSE_PROTOCOL            Variable = "RESPONSE_PROTOCOL"
	RESPONSE_STATUS              Variable = "RESPONSE_STATUS"
	RESOURCE                     Variable = "RESOURCE"
	SCRIPT_BASENAME              Variable = "SCRIPT_BASENAME"
	SCRIPT_FILENAME              Variable = "SCRIPT_FILENAME"
	SCRIPT_GID                   Variable = "SCRIPT_GID"
	SCRIPT_GROUPNAME             Variable = "SCRIPT_GROUPNAME"
	SCRIPT_MODE                  Variable = "SCRIPT_MODE"
	SCRIPT_UID                   Variable = "SCRIPT_UID"
	SCRIPT_USERNAME              Variable = "SCRIPT_USERNAME"
	SDBM_DELETE_ERROR            Variable = "SDBM_DELETE_ERROR"
	SERVER_ADDR                  Variable = "SERVER_ADDR"
	SERVER_NAME                  Variable = "SERVER_NAME"
	SERVER_PORT                  Variable = "SERVER_PORT"
	SESSIONID                    Variable = "SESSIONID"
	STATUS_LINE                  Variable = "STATUS_LINE"
	STREAM_INPUT_BODY            Variable = "STREAM_INPUT_BODY"
	STREAM_OUTPUT_BODY           Variable = "STREAM_OUTPUT_BODY"
	TIME                         Variable = "TIME"
	TIME_DAY                     Variable = "TIME_DAY"
	TIME_EPOCH                   Variable = "TIME_EPOCH"
	TIME_HOUR                    Variable = "TIME_HOUR"
	TIME_MIN                     Variable = "TIME_MIN"
	TIME_MON                     Variable = "TIME_MON"
	TIME_SEC                     Variable = "TIME_SEC"
	TIME_WDAY                    Variable = "TIME_WDAY"
	TIME_YEAR                    Variable = "TIME_YEAR"
	UNIQUE_ID                    Variable = "UNIQUE_ID"
	URLENCODED_ERROR             Variable = "URLENCODED_ERROR"
	USER                         Variable = "USER"
	USERAGENT_IP                 Variable = "USERAGENT_IP"
	USERID                       Variable = "USERID"
	WEBAPPID                     Variable = "WEBAPPID"
	WEBSERVER_ERROR_LOG          Variable = "WEBSERVER_ERROR_LOG"
	// Other vars
	MSC_PCRE_ERROR                   Variable = "MSC_PCRE_ERROR"
	MULTIPART_BOUNDARY_QUOTED        Variable = "MULTIPART_BOUNDARY_QUOTED"
	MULTIPART_BOUNDARY_WHITESPACE    Variable = "MULTIPART_BOUNDARY_WHITESPACE"
	MULTIPART_DATA_AFTER             Variable = "MULTIPART_DATA_AFTER"
	MULTIPART_DATA_BEFORE            Variable = "MULTIPART_DATA_BEFORE"
	MULTIPART_FILE_LIMIT_EXCEEDED    Variable = "MULTIPART_FILE_LIMIT_EXCEEDED"
	MULTIPART_HEADER_FOLDING         Variable = "MULTIPART_HEADER_FOLDING"
	MULTIPART_INVALID_HEADER_FOLDING Variable = "MULTIPART_INVALID_HEADER_FOLDING"
	MULTIPART_INVALID_PART           Variable = "MULTIPART_INVALID_PART"
	MULTIPART_INVALID_QUOTING        Variable = "MULTIPART_INVALID_QUOTING"
	MULTIPART_LF_LINE                Variable = "MULTIPART_LF_LINE"
	MULTIPART_MISSING_SEMICOLON      Variable = "MULTIPART_MISSING_SEMICOLON"
	MULTIPART_SEMICOLON_MISSING      Variable = "MULTIPART_SEMICOLON_MISSING"
	REQBODY_PROCESSOR_ERROR          Variable = "REQBODY_PROCESSOR_ERROR"
	REQBODY_PROCESSOR_ERROR_MSG      Variable = "REQBODY_PROCESSOR_ERROR_MSG"
	STATUS                           Variable = "STATUS"
)

var (
	allVariables = map[string]Variable{
		"ARGS":                             ARGS,
		"ARGS_GET":                         ARGS_GET,
		"ARGS_GET_NAMES":                   ARGS_GET_NAMES,
		"ARGS_NAMES":                       ARGS_NAMES,
		"ARGS_POST_NAMES":                  ARGS_POST_NAMES,
		"ARGS_POST":                        ARGS_POST,
		"ENV":                              ENV,
		"FILES":                            FILES,
		"GEO":                              GEO,
		"GLOBAL":                           GLOBAL,
		"IP":                               IP,
		"MATCHED_VARS_NAMES":               MATCHED_VARS_NAMES,
		"MATCHED_VARS":                     MATCHED_VARS,
		"MULTIPART_PART_HEADERS":           MULTIPART_PART_HEADERS,
		"PERF_RULES":                       PERF_RULES,
		"REQUEST_COOKIES_NAMES":            REQUEST_COOKIES_NAMES,
		"REQUEST_COOKIES":                  REQUEST_COOKIES,
		"REQUEST_HEADERS_NAMES":            REQUEST_HEADERS_NAMES,
		"REQUEST_HEADERS":                  REQUEST_HEADERS,
		"RESPONSE_HEADERS_NAMES":           RESPONSE_HEADERS_NAMES,
		"RESPONSE_HEADERS":                 RESPONSE_HEADERS,
		"RULE":                             RULE,
		"SESSION":                          SESSION,
		"TX":                               TX,
		"XML":                              XML,
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

func (v *Variables) ToString() string {
	result := ""
	for i, variable := range v.Variables {
		if i != len(v.Variables)-1 {
			result += string(variable) + "|"
		} else {
			result += string(variable)
		}
	}
	return result
}

func (v *Variables) AddVariable(variable string) error {
	/* 	constVariable, exists := allVariables[variable]
	   	if !exists {
	   		return fmt.Errorf("Invalid variable value: %s", variable)
	   	} */
	v.Variables = append(v.Variables, variable)
	return nil
}
