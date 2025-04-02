package validation

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
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
	Direction string // "asc" or "desc"
}

type PaginatorOptions struct {
	UseDefaults bool
	Fields      struct {
		Page     interface{} // string or FieldConfig
		PageSize interface{} // string or FieldConfig
		Sort     interface{} // string or FieldConfig
		Filter   interface{} // string or FieldConfig
	}
	Defaults struct {
		PageSize int
		Page     int
		Filter   map[string]interface{}
		Sort     []SortField
	}
	FilterSchema interface{} // validation schema
}

type Paginator struct {
	opts PaginatorOptions
	validate *validator.Validate
}

func NewPaginator(opts ...func(*PaginatorOptions)) *Paginator {
	// Set default options
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

	// Apply custom options
	for _, opt := range opts {
		opt(&config)
	}

	return &Paginator{
		opts: config,
		validate: validator.New(),
	}
}

// WithDefaults option sets whether to use default values
func WithDefaults(use bool) func(*PaginatorOptions) {
	return func(o *PaginatorOptions) {
		o.UseDefaults = use
	}
}

// WithPageField option configures the page field
func WithPageField(field interface{}) func(*PaginatorOptions) {
	return func(o *PaginatorOptions) {
		o.Fields.Page = field
	}
}

// WithPageSizeField option configures the page size field
func WithPageSizeField(field interface{}) func(*PaginatorOptions) {
	return func(o *PaginatorOptions) {
		o.Fields.PageSize = field
	}
}

// WithSortField option configures the sort field
func WithSortField(field interface{}) func(*PaginatorOptions) {
	return func(o *PaginatorOptions) {
		o.Fields.Sort = field
	}
}

// WithFilterField option configures the filter field
func WithFilterField(field interface{}) func(*PaginatorOptions) {
	return func(o *PaginatorOptions) {
		o.Fields.Filter = field
	}
}

// WithDefaultPageSize option sets the default page size
func WithDefaultPageSize(size int) func(*PaginatorOptions) {
	return func(o *PaginatorOptions) {
		o.Defaults.PageSize = size
	}
}

// WithDefaultPage option sets the default page number
func WithDefaultPage(page int) func(*PaginatorOptions) {
	return func(o *PaginatorOptions) {
		o.Defaults.Page = page
	}
}

// WithDefaultSort option sets the default sort fields
func WithDefaultSort(sort []SortField) func(*PaginatorOptions) {
	return func(o *PaginatorOptions) {
		o.Defaults.Sort = sort
	}
}

// WithDefaultFilter option sets the default filter
func WithDefaultFilter(filter map[string]interface{}) func(*PaginatorOptions) {
	return func(o *PaginatorOptions) {
		o.Defaults.Filter = filter
	}
}

// WithFilterSchema option sets the filter validation schema
func WithFilterSchema(schema interface{}) func(*PaginatorOptions) {
	return func(o *PaginatorOptions) {
		o.FilterSchema = schema
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

func (p *Paginator) fromRequest(c echo.Context, field FieldConfig) interface{} {
	switch field.From {
	case Query:
		return c.QueryParam(field.Name)
	case Params:
		return c.Param(field.Name)
	case Body:
		var body map[string]interface{}
		if err := c.Bind(&body); err == nil {
			return body[field.Name]
		}
		return nil
	default:
		return nil
	}
}

func (p *Paginator) GetPagination(c echo.Context) (page, pageSize int, err error) {
	pageField := p.getFieldConfig(p.opts.Fields.Page)
	pageSizeField := p.getFieldConfig(p.opts.Fields.PageSize)

	pageStr := p.fromRequest(c, pageField)
	pageSizeStr := p.fromRequest(c, pageSizeField)

	page = p.opts.Defaults.Page
	pageSize = p.opts.Defaults.PageSize

	if pageStr != nil && pageStr != "" {
		page, err = strconv.Atoi(fmt.Sprintf("%v", pageStr))
		if err != nil {
			if !p.opts.UseDefaults {
				return 0, 0, fmt.Errorf("%w: %s", ErrMissingField, pageField.Name)
			}
			page = p.opts.Defaults.Page
		}
	} else if !p.opts.UseDefaults {
		return 0, 0, fmt.Errorf("%w: %s", ErrMissingField, pageField.Name)
	}

	if pageSizeStr != nil && pageSizeStr != "" {
		pageSize, err = strconv.Atoi(fmt.Sprintf("%v", pageSizeStr))
		if err != nil {
			if !p.opts.UseDefaults {
				return 0, 0, fmt.Errorf("%w: %s", ErrMissingField, pageSizeField.Name)
			}
			pageSize = p.opts.Defaults.PageSize
		}
	} else if !p.opts.UseDefaults {
		return 0, 0, fmt.Errorf("%w: %s", ErrMissingField, pageSizeField.Name)
	}

	return page, pageSize, nil
}

func (p *Paginator) GetSorter(c echo.Context) ([]SortField, error) {
	sortField := p.getFieldConfig(p.opts.Fields.Sort)
	sortValue := p.fromRequest(c, sortField)

	if sortValue == nil || sortValue == "" {
		if !p.opts.UseDefaults {
			return nil, fmt.Errorf("%w: %s", ErrMissingField, sortField.Name)
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
		if !p.opts.UseDefaults {
			return nil, fmt.Errorf("%w: invalid sort format", ErrInvalidFilter)
		}
		return p.opts.Defaults.Sort, nil
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

		result = append(result, SortField{
			Field:     sort,
			Direction: direction,
		})
	}

	if len(result) == 0 {
		return p.opts.Defaults.Sort, nil
	}

	return result, nil
}

func (p *Paginator) GetFilter(c echo.Context) (map[string]interface{}, error) {
	filterField := p.getFieldConfig(p.opts.Fields.Filter)
	filterValue := p.fromRequest(c, filterField)

	if filterValue == nil || filterValue == "" {
		if !p.opts.UseDefaults {
			return nil, fmt.Errorf("%w: %s", ErrMissingField, filterField.Name)
		}
		return p.opts.Defaults.Filter, nil
	}

	filter, ok := filterValue.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%w: invalid filter format", ErrInvalidFilter)
	}

	// Validate filter against schema if provided
	if p.opts.FilterSchema != nil {
		if err := p.validate.Struct(filter); err != nil {
			return nil, fmt.Errorf("filter validation failed: %w", err)
		}
	}

	return filter, nil
}