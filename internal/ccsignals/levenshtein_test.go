package ccsignals

import (
	"testing"
)

func TestLevenshteinDistance(t *testing.T) {
	tests := []struct {
		name string
		s1   string
		s2   string
		want int
	}{
		{
			name: "identical strings",
			s1:   "hello",
			s2:   "hello",
			want: 0,
		},
		{
			name: "empty strings",
			s1:   "",
			s2:   "",
			want: 0,
		},
		{
			name: "first empty",
			s1:   "",
			s2:   "hello",
			want: 5,
		},
		{
			name: "second empty",
			s1:   "hello",
			s2:   "",
			want: 5,
		},
		{
			name: "one substitution",
			s1:   "hello",
			s2:   "hallo",
			want: 1,
		},
		{
			name: "one insertion",
			s1:   "hello",
			s2:   "helloo",
			want: 1,
		},
		{
			name: "one deletion",
			s1:   "hello",
			s2:   "helo",
			want: 1,
		},
		{
			name: "kitten to sitting",
			s1:   "kitten",
			s2:   "sitting",
			want: 3,
		},
		{
			name: "saturday to sunday",
			s1:   "saturday",
			s2:   "sunday",
			want: 3,
		},
		{
			name: "completely different",
			s1:   "abc",
			s2:   "xyz",
			want: 3,
		},
		{
			name: "case sensitive",
			s1:   "Hello",
			s2:   "hello",
			want: 1,
		},
		{
			name: "longer strings",
			s1:   "the quick brown fox",
			s2:   "the quick brown dog",
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LevenshteinDistance(tt.s1, tt.s2)
			if got != tt.want {
				t.Errorf("LevenshteinDistance(%q, %q) = %d, want %d", tt.s1, tt.s2, got, tt.want)
			}
		})
	}
}

func TestLevenshteinDistance_Symmetric(t *testing.T) {
	pairs := [][2]string{
		{"hello", "world"},
		{"kitten", "sitting"},
		{"abc", ""},
		{"test", "testing"},
	}

	for _, pair := range pairs {
		s1, s2 := pair[0], pair[1]
		d1 := LevenshteinDistance(s1, s2)
		d2 := LevenshteinDistance(s2, s1)

		if d1 != d2 {
			t.Errorf("LevenshteinDistance not symmetric: (%q, %q)=%d vs (%q, %q)=%d",
				s1, s2, d1, s2, s1, d2)
		}
	}
}

func TestLevenshteinDistance_TriangleInequality(t *testing.T) {
	// for any strings a, b, c: distance(a, c) <= distance(a, b) + distance(b, c)
	a := "hello"
	b := "hallo"
	c := "world"

	dAB := LevenshteinDistance(a, b)
	dBC := LevenshteinDistance(b, c)
	dAC := LevenshteinDistance(a, c)

	if dAC > dAB+dBC {
		t.Errorf("triangle inequality violated: d(%q,%q)=%d > d(%q,%q)=%d + d(%q,%q)=%d",
			a, c, dAC, a, b, dAB, b, c, dBC)
	}
}

func TestNormalizedEditDistance(t *testing.T) {
	tests := []struct {
		name string
		s1   string
		s2   string
		want float64
	}{
		{
			name: "identical",
			s1:   "hello",
			s2:   "hello",
			want: 0.0,
		},
		{
			name: "both empty",
			s1:   "",
			s2:   "",
			want: 0.0,
		},
		{
			name: "one empty",
			s1:   "hello",
			s2:   "",
			want: 1.0,
		},
		{
			name: "completely different same length",
			s1:   "abc",
			s2:   "xyz",
			want: 1.0,
		},
		{
			name: "50% different",
			s1:   "abcd",
			s2:   "abxy",
			want: 0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizedEditDistance(tt.s1, tt.s2)
			if got != tt.want {
				t.Errorf("NormalizedEditDistance(%q, %q) = %v, want %v", tt.s1, tt.s2, got, tt.want)
			}
		})
	}
}

func TestNormalizedEditDistance_Range(t *testing.T) {
	// normalized distance should always be between 0 and 1 for strings of equal length
	pairs := [][2]string{
		{"hello", "world"},
		{"test", "best"},
		{"aaaa", "bbbb"},
		{"abcd", "abcd"},
	}

	for _, pair := range pairs {
		s1, s2 := pair[0], pair[1]
		d := NormalizedEditDistance(s1, s2)

		if d < 0 || d > 1.0 {
			t.Errorf("NormalizedEditDistance(%q, %q) = %v, expected 0 <= d <= 1", s1, s2, d)
		}
	}
}

func BenchmarkLevenshteinDistance_Short(b *testing.B) {
	s1 := "hello"
	s2 := "world"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LevenshteinDistance(s1, s2)
	}
}

func BenchmarkLevenshteinDistance_Medium(b *testing.B) {
	s1 := "the quick brown fox jumps over the lazy dog"
	s2 := "the fast brown cat leaps over the lazy rat"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LevenshteinDistance(s1, s2)
	}
}

func BenchmarkLevenshteinDistance_Long(b *testing.B) {
	s1 := string(make([]byte, 500))
	s2 := string(make([]byte, 500))
	for i := range s1 {
		s1 = s1[:i] + "a" + s1[i+1:]
	}
	for i := range s2 {
		s2 = s2[:i] + "b" + s2[i+1:]
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LevenshteinDistance(s1, s2)
	}
}
