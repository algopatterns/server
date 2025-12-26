# Algorave CLI Tools

This document describes the command-line interface tools for Algorave.

## Quick Start

### Build All Binaries
```bash
make build
```

This creates:
- `bin/algorave` - Local interactive CLI
- `bin/algorave-ssh` - SSH server for remote access
- `bin/server` - HTTP API server
- `bin/ingester` - Documentation ingester

## Local CLI (`bin/algorave`)

Interactive terminal interface for local development.

### Usage
```bash
# Development mode (all commands available)
ALGORAVE_ENV=development ./bin/algorave

# Production mode (ingester hidden)
ALGORAVE_ENV=production ./bin/algorave
```

### Commands
- `start` - Start the Algorave HTTP server
- `ingest` - Run documentation ingester (dev mode only)
- `editor` - Interactive code editor with AI assistance
- `quit` - Exit the CLI

### Editor Mode
Press `Ctrl+S` to send code to AI for assistance
Press `Ctrl+L` to clear the editor
Press `Ctrl+C` to exit editor mode

## SSH Server (`bin/algorave-ssh`)

Remote terminal access for collaborative coding sessions.

### Usage
```bash
# Start SSH server (port 2222 by default)
./bin/algorave-ssh

# Custom configuration
ALGORAVE_SSH_PORT=2222 \
ALGORAVE_SSH_HOST_KEY=.ssh/algorave_host_key \
ALGORAVE_SSH_MAX_CONNECTIONS=50 \
./bin/algorave-ssh
```

### Connect
```bash
# From another terminal or machine
ssh localhost -p 2222
```

### Features
- Guest access (no authentication required)
- Production mode enforced (ingester hidden)
- Each connection gets isolated session
- Same TUI interface as local CLI

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `ALGORAVE_SSH_PORT` | 2222 | SSH server port |
| `ALGORAVE_SSH_HOST_KEY` | .ssh/algorave_host_key | Path to SSH host key |
| `ALGORAVE_SSH_MAX_CONNECTIONS` | 50 | Max concurrent connections |
| `ALGORAVE_AGENT_ENDPOINT` | http://localhost:8080/api/v1/generate | Agent API endpoint |

## Architecture

### Local CLI Flow
```
User → algorave binary → TUI (Bubbletea) → Agent API
                       → Server (background)
                       → Ingester (dev mode)
```

### SSH Server Flow
```
Remote User → SSH (port 2222) → algorave-ssh → TUI (per-session)
                                              → Agent API (shared)
```

### Production vs Development Mode

**Development Mode:**
- All commands available
- Can run ingester
- Full access to all features

**Production Mode (SSH default):**
- Ingester hidden
- Safe for public/guest access
- Users can still use editor and start server

## Makefile Commands

```bash
make build      # Build all binaries
make cli        # Build local CLI only
make ssh        # Build SSH server only
make clean      # Remove all binaries
make help       # Show all available commands
```

## File Locations

```
algorave/
├── bin/                    # All compiled binaries (gitignored)
│   ├── algorave            # Local CLI
│   ├── algorave-ssh        # SSH server
│   ├── server              # HTTP server
│   └── ingester            # Documentation ingester
│
├── cmd/
│   ├── algorave/           # Local CLI source
│   ├── algorave-ssh/       # SSH server source
│   ├── server/             # HTTP server source
│   └── ingester/           # Ingester source
│
└── internal/
    ├── tui/                # Terminal UI components
    └── ssh/                # SSH server wrapper
```

## Examples

### Local Development Workflow
```bash
# Build the CLI
make cli

# Run in dev mode
ALGORAVE_ENV=development ./bin/algorave

# At the prompt:
> ingest          # Ingest documentation
> start           # Start HTTP server
> editor          # Open code editor
> quit            # Exit
```

### Remote SSH Session
```bash
# Terminal 1: Start SSH server
./bin/algorave-ssh

# Terminal 2: Connect remotely
ssh localhost -p 2222

# In SSH session:
> editor          # Write code with AI
> quit            # Disconnect
```

### Using the Editor
```bash
./bin/algorave

> editor

# In editor:
# Type your musical idea:
"make a drum beat with reverb"

# Press Ctrl+S to send to AI
# AI generates Strudel code
# Continue editing or ask questions

# Press Ctrl+C to return to menu
```

## Troubleshooting

### SSH Server Won't Start
```bash
# Generate SSH host key
mkdir -p .ssh
ssh-keygen -t ed25519 -f .ssh/algorave_host_key -N ""

# Start server
./bin/algorave-ssh
```

### Agent Connection Errors
```bash
# Make sure HTTP server is running
./bin/server

# Or configure custom endpoint
ALGORAVE_AGENT_ENDPOINT=http://your-server:8080/api/v1/generate ./bin/algorave
```

### Port Already in Use
```bash
# Change SSH port
ALGORAVE_SSH_PORT=2223 ./bin/algorave-ssh
```

## Next Steps

See [CLI_ARCHITECTURE.md](docs/system-specs/CLI_ARCHITECTURE.md) for detailed technical architecture.

See [CLI_IMPLEMENTATION_PLAN.md](docs/CLI_IMPLEMENTATION_PLAN.md) for full implementation details.
