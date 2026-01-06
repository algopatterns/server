-- track when user strudels are used as RAG examples
CREATE TABLE IF NOT EXISTS rag_attributions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_strudel_id UUID NOT NULL REFERENCES user_strudels(id) ON DELETE CASCADE,
    target_strudel_id UUID REFERENCES user_strudels(id) ON DELETE SET NULL,
    requesting_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    similarity_score FLOAT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_rag_attributions_source ON rag_attributions(source_strudel_id);
CREATE INDEX IF NOT EXISTS idx_rag_attributions_created ON rag_attributions(created_at);

-- user notifications
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    title TEXT NOT NULL,
    body TEXT,
    data JSONB,
    read BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_user_unread ON notifications(user_id, read) WHERE read = false;

COMMENT ON TABLE rag_attributions IS 'Tracks when user strudels are retrieved as RAG examples';
COMMENT ON TABLE notifications IS 'User notifications for attribution and other events';
