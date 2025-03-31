// This is an MCP server for santactl.
package main

import (
	"fmt"
	"log"
	"os/exec"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
)

// EmptySantaCtlArgs is an empty struct used for commands that do not require
// any arguments.
type EmptySantaCtlArgs struct {
}

// runSantactlCommand runs the santactl command with the specified arguments.
func runSantactlCommand(asRoot bool, asJson bool, command string, args []string) (*mcp_golang.ToolResponse, error) {
	// Create a command with no arguments
	commandAndArgs := []string{}
	if asRoot {
		commandAndArgs = append(commandAndArgs, "sudo")
	}
	commandAndArgs = append(commandAndArgs, "/usr/local/bin/santactl")
	commandAndArgs = append(commandAndArgs, command)

	if asJson {
		commandAndArgs = append(commandAndArgs, "--json")
	}

	commandAndArgs = append(commandAndArgs, args...)

	cmd := exec.Command(commandAndArgs[0], commandAndArgs[1:]...)

	// Run the command and capture the output
	output, err := cmd.Output()
	if cmd == nil {
		msg := fmt.Sprintf("Error with %s: ", command)
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(msg + err.Error())), err
	}

	return mcp_golang.NewToolResponse(
		mcp_golang.NewTextContent(string(output))), nil
}

// santaVersion retrieves the version information for all of the santa
// components. It is a helper running for the santactl version command.
func santaVersion(args EmptySantaCtlArgs) (*mcp_golang.ToolResponse, error) {
	return runSantactlCommand(false, true, "version", nil)
}

// santaStatus retrieves the status information for the Santa daemon.
// This includes information about the rule database, watch items, and more.
func santaStatus(args EmptySantaCtlArgs) (*mcp_golang.ToolResponse, error) {
	return runSantactlCommand(false, false, "status", nil)
}

// SantaSyncArgs is the argument structure for the santaSyncRules function.
// It includes a flag for clean sync which will remove all existing rules.
type SantaSyncArgs struct {
	CleanSync bool `json:"clean_sync"`
}

// santaSyncRules forces a sync of the Santa rules database from a sync server
func santaSyncRules(args SantaSyncArgs) (*mcp_golang.ToolResponse, error) {
	if args.CleanSync {
		return runSantactlCommand(true, false, "sync", []string{"--clean"})
	}

	return runSantactlCommand(true, false, "sync", nil)
}

// santaMetrics retrieves the metrics information for the Santa daemon. It is
// a helper function for running santactl metrics.
func santaMetrics(args EmptySantaCtlArgs) (*mcp_golang.ToolResponse, error) {
	// Run the santactl command to get metrics information
	return runSantactlCommand(false, false, "metrics", nil)
}

type SantaFileinfoArgs struct {
	// Absolute path to the file on disk
	FilePath         string `json:"file_path"` // Absolute path to the file on disk
	ShowEntitlements bool   `json:"show_entitlements"`
}

// santaFileinfo retrieves information about a file using the santactl fileinfo command
func santaFileinfo(args SantaFileinfoArgs) (*mcp_golang.ToolResponse, error) {
	// Run the santactl command to get file information
	return runSantactlCommand(false, false, "fileinfo", []string{args.FilePath})
}

// Prompts below.
func santaFileInfoPrompt(args EmptySantaCtlArgs) (*mcp_golang.PromptResponse, error) {
	// Prompt for the file path
	prompt := `santactl fileinfo       

The details provided will be the same ones Santa uses to make a decision
about executables. This includes SHA-256, SHA-1, code signing information and
the type of file.
Usage: santactl fileinfo [options] [file-paths]
    --recursive (-r): Search directories recursively.
                      Incompatible with --bundleinfo.
    --json: Output in JSON format.
    --key: Search and return this one piece of information.
           You may specify multiple keys by repeating this flag.
           Valid Keys:
                       "Path"
                       "SHA-256"
                       "SHA-1"
                       "Bundle Name"
                       "Bundle Version"
                       "Bundle Version Str"
                       "Download Referrer URL"
                       "Download URL"
                       "Download Timestamp"
                       "Download Agent"
                       "Team ID"
                       "Signing ID"
                       "CDHash"
                       "Type"
                       "Page Zero"
                       "Code-signed"
                       "Rule"
                       "Entitlements"
                       "Signing Chain"
                       "Universal Signing Chain"

           Valid keys when using --cert-index:
                       "SHA-256"
                       "SHA-1"
                       "Common Name"
                       "Organization"
                       "Organizational Unit"
                       "Valid From"
                       "Valid Until"

    --cert-index: Supply an integer corresponding to a certificate of the
                  signing chain to show info only for that certificate.
                     0 up to n for the leaf certificate up to the root
                    -1 down to -n-1 for the root certificate down to the leaf
                  Incompatible with --bundleinfo.
    --filter: Use predicates of the form 'key=regex' to filter out which files
              are displayed. Valid keys are the same as for --key. Value is a
              case-insensitive regular expression which must match anywhere in
              the keyed property value for the file's info to be displayed.
              You may specify multiple filters by repeating this flag.
              If multiple filters are specified, any match will display the
              file.
    --filter-inclusive: If multiple filters are specified, they must all match
                        for the file to be displayed.
    --entitlements: If the file has entitlements, will also display them
    --bundleinfo: If the file is part of a bundle, will also display bundle
                  hash information and hashes of all bundle executables.
                  Incompatible with --recursive and --cert-index.

Examples: santactl fileinfo --cert-index 1 --key SHA-256 --json /usr/bin/yes
          santactl fileinfo --key SHA-256 --json /usr/bin/yes
          santactl fileinfo /usr/bin/yes /bin/*
          santactl fileinfo /usr/bin -r --key Path --key SHA-256 --key Rule
          santactl fileinfo /usr/bin/* --filter Type=Script --filter Path=zip`

	return mcp_golang.NewPromptResponse("santactl_fileinfo_prompt",
		mcp_golang.NewPromptMessage(mcp_golang.NewTextContent(prompt),
			mcp_golang.RoleUser)), nil
}

