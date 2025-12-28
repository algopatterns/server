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
			name: "Darkles by tzwaan",
			code: `
			setcpm(140/4)
			var scale = "d:iwato"
			function drumEffects(notes) {
				return notes.bank("rhythmace").distort("0.5:1").lpf(8000)
			}
			const KICK = drumEffects(sound("bd").struct("1*<1 2> 0 0 0 0 [[1 0] | [1 0 1 1]]@2 0"))
				.duck("2:3").duckattack(.1).duckdepth(.3).color("cyan")._punchcard()
			const HIGHHAT = drumEffects(sound("hh").struct("1 1*[2|1] 1 1*[2|4]@2 1 1*[3|1] 1"))
				.color("yellow")._punchcard()
			const SNARE = drumEffects(sound("sd").struct("0 0 1 0 0 0 <1 [0 1]> <0 0 [0 0 1]>"))
				.color("green")._punchcard()
			var bassNotes = n(irand(8).seg(2).sub("10").rib(2, 4))
			const BASS = bassNotes.sound("triangle").distort("<4:.3 3:.6>")
				.attack("[<0.2 0.1> 0 .5]/2").scale(scale).lpf("200").orbit(3)._punchcard()
			const MELODY = stack(
				"- - [g5 eb6 - eb5] [eb6 - g5 d6] - [- c6 d6 -] - -".slow(2).mask("<0 1 0 0>".slow(2)),
				"[c6 - - c6] - - [- - d6 -] [- eb6 - -] - - -".slow(2).mask("<0 0 0 1>".slow(2)),
			).note().sound("saw").distort("2:.4").echo(5, 1/8, .6)
			$: arrange([4, stack(SYNTH)], [2, stack(SYNTH, MOD_SPARKLES, MELODY)]).punchcard({ labels: true })`,
			expectedSounds:  []string{"drums", "percussion", "synth"},
			expectedEffects: []string{"distortion", "filter", "sidechain", "envelope"},
			expectedMusical: []string{"scales", "melodic", "rhythm", "layered", "arranged"},
			minComplexity:   5,
		},
		{
			name: "Memodries by Nicop",
			code: `
			setcps(.5)
			samples('github:switchangel/pad')
			samples('github:switchangel/breaks')
			$: note("[[c1 [eb1 bb0]]!3 [Ab0 [Bb0 G0]]]/8").add(note("12")).s("swpad:1").begin(.05).end(.15).loopb(.12).loop(1).loope(.15).att(.3).rel(.8).gain(.8)
			$: s("breaks:3").end(.5).fit().hpf(200).gain(.5).orbit(2)
			$: note("[[Ab1 [B1 F#1]]!3 [E1 [F#1 Eb1]]]/8").struct("<[x!2 ~@6] [[x!2 ~@2]!2]>").s("sawtooth").lpf(100).lpenv(4).gain(1.2)
			$: s("[~ sd:3]*2").begin(.0).gain(.6).dec(.2).orbit(2)
			$: s("bd:0").beat("0,5", 8).gain(.5).hpf(60).duckorbit(1).duckdepth(.8).duckattack(.125).orbit(2)
			$: n("0 .. 15").palindrome().chord("[[Ab2 [B2 F#2]]!3 [E2 [F#2 Eb2]]]/8").anchor("<G5!6 G4 B4>").voicing().s("sawtooth").lpf(tri.slow(2).rangex(600,3000)).lpq(.2).rel(0).gain(.4)`,
			expectedSounds:  []string{"drums", "percussion", "synth"},
			expectedEffects: []string{"filter", "filter-envelope", "sampler", "sidechain"},
			expectedMusical: []string{"melody", "melodic", "rhythm"},
			minComplexity:   3,
		},
		{
			name: "Coastline by eddyflux",
			code: `
			samples('github:eddyflux/crate')
			setcps(.75)
			let chords = chord("<Bbm9 Fm9>/4").dict('ireal')
			stack(
				stack(s("bd").struct("<[x*<1 2> [~@3 x]] x>"), s("~ [rim, sd:<2 3>]").room("<0 .2>"),
					n("[0 <1 3>]*<2!3 4>").s("hh"), s("rd:<1!3 2>*2").mask("<0 0 1 1>/16").gain(.5)
				).bank('crate').mask("<[0 1] 1 1 1>/16".early(.5)),
				chords.offset(-1).voicing().s("gm_epiano1:1").phaser(4).room(.5),
				n("<0!3 1*2>").set(chords).mode("root:g2").voicing().s("gm_acoustic_bass"),
				chords.n("[0 <4 3 <2 5>>*2](<3 5>,8)").anchor("D5").voicing().segment(4)
					.clip(rand.range(.4,.8)).room(.75).shape(.3).delay(.25).fm(sine.range(3,8).slow(8))
					.lpf(sine.range(500,1000).slow(8)).lpq(5).rarely(ply("2")).chunk(4, fast(2))
					.gain(perlin.range(.6, .9)).mask("<0 1 1 0>/16")
			).late("[0 .01]*4").late("[0 .01]*2").size(4)`,
			expectedSounds:  []string{"drums", "percussion"},
			expectedEffects: []string{"reverb", "modulation", "delay", "filter", "distortion", "fm-synthesis"},
			expectedMusical: []string{"rhythm", "layered"},
			minComplexity:   6,
		},
		{
			name: "Heliyatrel",
			code: `
			setCpm(135/4)
			$: note("<[f1 f2]*4 [ds1 ds2]*4 [db1 db2]*4 [c1 c2 c1 c2 d1 d2 e1 e2]>").sound("wt_digital").lpf("200 400").gain("2")
			$: sound("circuitsdrumtracks_bd*4").gain("0.75")
			$: sound("[~ circuitsdrumtracks_sd]*2").gain("0.75")
			$: sound("circuitsdrumtracks_hh*8").gain("[0.05 0.2]*4")
			$: sound("[~@12 circuitsdrumtracks_sd]").gain("0.5")
			$: note("[[f1 ~@1]!4 [g#1 ~@1]!4 [bb1 ~@1]!4 [c2 ~@1]!4]*0.25")
				.sound("supersaw").gain("0").room(1.5).delay("0.75").lpf("600 800 400 1000")
			$: note("[~ [<[f3,ab3,c4,eb4]!2 [f3,ab3,c4,d4]!2 [f3,ab3,c4,eb4]!2 [e3,f3,a3,c4,e4]!2> ~]]*2")
				.sound("gm_electric_guitar_muted").delay("0.8:0.6:0.5").gain("0.5")`,
			expectedSounds:  []string{"wavetable"},
			expectedEffects: []string{"filter", "delay", "reverb"},
			expectedMusical: []string{"melody", "melodic"},
			minComplexity:   3,
		},
		{
			name: "Music Theory by quteriss",
			code: `
			$: stack(
				sound("clave:5").struct("[1 0 0 0 0 0 1 0]").gain(".8 .6"),
				note("<<fs3,as3,cs4, fs2> <e3,gs3,b4, e2>>")
					.scale("b:major").sound("gm_epiano1:2").room(.6).delay(.3).gain(.6),
				note("<[<- as4> gs4]*2 [[fs4 cs4]/2 e4]*2>")
					.sound("<gm_ocarina, gm_epiano1:4>").delay(.4).room(.2).seg("<2 4 2 2>")
					.fast("<<0 1> 1 2 1>/2")
			);`,
			expectedSounds:  []string{},
			expectedEffects: []string{"reverb", "delay", "dynamics"},
			expectedMusical: []string{"melody", "melodic", "scales", "layered"},
			minComplexity:   4,
		},
		{
			name: "DnB by onefeather",
			code: `
			samples('github:0nefeather/akas-dnb-essentials')
			setGainCurve(x => Math.pow(x, 2))
			setCps(170/60/4)
			const BIG_GAIN = slider(1)

			$: s("hats:3!8").decay(.5).slow(2).color("mediumslateblue")
			  .almostNever(x=>x.ply("2 | 4").color("lime"))
			  .duck(2).gain(slider(0.739).mul(BIG_GAIN))
			  ._pianoroll({ playhead: 1, hideNegative: true })

			$: s("drum_loops/2:8").fit().color("mediumslateblue")
			  .scrub(irand(16).div(16).seg(8).rib("<12.25 32>", 2))
			  .sometimesBy(0.05, x=>x.ply("2 | 3").color("lime"))
			  .gain(slider(0.625).mul(BIG_GAIN))._scope()

			const VOX_GAIN = slider(0.813)
			$: s("pads_drones").scrub("0 .42 .14 .69").n("6").note("36")
			  .color("mediumslateblue")
			  .sometimesBy(slider(0.498), x=>x.ply("2 | 4 | 8").color("lime"))
			  .superimpose(x=>x.note("60"))
			  .gain(VOX_GAIN.mul(BIG_GAIN))._scope()

			const chops = [
			  { degrade: 0, chop: ".8@3 .1" },
			  { degrade: 0, chop: ".16!3 .7@5" },
			  { degrade: 0.25, chop: rand.seg(8) }
			]
			var ch = 2

			$: s("pads_drones").scrub(pick(chops.map(c => c.chop), ch)).n("3")
			  .degradeBy(pick(chops.map(c => c.degrade), ch))
			  .superimpose(x=>x.note("48")).rib(68, 2).color("mediumslateblue")
			  .sometimesBy(slider(0.2, 0, 1, 0.05), x=>x.ply("2 | 4").note("72").hpf(4000).decay(0.1).color("lime"))
			  .gain(slider(0.511).mul(BIG_GAIN))._pianoroll({ playhead: 1, hideNegative: true })

			$: note("<f1!8 c1!8>").s("tri").orbit(2).color("mediumslateblue")
			  .gain(slider(0.444, 0, 3).mul(BIG_GAIN))._pianoroll()

			$: note("<f2!8 c2!8>").s("sawtooth")
			  .lpf(sine.range(200, 15000).fast(24))
			  .gain(slider(1.022, 0, 2).mul(BIG_GAIN))

			$: note("<f2!8 c2!8, f3!8 c3!8>").s("z_triangle").orbit(2)
			  .color("mediumslateblue").gain(slider(1.125, 0, 3).mul(BIG_GAIN))
			  .phaser(2).phasercenter(800)
			  ._spiral({ steady: .96, colorizeInactive: true, inset: 2, cap: "round", logSpiral: true, playheadColor: "lime"})

			$: n(irand(16).add(8)).struct("x*8").s("z_sine")
			  .degradeBy(0.45).rib("13 | 7", 2).scale("<F:dorian!4>")
			  .echo(4, 1/16, .8).decay(0.1).room(0.9).pan(sine.slow(2))
			  .gain(slider(0.784).mul(BIG_GAIN)).color("mediumslateblue")
			  ._pianoroll({playhead: 0, playheadColor: "mediumslateblue", flipTime: true})

			$: note("<f1!8 c1!8>").s("tri").gain(1.75).orbit(2).gain(0)`,
			expectedSounds:  []string{"synth", "zzfx"},
			expectedEffects: []string{"envelope", "sidechain", "sampler", "filter", "modulation", "delay", "reverb", "spatial"},
			expectedMusical: []string{"melody", "melodic", "scales", "rhythm", "interactive"},
			minComplexity:   5,
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
