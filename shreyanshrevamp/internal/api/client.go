package api

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/Khan/genqlient/graphql"
)

// NewClient creates a new GraphQL client with authentication, version headers, and request logging
func NewClient(url string, token string, version string) graphql.Client {
	baseTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}

	authTransport := &transportWithAuth{
		Transport: baseTransport,
		Token:     token,
		Version:   version,
	}

	loggingTransport := &loggingRoundTripper{
		rt: authTransport, // wrap auth with logger
	}

	httpClient := &http.Client{
		Transport: loggingTransport,
	}

	client := graphql.NewClient(url, httpClient)
	return client
}

// transportWithAuth adds authentication and version headers
type transportWithAuth struct {
	Transport http.RoundTripper
	Token     string
	Version   string
}

func (t *transportWithAuth) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", t.Token)
	req.Header.Set("x-terraform-provider-version", t.Version)
	return t.Transport.RoundTrip(req)
}

// loggingRoundTripper logs the full request/response body
type loggingRoundTripper struct {
	rt http.RoundTripper
}

func (lrt *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// 🔽 Read and log request body
	var bodyCopy []byte
	if req.Body != nil {
		bodyCopy, _ = io.ReadAll(req.Body)
		log.Println("🔶 GraphQL Request Body:\n", string(bodyCopy))
		req.Body = io.NopCloser(bytes.NewBuffer(bodyCopy)) // reset body for downstream
	}

	// 🔽 Proceed with original request
	resp, err := lrt.rt.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// 🔽 Read and log response body
	if resp.Body != nil {
		respBody, _ := io.ReadAll(resp.Body)
		log.Println("🔷 GraphQL Response Body:\n", string(respBody))
		resp.Body = io.NopCloser(bytes.NewBuffer(respBody)) // reset for client use
	}

	return resp, nil
}
