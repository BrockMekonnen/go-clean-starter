package extension

type contextKey string

const (
	AuthContextKey contextKey = "auth"
)

type AuthContext struct {
	IsAuthenticated bool
	IsAuthorized    bool
	IsInjected      bool
	Credentials     struct {
		UID   string
		Scope []string
	}
	Artifacts struct {
		AccessToken string
	}
}
