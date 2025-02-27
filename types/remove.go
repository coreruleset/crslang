package types

import (
	"strconv"
)

type RemoveRuleDirective struct {
	Kind     Kind            `yaml:"kind"`
	Metadata CommentMetadata `yaml:"metadata,omitempty"`
	Ids      []int           `yaml:"ids,omitempty"`
	IdRanges []string        `yaml:"id_ranges,omitempty"`
	// IdRanges []IdRange       `yaml:"id_ranges,omitempty"`
	Tags []string `yaml:"tag,omitempty"`
	Msgs []string `yaml:"msg,omitempty"`
}

type IdRange struct {
	Start int `yaml:"start"`
	End   int `yaml:"end"`
}

func (r IdRange) ToString() string {
	start := strconv.Itoa(r.Start)
	end := strconv.Itoa(r.End)
	return start + "-" + end
}

// TODO: adjust
func (d RemoveRuleDirective) ToSeclang() string {
	result := ""
	if len(d.Ids) > 0 || len(d.IdRanges) > 0 {
		result += "SecRuleRemoveById"
		for _, id := range d.Ids {
			result += " " + strconv.Itoa(id)
		}
		for _, idRange := range d.IdRanges {
			/* 			ruleContent += " " + idRange.ToString() */
			result += " " + idRange
		}
		result += "\n"
	}
	if len(d.Tags) > 0 {
		for _, tag := range d.Tags {
			result += "SecRuleRemoveByTag " + tag + "\n"
		}
	}
	if len(d.Msgs) > 0 {
		for _, msg := range d.Msgs {
			result += "SecRuleRemoveByMsg " + msg + "\n"
		}
	}
	return result
}
