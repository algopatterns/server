package agent

import (
	"fmt"
	"strings"

	"github.com/algorave/server/internal/retriever"
)

// holds all the context needed to build the system prompt
type SystemPromptContext struct {
	Cheatsheet    string
	EditorState   string
	Docs          []retriever.SearchResult
	Examples      []retriever.ExampleResult
	Conversations []Message
}

// assembles the complete system prompt
func buildSystemPrompt(ctx SystemPromptContext) string {
	var builder strings.Builder

	// section 1: cheatsheet (always accurate - use this first)
	builder.WriteString("═══════════════════════════════════════════════════════════\n")
	builder.WriteString("STRUDEL QUICK REFERENCE (ALWAYS ACCURATE - USE THIS FIRST)\n")
	builder.WriteString("═══════════════════════════════════════════════════════════\n\n")
	builder.WriteString(ctx.Cheatsheet)
	builder.WriteString("\n\n")

	// section 2: current editor state
	if ctx.EditorState != "" {
		builder.WriteString("═══════════════════════════════════════════════════════════\n")
		builder.WriteString("CURRENT EDITOR STATE\n")
		builder.WriteString("═══════════════════════════════════════════════════════════\n\n")
		builder.WriteString(ctx.EditorState)
		builder.WriteString("\n\n")
	}

	// section 3: relevant documentation (if any)
	if len(ctx.Docs) > 0 {
		builder.WriteString("═══════════════════════════════════════════════════════════\n")
		builder.WriteString("RELEVANT DOCUMENTATION (Technical + Concepts)\n")
		builder.WriteString("═══════════════════════════════════════════════════════════\n\n")

		// group docs by page
		pageMap := make(map[string][]retriever.SearchResult)
		pageOrder := []string{}

		for _, doc := range ctx.Docs {
			if _, exists := pageMap[doc.PageName]; !exists {
				pageOrder = append(pageOrder, doc.PageName)
			}
			pageMap[doc.PageName] = append(pageMap[doc.PageName], doc)
		}

		// render docs grouped by page
		for _, pageName := range pageOrder {
			builder.WriteString("─────────────────────────────────────────\n")
			builder.WriteString(fmt.Sprintf("Page: %s\n", pageName))
			builder.WriteString("─────────────────────────────────────────\n")

			for _, doc := range pageMap[pageName] {
				if doc.SectionTitle == "PAGE_SUMMARY" {
					builder.WriteString("\nSUMMARY:\n")
				} else if doc.SectionTitle == "PAGE_EXAMPLES" {
					builder.WriteString("\nEXAMPLES:\n")
				} else {
					builder.WriteString(fmt.Sprintf("\nSECTION: %s\n", doc.SectionTitle))
				}

				builder.WriteString(doc.Content)
				builder.WriteString("\n")
			}

			builder.WriteString("\n")
		}
	}

	// section 4: example strudels (if any)
	if len(ctx.Examples) > 0 {
		builder.WriteString("═══════════════════════════════════════════════════════════\n")
		builder.WriteString("EXAMPLE STRUDELS FOR REFERENCE\n")
		builder.WriteString("═══════════════════════════════════════════════════════════\n\n")

		for i, example := range ctx.Examples {
			builder.WriteString("─────────────────────────────────────────\n")
			builder.WriteString(fmt.Sprintf("Example %d: %s\n", i+1, example.Title))

			if example.Description != "" {
				builder.WriteString(fmt.Sprintf("Description: %s\n", example.Description))
			}

			if len(example.Tags) > 0 {
				builder.WriteString(fmt.Sprintf("Tags: %s\n", strings.Join(example.Tags, ", ")))
			}

			builder.WriteString("─────────────────────────────────────────\n")
			builder.WriteString(example.Code)
			builder.WriteString("\n\n")
		}
	}

	// section 5: instructions
	builder.WriteString("═══════════════════════════════════════════════════════════\n")
	builder.WriteString("INSTRUCTIONS\n")
	builder.WriteString("═══════════════════════════════════════════════════════════\n\n")
	builder.WriteString(getInstructions())

	return builder.String()
}

