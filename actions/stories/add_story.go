package stories

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gopkg.in/yaml.v3"
	"log"
	"stories/actions/constants"
	"stories/database"
	"stories/generated/models"
	"stories/pkg/utils"
	"stories/types"
	"time"
)

const (
	awsRegion       = "us-east-1"
	bucketName      = "mystoriesbucket"
	objectKey       = "icecream"
	durationInSec   = 3600 // Duration for which the URL will be valid (in seconds)
	JsonUrlTemplate = `
	{ "url": "%s" }
`
)

type Ticket struct {
	Name    types.String `yaml:"name"`
	Summary types.String `yaml:"summary"`
}

type AddStoryStatus struct {
	Ticket   types.Int    `yaml:"ticket"`
	Status   types.String `yaml:"status"`
	Amount   types.String `yaml:"amount"`
	Image1   types.String `yaml:"image"`
	StepType types.String `yaml:"type"`
}

type DeleteTicket struct {
	Ticket types.Int `yaml:"ticket"`
}

func FlowProcessor(action string, order string, query *database.IQuery) (string, error) {
	var err error
	var value string
	ctx := context.Background()
	switch action {
	case "steps_schema":
		value, err = DoGetStepsSchema(order, query)
	case "new_story_flow":
		err = NewStoryFlow(order, query)
	case "add_step":
		_, err = DoAddStep(order, query)
	case "add_flow_step":
		err = DoAddStepFlowAtEnd[models.TaTracking](ctx, order, query)
	case "delete_ticket":
		err = DoDeleteStep(order, query)
	case "presign_image":
		value, err = DoPresignImage(order, query)
	case "get_presigned_image":
		value, err = DoGetPresignedImage(order, query)

	}
	if err != nil {
		return value, err
	}
	return value, nil
}

func DoGetStepsSchema(order string, query *database.IQuery) (string, error) {
	var value = `
 {
	"actions" : {
		"new_story_flow": {
			"description" : "create a new flow",
			"fields": [
				{ "name": "name", "type": "string", "validator" : "" },
				{ "name": "summary", "type": "string", "validator" : "" }
			]
		},
		"add_step": {
		"description" : "add step",
			"fields": [
				{ "name": "name", "type": "string", "validator" : "" },
				{ "name": "summary", "type": "string", "validator" : "" }
			]
		},
		"add_flow_step": {
			"description" : "add flow step",
			"fields": [
				{ "name": "name", "type": "string", "validator" : "" },
				{ "name": "summary", "type": "string", "validator" : "" }
			]
		},
		"delete_ticket": {
			"description" : "delete",
			"fields": [
				{ "name": "ticket", "type": "string", "validator" : "" },
				{ "name": "summary", "type": "string", "validator" : "" }
			]
		},
		"presign_image": {
		"description" : "presign",
			"fields": [
				{ "name": "image", "type": "string", "validator" : "" },
				{ "name": "summary", "type": "string", "validator" : "" }
			]
		},
		"get_presigned_image" : {
			"description" : "get presigned image",
			"fields": [
				{ "name": "image", "type": "string", "validator" : "" },
				{ "name": "summary", "type": "string", "validator" : "" }
			]
	}
 	}
}
`

	return value, nil
}

func DoGetPresignedImage(order string, query *database.IQuery) (string, error) {
	keyName := "icecream"
	fmt.Println("DoGetPresignedImage keyName", keyName)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		return "", err
	}
	// Create an S3 service client
	svc := s3.New(sess)

	// Specify the duration for which the URL should be valid (e.g., 1 hour)
	duration := time.Hour

	// Generate the pre-signed URL
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
	})
	url, err := req.Presign(duration)
	if err != nil {
		return "", err
	}
	value := fmt.Sprintf(JsonUrlTemplate, url)
	return value, nil
}

func DoPresignImage(order string, query *database.IQuery) (string, error) {
	log.Println("DoPresignImage executing")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		fmt.Println("Error creating session:", err)
		return "", err
	}

	// Create an S3 client
	s3Client := s3.New(sess)

	// Generate a pre-signed URL for uploading the file
	presignedURL, err := generatePresignedURL(s3Client, bucketName, objectKey, durationInSec)
	if err != nil {
		fmt.Println("Error generating pre-signed URL:", err)
		return "", err
	}
	value := fmt.Sprintf(JsonUrlTemplate, presignedURL)
	return value, nil
}

