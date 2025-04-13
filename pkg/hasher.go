package pkg

import (
	"crypto/sha256"
	"fmt"
)

type Hasher interface {
	Hash(password string) string
	Check(hashedPassword, password string) bool
}

type SHAhash struct {
	Salt string
}

func NewSHAhash(salt string) SHAhash {
	return SHAhash{salt}
}

func (s *SHAhash) Hash(password string) string {
	hashed := sha256.Sum256([]byte(password + s.Salt))
	return fmt.Sprintf("%x", hashed)
}

func (s *SHAhash) Check(hashedPassword, password string) bool {
	hashed := sha256.Sum256([]byte(password + s.Salt))
	return hashedPassword == fmt.Sprintf("%x", hashed)
}
