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

// assignDirectiveIDs assigns a base id (and an indexed suffix when there are
// multiple directive lists) to each entry produced by a single parse run.
func assignDirectiveIDs(directives []types.DirectiveList, id string) {
	for i := range directives {
		directives[i].Id = id
		if len(directives) > 1 {
			directives[i].Id += "_" + strconv.Itoa(i+1)
		}
	}
}

// LoadSeclangFromString loads seclang directives from a string and returns a ConfigurationList.
// The id parameter is used to name the resulting directive list.
func LoadSeclangFromString(content string, id string) (types.ConfigurationList, error) {
	input := antlr.NewInputStream(content)
	lexer := parser.NewSecLangLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewSecLangParser(stream)
	start := p.Configuration()
	var seclangListener listener.ExtendedSeclangParserListener
	antlr.ParseTreeWalkerDefault.Walk(&seclangListener, start)
	assignDirectiveIDs(seclangListener.ConfigurationList.DirectiveList, id)
	return seclangListener.ConfigurationList, nil
}

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
			id := strings.TrimSuffix(filepath.Base(info.Name()), filepath.Ext(info.Name()))
			assignDirectiveIDs(seclangListener.ConfigurationList.DirectiveList, id)
			resultConfigs = append(resultConfigs, seclangListener.ConfigurationList.DirectiveList...)
		}
		return nil
	})
	configList := types.ConfigurationList{DirectiveList: resultConfigs}
	return configList, nil
}

// PrintSeclang writes seclang directives to files specified in directive list ids.
func PrintSeclang(configList types.ConfigurationList, dir string) error {
	dir = filepath.Clean(dir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	unfDirs := types.FromCRSLangToUnformattedDirectives(configList)

	for _, group := range unfDirs.DirectiveList {
		seclangDirectives := group.ToSeclang()
		groupId := group.Id + ".conf"
		if strings.HasSuffix(group.Id, ".conf") {
			groupId = group.Id + ".conf.example"
		}
		err := writeToFile([]byte(seclangDirectives), filepath.Join(dir, groupId))
		if err != nil {
			return err
		}
	}

	return nil
}
