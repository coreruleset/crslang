package types

import (
	"strconv"
)

type RemoveRuleDirective struct {
	Kind     Kind            `yaml:"kind"`
	Metadata CommentMetadata `yaml:"metadata,omitempty"`
	Ids      []int           `yaml:"ids,omitempty"`
	IdRanges []IdRange       `yaml:"id_ranges,omitempty"`
	Tags     []string        `yaml:"tags,omitempty"`
	Msgs     []string        `yaml:"msgs,omitempty"`
}

func (d RemoveRuleDirective) GetKind() Kind {
	return d.Kind
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

func (d RemoveRuleDirective) ToSeclang() string {
	results := ""
	if len(d.Ids) > 0 || len(d.IdRanges) > 0 {
		results += "SecRuleRemoveById"
		for _, id := range d.Ids {
			results += " " + strconv.Itoa(id)
		}
		for _, idRange := range d.IdRanges {
			results += " " + idRange.ToString()
		}
		results += "\n"
	}
	if len(d.Tags) > 0 {
		for _, tag := range d.Tags {
			results += "SecRuleRemoveByTag " + tag + "\n"
		}
	}
	if len(d.Msgs) > 0 {
		for _, msg := range d.Msgs {
			results += "SecRuleRemoveByMsg " + msg + "\n"
		}
	}
	return results
}
