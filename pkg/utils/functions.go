package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"stories/types"
	"strconv"
	"strings"
	"text/template"
)

func ParseTemplate[T any](tmpls string, data *T) (string, error) {
	var outputBuffer bytes.Buffer
	var value string
	tmpl, err := template.New("Template").Parse(tmpls)
	if err != nil {
		return value, err
	}
	err = tmpl.Execute(&outputBuffer, data)
	if err != nil {
		return value, err
	}
	value = outputBuffer.String()
	return value, nil
}

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

func SetDefaults[T any](p *T) {
	for i := 0; i < reflect.TypeOf(*p).NumField(); i++ {
		field := reflect.TypeOf(*p).Field(i)
		fieldName := field.Name
		valueOf := reflect.ValueOf(p).Elem()
		valueField := valueOf.FieldByName(fieldName).Interface()
		if value, ok := field.Tag.Lookup("default"); ok {
			currentType := reflect.TypeOf(valueField).String()
			switch currentType {
			case "types.String":
				currentValue := valueField.(types.String)
				currentString := currentValue.String()
				if strings.Compare(currentString, "nil") == 0 {
					valueOf.FieldByName(field.Name).Set(reflect.ValueOf(types.NewString(value)))
				}
			case "int":
				currentValue := valueField.(int)
				if currentValue == 0 {
					valueInt, _ := strconv.Atoi(value)
					valueOf.FieldByName(field.Name).Set(reflect.ValueOf(valueInt))
				}
			}
		}
	}
}

func FormatYaml(input string) string {
	// Define the regular expression pattern
	pattern := `(\r?\n)?\w+:\s*`

	// Compile the regular expression
	regexpPattern := regexp.MustCompile(pattern)

	// Find the matching expression in the input string
	results := regexpPattern.FindAll([]byte(input), -1)
	fReplace := func(x []byte) string {
		pattern := `^\r?\n`
		regexpPattern := regexp.MustCompile(pattern)
		value := regexpPattern.ReplaceAllString(string(x), "")
		value = "\r\n" + value
		if !strings.HasSuffix(value, " ") {
			value = value + " "
		}
		return value
	}
	values := Map(results, fReplace)
	for i, v := range results {
		input = strings.ReplaceAll(input, string(v), values[i])
	}
	input = strings.TrimLeft(input, "\r\n")
	return input
}

func DownloadFromS3(bucketName string, objectKey string) (string, error) {
	// Create a new AWS session
	var value string
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		return value, err
	}

	// Create S3 service client
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "downloaded-file")
	if err != nil {
		return value, err
	}
	defer tempFile.Close()

	// Download the file from S3
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	downloader := s3manager.NewDownloader(sess)

	numBytes, err := downloader.Download(tempFile,
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		})
	if err != nil {
		fmt.Println(err)
	}
	absPath := tempFile.Name()
	log.Println("Downloaded", absPath, numBytes, "bytes")
	return absPath, nil
}

func ToJpeg(jpegPath string, imageBytes []byte) error {

	// DetectContentType detects the content type
	contentType := http.DetectContentType(imageBytes)

	switch contentType {
	case "image/png":
		// Decode the PNG image bytes
		img, err := jpeg.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			return err
		}

		buf := new(bytes.Buffer)

		// encode the image as a JPEG file
		if err := jpeg.Encode(buf, img, nil); err != nil {
			return err
		}
		jpegBytes := buf.Bytes()
		err = os.WriteFile(jpegPath, jpegBytes, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func EncodeImageBase64(path string) (string, error) {
	var base64String string
	imageData, err := ioutil.ReadFile(path)
	if err != nil {
		return base64String, err
	}
	// Convert the image data to base64
	base64String = base64.StdEncoding.EncodeToString(imageData)
	return base64String, nil
}
