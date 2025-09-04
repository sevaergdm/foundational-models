package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
)

func (cfg *apiConfig) handlerUpdateEntity(w http.ResponseWriter, r *http.Request) {
	entityName := r.PathValue("entityName")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to read request body", err)
		return
	}
	defer r.Body.Close()

	validationReader := bytes.NewReader(body)
	err = cfg.validateFoundationalEntity(validationReader)
	if err != nil {
		if verr, ok := err.(*jsonschema.ValidationError); ok {
			respondWithPrettyJSON(w, http.StatusBadRequest, map[string]any{
				"errors": FormatValidationError(verr),
			})
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Unknown error", err)
		return
	}

	cachedFilePath := fmt.Sprintf("entities/%s.json", strings.ToLower(entityName))
	cachedJSON, err := os.ReadFile(cachedFilePath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to open cached entity file", err)
		return
	}

	var mappedBody map[string]any
	var mappedCachedJSON map[string]any
	err = json.Unmarshal(body, &mappedBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to unmarshal body", err)
		return
	}

	err = json.Unmarshal(cachedJSON, &mappedCachedJSON)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to unmarshal cachedJSON", err)
		return
	}
	
	differ := gojsondiff.New()
	diff, err := differ.Compare(cachedJSON, body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to compare json files", err)
		return
	}
	
	if !diff.Modified() {
		respondWithError(w, http.StatusBadRequest, "No changes detected", nil)
		return
	}
	
	f := formatter.NewAsciiFormatter(mappedCachedJSON, formatter.AsciiFormatterDefaultConfig)
	diffReport, err := f.Format(diff)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to format output report", err)
		return
	}

	fmt.Println(diffReport)



//	pattern := `[+-][\s\x{A0}]+Version:`
//	re, err := regexp.Compile(pattern)
//	if err != nil {
//		respondWithError(w, http.StatusInternalServerError, "Unable to compile regexp", err)
//		return
//	}
//
//	if !re.MatchString(diff) {
//		respondWithError(w, http.StatusBadRequest, "Version was not changed with update", nil)
//		return
//	}
//
//	err = os.WriteFile(fmt.Sprintf("entities/%s.json", strings.ToLower(entity.Name)), body, 0644)
//	if err != nil {
//		respondWithError(w, http.StatusInternalServerError, "Unable to update entity file", err)
//	}
//
//	respondWithJSON(w, http.StatusOK, fmt.Sprintf("%s was successfully updated", entity.Name))
//	fmt.Println(diff)
}
