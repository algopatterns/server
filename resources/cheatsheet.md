# STRUDEL QUICK REFERENCE

## EDITOR SYNTAX (CRITICAL!)

Every pattern MUST start with `$:` or `$<name>:` in the Strudel editor.

```javascript
$: sound("bd sd")                 // Single pattern
$: note("c e g").sound("piano")   // Pattern with method chain
$melody: note("c e g").sound("piano")  // Named pattern
$drums: sound("bd sd hh cp")      // Named pattern
$bass: note("c2 e2").sound("sawtooth")

// Multiple patterns run simultaneously
$: sound("bd*4")
$: sound("hh*8")
$: note("c e g").sound("piano")

// Mute patterns with underscore
_$: sound("bd")                   // Muted
```

---

## BASICS

### Sounds & Notes
```javascript
sound("casio")                    // Play a sound
sound("bd hh sd hh")              // Sequence
sound("bd*4, hh*8")               // Parallel (comma)

note("c e g b")                   // Notes by letter
note("c2 e3 g4")                  // With octave (default: 3)
note("48 52 55")                  // By MIDI number
```

### Tempo
```javascript
setcpm(30)    // ~120 BPM
setcpm(60)    // ~240 BPM
hush()        // Stop all
```

### BPM - CPM - CPS Conversion
```
BPM ≈ CPM × 4
CPM = CPS × 60
CPS = CPM ÷ 60

CPM 15  = ~60 BPM  (slow)
CPM 30  = ~120 BPM (moderate)
CPM 45  = ~180 BPM (fast)
CPM 60  = ~240 BPM (very fast)
```
**Prefer setcpm() for tempo** - it's more intuitive than CPS.

---

## MINI-NOTATION

### Sequences & Rests
```javascript
sound("bd sd hh cp")              // Space-separated
sound("bd ~ sd ~")                // ~ = rest/silence
```

### Timing Symbols
```javascript
sound("<bd hh sd>")               // <> = play one, then next, then next...
sound("[bd hh sd]*2")             // *  = play faster (2x speed)
sound("[bd hh sd]/2")             // /  = play slower (half speed)
sound("bd@2 hh")                  // @  = hold longer (bd lasts 2x)
sound("bd!3")                     // !  = repeat (3 times, same speed)
sound("bd(3,8)")                  // () = spread evenly (3 hits in 8 slots)
```

### Subdivision & Parallel
```javascript
sound("bd [hh hh] sd cp")         // [brackets] subdivide
sound("bd, hh*8, sd")             // Comma = simultaneous
```

### Randomness
```javascript
note("[c|e|g]")                   // Random choice
note("[c e g]?")                  // 50% chance removal
note("[c e g]?0.2")               // 20% chance removal
```

### Sample Selection
```javascript
sound("hh:0 hh:1 hh:2 hh:3")      // Select sample by number
n("0 1 2 3").sound("jazz")        // Using n() for sample index
```

---

## SOUND BANKS

### Default Drums
```
bd=bass drum, sd=snare, hh=hi-hat, oh=open hat
cp=clap, rim=rimshot, cr=crash, rd=ride
ht/mt/lt=toms, cb=cowbell, sh=shaker
```

### Synths & Noise
```
sine, triangle, square, sawtooth
white, pink, brown (noise)
```

### Change Bank
```javascript
sound("bd hh sd").bank("RolandTR909")
// Banks: RolandTR808, RolandTR909, RolandTR707, AkaiLinn
```

---

## EFFECTS

### Filters
```javascript
// Low-pass (removes highs, sounds darker/muffled)
note("c2").lpf(800)               // Cutoff frequency
note("c2").lpf(800).lpq(5)        // With resonance (q)
note("c2").lpf("<200 800 2000>")  // Pattern cutoff

// High-pass (removes lows, sounds thinner)
note("c2").hpf(500)               // Cutoff frequency
note("c2").hpf(500).hpq(5)        // With resonance

// Band-pass (keeps only frequencies around cutoff)
sound("bd").bpf(1000)             // Center frequency
sound("bd").bpf(1000).bpq(5)      // With resonance
```

### Volume & Envelope
```javascript
sound("bd").gain(0.5)             // Volume (0-1)
sound("hh*8").gain("[.25 1]*4")   // Accent pattern
note("c3").attack(0.1).decay(0.2).sustain(0.5).release(0.3)
note("c3").adsr(".1:.2:.5:.3")    // Shorthand
```

### Reverb & Delay
```javascript
sound("bd").room(0.5)             // Reverb amount
sound("bd").room(0.8).roomsize(4) // With size
sound("bd").delay(0.5)            // Delay amount
sound("bd").delay(".5:.25:.7")    // delay:time:feedback
```

### Distortion & Shape
```javascript
sound("bd").distort(3)            // Distortion
sound("bd").crush(4)              // Bit crush
sound("bd").shape(0.5)            // Waveshape
sound("bd").coarse(8)             // Sample rate reduction
```

### Panning & Stereo
```javascript
sound("bd").pan(0)                // Left
sound("bd").pan("<0 0.5 1>")      // Pattern L-C-R
sound("bd").jux(rev)              // Stereo widening
```

