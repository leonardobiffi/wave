package radiogarden

import "testing"

func TestExtractID(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "place",
			args: args{
				url: "/visit/dourados-ms/8CKRVU3Z",
			},
			want: "8CKRVU3Z",
		},
		{
			name: "channel",
			args: args{
				url: "/listen/web-radio-jave/1z3W71T2",
			},
			want: "1z3W71T2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractID(tt.args.url); got != tt.want {
				t.Errorf("ExtractID() = %v, want %v", got, tt.want)
			}
		})
	}
}
