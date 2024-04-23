package auth

const (
	TokenResponseTemplate     = `{"token": "%s", "refresh_token": "%s"}`
	EndPointVerifyEmail       = "%s?code=%s"
	DefaultURLVerifyEmail     = "http://localhost:8080/"
	SenderSignupEmail         = "info@doitall.work"
	DefaultSubjectSignUpEmail = "Welcome to DoItAll"
	DefaultSignUpBodyEmail    = "Welcome to DoItAll %s %s!.<br> Please verify your email by clicking the link below <br> <a href='%s'>Verify Email</a> <br> <br> DoItAll Team"
	InvalidToken              = "Invalid token %s"
	UserNotActivated          = "User not activated"
)
