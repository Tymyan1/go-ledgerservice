package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

func DecodeJSON(w io.Writer, r io.Reader, dst interface{}) error {
	dec := json.NewDecoder(r)

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return errors.New(fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset))

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("Request body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			return errors.New(fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset))

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return errors.New(fmt.Sprintf("Request body contains unknown field %s", fieldName))

		case errors.Is(err, io.EOF):
			return errors.New("Request body must not be empty")

		default:
			return err
		}
	}
	return nil
}
