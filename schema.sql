-- Algorave RAG System - Supabase Schema Setup
-- Run this in your Supabase SQL Editor

-- Enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- Create document_chunks table
CREATE TABLE IF NOT EXISTS document_chunks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    page_name TEXT NOT NULL,
    page_url TEXT NOT NULL,
    section_title TEXT,
    content TEXT NOT NULL,
    embedding vector(1536),  -- OpenAI text-embedding-3-small dimension
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create vector similarity search index
-- ivfflat is good for datasets < 1M vectors
CREATE INDEX IF NOT EXISTS document_chunks_embedding_idx 
ON document_chunks 
USING ivfflat (embedding vector_cosine_ops)
WITH (lists = 100);

-- Create index on page_name for faster lookups
CREATE INDEX IF NOT EXISTS document_chunks_page_name_idx 
ON document_chunks(page_name);

-- Optional: Create index on created_at for maintenance queries
CREATE INDEX IF NOT EXISTS document_chunks_created_at_idx 
ON document_chunks(created_at DESC);

-- Verify the table was created
SELECT 
    tablename, 
    schemaname 
FROM pg_tables 
WHERE tablename = 'document_chunks';

-- Check if vector extension is enabled
SELECT * FROM pg_extension WHERE extname = 'vector';