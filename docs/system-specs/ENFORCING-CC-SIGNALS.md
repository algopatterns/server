# Enforcing CC Signals

## Overview

Behavioral detection system that protects CC Signal restrictions (especially `no-ai`) from being bypassed through copy-paste. When code is pasted, the AI agent is temporarily blocked until the user makes significant edits (30%+), demonstrating genuine engagement with the code.

## Problem Statement

Users can bypass CC Signal restrictions by:

1. Viewing public code with `no-ai` signal
2. Copy-pasting into their own editor
3. Requesting AI assistance on the copied code

This system detects paste behavior and temporarily blocks AI access until the user demonstrates genuine engagement with the code through significant edits.

## How It Works

```
User pastes code → WS detects large code_update → Validates against DB
                                                          ↓
                              ┌─────────────────────────────────────────┐
                              │ Is code from legitimate source?         │
                              │ - User's own strudel? → No lock         │
                              │ - Public strudel (allows AI)? → No lock │
                              │ - Public strudel (no-ai)? → LOCK        │
                              │ - External paste? → LOCK                │
                              └─────────────────────────────────────────┘
                                                          ↓
User requests AI → REST checks paste lock → If locked: reject with message
                                                          ↓
User makes 30%+ edits → WS detects significant change → Removes lock
                                                          ↓
User requests AI → REST checks paste lock → Unlocked → Proceed
```

## Paste Lock Decision Table

| Action                                | Large Delta? | DB Match?      | CC Signal          | Result  |
| ------------------------------------- | ------------ | -------------- | ------------------ | ------- |
| User loads their saved strudel        | Yes          | User owns it   | Any                | No lock |
| User forks public strudel (allows AI) | Yes          | Public strudel | cc-cr, cc-op, etc. | No lock |
| User forks public strudel (no-ai)     | Yes          | Public strudel | no-ai              | Locked  |
| User pastes external code             | Yes          | No match       | N/A                | Locked  |
| User types code gradually             | No           | N/A            | N/A                | No lock |

## Design Decisions

| Decision         | Value                         | Rationale                                      |
| ---------------- | ----------------------------- | ---------------------------------------------- |
| Paste threshold  | 200+ chars OR 50+ lines delta | Normal typing is 1-5 chars; paste is hundreds+ |
| Unlock threshold | 30% edit distance             | Requires genuine engagement with code          |
| Lock TTL         | 1 hour                        | Auto-cleanup for disconnected sessions         |
| No session_id    | Fail open                     | Backwards compatibility                        |
| Redis error      | Fail open                     | Better UX than blocking on infra issues        |

## Server-Side Validation

**IMPORTANT:** The server validates independently of the frontend `source` field. The `source` field is a hint only - the server performs its own validation:

