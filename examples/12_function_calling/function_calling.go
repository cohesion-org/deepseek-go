package main

import (
        "context"
        "fmt"
        "log"
        "os"
        "time"

        "github.com/cohesion-org/deepseek-go"
)

func main() {
        if len(os.Args) < 2 {
                fmt.Println("Usage: function-call <message>")
                return
        }
        userMessage := os.Args[1]
        FunctionCalling(userMessage)
}

// When you provide multiple utility functions, the model may tell you to call multiple utility functions
// in one request and multiple times in one conversation. Specifically, the model may ak you to pass the
// result of the previously called utility function as input to the second called utility function.
//
// Example:
//   User:      What the weather and relative humidity are at the current location?
//   Assistant: call GetLocation()
//   Tool:      "Beijing"
//   Assistant: call GetTemperature("Beijing") and call GetRelativeHumidity("Beijing")
//   Tool:      "11℃", "35%"
//   Assistant: The current temperature and humidity in Beijing are 11℃ and 35%.
//
// Yes, this happens in one request, so you can use it to solve complex problems.

var toolGetTime = deepseek.Tool{
        Type: "function",
        Function: deepseek.Function{
                Name: "GetTime",
                Description: "" +
                        "Get the current date and time. The returned time string format is RFC3339. " +
                        "Be careful not to abuse this function unless you really need to get the real world time.",
        },
}

func onGetTime() string {
        s := time.Now().Format(time.RFC3339)
        return "current time: " + s
}

func FunctionCalling(userMessage string) {
        client := deepseek.NewClient(os.Getenv("DEEPSEEK_API_KEY"))

        request := &deepseek.ChatCompletionRequest{
                Model: deepseek.DeepSeekChat,
                Messages: []deepseek.ChatCompletionMessage{
                        {Role: deepseek.ChatMessageRoleUser, Content: userMessage},
                },
                Tools: []deepseek.Tool{toolGetTime},
        }
        ctx := context.Background()
        response, err := client.CreateChatCompletion(ctx, request)
        if err != nil {
                log.Fatalf("error: %v", err)
        }

        msg := response.Choices[0].Message
        fmt.Println("response:", msg.Content)
        fmt.Println("tool calls:", msg.ToolCalls)

        toolCalls := msg.ToolCalls
        if len(toolCalls) == 0 {
                // No tool calls, just return the response
                return
        }

        question := deepseek.ChatCompletionMessage{
                Role:      deepseek.ChatMessageRoleAssistant,
                Content:   msg.Content,
                ToolCalls: toolCalls,
        }

        // Tool call dispatcher
        var answer deepseek.ChatCompletionMessage
        switch toolCalls[0].Function.Name {
        case "GetTime":
                answer = deepseek.ChatCompletionMessage{
                        Role:       deepseek.ChatMessageRoleTool,
                        Content:    onGetTime(),
                        ToolCallID: toolCalls[0].ID,
                }
        // Add more cases here for additional tools
        default:
                answer = deepseek.ChatCompletionMessage{
                        Role:       deepseek.ChatMessageRoleTool,
                        Content:    "Unknown tool call",
                        ToolCallID: toolCalls[0].ID,
                }
        }

        messages := request.Messages
        messages = append(messages, question, answer)
        toolReq := &deepseek.ChatCompletionRequest{
                Model:    request.Model,
                Messages: messages,

                // It is not recommended to use it unless it is a special case.
                // The official said that they are actively fixing the problem
                // of infinite loop calls.

                // Not using this field will force the model to call the
                // utility function only once per conversation.

                // Don't try to delete only the utility functions which may cause
                // infinite loop calls. I have tried this and the model will still
                // call the deleted utility functions unless you delete all the
                // utility functions, like the commented out code.

                // Tools: request.Tools, // This is the key to implement chain calls.
        }

        response, err = client.CreateChatCompletion(ctx, toolReq)
        if err != nil {
                log.Fatalf("error: %v", err)
        }

        // will return the current time
        fmt.Println("response:", response.Choices[0].Message.Content)
}
