package main

import (
	"slices"
	"testing"

	"github.com/antlr4-go/antlr/v4"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/parsing"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/types"
)

func TestLoadComment(t *testing.T) {
	testPayload := `#
# -- [[ Introduction ]] --------------------------------------------------------
#
# The OWASP ModSecurity Core Rule Set (CRS) is a set of generic attack
# detection rules that provide a base level of protection for any web
# application. They are written for the open source, cross-platform
# ModSecurity Web Application Firewall.
#
# See also:
# https://coreruleset.org/
# https://github.com/coreruleset/coreruleset
# https://owasp.org/www-project-modsecurity-core-rule-set/
#
`

	resultConfigs := []types.Configuration{}
	input := antlr.NewInputStream(testPayload)
	lexer := parsing.NewSecLangLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parsing.NewSecLangParser(stream)
	start := p.Configuration()
	var listener ExtendedSeclangParserListener
	antlr.ParseTreeWalkerDefault.Walk(&listener, start)
	resultConfigs = append(resultConfigs, listener.ConfigurationList.Configurations...)

	if len(resultConfigs) != 1 {
		t.Errorf("Expected 1 configuration, got %d", len(resultConfigs))
	}
	if len(resultConfigs[0].Directives) != 1 {
		t.Errorf("Expected 1 directive, got %d", len(resultConfigs[0].Directives))
	}
	_, ok := resultConfigs[0].Directives[0].(types.CommentMetadata)
	if !ok {
		t.Errorf("Expected comment, got %T", resultConfigs[0].Directives[0])
	}
	if resultConfigs[0].Directives[0].(types.CommentMetadata).Comment != testPayload {
		t.Errorf("Expected comment %s, got %s", testPayload, resultConfigs[0].Directives[0].(types.CommentMetadata).Comment)
	}
}

func TestLoadSecAction(t *testing.T) {
	testPayload := `
# Initialize anomaly scoring variables.
# All _score variables start at 0, and are incremented by the various rules
# upon detection of a possible attack.

SecAction \
    "id:901200,\
    phase:1,\
    pass,\
    t:none,\
    nolog,\
    tag:'OWASP_CRS',\
    ver:'OWASP_CRS/4.0.1-dev',\
    setvar:'tx.blocking_inbound_anomaly_score=0',\
    setvar:'tx.detection_inbound_anomaly_score=0',\
    setvar:'tx.inbound_anomaly_score_pl1=0',\
    setvar:'tx.inbound_anomaly_score_pl2=0',\
    setvar:'tx.inbound_anomaly_score_pl3=0',\
    setvar:'tx.inbound_anomaly_score_pl4=0',\
    setvar:'tx.sql_injection_score=0',\
    setvar:'tx.xss_score=0',\
    setvar:'tx.rfi_score=0',\
    setvar:'tx.lfi_score=0',\
    setvar:'tx.rce_score=0',\
    setvar:'tx.php_injection_score=0',\
    setvar:'tx.http_violation_score=0',\
    setvar:'tx.session_fixation_score=0',\
    setvar:'tx.blocking_outbound_anomaly_score=0',\
    setvar:'tx.detection_outbound_anomaly_score=0',\
    setvar:'tx.outbound_anomaly_score_pl1=0',\
    setvar:'tx.outbound_anomaly_score_pl2=0',\
    setvar:'tx.outbound_anomaly_score_pl3=0',\
    setvar:'tx.outbound_anomaly_score_pl4=0',\
    setvar:'tx.anomaly_score=0'"`

	resultConfigs := []types.Configuration{}
	input := antlr.NewInputStream(testPayload)
	lexer := parsing.NewSecLangLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parsing.NewSecLangParser(stream)
	start := p.Configuration()
	var listener ExtendedSeclangParserListener
	antlr.ParseTreeWalkerDefault.Walk(&listener, start)
	resultConfigs = append(resultConfigs, listener.ConfigurationList.Configurations...)

	if len(resultConfigs) != 1 {
		t.Errorf("Expected 1 configuration, got %d", len(resultConfigs))
	}
	if len(resultConfigs[0].Directives) != 1 {
		t.Errorf("Expected 1 directive, got %d", len(resultConfigs[0].Directives))
	}
	if resultConfigs[0].Directives[0].(*types.SecAction).Id != 901200 {
		t.Errorf("Expected id 901200, got %d", resultConfigs[0].Directives[0].(*types.SecAction).Id)
	}
	if resultConfigs[0].Directives[0].(*types.SecAction).Phase != "1" {
		t.Errorf("Expected phase 1, got %s", resultConfigs[0].Directives[0].(*types.SecAction).Phase)
	}
	if resultConfigs[0].Directives[0].(*types.SecAction).Tags[0] != "OWASP_CRS" {
		t.Errorf("Expected tag OWASP_CRS, got %s", resultConfigs[0].Directives[0].(*types.SecAction).Tags[0])
	}
	if resultConfigs[0].Directives[0].(*types.SecAction).Ver != "OWASP_CRS/4.0.1-dev" {
		t.Errorf("Expected version OWASP_CRS/4.0.1-dev, got %s", resultConfigs[0].Directives[0].(*types.SecAction).Ver)
	}
	if slices.Contains(resultConfigs[0].Directives[0].(*types.SecAction).GetActionKeys(), "nolog") == false {
		t.Errorf("Expected nolog action, not found")
	}
	if resultConfigs[0].Directives[0].(*types.SecAction).DisruptiveAction.Action != "pass" {
		t.Errorf("Expected disruptive action pass, got %s", resultConfigs[0].Directives[0].(*types.SecAction).DisruptiveAction.Action)
	} 
}

