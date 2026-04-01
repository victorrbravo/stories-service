package stories

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"stories/database"
	"stories/pkg/utils"
	"stories/types"
	"strings"

	"google.golang.org/genai"
)

type GeminiResponse struct {
	Summary         string   `json:"summary"`
	AssociatedWords []string `json:"associated_words"`
}

func uploadVideoFile(ctx context.Context, videoPath string, client genai.Client) (*genai.File, error) {
	fmt.Println("Uploading video...")

	fileOptions := genai.UploadFileConfig{
		MIMEType: "video/mov",
	}
	// Third argument is the config (can be nil)
	file, err := client.Files.UploadFromPath(ctx, videoPath, &fileOptions)
	if err != nil {
		log.Println("Error Google AI gemini err:", err)
		return nil, err
	}
	return file, err
}

func GemQuestion(order string, query *database.IQuery, question string) (string, error) {
	var cleanJSON string
	ctx := context.Background()
	handler, err := Unmarshall[SuggestAction]([]byte(order))
	if err != nil {
		log.Println("error GemQuestion unmarshalling")
		return cleanJSON, err
	}
	if !types.IsNil(handler.Question) {
		question = handler.Question.String()
	}
	if handler.ImageID.IsEmpty() {
		log.Println("image id is empty")
		return cleanJSON, err
	}
	imagePath, err := utils.DownloadFromS3(bucketName, handler.ImageID.String())
	if err != nil {
		log.Println("error GemQuestion download from S3 error", err)
		return cleanJSON, err
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GOOGLE_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return cleanJSON, err
	}
	var file *genai.File
	log.Println("imagePath", imagePath)
	file, err = uploadVideoFile(ctx, imagePath, *client)
	if err != nil {
		log.Println("err uploading video", err)
		return cleanJSON, err
	}
	fileName := file.Name
	// 2. Wait for processing
	fmt.Println("file.Name", fileName)
	fmt.Print("Processing video...")
	for {
		file, err = client.Files.Get(ctx, fileName, nil)
		if err != nil {
			log.Println("error GemQuestion processing video", err)
			return cleanJSON, err
		}
		if file.State == genai.FileStateActive {
			break
		}
		fmt.Print(".")
	}
	fmt.Println("\nVideo is ready!")

	// 3. Fixed Config: ThinkingLevel is inside ThinkingConfig
	config := &genai.GenerateContentConfig{
		ThinkingConfig: &genai.ThinkingConfig{
			IncludeThoughts: *(genai.Ptr(true)),
			// Use ThinkingLevel (minimal, low, medium, high)
		},
	}

	// 4. Fixed GenerateContent: FileData is passed inside the contents slice
	// Note: We use "gemini-3-flash-preview" as requested
	resp, err := client.Models.GenerateContent(ctx, "gemini-3-flash-preview",
		[]*genai.Content{
			{
				Parts: []*genai.Part{
					{Text: question},
					{FileData: &genai.FileData{FileURI: file.URI}}, // Fixed FileURI field name
				},
			},
		},
		config,
	)
	if err != nil {
		log.Println("error Generating GemQuestion processing video", err)
		return cleanJSON, err
	}
	jsonText := resp.Text()
	var geminiResponse GeminiResponse
	cleanJSON = strings.TrimPrefix(string(jsonText), "```json")
	cleanJSON = strings.TrimSuffix(cleanJSON, "```")
	cleanJSON = strings.TrimSpace(cleanJSON)
	err = json.Unmarshal([]byte(cleanJSON), &geminiResponse)
	if err != nil {
		return cleanJSON, err
	}

	fmt.Println("--- Analysis Results ---")
	fmt.Println(geminiResponse)
	return cleanJSON, err

}
