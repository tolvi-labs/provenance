package vault

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Decision is one vault/decisions/*.md file's frontmatter, plus the slug
// derived from its filename.
type Decision struct {
	Slug     string   // filename without the .md extension
	Tags     []string `yaml:"tags"`
	Date     string   `yaml:"date"`
	Status   string   `yaml:"status"`
	Repo     string   `yaml:"repo"`
	XGoverns []string `yaml:"x-governs"`
}

// LoadDecisions reads every *.md file under <vaultDir>/decisions and parses
// its frontmatter.
func LoadDecisions(vaultDir string) ([]Decision, error) {
	dir := filepath.Join(vaultDir, "decisions")
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading %s: %w", dir, err)
	}

	var decisions []Decision
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
		var d Decision
		if err := yaml.Unmarshal([]byte(fm), &d); err != nil {
			return nil, fmt.Errorf("unmarshaling frontmatter in %s: %w", path, err)
		}
		d.Slug = strings.TrimSuffix(entry.Name(), ".md")
		decisions = append(decisions, d)
	}
	return decisions, nil
}

// extractFrontmatter returns the YAML block between the first two `---`
// delimiter lines. Shared by decisions.go and captures.go.
func extractFrontmatter(content string) (string, error) {
	const delim = "---"
	lines := strings.Split(content, "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != delim {
		return "", fmt.Errorf("no opening frontmatter delimiter")
	}
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == delim {
			return strings.Join(lines[1:i], "\n"), nil
		}
	}
	return "", fmt.Errorf("no closing frontmatter delimiter")
}

// FilterByStatus returns only decisions whose Status is in allowed.
func FilterByStatus(decisions []Decision, allowed []string) []Decision {
	allowedSet := make(map[string]bool, len(allowed))
	for _, s := range allowed {
		allowedSet[s] = true
	}
	var out []Decision
	for _, d := range decisions {
		if allowedSet[d.Status] {
			out = append(out, d)
		}
	}
	return out
}
