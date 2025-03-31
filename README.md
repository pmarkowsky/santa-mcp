# santa-mcp

A PoC MCP Server for [Santa](https://github.com/northpolesec/santa). 

##  What is this?

> [!CAUTION]
> This is intended solely as a demonstration and is not production-ready. It is not an officially supported product of North Pole Security.

This is a Proof of Concept [Model Context Procotol server](https://modelcontextprotocol.io/introduction) for [santactl](https://northpole.dev/binaries/santactl.html). 

It allows you to connect and drive Santa via an LLM that implements an MCP
client.

[![Example session With Claude Desktop](https://img.youtube.com/vi/Q_bHdz3wFzQ/0.jpg)](https://youtu.be/Q_bHdz3wFzQ)

* [Session with Claude](https://claude.ai/share/9425ecd9-6dfb-40c7-adbd-6478ec857d4a)

## Quickstart

* Install [Claude Desktop](https://claude.ai/download) if you don't already have it.
* Run `make`
* Copy the `santa-mcp` binary somewhere on your system
* Edit the `claude_desktop_config.json` to point to the path from the previous
 step
* Copy the `claude_desktop_config.json` file to your `~/Library
  * `cp claude_desktop_config.json ~/Library/Application\ Support/Claude/`
* Open Claude desktop you should see 4 tools
 * santactl_fileinfo
 * santactl_status
 * santactl_sync
 * santactl_version

* Ask it questions about Santa e.g. `Why is Santa blocking osascript?`

## Dependencies

This depends on [MCP Golang](https://mcpgolang.com/).
