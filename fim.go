package deepseek

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	ConstFIMMaxTokensMin     = 1
	ConstFIMMaxTokensDefault = 4000
	ConstFIMMaxTokensMax     = 4000

	ConstFIMPresencePenaltyMin     = -2
	ConstFIMPresencePenaltyDefault = 0
	ConstFIMPresencePenaltyMax     = 2

	ConstFIMFrequencyPenaltyMin     = -2
	ConstFIMFrequencyPenaltyDefault = 0
	ConstFIMFrequencyPenaltyMax     = 2

	ConstFIMTemperatureDefault = 1
	ConstFIMTemperatureMax     = 2

	ConstFIMTopPDefault = 1
	ConstFIMTopPMax     = 1

	ConstFIMLogprobsMax = 20
)

// FIMCompletionRequest represents the request body for a Fill-In-the-Middle (FIM) completion.
type FIMCompletionRequest struct {
	Model            string   `json:"model"`                       // Model name to use for completion.
	Prompt           string   `json:"prompt"`                      // The prompt to start the completion from.
	Suffix           string   `json:"suffix,omitempty"`            // Optional: The suffix to complete the prompt with.
	MaxTokens        int      `json:"max_tokens,omitempty"`        // Optional: Maximum tokens to generate, > 1 and <= 4000.
	Temperature      float64  `json:"temperature,omitempty"`       // Optional: Sampling temperature, between 0 and 2.
	TopP             float64  `json:"top_p,omitempty"`             // Optional: Nucleus sampling probability threshold.
	N                int      `json:"n,omitempty"`                 // Optional: Number of completions to generate.
	Logprobs         int      `json:"logprobs,omitempty"`          // Optional: Number of log probabilities to return.
	Echo             bool     `json:"echo,omitempty"`              // Optional: Whether to echo the prompt in the completion.
	Stop             []string `json:"stop,omitempty"`              // Optional: List of stop sequences.
	PresencePenalty  float64  `json:"presence_penalty,omitempty"`  // Optional: Penalty for new tokens based on their presence in the text so far.
	FrequencyPenalty float64  `json:"frequency_penalty,omitempty"` // Optional: Penalty for new tokens based on their frequency in the text so far.
}

type FIMCompletionRequestOption func(*FIMCompletionRequest)

func NewDefaultFIMCompletionRequest(model string, prompt string, opts ...FIMCompletionRequestOption) *FIMCompletionRequest {
	locFIMCompletionReq := &FIMCompletionRequest{
		Model:            model,
		Prompt:           prompt,
		MaxTokens:        ConstFIMMaxTokensDefault,
		Temperature:      ConstFIMTemperatureDefault,
		TopP:             ConstFIMTopPDefault,
		PresencePenalty:  ConstFIMPresencePenaltyDefault,
		FrequencyPenalty: ConstFIMFrequencyPenaltyDefault,
	}
	for _, o := range opts {
		o(locFIMCompletionReq)
	}
	return locFIMCompletionReq
}

func WithModel(model string, prompt string) FIMCompletionRequestOption {
	return func(r *FIMCompletionRequest) {
		r.Model = model
		r.Prompt = prompt
	}
}

func WithMaxTokens(maxTokens int) FIMCompletionRequestOption {
	return func(r *FIMCompletionRequest) {
		r.MaxTokens = maxTokens
	}
}

func WithTemperature(temperature float64) FIMCompletionRequestOption {
	return func(r *FIMCompletionRequest) {
		r.Temperature = temperature
	}
}

func WithTopP(topP float64) FIMCompletionRequestOption {
	return func(r *FIMCompletionRequest) {
		r.TopP = topP
	}
}

func WithPresencePenalty(presencePenalty float64) FIMCompletionRequestOption {
	return func(r *FIMCompletionRequest) {
		r.PresencePenalty = presencePenalty
	}
}

func WithFrequencyPenalty(frequencyPenalty float64) FIMCompletionRequestOption {
	return func(r *FIMCompletionRequest) {
		r.FrequencyPenalty = frequencyPenalty
	}
}

func WithEcho(echo bool) FIMCompletionRequestOption {
	return func(r *FIMCompletionRequest) {
		r.Echo = echo
	}
}

