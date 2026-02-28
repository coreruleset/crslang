package translator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/coreruleset/crslang/types"
	"go.yaml.in/yaml/v4"
)

// ToCRSLang process previously loaded seclang directives to CRSLang schema directives
func ToCRSLang(configList types.Ruleset) *types.Ruleset {
	configListWithConditions := types.ToDirectiveWithConditions(configList)

	configListWithConditions.ExtractDefaultValues()

	for i := range configListWithConditions.Groups {
		configListWithConditions.Groups[i].ExtractDefaultValues()
	}

	return configListWithConditions
}

// WriteRuleSeparately writes each rule in a separate file, and creates a group.yaml file for each group and a ruleset.yaml file for the whole ruleset. The output is organized in the following structure:
// output/
// ├── ruleset.yaml
// ├── group1/
// │   ├── group.yaml
// │   └── rules/
// │       ├── 1.yaml
// │       ├── 2.yaml
// │       └── ...
func WriteRuleSeparately(ruleset types.Ruleset, output string) error {
	output = filepath.Clean(output)
	if err := os.MkdirAll(output, 0755); err != nil {
		return err
	}

	groups := []string{}

	// EXPERIMENTAL: output each group and rule in separate files
	for _, group := range ruleset.Groups {
		groups = append(groups, group.Id)

		groupFolder := filepath.Join(output, group.Id)
		ruleFolder := filepath.Join(groupFolder, "rules")
		err := os.MkdirAll(ruleFolder, 0755)
		if err != nil {
			return err
		}

		ruleIds := []string{}
		comments := []string{}
		configs := []types.ConfigurationDirective{}
		for _, directive := range group.Directives {
			if directive.GetKind() == types.RuleKind {
				rule, ok := directive.(*types.RuleWithCondition)
				if !ok {
					return fmt.Errorf("error casting directive of type %T in group %s to *types.RuleWithCondition", directive, group.Id)
				}

				fileName := filepath.Join(ruleFolder, strconv.Itoa(rule.Metadata.Id)+".yaml")
				err := PrintYAML(directive, fileName)
				if err != nil {
					return err
				}
				ruleIds = append(ruleIds, strconv.Itoa(rule.Metadata.Id))

			} else if directive.GetKind() == types.CommentKind {
				comment, ok := directive.(types.CommentDirective)
				if !ok {
					return fmt.Errorf("error casting directive of type %T in group %s to types.CommentDirective", directive, group.Id)
				}
				comments = append(comments, comment.Metadata.Comment)
			} else if directive.GetKind() == types.ConfigurationKind {
				config, ok := directive.(types.ConfigurationDirective)
				if !ok {
					return fmt.Errorf("error casting directive of type %T in group %s to types.ConfigurationDirective", directive, group.Id)
				}
				configs = append(configs, config)
			}
		}
		newGroup := types.Group{
			Id:             group.Id,
			Tags:           group.Tags,
			Comments:       comments,
			Rules:          ruleIds,
			Configurations: configs,
			Marker:         group.Marker,
		}
		err = PrintYAML(newGroup, filepath.Join(groupFolder, "group.yaml"))
		if err != nil {
			return err
		}
	}

	newRuleset := types.Ruleset{
		Global:    ruleset.Global,
		GroupsIds: groups,
	}
	err := PrintYAML(newRuleset, filepath.Join(output, "ruleset.yaml"))
	if err != nil {
		return err
	}
	return nil
}

// LoadRulesFromDirectory loads a ruleset from a directory containing a ruleset.yaml file and group subdirectories with group.yaml and rule yaml files. The structure of the directory should be as follows:
// dir/
// ├── ruleset.yaml
// ├── group1/
// │   ├── group.yaml
// │   └── rules/
// │       ├── 1.yaml
// │       ├── 2.yaml
// │       └── ...
func LoadRulesFromDirectory(dir string) (types.Ruleset, error) {
	info, err := os.Stat(dir)

	if err != nil {
		return types.Ruleset{}, err
	} else if !info.IsDir() {
		return types.Ruleset{}, fmt.Errorf("path is not a directory: %s", dir)
	}
	dir = filepath.Clean(dir)

	rFile, err := os.ReadFile(filepath.Join(dir, "ruleset.yaml"))

	if err != nil {
		return types.Ruleset{}, err
	}

	ruleset := types.Ruleset{}
	err = yaml.Unmarshal([]byte(rFile), &ruleset)

	if err != nil {
		return types.Ruleset{}, err
	}

	for _, groupId := range ruleset.GroupsIds {
		if err := verifyId(groupId); err != nil {
			return types.Ruleset{}, err
		}
		groupFile, err := os.ReadFile(filepath.Join(dir, groupId, "group.yaml"))
		if err != nil {
			return types.Ruleset{}, err
		}
		group := types.Group{}
		err = yaml.Unmarshal([]byte(groupFile), &group)
		if err != nil {
			return types.Ruleset{}, err
		}
		for _, comment := range group.Comments {
			group.Directives = append(group.Directives, types.CommentDirective{
				Metadata: types.CommentMetadata{
					Comment: comment,
				},
			})
		}
		for _, config := range group.Configurations {
			group.Directives = append(group.Directives, config)
		}
		for _, ruleId := range group.Rules {
			ruleFile, err := os.ReadFile(filepath.Join(dir, groupId, "rules", ruleId+".yaml"))
			if err != nil {
				return types.Ruleset{}, err
			}
			rule := types.RuleWithCondition{}
			err = yaml.Unmarshal([]byte(ruleFile), &rule)
			if err != nil {
				return types.Ruleset{}, err
			}
			group.Directives = append(group.Directives, &rule)
		}
		group.Rules = nil
		ruleset.Groups = append(ruleset.Groups, group)
	}
	ruleset.GroupsIds = nil
	return ruleset, nil
}

var validID = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

func verifyId(id string) error {
	if !validID.MatchString(id) {
		return fmt.Errorf("invalid id: %s. Ids can only contain letters, numbers, underscores and hyphens", id)
	}
	return nil
}
