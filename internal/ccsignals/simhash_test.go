package ccsignals

import (
	"testing"
)

func TestSimHasher_Hash(t *testing.T) {
	hasher := NewSimHasher(3)

	tests := []struct {
		name    string
		content string
		wantNonZero bool
	}{
		{
			name:    "normal text",
			content: "the quick brown fox jumps over the lazy dog",
			wantNonZero: true,
		},
		{
			name:    "empty string",
			content: "",
			wantNonZero: false,
		},
		{
			name:    "only whitespace",
			content: "   \t\n  ",
			wantNonZero: false,
		},
		{
			name:    "only punctuation",
			content: "!@#$%^&*()",
			wantNonZero: false,
		},
		{
			name:    "single word",
			content: "hello",
			wantNonZero: true,
		},
		{
			name:    "two words",
			content: "hello world",
			wantNonZero: true,
		},
		{
			name:    "code snippet",
			content: "function foo() { return 42; }",
			wantNonZero: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := hasher.Hash(tt.content)
			if tt.wantNonZero && fp == 0 {
				t.Error("expected non-zero fingerprint")
			}
			if !tt.wantNonZero && fp != 0 {
				t.Errorf("expected zero fingerprint, got %x", fp)
			}
		})
	}
}

func TestSimHasher_Deterministic(t *testing.T) {
	hasher := NewSimHasher(3)
	content := "the quick brown fox jumps over the lazy dog"

	fp1 := hasher.Hash(content)
	fp2 := hasher.Hash(content)

	if fp1 != fp2 {
		t.Errorf("hash not deterministic: %x != %x", fp1, fp2)
	}
}

func TestSimHasher_SimilarContent(t *testing.T) {
	hasher := NewSimHasher(3)

	original := "the quick brown fox jumps over the lazy dog"
	similar := "the quick brown fox leaps over the lazy dog"
	different := "completely unrelated content about programming"

	fpOriginal := hasher.Hash(original)
	fpSimilar := hasher.Hash(similar)
	fpDifferent := hasher.Hash(different)

	distSimilar := HammingDistance(fpOriginal, fpSimilar)
	distDifferent := HammingDistance(fpOriginal, fpDifferent)

	if distSimilar >= distDifferent {
		t.Errorf("similar content should have smaller distance: similar=%d, different=%d",
			distSimilar, distDifferent)
	}
}

func TestSimHasher_CaseInsensitive(t *testing.T) {
	hasher := NewSimHasher(3)

	lower := "hello world test"
	upper := "HELLO WORLD TEST"
	mixed := "HeLLo WoRLd TeST"

	fpLower := hasher.Hash(lower)
	fpUpper := hasher.Hash(upper)
	fpMixed := hasher.Hash(mixed)

	if fpLower != fpUpper || fpLower != fpMixed {
		t.Error("expected case insensitive hashing")
	}
}

func TestSimHasher_IgnoresPunctuation(t *testing.T) {
	hasher := NewSimHasher(3)

	plain := "hello world test"
	punctuated := "hello, world! test..."

	fpPlain := hasher.Hash(plain)
	fpPunctuated := hasher.Hash(punctuated)

	if fpPlain != fpPunctuated {
		t.Error("expected punctuation to be ignored")
	}
}

func TestSimHasher_ShingleSize(t *testing.T) {
	// different shingle sizes should produce different hashes
	hasher2 := NewSimHasher(2)
	hasher3 := NewSimHasher(3)
	hasher5 := NewSimHasher(5)

	content := "one two three four five six seven eight"

	fp2 := hasher2.Hash(content)
	fp3 := hasher3.Hash(content)
	fp5 := hasher5.Hash(content)

	// they should all be non-zero
	if fp2 == 0 || fp3 == 0 || fp5 == 0 {
		t.Error("expected all fingerprints to be non-zero")
	}

	// at least some should be different
	if fp2 == fp3 && fp3 == fp5 {
		t.Error("expected different shingle sizes to produce different hashes")
	}
}

