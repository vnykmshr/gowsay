package api

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestGetUsageString(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "t1",
			want: fmt.Sprintf("Usage: `%s [%s|%s] [cow] [mood] message`\n", CommandMoo, ActionThink, ActionSurprise),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUsageString(); got != tt.want {
				t.Errorf("GetUsageString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHelpString(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "t1",
			want: strings.Join([]string{
				"Cows: `apt`, `beavis.zen`, `bong`, `bud-frogs`, `bunny`, `calvin`, `cheese`, `cock`, `cower`,",
				"`daemon`, `default`, `dragon-and-cow`, `dragon`, `duck`, `elephant-in-snake`, `elephant`, `eyes`, `flaming-sheep`, `ghostbusters`,",
				"`gnu`, `hellokitty`, `kitty`, `koala`, `kosh`, `luke-koala`, `mech-and-cow`, `meow`, `milk`, `moofasa`,",
				"`moose`, `mutilated`, `pony-smaller`, `pony`, `random`, `ren`, `sheep`, `skeleton`, `snowman`, `stegosaurus`,",
				"`stimpy`, `suse`, `three-eyes`, `turkey`, `turtle`, `tux`, `unipony-smaller`, `unipony`, `vader-koala`, `vader`,",
				"`www`",
				"",
				"Moods: `borg`, `dead`, `greedy`, `paranoid`, `random`, `stoned`, `wired`, `young`",
				"",
			}, "\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetHelpString(); got != tt.want {
				t.Errorf("GetHelpString() = \n%v\nwant \n%v", got, tt.want)
			}
		})
	}
}

func Test_wrapString(t *testing.T) {
	type args struct {
		s     string
		limit int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "t1",
			args: args{
				s:     "        ",
				limit: 3,
			},
			want: "",
		},
		{
			name: "t2",
			args: args{
				s:     "All you see and all you touch is all you ever be",
				limit: 3,
			},
			want: strings.Join([]string{
				"All you see",
				"and all you",
				"touch is all",
				"you ever be",
				"",
			}, "\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := wrapString(tt.args.s, tt.args.limit); got != tt.want {
				t.Errorf("wrapString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getCows(t *testing.T) {
	type args struct {
		sorted bool
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "t1",
			args: args{
				sorted: true,
			},
			want: []string{
				"`apt`", "`beavis.zen`", "`bong`", "`bud-frogs`", "`bunny`", "`calvin`", "`cheese`", "`cock`", "`cower`", "`daemon`", "`default`", "`dragon-and-cow`", "`dragon`", "`duck`", "`elephant-in-snake`", "`elephant`", "`eyes`", "`flaming-sheep`", "`ghostbusters`", "`gnu`", "`hellokitty`", "`kitty`", "`koala`", "`kosh`", "`luke-koala`", "`mech-and-cow`", "`meow`", "`milk`", "`moofasa`", "`moose`", "`mutilated`", "`pony-smaller`", "`pony`", "`random`", "`ren`", "`sheep`", "`skeleton`", "`snowman`", "`stegosaurus`", "`stimpy`", "`suse`", "`three-eyes`", "`turkey`", "`turtle`", "`tux`", "`unipony-smaller`", "`unipony`", "`vader-koala`", "`vader`", "`www`",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCows(tt.args.sorted); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCows() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getMoods(t *testing.T) {
	type args struct {
		sorted bool
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "t1",
			args: args{
				sorted: true,
			},
			want: []string{
				"`borg`", "`dead`", "`greedy`", "`paranoid`", "`random`", "`stoned`", "`wired`", "`young`",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMoods(tt.args.sorted); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getMoods() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBanner(t *testing.T) {
	v := "test"
	type args struct {
		version string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "t1",
			args: args{
				version: v,
			},
			want: strings.Join([]string{
				"gowsay [" + v + "][" + os.Getenv(FieldEnv) + "]",
				fmt.Sprintf("Usage: `%s [%s|%s] [cow] [mood] message`", CommandMoo, ActionThink, ActionSurprise),
				"",
				"Cows: `apt`, `beavis.zen`, `bong`, `bud-frogs`, `bunny`, `calvin`, `cheese`, `cock`, `cower`,",
				"`daemon`, `default`, `dragon-and-cow`, `dragon`, `duck`, `elephant-in-snake`, `elephant`, `eyes`, `flaming-sheep`, `ghostbusters`,",
				"`gnu`, `hellokitty`, `kitty`, `koala`, `kosh`, `luke-koala`, `mech-and-cow`, `meow`, `milk`, `moofasa`,",
				"`moose`, `mutilated`, `pony-smaller`, `pony`, `random`, `ren`, `sheep`, `skeleton`, `snowman`, `stegosaurus`,",
				"`stimpy`, `suse`, `three-eyes`, `turkey`, `turtle`, `tux`, `unipony-smaller`, `unipony`, `vader-koala`, `vader`,",
				"`www`",
				"",
				"Moods: `borg`, `dead`, `greedy`, `paranoid`, `random`, `stoned`, `wired`, `young`",
				"",
			}, "\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetBanner(tt.args.version); got != tt.want {
				t.Errorf("GetBanner() = %v, want %v", got, tt.want)
			}
		})
	}
}
