package deepseek

func resolveBaseURL(customBaseURL, fallback string) string {
	if customBaseURL == "" {
		return fallback
	}

	return customBaseURL
}

func resolveBetaBaseURL(customBaseURL string) string {
	switch customBaseURL {
	case "", "https://api.deepseek.com/", "https://api.deepseek.com", BaseURL:
		return "https://api.deepseek.com/beta/"
	default:
		return customBaseURL
	}
}
