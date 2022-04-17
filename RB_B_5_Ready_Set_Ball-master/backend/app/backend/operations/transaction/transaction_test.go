package transaction

import (
	"reflect"
	"testing"
)

func TestReadableMap(t *testing.T) {
	type args struct {
		m map[string]interface{}
	}
	tests := []struct {
		name   string
		args   args
		want   map[string]string
		errMsg string
	}{
		//Test cases.
		{
			name:   "map[string]string is created",
			args:   args{map[string]interface{}{"one": "string", "two": "string"}},
			want:   map[string]string{"one": "string", "two": "string"},
			errMsg: "", //empty string denotes I'm not expecting an error
		},
		{
			name:   "fails when mixed with bools (non-string)",
			args:   args{map[string]interface{}{"one": "string", "two": false}},
			want:   map[string]string{"one": "string", "two": "false"},
			errMsg: "", //empty string denotes I'm not expecting an error
		},
		{
			name:   "fails when mixed with ints (non-string)",
			args:   args{map[string]interface{}{"one": "string", "two": 2}},
			want:   nil,
			errMsg: "Expected param type to be string/bool but recieved map[two]=2 as int",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadableMap(tt.args.m)
			if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("ReadableMap() error = %v, Expected error %s", err, tt.errMsg)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadableMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
