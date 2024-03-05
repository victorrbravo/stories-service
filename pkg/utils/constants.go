package utils

const (
	awsRegion       = "us-east-1"
	bucketName      = "mystoriesbucket"
	durationInSec   = 3600 // Duration for which the URL will be valid (in seconds)
	JsonUrlTemplate = `
	{ "url": "%s" }
`
)
