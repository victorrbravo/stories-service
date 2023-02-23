package crypto

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"log"
)

// Token cryptographically secure token
type Token struct {
	hash []byte
}

// Value implements the sql.Valuer interface
func (token Token) Value() (driver.Value, error) {
	return token.digest(), nil
}

// NewToken creates a new token with the given hash
func NewToken(hash []byte) Token {
	return Token{hash: hash}
}

// Equals compares two tokens for equality and returns true if they're found equal
func (token Token) Equals(other Token) bool {
	hash1 := token.hash
	hash2 := other.hash
	if hash1 == nil || hash2 == nil {
		log.Println("some token has a `nil' hash")
		return false
	}
	return bytes.Compare(hash1, hash2) == 0
}

// NewTokenFromString creates a new token given a string as a hash
func NewTokenFromString(hash string) (Token, error) {
	var data []byte
	var err error
	data, err = hex.DecodeString(hash)
	if err != nil {
		return Token{}, err
	}
	return Token{hash: data}, nil
}

func (token Token) digest() string {
	return hex.EncodeToString(token.hash)
}

// IsValid checks whether parameters are valid or returns an error
func (token Token) IsValid() bool {
	return token.hash != nil
}

// MarshalText returns a text representation of the token hash digest as a byte array
func (token Token) MarshalText() ([]byte, error) {
	return []byte(token.digest()), nil
}

// UnmarshalText builds a token given a text representation
func (token *Token) UnmarshalText(data []byte) error {
	var err error
	token.hash, err = hex.DecodeString(string(data))
	if err != nil {
		token.hash = nil
		// We do not want to error here, we instead will
		// let the caller think that the token is invalid
		// which it is
		return nil
	}
	return nil
}

func (token Token) String() string {
	return token.digest()
}
