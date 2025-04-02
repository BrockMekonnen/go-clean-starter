package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
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

func (v *Validator) GetBody(c echo.Context) (interface{}, error) {
	if v.schemas.Body == nil {
		var body map[string]interface{}
		if err := c.Bind(&body); err != nil {
			return nil, err
		}
		return body, nil
	}

	body := v.schemas.Body
	if err := c.Bind(&body); err != nil {
		return nil, err
	}

	if err := v.validate.Struct(body); err != nil {
		return nil, NewValidationError("body", err)
	}

	return body, nil
}

func (v *Validator) GetParams(c echo.Context) (map[string]string, error) {
	if v.schemas.Params == nil {
		params := make(map[string]string)
		for _, name := range c.ParamNames() {
			params[name] = c.Param(name)
		}
		return params, nil
	}

	// For params, we typically validate against the path parameters directly
	// since they're simple string values
	params := make(map[string]string)
	for _, name := range c.ParamNames() {
		params[name] = c.Param(name)
	}

	// If you need complex param validation, you would need to implement it here
	return params, nil
}

func (v *Validator) GetQuery(c echo.Context) (map[string]interface{}, error) {
	if v.schemas.Query == nil {
		query := make(map[string]interface{})
		for name, values := range c.QueryParams() {
			if len(values) == 1 {
				query[name] = values[0]
			} else {
				query[name] = values
			}
		}
		return query, nil
	}

	// Query validation would need custom implementation based on your needs
	// This is a simplified version
	query := make(map[string]interface{})
	for name, values := range c.QueryParams() {
		if len(values) == 1 {
			query[name] = values[0]
		} else {
			query[name] = values
		}
	}

	return query, nil
}

func (v *Validator) GetHeaders(c echo.Context) (map[string]string, error) {
	if v.schemas.Headers == nil {
		headers := make(map[string]string)
		for name := range c.Request().Header {
			headers[name] = c.Request().Header.Get(name)
		}
		return headers, nil
	}

	// Header validation would need custom implementation
	headers := make(map[string]string)
	for name := range c.Request().Header {
		headers[name] = c.Request().Header.Get(name)
	}

	return headers, nil
}

func (v *Validator) GetCookies(c echo.Context) (map[string]string, error) {
	if v.schemas.Cookies == nil {
		cookies := make(map[string]string)
		for _, cookie := range c.Cookies() {
			cookies[cookie.Name] = cookie.Value
		}
		return cookies, nil
	}

	// Cookie validation would need custom implementation
	cookies := make(map[string]string)
	for _, cookie := range c.Cookies() {
		cookies[cookie.Name] = cookie.Value
	}

	return cookies, nil
}

type ValidationError struct {
	Target string
	Err    error
}

func NewValidationError(target string, err error) *ValidationError {
	return &ValidationError{
		Target: target,
		Err:    err,
	}
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for %s: %v", e.Target, e.Err)
}