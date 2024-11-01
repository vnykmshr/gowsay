package api

import (
	"reflect"
	"strings"
	"testing"
)

func Test_newFace(t *testing.T) {
	f1 := newFace("apt", "tired")
	type args struct {
		cow  string
		mood string
	}
	tests := []struct {
		name string
		args args
		want *face
	}{
		{
			name: "t1",
			args: args{
				cow:  "apt",
				mood: "tired",
			},
			want: f1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newFace(tt.args.cow, tt.args.mood); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newFace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getGowsay(t *testing.T) {
	type args struct {
		action  string
		cow     string
		mood    string
		columns int32
		text    []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "t1",
			args: args{
				action: "help",
				text:   []string{"hewllo woold!"},
			},
			want: "```  _______________ < hewllo woold! >  ---------------  ``` ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := strings.Replace(getGowsay(tt.args.action, tt.args.cow, tt.args.mood, tt.args.columns, tt.args.text), "\n", " ", -1); got != tt.want {
				t.Errorf("getGowsay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_constructCow(t *testing.T) {
	f1 := newFace("apt", "wired")
	f2 := newFace("apt", "tired")
	type args struct {
		f *face
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "t1",
			args: args{
				f: f1,
			},
			want: `        (__)
         (OO)
   /------\/
  / |    ||
 *  /\---/\
    ~~   ~~
`,
		},
		{
			name: "t2",
			args: args{
				f: f2,
			},
			want: `        (__)
         (--)
   /------\/
  / |    ||
 *  /\---/\
    ~~   ~~
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := constructCow(tt.args.f); got != tt.want {
				t.Errorf("constructCow() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func Test_readInput(t *testing.T) {
	type args struct {
		args    []string
		columns int32
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "t1",
			args: args{
				args:    strings.Split("all you see", " "),
				columns: 5,
			},
			want: []string{"all", "you", "see"},
		},
		{
			name: "t2",
			args: args{
				args:    strings.Split("all\tyou see\tand", " "),
				columns: 5,
			},
			want: []string{"all", "you", "see", "and"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readInput(tt.args.args, tt.args.columns); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setPadding(t *testing.T) {
	type args struct {
		msgs  []string
		width int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "t1",
			args: args{
				msgs: []string{
					"abcde",
				},
				width: 10,
			},
			want: []string{
				"abcde     ",
			},
		},
		{
			name: "t2",
			args: args{
				msgs: []string{
					"abcde",
					"abcd",
					"abc",
				},
				width: 10,
			},
			want: []string{
				"abcde     ",
				"abcd      ",
				"abc       ",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := setPadding(tt.args.msgs, tt.args.width); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setPadding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_constructBallon(t *testing.T) {
	type args struct {
		f      *face
		action string
		msgs   []string
		width  int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := constructBallon(tt.args.f, tt.args.action, tt.args.msgs, tt.args.width); got != tt.want {
				t.Errorf("constructBallon() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_maxWidth(t *testing.T) {
	type args struct {
		msgs []string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "t1",
			args: args{},
			want: -1,
		},
		{
			name: "t2",
			args: args{
				msgs: []string{
					"all you see",
				},
			},
			want: 11,
		},
		{
			name: "t3",
			args: args{
				msgs: []string{
					"all you see",
					" and ",
					"all you touch is all you ever be",
				},
			},
			want: 32,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxWidth(tt.args.msgs); got != tt.want {
				t.Errorf("maxWidth() = %v, want %v", got, tt.want)
			}
		})
	}
}
