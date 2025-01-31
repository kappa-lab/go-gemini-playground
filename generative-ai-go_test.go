package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// see https://github.com/google/generative-ai-go/blob/7276bbc8524c52eed98d38f2169d14dab9d3289f/genai/internal/samples/docs-snippets_test.go#L1752
func printResponse(resp *genai.GenerateContentResponse) {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fmt.Println(part)
			}
		}
	}
	fmt.Println("=========")
}

func genClient_generative_ai_go(t *testing.T) *genai.Client {
	ctx := context.Background()
	require.NoError(t, godotenv.Load(".env"))
	apiKey := os.Getenv("GEMINI_API_KEY")

	client, err := genai.NewClient(
		ctx,
		option.WithAPIKey(apiKey),
	)
	require.NoError(t, err)

	return client
}

func Test_QuickStart_generative_ai_go(t *testing.T) {
	ctx := context.Background()
	require.NoError(t, godotenv.Load(".env"))
	apiKey := os.Getenv("GEMINI_API_KEY")

	client, err := genai.NewClient(
		ctx,
		option.WithAPIKey(apiKey),
	)
	require.NoError(t, err)

	model := client.GenerativeModel("gemini-2.0-flash-exp")
	resp, err := model.GenerateContent(ctx, genai.Text("1年は何日？"))
	require.NoError(t, err)

	printResponse(resp)
}
func Test_GenText_generative_ai_go(t *testing.T) {
	// see https://ai.google.dev/gemini-api/docs/text-generation?hl=ja&lang=go#generate-a-text
	ctx := context.Background()
	client := genClient_generative_ai_go(t)

	model := client.GenerativeModel("gemini-2.0-flash-exp")
	resp, err := model.GenerateContent(ctx, genai.Text("1年は何日？"))
	require.NoError(t, err)

	printResponse(resp)
}

func Test_GenTextStream_generative_ai_go(t *testing.T) {
	// see https://ai.google.dev/gemini-api/docs/text-generation?hl=ja&lang=go#generate-a-text-stream
	ctx := context.Background()
	client := genClient_generative_ai_go(t)

	model := client.GenerativeModel("gemini-2.0-flash-exp")
	iter := model.GenerateContentStream(ctx, genai.Text("1年は何日？"))
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		require.NoError(t, err)

		for _, cand := range resp.Candidates {
			if cand.Content != nil {
				for _, part := range cand.Content.Parts {
					fmt.Print(part)
				}
			}
		}
		// fmt.Print(" , ")//  streamの区切りを可視化
	}
}

func Test_Chat_generative_ai_go(t *testing.T) {
	//https://ai.google.dev/gemini-api/docs/text-generation?hl=ja&lang=go#chat
	ctx := context.Background()
	client := genClient_generative_ai_go(t)

	model := client.GenerativeModel("gemini-2.0-flash-exp")
	cs := model.StartChat()

	cs.History = []*genai.Content{
		{
			Parts: []genai.Part{
				genai.Text("まいど。おおきに。なんでも質問したってや。"),
			},
			Role: "user",
		},
	}

	resp, err := cs.SendMessage(ctx, genai.Text("1年は何日？"))
	require.NoError(t, err)
	printResponse(resp)

	resp, err = cs.SendMessage(ctx, genai.Text("ふーん、vなんで例外があるの？"))
	require.NoError(t, err)
	printResponse(resp)

	fmt.Println("HISTORY")
	for _, v := range cs.History {
		fmt.Println(v)
	}
}

func Test_Instruction_generative_ai_go(t *testing.T) {
	// see https://ai.google.dev/gemini-api/docs/text-generation?hl=ja&lang=go#system-instructions
	ctx := context.Background()
	client := genClient_generative_ai_go(t)

	model := client.GenerativeModel("gemini-2.0-flash-exp")
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text("京都弁で返答してください")},
	}
	resp, err := model.GenerateContent(ctx, genai.Text("1年は何日？"))
	require.NoError(t, err)

	printResponse(resp)
}
