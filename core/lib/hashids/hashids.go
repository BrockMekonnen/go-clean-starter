package hashids

import (
	"errors"

	"github.com/BrockMekonnen/go-clean-starter/core"
	"github.com/speps/go-hashids/v2"
)

// HashID defines the interface for encoding and decoding IDs.
type HashID interface {
	EncodeID(id uint) (string, error)
	DecodeID(hash string) (uint, error)
}

// hashIDService is the concrete implementation of HashID.
type hashIDService struct {
	hasher *hashids.HashID
}

// NewHashIDService initializes a new hashIDService with the given salt and minimum length.
func NewHashIDService(config core.EncryptionConfig) (HashID, error) {
	hd := hashids.NewData()
	hd.Salt = config.HashSalt
	hd.MinLength = 8

	hasher, err := hashids.NewWithData(hd)
	if err != nil {
		return nil, err
	}

	return &hashIDService{hasher: hasher}, nil
}

// EncodeID encodes a uint ID into a hashed string.
func (s *hashIDService) EncodeID(id uint) (string, error) {
	// Convert uint to int
	return s.hasher.Encode([]int{int(id)})
}

// DecodeID decodes a hashed string back to a uint ID.
func (s *hashIDService) DecodeID(hash string) (uint, error) {
	numbers, err := s.hasher.DecodeWithError(hash)
	if err != nil || len(numbers) == 0 {
		return 0, err
	}
	if numbers[0] < 0 {
		return 0, errors.New("decoded ID is negative, cannot convert to uint")
	}
	return uint(numbers[0]), nil
}
