package types

import (
	"reflect"
	"slices"
	"strconv"
	"strings"
)

var (
	metadataFields = []string{"Msg", "Ver", "Severity"}
	/* TODO: change action names to defined constants */
	fieldOrder = []string{"Id", "Phase", "allow", "block", "deny", "drop", "pass", "pause", "proxy", "redirect", "status", "capture", "Transformations", "log", "nolog", "auditlog", "noauditlog", "Msg", "logdata", "Tags", "sanitiseArg", "sanitiseRequestHeader", "sanitiseMatched", "sanitiseMatchedBytes", "ctl", "Ver", "Severity", "multiMatch", "initcol", "setenv", "setvar", "expirevar", "chain", "skip", "skipAfter"}
)

func sortActions(d ChainableDirective) []string {
	m := *d.GetMetadata().(*SecRuleMetadata)
	rM := reflect.ValueOf(m)
	a := *d.GetActions()
	aKeys := a.GetActionKeys()
	results := []string{}

	for _, fn := range fieldOrder {
		switch {
		case fn == "Id" && m.Id != 0:
			results = append(results, "id:"+strconv.Itoa(m.Id))
		case fn == "Phase" && m.Phase != "":
			results = append(results, "phase:"+m.Phase)
		case fn == "Tags" && len(m.Tags) > 0:
			for _, t := range m.Tags {
				results = append(results, "tag:'"+t+"'")
			}
		case slices.Contains(metadataFields, fn) && rM.FieldByName(fn).IsValid() && rM.FieldByName(fn).String() != "":
			results = append(results, strings.ToLower(fn)+":'"+rM.FieldByName(fn).String()+"'")
		case fn == "Transformations" && len(d.GetTransformations().Transformations) > 0:
			results = append(results, d.GetTransformations().ToString())
		case (fn == "setvar" || fn == "ctl") && slices.Contains(aKeys, fn):
			aList := a.GetActionsByKey(fn)
			for _, action := range aList {
				results = append(results, action.ToString())
			}
		case slices.Contains(aKeys, fn):
			results = append(results, a.GetActionByKey(fn).ToString())
		}
	}

	return results
}
