## gowsay

Implementation of cowsay in Go

Upstream for custom slack command `/moo` deployed at https://gowsay.vnykmshr.com/say

**Status:** Undergoing modernization for gowsay 2.0 - adding CLI tool and web interface while maintaining Slack compatibility.

### Usage
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

### Run
```bash
make run
# or
PORT=8080 ./bin/gowsay
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
