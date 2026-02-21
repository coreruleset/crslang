package translator

import (
	"io"
	"os"

	"go.yaml.in/yaml/v4"
)

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

// PrintYAML marshal and write structures to a yaml file
func PrintYAML(input any, filename string) error {
	yamlFile, err := yaml.Marshal(input)
	if err != nil {
		return err
	}

	err = writeToFile(yamlFile, filename)

	return err
}
