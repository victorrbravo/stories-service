package auth

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"log"
)

func SendEmail(sender, recipient, subject, body string) error {
	// Initialize a new AWS session using the default credentials
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), // Specify your AWS region here
	})
	if err != nil {
		log.Println("Error creating session:", err)
		return err
	}

	// Create a new SES client
	svc := ses.New(sess)

	// Define the email parameters
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(recipient)},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Data: aws.String(body),
				},
			},
			Subject: &ses.Content{
				Data: aws.String(subject),
			},
		},
		Source: aws.String(sender),
	}

	// Try to send the email
	_, err = svc.SendEmail(input)
	if err != nil {
		log.Println("Email sending failed:", err)
		return err
	}
	log.Println("Email sent successfully!")
	return nil

}
