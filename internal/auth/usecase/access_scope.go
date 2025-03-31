package usecase

// ScopeParams defines the structure for user and allowed scopes
type ScopeParams struct {
	UserScope    []string
	AllowedScope []string
}

// HasRole is the interface for the service
type HasRole interface {
	Execute(params ScopeParams) (bool, error)
}

// ScopeService implements the HasRole interface
type ScopeService struct{}

// NewScopeService returns a new instance of ScopeService
func NewScopeService() HasRole {
	return &ScopeService{}
}

// Execute checks if userScope includes any of the allowedScope
func (s *ScopeService) Execute(params ScopeParams) (bool, error) {
	for _, role := range params.UserScope {
		for _, allowed := range params.AllowedScope {
			if role == allowed {
				return true, nil
			}
		}
	}
	return false, nil
}