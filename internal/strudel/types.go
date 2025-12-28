package strudel

type CodeAnalysis struct {
	SoundTags      []string // ["drums", "synth", "bass"]
	EffectTags     []string // ["delay", "reverb", "filter"]
	MusicalTags    []string // ["melody", "chords", "rhythm"]
	ComplexityTags []string // ["layered", "advanced", "simple"]

	// metrics
	Complexity    int // 0-10 score
	LineCount     int
	FunctionCount int
	VariableCount int
}

type SoundDefinitions struct {
	Drums      []string
	Percussion []string
	Synth      []string
	Noise      []string
	ZZFX       []string
	Wavetable  []string
	Misc       []string
	Custom     []string
}

type EffectDefinitions struct {
	Filter         []string
	FilterEnvelope []string
	Distortion     []string
	Dynamics       []string
	Spatial        []string
	Delay          []string
	Reverb         []string
	Modulation     []string
	Envelope       []string
	PitchEnvelope  []string
	FMSynthesis    []string
	Sampler        []string
	Routing        []string
	Sidechain      []string
	Synthesis      []string
	ZZFX           []string
}

type ParsedCode struct {
	Sounds    []string       // sound sample names: ["bd", "hh", "sd"]
	Notes     []string       // note names: ["c", "e", "g"]
	Functions []string       // function names: ["fast", "slow", "stack"]
	Variables []string       // variable names: ["pat1", "rhythm"]
	Scales    []string       // scale/mode names: ["minor", "dorian"]
	Patterns  map[string]int // pattern counts: {"stack": 2, "arrange": 1}
}

type KeywordOptions struct {
	MaxKeywords      int  // limit total keywords (default: 10)
	IncludeSounds    bool // include sound names (default: true)
	IncludeNotes     bool // include note names (default: true)
	IncludeFunctions bool // include function names (default: true)
	IncludeScales    bool // include scale names (default: true)
	Deduplicate      bool // remove duplicates (default: true)
}
