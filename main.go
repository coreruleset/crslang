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
		resultConfigs := []types.DirectiveList{}
		filepath.Walk(pathArg, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Ext(info.Name()) == ".conf" {
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
				for i := range seclangListener.ConfigurationList.DirectiveList {
					seclangListener.ConfigurationList.DirectiveList[i].Id = strings.TrimSuffix(filepath.Base(info.Name()), filepath.Ext(info.Name()))
					if len(seclangListener.ConfigurationList.DirectiveList) > 1 {
						seclangListener.ConfigurationList.DirectiveList[i].Id += "_" + strconv.Itoa(i+1)
					}
				}
				resultConfigs = append(resultConfigs, seclangListener.ConfigurationList.DirectiveList...)
			}
			return nil
		})
		configList := types.ConfigurationList{DirectiveList: resultConfigs}

		if *output == "" {
			*output = "crslang"
		}
		err := printCRSLang(configList, *output+".yaml")
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		if filepath.Ext(pathArg) != ".yaml" {
			log.Fatal("Only .yaml files are allowed")
		}

		configList := types.LoadDirectivesWithConditionsFromFile(pathArg)

		err := printSeclang(configList, *output)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

// printSeclang writes seclang directives to files specified in directive list ids.
func printSeclang(configList types.ConfigurationList, dir string) error {
	unfDirs := types.FromCRSLangToUnformattedDirectives(configList)

	for _, dirList := range unfDirs.DirectiveList {
		f, err := os.Create(dir + dirList.Id + ".conf")
		if err != nil {
			return err
		}
		seclangDirectives := dirList.ToSeclang()

		_, err = io.WriteString(f, seclangDirectives)
		if err != nil {
			return err
		}
	}

	return nil
}

// printSeclangToFile writes seclang format directives to a file
func printSeclangToFile(configList types.ConfigurationList, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	seclangDirectives := types.ToSeclang(configList)

	_, err = io.WriteString(f, seclangDirectives)
	if err != nil {
		return err
	}

	return nil
}

// printCRSLang writes crslang format directives (directives with conditions) to a file
func printCRSLang(configList types.ConfigurationList, filename string) error {
	configListWithConditions := types.ToDirectiveWithConditions(configList)

	yamlFile, err := yaml.Marshal(configListWithConditions.DirectiveList)
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
