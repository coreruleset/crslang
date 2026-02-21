package translator

import (
	"os"
	"testing"

	"github.com/coreruleset/crslang/types"
	"go.yaml.in/yaml/v4"
)

var testFilepath = "../testdata/crs/"

func TestLoadCRS(t *testing.T) {

	configList, err := LoadSeclang(testFilepath)
	if err != nil {
		t.Errorf("Error loading seclang directives: %v", err)
	}

	configList = *ToCRSLang(configList)

	yamlFile, err := yaml.Marshal(configList)
	if err != nil {
		t.Errorf("Error marshalling yaml: %v", err)
	}

	err = writeToFile(yamlFile, "tmp_crslang.yaml")

	defer os.Remove("tmp_crslang.yaml")

	loadedConfigList := types.LoadDirectivesWithConditionsFromFile("tmp_crslang.yaml")
	yamlLoadedFile, err := yaml.Marshal(loadedConfigList)
	if err != nil {
		t.Errorf("Error writing file: %v", err)
	}

	if string(yamlFile) != string(yamlLoadedFile) {
		t.Errorf("Error: loaded file is different from original. Expected string length: %v, got: %v", len(string(yamlFile)), len(string(yamlLoadedFile)))
	}
}

func TestFromCRSLangToSeclang(t *testing.T) {
	configList, err := LoadSeclang(testFilepath)
	if err != nil {
		t.Errorf("Error loading seclang directives: %v", err)
	}

	seclangDirectives := types.ToSeclang(configList)

	crslangConfigList := ToCRSLang(configList)
	unformattedDirectives := types.FromCRSLangToUnformattedDirectives(*crslangConfigList)
	seclangDirectivesFromConditions := types.ToSeclang(*unformattedDirectives)

	if len(seclangDirectives) != len(seclangDirectivesFromConditions) {
		t.Errorf("Error in CRSLang to Seclang directives conversion. Expected length: %v, got: %v", len(seclangDirectives), len(seclangDirectivesFromConditions))
	}

}
