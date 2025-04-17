module github.com/cohesion-org/deepseek-go

go 1.24.0

toolchain go1.24.2

require (
	github.com/joho/godotenv v1.5.1
	github.com/ollama/ollama v0.6.5
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Retracting v1.1.0 and v1.0.1 because it was
// a premature release. Please get the supported version
// from the releases page.
retract (
	v1.1.0
	v1.0.1
)
