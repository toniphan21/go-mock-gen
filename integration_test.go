package mockgen

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func generateCodeForIntegrationTest(t *testing.T) {
	configs := []Config{
		{
			PackagePath: "nhatp.com/go/mock-gen/testdata/integration",
			Output: Output{
				TestFileName: "gen_test.go",
			},
			InterfaceName: "Repository",
		},
	
		{
			PackagePath: "nhatp.com/go/mock-gen/testdata/integration",
			Output: Output{
				TestFileName: "gen_test.go",
			},
			InterfaceName: "Service",
			SkipExpect:    true,
		},
	}

	_, filename, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(filename)
	dir := filepath.Join(testDir, "testdata", "integration")
	pkgs, err := LoadPackages(dir)
	require.NoError(t, err)

	fm := NewFileManager(dir, WithBinaryName(BinaryFullName), WithVersion(BinaryVersion))
	mocker := New(fm)
	for _, pkg := range pkgs {
		err = mocker.Generate(pkg, configs)
		require.NoError(t, err)
	}

	for _, v := range fm.Files() {
		err = os.WriteFile(v.FullPath, []byte(v.Content()), 0644)
		require.NoError(t, err)
	}
}

func Test_IntegrationTest(t *testing.T) {
	generateCodeForIntegrationTest(t)
}
