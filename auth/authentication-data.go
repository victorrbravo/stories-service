package auth

import "stories/crypto"

type AuthenticationData struct {
	Token           crypto.Token `json:"token,omitempty"`
	SessionCount    int          `json:"session_count"`
	User            User         `json:"user"`
	Invalid         bool         `json:"invalid"`
	Preregistration bool         `json:"pre_registration"`
}

type AuthData struct {
	AuthenticationData
}
