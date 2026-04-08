package mockgentest

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	genlib "nhatp.com/go/gen-lib"
	"nhatp.com/go/gen-lib/file"
	mockgen "nhatp.com/go/mock-gen"
)

type GoldenTestCase struct {
	Source      []file.File
	GoldenFiles []file.File
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
	chainer := mockgen.New(fm, options...)

	pkgs, err := genlib.LoadPackages(dir)
	if err != nil {
		return nil, err
	}

	for _, pkg := range pkgs {
		if err := chainer.Generate(pkg, configs); err != nil {
			return nil, err
		}
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