func TestSimHasher_InvalidShingleSize(t *testing.T) {
	hasher := NewSimHasher(0)
	if hasher.shingleSize != DefaultShingleSize {
		t.Errorf("expected default shingle size %d, got %d", DefaultShingleSize, hasher.shingleSize)
	}

	hasher = NewSimHasher(-5)
	if hasher.shingleSize != DefaultShingleSize {
		t.Errorf("expected default shingle size %d, got %d", DefaultShingleSize, hasher.shingleSize)
	}
}

func TestHammingDistance(t *testing.T) {
	tests := []struct {
		name string
		a    Fingerprint
		b    Fingerprint
		want int
	}{
		{
			name: "identical",
			a:    0xFFFFFFFFFFFFFFFF,
			b:    0xFFFFFFFFFFFFFFFF,
			want: 0,
		},
		{
			name: "one bit different",
			a:    0xFFFFFFFFFFFFFFFF,
			b:    0xFFFFFFFFFFFFFFFE,
			want: 1,
		},
		{
			name: "all bits different",
			a:    0xFFFFFFFFFFFFFFFF,
			b:    0x0000000000000000,
			want: 64,
		},
		{
			name: "half bits different",
			a:    0xFFFFFFFF00000000,
			b:    0x00000000FFFFFFFF,
			want: 64,
		},
		{
			name: "alternating bits",
			a:    0xAAAAAAAAAAAAAAAA,
			b:    0x5555555555555555,
			want: 64,
		},
		{
			name: "8 bits different",
			a:    0xFFFFFFFFFFFFFFFF,
			b:    0xFFFFFFFFFFFFFF00,
			want: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HammingDistance(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("HammingDistance(%x, %x) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestHammingDistance_Symmetric(t *testing.T) {
	a := Fingerprint(0xABCDEF1234567890)
	b := Fingerprint(0x1234567890ABCDEF)

	distAB := HammingDistance(a, b)
	distBA := HammingDistance(b, a)

	if distAB != distBA {
		t.Errorf("HammingDistance not symmetric: %d != %d", distAB, distBA)
	}
}

func TestIsSimilar(t *testing.T) {
	tests := []struct {
		name      string
		a         Fingerprint
		b         Fingerprint
		threshold int
		want      bool
	}{
		{
			name:      "identical within threshold",
			a:         0xFFFF,
			b:         0xFFFF,
			threshold: 5,
			want:      true,
		},
		{
			name:      "at threshold",
			a:         0xFFFFFFFFFFFFFFFF,
			b:         0xFFFFFFFFFFFFFFE0, // 5 bits different
			threshold: 5,
			want:      true,
		},
		{
			name:      "above threshold",
			a:         0xFFFFFFFFFFFFFFFF,
			b:         0xFFFFFFFFFFFFFFC0, // 6 bits different
			threshold: 5,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSimilar(tt.a, tt.b, tt.threshold)
			if got != tt.want {
				t.Errorf("IsSimilar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalizeText(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Hello World", "hello world"},
		{"UPPERCASE", "uppercase"},
		{"with  multiple   spaces", "with multiple spaces"},
		{"  leading and trailing  ", "leading and trailing"},
		{"with\ttabs\nand\nnewlines", "with tabs and newlines"},
		{"punctuation!@#$%test", "punctuation test"},
		{"numbers123test", "numbers123test"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := normalizeText(tt.input)
			if got != tt.want {
				t.Errorf("normalizeText(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestTokenize(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"", 0},
		{"hello", 1},
		{"hello world", 2},
		{"one two three four five", 5},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tokens := tokenize(tt.input)
			if len(tokens) != tt.want {
				t.Errorf("tokenize(%q) returned %d tokens, want %d", tt.input, len(tokens), tt.want)
			}
		})
	}
}

func BenchmarkSimHash(b *testing.B) {
	hasher := NewSimHasher(3)
	content := "the quick brown fox jumps over the lazy dog repeatedly many times"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hasher.Hash(content)
	}
}

func BenchmarkHammingDistance(b *testing.B) {
	a := Fingerprint(0xABCDEF1234567890)
	fp := Fingerprint(0x1234567890ABCDEF)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HammingDistance(a, fp)
	}
}
