package deepseek

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	utils "github.com/cohesion-org/deepseek-go/utils"
)

// ImageContent represents the content of an image in the chat completion request.
type ImageContent struct {
	URL any `json:"url,omitempty"` // URL of the image (can be a string or a base64 encoded string)
}

// ContentItem represents a single content item in the chat completion request.
type ContentItem struct {
	Type  string        `json:"type"`                // Type of content (e.g., "text", "image_url")
	Text  string        `json:"text,omitempty"`      // Text content (if type is "text")
	Image *ImageContent `json:"image_url,omitempty"` // Image content (if type is "image_url")
}

// ChatCompletionMessageWithImage represents a message in the chat completion request. It's not a deepseek feature, it's added to support images (files in the future).
type ChatCompletionMessageWithImage struct {
	Role             string     `json:"role"`                        // Role of the message sender (e.g., "user", "assistant")
	Content          any        `json:"content"`                     // Can accept both string and []ContentItem
	Prefix           bool       `json:"prefix,omitempty"`            // Whether to prefix the message
	ReasoningContent string     `json:"reasoning_content,omitempty"` // Reasoning content
	ToolCallID       string     `json:"tool_call_id,omitempty"`      // ID of the tool call
	ToolCalls        []ToolCall `json:"tool_calls,omitempty"`        // List of tool calls
}

// ChatCompletionRequest defines the structure for a chat completion request.
type ChatCompletionRequestWithImage struct {
	Model            string                           `json:"model"`                       // The ID of the model to use (required).
	Messages         []ChatCompletionMessageWithImage `json:"messages"`                    // A list of messages comprising the conversation (required).
	FrequencyPenalty float32                          `json:"frequency_penalty,omitempty"` // Penalty for new tokens based on their frequency in the text so far (optional).
	MaxTokens        int                              `json:"max_tokens,omitempty"`        // The maximum number of tokens to generate in the chat completion (optional).
	PresencePenalty  float32                          `json:"presence_penalty,omitempty"`  // Penalty for new tokens based on their presence in the text so far (optional).
	Temperature      float32                          `json:"temperature,omitempty"`       // The sampling temperature, between 0 and 2 (optional).
	TopP             float32                          `json:"top_p,omitempty"`             // The nucleus sampling parameter, between 0 and 1 (optional).
	ResponseFormat   *ResponseFormat                  `json:"response_format,omitempty"`   // The desired response format (optional).
	Stop             []string                         `json:"stop,omitempty"`              // A list of sequences where the model should stop generating further tokens (optional).
	Tools            []Tool                           `json:"tools,omitempty"`             // A list of tools the model may use (optional).
	ToolChoice       interface{}                      `json:"tool_choice,omitempty"`       // Controls which (if any) tool is called by the model (optional).
	LogProbs         bool                             `json:"logprobs,omitempty"`          // Whether to return log probabilities of the most likely tokens (optional).
	TopLogProbs      int                              `json:"top_logprobs,omitempty"`      // The number of top most likely tokens to return log probabilities for (optional).
	JSONMode         bool                             `json:"json,omitempty"`              // [deepseek-go feature] Optional: Enable JSON mode. If you're using the JSON mode, please mention "json" anywhere in your prompt, and also include the JSON schema in the request.
}

// StreamChatCompletionRequestWithImage represents the request body for a streaming chat completion API call with image support.
type StreamChatCompletionRequestWithImage struct {
	Stream           bool                             `json:"stream,omitempty"`            //Comments: Defaults to true, since it's "STREAM"
	StreamOptions    StreamOptions                    `json:"stream_options,omitempty"`    // Optional: Stream options for the request.
	Model            string                           `json:"model"`                       // Required: Model ID, e.g., "deepseek-chat"
	Messages         []ChatCompletionMessageWithImage `json:"messages"`                    // Required: List of messages
	FrequencyPenalty float32                          `json:"frequency_penalty,omitempty"` // Optional: Frequency penalty, >= -2 and <= 2
	MaxTokens        int                              `json:"max_tokens,omitempty"`        // Optional: Maximum tokens, > 1
	PresencePenalty  float32                          `json:"presence_penalty,omitempty"`  // Optional: Presence penalty, >= -2 and <= 2
	Temperature      float32                          `json:"temperature,omitempty"`       // Optional: Sampling temperature, <= 2
	TopP             float32                          `json:"top_p,omitempty"`             // Optional: Nucleus sampling parameter, <= 1
	ResponseFormat   *ResponseFormat                  `json:"response_format,omitempty"`   // Optional: Custom response format: just don't try, it breaks rn ;)
	Stop             []string                         `json:"stop,omitempty"`              // Optional: Stop signals
	Tools            []Tool                           `json:"tools,omitempty"`             // Optional: List of tools
	LogProbs         bool                             `json:"logprobs,omitempty"`          // Optional: Enable log probabilities
	TopLogProbs      int                              `json:"top_logprobs,omitempty"`      // Optional: Number of top tokens with log probabilities, <= 20
}

