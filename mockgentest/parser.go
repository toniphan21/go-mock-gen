package mockgentest

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexflint/go-arg"
	mockgen "nhatp.com/go/mock-gen"
	"nhatp.com/go/mock-gen/internal/cmd"
)

func hasFile(dir, filename string) bool {
	targetPath := filepath.Join(dir, filename)
	info, err := os.Stat(targetPath)
	if err != nil {
		return false
	}

	if info.IsDir() {
		return false
	}
	return true
}

func parseGenerateSH(dir, filename string) ([]mockgen.Config, error) {
	filePath := filepath.Join(dir, filename)
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(fileContent), "\n")
	idx := -1
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "go run nhatp.com/go/mock-gen/cmd/go-mock-gen") {
			idx = i
			break
		}
	}

	if idx == -1 || idx+1 == len(lines) {
		return nil, errors.New("invalid generate.sh file content")
	}

	lines = lines[idx+1:]
	var options []string
	for _, v := range lines {
		options = append(options, strings.ReplaceAll(strings.TrimSpace(v), "\\", ""))
	}
	rawArgs := strings.Fields(strings.Join(options, " "))

	var args cmd.Arguments
	p, err := arg.NewParser(arg.Config{}, &args)
	if err != nil {
		return nil, err
	}

	err = p.Parse(rawArgs)
	if err != nil {
		fmt.Printf("Error parsing: %v\n", err)
	}
	return cmd.ToConfigs(dir, args)
}
