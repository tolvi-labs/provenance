package gitutil

import (
	"fmt"
	"os/exec"
	"strings"
)

// ChangedFiles returns the list of file paths that differ between base and
// head in the git repository at repoDir, using `git diff --name-only`.
func ChangedFiles(repoDir, base, head string) ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", base+".."+head)
	cmd.Dir = repoDir
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git diff --name-only %s..%s: %w", base, head, err)
	}
	trimmed := strings.TrimSpace(string(out))
	if trimmed == "" {
		return nil, nil
	}
	return strings.Split(trimmed, "\n"), nil
}

// PrePushRef represents one line of pre-push hook stdin:
// <local-ref> <local-sha> <remote-ref> <remote-sha>
type PrePushRef struct {
	LocalRef  string
	LocalSHA  string
	RemoteRef string
	RemoteSHA string
}

// ParsePrePushStdin parses the lines git passes to a pre-push hook on stdin.
func ParsePrePushStdin(input string) ([]PrePushRef, error) {
	var refs []PrePushRef
	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 4 {
			return nil, fmt.Errorf("malformed pre-push line: %q", line)
		}
		refs = append(refs, PrePushRef{
			LocalRef:  fields[0],
			LocalSHA:  fields[1],
			RemoteRef: fields[2],
			RemoteSHA: fields[3],
		})
	}
	return refs, nil
}
