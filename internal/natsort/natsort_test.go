package natsort_test

import (
	"testing"

	. "github.com/dogmatiq/dapper/internal/natsort"
)

func Test_Less(t *testing.T) {
	tests := []struct {
		name     string
		a, b     string
		expected bool
	}{
		{
			"alphanumeric string 'a' sorts before 'b'",
			"alpha 1", "alpha 100",
			true,
		},
		{
			"alphanumeric string 'a' sorts after 'b'",
			"alpha 101", "alpha 100",
			false,
		},
		{
			"alpha-only string 'a' sorts after 'b'",
			"alpha beta", "alpha",
			false,
		},
		{
			"alpha-only string 'a' sorts before 'b'",
			"alpha beta", "alpha beta gamma",
			true,
		},
		{
			"numeric-only string 'a' sorts before 'b'",
			"1234.5678", "1234.56789",
			true,
		},
		{
			"numeric-only string 'a' sorts after 'b'",
			"1234.567810", "1234.56789",
			false,
		},
		{
			"shorter 'a' string sorts before 'b'",
			"alpha 100 200", "alpha 100 200 300",
			true,
		},
		{
			"longer 'a' string sorts after 'b'",
			"alpha 100 200 300", "alpha 100 200",
			false,
		},
		{
			"heterogeneous strings",
			"alpha beta", "alpha 100",
			false,
		},
		{
			"heterogeneous strings (reversed)",
			"alpha 100", "alpha beta",
			true,
		},
		{
			"equal strings",
			"alpha 100", "alpha 100",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := Less(tt.a, tt.b)

			t.Logf("natsort.Less(a:=%q, b:=%q) = %v, expected %v",
				tt.a,
				tt.b,
				actual,
				tt.expected,
			)

			if actual != tt.expected {
				t.Errorf(
					"natsort.Less(a:=%q, b:=%q) = %v, expected %v",
					tt.a,
					tt.b,
					actual,
					tt.expected,
				)
			}
		})
	}
}
