package auth

const (
	TokenResponseTemplate     = `{"token": "%s", "refresh_token": "%s"}`
	EndPointVerifyEmail       = "%s?code=%s"
	DefaultURLVerifyEmail     = "http://localhost:8080/confirm/"
	DefaultURLResendEmail     = "http://localhost:8080/resend/"
	SenderSignupEmail         = "info@doitall.work"
	DefaultSubjectSignUpEmail = "Welcome to DoitAll"
	DefaultResendSubjectEmail = "New Code to DoitAll"
	DefaultChangePassword     = "DoitAll - request for change your password"
	DefaultSignUpBodyEmail    = "Welcome to <b>DoitAll</b> %s %s!.<br> Please verify your email by clicking the link below <br> <a href='%s'>Verify now</a> <br> <br> DoitAll Team"
	InvalidToken              = "Invalid token %s"
	UserNotActivated          = "User not activated"
	ExpiredToken              = "Token %s expired"
)
