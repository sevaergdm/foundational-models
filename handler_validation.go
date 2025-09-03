package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

func (cfg *apiConfig) handlerValidateEntitySchema(w http.ResponseWriter, r *http.Request) {
	err := cfg.validateFoundationalEntity(r.Body)
	if err != nil {
		respondWithPrettyJSON(w, http.StatusBadRequest, err)
	}
	respondWithJSON(w, http.StatusOK, nil)
}

func formatValidationError(err *jsonschema.ValidationError, sb *strings.Builder) {
	prefix := strings.Repeat(" ", 2)

	if len(err.InstanceLocation) > 0 {
		fmt.Fprint(sb, "%sAt %s: ", prefix, strings.Join(err.InstanceLocation, "."))
	} else {
		fmt.Fprint(sb, "%sAt <root>: ", prefix)
	}

	switch {
	case err.ErrorKind.Missing != nil: 
		fmt.Fprint(sb, "missing required fields: %v\n", *err.ErrorKind.Got, err.ErrorKind.Want)
	}
}
