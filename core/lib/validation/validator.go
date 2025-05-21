package validation

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/errors"
)

type ValidationSchemas struct {
	Body    interface{}
	Params  interface{}
	Query   interface{}
	Headers interface{}
	Cookies interface{}
}

type Validator struct {
	validate *validator.Validate
	schemas  ValidationSchemas
}

func NewValidator(schemas ValidationSchemas) *Validator {
	return &Validator{
		validate: validator.New(),
		schemas:  schemas,
	}
}

func formatValidationError(fe validator.FieldError) string {
	return "Field validation for '" + fe.Field() + "' failed on the '" + fe.Tag() + "' tag"
}

func (v *Validator) GetBody(r *http.Request) (interface{}, error) {
	if v.schemas.Body == nil {
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			return nil, err
		}
		return body, nil
	}

	body := v.schemas.Body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	if err := v.validate.Struct(body); err != nil {
		if verr, ok := err.(validator.ValidationErrors); ok {
			msg := formatValidationError(verr[0])
			return nil, errors.NewValidationError(msg, "body", verr)
		}
		return nil, err
	}

	return body, nil
}

func (v *Validator) GetParams(r *http.Request) (map[string]string, error) {
	vars := mux.Vars(r)
	if v.schemas.Params == nil {
		return vars, nil
	}
	//TODO Optional: Add validation logic for params here
	return vars, nil
}

func (v *Validator) GetQuery(r *http.Request) (map[string]interface{}, error) {
	query := make(map[string]interface{})
	q := r.URL.Query()

	for key, values := range q {
		if len(values) == 1 {
			query[key] = values[0]
		} else {
			query[key] = values
		}
	}

	if v.schemas.Query == nil {
		return query, nil
	}

	data, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	querySchema := v.schemas.Query
	if err := json.Unmarshal(data, &querySchema); err != nil {
		return nil, err
	}

	if err := v.validate.Struct(querySchema); err != nil {
		if verr, ok := err.(validator.ValidationErrors); ok {
			msg := formatValidationError(verr[0]) // Get cleaned message
			return nil, errors.NewValidationError(msg, "query", verr)
		}
		return nil, err
	}

	return query, nil
}