func WithStop(stop []string) FIMCompletionRequestOption {
	return func(r *FIMCompletionRequest) {
		r.Stop = stop
	}
}

func WithLogprobs(logprobs int) FIMCompletionRequestOption {
	return func(r *FIMCompletionRequest) {
		r.Logprobs = logprobs
	}
}

func WithSuffix(suffix string) FIMCompletionRequestOption {
	return func(r *FIMCompletionRequest) {
		r.Suffix = suffix
	}
}

func CheckFIMCompletionRequest(request *FIMCompletionRequest) error {
	if request == nil {
		return fmt.Errorf("request cannot be nil")
	}
	if request.Model == "" {
		return fmt.Errorf("model cannot be empty")
	}
	if request.Prompt == "" {
		return fmt.Errorf("prompt cannot be empty")
	}
	if request.MaxTokens <= ConstFIMMaxTokensMin {
		return fmt.Errorf("max tokens must be > 1")
	}
	if request.MaxTokens > ConstFIMMaxTokensMax {
		return fmt.Errorf("max tokens must be <= 4000")
	}
	if request.FrequencyPenalty < ConstFIMFrequencyPenaltyMin {
		return fmt.Errorf("frequency penalty must be >= %v", ConstFIMFrequencyPenaltyMin)
	}
	if request.FrequencyPenalty > ConstFIMFrequencyPenaltyMax {
		return fmt.Errorf("frequency penalty must be < =%v", ConstFIMFrequencyPenaltyMax)
	}
	if request.Logprobs > ConstFIMLogprobsMax {
		return fmt.Errorf("logprobs must be <= %v", ConstFIMLogprobsMax)
	}
	if request.Temperature > ConstFIMTemperatureMax {
		return fmt.Errorf("temperature must be <= %v", ConstFIMTemperatureMax)
	}
	if request.TopP > ConstFIMTopPMax {
		return fmt.Errorf("top p must be <= %v", ConstFIMTopPMax)
	}

	return nil
}

