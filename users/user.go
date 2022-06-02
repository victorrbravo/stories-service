package users

import (
	"time"

	"stories/types"
)

// UserType a special type of int useful for validation only
type UserType int

// User types enumerator
const (
	UserTypeAdmin UserType = 1
	UserTypeWeb            = 2
)

// User a generic user interface to implement all types of users in the context of COMPA
type User interface {
	GetEmail() string
	GetPassword() string
}

// DummyUser is a no-op user
type DummyUser struct{}

// GetID returns the id of a user
func (dummyUser DummyUser) GetID() types.UUID {
	return types.UUID{}
}

// GetID returns the id of a user
func (dummyUser DummyUser) GetPeriodStartedAt() time.Time {
	return time.Now()
}

// GetEmail returns the email of a user
func (DummyUser) GetEmail() string {
	return ""
}

// GetUserName returns the user name of a user
func (DummyUser) GetUserName() string {
	return ""
}

// Gender returns the gender type of a user
func (DummyUser) GetGender() types.Gender {
	return 0
}

// GetFullName returns the full name of a user
func (DummyUser) GetFullName() string {
	return ""
}

// GetSessionCount returns the number of session the users has
func (DummyUser) GetSessionCount() int {
	return 0
}