func TestLoadSecRule(t *testing.T) {
	testPayload := `
#
# Validate request line against the format specified in the HTTP RFC
#
# -=[ Rule Logic ]=-
#
# Uses rule negation against the regex for positive security.   The regex specifies the proper
# construction of URI request lines such as:
#
#   "http" "://" authority path-abempty [ "?" query ]
#
# It also outlines proper construction for CONNECT, OPTIONS and GET requests.
#
# Regular expression generated from regex-assembly/920100.ra.
# To update the regular expression run the following shell script
# (consult https://coreruleset.org/docs/development/regex_assembly/ for details):
#   crs-toolchain regex update 920100
#
# -=[ References ]=-
# https://www.rfc-editor.org/rfc/rfc9110.html#section-4.2.1
# http://capec.mitre.org/data/definitions/272.html
#
SecRule REQUEST_LINE "!@rx (?i)^(?:get /[^#\?]*(?:\?[^\s\v#]*)?(?:#[^\s\v]*)?|(?:connect (?:(?:[0-9]{1,3}\.){3}[0-9]{1,3}\.?(?::[0-9]+)?|[\--9A-Z_a-z]+:[0-9]+)|options \*|[a-z]{3,10}[\s\v]+(?:[0-9A-Z_a-z]{3,7}?://[\--9A-Z_a-z]*(?::[0-9]+)?)?/[^#\?]*(?:\?[^\s\v#]*)?(?:#[^\s\v]*)?)[\s\v]+[\.-9A-Z_a-z]+)$" \
    "id:920100,\
    phase:1,\
    block,\
    t:none,\
    msg:'Invalid HTTP Request Line',\
    logdata:'%{request_line}',\
    tag:'application-multi',\
    tag:'language-multi',\
    tag:'platform-multi',\
    tag:'attack-protocol',\
    tag:'paranoia-level/1',\
    tag:'OWASP_CRS',\
    tag:'capec/1000/210/272',\
    ver:'OWASP_CRS/4.0.1-dev',\
    severity:'WARNING',\
    setvar:'tx.inbound_anomaly_score_pl1=+%{tx.warning_anomaly_score}'"`

	resultConfigs := []types.Configuration{}
	input := antlr.NewInputStream(testPayload)
	lexer := parsing.NewSecLangLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parsing.NewSecLangParser(stream)
	start := p.Configuration()
	var listener ExtendedSeclangParserListener
	antlr.ParseTreeWalkerDefault.Walk(&listener, start)
	resultConfigs = append(resultConfigs, listener.ConfigurationList.Configurations...)

	if len(resultConfigs) != 1 {
		t.Errorf("Expected 1 configuration, got %d", len(resultConfigs))
	}
	if len(resultConfigs[0].Directives) != 1 {
		t.Errorf("Expected 1 directive, got %d", len(resultConfigs[0].Directives))
	}
	if resultConfigs[0].Directives[0].(*types.SecRule).Id != 920100 {
		t.Errorf("Expected id 920100, got %d", resultConfigs[0].Directives[0].(*types.SecRule).Id)
	}
	if resultConfigs[0].Directives[0].(*types.SecRule).Phase != "1" {
		t.Errorf("Expected phase 1, got %s", resultConfigs[0].Directives[0].(*types.SecRule).Phase)
	}
	if resultConfigs[0].Directives[0].(*types.SecRule).Tags[0] != "application-multi" {
		t.Errorf("Expected tag application-multi, got %s", resultConfigs[0].Directives[0].(*types.SecRule).Tags[0])
	}
	if resultConfigs[0].Directives[0].(*types.SecRule).Ver != "OWASP_CRS/4.0.1-dev" {
		t.Errorf("Expected version OWASP_CRS/4.0.1-dev, got %s", resultConfigs[0].Directives[0].(*types.SecRule).Ver)
	}
	if resultConfigs[0].Directives[0].(*types.SecRule).Msg != "Invalid HTTP Request Line" {
		t.Errorf("Expected message Invalid HTTP Request Line, got %s", resultConfigs[0].Directives[0].(*types.SecRule).Msg)
	}
	if resultConfigs[0].Directives[0].(*types.SecRule).Severity != "WARNING" {
		t.Errorf("Expected severity WARNING, got %s", resultConfigs[0].Directives[0].(*types.SecRule).Severity)
	}
	if slices.Contains(resultConfigs[0].Directives[0].(*types.SecRule).GetActionKeys(), "logdata") == false {
		t.Errorf("Expected logdata action, not found")
	}
	if resultConfigs[0].Directives[0].(*types.SecRule).DisruptiveAction.Action != "block" {
		t.Errorf("Expected disruptive action block, got %s", resultConfigs[0].Directives[0].(*types.SecRule).DisruptiveAction.Action)
	}
}