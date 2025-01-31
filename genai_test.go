package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"google.golang.org/genai"
)

func genClient(t *testing.T) *genai.Client {
	require.NoError(t, godotenv.Load(".env"))
	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")

	client, err := genai.NewClient(
		ctx,
		&genai.ClientConfig{
			APIKey:  apiKey,
			Backend: genai.BackendGeminiAPI,
		},
	)
	require.NoError(t, err)
	return client
}

func Test_QuickStart(t *testing.T) {
	ctx := context.Background()
	client := genClient(t)

	// Call the GenerateContent method
	result, err := client.Models.GenerateContent(ctx, "gemini-2.0-flash-exp", genai.Text("Tell me about New York?"), nil)
	require.NoError(t, err)

	text, err := result.Text()
	require.NoError(t, err)

	fmt.Println(text)
}

func Test_GenText(t *testing.T) {
	ctx := context.Background()
	client := genClient(t)

	result, err := client.Models.GenerateContent(ctx, "gemini-2.0-flash-exp", genai.Text("1年は何日？"), nil)
	require.NoError(t, err)

	text, err := result.Text()
	require.NoError(t, err)

	fmt.Println(text)
}

func Test_GenTextStream(t *testing.T) {
	ctx := context.Background()
	client := genClient(t)

	stream := client.Models.GenerateContentStream(ctx, "gemini-2.0-flash-exp", genai.Text("1年は何日？"), nil)

	for result, err := range stream {
		if err != nil {
			require.NoError(t, err)
		}

		fmt.Print(result.Candidates[0].Content.Parts[0].Text)
	}
}

func Test_Chat(t *testing.T) {
	ctx := context.Background()
	client := genClient(t)

	contents := []*genai.Content{
		{
			Parts: []*genai.Part{
				{
					Text: "まいど。おおきに。なんでも質問したってや。",
				},
			},
			Role: "model",
		},
		{
			Parts: []*genai.Part{
				{
					Text: "1年は何日？",
				},
			},
			Role: "user",
		},
	}
	result, err := client.Models.GenerateContent(ctx, "gemini-2.0-flash-exp", contents, nil)
	require.NoError(t, err)

	text, err := result.Text()
	require.NoError(t, err)

	fmt.Println(text)

	contents = append(contents,
		&genai.Content{
			Parts: []*genai.Part{
				{
					Text: text,
				},
			},
			Role: "model",
		},
		&genai.Content{
			Parts: []*genai.Part{
				{
					Text: "なんで例外があるの？",
				},
			},
			Role: "user",
		},
	)

	result, err = client.Models.GenerateContent(ctx, "gemini-2.0-flash-exp", contents, nil)
	require.NoError(t, err)

	text, err = result.Text()
	require.NoError(t, err)

	fmt.Println(text)

}

func Test_Instruction(t *testing.T) {
	ctx := context.Background()
	client := genClient(t)

	result, err := client.Models.GenerateContent(ctx, "gemini-2.0-flash-exp", genai.Text("1年は何日？"),
		&genai.GenerateContentConfig{
			SystemInstruction: &genai.Content{
				Parts: []*genai.Part{
					{
						Text: "京都弁で返答してください",
					},
				},
			},
		})
	require.NoError(t, err)

	text, err := result.Text()
	require.NoError(t, err)

	fmt.Println(text)
}
