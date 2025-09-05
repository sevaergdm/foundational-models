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
	"github.com/sevaergdm/foundational-models/model_types"
)

func (cfg *apiConfig) handlerCreateEntity(w http.ResponseWriter, r *http.Request) {
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

	var entity model_types.FoundationalModel
	err = json.Unmarshal(body, &entity)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid format", err)
		return
	}
	
	if _, ok := cfg.entitiesCache[entity.Name]; ok {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Entity: %s already exists", entity.Name), err)
		return
	}

	err = os.WriteFile(fmt.Sprintf("entities/%s.json", strings.ToLower(entity.Name)), body, 0644)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to write to file", err)
	}

	err = cfg.loadEntities("entities")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to load entities after creation", err)
		return
	}

	respondWithJSON(w, http.StatusOK, fmt.Sprintf("New Foundational model: %s created", entity.Name))

}
