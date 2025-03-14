package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"merlin/pkg/mcp"
	"merlin/pkg/utils"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
)

func main() {
	// Find MCP config in ~/.merlin/mcp.json
	cfg := mcp.LoadConfig()

	ctx := context.Background()

	// Start MCP Servers
	clientSet, err := mcp.NewClientSetFromServerConfig(cfg.Servers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create MCP client %s", err)
		os.Exit(-1)
	}

	// Discover Tools
	tools, err := clientSet.FetchTools()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not fetch MCP tools %s", err)
		os.Exit(-1)
	}

	llm, err := anthropic.New(
		anthropic.WithModel("claude-3-5-sonnet-20240620"),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not initiate LLM  %s", err)
		os.Exit(-1)
	}

	history := []llms.MessageContent{}

	toolsCalled := false

	for {
		// Prompt user for question
		if !toolsCalled {
			fmt.Printf("> ")
			reader := bufio.NewReader(os.Stdin)
			userInput, _ := reader.ReadString('\n')
			if userInput == "/exit\n" {
				os.Exit(0)
			}

			if userInput == "/tools\n" {
				for i, t := range tools {
					fmt.Printf("%d. %s\n", i+1, t.Name)
				}
				continue
			}

			history = append(history, llms.TextParts(llms.ChatMessageTypeHuman, userInput))
		}

		schema := utils.GenerateLangChainToolsFromMCPTools(tools)
		res, err := llm.GenerateContent(
			ctx, history, llms.WithTools(schema),
			llms.WithMaxTokens(8192),
			llms.WithFunctionCallBehavior(llms.FunctionCallBehavior("auto")))
		if err != nil {
			log.Printf("Error calling LLM %s", err)
			os.Exit(-1)
		}

		toolsCalled = false

		for _, choice := range res.Choices {

			fmt.Printf("Agent says: %v\n", choice.Content)

			if choice.ToolCalls == nil {
				continue
			}

			toolsCalled = true
			for _, toolCall := range choice.ToolCalls {
				// Call tool
				args := map[string]interface{}{}
				err := json.Unmarshal([]byte(toolCall.FunctionCall.Arguments), &args)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error unmarshalling args  %s", err)
					os.Exit(-1)
				}

				fnResponse, err := clientSet.CallTool(ctx, toolCall.FunctionCall.Name, args)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error calling tool %s", err)
					os.Exit(-1)
				}
				// fmt.Printf("Tool call response %s\n", fnResponse)

				// Append tool_use to messageHistory
				assistantResponse := llms.MessageContent{
					Role: llms.ChatMessageTypeAI,
					Parts: []llms.ContentPart{
						llms.ToolCall{
							ID:   toolCall.ID,
							Type: toolCall.Type,
							FunctionCall: &llms.FunctionCall{
								Name:      toolCall.FunctionCall.Name,
								Arguments: toolCall.FunctionCall.Arguments,
							},
						},
					},
				}

				history = append(history, assistantResponse)

				// Append tool_result to history
				fnResponseInHistory := llms.MessageContent{
					Role: llms.ChatMessageTypeTool,
					Parts: []llms.ContentPart{
						llms.ToolCallResponse{
							ToolCallID: toolCall.ID,
							Name:       toolCall.FunctionCall.Name,
							Content:    fnResponse,
						},
					},
				}
				history = append(history, fnResponseInHistory)
			}
		}
	}
	// If question is /exit, exit
	// If questions is /install then install MCP server

	// Send question to LLM with MCP discovered tools

	// Get response and call tools on MCP server
}
