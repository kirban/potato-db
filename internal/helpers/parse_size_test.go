package helpers

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseSize(t *testing.T) {
	tests := map[string]struct {
		input       string
		want        int
		expectedErr error
	}{
		"int bytes": {
			input: "5b",
			want:  5,
		},
		"int kbytes": {
			input: "4kb",
			want:  4 << 10,
		},
		"int mbytes": {
			input: "13MB",
			want:  13 << 20,
		},
		"int gbytes": {
			input: "1gb",
			want:  1 << 30,
		},
		"float bytes": {
			input:       "1.5b",
			want:        0,
			expectedErr: errors.New("wrong size"),
		},
		"float kbytes": {
			input: "1.5kB",
			want:  1.5 * (1 << 10),
		},
		"invalid unit": {
			input:       "1QB",
			want:        0,
			expectedErr: errors.New("invalid unit"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ParseSize(tc.input)

			if tc.expectedErr != nil {
				assert.Errorf(t, err, tc.expectedErr.Error())
				//assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}
