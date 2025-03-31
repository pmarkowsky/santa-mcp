GO_SRCS = cmd/main.go
GO_BINS = santa-mcp

all: $(GO_BINS)

santa-mcp: cmd/main.go
	go build -o santa-mcp ./cmd/main.go
