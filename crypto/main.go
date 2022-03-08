package crypto

import (
	"crypto/rand"
	"crypto/sha256"

	"golang.org/x/crypto/bcrypt"
)

func random(size int) ([]byte, error) {
	var bytes = make([]byte, size)
	var err error
	// recv random bytes
	_, err = rand.Read(bytes)
	// Always check for errors
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// RandomHash generates a sha256 random hash
func RandomHash() ([]byte, error) {
	var sha = sha256.New()
	var err error
	var rnd []byte
	// get something random, or not
	rnd, err = random(100)
	if err != nil {
		return nil, err
	}
	// write some random stuff
	sha.Write(rnd)
	// Return it as a string
	return sha.Sum(nil), nil
}

// HashPassword builds the hash for a text representation of a password
func HashPassword(password string) ([]byte, error) {
	var hashBytes []byte
	var err error
	// Create the hash here
	hashBytes, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	return hashBytes, nil
}

// ComparePasswords compares the hash and the passwords returning an error if they're not equal
func ComparePasswords(hashed Password, plain string) error {
	var err error
	err = bcrypt.CompareHashAndPassword(hashed, []byte(plain))
	if err != nil {
		return err
	}
	return nil
}
