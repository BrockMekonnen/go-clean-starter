package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	customError "github.com/BrockMekonnen/go-clean-starter/core/lib/errors"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var (
	ErrMissingField    = errors.New("missing required field")
	ErrInvalidFilter   = errors.New("invalid filter format")
	DefaultPageSize    = 10
	DefaultPage        = 1
	DefaultSort        = []SortField{}
	DefaultFilterValue = map[string]interface{}{}
)

type FieldSource string

const (
	Query  FieldSource = "query"
	Params FieldSource = "params"
	Body   FieldSource = "body"
)

type FieldConfig struct {
	Name string
	From FieldSource
}

type SortField struct {
	Field     string
	Direction string
}

type PaginatorOptions struct {
	UseDefaults bool
	Fields      struct {
		Page     interface{}
		PageSize interface{}
		Sort     interface{}
		Filter   interface{}
	}
	Defaults struct {
		PageSize int
		Page     int
		Filter   map[string]interface{}
		Sort     []SortField
	}
	FilterSchema interface{}
}

type Paginator struct {
	opts     PaginatorOptions
	validate *validator.Validate
}

func NewPaginator(opts ...func(*PaginatorOptions)) *Paginator {
	config := PaginatorOptions{
		UseDefaults: true,
		Fields: struct {
			Page     interface{}
			PageSize interface{}
			Sort     interface{}
			Filter   interface{}
		}{
			Page:     "page",
			PageSize: "limit",
			Sort:     "sort",
			Filter:   "filter",
		},
		Defaults: struct {
			PageSize int
			Page     int
			Filter   map[string]interface{}
			Sort     []SortField
		}{
			PageSize: DefaultPageSize,
			Page:     DefaultPage,
			Sort:     DefaultSort,
			Filter:   DefaultFilterValue,
		},
	}

	for _, opt := range opts {
		opt(&config)
	}

	return &Paginator{
		opts:     config,
		validate: validator.New(),
	}
}

func (p *Paginator) getFieldConfig(field interface{}) FieldConfig {
	switch v := field.(type) {
	case string:
		return FieldConfig{Name: v, From: Query}
	case FieldConfig:
		return v
	default:
		return FieldConfig{Name: "", From: Query}
	}
}

func (p *Paginator) fromRequest(r *http.Request, field FieldConfig) interface{} {
	switch field.From {
	case Query:
		return r.URL.Query().Get(field.Name)
	case Body:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return nil
		}
		defer r.Body.Close()

		var parsed map[string]interface{}
		if err := json.Unmarshal(body, &parsed); err != nil {
			return nil
		}
		return parsed[field.Name]
	case Params:
		vars := mux.Vars(r)
		return vars[field.Name]
	default:
		return nil
	}
}

func (p *Paginator) GetPagination(r *http.Request) (int, int, error) {
	pageField := p.getFieldConfig(p.opts.Fields.Page)
	pageSizeField := p.getFieldConfig(p.opts.Fields.PageSize)

	pageStr := p.fromRequest(r, pageField)
	pageSizeStr := p.fromRequest(r, pageSizeField)

	page := p.opts.Defaults.Page
	pageSize := p.opts.Defaults.PageSize

	// Parse and validate page
	if pageStr != nil && pageStr != "" {
		pageInt, err := strconv.Atoi(fmt.Sprintf("%v", pageStr))
		if err != nil || pageInt <= 0 {
			return 0, 0, customError.NewBadRequestError[any](
				fmt.Sprintf("Invalid '%s.%s' value: must be a positive integer", pageField.From, pageField.Name), "", nil)
		}
		page = pageInt
	} else if !p.opts.UseDefaults {
		return 0, 0, customError.NewBadRequestError[any](
			fmt.Sprintf("Missing '%s.%s' value", pageField.From, pageField.Name), "", nil)
	}

	// Parse and validate pageSize
	if pageSizeStr != nil && pageSizeStr != "" {
		sizeInt, err := strconv.Atoi(fmt.Sprintf("%v", pageSizeStr))
		if err != nil || sizeInt <= 0 {
			return 0, 0, customError.NewBadRequestError[any](
				fmt.Sprintf("Invalid '%s.%s' value: must be a positive integer", pageSizeField.From, pageSizeField.Name), "", nil)
		}
		pageSize = sizeInt
	} else if !p.opts.UseDefaults {
		return 0, 0, customError.NewBadRequestError[any](
			fmt.Sprintf("Missing '%s.%s' value", pageSizeField.From, pageSizeField.Name), "", nil)
	}

	return page, pageSize, nil
}

func (p *Paginator) GetSorter(r *http.Request) ([]SortField, error) {
	sortField := p.getFieldConfig(p.opts.Fields.Sort)
	sortValue := p.fromRequest(r, sortField)

	if sortValue == nil || sortValue == "" {
		if !p.opts.UseDefaults {
			return nil, customError.NewBadRequestError[any](
				fmt.Sprintf("Missing '%s.%s' value", sortField.From, sortField.Name), "", nil)
		}
		return p.opts.Defaults.Sort, nil
	}

	var sortList []string
	switch v := sortValue.(type) {
	case string:
		sortList = strings.Split(v, ",")
	case []string:
		sortList = v
	default:
		return nil, customError.NewBadRequestError[any]("Invalid sort format", "", nil)
	}

	result := make([]SortField, 0, len(sortList))
	for _, sort := range sortList {
		sort = strings.TrimSpace(sort)
		if sort == "" {
			continue
		}

		direction := "asc"
		if strings.HasPrefix(sort, "-") {
			direction = "desc"
			sort = sort[1:]
		}

		result = append(result, SortField{Field: sort, Direction: direction})
	}

	if len(result) == 0 {
		return p.opts.Defaults.Sort, nil
	}

	return result, nil
}

func (p *Paginator) GetFilter(r *http.Request) (map[string]interface{}, error) {
	filterField := p.getFieldConfig(p.opts.Fields.Filter)
	filterValue := p.fromRequest(r, filterField)

	if filterValue == nil || filterValue == "" {
		if !p.opts.UseDefaults {
			return nil, customError.NewBadRequestError[any](
				fmt.Sprintf("Missing '%s.%s' value", filterField.From, filterField.Name), "", nil)
		}
		return p.opts.Defaults.Filter, nil
	}

	var filter map[string]interface{}
	switch v := filterValue.(type) {
	case map[string]interface{}:
		filter = v
	case string:
		if err := json.Unmarshal([]byte(v), &filter); err != nil {
			return nil, customError.NewBadRequestError[any](
				fmt.Sprintf("Invalid '%s.%s' format: not valid JSON", filterField.From, filterField.Name), "", nil)
		}
	default:
		return nil, customError.NewBadRequestError[any]("Invalid filter format", "", nil)
	}

	if p.opts.FilterSchema != nil {
		err := p.validate.Struct(p.opts.FilterSchema)
		if err != nil {
			if verr, ok := err.(validator.ValidationErrors); ok {
				return nil, customError.NewValidationError("Invalid filter parameters", "", verr)
			}
			return nil, err
		}
	}

	return filter, nil
}