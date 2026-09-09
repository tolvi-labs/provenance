package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
	"gopkg.in/yaml.v3"
)

// Config is the per-repo capture contract, provenance.yml at the repo root.
type Config struct {
	GovernedStatuses []string `yaml:"governed_statuses"`
	ExemptPaths      []string `yaml:"exempt_paths"`
}

// DefaultConfig is used when no provenance.yml exists.
func DefaultConfig() Config {
	return Config{
		GovernedStatuses: []string{"active"},
	}
}

// Load reads provenance.yml from repoDir, or returns DefaultConfig if it
// does not exist.
func Load(repoDir string) (Config, error) {
	path := filepath.Join(repoDir, "provenance.yml")
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return Config{}, fmt.Errorf("reading %s: %w", path, err)
	}
	var c Config
	if err := yaml.Unmarshal(content, &c); err != nil {
		return Config{}, fmt.Errorf("parsing %s: %w", path, err)
	}
	if len(c.GovernedStatuses) == 0 {
		c.GovernedStatuses = []string{"active"}
	}
	return c, nil
}

// Scaffold writes a default provenance.yml to repoDir, refusing to
// overwrite an existing one.
func Scaffold(repoDir string) error {
	path := filepath.Join(repoDir, "provenance.yml")
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("%s already exists", path)
	}
	const template = "governed_statuses: [active]\nexempt_paths: []\n"
	return os.WriteFile(path, []byte(template), 0644)
}

// ValidatePatterns reports the first syntactically invalid glob pattern in
// patterns, if any. MatchesAny treats an unparseable pattern as a non-match,
// so an unvalidated typo would silently leave a path ungoverned; callers
// validate upfront and fail closed instead.
func ValidatePatterns(patterns []string) error {
	for _, pattern := range patterns {
		if !doublestar.ValidatePattern(pattern) {
			return fmt.Errorf("invalid glob pattern %q", pattern)
		}
	}
	return nil
}

// MatchesAny reports whether path matches any of the given glob patterns.
// Patterns support "**" for recursive matching (via doublestar).
func MatchesAny(path string, patterns []string) bool {
	for _, pattern := range patterns {
		if ok, _ := doublestar.Match(pattern, path); ok {
			return true
		}
	}
	return false
}
