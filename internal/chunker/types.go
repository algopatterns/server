package chunker

type Chunk struct {
	PageName     string
	PageURL      string
	SectionTitle string
	Content      string
	Metadata     map[string]interface{}
}

type ChunkOptions struct {
	MaxTokens       int
	OverlapTokens   int
	PreserveHeaders bool
}

type Section struct {
	Title   string
	Level   int
	Content string
}
