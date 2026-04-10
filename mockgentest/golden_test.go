package mockgentest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_findFirstMismatch(t *testing.T) {
	cases := []struct {
		name         string
		want         string
		got          string
		expectedWant string
		expectedGot  string
	}{
		{
			name:         "returns both if they are less than 30 chars",
			want:         "source_test.go:10: Service",
			got:          "source_test.go:11: Service",
			expectedWant: "source_test.go:10: Service",
			expectedGot:  "source_test.go:11: Service",
		},

		{
			name:         "returns both if one of them less than 30 chars - want",
			want:         "source_test.go:10: ",
			got:          "source_test.go:10: Service was not created and failed",
			expectedWant: "source_test.go:10: ",
			expectedGot:  "source_test.go:10: Service was not created and failed",
		},

		{
			name:         "returns both if one of them less than 30 chars - got",
			want:         "source_test.go:10: Service was not created and failed",
			got:          "source_test.go:10: ",
			expectedWant: "source_test.go:10: Service was not created and failed",
			expectedGot:  "source_test.go:10: ",
		},

		{
			name:         "returns both if they are equal",
			want:         "source_test.go:10: Service was not created and failed",
			got:          "source_test.go:10: Service was not created and failed",
			expectedWant: "source_test.go:10: Service was not created and failed",
			expectedGot:  "source_test.go:10: Service was not created and failed",
		},

		{
			name:         "cut 20 chars before the diff point - case right in the beginning",
			want:         "sOurce_test.go:10: Service was not created and failed",
			got:          "source_test.go:10: Service was not created and failed",
			expectedWant: "sO...",
			expectedGot:  "so...",
		},

		{
			name:         "cut 20 chars before the diff point - case in the beginning",
			want:         "source_test.go:9: Service was not created and failed",
			got:          "source_test.go:10: Service was not created and failed",
			expectedWant: "source_test.go:9...",
			expectedGot:  "source_test.go:1...",
		},

		{
			name:         "cut 20 chars before the diff point - case in the middle",
			want:         "source_test.go:10: Service was nt created and failed",
			got:          "source_test.go:10: Service was not created and failed",
			expectedWant: "...o:10: Service was nt...",
			expectedGot:  "...o:10: Service was no...",
		},

		{
			name:         "cut 20 chars before the diff point - case close to the end",
			want:         "source_test.go:10: Service was not created and fAiled",
			got:          "source_test.go:10: Service was not created and failed",
			expectedWant: "...s not created and fA...",
			expectedGot:  "...s not created and fa...",
		},

		{
			name:         "cut 20 chars before the diff point - case right at the end",
			want:         "source_test.go:10: Service was not created and failed!",
			got:          "source_test.go:10: Service was not created and failed#",
			expectedWant: "... created and failed!",
			expectedGot:  "... created and failed#",
		},

		{
			name:         "mismatch length",
			want:         "source_test.go:10: Service was not created and failed then something happens",
			got:          "source_test.go:10: Service was not created and failed rather than passed",
			expectedWant: "...created and failed t...",
			expectedGot:  "...created and failed r...",
		},

		{
			name:         "mismatch length make anchor -1",
			want:         "source_test.go:10: Service was not created and failed then something happens",
			got:          "source_test.go:10: Service was not created and failed",
			expectedWant: "... then something happens",
			expectedGot:  "...",
		},

		{
			name:         "mismatch length make anchor -1",
			want:         "source_test.go:10: Service was not created and failed",
			got:          "source_test.go:10: Service was not created and failed then something happens",
			expectedWant: "...",
			expectedGot:  "... then something happens",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ew, eg := findFirstMismatch(tc.want, tc.got)

			assert.Equal(t, tc.expectedWant, ew)
			assert.Equal(t, tc.expectedGot, eg)
		})
	}
}
