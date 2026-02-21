//go:build js && wasm

// Package main provides a WebAssembly (WASM) interface for the crslang translation
// functions, exposing them as callable JavaScript functions in the browser.
package main

import (
	"syscall/js"

	"github.com/coreruleset/crslang/translator"
	"github.com/coreruleset/crslang/types"
	"go.yaml.in/yaml/v4"
)

// seclangToCRSLang is a JavaScript-callable function that translates a seclang
// configuration string into a CRSLang YAML string.
//
// JavaScript usage:
//
//	const result = seclangToCRSLang(seclangContent);
//	if (result.error) { console.error(result.error); }
//	else { console.log(result.yaml); }
func seclangToCRSLang(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return js.ValueOf(map[string]any{"error": "seclangToCRSLang requires one argument: the seclang content string"})
	}
	content := args[0].String()

	configList, err := translator.LoadSeclangFromString(content, "input")
	if err != nil {
		return js.ValueOf(map[string]any{"error": err.Error()})
	}

	crslangList := translator.ToCRSLang(configList)

	yamlBytes, err := yaml.Marshal(crslangList)
	if err != nil {
		return js.ValueOf(map[string]any{"error": err.Error()})
	}

	return js.ValueOf(map[string]any{"yaml": string(yamlBytes)})
}

// crslangToSeclang is a JavaScript-callable function that translates a CRSLang
// YAML string back into seclang format.
//
// JavaScript usage:
//
//	const result = crslangToSeclang(crslangYaml);
//	if (result.error) { console.error(result.error); }
//	else { console.log(result.seclang); }
func crslangToSeclang(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return js.ValueOf(map[string]any{"error": "crslangToSeclang requires one argument: the CRSLang YAML string"})
	}
	content := args[0].String()

	configList := types.LoadDirectivesWithConditions([]byte(content))
	unfDirs := types.FromCRSLangToUnformattedDirectives(configList)
	seclang := types.ToSeclang(*unfDirs)

	return js.ValueOf(map[string]any{"seclang": seclang})
}

func main() {
	js.Global().Set("seclangToCRSLang", js.FuncOf(seclangToCRSLang))
	js.Global().Set("crslangToSeclang", js.FuncOf(crslangToSeclang))

	// Keep the Go runtime alive until the page is closed.
	select {}
}
