package main

import (
	"github.com/centine/ticli/cmd"
	"github.com/centine/ticli/internal/webhosttmp"
	"github.com/pterm/pterm"
)

var Version = "dev"
var Build = "dev"

func main() {
	// Temporary hack to host the script bundles
	go webhosttmp.StartServer() // Start the server and its handling of exit signals in the background

	// Enable debug messages.
	pterm.EnableDebugMessages()

	pterm.Info.Printfln("ticli v%s, build %s.", Version, Build) // Print Info.
	checkForUpdates()

	cmd.Execute()

}
