package domain

type Credentials struct {
	Uid   uint  
	Scope []string
}

type AuthRepository interface {
	Generate(payload Credentials) (string, error)
	Decode(accessToken string) (*Credentials, error) // Changed return type
	Hash(plainText string) (string, error)
	Compare(plainText, hashedText string) (bool, error)
}