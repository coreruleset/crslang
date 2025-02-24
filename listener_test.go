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

	resultConfigs := []types.DirectiveList{}
	input := antlr.NewInputStream(testPayload)
	lexer := parsing.NewSecLangLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parsing.NewSecLangParser(stream)
	start := p.Configuration()
	var listener ExtendedSeclangParserListener
	antlr.ParseTreeWalkerDefault.Walk(&listener, start)
	resultConfigs = append(resultConfigs, listener.ConfigurationList.DirectiveList...)

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

func TestLoadConfigurationDirective(t *testing.T) {
	testPayload := `
SecRuleEngine On
`

	resultConfigs := []types.DirectiveList{}
	input := antlr.NewInputStream(testPayload)
	lexer := parsing.NewSecLangLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parsing.NewSecLangParser(stream)
	start := p.Configuration()
	var listener ExtendedSeclangParserListener
	antlr.ParseTreeWalkerDefault.Walk(&listener, start)
	resultConfigs = append(resultConfigs, listener.ConfigurationList.DirectiveList...)

	if len(resultConfigs) != 1 {
		t.Errorf("Expected 1 configuration, got %d", len(resultConfigs))
	}
	if len(resultConfigs[0].Directives) != 1 {
		t.Errorf("Expected 1 directive, got %d", len(resultConfigs[0].Directives))
	}
	configDirective, ok := resultConfigs[0].Directives[0].(types.ConfigurationDirective)
	if !ok {
		t.Errorf("Expected configuration directive, got %T", resultConfigs[0].Directives[0])
	}
	if configDirective.Name != "SecRuleEngine" {
		t.Errorf("Expected directive SecRuleEngine, got %s", configDirective.Name)
	}
	if configDirective.Parameter != "On" {
		t.Errorf("Expected parameter On, got %s", configDirective.Parameter)
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

	resultConfigs := []types.DirectiveList{}
	input := antlr.NewInputStream(testPayload)
	lexer := parsing.NewSecLangLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parsing.NewSecLangParser(stream)
	start := p.Configuration()
	var listener ExtendedSeclangParserListener
	antlr.ParseTreeWalkerDefault.Walk(&listener, start)
	resultConfigs = append(resultConfigs, listener.ConfigurationList.DirectiveList...)

	if len(resultConfigs) != 1 {
		t.Errorf("Expected 1 configuration, got %d", len(resultConfigs))
	}
	if len(resultConfigs[0].Directives) != 1 {
		t.Errorf("Expected 1 directive, got %d", len(resultConfigs[0].Directives))
	}
	secAction, ok := resultConfigs[0].Directives[0].(*types.SecAction)
	if !ok {
		t.Errorf("Expected SecAction, got %T", resultConfigs[0].Directives[0])
	}
	if secAction.Metadata.Id != 901200 {
		t.Errorf("Expected id 901200, got %d", secAction.Metadata.Id)
	}
	if secAction.Metadata.Phase != "1" {
		t.Errorf("Expected phase 1, got %s", secAction.Metadata.Phase)
	}
	if secAction.Metadata.Tags[0] != "OWASP_CRS" {
		t.Errorf("Expected tag OWASP_CRS, got %s", secAction.Metadata.Tags[0])
	}
	if secAction.Metadata.Ver != "OWASP_CRS/4.0.1-dev" {
		t.Errorf("Expected version OWASP_CRS/4.0.1-dev, got %s", secAction.Metadata.Ver)
	}
	if slices.Contains(secAction.Actions.GetActionKeys(), "nolog") == false {
		t.Errorf("Expected nolog action, not found")
	}
	if secAction.Actions.DisruptiveAction.Action != "pass" {
		t.Errorf("Expected disruptive action pass, got %s", secAction.Actions.DisruptiveAction.Action)
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

	resultConfigs := []types.DirectiveList{}
	input := antlr.NewInputStream(testPayload)
	lexer := parsing.NewSecLangLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parsing.NewSecLangParser(stream)
	start := p.Configuration()
	var listener ExtendedSeclangParserListener
	antlr.ParseTreeWalkerDefault.Walk(&listener, start)
	resultConfigs = append(resultConfigs, listener.ConfigurationList.DirectiveList...)

	if len(resultConfigs) != 1 {
		t.Errorf("Expected 1 configuration, got %d", len(resultConfigs))
	}
	if len(resultConfigs[0].Directives) != 1 {
		t.Errorf("Expected 1 directive, got %d", len(resultConfigs[0].Directives))
	}
	secRule, ok := resultConfigs[0].Directives[0].(*types.SecRule)
	if !ok {
		t.Errorf("Expected SecRule, got %T", resultConfigs[0].Directives[0])
	}
	if secRule.Metadata.Id != 920100 {
		t.Errorf("Expected id 920100, got %d", secRule.Metadata.Id)
	}
	if secRule.Metadata.Phase != "1" {
		t.Errorf("Expected phase 1, got %s", secRule.Metadata.Phase)
	}
	if secRule.Metadata.Tags[0] != "application-multi" {
		t.Errorf("Expected tag application-multi, got %s", secRule.Metadata.Tags[0])
	}
	if secRule.Metadata.Ver != "OWASP_CRS/4.0.1-dev" {
		t.Errorf("Expected version OWASP_CRS/4.0.1-dev, got %s", secRule.Metadata.Ver)
	}
	if secRule.Metadata.Msg != "Invalid HTTP Request Line" {
		t.Errorf("Expected message Invalid HTTP Request Line, got %s", secRule.Metadata.Msg)
	}
	if secRule.Metadata.Severity != "WARNING" {
		t.Errorf("Expected severity WARNING, got %s", secRule.Metadata.Severity)
	}
	if slices.Contains(secRule.Actions.GetActionKeys(), "logdata") == false {
		t.Errorf("Expected logdata action, not found")
	}
	if secRule.Actions.DisruptiveAction.Action != "block" {
		t.Errorf("Expected disruptive action block, got %s", secRule.Actions.DisruptiveAction.Action)
	}
	if len(secRule.Variables) != 1 {
		t.Errorf("Expected 1 variable, got %d", len(secRule.Variables))
	}
	if secRule.Variables[0] != "REQUEST_LINE" {
		t.Errorf("Expected variable REQUEST_LINE, got %s", secRule.Variables[0])
	}
}

func TestLoadSecRuleWithCollection(t *testing.T) {
	testPayload := `
#
# -=[ Exclusion rule for 942440 ]=-
#
# Prevent FPs against Facebook click identifier
#
SecRule ARGS_GET:fbclid "@rx [a-zA-Z0-9_-]{61,61}" \
    "id:942441,\
    phase:2,\
    pass,\
    t:none,t:urlDecodeUni,\
    nolog,\
    tag:'OWASP_CRS',\
    ctl:ruleRemoveTargetById=942440;ARGS:fbclid,\
    ver:'OWASP_CRS/4.0.1-dev'"`

	resultConfigs := []types.DirectiveList{}
	input := antlr.NewInputStream(testPayload)
	lexer := parsing.NewSecLangLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parsing.NewSecLangParser(stream)
	start := p.Configuration()
	var listener ExtendedSeclangParserListener
	antlr.ParseTreeWalkerDefault.Walk(&listener, start)
	resultConfigs = append(resultConfigs, listener.ConfigurationList.DirectiveList...)

	if len(resultConfigs) != 1 {
		t.Errorf("Expected 1 configuration, got %d", len(resultConfigs))
	}
	if len(resultConfigs[0].Directives) != 1 {
		t.Errorf("Expected 1 directive, got %d", len(resultConfigs[0].Directives))
	}
	secRule, ok := resultConfigs[0].Directives[0].(*types.SecRule)
	if !ok {
		t.Errorf("Expected SecRule, got %T", resultConfigs[0].Directives[0])
	}
	if secRule.Metadata.Id != 942441 {
		t.Errorf("Expected id 942441, got %d", secRule.Metadata.Id)
	}
	if secRule.Metadata.Phase != "2" {
		t.Errorf("Expected phase 2, got %s", secRule.Metadata.Phase)
	}
	if secRule.Metadata.Tags[0] != "OWASP_CRS" {
		t.Errorf("Expected tag OWASP_CRS, got %s", secRule.Metadata.Tags[0])
	}
	if secRule.Metadata.Ver != "OWASP_CRS/4.0.1-dev" {
		t.Errorf("Expected version OWASP_CRS/4.0.1-dev, got %s", secRule.Metadata.Ver)
	}
	if slices.Contains(secRule.Actions.GetActionKeys(), "nolog") == false {
		t.Errorf("Expected nolog action, not found")
	}
	if secRule.Actions.DisruptiveAction.Action != "pass" {
		t.Errorf("Expected disruptive action pass, got %s", secRule.Actions.DisruptiveAction.Action)
	}
	if len(secRule.Collections) != 1 {
		t.Errorf("Expected 1 collection, got %d", len(secRule.Collections))
	}
	if secRule.Collections[0].Name != "ARGS_GET" {
		t.Errorf("Expected variable ARGS_GET, got %s", secRule.Collections[0].Name)
	}
	if secRule.Collections[0].Argument != "fbclid" {
		t.Errorf("Expected argument fbclid, got %s", secRule.Collections[0].Argument)
	}
}

func TestLoadChain(t *testing.T) {
	testPayload := `# This file is used as an exception mechanism to remove common false positives
# that may be encountered.
#
# Exception for Apache SSL pinger
#
SecRule REQUEST_LINE "@streq GET /" \
    "id:905100,\
    phase:1,\
    pass,\
    t:none,\
    nolog,\
    tag:'application-multi',\
    tag:'language-multi',\
    tag:'platform-apache',\
    tag:'attack-generic',\
    tag:'OWASP_CRS',\
    ver:'OWASP_CRS/4.6.0-dev',\
    chain"
    SecRule REMOTE_ADDR "@ipMatch 127.0.0.1,::1" \
        "t:none,\
        ctl:ruleRemoveByTag=OWASP_CRS,\
        ctl:auditEngine=Off"`

	resultConfigs := []types.DirectiveList{}
	input := antlr.NewInputStream(testPayload)
	lexer := parsing.NewSecLangLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parsing.NewSecLangParser(stream)
	start := p.Configuration()
	var listener ExtendedSeclangParserListener
	antlr.ParseTreeWalkerDefault.Walk(&listener, start)
	resultConfigs = append(resultConfigs, listener.ConfigurationList.DirectiveList...)

	if len(resultConfigs) != 1 {
		t.Errorf("Expected 1 configuration, got %d", len(resultConfigs))
	}
	if len(resultConfigs[0].Directives) != 1 {
		t.Errorf("Expected 1 directive, got %d", len(resultConfigs[0].Directives))
	}
	secRule, ok := resultConfigs[0].Directives[0].(*types.SecRule)
	if !ok {
		t.Errorf("Expected SecRule, got %T", resultConfigs[0].Directives[0])
	}
	if secRule.Metadata.Id != 905100 {
		t.Errorf("Expected id 905100, got %d", secRule.Metadata.Id)
	}
	if secRule.Metadata.Phase != "1" {
		t.Errorf("Expected phase 1, got %s", secRule.Metadata.Phase)
	}
	if secRule.Metadata.Tags[0] != "application-multi" {
		t.Errorf("Expected tag application-multi, got %s", secRule.Metadata.Tags[0])
	}
	if secRule.Metadata.Ver != "OWASP_CRS/4.6.0-dev" {
		t.Errorf("Expected version OWASP_CRS/4.6.0-dev, got %s", secRule.Metadata.Ver)
	}
	if secRule.Actions.DisruptiveAction.Action != "pass" {
		t.Errorf("Expected disruptive action pass, got %s", secRule.Actions.DisruptiveAction.Action)
	}
	if slices.Contains(secRule.Actions.GetActionKeys(), "chain") == false {
		t.Errorf("Expected chain action, not found")
	}
	chainedRule, ok := secRule.ChainedRule.(*types.SecRule)
	if !ok {
		t.Errorf("Expected SecRule, got %T", secRule.ChainedRule)
	}
	if chainedRule.Operator.Name != "ipMatch" {
		t.Errorf("Expected operator ipMatch, got %s", chainedRule.Operator.Name)
	}
	if chainedRule.Operator.Value != "127.0.0.1,::1" {
		t.Errorf("Expected operator value 127.0.0.1,::1, got %s", chainedRule.Operator.Value)
	}
	if slices.Contains(chainedRule.Actions.GetActionKeys(), "ctl") == false {
		t.Errorf("Expected ctl action, not found")
	}
}
