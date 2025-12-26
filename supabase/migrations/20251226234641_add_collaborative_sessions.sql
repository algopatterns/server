-- Migration: Add Collaborative Sessions Infrastructure
-- Description: Tables for real-time collaborative code sessions with WebSocket support
-- Created: 2025-12-26

-- Collaborative sessions table (persistent multi-user sessions)
CREATE TABLE sessions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  host_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  code TEXT NOT NULL DEFAULT '',
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  ended_at TIMESTAMPTZ,
  last_activity TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_sessions_host ON sessions(host_user_id);
CREATE INDEX idx_sessions_active ON sessions(is_active) WHERE is_active = true;

-- Session participants with roles and status
CREATE TABLE session_participants (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
  user_id UUID REFERENCES users(id) ON DELETE SET NULL,
  display_name TEXT,  -- For anonymous users who don't have accounts
  role TEXT NOT NULL CHECK (role IN ('host', 'co-author', 'viewer')),
  status TEXT NOT NULL CHECK (status IN ('invited', 'active', 'left')) DEFAULT 'active',
  joined_at TIMESTAMPTZ DEFAULT NOW(),
  left_at TIMESTAMPTZ,
  UNIQUE(session_id, user_id)
);

CREATE INDEX idx_participants_session ON session_participants(session_id);
CREATE INDEX idx_participants_user ON session_participants(user_id);

-- Invite tokens for sharing sessions
CREATE TABLE invite_tokens (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
  token TEXT UNIQUE NOT NULL,
  role TEXT NOT NULL CHECK (role IN ('co-author', 'viewer')),
  max_uses INTEGER,  -- NULL = unlimited uses
  uses_count INTEGER DEFAULT 0,
  expires_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_invite_tokens_session ON invite_tokens(session_id);
CREATE INDEX idx_invite_tokens_token ON invite_tokens(token);

-- Session conversation history (for collaborative AI context)
CREATE TABLE session_messages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
  user_id UUID REFERENCES users(id) ON DELETE SET NULL,
  role TEXT NOT NULL CHECK (role IN ('user', 'assistant')),
  content TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_messages_session ON session_messages(session_id, created_at);

-- Comments for documentation
COMMENT ON TABLE sessions IS 'Collaborative coding sessions with real-time WebSocket support';
COMMENT ON TABLE session_participants IS 'Users participating in sessions with role-based permissions';
COMMENT ON TABLE invite_tokens IS 'Shareable tokens for joining sessions';
COMMENT ON TABLE session_messages IS 'Conversation history for collaborative AI context';

COMMENT ON COLUMN session_participants.role IS 'host: full control, co-author: can edit code, viewer: read-only';
COMMENT ON COLUMN session_participants.status IS 'invited: pending join, active: currently participating, left: no longer participating';
COMMENT ON COLUMN invite_tokens.max_uses IS 'NULL means unlimited uses, otherwise restricts number of times token can be used';
