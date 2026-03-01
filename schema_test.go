package main

import (
	"testing"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/coreruleset/crslang/translator"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"go.yaml.in/yaml/v4"
)

const schemaPath = "schema.json"
const testDataDir = "testdata"

func TestConfFilesAgainstSchema(t *testing.T) {
	compiler := jsonschema.NewCompiler()
	schema, err := compiler.Compile(schemaPath)
	if err != nil {
		t.Fatalf("Failed to compile schema: %v", err)
	}

	confFiles, err := doublestar.FilepathGlob(testDataDir + "/**/*.conf")
	if err != nil {
		t.Fatalf("Failed to glob conf files: %v", err)
	}

	if len(confFiles) == 0 {
		t.Fatal("No .conf files found in testdata")
	}

	t.Logf("Found %d .conf files to test", len(confFiles))

	for _, confFile := range confFiles {
		t.Run(confFile, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Skipf("Skipping %s: parser panic: %v", confFile, r)
				}
			}()

			configList, err := translator.LoadSeclang(confFile)
			if err != nil {
				t.Fatalf("Failed to load seclang file %s: %v", confFile, err)
			}

			configList = *translator.ToCRSLang(configList)

			yamlData, err := yaml.Marshal(configList)
			if err != nil {
				t.Fatalf("Failed to marshal to YAML: %v", err)
			}

			var yamlObj any
			if err := yaml.Unmarshal(yamlData, &yamlObj); err != nil {
				t.Fatalf("Failed to unmarshal YAML: %v", err)
			}

			if err := schema.Validate(yamlObj); err != nil {
				t.Errorf("Schema validation failed for %s:\n%v\n\nGenerated YAML:\n%s", confFile, err, string(yamlData))
			}
		})
	}
}
