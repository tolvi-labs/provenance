package cmdhook

import (
	"fmt"
	"io"
	"os"

	"github.com/tolvi-labs/provenance/internal/cmdcheck"
	"github.com/tolvi-labs/provenance/internal/gitutil"
	"github.com/tolvi-labs/provenance/internal/hook"
)

// Run implements `provenance hook install|uninstall|run`.
func Run(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "provenance hook: expected install, uninstall, or run")
		return 2
	}
	switch args[0] {
	case "install":
		force := len(args) > 1 && args[1] == "--force"
		if err := hook.Install(".", force); err != nil {
			fmt.Fprintf(os.Stderr, "provenance hook install: %v\n", err)
			return 1
		}
		fmt.Println("✓ Installed pre-push hook")
		return 0
	case "uninstall":
		if err := hook.Uninstall("."); err != nil {
			fmt.Fprintf(os.Stderr, "provenance hook uninstall: %v\n", err)
			return 1
		}
		fmt.Println("✓ Removed pre-push hook")
		return 0
	case "run":
		return runHook(os.Stdin)
	default:
		fmt.Fprintf(os.Stderr, "provenance hook: unknown subcommand %q\n", args[0])
		return 2
	}
}

const zeroSHA = "0000000000000000000000000000000000000000"

// runHook is invoked by the installed shim with git's pre-push stdin. It
// runs a check per ref being pushed and blocks (non-zero) if any is
// blocked. A brand-new branch (remote SHA all zeros) has nothing to diff
// against yet and is skipped.
func runHook(stdin io.Reader) int {
	raw, err := io.ReadAll(stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "provenance hook run: reading stdin: %v\n", err)
		return 2
	}
	refs, err := gitutil.ParsePrePushStdin(string(raw))
	if err != nil {
		fmt.Fprintf(os.Stderr, "provenance hook run: %v\n", err)
		return 2
	}

	blocked := false
	for _, ref := range refs {
		if ref.RemoteSHA == zeroSHA {
			continue
		}
		code := cmdcheck.Run([]string{"--base", ref.RemoteSHA, "--head", ref.LocalSHA})
		switch code {
		case 0:
			// pass, keep checking remaining refs
		case 1:
			blocked = true
		default:
			return code // a real error, not a gate block — surface it as-is
		}
	}
	if blocked {
		return 1
	}
	return 0
}
