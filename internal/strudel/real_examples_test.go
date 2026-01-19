package strudel

import (
	"testing"
)

// tests the strudel package with actual examples from strudel docs
func TestRealWorldExamples(t *testing.T) {
	tests := []struct {
		name            string
		code            string
		expectedSounds  []string
		expectedEffects []string
		expectedMusical []string
		minComplexity   int
	}{
		{
			name:            "Basic drum pattern from docs",
			code:            `s("bd sd [~ bd] sd,hh*16, misc")`,
			expectedSounds:  []string{"drums", "percussion"},
			expectedEffects: []string{},
			expectedMusical: []string{},
			minComplexity:   1,
		},
		{
			name:            "Synth with effects from docs",
			code:            `note("c2 <eb2 <g2 g1>>".fast(2)).sound("<sawtooth square triangle sine>")._scope()`,
			expectedSounds:  []string{},
			expectedEffects: []string{},
			expectedMusical: []string{"melody", "melodic"},
			minComplexity:   1,
		},
		{
			name:            "Noise hihat example",
			code:            `sound("bd*2,<white pink brown>*8").decay(.04).sustain(0)._scope()`,
			expectedSounds:  []string{"noise"},
			expectedEffects: []string{"envelope"},
			expectedMusical: []string{},
			minComplexity:   1,
		},
		{
			name:            "Filter sweep example",
			code:            `note("[c eb g <f bb>](3,8,<0 1>)".sub(12)).s("sawtooth").lpf(sine.range(300,2000).slow(16)).lpa(0.005).lpd(perlin.range(.02,.2)).lps(perlin.range(0,.5).slow(3)).lpq(sine.range(2,10).slow(32)).release(.5).lpenv(perlin.range(1,8).slow(2)).ftype('24db').room(1)`,
			expectedSounds:  []string{"synth"},
			expectedEffects: []string{"filter", "filter-envelope", "envelope", "reverb"},
			expectedMusical: []string{"melody", "melodic"},
			minComplexity:   2,
		},
		{
			name:            "TR808 drum pattern",
			code:            `s("bd sd,hh*16").bank("RolandTR808")`,
			expectedSounds:  []string{"drums", "percussion"},
			expectedEffects: []string{},
			expectedMusical: []string{},
			minComplexity:   1,
		},
		{
			name:            "Complex layered techno",
			code:            `s("bd").stack(s("hh")).stack(note("c e g").s("sawtooth")).lpf(2000).delay(0.25).room(0.5)`,
			expectedSounds:  []string{"drums", "percussion", "synth"},
			expectedEffects: []string{"filter", "delay", "reverb"},
			expectedMusical: []string{"melody", "melodic", "layered"},
			minComplexity:   3,
		},
		{
			name:            "ZZFX synth example",
			code:            `note("c2 eb2 f2 g2").s("{z_sawtooth z_tan z_noise z_sine z_square}%4").attack(0.001).decay(0.1).sustain(.8).release(.1).curve(1).slide(0)`,
			expectedSounds:  []string{"zzfx"},
			expectedEffects: []string{"envelope", "zzfx"},
			expectedMusical: []string{"melody", "melodic"},
			minComplexity:   1,
		},
		{
			name:            "Wavetable synthesis",
			code:            `note("<[g3,b3,e4]!2 [a3,c3,e4] [b3,d3,f#4]>").n("<1 2 3 4 5 6 7 8 9 10>/2").room(0.5).size(0.9).s('wt_flute').velocity(0.25).often(n => n.ply(2)).release(0.125).decay("<0.1 0.25 0.3 0.4>").sustain(0).cutoff(2000).cutoff("<1000 2000 4000>").fast(4)`,
			expectedSounds:  []string{"wavetable"},
			expectedEffects: []string{"reverb", "dynamics", "envelope"},
			expectedMusical: []string{"melody", "melodic", "chords", "harmony"},
			minComplexity:   2,
		},
		{
			name:            "Sampler effects chop",
			code:            `s("the_drum*2").chop(16).speed(rand.range(0.85,1.1))`,
			expectedSounds:  []string{},
			expectedEffects: []string{"sampler"},
			expectedMusical: []string{},
			minComplexity:   1,
		},

		{
			name: "Simple multi-layer",
			code: `
			samples('github:tidalcycles/dirt-samples')
			$: s("[bd <hh oh>]*4").bank("tr909").dec(.5)
			$: note("36 43, 52 59 62 64").sound("gm_acoustic_bass").mask("<0 0 0 0 0 0 0 0 1 1 1 1 1 1 1 1>")
			$: note("[36 -, 52 - - - ]/4 [-]*4").sound("sawtooth").lpf(2000)
			$: s("alphabet:24/4").delay("0.8")
			$: note("36 43, 52 59 62 64").sound("gm_blown_bottle").mask("<0 0 0 0 0 0 0 0 0 0 0 0 1 1 1 1>").lpf(500)`,
			expectedSounds:  []string{"drums", "percussion", "synth"},
			expectedEffects: []string{"envelope", "filter", "delay"},
			expectedMusical: []string{"melody", "melodic", "harmony"},
			minComplexity:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := AnalyzeCode(tt.code)

			// check sounds
			for _, expectedSound := range tt.expectedSounds {
				if !contains(analysis.SoundTags, expectedSound) {
					t.Errorf("Expected sound tag '%s' in %v", expectedSound, analysis.SoundTags)
				}
			}

			// check effects
			for _, expectedEffect := range tt.expectedEffects {
				if !contains(analysis.EffectTags, expectedEffect) {
					t.Errorf("Expected effect tag '%s' in %v", expectedEffect, analysis.EffectTags)
				}
			}

			// check musical elements
			for _, expectedMusical := range tt.expectedMusical {
				found := contains(analysis.MusicalTags, expectedMusical) ||
					contains(analysis.ComplexityTags, expectedMusical)
				if !found {
					t.Errorf("Expected musical/complexity tag '%s' in musical:%v or complexity:%v",
						expectedMusical, analysis.MusicalTags, analysis.ComplexityTags)
				}
			}

			// check complexity
			if analysis.Complexity < tt.minComplexity {
				t.Errorf("Expected complexity >= %d, got %d", tt.minComplexity, analysis.Complexity)
			}
		})
	}
}

