## gowsay
Implementation of cowsay in go

Upstream for custom slack command `/moo` deployed at https://gowsay.vnykmshr.com/say

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

## Development
### Build
```
go build -ldflags "-X 'main.version=`git log -1 --pretty=format:"%h"`'" -v
```

### Linux
```
GOOS=linux GOARCH=amd64
```

### Pi
```
GOOS=linux GOARCH=arm GOARM=5
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