// FIMCompletionResponse represents the response body for a Fill-In-the-Middle (FIM) completion.
type FIMCompletionResponse struct {
	ID      string `json:"id"`      // Unique ID for the completion.
	Object  string `json:"object"`  // The object type, e.g., "text_completion".
	Created int    `json:"created"` // Timestamp of when the completion was created.
	Model   string `json:"model"`   // Model used for the completion.
	Choices []struct {
		Text         string   `json:"text"`          // The generated completion text.
		Index        int      `json:"index"`         // Index of the choice.
		Logprobs     Logprobs `json:"logprobs"`      // Log probabilities of the generated tokens (if requested).
		FinishReason string   `json:"finish_reason"` // Reason for finishing the completion, e.g., "stop", "length".
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`     // Number of tokens in the prompt.
		CompletionTokens int `json:"completion_tokens"` // Number of tokens in the completion.
		TotalTokens      int `json:"total_tokens"`      // Total number of tokens used.
	} `json:"usage"`
}

// FIMStreamCompletionRequest represents the request body for a streaming Fill-In-the-Middle (FIM) completion.
// It's similar to FIMCompletionRequest but includes a `Stream` field.
type FIMStreamCompletionRequest struct {
	Model            string        `json:"model"`                       // Model name to use for completion.
	Prompt           string        `json:"prompt"`                      // The prompt to start the completion from.
	Stream           bool          `json:"stream"`                      // Whether to stream the completion.  This is the key difference.
	StreamOptions    StreamOptions `json:"stream_options,omitempty"`    // Optional: Options for streaming the completion.
	Suffix           string        `json:"suffix,omitempty"`            // Optional: The suffix to complete the prompt with.
	MaxTokens        int           `json:"max_tokens,omitempty"`        // Optional: Maximum tokens to generate, > 1 and <= 4000.
	Temperature      float64       `json:"temperature,omitempty"`       // Optional: Sampling temperature, between 0 and 2.
	TopP             float64       `json:"top_p,omitempty"`             // Optional: Nucleus sampling probability threshold.
	N                int           `json:"n,omitempty"`                 // Optional: Number of completions to generate.
	Logprobs         int           `json:"logprobs,omitempty"`          // Optional: Number of log probabilities to return.
	Echo             bool          `json:"echo,omitempty"`              // Optional: Whether to echo the prompt in the completion.
	Stop             []string      `json:"stop,omitempty"`              // Optional: List of stop sequences.
	PresencePenalty  float64       `json:"presence_penalty,omitempty"`  // Optional: Penalty for new tokens based on their presence in the text so far.
	FrequencyPenalty float64       `json:"frequency_penalty,omitempty"` // Optional: Penalty for new tokens based on their frequency in the text so far.
}

type FIMStreamCompletionRequestOption func(*FIMStreamCompletionRequest)

func NewDefaultFIMStreamCompletionRequest(model string, prompt string, opts ...FIMStreamCompletionRequestOption) *FIMStreamCompletionRequest {
	locFIMCompletionReq := &FIMStreamCompletionRequest{
		Model:            model,
		Prompt:           prompt,
		Stream:           true,
		MaxTokens:        ConstFIMMaxTokensDefault,
		Temperature:      ConstFIMTemperatureDefault,
		TopP:             ConstFIMTopPDefault,
		PresencePenalty:  ConstFIMPresencePenaltyDefault,
		FrequencyPenalty: ConstFIMFrequencyPenaltyDefault,
	}
	for _, o := range opts {
		o(locFIMCompletionReq)
	}
	return locFIMCompletionReq
}

func WithStreamModel(model string, prompt string) FIMStreamCompletionRequestOption {
	return func(r *FIMStreamCompletionRequest) {
		r.Model = model
		r.Prompt = prompt
	}
}

func WithStreamMaxTokens(maxTokens int) FIMStreamCompletionRequestOption {
	return func(r *FIMStreamCompletionRequest) {
		r.MaxTokens = maxTokens
	}
}

func WithStreamTemperature(temperature float64) FIMStreamCompletionRequestOption {
	return func(r *FIMStreamCompletionRequest) {
		r.Temperature = temperature
	}
}

func WithStreamTopP(topP float64) FIMStreamCompletionRequestOption {
	return func(r *FIMStreamCompletionRequest) {
		r.TopP = topP
	}
}

func WithStreamPresencePenalty(presencePenalty float64) FIMStreamCompletionRequestOption {
	return func(r *FIMStreamCompletionRequest) {
		r.PresencePenalty = presencePenalty
	}
}

func WithStreamFrequencyPenalty(frequencyPenalty float64) FIMStreamCompletionRequestOption {
	return func(r *FIMStreamCompletionRequest) {
		r.FrequencyPenalty = frequencyPenalty
	}
}

func WithStreamEcho(echo bool) FIMStreamCompletionRequestOption {
	return func(r *FIMStreamCompletionRequest) {
		r.Echo = echo
	}
}

func WithStreamStop(stop []string) FIMStreamCompletionRequestOption {
	return func(r *FIMStreamCompletionRequest) {
		r.Stop = stop
	}
}

func WithStreamLogprobs(logprobs int) FIMStreamCompletionRequestOption {
	return func(r *FIMStreamCompletionRequest) {
		r.Logprobs = logprobs
	}
}

func WithStreamSuffix(suffix string) FIMStreamCompletionRequestOption {
	return func(r *FIMStreamCompletionRequest) {
		r.Suffix = suffix
	}
}

func WithStreamUpdate(update bool) FIMStreamCompletionRequestOption {
	return func(r *FIMStreamCompletionRequest) {
		r.Stream = update
	}
}

func CheckFIMStreamCompletionRequest(request *FIMStreamCompletionRequest) error {
	if request == nil {
		return fmt.Errorf("request cannot be nil")
	}
	if request.Model == "" {
		return fmt.Errorf("model cannot be empty")
	}
	if request.Prompt == "" {
		return fmt.Errorf("prompt cannot be empty")
	}
	if request.MaxTokens <= ConstFIMMaxTokensMin {
		return fmt.Errorf("max tokens must be > 1")
	}
	if request.MaxTokens > ConstFIMMaxTokensMax {
		return fmt.Errorf("max tokens must be <= 4000")
	}
	if request.FrequencyPenalty < ConstFIMFrequencyPenaltyMin {
		return fmt.Errorf("frequency penalty must be >= %v", ConstFIMFrequencyPenaltyMin)
	}
	if request.FrequencyPenalty > ConstFIMFrequencyPenaltyMax {
		return fmt.Errorf("frequency penalty must be < =%v", ConstFIMFrequencyPenaltyMax)
	}
	if request.Logprobs > ConstFIMLogprobsMax {
		return fmt.Errorf("logprobs must be <= %v", ConstFIMLogprobsMax)
	}
	if request.Temperature > ConstFIMTemperatureMax {
		return fmt.Errorf("temperature must be <= %v", ConstFIMTemperatureMax)
	}
	if request.TopP > ConstFIMTopPMax {
		return fmt.Errorf("top p must be <= %v", ConstFIMTopPMax)
	}

	return nil
}

// FIMStreamChoice represents a single choice within a streaming Fill-In-the-Middle (FIM) completion response.
type FIMStreamChoice struct {
	// Text generated by the model for this choice.
	Text string `json:"text"`
	// Index of this choice within the list of choices.
	Index int `json:"index"`
	// Log probabilities for the generated tokens (if available).  May be `nil`.
	Logprobs Logprobs `json:"logprobs,omitempty"`
	// Reason why the generation finished (e.g., "stop", "length"). May be `nil`.
	FinishReason interface{} `json:"finish_reason,omitempty"`
}

// FIMStreamCompletionResponse represents the full response body for a streaming Fill-In-the-Middle (FIM) completion.
// It contains metadata about the completion request and a list of choices generated by the model.
type FIMStreamCompletionResponse struct {
	// Unique identifier for the completion response.
	ID string `json:"id"`
	// List of choices generated by the model.  Each choice represents a possible completion.
	Choices []FIMStreamChoice `json:"choices"`
	// Unix timestamp (seconds since the epoch) of when the completion was created.
	Created int64 `json:"created"`
	// Name of the model used for the completion.
	Model string `json:"model"`
	// Fingerprint of the system that generated the completion.
	SystemFingerprint string `json:"system_fingerprint"`
	// Type of object returned (always "text_completion" for FIM completions).
	Object string `json:"object"`
	// Usage statistics for the completion request (if available). May be `nil`.
	Usage *StreamUsage `json:"usage,omitempty"`
}

// fimCompletionStream implements the ChatCompletionStream interface.
type fimCompletionStream struct {
	ctx    context.Context    // Context for cancellation.
	cancel context.CancelFunc // Cancel function for the context.
	resp   *http.Response     // HTTP response from the API call.
	reader *bufio.Reader      // Reader for the response body.
}

// FIMChatCompletionStream is an interface for receiving streaming chat completion responses.
type FIMChatCompletionStream interface {
	FIMRecv() (*FIMStreamCompletionResponse, error)
	FIMClose() error
}

// FIMRecv receives the next response from the stream.
func (s *fimCompletionStream) FIMRecv() (*FIMStreamCompletionResponse, error) {
	reader := s.reader
	for {
		line, err := reader.ReadString('\n') // Read until newline
		if err != nil {
			if err == io.EOF {
				return nil, io.EOF
			}
			return nil, fmt.Errorf("error reading stream: %w", err)
		}

		line = strings.TrimSpace(line)
		if line == "data: [DONE]" {
			return nil, io.EOF // End of stream
		}
		if len(line) > 6 && line[:6] == "data: " {
			trimmed := line[6:] // Trim the "data: " prefix
			var response FIMStreamCompletionResponse
			if err := json.Unmarshal([]byte(trimmed), &response); err != nil {
				return nil, fmt.Errorf("unmarshal error: %w, raw data: %s", err, trimmed)
			}
			if response.Usage == nil {
				response.Usage = &StreamUsage{}
			}
			return &response, nil
		}
	}
}

// FIMClose terminates the stream.
func (s *fimCompletionStream) FIMClose() error {
	s.cancel()
	err := s.resp.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to close response body: %w", err)
	}
	return nil
}
