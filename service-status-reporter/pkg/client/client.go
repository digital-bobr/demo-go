package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type RequestBuilder struct {
	method      string
	urlTemplate string
	pathParams  map[string]string
	headers     map[string]string
	params      map[string]string
	body        string
}

func NewRequestBuilder() *RequestBuilder {
	return &RequestBuilder{
		headers:    make(map[string]string),
		params:     make(map[string]string),
		pathParams: make(map[string]string),
	}
}

// SetMethod sets the HTTP method for the request
func (rb *RequestBuilder) SetMethod(m Method) *RequestBuilder {
	rb.method = m.GetValue()
	return rb
}

// SetURLTemplate sets the URL template for the request
// The URL template can contain placeholders for path parameters, e.g., "/api/resource/{id}"
func (rb *RequestBuilder) SetURLTemplate(urlTemplate string) *RequestBuilder {
	rb.urlTemplate = urlTemplate
	return rb
}

// AddPathParam adds a path parameter to the request
// The key should match the placeholder in the URL template, e.g., "id"
func (rb *RequestBuilder) AddPathParam(key, value string) *RequestBuilder {
	rb.pathParams[key] = value
	return rb
}

// AddHeader adds a header to the request
func (rb *RequestBuilder) AddHeader(key, value string) *RequestBuilder {
	rb.headers[key] = value
	return rb
}

// AddParam adds a URL query parameter to the request
func (rb *RequestBuilder) AddParam(key, value string) *RequestBuilder {
	rb.params[key] = value
	return rb
}

// SetBody sets the body of the request
func (rb *RequestBuilder) SetBody(body string) *RequestBuilder {
	rb.body = body
	return rb
}

// Build builds the HTTP request and logs it as a curl command
func (rb *RequestBuilder) Build() (*http.Request, error) {
	// Replace path parameters in the URL template
	finalURL := rb.urlTemplate
	for key, value := range rb.pathParams {
		placeholder := fmt.Sprintf("{%s}", key)
		finalURL = strings.ReplaceAll(finalURL, placeholder, value)
	}

	// Parse the base URL
	parsedURL, err := url.Parse(finalURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %w", err)
	}

	// Add query parameters to the URL if provided
	if len(rb.params) > 0 {
		query := parsedURL.Query()
		for key, value := range rb.params {
			query.Add(key, value)
		}
		parsedURL.RawQuery = query.Encode()
	}

	// Create the HTTP request
	var req *http.Request
	if rb.method == http.MethodPost || rb.method == http.MethodPut || rb.method == http.MethodPatch {
		req, err = http.NewRequest(rb.method, parsedURL.String(), bytes.NewBufferString(rb.body))
	} else {
		req, err = http.NewRequest(rb.method, parsedURL.String(), nil)
	}

	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Add headers if provided
	for key, value := range rb.headers {
		req.Header.Set(key, value)
	}

	// Log the request as a curl command
	rb.logCurlCommand(parsedURL.String())

	return req, nil
}

func (rb *RequestBuilder) logCurlCommand(finalURL string) {
	curlCmd := fmt.Sprintf("curl -X %s '%s'", rb.method, finalURL)

	for key, value := range rb.headers {
		curlCmd += fmt.Sprintf(" -H '%s: %s'", key, value)
	}

	if rb.method != http.MethodGet && rb.method != http.MethodDelete {
		if rb.body != "" {
			curlCmd += fmt.Sprintf(" -d '%s'", rb.body)
		}
	}

	//fmt.Println("Generated curl command:")
	//fmt.Println(curlCmd)
}

func SendRequest(req *http.Request) (int, []byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("error reading response body: %w", err)
	}

	//fmt.Println("Response:")
	//fmt.Println(string(body))

	return resp.StatusCode, body, nil
}
