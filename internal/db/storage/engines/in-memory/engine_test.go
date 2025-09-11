package inmemory

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestNewInMemoryEngine(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		logger         *zap.Logger
		expectedErr    error
		expectedNilObj bool
	}{
		"create engine without logger": {
			expectedErr:    ErrInvalidLogger,
			expectedNilObj: true,
		},
		"create engine": {
			logger:      zap.NewNop(),
			expectedErr: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			engine, err := NewInMemoryEngine(tc.logger)
			assert.Equal(t, tc.expectedErr, err)

			if tc.expectedNilObj {
				assert.Nil(t, engine)
			} else {
				assert.NotNil(t, engine)
			}
		})
	}
}
