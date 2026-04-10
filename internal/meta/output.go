package meta

import (
	"regexp"
	"strings"
)

func ExtractOutput(in []byte) string {
	output := string(in)
	lines := strings.Split(output, "\n")
	trimmed := strings.Builder{}
	for _, line := range lines {
		if strings.HasPrefix(line, "=== RUN") {
			continue
		}
		if strings.HasPrefix(line, "--- FAIL") {
			continue
		}
		if strings.HasPrefix(line, "FAIL") {
			continue
		}
		if line == "" {
			continue
		}
		trimmed.WriteString(line)
		trimmed.WriteString("\n")
	}
	return extractPanicMessage(trimmed.String())
}

func extractPanicMessage(input string) string {
	re := regexp.MustCompile(`(?s)panic:\s+(.*?)\s+\[recovered`)

	match := re.FindStringSubmatch(input)

	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return input
}
