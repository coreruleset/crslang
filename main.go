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
	"gitlab.fing.edu.uy/gsi/seclang/crslang/listener"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/parsing"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/types"
	"gopkg.in/yaml.v3"
)

var (
	progName = filepath.Base(os.Args[0])
)

func main() {
	toSeclang := flag.Bool("s", false, "Transalates the specified CRSLang file to Seclang files, instead of the default Seclang to CRSLang.")
	output := flag.String("o", "crslang", "Output file name used in translation from Seclang to CRSLang.")

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
				lexer := parsing.NewSecLangLexer(input)
				stream := antlr.NewCommonTokenStream(lexer, 0)
				p := parsing.NewSecLangParser(stream)
				start := p.Configuration()
				var seclangListener listener.ExtendedSeclangParserListener
				antlr.ParseTreeWalkerDefault.Walk(&seclangListener, start)
				for i := range seclangListener.ConfigurationList.DirectiveList {
					seclangListener.ConfigurationList.DirectiveList[i].Id = strings.TrimSuffix(filepath.Base(info.Name()), filepath.Ext(info.Name()))
					if i > 0 {
						seclangListener.ConfigurationList.DirectiveList[i].Id += "_" + strconv.Itoa(i+1)
					}
				}
				resultConfigs = append(resultConfigs, seclangListener.ConfigurationList.DirectiveList...)
			}
			return nil
		})
		configList := types.ConfigurationList{DirectiveList: resultConfigs}

		err := printCRSLang(configList, *output+".yaml")
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

// printSeclang writes seclang format directives to a file
func printSeclang(configList types.ConfigurationList, filename string) error {
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
