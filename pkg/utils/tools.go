package utils

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/tmc/langchaingo/llms"
)

func GenerateLangChainToolsFromMCPTools(mcpTools []mcp.Tool) []llms.Tool {
	tools := make([]llms.Tool, 0, len(mcpTools))

	for _, mt := range mcpTools {
		tools = append(tools, llms.Tool{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        mt.Name,
				Description: mt.Description,
				Parameters:  mt.InputSchema,
			},
		})
	}

	return tools
}
