# go-a2a
Go client/server implementation for Google A2A(**A**gent **T**o **A**gent) Protocol.

# Stage
This project is in the early stage of development. The API may change in the future.
But you can use it for development now. Use in production is not recommended yet.

## Install

```go
go get github.com/zhengrenjie/go-a2a
```

# Roadmap

## protocol definition
- [x] A2A Protocol definition

## server implementation
- [x] A2A Server implementation with standard http server
- [ ] Support SSE for streaming responses
- [ ] More useful options for server configuration
- [ ] Server side logging

## client implementation
WIP

# Usage

## Server

```go
func main() {
    // create a new server with default options
    srv := server.NewServer(yourHandler)

    err := server.NewStandardA2AServerHost(":6789").Host(srv)
    if err != nil {
        log.Fatal(err)
    }
}
```

# Reference

See: https://developers.googleblog.com/en/a2a-a-new-era-of-agent-interoperability/

Protocol Specification: https://google.github.io/A2A/#/
