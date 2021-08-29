package main

import (

	"testing"
	"github.com/stretchr/testify/require"
)

func TestCase(t *testing.T) {
	type want struct {
		statusCode int
	}
	type params struct {
		statusCode int
	}
	tests := []struct {
		name   string
		want   want
		params params
	}{
		{
			name: "Success CASE",
			want: want{
				statusCode: 200,
			},
			params: params{
				statusCode: 200,
			},
		},
		{
			name: "Fail CASE",
			want: want{
				statusCode: 400,
			},
			params: params{
				statusCode: 200,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		require.Equal(t, tt.want.statusCode, tt.params.statusCode, tt.name)
	}
}
