package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/coreruleset/crslang/translator"
	"github.com/coreruleset/crslang/types"
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
		configList, err := translator.LoadSeclang(pathArg)
		if err != nil {
			log.Fatal(err.Error())
		}

		configList = *translator.ToCRSLang(configList)
		if !*dirMode {
			if *output == "" {
				*output = "crslang"
			}

			err = translator.PrintYAML(configList, *output+".yaml")
			if err != nil {
				log.Fatal(err.Error())
			}
		} else {
			err := translator.WriteRuleSeparately(configList, *output)
			if err != nil {
				log.Fatal(err.Error())
			}
		}

	} else {
		if !*dirMode {
			if filepath.Ext(pathArg) != ".yaml" {
				log.Fatal("Only .yaml files are allowed")
			}

			configList := types.LoadDirectivesWithConditionsFromFile(pathArg)
			err := translator.PrintSeclang(configList, *output)
			if err != nil {
				log.Fatal(err.Error())
			}
		} else {
			/* Load rule from dir */
			configList, err := translator.LoadRulesFromDirectory(pathArg)
			if err != nil {
				log.Fatal(err.Error())
			}
			err = translator.PrintSeclang(configList, *output)
			if err != nil {
				log.Fatal(err.Error())
			}
		}
	}
}
