package types

import "fmt"

type SecAction struct {
	Metadata        *SecRuleMetadata `yaml:"metadata,omitempty"`
	Transformations `yaml:",inline"`
	Actions         *SeclangActions    `yaml:"actions"`
	ChainedRule     ChainableDirective `yaml:"chainedRule,omitempty"`
}

func NewSecAction() *SecAction {
	secAction := new(SecAction)
	secAction.Metadata = new(SecRuleMetadata)
	secAction.Actions = new(SeclangActions)
	return secAction
}

func (d SecAction) GetMetadata() Metadata {
	return d.Metadata
}

func (d SecAction) GetActions() *SeclangActions {
	return d.Actions
}

func (s *SecAction) AppendChainedDirective(chainedDirective ChainableDirective) {
	s.ChainedRule = chainedDirective
}

func (s SecAction) GetChainedDirective() ChainableDirective {
	return s.ChainedRule
}

func (s SecAction) NonDisruptiveActionsCount() int {
	return len(s.Actions.NonDisruptiveActions)
}

func (s SecAction) ToSeclang() string {
	result := ""
	result += s.Metadata.Comment + "SecAction \"phase:" + s.Metadata.Phase
	actions := s.Actions.ToString()
	transformations := s.Transformations.ToString()
	if actions != "" {
		result += "," + actions
	}
	if transformations != "" {
		result += ", " + transformations
	}
	result += "\"\n"
	return result
}

func (s SecAction) ToSeclangWithIdent(initialString string) string {
	return initialString + s.ToSeclang()
}

func (s SecAction) Equal(s2 SecAction) error {
	err := s.Metadata.Equal(*s2.Metadata)
	if err != nil {
		return err
	}
	err = s.Transformations.Equal(s2.Transformations)
	if err != nil {
		return err
	}
	err = s.Actions.Equal(*s2.Actions)
	if err != nil {
		return err
	}
	if s.ChainedRule == nil && s2.ChainedRule != nil || s.ChainedRule != nil && s2.ChainedRule == nil {
		return fmt.Errorf("Expected chained rule: %s, got: %s", s.ChainedRule, s2.ChainedRule)
	}
	if c, ok := s.ChainedRule.(*SecRule); ok {
		if c2, ok := s2.ChainedRule.(*SecRule); ok {
			err = c.Equal(*c2)
		}
	} else if c, ok := s.ChainedRule.(*SecAction); ok {
		if c2, ok := s2.ChainedRule.(*SecAction); ok {
			err = c.Equal(*c2)
		}
	} else if c, ok := s.ChainedRule.(*SecRuleScript); ok {
		if c2, ok := s2.ChainedRule.(*SecRuleScript); ok {
			err = c.Equal(*c2)
		}
	}
	if err != nil {
		return err
	}
	return nil
}
