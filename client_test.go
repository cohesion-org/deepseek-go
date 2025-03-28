package deepseek_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cohesion-org/deepseek-go"
)

func TestCreateChatCompletion_NilRequest(t *testing.T) {
	client, err := deepseek.NewClientWithOptions("token")
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = client.CreateChatCompletion(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil request")
	}
	if err.Error() != "request cannot be nil" {
		t.Errorf("expected error 'request cannot be nil', got %q", err.Error())
	}
}

func TestCreateChatCompletion_ErrorHandling(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": {"message": "invalid request"}}`))
	}))
	defer testServer.Close()

	client, err := deepseek.NewClientWithOptions("token", deepseek.WithBaseURL(testServer.URL+"/"))
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = client.CreateChatCompletion(context.Background(), &deepseek.ChatCompletionRequest{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*deepseek.APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.Message != "Bad request" {
		t.Errorf("expected error message 'Bad request', got %s", apiErr.Message)
	}
}

func TestCreateChatCompletionStream_ErrorHandling(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": {"message": "stream error"}}`))
	}))
	defer testServer.Close()

	client, err := deepseek.NewClientWithOptions("token", deepseek.WithBaseURL(testServer.URL+"/"))
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = client.CreateChatCompletionStream(context.Background(), &deepseek.StreamChatCompletionRequest{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*deepseek.APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.Message != "Bad request" {
		t.Errorf("expected error message 'Bad request', got %s", apiErr.Message)
	}
}
