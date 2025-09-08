package compute

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		inputQuery    string
		expectedQuery *Query
		expectedErr   error
	}{
		"set query": {
			inputQuery:    "SET foo value",
			expectedQuery: NewQuery(SetCommand, []string{"foo", "value"}),
			expectedErr:   nil,
		},
		"get query": {
			inputQuery:    "GET foo",
			expectedQuery: NewQuery(GetCommand, []string{"foo"}),
			expectedErr:   nil,
		},
		"del query": {
			inputQuery:    "DEL foo",
			expectedQuery: NewQuery(DelCommand, []string{"foo"}),
			expectedErr:   nil,
		},
		"empty query": {
			inputQuery:    "",
			expectedQuery: nil,
			expectedErr:   ErrInvalidQuery,
		},
		"invalid command (wrong register)": {
			inputQuery:    "set foo value",
			expectedQuery: nil,
			expectedErr:   ErrUnknownCommand,
		},
		"invalid command": {
			inputQuery:    "СЕТ ФУ ВЭЛЬЮ",
			expectedQuery: nil,
			expectedErr:   ErrUnknownCommand,
		},
		"invalid n of args of GET": {
			inputQuery:    "GET foo bar baz",
			expectedQuery: nil,
			expectedErr:   ErrWrongNOfArgs,
		},
		"invalid n of args of SET": {
			inputQuery:    "SET foo",
			expectedQuery: nil,
			expectedErr:   ErrWrongNOfArgs,
		},
		"invalid n of args of DEL": {
			inputQuery:    "DEL baz bar boo faz",
			expectedQuery: nil,
			expectedErr:   ErrWrongNOfArgs,
		},
		"invalid query with empty symbols": {
			inputQuery:    "      ",
			expectedQuery: nil,
			expectedErr:   ErrInvalidQuery,
		},
	}

	p := NewQueryParser(zap.NewNop())

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			q, err := p.Parse(tc.inputQuery)

			assert.True(t, reflect.DeepEqual(tc.expectedQuery, q))
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
