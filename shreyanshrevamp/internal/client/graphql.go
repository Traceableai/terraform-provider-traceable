package client

import (
	"context"
)

// Client is an interface for making GraphQL requests
type Client interface {
	Query(ctx context.Context, query string, variables map[string]interface{}, response interface{}) error
	Mutate(ctx context.Context, mutation string, variables map[string]interface{}, response interface{}) error
}

// GraphQLClient implements the Client interface
type GraphQLClient struct {
	endpoint string
	token    string
}

// NewGraphQLClient creates a new GraphQL client
func NewGraphQLClient(endpoint, token string) *GraphQLClient {
	return &GraphQLClient{
		endpoint: endpoint,
		token:    token,
	}
}

// Query executes a GraphQL query
func (c *GraphQLClient) Query(ctx context.Context, query string, variables map[string]interface{}, response interface{}) error {
	// TODO: Implement actual GraphQL query logic
	return nil
}

// Mutate executes a GraphQL mutation
func (c *GraphQLClient) Mutate(ctx context.Context, mutation string, variables map[string]interface{}, response interface{}) error {
	// TODO: Implement actual GraphQL mutation logic
	return nil
}
