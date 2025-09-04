package main

import (
	"encoding/json"
	"regexp"

	"github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
)

const SchemasMatch = "Schemas are the same"
const VersionNotUptdated = "Version not updated"

func entityDiff(cachedEntity, requestEntity []byte) (string, error) {
	var mappedCachedJSON map[string]any
	err := json.Unmarshal(cachedEntity, &mappedCachedJSON)
	if err != nil {
		return "", err
	}

	var mappedBody map[string]any
	err = json.Unmarshal(requestEntity, &mappedBody)
	if err != nil {
		return "", err
	}
	
	differ := gojsondiff.New()
	diff := differ.CompareObjects(mappedCachedJSON, mappedBody)
	
	if !diff.Modified() {
		return SchemasMatch, nil
	}

	config := formatter.AsciiFormatterConfig{
		ShowArrayIndex: false,
		Coloring: true,
	}
	
	f := formatter.NewAsciiFormatter(mappedCachedJSON, config)
	diffReport, err := f.Format(diff)
	if err != nil {
		return "", err
	}


	pattern := `[+-][\s\x{A0}]+"version":`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}

	if !re.MatchString(diffReport) {
		return VersionNotUptdated, nil
	}
	
	return diffReport, nil
}
