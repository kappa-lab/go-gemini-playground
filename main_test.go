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

func Test_QuickStart(t *testing.T) {
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

	// Call the GenerateContent method
	result, err := client.Models.GenerateContent(ctx, "gemini-2.0-flash-exp", genai.Text("Tell me about New York?"), nil)
	require.NoError(t, err)

	text, err := result.Text()
	require.NoError(t, err)

	fmt.Println(text)
}
