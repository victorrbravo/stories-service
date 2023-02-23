package constants

// HTTP errors
const (
	UserNotFoundOrPasswordMismatch = "User/Password is wrong"
	UserBlockedByLoginAttempts     = "You've attempted to login with the wrong password 5 times (or more). Your account was blocked. We've sent you an email for the activation of your account."
	UserNotFound                   = "user not found"
	MalformedJSON                  = "malformed json"
)
