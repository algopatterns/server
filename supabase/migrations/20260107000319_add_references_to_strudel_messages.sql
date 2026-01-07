-- Add reference columns to strudel_messages for CC attribution display
-- Stores which strudels and docs were used as RAG context for each AI response

ALTER TABLE strudel_messages
ADD COLUMN strudel_references JSONB DEFAULT NULL,
ADD COLUMN doc_references JSONB DEFAULT NULL;

COMMENT ON COLUMN strudel_messages.strudel_references IS 'Array of {id, title, author_name, url} for strudels used as RAG context';
COMMENT ON COLUMN strudel_messages.doc_references IS 'Array of {page_name, section_title, url} for docs used as RAG context';
