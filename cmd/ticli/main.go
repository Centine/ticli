package main

import (
	"github.com/centine/ticli/cmd"
	"github.com/pterm/pterm"
)

var Version = "dev"
var Build = "dev"

func main() {

	// Enable debug messages.
	pterm.EnableDebugMessages()

	pterm.Info.Printfln("ticli v%s, build %s.", Version, Build) // Print Info.
	checkForUpdates()

	cmd.Execute()

}