// generatePresignedURL generates a pre-signed URL for uploading a file to S3.
func generatePresignedURL(s3Client *s3.S3, bucket, key string, seconds int) (string, error) {
	req, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	// Set the expiration time for the pre-signed URL
	duration, err := time.ParseDuration(fmt.Sprintf("%ds", seconds))
	if err != nil {
		return "", err
	}
	url, err := req.Presign(duration)
	if err != nil {
		return "", err
	}

	return url, nil
}

func UnmarshallTicket(order []byte) (Ticket, error) {
	var ticket Ticket
	var err error
	err = yaml.Unmarshal(order, &ticket)
	if err != nil {
		return ticket, err
	}

	return ticket, nil
}

func UnmarshallAddStep(order []byte, query *database.IQuery) (AddStoryStatus, error) {
	var change AddStoryStatus
	var err error
	err = yaml.Unmarshal(order, &change)
	if err != nil {
		return change, err
	}
	return change, nil
}

func NewStoryFlow(order string, query *database.IQuery) error {
	handler, err := UnmarshallTicket([]byte(order))
	if err != nil {
		return err
	}
	start := types.NewString(constants.Start)
	images := types.NewStringArray([]string{})
	roles := types.NewStringArray([]string{})
	taStory, err := models.CreateTaStory(*query,
		handler.Name,
		start,
		handler.Summary,
	)
	if err != nil {
		return err
	}
	stringArr := []string{}
	models.CreateTaTracking(*query,
		taStory.PkStory,
		start,
		types.NilString(),
		types.NewStringArray(stringArr),
		images,
		roles,
	)
	return nil
}

func DoAddStepFlowAtEnd[T database.Scanner](ctx context.Context, order string, query *database.IQuery) error {
	addStep, err := UnmarshallAddStep([]byte(order), query)
	if err != nil {
		return err
	}
	var addTask AddTask
	wf, err := GetWorkflow[T](1, query)
	if err != nil {
		return err
	}
	addTask.Name = types.NewString(addStep.Status.String())
	lastTask, err := wf.Last(ctx)
	if err != nil {
		return err
	}
	addTask.PkFlow = 1
	addTask.Source = types.NewString(lastTask.Name)
	addTask.Sentence = types.NewString(wf.Sentence(addTask.Name.String()))
	utils.SetDefaults(&addTask)
	wf, err = DoAddTaskIntl[T](addTask, query)
	if err != nil {
		return err
	}
	err = DoAddStepIntl(addStep, query)
	if err != nil {
		return err
	}
	return nil
}

func DoAddStepFlow(order string, query *database.IQuery) error {
	handler, err := UnmarshallAddStep([]byte(order), query)
	if err != nil {
		return err
	}

	taStory, err := models.GetTaStoryByPkStory(*query,
		handler.Ticket)
	if err != nil {
		return err
	}

	images := types.NewStringArray([]string{handler.Image1.String()})
	roles := types.NewStringArray([]string{handler.Amount.String(), handler.StepType.String()})
	status := types.NewString(handler.Status.String())
	taStory.Update(*query,
		types.NilString(),
		status,
		types.NilString())

	models.CreateTaTracking(*query,
		taStory.PkStory,
		status,
		types.NilString(),
		types.NewStringArray([]string{}),
		images,
		roles,
	)

	return nil
}

func DoAddStepIntl(status AddStoryStatus, query *database.IQuery) error {
	taStory, err := models.GetTaStoryByPkStory(*query,
		status.Ticket)
	if err != nil {
		return err
	}

	images := types.NewStringArray([]string{status.Image1.String()})
	roles := types.NewStringArray([]string{status.Amount.String(), status.StepType.String()})
	statusStr := types.NewString(status.Status.String())
	taStory.Update(*query,
		types.NilString(),
		statusStr,
		types.NilString())

	models.CreateTaTracking(*query,
		taStory.PkStory,
		statusStr,
		types.NilString(),
		types.NewStringArray([]string{}),
		images,
		roles,
	)

	return nil
}

func DoAddStep(order string, query *database.IQuery) (AddStoryStatus, error) {
	handler, err := UnmarshallAddStep([]byte(order), query)
	if err != nil {
		return handler, err
	}
	err = DoAddStepIntl(handler, query)
	if err != nil {
		return handler, err
	}
	return handler, nil
}

func DoDeleteStep(order string, query *database.IQuery) error {
	handler, err := UnmarshallAddStep([]byte(order), query)
	if err != nil {
		return err
	}
	err = models.DeleteTaTrackingByFkStory(*query, handler.Ticket)
	if err != nil {
		return err
	}
	taStory, err := models.GetTaStoryByPkStory(*query, handler.Ticket)
	if err != nil {
		return err
	}
	err = taStory.Delete(*query)
	if err != nil {
		return err
	}
	return nil
}
