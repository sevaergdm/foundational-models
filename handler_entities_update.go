package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
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

	diffReport, err := entityDiff(cachedJSON, body)
	if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Unable to create diff report", err)
			return
	}
	if diffReport == SchemasMatch {
		respondWithError(w, http.StatusBadRequest, "Couldn't update because schemas are the same", nil)
		return
	}
	if diffReport == VersionNotUptdated {
		respondWithError(w, http.StatusBadRequest, "Couldn't update schema because the version was not upated", nil)
		return
	}

	err = os.WriteFile(fmt.Sprintf("entities/%s.json", strings.ToLower(entityName)), body, 0644)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update entity file", err)
		return
	}

	respondWithJSON(w, http.StatusOK, fmt.Sprintf("%s was successfully updated", entityName))
	fmt.Println(diffReport)
}

