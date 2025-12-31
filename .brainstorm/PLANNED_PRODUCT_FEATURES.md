# Algorave - Future Features & Improvements

Documentation for planned new features for front facing app.

## Phase 1 Collaboration

- [x] WebSocket infrastructure
- [x] Session management (create, join, invite tokens)
- [x] Real-time code synchronization
- [x] Multi-user agent requests
- [x] Testing & validation
- [x] Documentation

---

## Phase 2: Advanced Collaboration Features

### Real-time Cursor Tracking
**Status:** Deferred from Phase 1
**Priority:** Medium
**Effort:** 2-3 days

**Description:**
Enable real-time cursor position sharing between co-authors to show where each participant is editing.

**Implementation Notes:**
- Handler exists: `internal/websocket/handlers/cursor_position.go`
- Message type: `TypeCursorPosition`
- Needs throttling (100-200ms debounce) to prevent spam
- Client-side UI needed to render colored cursors with names
- Consider: Only show cursors for co-authors, not viewers

**Tasks:**
- [ ] Add throttling/debouncing to cursor updates
- [ ] Implement cursor UI in web frontend
- [ ] Add cursor color assignment per user
- [ ] Test with 5+ concurrent users
- [ ] Add user preference to hide cursors

**Files:**
- `internal/websocket/handlers/cursor_position.go` (already exists)
- `internal/websocket/message.go` (payload defined)

---

### Streaming Agent Responses
**Status:** Planned
**Priority:** Medium
**Effort:** 3-5 days

**Description:**
Stream LLM token generation in real-time instead of waiting for complete response.

**Benefits:**
- Faster perceived response time
- See code being generated live
- Can stop generation mid-stream

**Implementation Notes:**
- Agent needs to support streaming (check if Claude API supports)
- WebSocket message type: `TypeAgentResponseChunk`
- Buffer partial responses on client
- Handle connection drops mid-stream

**Tasks:**
- [ ] Investigate Claude API streaming support
- [ ] Create chunk message types
- [ ] Implement server-side streaming
- [ ] Update agent handlers to stream
- [ ] Client-side buffering and rendering
- [ ] Add "Stop generation" button

---

### Session History & Replay
**Status:** Planned
**Priority:** Low
**Effort:** 1 week

**Description:**
Record all session events (code changes, messages) for playback and debugging.

**Features:**
- Timeline of all changes
- Replay session from any point
- Download session transcript
- Session recording permission (host-controlled)

**Database:**
- Table: `session_events`
  - Columns: `id, session_id, event_type, user_id, payload, timestamp`

**Tasks:**
- [ ] Create session_events table
- [ ] Record all WebSocket events
- [ ] Build replay API endpoint
- [ ] Create timeline UI component
- [ ] Add privacy controls

---

### Presence & Typing Indicators
**Status:** Planned
**Priority:** Low
**Effort:** 2-3 days

**Description:**
Show who's currently active and who's typing.

**Features:**
- "Alice is typing..." indicator
- Active/idle status (based on last activity)
- "5 people viewing" counter

**Implementation:**
- Message type: `TypeUserTyping`
- Timeout-based idle detection (5 minutes)
- Heartbeat every 30 seconds

**Tasks:**
- [ ] Add typing event handlers
- [ ] Implement idle timeout
- [ ] Create presence UI components
- [ ] Add user avatars/profile pictures

---

## Phase 3: Scheduled Events

### Event System
**Status:** Planned
**Priority:** High
**Effort:** 2-3 weeks

**Description:**
Live coding events with waiting rooms, RSVP, and scheduled sessions.

**Features:**
- Create scheduled events with start/end times
- RSVP system with capacity limits
- Waiting room before event starts
- Automatic session creation at event time
- Post-event recording/gallery

**Database Tables:**
- `events` (title, description, start_time, end_time, capacity, host_id)
- `event_rsvps` (event_id, user_id, status)

**Tasks:**
- [ ] Design event schema
- [ ] Event CRUD API endpoints
- [ ] RSVP system
- [ ] Waiting room UI
- [ ] Event calendar view
- [ ] Email notifications
- [ ] Event recordings

---

### Performance Optimizations

#### Query Caching
**Status:** Planned
**Priority:** Medium
**Effort:** 1 week

- Cache common query transformations
- Semantic caching for repeated queries
- Redis integration for distributed cache

#### Connection Pooling
**Status:** Review needed
**Priority:** Medium

- Review pgx connection pool settings
- Tune pool size based on load testing
- Monitor connection leaks

#### WebSocket Scaling
**Status:** Planned
**Priority:** High (before production)
**Effort:** 1-2 weeks

**Current:** Single-server, in-memory hub
**Target:** Multi-server with Redis pub/sub

**Tasks:**
- [ ] Add Redis adapter for hub
- [ ] Implement pub/sub for cross-server messages
- [ ] Sticky sessions in load balancer
- [ ] Test failover scenarios
- [ ] Benchmark: 100 concurrent sessions

---

## Testing & Monitoring

### Integration Tests
**Status:** Minimal
**Priority:** High
**Effort:** 1 week

**Needed:**
- [ ] WebSocket connection tests
- [ ] Multi-client broadcast tests
- [ ] Session CRUD tests
- [ ] Invite token tests
- [ ] Permission enforcement tests

### Load Testing
**Status:** Not started
**Priority:** Medium

**Targets:**
- 100 concurrent sessions
- 500 concurrent WebSocket connections
- 1000 messages/second throughput

**Tools:**
- k6 for load generation
- Prometheus + Grafana for metrics

### Monitoring
**Status:** Not started
**Priority:** High (before production)

**Metrics to track:**
- Active sessions count
- WebSocket connection count
- Message throughput
- Database query latency
- Error rates

---

## Documentation

### API Documentation
**Status:** None
**Priority:** High
**Effort:** 2-3 days

**Needed:**
- [x] REST API reference (OpenAPI/Swagger)
- [x] WebSocket message protocol spec
- [ ] Authentication flow diagram
- [ ] Session lifecycle diagram

### User Documentation
**Status:** None
**Priority:** Medium

- [ ] Getting started guide
- [ ] Collaboration tutorial
- [ ] Event hosting guide
- [ ] Strudel syntax reference

---

## Security

### Rate Limiting Improvements
**Status:** Basic implementation
**Priority:** High

**Current:** requests/minute rate limits per IP (HTTP only)
**Needed:**
- [x] WebSocket message rate limiting
- [  ] Per-session message limits
- [ ] IP ban list for abuse

### Security Audit
**Status:** Not started
**Priority:** High (before production)

**Focus Areas:**
- [ ] SQL injection prevention review
- [ ] XSS vulnerability check
- [ ] WebSocket message validation
- [ ] Invite token security
- [ ] Session hijacking prevention
- [ ] CORS configuration review

---

### WebSocket Origin Validation
**Status:** TODO in production
**Location:** `api/websocket/handler.go:21`
**Current:** `CheckOrigin` returns true for all origins
**Needed:** Whitelist allowed origins in production
