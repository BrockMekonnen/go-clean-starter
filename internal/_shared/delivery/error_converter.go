package delivery

import (
	"net/http"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/errors"
	sharedDomain "github.com/BrockMekonnen/go-clean-starter/internal/_shared/domain"
)

// ErrorConverter defines the interface for converting errors to HTTP responses
type ErrorConverter interface {
	Test(err error) bool
	Convert(err error) (int, interface{})
}

// errorConverter implements the ErrorConverter interface
type errorConverter struct {
	test    func(error) bool
	convert func(error) (int, interface{})
}

// NewErrorConverter creates a new ErrorConverter
func NewErrorConverter(
	test func(error) bool,
	convert func(error) (int, interface{}),
) ErrorConverter {
	return &errorConverter{
		test:    test,
		convert: convert,
	}
}

func (ec *errorConverter) Test(err error) bool {
	return ec.test(err)
}

func (ec *errorConverter) Convert(err error) (int, interface{}) {
	return ec.convert(err)
}

// DefaultErrorConverters provides the standard set of error converters
var DefaultErrorConverters = []ErrorConverter{
	// ValidationError converter
	NewErrorConverter(
		func(err error) bool {
			_, ok := err.(*errors.ValidationError)
			return ok
		},
		func(err error) (int, interface{}) {
			ve := err.(*errors.ValidationError)
			status := http.StatusUnprocessableEntity
			if ve.Meta.Target != "body" {
				status = http.StatusBadRequest
			}

			details := make([]map[string]interface{}, len(ve.Meta.Error))
			for i, d := range ve.Meta.Error {
				details[i] = map[string]interface{}{
					"field": d.Field,
					"path":  d.Value(),
				}
			}

			return status, map[string]interface{}{
				"error":   ve.Name,
				"code":    ve.Code,
				"status":  status,
				"message": ve.Message,
			}
		},
	),

	// BadRequestError converter
	NewErrorConverter(
		func(err error) bool {
			_, ok := err.(*errors.BadRequestError)
			return ok
		},
		func(err error) (int, interface{}) {
			be := err.(*errors.BadRequestError)
			return http.StatusBadRequest, map[string]interface{}{
				"error":   be.Name,
				"code":    be.Code,
				"status":  http.StatusBadRequest,
				"message": be.Message,
			}
		},
	),

	// NotFoundError converter
	NewErrorConverter(
		func(err error) bool {
			_, ok := err.(*errors.NotFoundError)
			return ok
		},
		func(err error) (int, interface{}) {
			ne := err.(*errors.NotFoundError)
			return http.StatusNotFound, map[string]interface{}{
				"error":   ne.Name,
				"code":    ne.Code,
				"status":  http.StatusNotFound,
				"message": ne.Message,
			}
		},
	),

	// UnauthorizedError converter
	NewErrorConverter(
		func(err error) bool {
			_, ok := err.(*errors.UnauthorizedError)
			return ok
		},
		func(err error) (int, interface{}) {
			ue := err.(*errors.UnauthorizedError)
			return http.StatusUnauthorized, map[string]interface{}{
				"error":   ue.Name,
				"code":    ue.Code,
				"status":  http.StatusUnauthorized,
				"message": ue.Message,
			}
		},
	),

	// ForbiddenError converter
	NewErrorConverter(
		func(err error) bool {
			_, ok := err.(*errors.ForbiddenError)
			return ok
		},
		func(err error) (int, interface{}) {
			fe := err.(*errors.ForbiddenError)
			return http.StatusForbidden, map[string]interface{}{
				"error":   fe.Name,
				"code":    fe.Code,
				"status":  http.StatusForbidden,
				"message": fe.Message,
			}
		},
	),

	// BusinessError converter
	NewErrorConverter(
		func(err error) bool {
			_, ok := err.(*sharedDomain.BusinessError)
			return ok
		},
		func(err error) (int, interface{}) {
			be := err.(*sharedDomain.BusinessError)
			return http.StatusConflict, map[string]interface{}{
				"error":   be.Name,
				"code":    be.Code,
				"status":  http.StatusConflict,
				"message": be.Message,
			}
		},
	),

	// BaseError fallback converter
	NewErrorConverter(
		func(err error) bool {
			_, ok := err.(*errors.BaseError[any])
			return ok
		},
		func(err error) (int, interface{}) {
			be := err.(*errors.BaseError[any])
			return http.StatusInternalServerError, map[string]interface{}{
				"error":   be.Name,
				"code":    be.Code,
				"status":  http.StatusInternalServerError,
				"message": be.Message,
			}
		},
	),
}