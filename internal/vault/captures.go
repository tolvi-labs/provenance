package vault

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

var validAttributions = map[string]bool{
	"human":         true,
	"ai-assisted":   true,
	"ai-and-edited": true,
}

// Capture is one vault/captures/*.md file's frontmatter, plus its path.
type Capture struct {
	Path        string
	Tags        []string `yaml:"tags"`
	Date        string   `yaml:"date"`
	Status      string   `yaml:"status"`
	Repo        string   `yaml:"repo"`
	Decisions   []string `yaml:"decisions"`
	Attribution string   `yaml:"attribution"`
}

// LoadCaptures reads every *.md file under <vaultDir>/captures.
func LoadCaptures(vaultDir string) ([]Capture, error) {
	dir := filepath.Join(vaultDir, "captures")
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading %s: %w", dir, err)
	}

	var captures []Capture
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", path, err)
		}
		fm, err := extractFrontmatter(string(content))
		if err != nil {
			return nil, fmt.Errorf("parsing frontmatter in %s: %w", path, err)
		}
		var c Capture
		if err := yaml.Unmarshal([]byte(fm), &c); err != nil {
			return nil, fmt.Errorf("unmarshaling frontmatter in %s: %w", path, err)
		}
		c.Path = filepath.Join("vault", "captures", entry.Name())
		captures = append(captures, c)
	}
	return captures, nil
}

// ValidateCapture checks a capture fragment is well-formed: required fields
// present, a valid attribution value, and every referenced decision slug
// actually exists among knownDecisions. Returns a list of problems, empty
// if the capture is well-formed.
func ValidateCapture(c Capture, knownDecisions []Decision) []string {
	var problems []string
	if c.Status == "" {
		problems = append(problems, "missing status")
	}
	if len(c.Decisions) == 0 {
		problems = append(problems, "missing decisions (must reference at least one)")
	}
	if !validAttributions[c.Attribution] {
		problems = append(problems, fmt.Sprintf("invalid attribution %q (must be human, ai-assisted, or ai-and-edited)", c.Attribution))
	}

	known := make(map[string]bool, len(knownDecisions))
	for _, d := range knownDecisions {
		known[d.Slug] = true
	}
	for _, slug := range c.Decisions {
		if !known[slug] {
			problems = append(problems, fmt.Sprintf("references unknown decision %q", slug))
		}
	}
	return problems
}
