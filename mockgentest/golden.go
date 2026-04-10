package mockgentest

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	genlib "nhatp.com/go/gen-lib"
	"nhatp.com/go/gen-lib/file"
	mockgen "nhatp.com/go/mock-gen"
	"nhatp.com/go/mock-gen/internal/meta"
)

type GoldenTestCase struct {
	Source        []file.File
	GoldenFiles   []file.File
	ExecutedTests [][]string
}

type integrationTestCase struct {
	TestName           string
	ExpectedPass       bool
	ExpectedFailOutput []string
}

func (tc *GoldenTestCase) Setup(dir string) error {
	var files []file.File

	if tc.Source != nil {
		for _, f := range tc.Source {
			files = append(files, f)
		}
	}

	err := genlib.SetupSourceCode(dir, files)
	if err != nil {
		return err
	}
	return nil
}

func (tc *GoldenTestCase) ParseConfigs(dir string) ([]mockgen.Config, error) {
	if hasFile(dir, "generate.sh") {
		return parseGenerateSH(dir, "generate.sh")
	}
	return nil, errors.New("there is no configure files")
}

func (tc *GoldenTestCase) ExecuteMockGen(dir string, configs []mockgen.Config, options ...mockgen.Option) (genlib.FileManager, error) {
	fm := mockgen.NewFileManager(dir, genlib.WithBinaryName(mockgen.BinaryName))
	ignoredFM := mockgen.NewFileManager(dir, genlib.WithHeaderCommentTemplate("// ignored-content"))

	generator := mockgen.New(fm, options...)

	pkgs, err := genlib.LoadPackages(dir)
	if err != nil {
		return nil, err
	}

	ignoredGoldenFiles := false
	for _, pkg := range pkgs {
		if err := generator.Generate(pkg, configs); err != nil {
			return nil, err
		}

		// run integration test if the golden-files are ignored
		for _, v := range tc.GoldenFiles {
			hasGenerated := false
			for _, gf := range fm.Files() {
				if gf.RelPath == v.FilePath() {
					hasGenerated = true
				}
			}

			if !hasGenerated {
				continue
			}

			content := string(v.FileContent())
			firstIdx := strings.Index(content, "\n")
			if firstIdx == -1 {
				continue
			}

			firstLine := content[:firstIdx]
			if strings.TrimSpace(firstLine) != "// ignored-content" {
				continue
			}

			if _, err = ignoredFM.Make(pkg, "", v.FilePath()); err != nil {
				return nil, err
			}
			ignoredGoldenFiles = true
		}

		if ignoredGoldenFiles {
			relPath, err := filepath.Rel(fm.RootDir(), pkg.Dir)
			if err != nil {
				return nil, err
			}

			if relPath == "." {
				relPath = ""
			}

			if err = tc.saveGeneratedFiles(fm); err != nil {
				return nil, err
			}

			if err = tc.runIntegrationTests(pkg.Dir, "expected_test_result.txt"); err != nil {
				return nil, err
			}
		}
	}

	if ignoredGoldenFiles {
		return ignoredFM, nil
	}
	return fm, nil
}

func (tc *GoldenTestCase) Run(t *testing.T, dir string, configs []mockgen.Config, options ...mockgen.Option) {
	fm, err := tc.ExecuteMockGen(dir, configs, options...)
	require.NoError(t, err)

	var out = make(map[string]string) // path -> content
	for _, f := range fm.Files() {
		out[f.RelPath] = f.Content()
	}

	for _, f := range tc.GoldenFiles {
		result, ok := out[f.FilePath()]
		if !ok {
			t.Errorf(`expected file "%s" but file is not generated`, f.FilePath())
		}

		expected := string(f.FileContent())
		assert.Equal(t, expected, result)
	}
}

func (tc *GoldenTestCase) saveGeneratedFiles(fm mockgen.FileManager) error {
	for _, out := range fm.Files() {
		content := out.Content()
		outer := filepath.Dir(out.RelPath)
		if err := os.MkdirAll(outer, 0755); err != nil {
			return err
		}

		if err := os.WriteFile(out.FullPath, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}

func (tc *GoldenTestCase) runIntegrationTests(dir, filename string) error {
	if !hasFile(dir, filename) {
		return nil
	}

	filePath := filepath.Join(dir, filename)
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(fileContent), "\n")
	var iTestCases []integrationTestCase
	var current *integrationTestCase
	for _, v := range lines {
		line := strings.TrimSpace(v)
		if strings.HasPrefix(line, "=== RUN") {
			current = &integrationTestCase{}
			current.TestName = strings.TrimSpace(strings.TrimPrefix(line, "=== RUN"))
			continue
		}

		if strings.HasPrefix(line, "--- PASS") && current != nil {
			current.ExpectedPass = true
			iTestCases = append(iTestCases, *current)
			current = nil
			continue
		}

		if strings.HasPrefix(line, "--- FAIL") && current != nil {
			current.ExpectedPass = false
			iTestCases = append(iTestCases, *current)
			current = nil
			continue
		}

		if current != nil {
			current.ExpectedFailOutput = append(current.ExpectedFailOutput, v)
		}
	}

	for _, itc := range iTestCases {
		out, err := tc.runIntegrationTest(dir, itc.TestName)
		if itc.ExpectedPass {
			if err != nil {
				return fmt.Errorf(`expected test "%s" PASS but it failed, check %s`, itc.TestName, filename)
			}
			tc.ExecutedTests = append(tc.ExecutedTests, []string{filename, itc.TestName})
			continue
		}

		if err == nil {
			return fmt.Errorf(`expected test "%s" FAIL but it passed, check %s`, itc.TestName, filename)
		}

		expected := strings.Join(itc.ExpectedFailOutput, "\n")
		if expected != out {
			msg := []string{
				fmt.Sprintf("expected output of test %s doesn't match, check %s", itc.TestName, filename),
			}

			want, got := findFirstMismatch(expected, out)
			msg = append(msg, fmt.Sprintf("\twant: %#v", want))
			msg = append(msg, fmt.Sprintf("\t got: %#v", got))
			return errors.New(strings.Join(msg, "\n"))
		}

		tc.ExecutedTests = append(tc.ExecutedTests, []string{filename, itc.TestName})
	}

	return nil
}

func (tc *GoldenTestCase) runIntegrationTest(dir, name string) (string, error) {
	testName := "^" + name + "$"

	cmd := exec.Command("go", "test", "-v", "-count=1", "-run", testName)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()

	output := meta.ExtractOutput(out)
	idx := strings.LastIndex(output, "\nexit status")
	if idx == -1 {
		return output, err
	}
	return output[:idx], err
}

func findFirstMismatch(want, got string) (string, string) {
	if want == got {
		return want, got
	}

	lw := len(want)
	lg := len(got)
	if lw < 30 || lg < 30 {
		return want, got
	}

	lm := min(lw, lg)
	anchor := -1
	for i := 0; i < lm; i++ {
		if want[i] != got[i] {
			anchor = i
			break
		}
	}

	if anchor == -1 {
		if lg < lw {
			return "..." + strings.TrimPrefix(want, got), "..."
		}
		return "...", "..." + strings.TrimPrefix(got, want)
	}

	w := want[:anchor+1]
	g := got[:anchor+1]
	if len(w) > 20 {
		w = "..." + w[len(w)-20:]
	}
	if len(g) > 20 {
		g = "..." + g[len(g)-20:]
	}

	if anchor != lw-1 {
		w = w + "..."
	}
	if anchor != lg-1 {
		g = g + "..."
	}
	return w, g
}
