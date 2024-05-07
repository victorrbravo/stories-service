package auth

const (
	TokenResponseTemplate     = `{"token": "%s", "refresh_token": "%s"}`
	EndPointVerifyEmail       = "%s?code=%s"
	DefaultURLVerifyEmail     = "http://localhost:8080/confirm/"
	DefaultURLResendEmail     = "http://localhost:8080/resend/"
	SenderSignupEmail         = "info@doitall.work"
	DefaultSubjectSignUpEmail = "Welcome to Use2"
	DefaultResendSubjectEmail = "New Code to Use2"
	DefaultChangePassword     = "Use2 - request for change your password"
	DefaultSignUpBodyEmail    = "Welcome to <b>Use2</b> %s %s!.<br> Please verify your email by clicking the link below <br> <a href='%s'>Verify now</a> <br> <br> Use2 Team"
	InvalidToken              = "Invalid token %s"
	UserNotActivated          = "User not activated"
	ExpiredToken              = "Token %s expired"
)