// CreateChatCompletion sends a chat completion request and returns the generated response.
func (c *Client) CreateChatCompletionWithImage(
	ctx context.Context,
	request *ChatCompletionRequestWithImage,
) (*ChatCompletionResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	ctx, tcancel, err := getTimeoutContext(ctx, c.Timeout)
	if err != nil {
		return nil, err
	}
	defer tcancel()

	req, err := utils.NewRequestBuilder(c.AuthToken).
		SetBaseURL(c.BaseURL).
		SetPath(c.Path).
		SetBodyFromStruct(request).
		Build(ctx)

	if err != nil {
		return nil, fmt.Errorf("error building request: %w", err)
	}
	resp, err := HandleSendChatCompletionRequest(*c, req)

	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, HandleAPIError(resp)
	}

	updatedResp, err := HandleChatCompletionResponse(resp)

	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return updatedResp, err
}

// CreateChatCompletionStream sends a chat completion request with stream = true and returns the delta
func (c *Client) CreateChatCompletionStreamWithImage(
	ctx context.Context,
	request *StreamChatCompletionRequestWithImage,
) (ChatCompletionStream, error) {
	if request == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	ctx, _, err := getTimeoutContext(ctx, c.Timeout)
	if err != nil {
		return nil, fmt.Errorf("error getting timeout context: %w", err)
	}

	request.Stream = true
	req, err := utils.NewRequestBuilder(c.AuthToken).
		SetBaseURL(c.BaseURL).
		SetPath(c.Path).
		SetBodyFromStruct(request).
		BuildStream(ctx)

	if err != nil {
		return nil, fmt.Errorf("error building request: %w", err)
	}

	resp, err := HandleSendChatCompletionRequest(*c, req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, HandleAPIError(resp)
	}

	ctx, cancel := context.WithCancel(ctx)
	stream := &chatCompletionStream{
		ctx:    ctx,
		cancel: cancel,
		resp:   resp,
		reader: bufio.NewReader(resp.Body),
	}
	return stream, nil
}

// NewImageMessage creates a new image message for the chat completion request.
func NewImageMessage(role string, text string, imageURL string) ChatCompletionMessageWithImage {

	return ChatCompletionMessageWithImage{
		Role: role,
		Content: []ContentItem{
			{
				Type: "text",
				Text: text,
			},
			{
				Type: "image_url",
				Image: &ImageContent{
					URL: imageURL,
				},
			},
		},
	}
}

// ImageToBase64 converts an image URL to a base64 encoded string.
func ImageToBase64(imageURL string) (string, error) {
	if imageURL == "" {
		return "", fmt.Errorf("imageURL cannot be empty")
	}
	ext := strings.ToLower(filepath.Ext(imageURL))

	validExtensions := map[string]bool{
		".png":  true,
		".jpg":  true,
		".jpeg": true,
		".webp": true,
	}
	if !validExtensions[ext] {
		return "", fmt.Errorf("unsupported image format: %s", ext)
	}

	// Check if the imageURL is a web URL by looking for http(s) prefix
	if strings.HasPrefix(imageURL, "http://") || strings.HasPrefix(imageURL, "https://") {
		return handleImageFromURL(imageURL)
	}

	// Read and encode the file
	imgData, err := os.ReadFile(imageURL)

	if err != nil {
		return "", fmt.Errorf("failed to read image file: %w", err)
	}

	base64Str := base64.StdEncoding.EncodeToString(imgData)
	contentType := createContentType(ext, base64Str)
	if contentType == "" {
		return "", fmt.Errorf("unsupported image format: %s", ext)
	}
	return contentType, nil
}

// handleImageFromURL downloads an image from a URL and converts it to a base64 encoded string.
func handleImageFromURL(url string) (string, error) {
	resp, err := http.Get(url)

	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image: %s", resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("invalid content type: %s", contentType)
	}

	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read image data: %w", err)
	}

	base64Str := base64.StdEncoding.EncodeToString(imgData)
	ext := strings.ToLower(filepath.Ext(url))
	contentType = createContentType(ext, base64Str)
	if contentType == "" {
		return "", fmt.Errorf("unsupported image format: %s", ext)
	}
	return contentType, nil
}

func createContentType(ext string, base64 string) string {
	switch ext {
	case ".png":
		return fmt.Sprintf("data:image/png;base64,%s", base64)
	case ".jpg", ".jpeg":
		return fmt.Sprintf("data:image/jpeg;base64,%s", base64)
	case ".webp":
		return fmt.Sprintf("data:image/webp;base64,%s", base64)
	default:
		return ""
	}

}
