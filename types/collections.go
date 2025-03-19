package types

import "fmt"

type Collection struct {
	Name      CollectionName `yaml:"name,omitempty"`
	Arguments []string       `yaml:"arguments,omitempty"`
	Excluded  []string       `yaml:"excludeds,omitempty"`
	Count     bool           `yaml:"count,omitempty"`
}

type CollectionName string

const (
	// Collections
	ARGS                   CollectionName = "ARGS"
	ARGS_GET               CollectionName = "ARGS_GET"
	ARGS_GET_NAMES         CollectionName = "ARGS_GET_NAMES"
	ARGS_NAMES             CollectionName = "ARGS_NAMES"
	ARGS_POST_NAMES        CollectionName = "ARGS_POST_NAMES"
	ARGS_POST              CollectionName = "ARGS_POST"
	ENV                    CollectionName = "ENV"
	FILES                  CollectionName = "FILES"
	GEO                    CollectionName = "GEO"
	GLOBAL                 CollectionName = "GLOBAL"
	IP                     CollectionName = "IP"
	MATCHED_VARS_NAMES     CollectionName = "MATCHED_VARS_NAMES"
	MATCHED_VARS           CollectionName = "MATCHED_VARS"
	MULTIPART_PART_HEADERS CollectionName = "MULTIPART_PART_HEADERS"
	PERF_RULES             CollectionName = "PERF_RULES"
	REQUEST_COOKIES_NAMES  CollectionName = "REQUEST_COOKIES_NAMES"
	REQUEST_COOKIES        CollectionName = "REQUEST_COOKIES"
	REQUEST_HEADERS_NAMES  CollectionName = "REQUEST_HEADERS_NAMES"
	REQUEST_HEADERS        CollectionName = "REQUEST_HEADERS"
	RESPONSE_HEADERS_NAMES CollectionName = "RESPONSE_HEADERS_NAMES"
	RESPONSE_HEADERS       CollectionName = "RESPONSE_HEADERS"
	RULE                   CollectionName = "RULE"
	SESSION                CollectionName = "SESSION"
	TX                     CollectionName = "TX"
	XML                    CollectionName = "XML"
)

var (
	allCollections = map[string]CollectionName{
		"ARGS":                   ARGS,
		"ARGS_GET":               ARGS_GET,
		"ARGS_GET_NAMES":         ARGS_GET_NAMES,
		"ARGS_NAMES":             ARGS_NAMES,
		"ARGS_POST_NAMES":        ARGS_POST_NAMES,
		"ARGS_POST":              ARGS_POST,
		"ENV":                    ENV,
		"FILES":                  FILES,
		"GEO":                    GEO,
		"GLOBAL":                 GLOBAL,
		"IP":                     IP,
		"MATCHED_VARS_NAMES":     MATCHED_VARS_NAMES,
		"MATCHED_VARS":           MATCHED_VARS,
		"MULTIPART_PART_HEADERS": MULTIPART_PART_HEADERS,
		"PERF_RULES":             PERF_RULES,
		"REQUEST_COOKIES_NAMES":  REQUEST_COOKIES_NAMES,
		"REQUEST_COOKIES":        REQUEST_COOKIES,
		"REQUEST_HEADERS_NAMES":  REQUEST_HEADERS_NAMES,
		"REQUEST_HEADERS":        REQUEST_HEADERS,
		"RESPONSE_HEADERS_NAMES": RESPONSE_HEADERS_NAMES,
		"RESPONSE_HEADERS":       RESPONSE_HEADERS,
		"RULE":                   RULE,
		"SESSION":                SESSION,
		"TX":                     TX,
		"XML":                    XML,
	}
)

func CollectionsToString(collections []Collection, separator string) string {
	result := ""
	for i, collection := range collections {
		if len(collection.Arguments) == 0 && len(collection.Excluded) == 0 {
			if collection.Count {
				result += "&"
			}
			result += string(collection.Name)
		} else {
			for _, arg := range collection.Arguments {
				if collection.Count {
					result += "&"
				}
				result += string(collection.Name) + ":" + arg
				if i != len(collection.Arguments)-1 || len(collection.Excluded) > 0 {
					result += separator
				}
			}
			for _, excluded := range collection.Excluded {
				result += "!" + string(collection.Name) + ":" + excluded
				if i != len(collection.Excluded)-1 {
					result += separator
				}
			}
		}
		if i != len(collections)-1 {
			result += separator
		}
	}
	return result
}

func GetCollection(name string) (CollectionName, error) {
	col, exists := allCollections[name]
	if !exists {
		return "", fmt.Errorf("Invalid collection name: %s", name)
	}
	return col, nil
}
