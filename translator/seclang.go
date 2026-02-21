package translator

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/coreruleset/crslang/listener"
	"github.com/coreruleset/crslang/types"
	"github.com/coreruleset/seclang_parser/parser"
)

// LoadSeclang loads seclang directives from an input file or folder and returns a ConfigurationList
// if a folder is specified it loads all .conf files in the folder
func LoadSeclang(input string) (types.ConfigurationList, error) {
	resultConfigs := []types.DirectiveList{}
	filepath.Walk(input, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if ext := filepath.Ext(info.Name()); !info.IsDir() && (ext == ".conf" || (ext == ".example" && filepath.Ext(strings.TrimSuffix(info.Name(), ext)) == ".conf")) {
			input, err := antlr.NewFileStream(path)
			if err != nil {
				return err
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
	return configList, nil
}

// PrintSeclang writes seclang directives to files specified in directive list ids.
func PrintSeclang(configList types.ConfigurationList, dir string) error {
	unfDirs := types.FromCRSLangToUnformattedDirectives(configList)

	for _, dirList := range unfDirs.DirectiveList {
		seclangDirectives := dirList.ToSeclang()
		dirId := dirList.Id + ".conf"
		if strings.HasSuffix(dirId, ".conf") {
			dirId = dirList.Id + ".conf.example"
		}
		err := writeToFile([]byte(seclangDirectives), dir+dirId)
		if err != nil {
			return err
		}
	}

	return nil
}