func santaSubCommandPrompt(args EmptySantaCtlArgs) (*mcp_golang.PromptResponse, error) {
	// Prompt for the santactl subcommand
	santactlPrompt := `santactl has the following sub commands

		Usage: santactl:
	  fileinfo - Prints information about a file.
	   metrics - Show Santa metric information.
	  printlog - Prints the contents of Santa protobuf log files as JSON.
		  rule - Manually add/remove/check rules.
		status - Show Santa status information.
		  sync - Synchronizes Santa with a configured server.
	   version - Show Santa component versions.
	   `

	return mcp_golang.NewPromptResponse("santactl subcommands",
		mcp_golang.NewPromptMessage(mcp_golang.NewTextContent(santactlPrompt),
			mcp_golang.RoleAssistant)), nil
}

func santaStatusCommandPrompt(args EmptySantaCtlArgs) (*mcp_golang.PromptResponse, error) {
	// Prompt for the status command
	statusPrompt := `santactl status produces a set of key value pairs
	describing how santa operates. The database section is for execution rules
	and the watch items related to file access authorization rules.

	Do not conflate the watch items and the database items keep all statistics
	separate.

	Additionally if santa is not configured to use a sync service then local
	admin users may add rules using the santactl rule command. This is also true
	if santa does not have any static rules.
	`

	return mcp_golang.NewPromptResponse("santactl_status_prompt",
		mcp_golang.NewPromptMessage(mcp_golang.NewTextContent(statusPrompt),
			mcp_golang.RoleAssistant)), nil
}

func main() {
	done := make(chan struct{})

	server := mcp_golang.NewServer(stdio.NewStdioServerTransport())

	err := server.RegisterTool("santactl_version",
		"List santa daemon versions", santaVersion)
	if err != nil {
		log.Fatalf("Failed to register santactl_version tool: %v", err)
	}

	err = server.RegisterTool("santactl_status",
		"Get the status of the Santa daemon", santaStatus)
	if err != nil {
		log.Fatalf("Failed to register santactl_status tool: %v", err)
	}

	err = server.RegisterTool("santactl_metrics",
		"Get performance metrics from the Santa daemon", santaMetrics)

	if err != nil {
		log.Fatalf("Failed to register santactl_metrics tool: %v", err)
	}

	err = server.RegisterTool("santactl_sync",
		"Have Santa perform a sync of the Santa rules database from the sync server",
		santaSyncRules)

	if err != nil {
		log.Fatalf("Failed to register santactl_sync tool: %v", err)
	}

	err = server.RegisterTool("santactl_fileinfo",
		"Get file information for a given absolute file path", santaFileinfo)
	if err != nil {
		log.Fatalf("Failed to register santactl_fileinfo tool: %v", err)
	}

	err = server.RegisterPrompt("santactl_status_prompt",
		"explain the data in a santactl_status output", santaStatusCommandPrompt)

	if err != nil {
		log.Fatalf("Failed to register prompt: %v", err)
	}

	err = server.RegisterPrompt("santactl_subcommand_prompt",
		"list of valid santactl subcommands", santaSubCommandPrompt)

	if err != nil {
		log.Fatalf("Failed to register prompt: %v", err)
	}

	err = server.RegisterPrompt("santactl_fileinfo_prompt",
		"list of valid santactl subcommands", santaFileInfoPrompt)

	if err != nil {
		log.Fatalf("Failed to register prompt: %v", err)
	}

	err = server.Serve()
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
	<-done
}
