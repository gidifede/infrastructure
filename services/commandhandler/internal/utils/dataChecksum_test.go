package utils

import "testing"

func TestComputeChecksum(t *testing.T) {
	type args struct {
		data              []byte
		removeWhiteSpaces bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "No space, remove disabled", args: args{data: []byte("AAAA"), removeWhiteSpaces: false}, want: "098890dde069e9abad63f19a0d9e1f32"},
		{name: "No space, remove enabled", args: args{data: []byte("AAAA"), removeWhiteSpaces: true}, want: "098890dde069e9abad63f19a0d9e1f32"},
		{name: "Space, remove disabled", args: args{data: []byte(`AAAA    HELLO 
		 bye bye`), removeWhiteSpaces: false}, want: "3e4a17e61cbb9e92b5ba5c3d60592541"},
		{name: "Space, remove enabled", args: args{data: []byte(`AAAA    HELLO 
		 bye bye`), removeWhiteSpaces: true}, want: "575634a18a0790c99d9aaba5fe93b71a"},
		{name: "Space, remove enabled", args: args{data: []byte(`AAAAHELLObyebye`), removeWhiteSpaces: true}, want: "575634a18a0790c99d9aaba5fe93b71a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ComputeChecksum(tt.args.data, tt.args.removeWhiteSpaces); got != tt.want {
				t.Errorf("ComputeChecksum() = %v, want %v", got, tt.want)
			}
		})
	}
}
