# Algorave

Create music using human language.

## Quick Start

### Prerequisites

- Go 1.21+
- Supabase account
- OpenAI API key
- Anthropic API key (for server phase)

### Setup

1. **Clone and enter the project:**

   ```bash
   cd algorave
   ```

2. **Install dependencies:**

   ```bash
   go mod download
   ```

3. **Set up Supabase:**

   - Run the SQL in `schema.sql` in your Supabase SQL Editor
   - Get your connection string from Supabase Dashboard → Settings → Database

4. **Configure environment:**

   ```bash
   cp .env.example .env
   # Edit .env with the actual API keys
   ```

### Running Ingestion

```bash
go run cmd/ingest/main.go --docs ./docs/project-docs --clear
```

Options:

- `--docs`: Path to documentation directory (default: `./docs/project-docs`)
- `--clear`: Clear existing embeddings before ingesting

### Development

See `agent.md` for full architectural details.

For detailed coding standards, see `.clinerules`.

## Architecture

```
User Query → Query Transformation → Vector Search → Retrieved Docs
                                                            ↓
                                    Code Generation ← Docs + History
```

See `agent.md` for complete architecture documentation.
