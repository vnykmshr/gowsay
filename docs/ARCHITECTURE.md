# Architecture

## Overview

gowsay is a single Go binary with three interfaces (CLI, Web UI, API) sharing a common rendering engine.

```
┌─────────────┐
│   main.go   │  Entry point - routes to CLI or server
└──────┬──────┘
       │
   ┌───┴────┐
   │        │
   ▼        ▼
┌──────┐  ┌────────┐
│ CLI  │  │ Server │
└──┬───┘  └────┬───┘
   │           │
   │      ┌────┴─────┬─────────┐
   │      │          │         │
   │      ▼          ▼         ▼
   │   ┌────┐    ┌─────┐   ┌─────┐
   │   │API │    │Slack│   │Web  │
   │   │    │    │     │   │UI   │
   │   └─┬──┘    └──┬──┘   └──┬──┘
   │     │          │         │
   └─────┴──────────┴─────────┘
                │
                ▼
         ┌────────────┐
         │ cow/       │  Core rendering engine
         │ - render() │
         │ - cows     │
         │ - moods    │
         └────────────┘
```

## Package Responsibilities

### `main.go`
- Parses command-line flags
- Routes to CLI mode or server mode
- Version injection point

### `cow/`
Core rendering logic - **no external dependencies**
- `render.go` - Main `Render()` function, balloon building, text wrapping
- `cows.go` - 52 cow templates as embedded strings
- `moods.go` - 8 facial expression configurations
- `messages.go` - Random moo messages

### `api/`
HTTP server and handlers
- `init.go` - Module initialization, Slack handler
- `handlers.go` - REST API endpoints (moo, cows, moods, health)
- `middleware.go` - CORS handling
- `types.go` - Request/response structs
- `common.go` - Shared utilities
- `help.go` - Banner and usage text generation

### `web/`
Static assets embedded at compile time
- `embed.go` - Go embed directives
- HTML/CSS/JS served at `/` route

## Request Flow

### CLI Mode
```
User input → flag parse → cow.Render() → stdout
```

### API Mode
```
HTTP request → handlers.APIMoo →
  validate params →
  cow.Render() →
  JSON response
```

### Web UI Mode
```
Browser → ServeWeb() → static HTML/CSS/JS →
  User interaction → fetch /api/moo →
  Display result
```

## Data Flow

### Rendering Pipeline
```
1. Text input (string or []string)
2. Select cow template (by name or random)
3. Apply mood (changes eyes/tongue)
4. Choose action (say vs think - changes bubble connectors)
5. Wrap text to column width
6. Build balloon (border + wrapped text)
7. Substitute placeholders in cow template
8. Return ASCII art string
```

### Template Variables
- `{{.Thoughts}}` - Speech/thought bubble connector (\ or o)
- `{{.Eyes}}` - Cow eyes (oo, xx, $$, etc.)
- `{{.Tongue}}` - Cow tongue (usually empty or "U ")

## Configuration

All config via environment variables:
- `PORT` - HTTP server port (default: 9000)
- `GOWSAY_TOKEN` - Auth token for /say endpoint (default: "devel")
- `GOWSAY_COLUMNS` - Text wrapping width (default: 40)

## Deployment

### Single Binary
- Statically linked Go binary (~11MB with embedded assets)
- No runtime dependencies
- Works on Linux, macOS, Windows

### Docker
- Multi-stage build (Go 1.23 → scratch)
- Final image ~11MB
- Default: runs server on port 9000
- Override: `docker run gowsay --help` for CLI mode

## Testing Strategy

- **Unit tests**: Core rendering logic in `cow/`
- **Handler tests**: HTTP endpoint behavior in `api/`
- **Integration tests**: Full server lifecycle
- **Edge cases**: Unicode, large inputs, concurrent requests
- Coverage: 89.7% API, 97.6% cow rendering

## Design Decisions

### Why Embedded Assets?
Single binary deployment, no external files to manage.

### Why No Framework?
Stdlib sufficient for simple routing, reduces dependencies.

### Why Maps for Cows/Moods?
Fast lookup, easy to extend, no file I/O at runtime.

### Why Separate Packages?
- `cow/` is reusable, has no HTTP knowledge
- `api/` depends on `cow/`, not vice versa
- `main.go` orchestrates, doesn't contain logic

## Adding Features

### New Cow
Add to `cow/cows.go` - update `cowNames` slice and `cows` map.

### New Endpoint
Add handler to `api/handlers.go`, register route in `main.go`.

### New Mood
Add to `cow/moods.go` - update `moods` map and `moodNames` slice.

## Performance

- Typical render: <1ms
- API response: <10ms (includes network)
- Memory: ~10MB resident
- Concurrency: Stateless handlers, safe for concurrent use
