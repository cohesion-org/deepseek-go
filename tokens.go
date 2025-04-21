package deepseek

import (
	"unicode"
)

// TokenEstimate represents an estimated token count
type TokenEstimate struct {
	EstimatedTokens int `json:"estimated_tokens"` // The total estimated prompt tokens. These are different form total tokens used.
}

// EstimateTokenCount estimates the number of tokens in a text based on character type ratios
func EstimateTokenCount(text string) *TokenEstimate {
	var total float64
	for _, r := range text {
		if unicode.Is(unicode.Han, r) {
			// Chinese character ≈ 0.6 token
			total += 0.6
		} else if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsPunct(r) || unicode.IsSymbol(r) {
			// English character/number/symbol ≈ 0.3 token
			total += 0.3
		}
		// Skip whitespace and other characters
	}

	// Round up to nearest integer
	estimatedTokens := int(total + 0.5)
	if estimatedTokens < 1 {
		estimatedTokens = 1
	}

	return &TokenEstimate{
		EstimatedTokens: estimatedTokens,
	}
}

// EstimateTokensFromMessages estimates the number of tokens in a list of chat messages
func EstimateTokensFromMessages(messages *ChatCompletionRequest) *TokenEstimate {
	var totalTokens int

	for _, msg := range messages.Messages {
		// Add tokens for role (system/user/assistant)
		totalTokens += 2 // Approximate tokens for role

		// Add tokens for content
		totalTokens += EstimateTokenCount(msg.Content).EstimatedTokens
	}

	for _, tool := range messages.Tools {
		// Add tokens for function name and description
		totalTokens += EstimateTokenCount(tool.Function.Name).EstimatedTokens
		totalTokens += EstimateTokenCount(tool.Function.Description).EstimatedTokens

		// Add tokens for function parameters if present
		if tool.Function.Parameters != nil {
			totalTokens += EstimateTokenCount(tool.Function.Parameters.Type).EstimatedTokens
			for key, value := range tool.Function.Parameters.Properties {
				totalTokens += EstimateTokenCount(key).EstimatedTokens
				totalTokens += EstimateTokenCount(value.(string)).EstimatedTokens
			}
			for _, req := range tool.Function.Parameters.Required {
				totalTokens += EstimateTokenCount(req).EstimatedTokens
			}
		}
	}

	return &TokenEstimate{
		EstimatedTokens: totalTokens,
	}
}