// tests keyword extraction with real Strudel code
func TestKeywordExtraction(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		mustContain []string
	}{
		{
			name:        "Drum pattern",
			code:        `s("bd sd hh")`,
			mustContain: []string{"bd", "sd", "hh"},
		},
		{
			name:        "Synth with notes",
			code:        `note("c e g").sound("sawtooth")`,
			mustContain: []string{"c", "e", "g", "sawtooth", "sound"},
		},
		{
			name:        "Complex effects chain",
			code:        `s("bd").lpf(1000).delay(0.25).room(0.5).fast(2)`,
			mustContain: []string{"bd", "lpf", "delay", "room", "fast"},
		},
		{
			name:        "Noise generators",
			code:        `sound("white pink brown")`,
			mustContain: []string{"white", "pink", "brown"},
		},
		{
			name:        "ZZFX synths",
			code:        `s("z_sawtooth z_sine").attack(0.1)`,
			mustContain: []string{"z_sawtooth", "z_sine", "attack"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keywords := ExtractKeywords(tt.code)

			for _, mustContain := range tt.mustContain {
				if !stringContains(keywords, mustContain) {
					t.Errorf("Keywords '%s' should contain '%s'", keywords, mustContain)
				}
			}
		})
	}
}

func stringContains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr)) ||
		containsWord(s, substr))
}

func containsWord(s, word string) bool {
	// simple word boundary check
	for i := 0; i <= len(s)-len(word); i++ {
		if s[i:i+len(word)] == word {
			// check word boundaries
			if (i == 0 || s[i-1] == ' ') && (i+len(word) == len(s) || s[i+len(word)] == ' ') {
				return true
			}
		}
	}
	return false
}
