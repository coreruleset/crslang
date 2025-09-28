package types

import "fmt"

type Collection struct {
	Name      CollectionName `yaml:"name,omitempty"`
	Arguments []string       `yaml:"arguments,omitempty"`
	Excluded  []string       `yaml:"excludeds,omitempty"`
	Count     bool           `yaml:"count,omitempty"`
}

type CollectionName int

const (
	// Collections
	UNKNOWN_COLLECTION CollectionName = iota
	ARGS
	ARGS_GET
	ARGS_GET_NAMES
	ARGS_NAMES
	ARGS_POST_NAMES
	ARGS_POST
	ENV
	FILES
	GEO
	GLOBAL
	IP
	MATCHED_VARS_NAMES
	MATCHED_VARS
	MULTIPART_PART_HEADERS
	PERF_RULES
	REQUEST_COOKIES_NAMES
	REQUEST_COOKIES
	REQUEST_HEADERS_NAMES
	REQUEST_HEADERS
	RESPONSE_HEADERS_NAMES
	RESPONSE_HEADERS
	RULE
	SESSION
	TX
	XML
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

func (c CollectionName) String() string {
	switch c {
	case ARGS:
		return "ARGS"
	case ARGS_GET:
		return "ARGS_GET"
	case ARGS_GET_NAMES:
		return "ARGS_GET_NAMES"
	case ARGS_NAMES:
		return "ARGS_NAMES"
	case ARGS_POST_NAMES:
		return "ARGS_POST_NAMES"
	case ARGS_POST:
		return "ARGS_POST"
	case ENV:
		return "ENV"
	case FILES:
		return "FILES"
	case GEO:
		return "GEO"
	case GLOBAL:
		return "GLOBAL"
	case IP:
		return "IP"
	case MATCHED_VARS_NAMES:
		return "MATCHED_VARS_NAMES"
	case MATCHED_VARS:
		return "MATCHED_VARS"
	case MULTIPART_PART_HEADERS:
		return "MULTIPART_PART_HEADERS"
	case PERF_RULES:
		return "PERF_RULES"
	case REQUEST_COOKIES_NAMES:
		return "REQUEST_COOKIES_NAMES"
	case REQUEST_COOKIES:
		return "REQUEST_COOKIES"
	case REQUEST_HEADERS_NAMES:
		return "REQUEST_HEADERS_NAMES"
	case REQUEST_HEADERS:
		return "REQUEST_HEADERS"
	case RESPONSE_HEADERS_NAMES:
		return "RESPONSE_HEADERS_NAMES"
	case RESPONSE_HEADERS:
		return "RESPONSE_HEADERS"
	case RULE:
		return "RULE"
	case SESSION:
		return "SESSION"
	case TX:
		return "TX"
	case XML:
		return "XML"
	default:
		return "unknown"
	}
}

func (c CollectionName) MarshalYAML() (interface{}, error) {
	if c == UNKNOWN_COLLECTION {
		return nil, fmt.Errorf("Unknown collection name")
	}
	return c.String(), nil
}

func (c *CollectionName) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var name string
	if err := unmarshal(&name); err != nil {
		return err
	}
	*c = stringToCollectionName(name)
	if *c == UNKNOWN_COLLECTION {
		return fmt.Errorf("Collection name %s is not valid", name)
	}
	return nil
}

func CollectionsToString(collections []Collection, separator string) string {
	result := ""
	for i, collection := range collections {
		if len(collection.Arguments) == 0 && len(collection.Excluded) == 0 {
			if collection.Count {
				result += "&"
			}
			result += collection.Name.String()
		} else {
			for j, arg := range collection.Arguments {
				if collection.Count {
					result += "&"
				}
				result += collection.Name.String() + ":" + arg
				if j != len(collection.Arguments)-1 || len(collection.Excluded) > 0 {
					result += separator
				}
			}
			for j, excluded := range collection.Excluded {
				result += "!" + collection.Name.String() + ":" + excluded
				if j != len(collection.Excluded)-1 {
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

func stringToCollectionName(name string) CollectionName {
	col, exists := allCollections[name]
	if !exists {
		return UNKNOWN_COLLECTION
	}
	return col
}
