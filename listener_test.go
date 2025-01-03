package main

import (
	"slices"
	"testing"

	"github.com/antlr4-go/antlr/v4"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/parsing"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/types"
)

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