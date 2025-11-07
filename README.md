## gowsay

Implementation of cowsay in Go

**Features:**
- 🖥️ Command-line tool (like original cowsay)
- 🌐 HTTP API server for Slack integration
- 🐮 41 different cows
- 😀 7 moods (borg, dead, greedy, paranoid, stoned, wired, young)

**Status:** Undergoing modernization for gowsay 2.0 - CLI tool ✅, web interface (coming soon)

## Usage

### CLI Tool

```bash
# Basic usage
gowsay "Hello World"

# Use a different cow
gowsay -c dragon "Fire!"

# Make the cow think instead of speak
gowsay -t "Hmm..."

# Random cow and mood
gowsay -r "Surprise!"

# Use mood
gowsay -c tux -m dead "System crashed"

# From pipe
echo "Hello from pipe" | gowsay

# List available cows and moods
gowsay -l

# Help
gowsay --help
```

### HTTP API

Start the server:
```bash
./bin/gowsay serve
# or
PORT=8080 ./bin/gowsay serve
```

**Endpoints:**

```bash
# Generate cowsay (query params)
curl 'http://localhost:9000/api/moo?text=Hello&cow=dragon&action=say'

# Generate cowsay (JSON)
curl -X POST http://localhost:9000/api/moo \
  -H 'Content-Type: application/json' \
  -d '{"text":"Hello","cow":"dragon","mood":"wired"}'

# List all cows
curl http://localhost:9000/api/cows

# List all moods
curl http://localhost:9000/api/moods

# Health check
curl http://localhost:9000/health
```

**API Parameters:**
- `text` - Message to display (required)
- `cow` - Cow name (default: "default", or "random")
- `mood` - Mood name (optional, or "random")
- `action` - "say" or "think" (default: "say")
- `columns` - Text width for wrapping (default: 40)

### Slack Command

Deployed at https://gowsay.vnykmshr.com/say

```
/moo [think|surprise] [cow] [mood] message
```

### Cows
```
`apt`, `beavis.zen`, `bong`, `bud-frogs`, `bunny`, `calvin`, `cheese`, `cock`, `cower`,
`daemon`, `default`, `dragon-and-cow`, `dragon`, `duck`, `elephant-in-snake`, `elephant`, `eyes`, `flaming-sheep`, `ghostbusters`,
`gnu`, `hellokitty`, `kitty`, `koala`, `kosh`, `luke-koala`, `mech-and-cow`, `meow`, `milk`, `moofasa`,
`moose`, `mutilated`, `pony-smaller`, `pony`, `random`, `ren`, `sheep`, `skeleton`, `snowman`, `stegosaurus`,
`stimpy`, `suse`, `three-eyes`, `turkey`, `turtle`, `tux`, `unipony-smaller`, `unipony`, `vader-koala`, `vader`,
`www`
```
### Moods
```
`borg`, `dead`, `greedy`, `paranoid`, `random`, `stoned`, `wired`, `young`
```

## Configuration

Configuration via environment variables:

- `PORT` - Server port (default: `9000`)
- `GOWSAY_TOKEN` - Authentication token (default: `devel`)
- `GOWSAY_COLUMNS` - Text column width (default: `40`)

## Development

### Build
```bash
make build
# or
go build -ldflags "-X 'main.version=`git log -1 --pretty=format:"%h"`'" -v
```

### Run CLI
```bash
./bin/gowsay "Hello World"
```

### Run Server
```bash
make run/server
# or
./bin/gowsay serve
# or with custom port
PORT=8080 ./bin/gowsay serve
```

### Test
```bash
make test
```

### Cross-compile

Linux:
```bash
GOOS=linux GOARCH=amd64 go build -o bin/gowsay-linux
```

Raspberry Pi:
```bash
GOOS=linux GOARCH=arm GOARM=5 go build -o bin/gowsay-pi
```

### Slack commands
### Sample Request
```
curl --location --request POST 'https://gowsay.vnykmshr.com/say' \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data-urlencode 'token=xxx' \
  --data-urlencode 'text=Puff! Puff!!'
```

### Sample Response
```
{
  "response_type": "in_channel",
  "text": "```\n ________\n/ Puff!  \\\n\\ Puff!! /\n --------\n        \\   ^__^\n         \\  (oo)\\_______\n            (__)\\       )\\/\\\n                ||----w |\n                ||     ||\n\n```\n"
}
```
