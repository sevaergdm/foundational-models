package main

import (
	"net/http"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

func (cfg *apiConfig) handlerValidateEntitySchema(w http.ResponseWriter, r *http.Request) {
	err := cfg.validateFoundationalEntity(r.Body)
	if err != nil {
		if verr, ok := err.(*jsonschema.ValidationError); ok {
			respondWithPrettyJSON(w, http.StatusBadRequest, map[string]any{
				"errors": FormatValidationError(verr),
			})
			return
		}
		respondWithError(w, http.StatusBadRequest, "Unkown Error", err)
		return
	}
	respondWithJSON(w, http.StatusOK, "Schema is valid")
}