### Sidechain / Ducking
```javascript
// Make other sounds "duck" when kick hits
$: sound("bd*4").duckorbit(1)                    // Kick triggers duck
$: note("c2 c2 c2 c2").sound("sawtooth")
    .orbit(1).duckdepth(0.8).duckattack(0.01).duckrelease(0.2)
```

### Tremolo & Vibrato
```javascript
note("c3").vib(4)                 // Vibrato speed Hz
note("c3").vib(4).vibmod(0.5)     // With depth (semitones)
note("c3").tremolo(4)             // Tremolo speed
note("c3").tremolo(4).tremdp(0.5) // With depth
```

### Vowel & Phaser
```javascript
note("c3").vowel("<a e i o u>")   // Vowel sounds
note("c3").phaser(2)              // Phaser speed Hz
note("c3").phaser(2).phaserdepth(0.5)
```

---

## TIME MODIFIERS

```javascript
sound("bd hh").fast(2)            // Speed up
sound("bd hh").slow(2)            // Slow down
sound("bd").early(0.1)            // Shift earlier
sound("bd").late(0.1)             // Shift later
sound("bd hh sd").rev()           // Reverse
```

---

## SAMPLE CONTROL

```javascript
sound("bd").speed(2)              // Double speed (pitch up)
sound("bd").speed(0.5)            // Half speed (pitch down)
sound("bd").speed(-1)             // Reverse playback
sound("bd").begin(0.25)           // Start at 25% of sample
sound("bd").end(0.75)             // End at 75% of sample
sound("hh oh").cut(1)             // Cut group (hh chokes oh)
sound("breaks").loopAt(2)         // Fit sample into 2 cycles
```

---

## PROBABILITY FUNCTIONS

```javascript
sound("hh*8").sometimes(x => x.gain(0.5))    // 50% chance
sound("hh*8").often(x => x.speed(2))         // 75% chance
sound("hh*8").rarely(x => x.crush(4))        // 25% chance
sound("hh*8").degrade()                      // Remove 50% randomly
sound("hh*8").degradeBy(0.3)                 // Remove 30% randomly
```

---

## SCALES & CHORDS

```javascript
n("0 2 4 6").scale("C:major")     // Scale degrees
n("0 2 4").scale("C:minor")
n("0 2 4").scale("C:pentatonic")

// Common scales: major, minor, dorian, pentatonic, blues
```

### Transpose
```javascript
note("c e g").transpose(12)       // Up an octave
n("0 2 4").scale("C:major").scaleTranspose(1)
```

### Common Chord Progressions
Scale degrees: I=0, ii=1, iii=2, IV=3, V=4, vi=5, vii=6

```javascript
// Pop (I-V-vi-IV)
n("0 4 5 3").scale("C:major").sound("piano")

// Jazz ii-V-I
n("1 4 0").scale("C:major").sound("piano")

// Rock (I-bVII-IV)
n("0 -1 3").scale("C:major").sound("sawtooth")

// EDM (i-VI-III-VII in minor)
n("0 5 2 6").scale("C:minor").sound("sawtooth")
```

---

## PATTERN RULES

### Separate Sound Types
```javascript
// CORRECT - separate patterns for drums vs synths
$: sound("bd*4, hh*8").bank("RolandTR909")
$: note("c2 e2").sound("sawtooth").lpf(400)

// CORRECT - stack same sound types
$: stack(
  sound("bd*4"),
  sound("hh*8"),
  sound("cp*2")
).bank("RolandTR909")

// WRONG - mixing drums and synths in same stack
$: stack(sound("bd*4"), note("c1").sound("sawtooth")).bank("RolandTR909")
```

---

## TIPS & BEST PRACTICES

### Working with Scales
```javascript
// Always use valid scale names:
"C:major", "A:minor", "D:dorian", "G:mixolydian"
"F:pentatonic", "Bb:blues"
"C4:major"                        // With octave
"C:minor:pentatonic"              // Compound scales
```

### Scale Types
```
major, minor, dorian, phrygian, lydian, mixolydian
pentatonic, blues, harmonic_minor
```
Available for all keys: C, C#, D, Eb, E, F, F#, G, Ab, A, Bb, B

### Combining Parameters
```javascript
note("c e g").sound("piano")
  .lpf(800)
  .room(0.5)
  .delay(0.25)
  .gain(0.8)
```

### Modulation Patterns
```javascript
note("c3").lpf(sine.range(200, 2000).slow(4))   // LFO filter
sound("hh*16").gain(sine.range(0.3, 1).slow(2)) // Varying gain
```

---

## COMMON PATTERNS

```javascript
// Basic drum loop
$: sound("bd*4, [~ sd]*2, hh*8").bank("RolandTR909")

// Bassline
$: note("c2 c2 g1 g1").sound("sawtooth").lpf(400)

// Melody with scale
$: n("0 2 4 [6 8]").scale("C:minor").sound("piano")

// Filter sweep
$: note("c2").sound("sawtooth").lpf(sine.range(200, 2000).slow(4))

// Full loop example
$: sound("bd*4").bank("RolandTR909")
$: sound("[~ hh]*8")
$: note("<0 2 4 [6 8]>").scale("C:minor").sound("sawtooth").lpf(800)
```
