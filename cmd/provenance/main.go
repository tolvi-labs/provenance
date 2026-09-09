package main

import (
	"fmt"
	"os"

	"github.com/tolvi-labs/provenance/internal/cmdcheck"
	"github.com/tolvi-labs/provenance/internal/cmdhook"
	"github.com/tolvi-labs/provenance/internal/cmdinit"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "version":
		fmt.Println("provenance " + version)
	case "check":
		os.Exit(cmdcheck.Run(os.Args[2:]))
	case "init":
		os.Exit(cmdinit.Run(os.Args[2:]))
	case "hook":
		os.Exit(cmdhook.Run(os.Args[2:]))
	case "-h", "--help", "help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "provenance: unknown command %q\n", os.Args[1])
		printUsage()
		os.Exit(2)
	}
}

func printUsage() {
	fmt.Println(`provenance — the capture-enforcement gate at the code→push boundary

Usage:
  provenance version
  provenance check --base <ref> --head <ref> [--json-out <path>]
  provenance init
  provenance hook install|uninstall`)
}
