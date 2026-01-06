package contribution

import (
	"strings"
)

// calculates AI contribution score for a strudel based on conversation history
// returns a score from 0.0 to 1.0 indicating how much of the final code came from AI
func Calculate(finalCode string, aiCodeResponses []string) float64 {
	if len(finalCode) == 0 || len(aiCodeResponses) == 0 {
		return 0.0
	}

	// normalize code for comparison
	finalCode = normalizeCode(finalCode)
	if len(finalCode) == 0 {
		return 0.0
	}

	// find the maximum overlap between any AI response and the final code
	var maxOverlap int
	for _, aiCode := range aiCodeResponses {
		aiCode = normalizeCode(aiCode)
		if len(aiCode) == 0 {
			continue
		}

		overlap := lcsLength(finalCode, aiCode)
		if overlap > maxOverlap {
			maxOverlap = overlap
		}
	}

	// calculate score as ratio of overlap to final code length
	score := float64(maxOverlap) / float64(len(finalCode))

	// clamp to 0.0-1.0
	if score > 1.0 {
		score = 1.0
	}
	if score < 0.0 {
		score = 0.0
	}

	return score
}

// normalizes code for comparison by removing insignificant whitespace
func normalizeCode(code string) string {
	// remove leading/trailing whitespace from each line
	lines := strings.Split(code, "\n")
	var normalized []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			normalized = append(normalized, trimmed)
		}
	}

	return strings.Join(normalized, "\n")
}

// calculates the length of the longest common subsequence between two strings
// uses space-optimized DP approach
func lcsLength(a, b string) int {
	m, n := len(a), len(b)
	if m == 0 || n == 0 {
		return 0
	}

	// use two rows instead of full matrix for space efficiency
	prev := make([]int, n+1)
	curr := make([]int, n+1)

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				curr[j] = prev[j-1] + 1
			} else {
				curr[j] = max(prev[j], curr[j-1])
			}
		}
		prev, curr = curr, prev
	}

	return prev[n]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
