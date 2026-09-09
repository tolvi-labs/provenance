package main

import (
	"fmt"
	"os"
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
    provenance check --base <ref> --head <ref> [--json-out <path>]   (not yet available)
    provenance init                                                    (not yet available)
    provenance hook install|uninstall                                  (not yet available)`)
}
