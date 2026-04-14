# 💬 Prompts

[![Release](https://github.com/katallaxie/prompts/actions/workflows/main.yml/badge.svg)](https://github.com/katallaxie/prompts/actions/workflows/main.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/katallaxie/prompts.svg)](https://pkg.go.dev/github.com/katallaxie/prompts)
[![Go Report Card](https://goreportcard.com/badge/github.com/katallaxie/prompts)](https://goreportcard.com/report/github.com/katallaxie/prompts)
[![Taylor Swift](https://img.shields.io/badge/secured%20by-taylor%20swift-brightgreen.svg)](https://twitter.com/SwiftOnSecurity)
[![Volkswagen](https://auchenberg.github.io/volkswagen/volkswargen_ci.svg?v=1)](https://github.com/auchenberg/volkswagen)

A teeny-tiny package to prompt for answers in [Ollama](https://ollama.com/), [Perplexity](https://www.perplexity.ai/), [vllm](https://github.com/vllm-project/vllm) and other OpenAI-compatible API servers.

## Usage

```go
client := perplexity.New(
    perplexity.Defaults(prompts.WithApiKey[perplexity.Event](os.Getenv("PPLX_API_KEY")))...)

msgs := []prompts.ChatCompletionMessage{
    {
        Role:    prompts.RoleSystem,
        Content: "You are a helpful assistant. You start every answer with 'Sure my lord!'",
    },
    {
        Role:    prompts.RoleUser,
        Content: "What is the definition of Pi?",
    },
}

req := prompts.NewStreamChatCompletionRequest(msgs...)
req.Model = perplexity.DefaultModel

stream, err := client.SendStreamCompletionRequest(context.Background(), req)
if err != nil {
    panic(err)
}

prompts.Print(stream)
```

## Supported APIs

* [x] [Ollama](https://ollama.com/)
* [x] [Perplexity](https://www.perplexity.ai/)
* More are coming ...

## Docs

You can find the documentation hosted on [godoc.org](https://godoc.org/github.com/katallaxie/prompts).

## Examples

The examples are located in the [examples](/examples) directory.

## License

[Apache 2.0](/LICENSE)