1. Does code match user's own strudels? → No lock (legitimate load)
2. Does code match any public strudel **that allows AI**? → No lock (legitimate fork)
3. Does code match public strudel **with `no-ai` CC signal**? → **LOCK** (protect creator's wishes)
4. Otherwise with large delta → Paste lock

This prevents malicious clients from spoofing the `source` field to bypass detection.

## Architecture

### Design Principles

- **Decoupled from WebSocket**: Paste lock validation only requires Redis, not an active WebSocket connection. This makes the ccsignals module portable to other projects.
- **Fail Open on Infrastructure Errors**: Redis failures don't block users, but are logged for visibility.
- **Fail Closed on Missing Data**: If a fork's parent strudel can't be found, AI is blocked (can't verify it wasn't no-ai).

### Redis Keys

```
paste_lock:{sessionID}     → "1" (TTL: 1 hour)
paste_baseline:{sessionID} → <code at time of paste> (TTL: 1 hour)
```

### Data Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              FRONTEND                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│  Editor paste event → Mark next update source as 'paste'                    │
│  Code change → sendCodeUpdate({ code, source })                             │
│  AI request → POST /agent/generate { session_id, ... }                      │
└──────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              BACKEND                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│  WebSocket Hub:                                                              │
│    1. Receive code_update                                                    │
│    2. Detect large delta (behavioral detection)                              │
│    3. Validate against DB:                                                   │
│       - Check user's strudels                                                │
│       - Check public strudels (respecting CC signals)                        │
│    4. Set/remove paste lock in Redis                                         │
│                                                                              │
│  REST API:                                                                   │
│    1. Receive /agent/generate with session_id                               │
│    2. Validate session is active                                             │
│    3. Check paste lock in Redis                                              │
│    4. If locked: return 403                                                  │
│    5. If unlocked: proceed with AI generation                               │
└──────────────────────────────────────────────────────────────────────────────┘
```

## Files Modified

### Backend

| File                             | Changes                                                                        |
| -------------------------------- | ------------------------------------------------------------------------------ |
| `internal/buffer/buffer.go`      | `SetPasteLock`, `IsPasteLocked`, `GetPasteBaseline`, `RemovePasteLock` methods |
| `internal/buffer/types.go`       | `keyPasteLock`, `keyPasteBaseline` constants                                   |
| `internal/buffer/utils.go`       | `LevenshteinDistance`, `IsLargeDelta`, `IsSignificantEdit`                     |
| `algorave/strudels/queries.go`   | `queryUserOwnsStrudelWithCode`, `queryPublicStrudelExistsWithCodeAllowsAI`     |
| `algorave/strudels/strudels.go`  | `UserOwnsStrudelWithCode`, `PublicStrudelExistsWithCodeAllowsAI` methods       |
| `internal/websocket/types.go`    | `Source` field in `CodeUpdatePayload`                                          |
| `internal/websocket/handlers.go` | `CodeUpdateHandler` with paste detection and DB validation                     |
| `internal/websocket/hub.go`      | `IsSessionActive` method                                                       |
| `api/rest/agent/types.go`        | `SessionID` in `GenerateRequest`                                               |
| `api/rest/agent/handlers.go`     | Paste lock validation before AI generation                                     |

### Frontend

| File                                   | Changes                                      |
| -------------------------------------- | -------------------------------------------- |
| `lib/websocket/types.ts`               | `source` field in `CodeUpdatePayload`        |
| `lib/websocket/client.ts`              | `sendCodeUpdate` accepts source parameter    |
| `lib/stores/editor.ts`                 | `nextUpdateSource` state for paste tracking  |
| `lib/api/agent/types.ts`               | `session_id` in `GenerateRequest`            |
| `lib/hooks/use-agent.ts`               | Pass `session_id`, handle paste_locked error |
| `components/shared/strudel-editor.tsx` | Paste event listener                         |

## API Changes

### WebSocket: code_update

```json
{
  "type": "code_update",
  "payload": {
    "code": "...",
    "cursor_line": 10,
    "cursor_col": 5,
    "source": "typed"
  }
}
```

### REST: POST /api/v1/agent/generate

Request now includes `session_id`:

```json
{
  "user_query": "...",
  "editor_state": "...",
  "session_id": "abc123"
}
```

New error response when paste locked:

```json
{
  "error": "paste_locked",
  "message": "AI assistant temporarily disabled - please make significant edits to the pasted code before using AI. This helps protect code shared with 'no-ai' restrictions."
}
```

## Edge Cases

| Case                               | Handling                                                              |
| ---------------------------------- | --------------------------------------------------------------------- |
| Session disconnect during lock     | TTL auto-expires after 1 hour                                         |
| Multiple pastes                    | Each paste resets baseline to new code                                |
| Loading own strudel                | Server validates: code matches user's strudel in DB → no lock         |
| Forking public strudel (allows AI) | Server validates: code matches public strudel without no-ai → no lock |
| Forking public strudel (no-ai)     | Server validates: code matches but has no-ai → **LOCK**               |
| No session_id in request           | Fail open (backwards compat)                                          |
| Redis unavailable                  | Fail open, log error                                                  |
| Anonymous user loads large code    | Still checks public strudels; if no match or no-ai, locks             |
| User spoofs `source` field         | Server ignores it, validates against DB                               |
| Parent strudel deleted             | AI blocked - can't verify CC signal was not no-ai                     |
| Fake/invalid forked_from_id        | AI blocked - must reference valid strudel to use AI on forks          |

## Security Considerations

### What This Catches

- Casual copy-paste from public strudels with `no-ai`
- Accidental use of restricted code
- External code pasted without authorization

### What This Doesn't Catch

- Simulated gradual typing (requires custom tooling)
- Code memorized and retyped
- Heavily modified copied code (30%+ changes = transformative use)

This is a **first layer of defense** that raises the bar for circumvention while maintaining good UX for honest users.

## Testing Checklist

- [ ] Paste 200+ chars from external source → lock is set
- [ ] Paste 50+ lines from external source → lock is set
- [ ] Normal typing → no lock
- [ ] Load own saved strudel → no lock (server validates via DB)
- [ ] Fork public strudel (allows AI) → no lock (server validates via DB)
- [ ] Fork public strudel with `no-ai` CC signal → LOCK (server respects creator's wishes)
- [ ] AI request while locked → 403 error
- [ ] Edit 30%+ of pasted code → lock removed
- [ ] AI request after unlock → succeeds
- [ ] Session disconnect → lock auto-expires
- [ ] Redis down → AI still works (fail open)
- [ ] Spoof `source='loaded_strudel'` with external code → lock is set (server ignores frontend hint)
- [ ] Anonymous user loads public strudel code (allows AI) → no lock
- [ ] Anonymous user pastes random code → lock is set
