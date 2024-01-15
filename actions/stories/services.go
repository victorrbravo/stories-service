package stories

import (
	"context"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
	"os"
	"stories/pkg/utils"
)

func GetTwoWords(imagePath string) (string, error) {
	var err error
	client := openai.NewClient(os.Getenv("OPEN_API_KEY"))

	image64, err := utils.EncodeImageBase64(imagePath)
	if err != nil {
		fmt.Println("Error processing image ", imagePath, " error", err)
		return "", nil
	}

	url := fmt.Sprintf("data:image/jpeg;base64,%s", image64)
	imageURL := openai.ChatMessageImageURL{
		URL: url,
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4VisionPreview,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleUser,
					MultiContent: []openai.ChatMessagePart{
						{
							Type: "text",
							Text: "Summarize the image in two words or less",
						},
						{
							Type:     "image_url",
							ImageURL: &imageURL,
						},
					},
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	value := resp.Choices[0].Message.Content
	return value, nil
}
