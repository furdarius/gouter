package gouter

import (
	"testing"
)

func TestMin(t *testing.T) {
	tests := []struct {
		a, b, want int
	}{
		{
			a:    5,
			b:    2,
			want: 2,
		},
		{
			a:    0,
			b:    0,
			want: 0,
		},
		{
			a:    10,
			b:    2,
			want: 2,
		},
		{
			a:    -123,
			b:    634,
			want: -123,
		},
		{
			a:    1000000000000,
			b:    10000000000000,
			want: 1000000000000,
		},
	}

	for _, test := range tests {
		res := min(test.a, test.b)

		if res != test.want {
			t.Fatalf("min(%d, %d) got %d, want %d", test.a, test.b, res, test.want)
		}
	}
}

func TestLcp(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{
			a:    "wef",
			b:    "grgrer",
			want: 0,
		},
		{
			a:    "/",
			b:    "/",
			want: 1,
		},
		{
			a:    "/popup",
			b:    "/popdup",
			want: 4,
		},
		{
			a:    "",
			b:    "",
			want: 0,
		},
		{
			a:    "123123123123123",
			b:    "123123123123123",
			want: 15,
		},
	}

	for _, test := range tests {
		res := lcp(test.a, test.b)

		if res != test.want {
			t.Fatalf("lcp(\"%s\", \"%s\") got %d, want %d", test.a, test.b, res, test.want)
		}
	}
}

func TestFindFirstParam(t *testing.T) {
	const symbol = ':'

	type Result struct {
		pos    int
		length int
	}

	type Test struct {
		path string
		want Result
	}

	tests := []Test{
		{
			path: "/test",
			want: Result{
				pos:    -1,
				length: 0,
			},
		},
		{
			path: "/test/:id",
			want: Result{
				pos:    6,
				length: 3,
			},
		},
		{
			path: "/test/it/:paramname",
			want: Result{
				pos:    9,
				length: 10,
			},
		},
	}

	for _, test := range tests {
		pos, length := findFirstParam(test.path, symbol)

		if pos != test.want.pos || length != test.want.length {
			t.Fatalf("findFirstParam(\"%s\", \"%v\") got (%d, %d), want (%d, %d)", test.path, symbol, pos, length, test.want.pos, test.want.length)
		}
	}
}
