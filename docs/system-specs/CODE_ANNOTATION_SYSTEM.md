# Code Annotation System

## Overview

The agent uses enhanced system prompt instructions to achieve surgical precision when modifying user code, without requiring additional LLM calls.

## Implementation

**Location**: `/internal/agent/prompt.go`

The system prompt includes:

1. **Request Type Classification**
   - ADDITIVE: Adding new elements
   - MODIFICATION: Changing existing elements
   - DELETION: Removing elements
   - QUESTIONS: Asking for help/information

2. **Surgical Precision Process**
   - Step 1: Identify the target element
   - Step 2: Locate the target in current editor state
   - Step 3: Make the change surgically
   - Step 4: Preserve everything else

3. **Few-shot Examples**
   - Concrete examples for each request type
   - Demonstrates correct before/after patterns

## Benefits

- No additional API calls (no latency or cost overhead)
- Request classification happens within the main generation call
- Surgical accuracy through detailed instructions and examples

## Related

- [RAG_ARCHITECTURE.md](./RAG_ARCHITECTURE.md)
- [HYBRID_RETRIEVAL_GUIDE.md](./HYBRID_RETRIEVAL_GUIDE.md)
