package alias

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRandomString(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{
			name: "size = 5",
			size: 5,
		},
		{
			name: "size = 10",
			size: 10,
		},
		{
			name: "size = 20",
			size: 20,
		},
		{
			name: "size = 30",
			size: 30,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			str1 := GenerateAlias(tc.size)
			str2 := GenerateAlias(tc.size)

			assert.Len(t, str1, tc.size)
			assert.Len(t, str2, tc.size)

			// Check that two generated strings are different
			// This is not an absolute guarantee that the function works correctly,
			// but this is a good heuristic for a simple random generator.
			assert.NotEqual(t, str1, str2)
		})
	}
}
