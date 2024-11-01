package api

import (
	"errors"
	"testing"
)

func Test_wrap(t *testing.T) {
	type args struct {
		msg string
		err error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "t1",
			args: args{
				msg: "failed to process",
				err: errors.New("test"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := wrap(tt.args.msg, tt.args.err); (err != nil) != tt.wantErr {
				t.Errorf("wrap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
