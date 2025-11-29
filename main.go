package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/coreruleset/crslang/listener"
	"github.com/coreruleset/crslang/types"
	"github.com/coreruleset/seclang_parser/parser"
	"go.yaml.in/yaml/v4"
)

var (
	progName = filepath.Base(os.Args[0])
)

func main() {
	toSeclang := flag.Bool("s", false, "Transalates the specified CRSLang file to Seclang files, instead of the default Seclang to CRSLang.")
	// Experimental flag
	dirMode := flag.Bool("d", false, "If set, the script output will be divided into multiple files when translating from Seclang to CRSLang.")
	output := flag.String("o", "", "Output file name used in translation from Seclang to CRSLang. Output folder used in translation from CRSLang to Seclang.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `usage:
	%s [flags] filepath
 
Flags:
`, progName)
		flag.PrintDefaults()
	}

	flag.Parse()
	args := flag.Args()
	var pathArg string
	switch len(args) {
	case 0:
		log.Fatal("filepath is required")
	case 1:
		pathArg = args[0]
	default:
		log.Fatal("Only filepath is allowed")
	}

	if !*toSeclang {
		configList := LoadSeclang(pathArg)

		configList = *ToCRSLang(configList)
		if !*dirMode {
			if *output == "" {
				*output = "crslang"
			}

			err := printYAML(configList, *output+".yaml")
			if err != nil {
				log.Fatal(err.Error())
			}
		} else {
			// EXPERIMENTAL: output each group and rule in separate files
			for _, dirList := range configList.Groups {
				groupFolder := *output + "/" + dirList.Id + "/"
				ruleFolder := groupFolder + "/rules/"
				err := os.MkdirAll(ruleFolder, os.ModePerm)
				if err != nil {
					log.Fatal(err.Error())
				}
				listWithoutRules := []types.SeclangDirective{}
				for _, directive := range dirList.Directives {
					if directive.GetKind() == types.RuleKind {
						rule, ok := directive.(*types.RuleWithCondition)
						lastDigits := rule.Metadata.Id % 1000
						if lastDigits/100 != 0 {
							if !ok {
								log.Fatal("Error casting to RuleDirective")
							}
							fileName := ruleFolder + strconv.Itoa(rule.Metadata.Id) + ".yaml"
							err := printYAML(directive, fileName)
							if err != nil {
								log.Fatal(err.Error())
							}
						} else {
							listWithoutRules = append(listWithoutRules, directive)
						}
					} else {
						listWithoutRules = append(listWithoutRules, directive)
					}
				}
				group := types.Group{
					Id:         dirList.Id,
					Tags:       dirList.Tags,
					Directives: listWithoutRules,
					Marker:     dirList.Marker,
				}
				err = printYAML(group, groupFolder+"group.yaml")
				if err != nil {
					log.Fatal(err.Error())
				}
			}
		}
	} else {
		if filepath.Ext(pathArg) != ".yaml" {
			log.Fatal("Only .yaml files are allowed")
		}

		configList := types.LoadDirectivesWithConditionsFromFile(pathArg)

		err := PrintSeclang(configList, *output)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

// LoadSeclang loads seclang directives from an input file or folder and returns a ConfigurationList
// if a folder is specified it loads all .conf files in the folder
func LoadSeclang(input string) types.Ruleset {
	resultConfigs := []types.Group{}
	filepath.Walk(input, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if ext := filepath.Ext(info.Name()); !info.IsDir() && (ext == ".conf" || (ext == ".example" && filepath.Ext(strings.TrimSuffix(info.Name(), ext)) == ".conf")) {
			input, err := antlr.NewFileStream(path)
			if err != nil {
				panic("Error reading file" + path)
			}
			lexer := parser.NewSecLangLexer(input)
			stream := antlr.NewCommonTokenStream(lexer, 0)
			p := parser.NewSecLangParser(stream)
			start := p.Configuration()
			var seclangListener listener.ExtendedSeclangParserListener
			antlr.ParseTreeWalkerDefault.Walk(&seclangListener, start)
			for i := range seclangListener.ConfigurationList.Groups {
				seclangListener.ConfigurationList.Groups[i].Id = strings.TrimSuffix(filepath.Base(info.Name()), filepath.Ext(info.Name()))
				if len(seclangListener.ConfigurationList.Groups) > 1 {
					seclangListener.ConfigurationList.Groups[i].Id += "_" + strconv.Itoa(i+1)
				}
			}
			resultConfigs = append(resultConfigs, seclangListener.ConfigurationList.Groups...)
		}
		return nil
	})
	configList := types.Ruleset{Groups: resultConfigs}
	return configList
}

// PrintSeclang writes seclang directives to files specified in directive list ids.
func PrintSeclang(configList types.Ruleset, dir string) error {
	unfDirs := types.FromCRSLangToUnformattedDirectives(configList)

	for _, dirList := range unfDirs.Groups {
		seclangDirectives := dirList.ToSeclang()
		err := writeToFile([]byte(seclangDirectives), dir+dirList.Id+".conf")
		if err != nil {
			return err
		}
	}

	return nil
}

// ToCRSLang process previously loaded seclang directives to CRSLang schema directives
func ToCRSLang(configList types.Ruleset) *types.Ruleset {
	configListWithConditions := types.ToDirectiveWithConditions(configList)

	configListWithConditions.ExtractDefaultValues()
	// EXPERIMENTAL: extract default values for each group
	for i := range configListWithConditions.Groups {
		configListWithConditions.Groups[i].ExtractDefaultValues()
	}
	return configListWithConditions
}

// printYAML marshal and write structures to a yaml file
func printYAML(input any, filename string) error {
	yamlFile, err := yaml.Marshal(input)
	if err != nil {
		return err
	}

	err = writeToFile(yamlFile, filename)

	return err
}

func writeToFile(payload []byte, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.WriteString(f, string(payload))
	if err != nil {
		return err
	}

	return nil
}
