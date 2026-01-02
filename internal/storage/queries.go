package storage

const (
	getChunkCountQuery   = "SELECT COUNT(*) FROM doc_embeddings"
	deleteAllChunksQuery = "DELETE FROM doc_embeddings"

	insertChunkQuery = `
		INSERT INTO doc_embeddings (page_name, page_url, section_title, content, embedding, metadata)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
)
