package main

import (
	"github.com/santhosh-tekuri/jsonschema/v6"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io"
	"strings"
)

var defaultPrinter = message.NewPrinter(language.English)

func (cfg *apiConfig) validateFoundationalEntity(reader io.Reader) error {
	data, err := jsonschema.UnmarshalJSON(reader)
	if err != nil {
		return err
	}

	err = cfg.compiledCanonicalSchema.Validate(data)
	if err != nil {
		return err
	}
	return nil
}

type ValidationErrorMessage struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

func FormatValidationError(err *jsonschema.ValidationError) []ValidationErrorMessage {
	var errors []ValidationErrorMessage
	collectErrors(err, &errors)
	return errors
}

func collectErrors(err *jsonschema.ValidationError, errors *[]ValidationErrorMessage) {
	var msg string
	if err.ErrorKind != nil {
		msg = err.ErrorKind.LocalizedString(defaultPrinter)
	} else {
		msg = "(unspecified error)"
	}

	path := "root"
	if len(err.InstanceLocation) > 0 {
		path = strings.Join(err.InstanceLocation, ".")
	}

	if len(err.InstanceLocation) == 0 && len(err.Causes) > 0 {
	} else {
		*errors = append(*errors, ValidationErrorMessage{
			Path: path,
			Message: msg,
		})
	} 

	for _, cause := range err.Causes {
		collectErrors(cause, errors)
	}
}
