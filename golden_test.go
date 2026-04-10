package mockgen_test

import (
	"embed"
	"testing"

	"github.com/stretchr/testify/require"
	"nhatp.com/go/gen-lib/gentest"
	"nhatp.com/go/mock-gen/mockgentest"
)

//go:embed features/*.md testdata/*.md
var goldenMarkdownFiles embed.FS

func TestGoldenFiles(t *testing.T) {
	gentest.RunEmbedGoldenFiles(t, goldenMarkdownFiles, func(testCase gentest.MarkdownTestCase) {
		runGoldenTest(t, testCase)
	})
}

func TestGoldenFiles_Dev(t *testing.T) {
	cases := []struct {
		file string
	}{
		{file: "testdata/integration-test.md"},
	}

	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			gentest.RunEmbedGoldenFile(t, goldenMarkdownFiles, tc.file, func(testCase gentest.MarkdownTestCase) {
				runGoldenTest(t, testCase)
			})
		})
	}
}

func runGoldenTest(t *testing.T, tc gentest.MarkdownTestCase) {
	gtc := mockgentest.GoldenTestCase{
		Source:      tc.SourceFiles,
		GoldenFiles: tc.GoldenFiles,
	}

	dir := t.TempDir()
	err := gtc.Setup(dir)
	require.NoError(t, err)

	cf, err := gtc.ParseConfigs(dir)
	require.NoError(t, err)

	gtc.Run(t, dir, cf)
}
