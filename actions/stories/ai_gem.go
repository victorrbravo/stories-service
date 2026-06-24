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
	log.Println("imagePath", imagePath)
	// Read the file downloaded from S3 and send it as inline data (base64).
	// genai.Blob.Data is a []byte; encoding/json base64-encodes it automatically
	// when the SDK serializes the request, so we pass the raw bytes here.
	videoBytes, err := os.ReadFile(imagePath)
	if err != nil {
		log.Println("error GemQuestion reading downloaded video file", err)
		return cleanJSON, err
	}

	// 3. Fixed Config: ThinkingLevel is inside ThinkingConfig
	config := &genai.GenerateContentConfig{
		ThinkingConfig: &genai.ThinkingConfig{
			IncludeThoughts: *(genai.Ptr(true)),
			// Use ThinkingLevel (minimal, low, medium, high)
		},
	}

	// 4. GenerateContent: inline video data is passed inside the contents slice.
	// gemini-2.5-flash handles video understanding and is available on the free tier
	// (gemini-2.5-pro typically requires billing).
	resp, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash",
		[]*genai.Content{
			{
				Parts: []*genai.Part{
					{Text: question},
					{InlineData: &genai.Blob{MIMEType: "video/mov", Data: videoBytes}},
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
