package web

import "stories/crypto"
import "stories/generated/models"

type AuthenticationData struct {
	Token           crypto.Token `json:"token,omitempty"`
	SessionCount    int          `json:"session_count"`
	User            models.User  `json:"user"`
	Invalid         bool         `json:"invalid"`
	Preregistration bool         `json:"pre_registration"`
}

type AuthData struct {
	AuthenticationData
}