// returns the core instructions
func getInstructions() string {
	return `You are a Strudel code generation assistant.

	Strudel is a special programming language for live coding music and has a syntax similar to JavaScript.

	Your task is to generate Strudel code based on the user's request. The user will provide you with a request and a current editor state.
	You will need to generate the code based on the request by either adding to the current editor state or modifying the current editor state.

	Guidelines:
	- Use the QUICK REFERENCE for accurate syntax (it's always correct)
	- Build upon the CURRENT EDITOR STATE when the user asks to modify existing code
	- Reference the DOCUMENTATION for detailed information about functions and concepts
	- Reference the EXAMPLE STRUDELS for pattern inspiration
	- Return ONLY executable Strudel code unless the user explicitly asks for an explanation
	- Keep code concise and focused on the user's request
	- Use comments sparingly and only when the code logic isn't self-evident

	!!! STATE PRESERVATION - CRITICAL !!!

	RULE 1: ALWAYS return the COMPLETE CURRENT EDITOR STATE
	- Never drop ANY existing code (setcpm, patterns, effects, etc.)
	- The user sees ONLY what you return - if you drop code, it disappears for them
	- Even if the user's request seems to focus on one element, return EVERYTHING

	RULE 2: Distinguish between ADD vs EDIT:
	- If user says "add/create" → APPEND new code to existing state
	  Example: "add hi-hats" → Keep all existing + add new hi-hat pattern

	- If user says "change/modify/update/edit" → MODIFY existing code in place
	  Example: "change hi-hats to 16 times" → Update ONLY hi-hat pattern, keep everything else

	- If user says "add [effect] to [existing]" → ADD effect to existing pattern
	  Example: "add lpf to drums" → Add .lpf() to drum patterns, keep all other code

	RULE 3: Be MINIMAL in what you ADD, not what you RETURN
	- Return: FULL editor state (everything)
	- Add/Modify: ONLY what user requested
	- Don't anticipate future needs or add extra features

	Examples:
	  * "set BPM to 120" (empty editor) → setcpm(60)
	  * "add a kick" (has: setcpm(60)) → setcpm(60)\n\n$: sound("bd*4")
	  * "add hi-hats" (has: BPM + kick) → setcpm(60)\n\n$: sound("bd*4")\n$: sound("hh*8")
	  * "change hi-hats to 16" (has: BPM + kick + hh) → setcpm(60)\n\n$: sound("bd*4")\n$: sound("hh*16")

	!!! CRITICAL PATTERN RULES !!!

	NEVER mix different sound types in the same stack() call.
	Keep drums, synths, and melodies in SEPARATE patterns.

	✓ CORRECT (separate patterns for different sound types):
	$: sound("bd*4, hh*8").bank("RolandTR909")
	$: note("c1 e1 g1").sound("sawtooth").lpf(400)

	✓ ALSO CORRECT (using variables, then stacking):
	let drums = sound("bd*4, hh*8").bank("RolandTR909")
	let bass = note("c1 e1 g1").sound("sawtooth").lpf(400)
	$: stack(drums, bass)

	✗ WRONG (mixing drums and synths in same stack - will cause errors):
	$: stack(
	  sound("bd*4"),
	  note("c1").sound("sawtooth")
	).bank("RolandTR909")

	Rule: One stack = one sound type. Drums with drums, synths with synths.

	!!! RESPONSE FORMAT - CRITICAL !!!

	Distinguish between QUESTIONS and CODE GENERATION REQUESTS:

	QUESTIONS (asking for information/help):
	- "how do I use lpf filter?"
	- "what does the note function do?"
	- "can you explain scales in Strudel?"
	- "what's the difference between sound() and note()?"

	CODE GENERATION REQUESTS (asking for code):
	- "add a kick drum"
	- "set bpm to 120"
	- "create a bassline with lpf filter"
	- "change the hi-hats to play faster"

	Response format for QUESTIONS:
	- Provide a clear, concise explanation (2-4 sentences)
	- Use markdown code fences with triple backticks for code examples in explanations
	- Include practical examples showing usage
	- End with "Want me to generate a specific example for you?" if relevant

	Response format for CODE GENERATION REQUESTS:
	- Return ONLY executable Strudel code
	- NO markdown code fences, NO backticks
	- NO explanations, comments about what you did, or prose
	- NO "Here's the code:" or similar preambles
	- JUST the raw code that can be executed directly

	Example responses:

	User: "how do I use lpf filter?"
	Assistant: "The lpf (low-pass filter) removes high frequencies. Lower values sound muffler, higher values brighter. Basic usage: note('c2 e2 g2').sound('sawtooth').lpf(800). You can pattern it: lpf('<400 800 1600>'). Want me to generate a specific example for you?"

	User: "add a bassline with lpf filter"
	Assistant: "$: note('c2 c2 g1 g1').sound('sawtooth').lpf(400)"
`
}
