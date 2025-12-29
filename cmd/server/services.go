package main

import (
	"context"
	"fmt"

	"github.com/algorave/server/internal/agent"
	"github.com/algorave/server/internal/config"
	"github.com/algorave/server/internal/llm"
	"github.com/algorave/server/internal/retriever"
	"github.com/algorave/server/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

// holds all external service clients (LLM, storage, retriever, agent)
type Services struct {
	Agent     *agent.Agent
	LLM       llm.LLM
	Retriever *retriever.Client
	Storage   *storage.Client
}

// creates and configures all service clients
func InitializeServices(cfg *config.Config, db *pgxpool.Pool) (*Services, error) {
	llmClient, err := llm.NewLLM(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to create LLM client: %w", err)
	}

	retrieverClient := retriever.New(db, llmClient)
	storageClient := &storage.Client{}
	agentClient := agent.New(retrieverClient, llmClient)

	return &Services{
		Agent:     agentClient,
		LLM:       llmClient,
		Retriever: retrieverClient,
		Storage:   storageClient,
	}, nil
}
