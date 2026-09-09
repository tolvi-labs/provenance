package hook

import (
	"fmt"
	"os"
	"path/filepath"
)

const shim = "#!/usr/bin/env sh\n" +
	"# provenance pre-push hook — installed by `provenance hook install`.\n" +
	"# Reads git's pre-push stdin and runs provenance check per ref pushed.\n" +
	"# Blocks the push (non-zero exit) when a governed path is unaddressed.\n" +
	"command -v provenance >/dev/null 2>&1 || exit 0\n" +
	"provenance hook run\n"

// hookPath returns the path to .git/hooks/pre-push for the repo at repoDir.
func hookPath(repoDir string) string {
	return filepath.Join(repoDir, ".git", "hooks", "pre-push")
}

// Install writes the pre-push shim to repoDir's .git/hooks/pre-push.
// It refuses to overwrite an existing hook unless force is true.
func Install(repoDir string, force bool) error {
	path := hookPath(repoDir)
	if _, err := os.Stat(path); err == nil && !force {
		return fmt.Errorf("%s already exists (use --force to overwrite)", path)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("creating hooks dir: %w", err)
	}
	if err := os.WriteFile(path, []byte(shim), 0755); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}

// Uninstall removes the pre-push shim if it was installed by this tool.
func Uninstall(repoDir string) error {
	path := hookPath(repoDir)
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("reading %s: %w", path, err)
	}
	if string(content) != shim {
		return fmt.Errorf("%s was not installed by provenance (leaving it in place)", path)
	}
	return os.Remove(path)
}
