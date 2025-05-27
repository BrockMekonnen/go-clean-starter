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

	if err := json.NewDecoder(r.Body).Decode(v.schemas.Body); err != nil {
		return nil, err
	}

	if err := v.validate.Struct(v.schemas.Body); err != nil {
		if verr, ok := err.(validator.ValidationErrors); ok {
			msg := formatValidationError(verr[0])
			return nil, errors.NewValidationError(msg, "body", verr)
		}
		return nil, err
	}

	return v.schemas.Body, nil
}

func (v *Validator) GetParams(r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	if v.schemas.Params == nil {
		return vars, nil
	}

	data, err := json.Marshal(vars)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, v.schemas.Params); err != nil {
		return nil, err
	}

	if err := v.validate.Struct(v.schemas.Params); err != nil {
		if verr, ok := err.(validator.ValidationErrors); ok {
			msg := formatValidationError(verr[0])
			return nil, errors.NewValidationError(msg, "params", verr)
		}
		return nil, err
	}

	return v.schemas.Params, nil
}

func (v *Validator) GetQuery(r *http.Request) (interface{}, error) {
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

	if err := json.Unmarshal(data, v.schemas.Query); err != nil {
		return nil, err
	}

	if err := v.validate.Struct(v.schemas.Query); err != nil {
		if verr, ok := err.(validator.ValidationErrors); ok {
			msg := formatValidationError(verr[0])
			return nil, errors.NewValidationError(msg, "query", verr)
		}
		return nil, err
	}

	return v.schemas.Query, nil
}

func (v *Validator) BindAndValidateQuery(r *http.Request, target interface{}) error {
	q := r.URL.Query()
	query := make(map[string]interface{})

	for key, values := range q {
		if len(values) == 1 {
			query[key] = values[0]
		} else {
			query[key] = values
		}
	}

	data, err := json.Marshal(query)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, target); err != nil {
		return err
	}

	if err := v.validate.Struct(target); err != nil {
		if verr, ok := err.(validator.ValidationErrors); ok {
			msg := formatValidationError(verr[0])
			return errors.NewValidationError(msg, "query", verr)
		}
		return err
	}

	return nil
}

func (v *Validator) BindAndValidateBody(r *http.Request, target interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		return err
	}

	if err := v.validate.Struct(target); err != nil {
		if verr, ok := err.(validator.ValidationErrors); ok {
			msg := formatValidationError(verr[0])
			return errors.NewValidationError(msg, "body", verr)
		}
		return err
	}

	return nil
}

func (v *Validator) BindAndValidateParams(r *http.Request, target interface{}) error {
	vars := mux.Vars(r)
	data, err := json.Marshal(vars)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, target); err != nil {
		return err
	}

	if err := v.validate.Struct(target); err != nil {
		if verr, ok := err.(validator.ValidationErrors); ok {
			msg := formatValidationError(verr[0])
			return errors.NewValidationError(msg, "params", verr)
		}
		return err
	}

	return nil
}
