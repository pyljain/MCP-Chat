package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func startServer(server *Server) (*client.StdioMCPClient, error) {
	client, err := client.NewStdioMCPClient(server.Command, []string{}, server.Args...)
	if err != nil {
		return nil, err
	}

	// Initialize the client
	// fmt.Println("Initializing client...")
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "blah",
		Version: "1.0.0",
	}

	ctx := context.Background()
	_, err = client.Initialize(ctx, initRequest)
	if err != nil {
		return nil, err
	}

	return client, nil
}
