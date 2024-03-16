package dribbble

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

// ------------------------------------------------------------------------
// outputs

func ToToml(v any) ([]byte, error) {
	return toml.Marshal(v)
}

func ToYaml(v any) ([]byte, error) {
	return yaml.Marshal(v)
}

func ToFile(filename string, content []byte) error {
	return os.WriteFile(filename, content, 0644)
}

// ------------------------------------------------------------------------
// string formatting

// paginationQueryString returns a query string for setting the page and per_page parameters.
func paginationQueryString(page, perPage int) string {
	if page <= 0 {
		page = 1
	}

	if perPage <= 0 {
		perPage = 100
	}

	return fmt.Sprintf("?page=%d&per_page=%d", page, perPage)
}

// writeIfNotEmpty writes formatted key-value pairs to a StringBuilder if the value is not empty.
func writeIfNotEmpty(sb *strings.Builder, key, value string) {
	if value != "" {
		sb.WriteString(fmt.Sprintf("%s: %s\n", color.HiBlackString(key), value))
	}
}

func writeTitleString(sb *strings.Builder, key string) {
	sb.WriteString(fmt.Sprintf("--%s--\n", color.HiBlackString(key)))
}

// formatTags formats a slice of tags into a comma-separated string.
func formatTags(tags []string) string {
	if len(tags) > 0 {
		return strings.Join(tags, ", ")
	}
	return ""
}
