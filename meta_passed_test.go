package mockgen

import (
	"os/exec"
	"strings"
	"testing"
)

func Test_Meta_Passed(t *testing.T) {
	cases := []struct {
		name string
		test string
	}{
		{name: "Test_STUB_Via_Struct"},
		{name: "Test_STUB_Via_Ctor"},
		{name: "Test_STUB_Call_Has_Location"},
		{name: "Test_EXPECT_Via_Struct"},
		{name: "Test_EXPECT_Via_Ctor"},
		{name: "Test_EXPECT_Partial_Arg"},
		{name: "Test_EXPECT_Return"},
		{name: "Test_SubTests"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			testDir := "./testdata/meta"
			want := tc.test
			testName := tc.test
			if testName == "" {
				testName = "^" + tc.name + "$"
				want = tc.name
			}

			cmd := exec.Command("go", "test", "-v", "-count=1", testDir, "-run", testName)
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatal("Expected the meta-test to pass, but it failed!")
			}

			lines := strings.Split(string(out), "\n")
			isPassed := false
			for _, line := range lines {
				if strings.HasPrefix(line, "--- PASS:") {
					cut := strings.TrimSpace(strings.TrimPrefix(line, "--- PASS:"))
					if strings.HasPrefix(cut, want) {
						isPassed = true
						break
					}
				}
			}

			if !isPassed {
				t.Fatal("Expected the meta-test to pass, but it failed!")
			}
		})
	}
}
