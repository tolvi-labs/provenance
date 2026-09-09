package cmdhook

import (
	"strings"
	"testing"
)

func TestRunHook_SkipsNewBranchPush(t *testing.T) {
	// A brand-new branch push has an all-zero remote SHA; nothing to diff
	// against, so it must not error or attempt a check.
	input := "refs/heads/new-branch abc123 refs/heads/new-branch 0000000000000000000000000000000000000000\n"
	code := runHook(strings.NewReader(input))
	if code != 0 {
		t.Fatalf("expected exit 0 for a new-branch push, got %d", code)
	}
}

func TestRunHook_RejectsMalformedInput(t *testing.T) {
	code := runHook(strings.NewReader("not-well-formed\n"))
	if code != 2 {
		t.Fatalf("expected exit 2 for malformed input, got %d", code)
	}
}
