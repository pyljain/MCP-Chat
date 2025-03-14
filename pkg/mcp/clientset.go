package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

type clientSet struct {
	clients         []*client.StdioMCPClient
	toolToClientMap map[string]*client.StdioMCPClient
}

func NewClientSetFromServerConfig(servers []Server) (*clientSet, error) {

	cs := &clientSet{
		clients:         []*client.StdioMCPClient{},
		toolToClientMap: map[string]*client.StdioMCPClient{},
	}
	for _, server := range servers {
		c, err := startServer(&server)
		if err != nil {
			return nil, fmt.Errorf("failed to start server: %v", err)
		}

		cs.clients = append(cs.clients, c)
	}

	return cs, nil
}

func (cs *clientSet) FetchTools() ([]mcp.Tool, error) {
	toolsRequest := mcp.ListToolsRequest{}
	ctx := context.Background()
	toolsSet := make([]mcp.Tool, 0, 10*len(cs.clients))

	for _, c := range cs.clients {
		tools, err := c.ListTools(ctx, toolsRequest)
		if err != nil {
			return nil, err
		}

		toolsSet = append(toolsSet, tools.Tools...)
		for _, t := range tools.Tools {
			cs.toolToClientMap[t.Name] = c
		}
	}

	return toolsSet, nil
}

func (cs *clientSet) CallTool(ctx context.Context, toolName string, args map[string]interface{}) (string, error) {
	client, exists := cs.toolToClientMap[toolName]
	if !exists {
		return "", fmt.Errorf("unknown tool call: %s", toolName)
	}

	req := mcp.CallToolRequest{}
	req.Params.Name = toolName
	req.Params.Arguments = args

	res, err := client.CallTool(ctx, req)
	if err != nil {
		return "", err
	}

	return res.Content[0].(mcp.TextContent).Text, nil
}
