package introspection

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ClientOptions configures the introspection client.
type ClientOptions struct {
	Headers map[string]string
	Timeout time.Duration
}

// DefaultClientOptions returns default client configuration.
func DefaultClientOptions() ClientOptions {
	return ClientOptions{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Timeout: 30 * time.Second,
	}
}

// IsURL checks if the argument is a URL (http:// or https://).
func IsURL(arg string) bool {
	return strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://")
}

// FetchSchema fetches a GraphQL schema via introspection from the given endpoint.
func FetchSchema(ctx context.Context, endpoint string, opts ClientOptions) (*Response, error) {
	// Build request body
	reqBody := map[string]string{
		"query": Query,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range opts.Headers {
		req.Header.Set(key, value)
	}

	// Execute request
	client := &http.Client{Timeout: opts.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var introspectionResp Response
	if err := json.Unmarshal(body, &introspectionResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for GraphQL errors
	if len(introspectionResp.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", introspectionResp.Errors[0].Message)
	}

	// Validate response has data
	if introspectionResp.Data == nil {
		return nil, fmt.Errorf("introspection response missing data")
	}

	return &introspectionResp, nil
}

// ParseHeaders parses header strings in "Key: Value" format.
func ParseHeaders(headers []string) (map[string]string, error) {
	result := make(map[string]string)
	for _, h := range headers {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid header format '%s': expected 'Key: Value'", h)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		result[key] = value
	}
	return result, nil
}
