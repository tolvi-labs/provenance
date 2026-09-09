package cmdinit

import (
	"flag"
	"fmt"
	"os"

	"github.com/tolvi-labs/provenance/internal/config"
)

// Run implements `provenance init`.
func Run(args []string) int {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	repoDir := fs.String("repo", ".", "path to the repository")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if err := config.Scaffold(*repoDir); err != nil {
		fmt.Fprintf(os.Stderr, "provenance init: %v\n", err)
		return 1
	}
	fmt.Println("✓ Wrote provenance.yml")
	return 0
}
